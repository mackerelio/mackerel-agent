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
		status, msg, err := checkerOK.Check()
		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}
		if status != StatusOK {
			t.Errorf("status should be OK: %v", status)
		}
		if msg != "OK\n" {
			t.Errorf("wrong message: %q", msg)
		}
	}

	{
		status, msg, err := checkerWarning.Check()
		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}
		if status != StatusWarning {
			t.Errorf("status should be WARNING: %v", status)
		}
		if msg != "something is going wrong\n" {
			t.Errorf("wrong message: %q", msg)
		}
	}
}
