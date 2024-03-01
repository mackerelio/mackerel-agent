//go:build linux
// +build linux

package linux

import (
	"fmt"
	"strconv"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUGenerator collects CPU specs
type CPUGenerator struct {
}

var cpuLogger = logging.GetLogger("spec.cpu")

// Generate cpu specs
func (g *CPUGenerator) Generate() (interface{}, error) {
	infoStats, err := cpu.Info()
	if err != nil {
		cpuLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}
	results := make(mackerel.CPU, 0, len(infoStats))
	for _, infoStat := range infoStats {
		result := map[string]any{
			"vendor_id":   infoStat.VendorID,
			"model":       infoStat.Model,
			"stepping":    strconv.Itoa(int(infoStat.Stepping)),
			"physical_id": infoStat.PhysicalID,
			"core_id":     infoStat.CoreID,
			"cache_size":  fmt.Sprintf("%d KB", infoStat.CacheSize),
			"model_name":  infoStat.ModelName,
			"family":      infoStat.Family,
			"cores":       strconv.Itoa(int(infoStat.Cores)),
			"mhz":         strconv.FormatFloat(infoStat.Mhz, 'f', -1, 64),
		}
		results = append(results, result)
	}

	return results, nil
}
