// +build linux

package linux

import (
	"bufio"
	"os"
	"regexp"

	"github.com/mackerelio/mackerel-agent/logging"
)

type CPUGenerator struct {
}

func (g *CPUGenerator) Key() string {
	return "cpu"
}

var cpuLogger = logging.GetLogger("spec.cpu")

func (g *CPUGenerator) Generate() (interface{}, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		cpuLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	results := make([]map[string]interface{}, 0)
	var cur map[string]interface{}
	for scanner.Scan() {
		line := scanner.Text()

		if matches := regexp.MustCompile(`^processor\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur = make(map[string]interface{})
			results = append(results, cur)
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
		} else if matches := regexp.MustCompile(`^flags\s+:\s+(.*)$`).FindStringSubmatch(line); matches != nil {
			cur["flags"] = regexp.MustCompile(` `).Split(matches[1], -1)
		}
	}
	if err := scanner.Err(); err != nil {
		cpuLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}

	return results, nil
}
