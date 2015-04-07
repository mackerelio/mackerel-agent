package checks

import (
	"fmt"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/util"
)

type Status string

const (
	StatusOK       Status = "OK"
	StatusWarning         = "WARNING"
	StatusCritical        = "CRITICAL"
	StatusUnknown         = "UNKNOWN"
)

var exitCodeToStatus = map[int]Status{
	0: StatusOK,
	1: StatusWarning,
	2: StatusCritical,
	3: StatusUnknown,
}

type Checker struct {
	Name string
	// NOTE(motemen): We make use of config.PluginConfig as it happens
	// to have Command field which was used by metrics.pluginGenerator.
	// If the configuration of checks.Checker and metrics.pluginGenerator goes different ones,
	// we should reconcider using config.PluginConfig.
	Config config.PluginConfig
}

func (c Checker) String() string {
	return fmt.Sprintf("checker %q command=[%s]", c.Name, c.Config.Command)
}

func (c Checker) Check() (status Status, message string, err error) {
	command := c.Config.Command
	stdout, stderr, exitCode, err := util.RunCommand(command)
	if err != nil {
		return StatusUnknown, "", err
	}

	fmt.Printf("%q %q", stdout, stderr)

	if s, ok := exitCodeToStatus[exitCode]; ok {
		status = s
	} else {
		status = StatusUnknown
	}

	return status, stdout, nil
}
