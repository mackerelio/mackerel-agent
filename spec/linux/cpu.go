// +build linux

package linux

import (
	"bufio"
	"io"
	"os"
	"regexp"

	"github.com/mackerelio/mackerel-agent/logging"
)

// CPUGenerator Collects CPU specs
type CPUGenerator struct {
}

// Key XXX
func (g *CPUGenerator) Key() string {
	return "cpu"
}

var cpuLogger = logging.GetLogger("spec.cpu")

func (g *CPUGenerator) generate(file io.Reader) (interface{}, error) {
	scanner := bufio.NewScanner(file)

	var results []map[string]interface{}
	var cur map[string]interface{}
	var modelName string

	for scanner.Scan() {
		line := scanner.Text()

		if matches := regexp.MustCompile(`^processor\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur = make(map[string]interface{})
			if modelName != "" {
				cur["model_name"] = modelName
			}
			results = append(results, cur)
		} else if matches := regexp.MustCompile(`^Processor\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			modelName = matches[1]
		} else if matches := regexp.MustCompile(`^vendor_id\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["vendor_id"] = matches[1]
		} else if matches := regexp.MustCompile(`^cpu family\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["family"] = matches[1]
		} else if matches := regexp.MustCompile(`^model\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["model"] = matches[1]
		} else if matches := regexp.MustCompile(`^stepping\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["stepping"] = matches[1]
		} else if matches := regexp.MustCompile(`^physical id\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["physical_id"] = matches[1]
		} else if matches := regexp.MustCompile(`^core id\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["core_id"] = matches[1]
		} else if matches := regexp.MustCompile(`^cpu cores\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["cores"] = matches[1]
		} else if matches := regexp.MustCompile(`^model name\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["model_name"] = matches[1]
		} else if matches := regexp.MustCompile(`^cpu MHz\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["mhz"] = matches[1]
		} else if matches := regexp.MustCompile(`^cache size\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["cache_size"] = matches[1]
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

// Generate XXX
func (g *CPUGenerator) Generate() (interface{}, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		cpuLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}
	defer file.Close()

	return g.generate(file)
}
