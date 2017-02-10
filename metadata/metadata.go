package metadata

import (
	"github.com/mackerelio/mackerel-agent/config"
)

// MetadataGenerator generates metadata
type MetadataGenerator struct {
	Name   string
	Config *config.MetadataPlugin
}
