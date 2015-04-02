package metrics

import "github.com/mackerelio/mackerel-agent/mackerel"

// Values XXX
type Values map[string]float64

// Merge XXX
func (vs *Values) Merge(other Values) {
	for k, v := range (map[string]float64)(other) {
		(*vs)[k] = v
	}
}

// Generator XXX
type Generator interface {
	Generate() (Values, error)
}

// PluginGenerator XXX
type PluginGenerator interface {
	Generate() (Values, error)
	PrepareGraphDefs() ([]mackerel.CreateGraphDefsPayload, error)
}
