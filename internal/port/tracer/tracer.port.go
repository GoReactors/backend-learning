package port

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	InitTracer()
	StartSpan(ctx context.Context, name string) (context.Context, trace.Span)
	SetAttribute(span trace.Span, key string, value any)
}
