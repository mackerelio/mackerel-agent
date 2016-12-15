package agent

import (
	"time"

	"github.com/mackerelio/mackerel-agent/checks"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// Agent is the root of metrics collectors
type Agent struct {
	MetricsGenerators []metrics.Generator
	PluginGenerators  []metrics.PluginGenerator
	Checkers          []*checks.Checker
}

// MetricsResult XXX
type MetricsResult struct {
	Created time.Time
	Values  []*metrics.ValuesCustomIdentifier
}

// CollectMetrics collects metrics with generators.
func (agent *Agent) CollectMetrics(collectedTime time.Time) *MetricsResult {
	generators := agent.MetricsGenerators
	for _, g := range agent.PluginGenerators {
		generators = append(generators, g)
	}
	values := generateValues(generators)
	return &MetricsResult{Created: collectedTime, Values: values}
}

// Watch XXX
func (agent *Agent) Watch(quit chan struct{}) chan *MetricsResult {

	metricsResult := make(chan *MetricsResult)
	ticker := make(chan time.Time)
	interval := config.PostMetricsInterval

	go func() {
		t := time.NewTicker(1 * time.Second)

		last := time.Now()
		ticker <- last // sends tick once at first

		for {
			select {
			case <-quit:
				close(ticker)
				t.Stop()
				return
			case t := <-t.C:
				// Fire an event at 0 second per minute.
				// Because ticks may not be accurate,
				// fire an event if t - last is more than 1 minute
				if t.Second()%int(interval.Seconds()) == 0 || t.After(last.Add(interval)) {
					last = t
					ticker <- t
				}
			}
		}
	}()

	const collectMetricsWorkerMax = 3

	go func() {
		// Start collectMetrics concurrently
		// so that it does not prevent runnnig next collectMetrics.
		sem := make(chan struct{}, collectMetricsWorkerMax)
		for tickedTime := range ticker {
			ti := tickedTime
			sem <- struct{}{}
			go func() {
				metricsResult <- agent.CollectMetrics(ti)
				<-sem
			}()
		}
	}()

	return metricsResult
}

// CollectGraphDefsOfPlugins collects GraphDefs of Plugins
func (agent *Agent) CollectGraphDefsOfPlugins() []mackerel.CreateGraphDefsPayload {
	payloads := []mackerel.CreateGraphDefsPayload{}

	for _, g := range agent.PluginGenerators {
		p, err := g.PrepareGraphDefs()
		if err != nil {
			logger.Debugf("Failed to fetch meta information from plugin %s (non critical); seems that this plugin does not have meta information: %s", g, err)
		}
		if p != nil {
			payloads = append(payloads, p...)
		}
	}

	return payloads
}

// InitPluginGenerators XXX
func (agent *Agent) InitPluginGenerators(api *mackerel.API) {
	payloads := agent.CollectGraphDefsOfPlugins()

	if len(payloads) > 0 {
		err := api.CreateGraphDefs(payloads)
		if err != nil {
			logger.Errorf("Failed to create graphdefs: %s", err)
		}
	}
}
