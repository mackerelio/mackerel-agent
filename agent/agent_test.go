package agent

import (
	"context"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metrics"
	mkr "github.com/mackerelio/mackerel-client-go"
)

type fakeGenerator struct {
	metrics.Generator
	FakeGenerate func() (metrics.Values, error)
}

func (f *fakeGenerator) Generate() (metrics.Values, error) {
	return f.FakeGenerate()
}

type fakePluginGenerator struct {
	metrics.PluginGenerator
	FakeGenerate         func() (metrics.Values, error)
	FakeCustomIdentifier *string
}

func (f *fakePluginGenerator) Generate() (metrics.Values, error) {
	return f.FakeGenerate()
}

func (f *fakePluginGenerator) PrepareGraphDefs() ([]*mkr.GraphDefsParam, error) {
	return nil, nil
}

func (f *fakePluginGenerator) CustomIdentifier() *string {
	return f.FakeCustomIdentifier
}

func TestAgent_Watch(t *testing.T) {
	g1Cnt := 0
	g1 := &fakeGenerator{
		FakeGenerate: func() (metrics.Values, error) {
			g1Cnt++
			return map[string]float64{"g1.a": float64(g1Cnt)}, nil
		},
	}
	g2i := "g2"
	g2 := &fakePluginGenerator{
		FakeCustomIdentifier: &g2i,
		FakeGenerate: func() (metrics.Values, error) {
			return map[string]float64{"g2.a": float64(14), "g2.b": float64(42)}, nil
		},
	}

	ag := &Agent{MetricsGenerators: []metrics.Generator{g1, g2}}

	defer func(i time.Duration) { config.PostMetricsInterval = i }(config.PostMetricsInterval)
	// we cannot set interval less than 1 second
	config.PostMetricsInterval = 1 * time.Second

	ctx := context.Background()
	metricsResult := ag.Watch(ctx)

	cnt := 0
	end := time.After(5 * time.Second)
	for {
		select {
		case <-end:
			t.Error("timeout")
			return
		case mr := <-metricsResult:
			cnt++
			if cnt > 3 {
				return
			}
			v1, v2 := mr.Values[0], mr.Values[1]

			if v1.CustomIdentifier != nil {
				v1, v2 = v2, v1
			}

			if got, ok := v1.Values["g1.a"]; ok {
				if want := float64(cnt); got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			} else {
				t.Errorf("generator1: MetricResult should have 'g1.a': %v", v1)
			}

			if got, ok := v2.Values["g2.a"]; ok {
				if want := float64(14); got != want {
					t.Errorf("generator2: g2.a: got %v, want %v", got, want)
				}
			}
			if got, ok := v2.Values["g2.b"]; ok {
				if want := float64(42); got != want {
					t.Errorf("generator2: g2.b: got %v, want %v", got, want)
				}
			}
		}
	}
}
