package main

import (
	"errors"
	"fmt"
	// "fmt"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	ErrInvalidSampleRate = errors.New("sample rate must be >= 1")
)

type deterministicSampler struct {
	sampleRate          int
	traceIDRatioSampler sdktrace.Sampler
}

func DeterministicSampler(sampleRate int) (*deterministicSampler, error) {
	if sampleRate < 1 {
		return nil, ErrInvalidSampleRate
	}

	return &deterministicSampler{
		sampleRate:          sampleRate,
		traceIDRatioSampler: sdktrace.TraceIDRatioBased(1.0 / float64(sampleRate)),
	}, nil
}

func (ds *deterministicSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	attrs := []attribute.KeyValue{
		attribute.Int("SampleRate", int(ds.sampleRate)),
	}

	for _, attr := range p.Attributes {
		fmt.Println(attr)
	}

	delegatedResult := ds.traceIDRatioSampler.ShouldSample(p)

	return sdktrace.SamplingResult{
		Decision:   delegatedResult.Decision,
		Attributes: attrs,
		Tracestate: delegatedResult.Tracestate,
	}
}

func (ds *deterministicSampler) Description() string {
	return "HoneycombDeterministicSampler"
}
