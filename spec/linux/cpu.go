// +build linux

package linux

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
)

// CPUGenerator collects CPU specs
type CPUGenerator struct {
}

var cpuLogger = logging.GetLogger("spec.cpu")

func (g *CPUGenerator) generate(file io.Reader) (interface{}, error) {
	scanner := bufio.NewScanner(file)

	var results mackerel.CPU
	var cur map[string]interface{}
	var modelName string

	for scanner.Scan() {
		line := scanner.Text()
		kv := strings.SplitN(line, ":", 2)
		if len(kv) < 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		switch key {
		case "processor":
			cur = make(map[string]interface{})
			if modelName != "" {
				cur["model_name"] = modelName
			}
			results = append(results, cur)
		case "Processor", "system type":
			modelName = val
		case "vendor_id", "model", "stepping", "physical id", "core id", "model name", "cache size":
			cur[strings.Replace(key, " ", "_", -1)] = val
		case "cpu family":
			cur["family"] = val
		case "cpu cores":
			cur["cores"] = val
		case "cpu MHz":
			cur["mhz"] = val
		}
	}
	if err := scanner.Err(); err != nil {
		cpuLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}

	// Old kernels with CONFIG_SMP disabled has no "processor: " line
	if len(results) == 0 && modelName != "" {
		cur = make(map[string]interface{})
		cur["model_name"] = modelName
		results = append(results, cur)
	}

	return results, nil
}

// Generate cpu specs
func (g *CPUGenerator) Generate() (interface{}, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		cpuLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}
	defer file.Close()

	return g.generate(file)
}
