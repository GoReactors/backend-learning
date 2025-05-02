package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceFieldsFromContext(ctx context.Context) []Field {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return nil
	}
	return []Field{
		{Key: "trace_id", Value: sc.TraceID().String()},
		{Key: "span_id", Value: sc.SpanID().String()},
	}
}
