package spec

import (
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("spec")

// Generator interface for generating spec values
type Generator interface {
	Key() string
	Generate() (interface{}, error)
}

// Collect spec values
func Collect(specGenerators []Generator) map[string]interface{} {
	specs := make(map[string]interface{})
	for _, g := range specGenerators {
		value, err := g.Generate()
		if err != nil {
			logger.Warningf("Failed to collect meta in %T (skip this spec): %s", g, err.Error())
			continue
		}
		specs[g.Key()] = value
	}
	return specs
}
