package metadata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metadata")

// Generator generates metadata
type Generator struct {
	Name         string
	Config       *config.MetadataPlugin
	PrevMetadata interface{}
}

// Fetch invokes the command and returns the result
func (g *Generator) Fetch() (interface{}, error) {
	message, stderr, exitCode, err := g.Config.Run()

	if err != nil {
		logger.Warningf("Error occurred while executing a metadata plugin %q: %s", g.Name, err.Error())
		return nil, err
	}

	if stderr != "" {
		logger.Warningf("Metadata plugin %q outputs stderr: %s", g.Name, stderr)
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("Metadata plugin %q exits with: %d", g.Name, exitCode)
	}

	var metadata interface{}
	if err := json.Unmarshal([]byte(message), &metadata); err != nil {
		return nil, fmt.Errorf("Metadata plugin %q outputs invalid JSON: %v", g.Name, message)
	}

	return metadata, nil
}

// Differs returns whether the metadata has been changed or not
func (g *Generator) Differs(metadata interface{}) bool {
	return !reflect.DeepEqual(g.PrevMetadata, metadata)
}

// Save stores the metadata locally
func (g *Generator) Save(metadata interface{}) error {
	g.PrevMetadata = metadata
	return nil
}

const defaultExecutionInterval = 10 * time.Minute

// Interval calculates the time interval of command execution
func (g *Generator) Interval() time.Duration {
	if g.Config.ExecutionInterval == nil {
		return defaultExecutionInterval
	}
	interval := time.Duration(*g.Config.ExecutionInterval) * time.Minute
	if interval < 1*time.Minute {
		return 1 * time.Minute
	}
	return interval
}
