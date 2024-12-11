package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func runOTelLoop(ctx context.Context, app *App, termC <-chan struct{}) {
	logger.Debugf("Initializing OpenTelemetry...")
	shutdown, err := setupOTel(ctx, app)
	if err != nil {
		logger.Warningf("OpenTelemetry: %s", err.Error())
		return
	}
	logger.Debugf("OpenTelemetry initialized.")

	<-termC

	logger.Debugf("Stopping OpenTelemetry...")
	if err := shutdown(ctx); err != nil {
		logger.Warningf("OpenTelemetry: %s", err.Error())
		return
	}
	logger.Debugf("OpenTelemetry stopped.")
}

func setupOTel(ctx context.Context, app *App) (shutdown func(context.Context) error, err error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName("mackerel-agent"),
			semconv.HostName(app.Host.Name),
			attribute.Key("mackerelio.host.id").String(app.Host.ID),
		),
	)
	if err != nil {
		if !errors.Is(err, resource.ErrPartialResource) && !errors.Is(err, resource.ErrSchemaURLConflict) {
			return nil, err
		}
	}

	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint("otlp.mackerelio.com:4317"),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithHeaders(map[string]string{
			"Mackerel-Api-Key": app.API.APIKey,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(60*time.Second),
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}
