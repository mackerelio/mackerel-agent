package command

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mackerelio/mackerel-agent/agent"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/spec"
)

var logger = logging.GetLogger("command")

const idFileName = "id"

func IdFilePath(root string) string {
	return filepath.Join(root, idFileName)
}

func LoadHostId(root string) (string, error) {
	content, err := ioutil.ReadFile(IdFilePath(root))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SaveHostId(root string, id string) error {
	err := os.MkdirAll(root, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(IdFilePath(root))
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

// prepareHost collects specs of the host and sends them to Mackerel server.
// A unique host-id is returned by the server if one is not specified.
func prepareHost(root string, api *mackerel.API, roleFullnames []string) (*mackerel.Host, error) {
	// XXX this configuration should be moved to under spec/linux
	os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+os.Getenv("PATH"))
	os.Setenv("LANG", "C") // prevent changing outputs of some command, e.g. ifconfig.

	hostname, meta, interfaces, err := collectHostSpecs()
	if err != nil {
		return nil, fmt.Errorf("error while collecting host specs: %s", err.Error())
	}

	var result *mackerel.Host
	if hostId, err := LoadHostId(root); err != nil { // create
		logger.Debugf("Registering new host on mackerel...")
		createdHostId, err := api.CreateHost(hostname, meta, interfaces, roleFullnames)
		if err != nil {
			return nil, fmt.Errorf("Failed to register this host: %s", err.Error())
		}

		result, err = api.FindHost(createdHostId)
		if err != nil {
			return nil, fmt.Errorf("Failed to find this host on mackerel: %s", err.Error())
		}
	} else { // update
		result, err = api.FindHost(hostId)
		if err != nil {
			return nil, fmt.Errorf("Failed to find this host on mackerel (You may want to delete file \"%s\" to register this host to an another organization): %s", IdFilePath(root), err.Error())
		}
		err := api.UpdateHost(hostId, hostname, meta, interfaces, roleFullnames)
		if err != nil {
			return nil, fmt.Errorf("Failed to update this host: %s", err.Error())
		}
	}

	err = SaveHostId(root, result.Id)
	if err != nil {
		return nil, fmt.Errorf("Failed to save host ID: %s", err.Error())
	}

	return result, nil
}

// Interval between each updating host specs.
var specsUpdateInterval = 1 * time.Hour

func delayByHost(host *mackerel.Host) int {
	s := sha1.Sum([]byte(host.Id))
	return int(s[len(s)-1]) % int(config.PostMetricsInterval.Seconds())
}

type loopState uint8

const (
	loopStateFirst loopState = iota
	loopStateDefault
	loopStateQueued
	loopStateHadError
)

func loop(ag *agent.Agent, conf *config.Config, api *mackerel.API, host *mackerel.Host) {
	metricsResult := ag.Watch()

	postQueue := make(chan *postValue, conf.Connection.Post_Metrics_Buffer_Size)
	go func() {
		postDelaySeconds := delayByHost(host)
		lState := loopStateFirst
		for v := range postQueue {
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
				delaySeconds = conf.Connection.Post_Metrics_Dequeue_Delay_Seconds
			case loopStateHadError:
				delaySeconds = conf.Connection.Post_Metrics_Retry_Delay_Seconds
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

			// update loopState before sleeping
			if len(postQueue) > 0 {
				lState = loopStateQueued
			} else {
				lState = loopStateDefault
			}

			time.Sleep(time.Duration(delaySeconds) * time.Second)

			postValues := [](*mackerel.CreatingMetricsValue){}
			for _, v := range origPostValues {
				postValues = append(postValues, v.values...)
			}
			err := api.PostMetricsValues(postValues)
			if err != nil {
				logger.Errorf("Failed to post metrics value (will retry): %s", err.Error())
				lState = loopStateHadError
				go func() {
					for _, v := range origPostValues {
						v.retryCnt++
						if v.retryCnt <= conf.Connection.Post_Metrics_Retry_Max {
							postQueue <- v
						}
					}
				}()
				continue
			}
			logger.Debugf("Posting metrics succeeded.")
		}
	}()

	// Periodically update host specs.
	go func() {
		for {
			time.Sleep(specsUpdateInterval)
			UpdateHostSpecs(conf, api, host)
		}
	}()

	for {
		select {
		case result := <-metricsResult:
			created := float64(result.Created.Unix())
			creatingValues := [](*mackerel.CreatingMetricsValue){}
			for name, value := range (map[string]float64)(result.Values) {
				creatingValues = append(
					creatingValues,
					&mackerel.CreatingMetricsValue{host.Id, name, created, value},
				)
			}
			logger.Debugf("Enqueuing task to post metrics.")
			postQueue <- &postValue{creatingValues, 0}
		}
	}
}

type postValue struct {
	values   []*mackerel.CreatingMetricsValue
	retryCnt int
}

// collectHostSpecs collects host specs (correspond to "name", "meta" and "interfaces" fields in API v0)
func collectHostSpecs() (string, map[string]interface{}, []map[string]interface{}, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to obtain hostname: %s", err.Error())
	}

	meta := spec.Collect(specGenerators())

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

	err = api.UpdateHost(host.Id, hostname, meta, interfaces, conf.Roles)
	if err != nil {
		logger.Errorf("Error while updating host specs: %s", err)
	} else {
		logger.Debugf("Host specs sent.")
	}
}

// Prepare sets up API and registers the host data to the Mackerel server.
// Use returned values to call Run().
func Prepare(conf *config.Config) (*mackerel.API, *mackerel.Host, error) {
	api, err := mackerel.NewApi(conf.Apibase, conf.Apikey, conf.Verbose)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to prepare an api: %s", err.Error())
	}

	host, err := prepareHost(conf.Root, api, conf.Roles)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to preapre host: %s", err.Error())
	}

	return api, host, nil
}

// Run starts the main metric collecting logic and this function will never return.
func Run(conf *config.Config, api *mackerel.API, host *mackerel.Host) {
	logger.Infof("Start: apibase = %s, hostName = %s, hostId = %s", conf.Apibase, host.Name, host.Id)

	ag := &agent.Agent{
		MetricsGenerators: metricsGenerators(conf),
		PluginGenerators:  pluginGenerators(conf),
	}
	ag.InitPluginGenerators(api)

	loop(ag, conf, api, host)
}

func pluginGenerators(conf *config.Config) []metrics.PluginGenerator {
	generators := []metrics.PluginGenerator{}

	for _, pluginConfig := range conf.Plugin["metrics"] {
		generators = append(generators, metrics.NewPluginGenerator(pluginConfig))
	}

	return generators
}
