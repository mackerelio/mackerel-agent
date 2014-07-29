// +build linux

package linux

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"

	"github.com/BurntSushi/toml"
)

// PluginGenerator collects user-defined metrics.
// mackerel-agent runs specified command and parses the result for the metric names and values.
type PluginGenerator struct {
	Config    config.PluginConfig
	ConfigKey string
	Meta      *pluginMeta
	Version   int
}

// pluginMeta is generated from plugin command. (not the configuration file)
type pluginMeta struct {
	Graphs map[string]customGraphDef
}

type customGraphDef struct {
	Prefix  string
	Label   string
	Unit    string
	Metrics map[string]customGraphMetricDef
}

type customGraphMetricDef struct {
	Name    string
	Label   string
	Stacked bool
}

var pluginLogger = logging.GetLogger("metrics.plugin")
var PLUGIN_PREFIX = "custom."

var pluginConfigurationEnvName = "MACKEREL_AGENT_PLUGIN_META"

func (g *PluginGenerator) Generate() (metrics.Values, error) {
	results, err := g.collectValues()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (g *PluginGenerator) InitWithAPI(api *mackerel.API) error {
	var err error
	g.Meta, err = g.loadPluginMeta()
	if err != nil {
		return err
	}

	payload := g.makeCreateGraphDefsPayload()
	if payload == nil {
		// this plugin does not provide graph definitions
		return nil
	}

	return api.CreateGraphDefs(payload)
}

// loadPluginMeta obtains plugin information (e.g. graph visuals, metric
// namespaces, etc) from the command specified.
// mackerel-agent runs the command with MACKEREL_AGENT_PLUGIN_META
// environment variable set.  The command is supposed to output like below:
//
//	# mackerel-agent-plugin
//	TBD
//
// The output should start with a line beginning with '#', which contains
// meta-info of the configuration. (eg. plugin schema version)
func (g *PluginGenerator) loadPluginMeta() (*pluginMeta, error) {
	command := g.Config.Command
	pluginLogger.Debugf("Obtaining plugin configuration: %q", command)

	// Set environment variable to make the plugin command generate its configuration
	os.Setenv(pluginConfigurationEnvName, "1")
	defer os.Setenv(pluginConfigurationEnvName, "")

	var outBuffer, errBuffer bytes.Buffer

	os.Setenv(pluginConfigurationEnvName, "")
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("running %q failed: %s, stderr=%q", command, err, string(errBuffer.Bytes()))
	}

	// Read the plugin configuration meta (version etc)

	headerLine, err := outBuffer.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("while reading the first line of command %q: %s", command, err)
	}

	// Parse the header line of format:
	// # mackerel-agent-plugin [key=value]...
	pluginMetaHeader := map[string]string{}

	re := regexp.MustCompile(`^#\s*mackerel-agent-plugin\b(.*)`)
	m := re.FindStringSubmatch(headerLine)
	if m == nil {
		return nil, fmt.Errorf("bad format of first line: %q", headerLine)
	}

	for _, field := range strings.Fields(m[1]) {
		keyValue := strings.Split(field, "=")
		var value string
		if len(keyValue) > 1 {
			value = keyValue[1]
		} else {
			value = ""
		}
		pluginMetaHeader[keyValue[0]] = value
	}

	// Check schema version
	version, ok := pluginMetaHeader["version"]
	if !ok {
		version = "1"
	}

	if version != "1" {
		return nil, fmt.Errorf("unsupported plugin meta version: %q", version)
	}

	g.Version = 1

	conf := &pluginMeta{}
	_, err = toml.DecodeReader(&outBuffer, conf)

	if err != nil {
		return nil, fmt.Errorf("while reading plugin configuration: %s", err)
	}

	return conf, nil
}

func (g *PluginGenerator) makeCreateGraphDefsPayload() []mackerel.CreateGraphDefsPayload {
	if g.Meta == nil {
		return nil
	}

	payloads := []mackerel.CreateGraphDefsPayload{}

	for key, graph := range g.Meta.Graphs {
		payload := mackerel.CreateGraphDefsPayload{
			Name:        PLUGIN_PREFIX + key,
			DisplayName: graph.Label,
			Unit:        graph.Unit,
		}
		if payload.Unit == "" {
			payload.Unit = "float"
		}

		for metricKey, metric := range graph.Metrics {
			metricPayload := mackerel.CreateGraphDefsPayloadMetric{
				Name:        PLUGIN_PREFIX + key + "." + metricKey,
				DisplayName: metric.Label,
				IsStacked:   metric.Stacked,
			}
			payload.Metrics = append(payload.Metrics, metricPayload)
		}

		payloads = append(payloads, payload)
	}

	return payloads
}

func (g *PluginGenerator) collectValues() (metrics.Values, error) {
	command := g.Config.Command
	pluginLogger.Debugf("Executing plugin: command = \"%s\"", command)

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	os.Setenv(pluginConfigurationEnvName, "")
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

		key := items[0]

		if g.Version == 1 {
			key = g.ConfigKey + "." + key
		}

		results[PLUGIN_PREFIX+key] = value
	}

	return results, nil
}
