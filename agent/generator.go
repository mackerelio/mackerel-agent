package agent

import (
	"sync"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

var logger = logging.GetLogger("agent")

func generateValues(generators []metrics.Generator) chan metrics.Values {
	processed := make(chan metrics.Values)
	finish := make(chan bool)
	result := make(chan metrics.Values)

	go func() {
		allValues := metrics.Values(make(map[string]float64))
		for {
			select {
			case values := <-processed:
				allValues.Merge(values)
			case <-finish:
				result <- allValues
				return
			}
		}
	}()

	go func() {
		var wg sync.WaitGroup
		for _, g := range generators {
			wg.Add(1)
			go func(g metrics.Generator) {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf("Panic: generating value in %T (skip this metric): %s", g, r)
					}
					wg.Done()
				}()

				values, err := g.Generate()
				if err != nil {
					logger.Errorf("Failed to generate value in %T (skip this metric): %s", g, err.Error())
					return
				}
				processed <- values
			}(g)
		}
		wg.Wait()
		finish <- true // processed all jobs
	}()

	return result
}
