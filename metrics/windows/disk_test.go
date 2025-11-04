//go:build windows

package windows

import (
	"testing"
	"time"
)

func TestDiskGenerator(t *testing.T) {
	g, err := NewDiskGenerator(nil, 1*time.Second)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
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
		if value, ok := values["disk.sda."+metric]; !ok {
			t.Errorf("Value for disk.sda.%s should be collected", metric)
		} else {
			t.Logf("Disk '%s' collected: %+v", metric, value)
		}
	}

	for _, metric := range metrics {
		if value, ok := values["disk.sda."+metric+".delta"]; !ok {
			t.Errorf("Value for disk.sda.%s.delta should be collected", metric)
		} else {
			t.Logf("Disk '%s' delta collected: %+v", metric, value)
		}
	}

	for _, key := range metrics {
		if value, ok := values["disk.loop0."+key]; ok {
			t.Errorf("Value for disk.loop0.%s should not be collected but got %v. The value won't change.", key, value)
		}
	}
}
