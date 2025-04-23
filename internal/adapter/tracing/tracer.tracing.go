package tracingadapter

import (
	"context"
	"fmt"

	port "github.com/GoReactors/backend-learning/internal/port/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracerAdapter struct {
	name   string
	tracer trace.Tracer
}

func NewTracer(name string) port.Tracer {
	tracer_adapter := &TracerAdapter{
		name: name,
	}
	tracer_adapter.tracer = otel.Tracer(name)
	return tracer_adapter
}

func (t *TracerAdapter) GetTracer() trace.Tracer {
	return t.tracer
}

func (t *TracerAdapter) StartSpan(ctx context.Context, span_name string) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, span_name)
	span.SetAttributes(attribute.String("service.name", t.name))
	return ctx, span
}

func (t *TracerAdapter) SetAttribute(span trace.Span, key string, value any) {
	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, v))
	case int:
		span.SetAttributes(attribute.Int(key, v))
	case int64:
		span.SetAttributes(attribute.Int64(key, v))
	case float64:
		span.SetAttributes(attribute.Float64(key, v))
	case bool:
		span.SetAttributes(attribute.Bool(key, v))
	default:
		// fallback to string conversion
		span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
	}
}
