package command

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/mackerelio/mackerel-agent/agent"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/spec"
)

var logger = logging.GetLogger("command")

func prepareHost(root string, api *mackerel.API, metaGenerators []spec.Generator, interfaceGenerator spec.Generator, roleFullnames []string) (*mackerel.Host, error) {
	os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+os.Getenv("PATH"))
	os.Setenv("LANG", "C") // prevent changing outputs of some command, e.g. ifconfig.
	meta := spec.CollectMeta(metaGenerators)
	interfaces := spec.CollectInterfaces(interfaceGenerator)

	hostname, err := spec.GetHostname()
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain hostname: %s", err.Error())
	}

	var result *mackerel.Host
	if hostId, err := mackerel.LoadHostId(root); err != nil { // create
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
			return nil, fmt.Errorf("Failed to find this host on mackerel (You may want to delete file \"%s\" to register this host to an another organization): %s", mackerel.IdFilePath(root), err.Error())
		}
		err := api.UpdateHost(hostId, hostname, meta, interfaces, roleFullnames)
		if err != nil {
			return nil, fmt.Errorf("Failed to update this host: %s", err.Error())
		}
	}

	err = mackerel.SaveHostId(root, result.Id)
	if err != nil {
		logger.Criticalf("Failed to save host ID: %s", err.Error())
		os.Exit(1)
	}

	return result, nil
}

const METRICS_POST_DEQUEUE_DELAY = 30 * time.Second // delay for dequeuing from buffer queue
const METRICS_POST_RETRY_DELAY = 1 * time.Minute    // delay for retring a request that causes errors
const METRICS_POST_RETRY_MAX = 10                   // max numbers of retries for a request that causes errors
const METRICS_POST_BUFFER_SIZE = 30                 // max numbers of requests stored in buffer queue.

func delayByHost(host *mackerel.Host) time.Duration {
	s := sha1.Sum([]byte(host.Id))
	return time.Duration(int(s[len(s)-1])%60) * time.Second
}

func loop(ag *agent.Agent, api *mackerel.API, host *mackerel.Host) {
	metricsResult := ag.Watch()

	postQueue := make(chan []*mackerel.CreatingMetricsValue, METRICS_POST_BUFFER_SIZE)

	go func() {
		for values := range postQueue {
			tries := METRICS_POST_RETRY_MAX
			for {
				err := api.PostMetricsValues(values)
				if err == nil {
					logger.Debugf("Posting metrics succeeded.")
					break
				}
				logger.Errorf("Failed to post metrics value (will retry): %s", err.Error())

				tries -= 1
				if tries <= 0 {
					logger.Errorf("Give up retrying to post metrics.")
					break
				}

				logger.Debugf("Retrying to post metrics...")
				time.Sleep(METRICS_POST_RETRY_DELAY)
			}

			time.Sleep(METRICS_POST_DEQUEUE_DELAY)
		}
	}()

	postDelay := delayByHost(host)
	isFirstTime := true
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
			if isFirstTime { // request immediately to create graph defs of host
				isFirstTime = false
			} else {
				// Sending data at every 0 second from all hosts causes request flooding.
				// To prevent flooding, this loop sleeps for some seconds
				// which is specific to the ID of the host running agent on.
				// The sleep second is up to 60s.
				logger.Debugf("Sleeping %v to enqueue post request...", postDelay)
				time.Sleep(postDelay)
			}
			logger.Debugf("Enqueuing task to post metrics.")
			postQueue <- creatingValues
		}
	}
}

func metaGenerators() []spec.Generator {
	return  []spec.Generator{
		&spec.KernelGenerator{},
		&spec.CPUGenerator{},
		&spec.MemoryGenerator{},
		&spec.BlockDeviceGenerator{},
		&spec.FilesystemGenerator{},
	}
}

func interfaceGenerator() spec.Generator {
	return &spec.InterfaceGenerator{}
}

func metricsGenerators() []metrics.Generator {
	return []metrics.Generator{
		&metrics.Loadavg5Generator{},
		&metrics.CpuusageGenerator{Interval: 60},
		&metrics.MemoryGenerator{},
		&metrics.UptimeGenerator{},
		&metrics.InterfaceGenerator{Interval: 60},
		&metrics.DiskGenerator{Interval: 60},
	}
}

func Run(config mackerel.Config) {
	api, err := mackerel.NewApi(config.Apibase, config.Apikey, config.Verbose)
	if err != nil {
		logger.Criticalf("Failed to prepare an api: %s", err.Error())
		os.Exit(1)
	}

	host, err := prepareHost(config.Root, api, metaGenerators(), interfaceGenerator(), config.Roles)
	if err != nil {
		logger.Criticalf("Failed to run this agent: %s", err.Error())
		os.Exit(1)
	}

	logger.Infof("Start: apibase = %s, hostName = %s, hostId = %s", config.Apibase, host.Name, host.Id)

	metricsGenerators := metricsGenerators()
	for _, pluginConfig := range config.Plugin["metrics"] {
		metricsGenerators = append(metricsGenerators, &metrics.PluginGenerator{pluginConfig})
	}

	ag := &agent.Agent{metricsGenerators}
	loop(ag, api, host)
}
