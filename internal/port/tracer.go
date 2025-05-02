package port

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Tracing interface {
	StartSpan(ctx context.Context, name string) (context.Context, Span)
	SetAttribute(span Span, key string, value any)
}

type Span = trace.Span
