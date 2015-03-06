// +build linux

package linux

import (
	"os"
	"testing"
)

func TestDiskGenerator(t *testing.T) {
	if os.Getenv("TRAVIS") != "" {
		// ref. https://github.com/travis-ci/travis-ci/issues/2627
		t.Skipf("Skip: can't access `/proc/diskstats` in Travis environment.")
	}

	g := &DiskGenerator{1}
	values, err := g.Generate()
	if err != nil {
		t.Error("should not raise error")
	}

	metrics := []string{
		"reads", "writes",
	}

	if _, ok := values["disk.sda.reads.delta"]; !ok {
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
