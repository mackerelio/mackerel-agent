package spec

import (
	"fmt"
	"testing"

	"github.com/mackerelio/mackerel-agent/version"
)

type testStructOK struct{}

func (tok *testStructOK) Key() string {
	return "ok"
}

func (tok *testStructOK) Generate() (interface{}, error) {
	return 15, nil
}

type testStructErr struct{}

func (tok *testStructErr) Key() string {
	return "error"
}

func (tok *testStructErr) Generate() (interface{}, error) {
	return nil, fmt.Errorf("error")
}

func TestCollect(t *testing.T) {
	version.VERSION = "1.0.0"
	version.GITCOMMIT = "1234beaf"

	generators := []Generator{
		&testStructOK{},
		&testStructErr{},
	}
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

	if specs["ok"] != 15 {
		t.Error("metric value of ok should be 15")
	}

	_, ok := specs["error"]
	if ok {
		t.Error("when error, metric should not be collected")
	}

}
