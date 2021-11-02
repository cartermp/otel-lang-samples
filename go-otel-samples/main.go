package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	tracer trace.Tracer
)

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint("api.honeycomb.io:443"),
		otlptracegrpc.WithHeaders(map[string]string{
			"x-honeycomb-team":    os.Getenv("HONEYCOMB_TEAM"),
			"x-honeycomb-dataset": os.Getenv("HONEYCOMB_DATASET"),
		}),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}

	client := otlptracegrpc.NewClient(opts...)
	return otlptrace.New(ctx, client)
}

func newTraceProvider(exp *otlptrace.Exporter) *sdktrace.TracerProvider {
	resource :=
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ExampleService"), // lol no generics
		)

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
	)
}

func fib(n uint) uint {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fib := func(ctx context.Context) uint {
		_, span := tracer.Start(ctx, "fib")
		defer span.End()

		return fib(25)
	}(ctx)

	fmt.Fprintf(w, "%d", fib)
}

func main() {
	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := newTraceProvider(exp)

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()

	// Set the Tracer Provider and the W3C Trace Context propagator as globals
	// I have no idea if this is actually required or not
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	// Initialize the module-level tracer. One tracer per module!
	tracer = tp.Tracer("fib-http-service")

	// Wire up http auto-instrumentation
	handler := http.HandlerFunc(httpHandler)
	wrapedHandler := otelhttp.NewHandler(handler, "hello")
	http.Handle("/hello", wrapedHandler)

	log.Fatal(http.ListenAndServe(":3030", nil))
}
