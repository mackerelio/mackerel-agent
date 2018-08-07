// +build linux

package linux

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
)

// BlockDeviceGenerator XXX
type BlockDeviceGenerator struct {
}

var blockDeviceLogger = logging.GetLogger("spec.block_device")

// Generate generate metric values
func (g *BlockDeviceGenerator) Generate() (interface{}, error) {
	fileInfos, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		// /sys/block is unavailable on some PaaS (Heroku, for example)
		if os.IsNotExist(err) {
			blockDeviceLogger.Debugf("Failed (skip this spec): %s", err)
			return nil, nil
		}
		blockDeviceLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}

	results := make(mackerel.BlockDevice)

	for _, fileInfo := range fileInfos {
		deviceName := fileInfo.Name()
		result := map[string]interface{}{}

		for _, key := range []string{"size", "removable"} {
			filename := path.Join("/sys/block", deviceName, key)
			if _, err := os.Stat(filename); err == nil {
				bytes, err := ioutil.ReadFile(filename)
				if err != nil {
					blockDeviceLogger.Errorf("Failed (skip this spec): %s", err)
					return nil, err
				}
				result[key] = strings.TrimSpace(string(bytes))
			}
		}

		for _, key := range []string{"model", "rev", "state", "timeout", "vendor"} {
			filename := path.Join("/sys/block", deviceName, "device", key)
			if _, err := os.Stat(filename); err == nil {
				bytes, err := ioutil.ReadFile(filename)
				if err != nil {
					blockDeviceLogger.Errorf("Failed (skip this spec): %s", err)
					return nil, err
				}
				result[key] = strings.TrimSpace(string(bytes))
			}
		}

		results[fileInfo.Name()] = result
	}

	return results, nil
}
