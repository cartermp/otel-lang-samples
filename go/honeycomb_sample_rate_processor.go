package main

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type honeycombSamplerSpanProcessor struct {
	sampleRate float64
}

func (h *honeycombSamplerSpanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
	// Honeycomb accepts a whole number N where 1/N corresponds to the given "sampleRate" in the struct.
	// This isn't perfect at all - yay floating point math - but it seems to work.
	// Moreover, this just documents how you might process spans like this.
	sampleRateForHNY := int(1.0 / h.sampleRate)
	s.SetAttributes(attribute.Int("sampleRate", sampleRateForHNY))
}

func (h *honeycombSamplerSpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	// Do nothing, since we've already added the attribute to the span via OnStart
}

func (h *honeycombSamplerSpanProcessor) Shutdown(ctx context.Context) error {
	// N/A
	return nil
}

func (h *honeycombSamplerSpanProcessor) ForceFlush(ctx context.Context) error {
	// N/A
	return nil
}
