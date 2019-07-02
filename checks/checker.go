package checks

import (
	"fmt"
	"time"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/config"
)

var logger = logging.GetLogger("checks")

// Status is a status that is produced by periodical checking.
// It is currently compatible with Nagios.
type Status string

// Current possible statuses, which is taken from command's exit code.
// the mapping is given as exitCodeToStatus.
const (
	StatusUndefined Status = ""
	StatusOK        Status = "OK"
	StatusWarning   Status = "WARNING"
	StatusCritical  Status = "CRITICAL"
	StatusUnknown   Status = "UNKNOWN"
)

const defaultCheckInterval = 1 * time.Minute

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
	Name   string
	Config *config.CheckPlugin
}

// Report is what Checker produces by invoking its command.
type Report struct {
	Name                 string
	Status               Status
	Message              string
	OccurredAt           time.Time
	NotificationInterval *int32
	MaxCheckAttempts     *int32
	CustomIdentfier      *string
}

func (c *Checker) String() string {
	return fmt.Sprintf("checker %q command=[%s]", c.Name, c.Config.Command)
}

// Check invokes the command and transforms its result to a Report.
func (c *Checker) Check() *Report {
	now := time.Now()
	message, stderr, exitCode, err := c.Config.Command.Run()
	if stderr != "" {
		logger.Warningf("Checker %q output stderr: %s", c.Name, stderr)
	}

	status := StatusUnknown
	if err != nil {
		message = err.Error()
	} else {
		if s, ok := exitCodeToStatus[exitCode]; ok {
			status = s
		}

		logger.Debugf("Checker %q status=%s message=%q", c.Name, status, message)
	}

	return &Report{
		Name:                 c.Name,
		Status:               status,
		Message:              message,
		OccurredAt:           now,
		NotificationInterval: c.Config.NotificationInterval,
		MaxCheckAttempts:     c.Config.MaxCheckAttempts,
		CustomIdentfier:      c.Config.CustomIdentifier,
	}
}

// Interval is the interval where the command is invoked.
func (c *Checker) Interval() time.Duration {
	if c.Config.CheckInterval != nil {
		interval := time.Duration(*c.Config.CheckInterval) * time.Minute
		if interval < 1*time.Minute {
			interval = 1 * time.Minute
		} else if interval > 60*time.Minute {
			interval = 60 * time.Minute
		}
		return interval
	}
	return defaultCheckInterval
}
