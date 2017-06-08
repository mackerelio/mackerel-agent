package spec

import (
	"fmt"
	"testing"
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
	generators := []Generator{
		&testStructOK{},
		&testStructErr{},
	}
	specs := Collect(generators)

	if specs["ok"] != 15 {
		t.Error("metric value of ok should be 15")
	}

	_, ok := specs["error"]
	if ok {
		t.Error("when error, metric should not be collected")
	}
}
