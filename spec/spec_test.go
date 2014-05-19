package spec

import (
	"github.com/mackerelio/mackerel-agent/version"
	"regexp"
	"testing"
)

func TestCollectMetas(t *testing.T) {
	version.VERSION = "1.0.0"
	version.GITCOMMIT = "1234beaf"

	generators := []Generator{}
	meta := CollectMeta(generators)

	if meta["agent-version"] != "1.0.0" {
		t.Error("version should be 1.0.0")
	}
	if meta["agent-revision"] != "1234beaf" {
		t.Error("revision should be 1234beaf")
	}
	if meta["agent-name"] != "mackerel-agent/1.0.0 (Revision 1234beaf)" {
		t.Error("agent-name should be 'mackerel-agent/1.0.0 Revision/1234beaf'")
	}
}

func TestGetHostname(t *testing.T) {
	hostname, err := GetHostname()
	if err != nil {
		t.Error("should not raise error")
	}

	if !regexp.MustCompile(`\w+`).MatchString(hostname) {
		t.Error("hostname should have length and not contains whitespace but:", hostname)
	}
}
