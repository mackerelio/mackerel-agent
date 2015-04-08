package checks

import (
	"fmt"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/util"
)

var logger = logging.GetLogger("checks")

// Status is a status that is produced by periodical checking.
// It is currently compatible with Nagios.
type Status string

// Current possible statuses, which is taken from command's exit code.
// the mapping is given as exitCodeToStatus.
const (
	StatusOK       Status = "OK"
	StatusWarning  Status = "WARNING"
	StatusCritical Status = "CRITICAL"
	StatusUnknown  Status = "UNKNOWN"
)

var exitCodeToStatus = map[int]Status{
	0: StatusOK,
	1: StatusWarning,
	2: StatusCritical,
	3: StatusUnknown,
}

// Checker is the main interface of check monitoring.
// It invokes its given command and transforms the result to a Report
// to be sent to Mackerel periodically.
type Checker struct {
	Name string
	// NOTE(motemen): We make use of config.PluginConfig as it happens
	// to have Command field which was used by metrics.pluginGenerator.
	// If the configuration of checks.Checker and metrics.pluginGenerator goes different ones,
	// we should reconcider using config.PluginConfig.
	Config config.PluginConfig
}

// Report is what Checker produces by invoking its command.
type Report struct {
	Name       string
	Status     Status
	Message    string
	OccurredAt time.Time
}

func (c Checker) String() string {
	return fmt.Sprintf("checker %q command=[%s]", c.Name, c.Config.Command)
}

// Check invokes the command and transforms its result to a Report.
func (c Checker) Check() (*Report, error) {
	now := time.Now()

	command := c.Config.Command
	logger.Debugf("Checker %q executing command %q", c.Name, command)
	stdout, _, exitCode, err := util.RunCommand(command)
	if err != nil {
		return nil, err
	}

	status := StatusUnknown
	if s, ok := exitCodeToStatus[exitCode]; ok {
		status = s
	}

	logger.Debugf("Checker %q status=%s", c.Name, status)

	return &Report{
		Name:       c.Name,
		Status:     status,
		Message:    stdout,
		OccurredAt: now,
	}, nil
}

// Interval is the interval where the command is invoked.
// (Will be configurable in the future)
func (c Checker) Interval() time.Duration {
	return 1 * time.Minute
}
