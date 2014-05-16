package metrics

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
)

type PluginGenerator struct {
	Config mackerel.PluginConfig
}

var pluginLogger = logging.GetLogger("metrics.plugin")
var PLUGIN_PREFIX = "custom."

func (g *PluginGenerator) Generate() (Values, error) {
	results, err := g.collectValues(g.Config.Command)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (g *PluginGenerator) collectValues(command string) (Values, error) {
	pluginLogger.Debugf("Executing plugin: command = \"%s\"", command)

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		pluginLogger.Errorf("Failed to execute command \"%s\" (skip these metrics):\n%s", command, string(errBuffer.Bytes()))
		return nil, err
	}

	results := make(map[string]float64, 0)
	for _, line := range strings.Split(string(outBuffer.Bytes()), "\n") {
		// Key, value, timestamp
		// ex.) localhost.localdomain.tcp.CLOSING 0 1397031808
		items := strings.Split(line, "\t")
		if len(items) != 3 {
			continue
		}
		value, err := strconv.ParseFloat(items[1], 64)
		if err != nil {
			pluginLogger.Warningf("Failed to parse values: %s", err)
			continue
		}
		results[PLUGIN_PREFIX + items[0]] = value
	}

	return results, nil
}
