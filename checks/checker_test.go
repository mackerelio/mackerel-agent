package checks

import (
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
)

func TestChecker_Check(t *testing.T) {
	checkerOK := Checker{
		Config: config.PluginConfig{
			Command: "go run testdata/exit.go -code 0 -message OK",
		},
	}

	checkerWarning := Checker{
		Config: config.PluginConfig{
			Command: "go run testdata/exit.go -code 1 -message 'something is going wrong'",
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
		if report.Message != "something is going wrong\n" {
			t.Errorf("wrong message: %q", report.Message)
		}
	}
}
