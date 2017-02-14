package command

import (
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metadata"
)

func metadataGenerators(conf *config.Config) []*metadata.Generator {
	generators := make([]*metadata.Generator, 0, len(conf.MetadataPlugins))

	for name, pluginConfig := range conf.MetadataPlugins {
		generator := &metadata.Generator{
			Name:   name,
			Config: pluginConfig,
		}
		logger.Debugf("Metadata plugin generator created: %v", generator)
		generators = append(generators, generator)
	}

	return generators
}
