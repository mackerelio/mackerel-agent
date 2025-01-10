//go:build darwin
// +build darwin

package darwin

import "go.opentelemetry.io/otel"

var meter = otel.Meter("github.com/mackerelio/mackerel-agent/metrics/darwin")
