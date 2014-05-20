package spec

import (
	"testing"

	"github.com/mackerelio/mackerel-agent/version"
)

func TestCollect(t *testing.T) {
	version.VERSION = "1.0.0"
	version.GITCOMMIT = "1234beaf"

	generators := []Generator{}
	specs := Collect(generators)

	if specs["agent-version"] != "1.0.0" {
		t.Error("version should be 1.0.0")
	}
	if specs["agent-revision"] != "1234beaf" {
		t.Error("revision should be 1234beaf")
	}
	if specs["agent-name"] != "mackerel-agent/1.0.0 (Revision 1234beaf)" {
		t.Error("agent-name should be 'mackerel-agent/1.0.0 Revision/1234beaf'")
	}
}
