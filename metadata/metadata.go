package metadata

import (
	"github.com/mackerelio/mackerel-agent/config"
)

type MetadataGenerator struct {
	Name   string
	Config *config.MetadataPlugin
}
