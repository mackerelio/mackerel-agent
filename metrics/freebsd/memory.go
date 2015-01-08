// +build freebsd

package freebsd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// MemoryGenerator XXX
type MemoryGenerator struct {
}

var memoryLogger = logging.GetLogger("metrics.memory")

// Generate generate metrics values
func (g *MemoryGenerator) Generate() (metrics.Values, error) {
	var errRet error
	outBytes, err := exec.Command("top", "-bn", "1").Output()
	if err != nil {
		logger.Warningf("'top -bn 1' command exited with a non-zero status: '%s'", err)
		errRet = err
	}
	ret := make(map[string]float64)

	out := string(outBytes)
	lines := strings.Split(out, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}

		if fields[0] == "Mem:" {

			if v, err := getValue(fields[1]); err == nil {
				ret["memory.active"] = v
			} else {
				errRet = err
			}
			if v, err := getValue(fields[3]); err == nil {
				ret["memory.inactive"] = v
			} else {
				errRet = err
			}
			/*
			         if v, err := getValue(fields[5]); err == nil {
			           ret["memory.wired"] = v
			         } else {
			   	errRet = err
			         }
			*/
			if v, err := getValue(fields[7]); err == nil {
				ret["memory.cached"] = v
			} else {
				errRet = err
			}
			if v, err := getValue(fields[9]); err == nil {
				ret["memory.buffers"] = v
			} else {
				errRet = err
			}
			if v, err := getValue(fields[11]); err == nil {
				ret["memory.free"] = v
			} else {
				errRet = err
			}
		}
		if fields[0] == "Swap:" {
			if v, err := getValue(fields[1]); err == nil {
				ret["memory.swap_total"] = v
			} else {
				errRet = err
			}
			if v, err := getValue(fields[5]); err == nil {
				ret["memory.swap_free"] = v
			} else {
				errRet = err
			}
		}
	}
	if v, err := getTotalMem(); err != nil {
		return nil, err
	}
	ret["memory.total"] = v

	if errRet == nil {
		return metrics.Values(ret), nil
	}
	return nil, errRet
}

func getValue(strValue string) (float64, error) {
	parseRegexp := regexp.MustCompile(`^(\d+)(.?)$`)
	match := parseRegexp.FindStringSubmatch(strValue)
	var unit float64 = 1
	switch {
	case match[2] == "G":
		unit = 1024 * 1024 * 1024
	case match[2] == "M":
		unit = 1024 * 1024
	case match[2] == "K":
		unit = 1024
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(match[1]), 64)
	if err != nil {
		return 0, fmt.Errorf("while parsing %q: %s", match[1], err)
	}
	return value * unit, nil
}

func getTotalMem() (float64, error) {
	cmd := exec.Command("sysctl", "-n", "hw.physmem")
	outputBytes, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("sysctl -n hw.physmem: %s", err)
	}

	output := string(outputBytes)

	memsizeInBytes, err := strconv.ParseFloat(strings.TrimSpace(output), 64)
	if err != nil {
		return 0, fmt.Errorf("while parsing %q: %s", output, err)
	}

	return memsizeInBytes, nil
}
