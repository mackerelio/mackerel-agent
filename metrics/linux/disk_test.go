// +build linux

package linux

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/metrics"
)

func TestDiskGenerator(t *testing.T) {
	if os.Getenv("TRAVIS") != "" {
		// ref. https://github.com/travis-ci/travis-ci/issues/2627
		t.Skipf("Skip: can't access `/proc/diskstats` in Travis environment.")
	}

	g := &DiskGenerator{1 * time.Second}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
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

func TestParseDiskStats(t *testing.T) {
	// insert empty line intentionally
	out := []byte(`202       1 xvda1 750193 3037 28116978 368712 16600606 7233846 424712632 23987908 0 2355636 24345740

202       2 xvda2 1641 9310 87552 1252 6365 3717 80664 24192 0 15040 25428`)

	result, err := parseDiskStats(out)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expect := metrics.Values{
		"disk.xvda1.reads":          750193,
		"disk.xvda1.readsMerged":    3037,
		"disk.xvda1.sectorsRead":    28116978,
		"disk.xvda1.readTime":       368712,
		"disk.xvda1.writes":         16600606,
		"disk.xvda1.writesMerged":   7233846,
		"disk.xvda1.sectorsWritten": 424712632,
		"disk.xvda1.writeTime":      23987908,
		"disk.xvda1.ioInProgress":   0,
		"disk.xvda1.ioTime":         2355636,
		"disk.xvda1.ioTimeWeighted": 24345740,
		"disk.xvda2.reads":          1641,
		"disk.xvda2.readsMerged":    9310,
		"disk.xvda2.sectorsRead":    87552,
		"disk.xvda2.readTime":       1252,
		"disk.xvda2.writes":         6365,
		"disk.xvda2.writesMerged":   3717,
		"disk.xvda2.sectorsWritten": 80664,
		"disk.xvda2.writeTime":      24192,
		"disk.xvda2.ioInProgress":   0,
		"disk.xvda2.ioTime":         15040,
		"disk.xvda2.ioTimeWeighted": 25428,
	}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result is not expected one: %+v", result)
	}
}
