package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	Tempfile     string
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
		logger.Warningf("metadata plugin %q outputs stderr: %s", g.Name, stderr)
		// metadata plugin can output message to stderr for debugging and json to stdout
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("exits with: %d", exitCode)
	}

	var metadata interface{}
	if err := json.Unmarshal([]byte(message), &metadata); err != nil {
		return nil, fmt.Errorf("outputs invalid JSON: %v", message)
	}

	return metadata, nil
}

// Differs returns whether the metadata has been changed or not
func (g *Generator) Differs(metadata interface{}) bool {
	if g.PrevMetadata == nil {
		g.LoadFromFile()
	}
	return !reflect.DeepEqual(g.PrevMetadata, metadata)
}

// LoadFromFile loads the previous metadata from file
func (g *Generator) LoadFromFile() {
	data, err := ioutil.ReadFile(g.Tempfile)
	if err != nil { // maybe initial state
		return
	}
	var metadata interface{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		logger.Warningf("metadata plugin %q detected a invalid json in temporary file: %s", g.Name, string(data))
		// ignore errors, the file will be overwritten by Save()
		return
	}
	g.PrevMetadata = metadata
}

// Save stores the metadata locally
func (g *Generator) Save(metadata interface{}) error {
	g.PrevMetadata = metadata
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal the metadata to json: %v %s", metadata, err.Error())
	}
	if g.Tempfile == "" {
		return fmt.Errorf("specify the name of temporary file")
	}
	if err = writeFileAtomically(g.Tempfile, data); err != nil {
		return fmt.Errorf("failed to write the metadata to temporary file: %v %s", metadata, err.Error())
	}
	return nil
}

// Clear destroys the metadata cache
func (g *Generator) Clear() error {
	g.PrevMetadata = nil
	return os.Remove(g.Tempfile)
}

// writeFileAtomically writes contents to the file atomically
func writeFileAtomically(f string, contents []byte) error {
	tmpf, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(tmpf.Name())
	_, err = tmpf.Write(contents)
	if err != nil {
		return err
	}
	tmpf.Close()
	return os.Rename(tmpf.Name(), f)
}

const defaultExecutionInterval = 10 * time.Minute

// Interval calculates the time interval of command execution
func (g *Generator) Interval() time.Duration {
	if g.Config.ExecutionInterval == nil {
		return defaultExecutionInterval
	}
	interval := time.Duration(*g.Config.ExecutionInterval) * time.Minute
	if interval < defaultExecutionInterval {
		return defaultExecutionInterval
	}
	return interval
}
