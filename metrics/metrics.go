package metrics

import mkr "github.com/mackerelio/mackerel-client-go"

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
	PrepareGraphDefs() ([]*mkr.GraphDefsParam, error)
	CustomIdentifier() *string
}

// PluginFaultError may be returned by [PluginGenerator.PrepareGraphDefs].
// This error indicates a bug in a plugin and should be logged for a user.
// Note that [PluginGenerator.PrepareGraphDefs] can also return other error types.
type PluginFaultError struct {
	Err error
}

func (e *PluginFaultError) Error() string {
	return e.Err.Error()
}

func (e *PluginFaultError) Unwrap() error {
	return e.Err
}
