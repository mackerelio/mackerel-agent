// +build windows

package windows

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// PluginGenerator XXX
type PluginGenerator struct {
	Config config.PluginConfig
}

var pluginLogger = logging.GetLogger("metrics.plugin")

const pluginPrefix = "custom."

// NewPluginGenerator XXX
func NewPluginGenerator(c config.PluginConfig) (*PluginGenerator, error) {
	return &PluginGenerator{c}, nil
}

// Generate XXX
func (g *PluginGenerator) Generate() (metrics.Values, error) {
	if g == nil {
		err := errors.New("PluginGenerator is not initialized")
		pluginLogger.Criticalf(err.Error())
		return nil, err
	}
	results, err := g.collectValues(g.Config.Command)
	if err != nil {
		pluginLogger.Criticalf(err.Error())
		return nil, err
	}
	return results, nil
}

var delimReg = regexp.MustCompile(`[\s\t]+`)

func (g *PluginGenerator) collectValues(command string) (metrics.Values, error) {
	pluginLogger.Debugf("Executing plugin: command = \"%s\"", command)

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	cmd := exec.Command("cmd", "/c", command)
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
		items := delimReg.Split(line, 3)
		if len(items) != 3 {
			continue
		}
		value, err := strconv.ParseFloat(items[1], 64)
		if err != nil {
			pluginLogger.Warningf("Failed to parse values: %s", err)
			continue
		}
		results[pluginPrefix+items[0]] = value
	}

	return results, nil
}
