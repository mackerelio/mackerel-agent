package agent

import (
	"time"

	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
)

type Agent struct {
	MetricsGenerators []metrics.Generator
	PluginGenerators  []metrics.PluginGenerator
}

type MetricsResult struct {
	Created time.Time
	Values  metrics.Values
}

func (agent *Agent) collectMetrics(collectedTime time.Time) *MetricsResult {
	generators := agent.MetricsGenerators
	for _, g := range agent.PluginGenerators {
		generators = append(generators, g)
	}
	result := generateValues(generators)
	values := <-result
	return &MetricsResult{Created: collectedTime, Values: values}
}

func (agent *Agent) Watch() chan *MetricsResult {

	metricsResult := make(chan *MetricsResult)
	ticker := make(chan time.Time)

	go func() {
		c := time.Tick(1 * time.Second)

		last := time.Now()
		ticker <- last // sends tick once at first

		for t := range c {
			// Fire an event at 0 second per minute.
			// Because ticks may not be accurate,
			// fire an event if t - last is more than 1 minute
			if t.Second() == 0 || t.After(last.Add(1*time.Minute)) {
				last = t
				ticker <- t
			}
		}
	}()

	const COLLECT_METRICS_WORKER_MAX = 3

	go func() {
		// Start collectMetrics concurrently
		// so that it does not prevent runnnig next collectMetrics.
		sem := make(chan uint, COLLECT_METRICS_WORKER_MAX)
		for tickedTime := range ticker {
			sem <- 1
			go func() {
				metricsResult <- agent.collectMetrics(tickedTime)
				<-sem
			}()
		}
	}()

	return metricsResult
}

func (agent *Agent) InitPluginGenerators(api *mackerel.API) {
	for _, g := range agent.PluginGenerators {
		err := g.InitWithAPI(api)
		if err != nil {
			logger.Errorf("Error while initializing plugin %s: %s", g, err)
		}
	}
}
