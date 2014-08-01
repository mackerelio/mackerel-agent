package metrics

import "github.com/mackerelio/mackerel-agent/mackerel"

type Values map[string]float64

func (vs *Values) Merge(other Values) {
	for k, v := range (map[string]float64)(other) {
		(*vs)[k] = v
	}
}

type Generator interface {
	Generate() (Values, error)
}

type PluginGenerator interface {
	Generate() (Values, error)
	InitWithAPI(api *mackerel.API) error
}
