package checks

import (
	"runtime"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/util"
)

func TestChecker_Check(t *testing.T) {
	checkerOK := Checker{
		Config: config.PluginConfig{
			Command: "go run testdata/exit.go -code 0 -message OK",
		},
	}

	checkerWarning := Checker{
		Config: config.PluginConfig{
			Command: "go run testdata/exit.go -code 1 -message something_is_going_wrong",
		},
	}

	{
		report, err := checkerOK.Check()
		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}
		if report.Status != StatusOK {
			t.Errorf("status should be OK: %v", report.Status)
		}
		if report.Message != "OK\n" {
			t.Errorf("wrong message: %q", report.Message)
		}
	}

	{
		report, err := checkerWarning.Check()
		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}
		if report.Status != StatusWarning {
			t.Errorf("status should be WARNING: %v", report.Status)
		}
		if report.Message != "something_is_going_wrong\n" {
			t.Errorf("wrong message: %q", report.Message)
		}
	}

}

func TestChecker_CheckTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skip check timeout test in windows")
	}
	origTimeoutDur := util.TimeoutDuration
	defer func() {
		util.TimeoutDuration = origTimeoutDur
	}()
	util.TimeoutDuration = 1 * time.Second

	checkerTimeout := Checker{
		Config: config.PluginConfig{
			Command: "sleep 2",
		},
	}

	{
		report, err := checkerTimeout.Check()
		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}
		if report.Status != StatusUnknown {
			t.Errorf("status should be UNKNOWN: %v", report.Status)
		}
		if report.Message != "command timed out" {
			t.Errorf("wrong message: %q", report.Message)
		}
	}
}
