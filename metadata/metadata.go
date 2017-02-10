package metadata

import (
	"errors"
	"fmt"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metadata")

// MetadataGenerator generates metadata
type MetadataGenerator struct {
	Name   string
	Config *config.MetadataPlugin
}

// Fetch invokes the command and returns the result
func (g *MetadataGenerator) Fetch() (string, error) {
	message, stderr, exitCode, err := g.Config.Run()

	if err != nil {
		logger.Warningf("Error occurred while executing a metadata plugin %q: %s", g.Name, err.Error())
		return "", err
	}

	if stderr != "" {
		logger.Warningf("Metadata generator %q outputs stderr: %s", g.Name, stderr)
	}

	if exitCode != 0 {
		return "", errors.New(fmt.Sprintf("Metadata plugin command exits with: %d", exitCode))
	}

	return message, nil
}

const defaultExecutionInterval = 10 * time.Minute

// Interval calculates the time interval of command execution
func (g *MetadataGenerator) Interval() time.Duration {
	if g.Config.ExecutionInterval == nil {
		return defaultExecutionInterval
	}
	interval := time.Duration(*g.Config.ExecutionInterval) * time.Minute
	if interval < 1*time.Minute {
		return 1 * time.Minute
	}
	return interval
}
