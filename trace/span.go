package trace

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
)

// Span is a single dd-trace span
//
// https://github.com/DataDog/dd-trace-go/blob/v1/ddtrace/tracer/span.go
type Span struct {
	Operation string             `msgpack:"name"`              // operation name
	Service   string             `msgpack:"service"`           // service name (i.e. "grpc.server", "http.request")
	Resource  string             `msgpack:"resource"`          // resource name (i.e. "/user?id=123", "SELECT * FROM users")
	Type      string             `msgpack:"type"`              // protocol associated with the span (i.e. "web", "db", "cache")
	Start     int64              `msgpack:"start"`             // span start time expressed in nanoseconds since epoch
	Duration  int64              `msgpack:"duration"`          // duration of the span expressed in nanoseconds
	Meta      map[string]string  `msgpack:"meta,omitempty"`    // arbitrary map of metadata
	Metrics   map[string]float64 `msgpack:"metrics,omitempty"` // arbitrary map of numeric metrics
	SpanID    uint64             `msgpack:"span_id"`           // identifier of this span
	TraceID   uint64             `msgpack:"trace_id"`          // identifier of the root span
	ParentID  uint64             `msgpack:"parent_id"`         // identifier of the span's direct parent
	Error     int32              `msgpack:"error"`             // error status of the span; 0 means no errors
}

// Trace is a collection of spans with the same trace ID
type Trace []*Span

// Traces is a collection of traces for unmarshalling
type Traces []Trace

// DecodeSpans decodes an array of spans from a msgpack encoded array
func DecodeSpans(data []byte) ([]Span, error) {
	var traces Traces
	spans := []Span{}

	err := msgpack.Unmarshal(data, &traces)
	if err != nil {
		return spans, fmt.Errorf("error decoding spans: %v", err)
	}

	for _, trace := range traces {
		for _, span := range trace {
			spans = append(spans, *span)
		}
	}

	return spans, nil
}
