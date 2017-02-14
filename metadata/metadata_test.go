package metadata

import (
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
)

func TestMetadataGenerator(t *testing.T) {
	tests := []struct {
		command string
		message string
		err     bool
	}{
		{
			command: `go run testdata/json.go -exit-code 0 -message "{}"`,
			message: `{}`,
			err:     false,
		},
		{
			command: `go run testdata/json.go -exit-code 1 -message "{}"`,
			message: ``,
			err:     true,
		},
		{
			command: `go run testdata/json.go -exit-code 0 -message '{"example": "message"}'`,
			message: `{"example": "message"}`,
			err:     false,
		},
		{
			command: `go run testdata/json.go -exit-code 0 -message '{"example": message"}'`,
			message: ``,
			err:     true,
		},
		{
			command: `go run testdata/json.go -exit-code 0 -message '"foobar"'`,
			message: `"foobar"`,
			err:     false,
		},
		{
			command: `go run testdata/json.go -exit-code 0 -message foobar`,
			message: ``,
			err:     true,
		},
		{
			command: `go run testdata/json.go -exit-code 0 -message 16777216`,
			message: `16777216`,
			err:     false,
		},
		{
			command: `go run testdata/json.go -exit-code 0 -message true`,
			message: `true`,
			err:     false,
		},
		{
			command: `go run testdata/json.go -exit-code 0 -message null`,
			message: `null`,
			err:     false,
		},
	}
	for _, test := range tests {
		g := Generator{
			Config: &config.MetadataPlugin{
				Command: test.command,
			},
		}
		message, err := g.Fetch()
		if err != nil {
			if !test.err {
				t.Errorf("error occurred unexpectedly on command %q", test.command)
			}
		} else {
			if test.err {
				t.Errorf("error did not occurr but error expected on command %q", test.command)
			}
			if message != test.message {
				t.Errorf("message should be %q but got %q", test.message, message)
			}
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
			expected: 1 * time.Minute,
		},
		{
			interval: pint(1),
			expected: 1 * time.Minute,
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
