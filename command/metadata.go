package command

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metadata"
)

func metadataGenerators(conf *config.Config) []*metadata.Generator {
	generators := make([]*metadata.Generator, 0, len(conf.MetadataPlugins))

	workdir := os.Getenv("MACKEREL_PLUGIN_WORKDIR")
	if workdir == "" {
		workdir = os.TempDir()
	}

	for name, pluginConfig := range conf.MetadataPlugins {
		generator := &metadata.Generator{
			Name:     name,
			Config:   pluginConfig,
			Tempfile: filepath.Join(workdir, "mackerel-metadata", name),
		}
		logger.Debugf("Metadata plugin generator created: %#v %#v", generator, generator.Config)
		generators = append(generators, generator)
	}

	return generators
}

type metadataResult struct {
	namespace string
	metadata  interface{}
}

func runMetadataLoop(c *Context, termMetadataCh <-chan struct{}, quit <-chan struct{}) {
	resultCh := make(chan *metadataResult)
	for _, g := range c.Agent.MetadataGenerators {
		go runEachMetadataLoop(g, resultCh, quit)
	}

	exit := false
	for !exit {
		select {
		case <-time.After(1 * time.Minute):
		case <-termMetadataCh:
			logger.Debugf("received 'term' chan for metadata loop")
			exit = true
		}

		results := []*metadataResult{}
	ConsumeResults:
		for {
			select {
			case result := <-resultCh:
				results = append(results, result)
			default:
				break ConsumeResults
			}
		}

		for _, result := range results {
			resp, err := c.API.PutMetadata(c.Host.ID, result.namespace, result.metadata)
			// retry on 5XX errors
			if resp != nil && resp.StatusCode >= 500 {
				logger.Errorf("put metadata %q failed: status %s", result.namespace, resp.Status)
				resultCh <- &metadataResult{
					namespace: result.namespace,
					metadata:  result.metadata,
				}
				continue
			}
			if err != nil {
				logger.Errorf("put metadata %q failed: %v", result.namespace, err)
				clearMetadataCache(c.Agent.MetadataGenerators, result.namespace)
				continue
			}
		}
	}
}

func clearMetadataCache(generators []*metadata.Generator, namespace string) {
	for _, g := range generators {
		if g.Name == namespace {
			g.Clear()
			return
		}
	}
}

func runEachMetadataLoop(g *metadata.Generator, resultCh chan<- *metadataResult, quit <-chan struct{}) {
	interval := g.Interval()
	nextInterval := 10 * time.Second
	nextTime := time.Now()

	for {
		select {
		case <-time.After(nextInterval):
			metadata, err := g.Fetch()

			// case for laptop sleep mode (now >> nextTime + interval)
			now := time.Now()
			nextInterval = interval - (now.Sub(nextTime) % interval)
			nextTime = now.Add(nextInterval)

			if err != nil {
				logger.Warningf("metadata plugin %q: %s", g.Name, err.Error())
				continue
			}

			if !g.IsChanged(metadata) {
				logger.Debugf("skipping metadata %q, metadata does not change", g.Name)
				continue
			}

			if err := g.Save(metadata); err != nil {
				logger.Warningf("metadata plugin %q: %s", g.Name, err.Error())
				continue
			}

			logger.Debugf("generated metadata %q (saved cache to file: %s)", g.Name, g.Tempfile)
			resultCh <- &metadataResult{
				namespace: g.Name,
				metadata:  metadata,
			}

		case <-quit:
			return
		}
	}
}
