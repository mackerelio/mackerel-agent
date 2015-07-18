package command

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/Songmu/retry"
	"github.com/mackerelio/mackerel-agent/agent"
	"github.com/mackerelio/mackerel-agent/checks"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/spec"
	"github.com/mackerelio/mackerel-agent/util"
)

var logger = logging.GetLogger("command")
var metricsInterval = 60 * time.Second

const idFileName = "id"

func idFilePath(root string) string {
	return filepath.Join(root, idFileName)
}

func loadHostID(root string) (string, error) {
	content, err := ioutil.ReadFile(idFilePath(root))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func saveHostID(root string, id string) error {
	err := os.MkdirAll(root, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(idFilePath(root))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(id))
	if err != nil {
		return err
	}

	return nil
}

var retryNum uint = 20
var retryInterval = 3 * time.Second

// prepareHost collects specs of the host and sends them to Mackerel server.
// A unique host-id is returned by the server if one is not specified.
func prepareHost(root string, api *mackerel.API, roleFullnames []string, checks []string, displayName string, hostSt string) (*mackerel.Host, error) {
	// XXX this configuration should be moved to under spec/linux
	os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+os.Getenv("PATH"))
	os.Setenv("LANG", "C") // prevent changing outputs of some command, e.g. ifconfig.

	doRetry := func(f func() error) {
		retry.Retry(retryNum, retryInterval, f)
	}

	filterErrorForRetry := func(err error) error {
		if err != nil {
			logger.Warningf("%s", err.Error())
		}
		if apiErr, ok := err.(*mackerel.Error); ok && apiErr.IsClientError() {
			// don't retry when client error (APIKey error etc.) occurred
			return nil
		}
		return err
	}

	hostname, meta, interfaces, lastErr := collectHostSpecs()
	if lastErr != nil {
		return nil, fmt.Errorf("error while collecting host specs: %s", lastErr.Error())
	}

	var result *mackerel.Host
	if hostID, err := loadHostID(root); err != nil { // create
		logger.Debugf("Registering new host on mackerel...")

		doRetry(func() error {
			hostID, lastErr = api.CreateHost(hostname, meta, interfaces, roleFullnames, displayName)
			return filterErrorForRetry(lastErr)
		})

		if lastErr != nil {
			return nil, fmt.Errorf("Failed to register this host: %s", lastErr.Error())
		}

		doRetry(func() error {
			result, lastErr = api.FindHost(hostID)
			return filterErrorForRetry(lastErr)
		})
		if lastErr != nil {
			return nil, fmt.Errorf("Failed to find this host on mackerel: %s", lastErr.Error())
		}
	} else { // check the hostID is valid or not
		doRetry(func() error {
			result, lastErr = api.FindHost(hostID)
			return filterErrorForRetry(lastErr)
		})
		if lastErr != nil {
			return nil, fmt.Errorf("Failed to find this host on mackerel (You may want to delete file \"%s\" to register this host to an another organization): %s", idFilePath(root), lastErr.Error())
		}
	}

	if hostSt != "" && hostSt != result.Status {
		doRetry(func() error {
			lastErr = api.UpdateHostStatus(result.ID, hostSt)
			return filterErrorForRetry(lastErr)
		})
		if lastErr != nil {
			return nil, fmt.Errorf("Failed to set default host status: %s, %s", hostSt, lastErr.Error())
		}
	}

	lastErr = saveHostID(root, result.ID)
	if lastErr != nil {
		return nil, fmt.Errorf("Failed to save host ID: %s", lastErr.Error())
	}

	return result, nil
}

// Interval between each updating host specs.
var specsUpdateInterval = 1 * time.Hour

func delayByHost(host *mackerel.Host) int {
	s := sha1.Sum([]byte(host.ID))
	return int(s[len(s)-1]) % int(config.PostMetricsInterval.Seconds())
}

type Context struct {
	ag   *agent.Agent
	Conf *config.Config
	Host *mackerel.Host
	API  *mackerel.API
}

type postValue struct {
	values   []*mackerel.CreatingMetricsValue
	retryCnt int
}

func newPostValue(values []*mackerel.CreatingMetricsValue) *postValue {
	return &postValue{values, 0}
}

type loopState uint8

const (
	loopStateFirst loopState = iota
	loopStateDefault
	loopStateQueued
	loopStateHadError
	loopStateTerminating
)

func loop(c *Context, termCh chan struct{}) int {
	quit := make(chan struct{})
	defer close(quit) // broadcast terminating

	// Periodically update host specs.
	go updateHostSpecsLoop(c, quit)

	postQueue := make(chan *postValue, c.Conf.Connection.PostMetricsBufferSize)
	go enqueueLoop(c, postQueue, quit)

	postDelaySeconds := delayByHost(c.Host)
	initialDelay := postDelaySeconds / 2
	logger.Debugf("wait %d seconds before initial posting.", initialDelay)
	select {
	case <-termCh:
		return 0
	case <-time.After(time.Duration(initialDelay) * time.Second):
		c.ag.InitPluginGenerators(c.API)
	}

	termCheckerCh := make(chan struct{})
	termMetricsCh := make(chan struct{})

	// fan-out termCh
	go func() {
		for _ = range termCh {
			termCheckerCh <- struct{}{}
			termMetricsCh <- struct{}{}
		}
	}()

	runCheckersLoop(c, termCheckerCh, quit)

	lState := loopStateFirst
	for {
		select {
		case <-termMetricsCh:
			if lState == loopStateTerminating {
				return 1
			}
			lState = loopStateTerminating
			if len(postQueue) <= 0 {
				return 0
			}
		case v := <-postQueue:
			origPostValues := [](*postValue){v}
			if len(postQueue) > 0 {
				// Bulk posting. However at most "two" metrics are to be posted, so postQueue isn't always empty yet.
				logger.Debugf("Merging datapoints with next queued ones")
				nextValues := <-postQueue
				origPostValues = append(origPostValues, nextValues)
			}

			delaySeconds := 0
			switch lState {
			case loopStateFirst: // request immediately to create graph defs of host
				// nop
			case loopStateQueued:
				delaySeconds = c.Conf.Connection.PostMetricsDequeueDelaySeconds
			case loopStateHadError:
				// TODO: better interval calculation. exponential backoff or so.
				delaySeconds = c.Conf.Connection.PostMetricsRetryDelaySeconds
			case loopStateTerminating:
				// dequeue and post every one second when terminating.
				delaySeconds = 1
			default:
				// Sending data at every 0 second from all hosts causes request flooding.
				// To prevent flooding, this loop sleeps for some seconds
				// which is specific to the ID of the host running agent on.
				// The sleep second is up to 60s (to be exact up to `config.Postmetricsinterval.Seconds()`.
				elapsedSeconds := int(time.Now().Unix() % int64(config.PostMetricsInterval.Seconds()))
				if postDelaySeconds > elapsedSeconds {
					delaySeconds = postDelaySeconds - elapsedSeconds
				}
			}

			// determin next loopState before sleeping
			if lState != loopStateTerminating {
				if len(postQueue) > 0 {
					lState = loopStateQueued
				} else {
					lState = loopStateDefault
				}
			}

			logger.Debugf("Sleep %d seconds before posting.", delaySeconds)
			select {
			case <-time.After(time.Duration(delaySeconds) * time.Second):
				// nop
			case <-termMetricsCh:
				if lState == loopStateTerminating {
					return 1
				}
				lState = loopStateTerminating
			}

			postValues := [](*mackerel.CreatingMetricsValue){}
			for _, v := range origPostValues {
				postValues = append(postValues, v.values...)
			}
			err := c.API.PostMetricsValues(postValues)
			if err != nil {
				logger.Errorf("Failed to post metrics value (will retry): %s", err.Error())
				if lState != loopStateTerminating {
					lState = loopStateHadError
				}
				go func() {
					for _, v := range origPostValues {
						v.retryCnt++
						// It is difficult to distinguish the error is server error or data error.
						// So, if retryCnt exceeded the configured limit, postValue is considered invalid and abandoned.
						if v.retryCnt > c.Conf.Connection.PostMetricsRetryMax {
							json, err := json.Marshal(v.values)
							if err != nil {
								logger.Errorf("Something wrong with post values. marshaling failed.")
							} else {
								logger.Errorf("Post values may be invalid and abandoned: %s", string(json))
							}
							continue
						}
						postQueue <- v
					}
				}()
				continue
			}
			logger.Debugf("Posting metrics succeeded.")

			if lState == loopStateTerminating && len(postQueue) <= 0 {
				return 0
			}
		}
	}
}

func updateHostSpecsLoop(c *Context, quit chan struct{}) {
	for {
		UpdateHostSpecs(c.Conf, c.API, c.Host)
		select {
		case <-quit:
			return
		case <-time.After(specsUpdateInterval):
			// nop
		}
	}
}

func enqueueLoop(c *Context, postQueue chan *postValue, quit chan struct{}) {
	metricsResult := c.ag.Watch()
	for {
		select {
		case <-quit:
			return
		case result := <-metricsResult:
			created := float64(result.Created.Unix())
			creatingValues := [](*mackerel.CreatingMetricsValue){}
			for name, value := range (map[string]float64)(result.Values) {
				if math.IsNaN(value) || math.IsInf(value, 0) {
					logger.Warningf("Invalid value: hostID = %s, name = %s, value = %f\n is not sent.", c.Host.ID, name, value)
					continue
				}

				creatingValues = append(
					creatingValues,
					&mackerel.CreatingMetricsValue{
						HostID: c.Host.ID,
						Name:   name,
						Time:   created,
						Value:  value,
					},
				)
			}
			logger.Debugf("Enqueuing task to post metrics.")
			postQueue <- newPostValue(creatingValues)
		}
	}
}

// runCheckersLoop generates "checker" goroutines
// which run for each checker commands and one for HTTP POSTing
// the reports to Mackerel API.
func runCheckersLoop(c *Context, termCheckerCh <-chan struct{}, quit <-chan struct{}) {
	var (
		checkReportCh          chan *checks.Report
		reportCheckImmediateCh chan struct{}
	)
	for _, checker := range c.ag.Checkers {
		if checkReportCh == nil {
			checkReportCh = make(chan *checks.Report)
			reportCheckImmediateCh = make(chan struct{})
		}

		go func(checker checks.Checker) {
			var (
				lastStatus  = checks.StatusUndefined
				lastMessage = ""
			)

			util.Periodically(
				func() {
					report, err := checker.Check()
					if err != nil {
						logger.Errorf("checker %v: %s", checker, err)
						return
					}

					logger.Debugf("checker %q: report=%v", checker.Name, report)

					if report.Status == lastStatus && report.Message == lastMessage {
						// Do not report if nothing has changed
						return
					}

					checkReportCh <- report

					// If status has changed, send it immediately
					// but if the status was OK and it's first invocation of a check, do not
					if report.Status != lastStatus && !(report.Status == checks.StatusOK && lastStatus == checks.StatusUndefined) {
						logger.Debugf("checker %q: status has changed %v -> %v: send it immediately", checker.Name, lastStatus, report.Status)
						reportCheckImmediateCh <- struct{}{}
					}

					lastStatus = report.Status
					lastMessage = report.Message
				},
				checker.Interval(),
				quit,
			)
		}(checker)
	}
	if checkReportCh != nil {
		go func() {
			exit := false
			for !exit {
				select {
				case <-time.After(1 * time.Minute):
				case <-termCheckerCh:
					logger.Debugf("received 'term' chan")
					exit = true
				case <-reportCheckImmediateCh:
					logger.Debugf("received 'immediate' chan")
				}

				reports := []*checks.Report{}
			DrainCheckReport:
				for {
					select {
					case report := <-checkReportCh:
						reports = append(reports, report)
					default:
						break DrainCheckReport
					}
				}

				for i, report := range reports {
					logger.Debugf("reports[%d]: %#v", i, report)
				}

				if len(reports) == 0 {
					continue
				}

				err := c.API.ReportCheckMonitors(c.Host.ID, reports)
				if err != nil {
					logger.Errorf("ReportCheckMonitors: %s", err)

					// queue back the reports
					go func() {
						for _, report := range reports {
							logger.Debugf("queue back report: %#v", report)
							checkReportCh <- report
						}
					}()
				}
			}
		}()
	} else {
		// consume termCheckerCh
		go func() {
			for _ = range termCheckerCh {
			}
		}()
	}
}

// collectHostSpecs collects host specs (correspond to "name", "meta" and "interfaces" fields in API v0)
func collectHostSpecs() (string, map[string]interface{}, []map[string]interface{}, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to obtain hostname: %s", err.Error())
	}

	specGens := specGenerators()
	cGen := spec.SuggestCloudGenerator()
	if cGen != nil {
		specGens = append(specGens, cGen)
	}
	meta := spec.Collect(specGens)

	interfacesSpec, err := interfaceGenerator().Generate()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to collect interfaces: %s", err.Error())
	}

	interfaces, _ := interfacesSpec.([]map[string]interface{})

	return hostname, meta, interfaces, nil
}

// UpdateHostSpecs updates the host information that is already registered on Mackerel.
func UpdateHostSpecs(conf *config.Config, api *mackerel.API, host *mackerel.Host) {
	logger.Debugf("Updating host specs...")

	hostname, meta, interfaces, err := collectHostSpecs()
	if err != nil {
		logger.Errorf("While collecting host specs: %s", err)
		return
	}

	err = api.UpdateHost(host.ID, mackerel.HostSpec{
		Name:          hostname,
		Meta:          meta,
		Interfaces:    interfaces,
		RoleFullnames: conf.Roles,
		Checks:        conf.CheckNames(),
		DisplayName:   conf.DisplayName,
	})

	if err != nil {
		logger.Errorf("Error while updating host specs: %s", err)
	} else {
		logger.Debugf("Host specs sent.")
	}
}

// Prepare sets up API and registers the host data to the Mackerel server.
// Use returned values to call Run().
func Prepare(conf *config.Config) (*Context, error) {
	api, err := mackerel.NewAPI(conf.Apibase, conf.Apikey, conf.Verbose)
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare an api: %s", err.Error())
	}

	host, err := prepareHost(conf.Root, api, conf.Roles, conf.CheckNames(), conf.DisplayName, conf.HostStatus.OnStart)
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare host: %s", err.Error())
	}

	return &Context{
		ag:   NewAgent(conf),
		Conf: conf,
		Host: host,
		API:  api,
	}, nil
}

// RunOnce collects specs and metrics, then output them to stdout.
func RunOnce(conf *config.Config) error {
	graphdefs, hostSpec, metrics, err := runOncePayload(conf)
	if err != nil {
		return err
	}

	json, err := json.Marshal(map[string]interface{}{
		"host":    hostSpec,
		"metrics": metrics,
	})
	if err != nil {
		logger.Warningf("Error while marshaling graphdefs: err = %s, graphdefs = %s.", err.Error(), graphdefs)
		return err
	}
	fmt.Println(string(json))
	return nil
}

func runOncePayload(conf *config.Config) ([]mackerel.CreateGraphDefsPayload, *mackerel.HostSpec, *agent.MetricsResult, error) {
	hostname, meta, interfaces, err := collectHostSpecs()
	if err != nil {
		logger.Errorf("While collecting host specs: %s", err)
		return nil, nil, nil, err
	}

	origInterval := metricsInterval
	metricsInterval = 1 * time.Second
	defer func() {
		metricsInterval = origInterval
	}()
	ag := NewAgent(conf)
	graphdefs := ag.CollectGraphDefsOfPlugins()
	metrics := ag.CollectMetrics(time.Now())
	return graphdefs, &mackerel.HostSpec{
		Name:          hostname,
		Meta:          meta,
		Interfaces:    interfaces,
		RoleFullnames: conf.Roles,
		Checks:        conf.CheckNames(),
		DisplayName:   conf.DisplayName,
	}, metrics, nil
}

// NewAgent creates a new instance of agent.Agent from its configuration conf.
func NewAgent(conf *config.Config) *agent.Agent {
	return &agent.Agent{
		MetricsGenerators: prepareGenerators(conf),
		PluginGenerators:  pluginGenerators(conf),
		Checkers:          createCheckers(conf),
	}
}

// Run starts the main metric collecting logic and this function will never return.
func Run(c *Context, termCh chan struct{}) int {
	logger.Infof("Start: apibase = %s, hostName = %s, hostID = %s", c.Conf.Apibase, c.Host.Name, c.Host.ID)

	exitCode := loop(c, termCh)
	if exitCode == 0 && c.Conf.HostStatus.OnStop != "" {
		// TOOD error handling. supoprt retire(?)
		err := c.API.UpdateHostStatus(c.Host.ID, c.Conf.HostStatus.OnStop)
		if err != nil {
			logger.Errorf("Failed update host status on stop: %s", err)
		}
	}
	return exitCode
}

func createCheckers(conf *config.Config) []checks.Checker {
	checkers := []checks.Checker{}

	for name, pluginConfig := range conf.Plugin["checks"] {
		checker := checks.Checker{
			Name:   name,
			Config: pluginConfig,
		}
		logger.Debugf("Checker created: %v", checker)
		checkers = append(checkers, checker)
	}

	return checkers
}

func prepareGenerators(conf *config.Config) []metrics.Generator {
	diagnostic := conf.Diagnostic
	generators := metricsGenerators(conf)
	if diagnostic {
		generators = append(generators, &metrics.AgentGenerator{})
	}
	return generators
}
