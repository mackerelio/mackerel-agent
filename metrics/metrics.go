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

// ValuesCustomIdentifier holds the metric values with the optional custom identifier
type ValuesCustomIdentifier struct {
	Values           Values
	CustomIdentifier *string
}

// MergeValuesCustomIdentifiers merges the metric values and custom identifiers
func MergeValuesCustomIdentifiers(values []ValuesCustomIdentifier, newValue ValuesCustomIdentifier) []ValuesCustomIdentifier {
	for _, value := range values {
		if value.CustomIdentifier == newValue.CustomIdentifier ||
			(value.CustomIdentifier != nil && newValue.CustomIdentifier != nil &&
				*value.CustomIdentifier == *newValue.CustomIdentifier) {
			value.Values.Merge(newValue.Values)
			return values
		}
	}
	return append(values, newValue)
}

// Generator XXX
type Generator interface {
	Generate() (Values, error)
}

// PluginGenerator XXX
type PluginGenerator interface {
	Generate() (Values, error)
	PrepareGraphDefs() ([]mackerel.CreateGraphDefsPayload, error)
	CustomIdentifier() *string
}
