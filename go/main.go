package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
			"x-honeycomb-team": "lol no u",
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
			semconv.ServiceNameKey.String("phillips-happy-fun-time"), // lol no generics
		)

	sampler, err := DeterministicSampler(1)
	if err != nil {
		panic(err) // idk lol
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sampler),
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

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("isTrue", true), attribute.String("stringAttr", "hi!"))

	fib := func(ctx context.Context) uint {
		_, span := tracer.Start(ctx, "fib")
		defer span.End()

		span.SetAttributes(attribute.Bool("isTrue", false), attribute.String("stringAttr", "bye!"))

		return fib(2)
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

	// Set the Tracer Provider and the W3C Trace Context propagator as globals. Important, otherwise this won't let you see distributed traces be connected!
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

	fmt.Println("doing it")

	log.Fatal(http.ListenAndServe(":3030", nil))
}
