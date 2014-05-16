package spec

type Generator interface {
	Key() string
	Generate() (interface{}, error)
}
