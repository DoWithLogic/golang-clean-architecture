package instrumentation

import (
	"context"
	"net/http/httptrace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// tracer: Global instance of the OpenTelemetry tracer, initialized with the service name "btb-cms".
var tracer = otel.Tracer("golang-clean-architecture")

// NewTraceSpan: Starts a new trace span with the given context and name, incorporating any baggage attributes into the span attributes.
// Parameters:
// - ctx: Context containing baggage attributes.
// - name: Name of the new trace span.
// Returns:
// - Context with the new span added.
// - The newly created trace span.
func NewTraceSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(
		ctx,
		name,
		trace.WithAttributes(ctxBaggageToAttributes(ctx)...),
	)
}

// NewTraceSpanWithoutBaggage: Starts a new trace span without incorporating any baggage attributes.
// Parameters:
// - ctx: Context for the new trace span.
// - name: Name of the new trace span.
// Returns:
// - Context with the new span added.
// - The newly created trace span.
func NewTraceSpanWithoutBaggage(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(
		ctx,
		name,
	)
}

// NewTraceHttpClient: Returns a new HTTP client trace for OpenTelemetry using the provided context.
// Parameters:
// - ctx: Context for the HTTP client trace.
// Returns:
// - A new HTTP client trace.
func NewTraceHttpClient(ctx context.Context) *httptrace.ClientTrace {
	return otelhttptrace.NewClientTrace(ctx)
}

// RecordSpanError: Records an error on the provided trace span and sets the status of the span to 'Error'.
// Parameters:
// - span: Trace span on which to record the error.
// - err: Error to record.
func RecordSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// ctxBaggageToAttributes: Converts baggage attributes from the context into OpenTelemetry attribute key-values.
// Parameters:
// - ctx: Context containing baggage attributes.
// Returns:
// - Slice of attribute key-values converted from the baggage attributes.
func ctxBaggageToAttributes(ctx context.Context) []attribute.KeyValue {
	var attributes []attribute.KeyValue

	bag := baggage.FromContext(ctx)
	for _, member := range bag.Members() {
		attributes = append(attributes, attribute.String(member.Key(), member.Value()))
	}

	return attributes
}
