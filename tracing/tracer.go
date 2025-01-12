package tracing

import (
	"context"
	"fmt"

	"github.com/orders-app/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer() func() {
	// Create a console exporter
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("failed to create stdouttrace exporter: %v", err))

	}

	// Create a tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("orders-app"),
		)),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)

	// Return a function to clean up the tracer provider
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Log.Fatal(fmt.Sprintf("failed to shutdown tracer provider: %v", err))
		}
	}
}
