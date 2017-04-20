package command

import (
	"path/filepath"
	"time"

	"github.com/mackerelio/golib/pluginutil"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metadata"
)

func metadataGenerators(conf *config.Config) []*metadata.Generator {
	generators := make([]*metadata.Generator, 0, len(conf.MetadataPlugins))

	workdir := pluginutil.PluginWorkDir()
	for name, pluginConfig := range conf.MetadataPlugins {
		generator := &metadata.Generator{
			Name:      name,
			Config:    pluginConfig,
			Cachefile: filepath.Join(workdir, "mackerel-metadata", name),
		}
		logger.Debugf("Metadata plugin generator created: %#v %#v", generator, generator.Config)
		generators = append(generators, generator)
	}

	return generators
}

type metadataResult struct {
	namespace string
	metadata  interface{}
	createdAt time.Time
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

		results := make(map[string]*metadataResult)
	ConsumeResults:
		for {
			select {
			case result := <-resultCh:
				// prefer new result to avoid infinite number of retries
				if prev, ok := results[result.namespace]; ok {
					if result.createdAt.After(prev.createdAt) {
						results[result.namespace] = result
					}
				} else {
					results[result.namespace] = result
				}
			default:
				break ConsumeResults
			}
		}

		for _, result := range results {
			resp, err := c.API.PutMetadata(c.Host.ID, result.namespace, result.metadata)
			// retry on 5XX errors
			if resp != nil && resp.StatusCode >= 500 {
				logger.Errorf("put metadata %q failed: status %s", result.namespace, resp.Status)
				go func() {
					resultCh <- result
				}()
				continue
			}
			if err != nil {
				logger.Errorf("put metadata %q failed: %v", result.namespace, err)
				clearMetadataCache(c.Agent.MetadataGenerators, result.namespace)
				continue
			}
		}
		results = nil
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
				logger.Debugf("metadata plugin %q: metadata does not change", g.Name)
				continue
			}

			if err := g.Save(metadata); err != nil {
				logger.Warningf("metadata plugin %q: %s", g.Name, err.Error())
				continue
			}

			logger.Debugf("metadata plugin %q: generated metadata (and saved cache to file: %s)", g.Name, g.Cachefile)
			resultCh <- &metadataResult{
				namespace: g.Name,
				metadata:  metadata,
				createdAt: time.Now(),
			}

		case <-quit:
			return
		}
	}
}
