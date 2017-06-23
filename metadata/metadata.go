package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/config"
)

var logger = logging.GetLogger("metadata")

// Generator generates metadata
type Generator struct {
	Name         string
	Config       *config.MetadataPlugin
	Cachefile    string
	PrevMetadata interface{}
}

// Fetch invokes the command and returns the result
func (g *Generator) Fetch() (interface{}, error) {
	message, stderr, exitCode, err := g.Config.Run()

	if err != nil {
		logger.Warningf("Error occurred while executing a metadata plugin %q: %v", g.Name, err)
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

// IsChanged returns whether the metadata has been changed or not
func (g *Generator) IsChanged(metadata interface{}) bool {
	if g.PrevMetadata == nil {
		g.LoadFromFile()
	}
	return !reflect.DeepEqual(g.PrevMetadata, metadata)
}

// LoadFromFile loads the previous metadata from file
func (g *Generator) LoadFromFile() {
	data, err := ioutil.ReadFile(g.Cachefile)
	if err != nil { // maybe initial state
		return
	}
	var metadata interface{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		logger.Warningf("metadata plugin %q detected a invalid json in the cache file: %s", g.Name, string(data))
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
		return fmt.Errorf("failed to marshal the metadata to json: %v", err)
	}
	if g.Cachefile == "" {
		return fmt.Errorf("specify the name of the metadata cache file")
	}
	if err := os.MkdirAll(filepath.Dir(g.Cachefile), 0755); err != nil {
		return err
	}
	if err := writeFileAtomically(g.Cachefile, data); err != nil {
		return fmt.Errorf("failed to write the metadata to the cache file: %v", err)
	}
	return nil
}

// Clear destroys the metadata cache
func (g *Generator) Clear() error {
	g.PrevMetadata = nil
	return os.Remove(g.Cachefile)
}

// writeFileAtomically writes contents to the file atomically
func writeFileAtomically(f string, contents []byte) error {
	// MUST be located on same disk partition
	tmpf, err := ioutil.TempFile(filepath.Dir(f), "tmp")
	if err != nil {
		return err
	}
	// os.Remove here works successfully when tmpf.Write fails or os.Rename fails.
	// In successful case, os.Remove fails because the temporary file is already renamed.
	defer os.Remove(tmpf.Name())
	_, err = tmpf.Write(contents)
	tmpf.Close() // should be called before rename
	if err != nil {
		return err
	}
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
