package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

type writer struct{}

func (*writer) Write([]byte) (int, error) {
	return 0, nil
}

func init() {
	exporter, _ := stdouttrace.New(stdouttrace.WithWriter(&writer{}))
	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	))
}
