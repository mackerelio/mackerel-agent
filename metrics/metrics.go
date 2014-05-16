package metrics

type Values map[string]float64

func (vs *Values) Merge(other Values) {
	for k, v := range (map[string]float64)(other) {
		(*vs)[k] = v
	}
}

type Generator interface {
	Generate() (Values, error)
}
