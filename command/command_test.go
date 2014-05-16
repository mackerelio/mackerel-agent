package command

import (
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/spec"
	"github.com/mackerelio/mackerel-agent/version"
	"regexp"
	"testing"
)

func TestCollectSpecs(t *testing.T) {
	version.VERSION = "1.0.0"
	version.GITCOMMIT = "1234beaf"

	specGenerators := []spec.Generator{}
	specs := collectSpecs(specGenerators)

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

func TestGetHostname(t *testing.T) {
	hostname, err := getHostname()
	if err != nil {
		t.Error("should not raise error")
	}

	if !regexp.MustCompile(`\w+`).MatchString(hostname) {
		t.Error("hostname should have length and not contains whitespace but:", hostname)
	}
}

func TestDelayByHost(t *testing.T) {
	delay1 := delayByHost(&mackerel.Host{
		Id:     "246PUVUngPo",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})

	delay2 := delayByHost(&mackerel.Host{
		Id:     "21GZjCE5Etb",
		Name:   "hogehoge2.host.h",
		Type:   "unknown",
		Status: "working",
	})

	if !(0 <= delay1.Seconds() && delay1.Seconds() < 60) {
		t.Errorf("delay shoud be between 0 and 60 but %v", delay1)
	}

	if delay1 == delay2 {
		t.Error("delays shoud be different")
	}
}
