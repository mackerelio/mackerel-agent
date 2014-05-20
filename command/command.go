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

func prepareHost(root string, api *mackerel.API, specGenerators []spec.Generator, roleFullnames []string) (*mackerel.Host, error) {
	os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+os.Getenv("PATH"))
	os.Setenv("LANG", "C") // prevent changing outputs of some command, e.g. ifconfig.
	meta := spec.Collect(specGenerators)

	// retrieve intaface
	interfaces, _ := meta["interface"].([]map[string]interface{})
	delete(meta, "interface")

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain hostname: %s", err.Error())
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

func Run(conf config.Config) {
	api, err := mackerel.NewApi(conf.Apibase, conf.Apikey, conf.Verbose)
	if err != nil {
		logger.Criticalf("Failed to prepare an api: %s", err.Error())
		os.Exit(1)
	}

	host, err := prepareHost(conf.Root, api, specGenerators(), conf.Roles)
	if err != nil {
		logger.Criticalf("Failed to run this agent: %s", err.Error())
		os.Exit(1)
	}

	logger.Infof("Start: apibase = %s, hostName = %s, hostId = %s", conf.Apibase, host.Name, host.Id)

	ag := &agent.Agent{metricsGenerators(conf)}
	loop(ag, api, host)
}
