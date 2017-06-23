// +build freebsd

package freebsd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mackerelio/golib/logging"
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
		memoryLogger.Warningf("'top -bn 1' command exited with a non-zero status: '%s'", err)
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
			for i := 1; i < len(fields); i += 2 {
				v, err := getValue(fields[i])
				if err == nil {
					k := fields[i+1]
					if strings.HasSuffix(k, ",") {
						k = k[0 : len(k)-1]
					}
					switch k {
					case "Active":
						ret["memory.active"] = v
					case "Inact":
						ret["memory.inactive"] = v
					/*
						case "Wired":
							ret["memory.wired"] = v
					*/
					case "Cache":
						ret["memory.cached"] = v
					case "Buf":
						ret["memory.buffers"] = v
					case "Free":
						ret["memory.free"] = v
					}
				} else {
					errRet = err
				}
			}
		}
		if fields[0] == "Swap:" {
			if v, err := getValue(fields[1]); err == nil {
				ret["memory.swap_total"] = v
			} else {
				errRet = err
			}
			swapFreeIndex := 5
			if len(fields) == 5 {
				swapFreeIndex = 3
			}
			if v, err := getValue(fields[swapFreeIndex]); err == nil {
				ret["memory.swap_free"] = v
			} else {
				errRet = err
			}
		}
	}

	v, err := getTotalMem()
	if err != nil {
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
