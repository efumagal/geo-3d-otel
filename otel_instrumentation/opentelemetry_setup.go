package otel_instrumentation

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitializeGlobalTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, *otlptrace.Exporter, error) {
	// Configure a new OTLP exporter using environment variables for sending data to Honeycomb over gRPC
	clientOTel := otlptracegrpc.NewClient()
	exp, err := otlptrace.New(ctx, clientOTel)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %e", err)
	}

	// Create a new tracer provider with a batch span processor and the otlp exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Register the global Tracer provider
	otel.SetTracerProvider(tp)

	// Register the W3C trace context and baggage propagators so data is propagated across services/processes
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, exp, nil
}
