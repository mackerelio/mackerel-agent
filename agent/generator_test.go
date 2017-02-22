package agent

import (
	"testing"

	"github.com/mackerelio/mackerel-agent/metrics"
)

type testGenerator struct{}

func (g *testGenerator) Generate() (metrics.Values, error) {
	values := make(map[string]float64)
	values["test"] = 10
	return values, nil
}

type testPanicGenerator struct{}

func (g *testPanicGenerator) Generate() (metrics.Values, error) {
	panic("sudden panic!!")
}

func TestGenerateValues(t *testing.T) {
	tg := &testGenerator{}
	tpg := &testPanicGenerator{}
	generators := []metrics.Generator{tg, tpg}
	values := generateValues(generators)

	if len(values) != 1 {
		t.Errorf("Num of results should be 1, but %d", len(values))
	}
}
