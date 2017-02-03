package checks

import (
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
)

func TestChecker_Check(t *testing.T) {
	checkerOK := Checker{
		Config: &config.CheckPlugin{
			Command: "go run testdata/exit.go -code 0 -message OK",
		},
	}

	checkerWarning := Checker{
		Config: &config.CheckPlugin{
			Command: "go run testdata/exit.go -code 1 -message something_is_going_wrong",
		},
	}

	{
		report := checkerOK.Check()
		if report.Status != StatusOK {
			t.Errorf("status should be OK: %v", report.Status)
		}
		if report.Message != "OK\n" {
			t.Errorf("wrong message: %q", report.Message)
		}
	}

	{
		report := checkerWarning.Check()
		if report.Status != StatusWarning {
			t.Errorf("status should be WARNING: %v", report.Status)
		}
		if report.Message != "something_is_going_wrong\n" {
			t.Errorf("wrong message: %q", report.Message)
		}
	}
}
