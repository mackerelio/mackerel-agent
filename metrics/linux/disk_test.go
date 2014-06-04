// +build linux

package linux

import (
	"testing"
)

func TestDiskGenerator(t *testing.T) {
	g := &DiskGenerator{1}
	values, err := g.Generate()
	if err != nil {
		t.Error("should not raise error")
	}

	metrics := []string{
		"reads", "readsMerged", "sectorsRead", "readTime",
		"writes", "writesMerged", "sectorsWritten", "writeTime",
		"ioInProgress", "ioTime", "ioTimeWeighted",
	}

	if _, ok := values["disk.sda.reads"]; !ok {
		t.Skipf("Skip: this node does not have sda device")
	}

	for _, metric := range metrics {
		if value, ok := values["disk.sda."+metric+".delta"]; !ok {
			t.Errorf("Value for disk.sda.%s.delta should be collected", metric)
		} else {
			t.Logf("Disk '%s' delta collected: %+v", metric, value)
		}
	}

	for _, key := range metrics {
		if value, ok := values["disk.loop0."+key+".delta"]; ok {
			t.Errorf("Value for disk.loop0.%s should not be collected but got %v. The value won't change.", key, value)
		}
	}
}
