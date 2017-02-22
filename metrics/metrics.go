package metrics

import "github.com/mackerelio/mackerel-agent/mackerel"

// Values represents metric values
type Values map[string]float64

func merge(v1, v2 Values) Values {
	for k, v := range v2 {
		v1[k] = v
	}
	return v1
}

// ValuesCustomIdentifier holds the metric values with the optional custom identifier
type ValuesCustomIdentifier struct {
	Values           Values
	CustomIdentifier *string
}

// MergeValuesCustomIdentifiers merges the metric values and custom identifiers
func MergeValuesCustomIdentifiers(values []*ValuesCustomIdentifier, newValue *ValuesCustomIdentifier) []*ValuesCustomIdentifier {
	for _, value := range values {
		if value.CustomIdentifier == newValue.CustomIdentifier ||
			(value.CustomIdentifier != nil && newValue.CustomIdentifier != nil &&
				*value.CustomIdentifier == *newValue.CustomIdentifier) {
			value.Values = merge(value.Values, newValue.Values)
			return values
		}
	}
	return append(values, newValue)
}

// Generator generates metrics
type Generator interface {
	Generate() (Values, error)
}

// PluginGenerator generates metrics of plugin
type PluginGenerator interface {
	Generator
	PrepareGraphDefs() ([]mackerel.CreateGraphDefsPayload, error)
	CustomIdentifier() *string
}
