package metadata

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
)

func TestMetadataGeneratorFetch(t *testing.T) {
	tests := []struct {
		command  []string
		metadata string
		err      bool
	}{
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", "{}"},
			metadata: `{}`,
			err:      false,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "1", "-metadata", "{}"},
			metadata: ``,
			err:      true,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", `{"example": "metadata", "foo": [100, 200, {}, null]}`},
			metadata: `{"example":"metadata","foo":[100,200,{},null]}`,
			err:      false,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", `{"example": metadata"}`},
			metadata: ``,
			err:      true,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", `"foobar"`},
			metadata: `"foobar"`,
			err:      false,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", `foobar`},
			metadata: ``,
			err:      true,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", "262144"},
			metadata: `262144`,
			err:      false,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", "true"},
			metadata: `true`,
			err:      false,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0", "-metadata", "null"},
			metadata: `null`,
			err:      false,
		},
		{
			command:  []string{"go", "run", "testdata/json.go", "-exit-code", "0"},
			metadata: ``,
			err:      true,
		},
	}
	for _, test := range tests {
		g := Generator{
			Config: &config.MetadataPlugin{
				CommandArgs: test.command,
			},
		}
		metadata, err := g.Fetch()
		if err != nil {
			if !test.err {
				t.Errorf("error occurred unexpectedly on command %v %s", test.command, err.Error())
			}
		} else {
			if test.err {
				t.Errorf("error did not occurr but error expected on command %v", test.command)
			}
			metadataStr, _ := json.Marshal(metadata)
			if string(metadataStr) != test.metadata {
				t.Errorf("metadata should be %q but got %q", test.metadata, string(metadataStr))
			}
		}
	}
}

func TestMetadataGeneratorSaveIsChanged(t *testing.T) {
	tests := []struct {
		prevmetadata string
		metadata     string
		ischanged    bool
	}{
		{
			prevmetadata: `{}`,
			metadata:     `{}`,
			ischanged:    false,
		},
		{
			prevmetadata: `{ "foo": [ 100, 200, null, {} ] }`,
			metadata:     `{"foo":[100,200,null,{}]}`,
			ischanged:    false,
		},
		{
			prevmetadata: `null`,
			metadata:     `{}`,
			ischanged:    true,
		},
		{
			prevmetadata: `[]`,
			metadata:     `{}`,
			ischanged:    true,
		},
	}
	for i, test := range tests {
		g := Generator{Cachefile: filepath.Join("testdata", ".mackerel-metadata-test-"+strconv.Itoa(i))}
		var prevmetadata interface{}
		_ = json.Unmarshal([]byte(test.prevmetadata), &prevmetadata)

		if err := g.Save(prevmetadata); err != nil {
			t.Errorf("Error should not occur in Save() but got: %s", err.Error())
		}

		var metadata interface{}
		_ = json.Unmarshal([]byte(test.metadata), &metadata)

		got := g.IsChanged(metadata)
		if got != test.ischanged {
			t.Errorf("IsChanged() should return %t but got %t for %v, %v", test.ischanged, got, prevmetadata, metadata)
		}

		if err := g.Clear(); err != nil {
			t.Errorf("Error should not occur in Clear() but got: %s", err.Error())
		}
		if g.PrevMetadata != nil {
			t.Errorf("metadata cache should be cleared but got: %v", g.PrevMetadata)
		}
	}
}

func TestMetadataGeneratorInterval(t *testing.T) {
	tests := []struct {
		interval *int32
		expected time.Duration
	}{
		{
			interval: pint(0),
			expected: 10 * time.Minute,
		},
		{
			interval: pint(1),
			expected: 10 * time.Minute,
		},
		{
			interval: pint(30),
			expected: 30 * time.Minute,
		},
		{
			interval: nil,
			expected: 10 * time.Minute,
		},
	}
	for _, test := range tests {
		g := Generator{
			Config: &config.MetadataPlugin{
				ExecutionInterval: test.interval,
			},
		}
		if g.Interval() != test.expected {
			t.Errorf("interval should be %v but got: %v", test.expected, g.Interval())
		}
	}
}

func pint(i int32) *int32 {
	return &i
}
