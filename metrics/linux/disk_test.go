//go:build linux

package linux

import (
	"reflect"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/metrics"
)

func TestDiskGenerator(t *testing.T) {
	g := &DiskGenerator{Interval: 1 * time.Second}
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
}

func TestParseDiskStats(t *testing.T) {
	g := &DiskGenerator{Interval: 1 * time.Second}
	// insert empty line intentionally
	out := []byte(`202       1 xvda1 750193 3037 28116978 368712 16600606 7233846 424712632 23987908 0 2355636 24345740

202       2 xvda2 1641 9310 87552 1252 6365 3717 80664 24192 0 15040 25428`)

	var emptyMapping map[string]string
	result, err := g.parseDiskStats(out, emptyMapping)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expect := metrics.Values{
		"disk.xvda1.reads":          metrics.NewValueAttribute(750193),
		"disk.xvda1.readsMerged":    metrics.NewValueAttribute(3037),
		"disk.xvda1.sectorsRead":    metrics.NewValueAttribute(28116978),
		"disk.xvda1.readTime":       metrics.NewValueAttribute(368712),
		"disk.xvda1.writes":         metrics.NewValueAttribute(16600606),
		"disk.xvda1.writesMerged":   metrics.NewValueAttribute(7233846),
		"disk.xvda1.sectorsWritten": metrics.NewValueAttribute(424712632),
		"disk.xvda1.writeTime":      metrics.NewValueAttribute(23987908),
		"disk.xvda1.ioInProgress":   metrics.NewValueAttribute(0),
		"disk.xvda1.ioTime":         metrics.NewValueAttribute(2355636),
		"disk.xvda1.ioTimeWeighted": metrics.NewValueAttribute(24345740),
		"disk.xvda2.reads":          metrics.NewValueAttribute(1641),
		"disk.xvda2.readsMerged":    metrics.NewValueAttribute(9310),
		"disk.xvda2.sectorsRead":    metrics.NewValueAttribute(87552),
		"disk.xvda2.readTime":       metrics.NewValueAttribute(1252),
		"disk.xvda2.writes":         metrics.NewValueAttribute(6365),
		"disk.xvda2.writesMerged":   metrics.NewValueAttribute(3717),
		"disk.xvda2.sectorsWritten": metrics.NewValueAttribute(80664),
		"disk.xvda2.writeTime":      metrics.NewValueAttribute(24192),
		"disk.xvda2.ioInProgress":   metrics.NewValueAttribute(0),
		"disk.xvda2.ioTime":         metrics.NewValueAttribute(15040),
		"disk.xvda2.ioTimeWeighted": metrics.NewValueAttribute(25428),
	}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result is not expected one: %+v", result)
	}

	mapping := map[string]string{
		"xvda1": "_some_mount",
		"xvda3": "_nonused_mount",
	}
	resultWithMapping, err := g.parseDiskStats(out, mapping)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expectWithMapping := metrics.Values{
		"disk._some_mount.reads":          metrics.NewValueAttribute(750193),
		"disk._some_mount.readsMerged":    metrics.NewValueAttribute(3037),
		"disk._some_mount.sectorsRead":    metrics.NewValueAttribute(28116978),
		"disk._some_mount.readTime":       metrics.NewValueAttribute(368712),
		"disk._some_mount.writes":         metrics.NewValueAttribute(16600606),
		"disk._some_mount.writesMerged":   metrics.NewValueAttribute(7233846),
		"disk._some_mount.sectorsWritten": metrics.NewValueAttribute(424712632),
		"disk._some_mount.writeTime":      metrics.NewValueAttribute(23987908),
		"disk._some_mount.ioInProgress":   metrics.NewValueAttribute(0),
		"disk._some_mount.ioTime":         metrics.NewValueAttribute(2355636),
		"disk._some_mount.ioTimeWeighted": metrics.NewValueAttribute(24345740),
		"disk.xvda2.reads":                metrics.NewValueAttribute(1641),
		"disk.xvda2.readsMerged":          metrics.NewValueAttribute(9310),
		"disk.xvda2.sectorsRead":          metrics.NewValueAttribute(87552),
		"disk.xvda2.readTime":             metrics.NewValueAttribute(1252),
		"disk.xvda2.writes":               metrics.NewValueAttribute(6365),
		"disk.xvda2.writesMerged":         metrics.NewValueAttribute(3717),
		"disk.xvda2.sectorsWritten":       metrics.NewValueAttribute(80664),
		"disk.xvda2.writeTime":            metrics.NewValueAttribute(24192),
		"disk.xvda2.ioInProgress":         metrics.NewValueAttribute(0),
		"disk.xvda2.ioTime":               metrics.NewValueAttribute(15040),
		"disk.xvda2.ioTimeWeighted":       metrics.NewValueAttribute(25428),
	}
	if !reflect.DeepEqual(resultWithMapping, expectWithMapping) {
		t.Errorf("result is not expected one: %+v", resultWithMapping)
	}
}

func TestParseDiskStats_MoreFields(t *testing.T) {
	g := &DiskGenerator{Interval: 1 * time.Second}
	// There are 18 columns since Linux 4.18+.
	out := []byte(`202       1 xvda1 750193 3037 28116978 368712 16600606 7233846 424712632 23987908 0 2355636 24345740 0 0 0 0
  7       0 loop0 15 0 0 0 0 0 0 0 0 0 0 0 0 0 0`)

	var emptyMapping map[string]string
	result, err := g.parseDiskStats(out, emptyMapping)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expect := metrics.Values{
		"disk.xvda1.reads":          metrics.NewValueAttribute(750193),
		"disk.xvda1.readsMerged":    metrics.NewValueAttribute(3037),
		"disk.xvda1.sectorsRead":    metrics.NewValueAttribute(28116978),
		"disk.xvda1.readTime":       metrics.NewValueAttribute(368712),
		"disk.xvda1.writes":         metrics.NewValueAttribute(16600606),
		"disk.xvda1.writesMerged":   metrics.NewValueAttribute(7233846),
		"disk.xvda1.sectorsWritten": metrics.NewValueAttribute(424712632),
		"disk.xvda1.writeTime":      metrics.NewValueAttribute(23987908),
		"disk.xvda1.ioInProgress":   metrics.NewValueAttribute(0),
		"disk.xvda1.ioTime":         metrics.NewValueAttribute(2355636),
		"disk.xvda1.ioTimeWeighted": metrics.NewValueAttribute(24345740),
		"disk.loop0.reads":          metrics.NewValueAttribute(15),
		"disk.loop0.readsMerged":    metrics.NewValueAttribute(0),
		"disk.loop0.sectorsRead":    metrics.NewValueAttribute(0),
		"disk.loop0.readTime":       metrics.NewValueAttribute(0),
		"disk.loop0.writes":         metrics.NewValueAttribute(0),
		"disk.loop0.writesMerged":   metrics.NewValueAttribute(0),
		"disk.loop0.sectorsWritten": metrics.NewValueAttribute(0),
		"disk.loop0.writeTime":      metrics.NewValueAttribute(0),
		"disk.loop0.ioInProgress":   metrics.NewValueAttribute(0),
		"disk.loop0.ioTime":         metrics.NewValueAttribute(0),
		"disk.loop0.ioTimeWeighted": metrics.NewValueAttribute(0),
	}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result is not expected one: %+v", result)
	}
}

func TestParseDiskStats_ShouldIgnoreIfAllFieldsAreZeroOrSpecificDeviceName(t *testing.T) {
	g := &DiskGenerator{Interval: 1 * time.Second}
	out := []byte(`253       0 dm-0 2 0 40 0 314 0 2512 2136 0 236 2136
253       1 dm-1 964 0 57886 944 74855 0 644512 5421192 0 1580 5422136
  7       0 loop0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0`)

	var emptyMapping map[string]string
	result, err := g.parseDiskStats(out, emptyMapping)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	expect := metrics.Values{}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result is not expected one: %+v", result)
	}
}
