package command

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/mackerelio/golib/pluginutil"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metadata"
	mkr "github.com/mackerelio/mackerel-client-go"
)

func metadataGenerators(conf *config.Config) []*metadata.Generator {
	generators := make([]*metadata.Generator, 0, len(conf.MetadataPlugins))

	workdir := filepath.Join(pluginutil.PluginWorkDir(), "mackerel-metadata")
	for name, pluginConfig := range conf.MetadataPlugins {
		generator := &metadata.Generator{
			Name:      name,
			Config:    pluginConfig,
			Cachefile: getCacheFileName(name, workdir, pluginConfig),
		}
		logger.Debugf("Metadata plugin generator created: %#v %#v", generator, generator.Config)
		generators = append(generators, generator)
	}

	return generators
}

// The directory configuration in the env config of metadata should work as
// same as metric plugins. Since the working directory of metadata plugin is
// handled by mackerel-agent (not the plugin process), we have to lookup here.
func lookupPluginWorkDir(env []string) string {
	workDirPrefix := "MACKEREL_PLUGIN_WORKDIR="
	for _, e := range env {
		if strings.HasPrefix(e, workDirPrefix) {
			return strings.TrimPrefix(e, workDirPrefix)
		}
	}
	return ""
}

func getCacheFileName(name, defaultWorkDir string, plugin *config.MetadataPlugin) string {
	if dir := lookupPluginWorkDir(plugin.Command.Env); dir != "" {
		return filepath.Join(dir, "mackerel-plugin-metadata-"+name)
	}
	return filepath.Join(defaultWorkDir, name)
}

type metadataResult struct {
	namespace string
	metadata  interface{}
	createdAt time.Time
}

func runMetadataLoop(ctx context.Context, app *App, termMetadataCh <-chan struct{}) {
	resultCh := make(chan *metadataResult)
	for _, g := range app.Agent.MetadataGenerators {
		go runEachMetadataLoop(ctx, g, resultCh)
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
			err := app.API.PutHostMetaData(app.Host.ID, result.namespace, result.metadata)
			// retry on 5XX errors
			if mackerel.IsServerError(err) {
				e := err.(*mkr.APIError)
				logger.Errorf("put metadata %q failed: status %s", result.namespace, e.StatusCode)
				go func() {
					resultCh <- result
				}()
				continue
			}
			if err != nil {
				logger.Errorf("put metadata %q failed: %v", result.namespace, err)
				clearMetadataCache(app.Agent.MetadataGenerators, result.namespace)
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

func runEachMetadataLoop(ctx context.Context, g *metadata.Generator, resultCh chan<- *metadataResult) {
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

		case <-ctx.Done():
			return
		}
	}
}
