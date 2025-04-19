package tracing

import (
	"context"
	"fmt"
	"time"

	"github.com/GoReactors/backend-learning/config"
	"github.com/GoReactors/backend-learning/pkg/constants"
	"github.com/GoReactors/backend-learning/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(cfg config.Config, l logger.Logger) (*sdktrace.TracerProvider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	serviceName := cfg.ServiceName

	collectorEndpoint := cfg.OtelCollectorEndpoint
	var exporter *otlptrace.Exporter
	var err error
	if cfg.Environment == constants.EnvDevelopment {
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(collectorEndpoint),
			otlptracehttp.WithInsecure(),
		)
	} else {
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(collectorEndpoint),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	l.Info(fmt.Sprintf("[Tracing] Initialized OTLP exporter -> %s", collectorEndpoint))

	return tp, nil
}
