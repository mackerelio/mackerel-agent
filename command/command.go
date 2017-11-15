package command

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Songmu/retry"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/agent"
	"github.com/mackerelio/mackerel-agent/checks"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/spec"
)

var logger = logging.GetLogger("command")
var metricsInterval = 60 * time.Second

var retryNum uint = 20
var retryInterval = 3 * time.Second

// AgentMeta contains meta information about mackerel-agent
type AgentMeta struct {
	Version  string
	Revision string
}

// prepareHost collects specs of the host and sends them to Mackerel server.
// A unique host-id is returned by the server if one is not specified.
func prepareHost(conf *config.Config, api *mackerel.API) (*mackerel.Host, error) {
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

	hostname, meta, interfaces, customIdentifier, lastErr := collectHostSpecs()
	if lastErr != nil {
		return nil, fmt.Errorf("error while collecting host specs: %s", lastErr.Error())
	}

	var result *mackerel.Host
	if hostID, err := conf.LoadHostID(); err != nil { // create

		if customIdentifier != "" {
			retry.Retry(3, 2*time.Second, func() error {
				result, lastErr = api.FindHostByCustomIdentifier(customIdentifier)
				return filterErrorForRetry(lastErr)
			})
			if result != nil {
				hostID = result.ID
			}
		}

		if result == nil {
			logger.Debugf("Registering new host on mackerel...")

			doRetry(func() error {
				hostID, lastErr = api.CreateHost(mackerel.HostSpec{
					Name:             hostname,
					Meta:             meta,
					Interfaces:       interfaces,
					RoleFullnames:    conf.Roles,
					DisplayName:      conf.DisplayName,
					CustomIdentifier: customIdentifier,
				})
				return filterErrorForRetry(lastErr)
			})

			if lastErr != nil {
				return nil, fmt.Errorf("failed to register this host: %s", lastErr.Error())
			}

			doRetry(func() error {
				result, lastErr = api.FindHost(hostID)
				return filterErrorForRetry(lastErr)
			})
			if lastErr != nil {
				return nil, fmt.Errorf("failed to find this host on mackerel: %s", lastErr.Error())
			}
		}
	} else { // check the hostID is valid or not
		doRetry(func() error {
			result, lastErr = api.FindHost(hostID)
			return filterErrorForRetry(lastErr)
		})
		if lastErr != nil {
			if fsStorage, ok := conf.HostIDStorage.(*config.FileSystemHostIDStorage); ok {
				return nil, fmt.Errorf("failed to find this host on mackerel (You may want to delete file \"%s\" to register this host to an another organization): %s", fsStorage.HostIDFile(), lastErr.Error())
			}
			return nil, fmt.Errorf("failed to find this host on mackerel: %s", lastErr.Error())
		}
		if result.CustomIdentifier != "" && result.CustomIdentifier != customIdentifier {
			if fsStorage, ok := conf.HostIDStorage.(*config.FileSystemHostIDStorage); ok {
				return nil, fmt.Errorf("custom identifiers mismatch: this host = \"%s\", the host whose id is \"%s\" on mackerel.io = \"%s\" (File \"%s\" may be copied from another host. Try deleting it and restarting agent)", customIdentifier, hostID, result.CustomIdentifier, fsStorage.HostIDFile())
			}
			return nil, fmt.Errorf("custom identifiers mismatch: this host = \"%s\", the host whose id is \"%s\" on mackerel.io = \"%s\" (Host ID file may be copied from another host. Try deleting it and restarting agent)", customIdentifier, hostID, result.CustomIdentifier)
		}
	}

	hostSt := conf.HostStatus.OnStart
	if hostSt != "" && hostSt != result.Status {
		doRetry(func() error {
			lastErr = api.UpdateHostStatus(result.ID, hostSt)
			return filterErrorForRetry(lastErr)
		})
		if lastErr != nil {
			return nil, fmt.Errorf("failed to set default host status: %s, %s", hostSt, lastErr.Error())
		}
	}

	lastErr = conf.SaveHostID(result.ID)
	if lastErr != nil {
		return nil, fmt.Errorf("failed to save host ID: %s", lastErr.Error())
	}

	return result, nil
}

// prepareCustomIdentiferHosts collects the host information based on the
// configuration of the custom_identifier fields.
func prepareCustomIdentiferHosts(conf *config.Config, api *mackerel.API) map[string]*mackerel.Host {
	customIdentifierHosts := make(map[string]*mackerel.Host)
	for _, customIdentifier := range conf.ListCustomIdentifiers() {
		host, err := api.FindHostByCustomIdentifier(customIdentifier)
		if err != nil {
			logger.Warningf("Failed to retrieve the host of custom_identifier: %s, %s", customIdentifier, err)
			continue
		}
		customIdentifierHosts[customIdentifier] = host
	}
	return customIdentifierHosts
}

// Interval between each updating host specs.
var specsUpdateInterval = 1 * time.Hour

func delayByHost(host *mackerel.Host) int {
	s := sha1.Sum([]byte(host.ID))
	return int(s[len(s)-1]) % int(config.PostMetricsInterval.Seconds())
}

// App contains objects for running main loop of mackerel-agent
type App struct {
	Agent                 *agent.Agent
	Config                *config.Config
	Host                  *mackerel.Host
	API                   *mackerel.API
	CustomIdentifierHosts map[string]*mackerel.Host
	AgentMeta             *AgentMeta
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

func loop(app *App, termCh chan struct{}) error {
	quit := make(chan struct{})
	defer close(quit) // broadcast terminating

	// Periodically update host specs.
	go updateHostSpecsLoop(app, quit)

	postQueue := make(chan *postValue, app.Config.Connection.PostMetricsBufferSize)
	go enqueueLoop(app, postQueue, quit)

	postDelaySeconds := delayByHost(app.Host)
	initialDelay := postDelaySeconds / 2
	logger.Debugf("wait %d seconds before initial posting.", initialDelay)
	select {
	case <-termCh:
		return nil
	case <-time.After(time.Duration(initialDelay) * time.Second):
		app.Agent.InitPluginGenerators(app.API)
	}

	termMetricsCh := make(chan struct{})
	var termCheckerCh chan struct{}
	var termMetadataCh chan struct{}

	hasChecks := len(app.Agent.Checkers) > 0
	if hasChecks {
		termCheckerCh = make(chan struct{})
	}

	hasMetadataPlugins := len(app.Agent.MetadataGenerators) > 0
	if hasMetadataPlugins {
		termMetadataCh = make(chan struct{})
	}

	// fan-out termCh
	go func() {
		for range termCh {
			termMetricsCh <- struct{}{}
			if termCheckerCh != nil {
				termCheckerCh <- struct{}{}
			}
			if termMetadataCh != nil {
				termMetadataCh <- struct{}{}
			}
		}
	}()

	if hasChecks {
		go runCheckersLoop(app, termCheckerCh, quit)
	}

	if hasMetadataPlugins {
		go runMetadataLoop(app, termMetadataCh, quit)
	}

	lState := loopStateFirst
	for {
		select {
		case <-termMetricsCh:
			if lState == loopStateTerminating {
				return fmt.Errorf("received terminate instruction again. force return")
			}
			lState = loopStateTerminating
			if len(postQueue) <= 0 {
				return nil
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
				delaySeconds = app.Config.Connection.PostMetricsDequeueDelaySeconds
			case loopStateHadError:
				// TODO: better interval calculation. exponential backoff or so.
				delaySeconds = app.Config.Connection.PostMetricsRetryDelaySeconds
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

			// determine next loopState before sleeping
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
					return fmt.Errorf("received terminate instruction again. force return")
				}
				lState = loopStateTerminating
			}

			postValues := [](*mackerel.CreatingMetricsValue){}
			for _, v := range origPostValues {
				postValues = append(postValues, v.values...)
			}
			err := app.API.PostMetricsValues(postValues)
			if err != nil {
				logger.Warningf("Failed to post metrics value (will retry): %s", err.Error())
				if lState != loopStateTerminating {
					lState = loopStateHadError
				}
				go func() {
					for _, v := range origPostValues {
						v.retryCnt++
						// It is difficult to distinguish the error is server error or data error.
						// So, if retryCnt exceeded the configured limit, postValue is considered invalid and abandoned.
						if v.retryCnt > app.Config.Connection.PostMetricsRetryMax {
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
				return nil
			}
		}
	}
}

func updateHostSpecsLoop(app *App, quit chan struct{}) {
	for {
		app.UpdateHostSpecs()
		select {
		case <-quit:
			return
		case <-time.After(specsUpdateInterval):
			// nop
		}
	}
}

func enqueueLoop(app *App, postQueue chan *postValue, quit chan struct{}) {
	metricsResult := app.Agent.Watch(quit)
	for {
		select {
		case <-quit:
			return
		case result := <-metricsResult:
			created := float64(result.Created.Unix())
			creatingValues := [](*mackerel.CreatingMetricsValue){}
			for _, values := range result.Values {
				hostID := app.Host.ID
				if values.CustomIdentifier != nil {
					if host, ok := app.CustomIdentifierHosts[*values.CustomIdentifier]; ok {
						hostID = host.ID
					} else {
						continue
					}
				}
				for name, value := range values.Values {
					if math.IsNaN(value) || math.IsInf(value, 0) {
						logger.Warningf("Invalid value: hostID = %s, name = %s, value = %f\n is not sent.", hostID, name, value)
						continue
					}

					creatingValues = append(
						creatingValues,
						&mackerel.CreatingMetricsValue{
							HostID: hostID,
							Name:   name,
							Time:   created,
							Value:  value,
						},
					)
				}
			}
			logger.Debugf("Enqueuing task to post metrics.")
			postQueue <- newPostValue(creatingValues)
		}
	}
}

func runChecker(checker *checks.Checker, checkReportCh chan *checks.Report, reportImmediateCh chan struct{}, quit <-chan struct{}) {
	lastStatus := checks.StatusUndefined
	lastMessage := ""
	interval := checker.Interval()
	nextInterval := time.Duration(0)
	nextTime := time.Now()

	for {
		select {
		case <-time.After(nextInterval):
			report := checker.Check()
			logger.Debugf("checker %q: report=%v", checker.Name, report)

			// It is possible that `now` is much bigger than `nextTime` because of
			// laptop sleep mode or any reason.
			now := time.Now()
			nextInterval = interval - (now.Sub(nextTime) % interval)
			nextTime = now.Add(nextInterval)

			if checker.Config.Action != nil {
				env := []string{fmt.Sprintf("MACKEREL_STATUS=%s", report.Status), fmt.Sprintf("MACKEREL_PREVIOUS_STATUS=%s", lastStatus)}
				go func() {
					_, stderr, _, _ := checker.Config.Action.RunWithEnv(env)

					if stderr != "" {
						logger.Warningf("Checker %q action output stderr: %s", checker.Name, stderr)
					}
				}()
			}

			if report.Status == checks.StatusOK && report.Status == lastStatus && report.Message == lastMessage {
				// Do not report if nothing has changed
				continue
			}
			if report.Status == checks.StatusOK && checker.Config.PreventAlertAutoClose {
				// Do not report `OK` if `PreventAlertAutoClose`
				lastStatus = report.Status
				lastMessage = report.Message
				continue
			}
			checkReportCh <- report

			// If status has changed, send it immediately
			// but if the status was OK and it's first invocation of a check, do not
			if report.Status != lastStatus && !(report.Status == checks.StatusOK && lastStatus == checks.StatusUndefined) {
				logger.Debugf("checker %q: status has changed %v -> %v: send it immediately", checker.Name, lastStatus, report.Status)
				reportImmediateCh <- struct{}{}
			}

			lastStatus = report.Status
			lastMessage = report.Message
		case <-quit:
			return
		}
	}
}

// runCheckersLoop generates "checker" goroutines
// which run for each checker commands and one for HTTP POSTing
// the reports to Mackerel API.
func runCheckersLoop(app *App, termCheckerCh <-chan struct{}, quit <-chan struct{}) {
	checkReportCh := make(chan *checks.Report)
	reportImmediateCh := make(chan struct{})

	for _, checker := range app.Agent.Checkers {
		go runChecker(checker, checkReportCh, reportImmediateCh, quit)
	}

	exit := false
	for !exit {
		select {
		case <-time.After(1 * time.Minute):
		case <-termCheckerCh:
			logger.Debugf("received 'term' chan for checkers loop")
			exit = true
		case <-reportImmediateCh:
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

		if len(reports) == 0 {
			continue
		}

		// Do not report too many reports at once.
		const checkReportMaxSize = 10
		partialReports := make([]*checks.Report, 0, checkReportMaxSize)
		for i, report := range reports {
			logger.Debugf("reports[%d]: %#v", i, report)
			partialReports = append(partialReports, report)
			if len(partialReports) >= checkReportMaxSize {
				reportCheckMonitors(app, checkReportCh, partialReports)
				partialReports = make([]*checks.Report, 0, checkReportMaxSize)
			}
		}
		reportCheckMonitors(app, checkReportCh, partialReports)
	}
}

func reportCheckMonitors(app *App, checkReportCh chan *checks.Report, reports []*checks.Report) {
	err := app.API.ReportCheckMonitors(app.Host.ID, reports)
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

// collectHostSpecs collects host specs (correspond to "name", "meta", "interfaces" and "customIdentifier" fields in API v0)
func collectHostSpecs() (string, map[string]interface{}, []spec.NetInterface, string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", nil, nil, "", fmt.Errorf("failed to obtain hostname: %s", err.Error())
	}

	specGens := specGenerators()
	cGen := spec.SuggestCloudGenerator()
	if cGen != nil {
		specGens = append(specGens, cGen)
	}
	meta := spec.Collect(specGens)

	var customIdentifier string
	if cGen != nil {
		customIdentifier, err = cGen.SuggestCustomIdentifier()
		if err != nil {
			logger.Warningf("Error while suggesting custom identifier. err: %s", err.Error())
		}
	}

	interfaces, err := interfaceGenerator().Generate()
	if err != nil {
		return "", nil, nil, "", fmt.Errorf("failed to collect interfaces: %s", err.Error())
	}
	return hostname, meta, interfaces, customIdentifier, nil
}

func fillUpSpecMeta(meta map[string]interface{}, ver, rev string) map[string]interface{} {
	meta["agent-version"] = ver
	meta["agent-revision"] = rev
	meta["agent-name"] = buildUA(ver, rev)
	return meta
}

// UpdateHostSpecs updates the host information that is already registered on Mackerel.
func (app *App) UpdateHostSpecs() {
	logger.Debugf("Updating host specs...")

	hostname, meta, interfaces, customIdentifier, err := collectHostSpecs()
	if err != nil {
		logger.Errorf("While collecting host specs: %s", err)
		return
	}
	meta = fillUpSpecMeta(meta, app.AgentMeta.Version, app.AgentMeta.Revision)

	err = app.API.UpdateHost(app.Host.ID, mackerel.HostSpec{
		Name:             hostname,
		Meta:             meta,
		Interfaces:       interfaces,
		RoleFullnames:    app.Config.Roles,
		Checks:           app.Config.CheckNames(),
		DisplayName:      app.Config.DisplayName,
		CustomIdentifier: customIdentifier,
	})

	if err != nil {
		logger.Errorf("Error while updating host specs: %s", err)
	} else {
		logger.Debugf("Host specs sent.")
	}
}

func buildUA(ver, rev string) string {
	return fmt.Sprintf("mackerel-agent/%s (Revision %s)", ver, rev)
}

// NewMackerelClient returns Mackerel API client for mackerel-agent
func NewMackerelClient(apibase, apikey, ver, rev string, verbose bool) (*mackerel.API, error) {
	api, err := mackerel.NewAPI(apibase, apikey, verbose)
	if err != nil {
		return nil, err
	}
	api.UA = buildUA(ver, rev)
	api.DefaultHeaders = http.Header{}
	api.DefaultHeaders.Add("X-Agent-Version", ver)
	api.DefaultHeaders.Add("X-Revision", rev)
	return api, nil
}

// Prepare sets up API and registers the host data to the Mackerel server.
// Use returned values to call Run().
func Prepare(conf *config.Config, ameta *AgentMeta) (*App, error) {
	api, err := NewMackerelClient(conf.Apibase, conf.Apikey, ameta.Version, ameta.Revision, conf.Verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare an api: %s", err.Error())
	}

	host, err := prepareHost(conf, api)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare host: %s", err.Error())
	}

	return &App{
		Agent:  NewAgent(conf),
		Config: conf,
		Host:   host,
		API:    api,
		CustomIdentifierHosts: prepareCustomIdentiferHosts(conf, api),
		AgentMeta:             ameta,
	}, nil
}

// RunOnce collects specs and metrics, then output them to stdout.
func RunOnce(conf *config.Config, ameta *AgentMeta) error {
	graphdefs, hostSpec, metrics, err := runOncePayload(conf, ameta)
	if err != nil {
		return err
	}

	json, err := json.Marshal(map[string]interface{}{
		"host":    hostSpec,
		"metrics": metrics,
	})
	if err != nil {
		logger.Warningf("Error while marshaling graphdefs: err = %s, graphdefs = %v.", err.Error(), graphdefs)
		return err
	}
	fmt.Println(string(json))
	return nil
}

func runOncePayload(conf *config.Config, ameta *AgentMeta) ([]mackerel.CreateGraphDefsPayload, *mackerel.HostSpec, *agent.MetricsResult, error) {
	hostname, meta, interfaces, customIdentifier, err := collectHostSpecs()
	if err != nil {
		logger.Errorf("While collecting host specs: %s", err)
		return nil, nil, nil, err
	}
	meta = fillUpSpecMeta(meta, ameta.Version, ameta.Revision)

	origInterval := metricsInterval
	metricsInterval = 1 * time.Second
	defer func() {
		metricsInterval = origInterval
	}()
	ag := NewAgent(conf)
	graphdefs := ag.CollectGraphDefsOfPlugins()
	metrics := ag.CollectMetrics(time.Now())
	return graphdefs, &mackerel.HostSpec{
		Name:             hostname,
		Meta:             meta,
		Interfaces:       interfaces,
		RoleFullnames:    conf.Roles,
		Checks:           conf.CheckNames(),
		DisplayName:      conf.DisplayName,
		CustomIdentifier: customIdentifier,
	}, metrics, nil
}

// NewAgent creates a new instance of agent.Agent from its configuration conf.
func NewAgent(conf *config.Config) *agent.Agent {
	return &agent.Agent{
		MetricsGenerators:  prepareGenerators(conf),
		PluginGenerators:   pluginGenerators(conf),
		Checkers:           createCheckers(conf),
		MetadataGenerators: metadataGenerators(conf),
	}
}

// Run starts the main metric collecting logic and this function will never return.
func Run(app *App, termCh chan struct{}) error {
	logger.Infof("Start: apibase = %s, hostName = %s, hostID = %s", app.Config.Apibase, app.Host.Name, app.Host.ID)

	err := loop(app, termCh)
	if err == nil && app.Config.HostStatus.OnStop != "" {
		// TODO error handling. support retire(?)
		e := app.API.UpdateHostStatus(app.Host.ID, app.Config.HostStatus.OnStop)
		if e != nil {
			logger.Errorf("Failed update host status on stop: %s", e)
		}
	}
	return err
}

func createCheckers(conf *config.Config) []*checks.Checker {
	checkers := []*checks.Checker{}

	for name, pluginConfig := range conf.CheckPlugins {
		checker := &checks.Checker{
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

func pluginGenerators(conf *config.Config) []metrics.PluginGenerator {
	generators := []metrics.PluginGenerator{}

	for _, pluginConfig := range conf.MetricPlugins {
		generators = append(generators, metrics.NewPluginGenerator(pluginConfig))
	}

	return generators
}
