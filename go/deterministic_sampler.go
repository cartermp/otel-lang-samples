package main

import (
	"errors"
	"fmt"

	// "fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	ErrInvalidSampleRate = errors.New("sample rate must be >= 1")
)

type deterministicSampler struct {
	sampleRate int
	sampler    sdktrace.Sampler
}

func DeterministicAlwaysSampleSampler(sampleRate int) (*deterministicSampler, error) {
	return &deterministicSampler{
		sampleRate: 1,
		sampler:    sdktrace.AlwaysSample(),
	}, nil
}

func (ds *deterministicSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	for _, attr := range p.Attributes {
		fmt.Println(attr)
	}

	delegatedResult := ds.sampler.ShouldSample(p)

	return sdktrace.SamplingResult{
		Decision:   delegatedResult.Decision,
		Attributes: delegatedResult.Attributes,
		Tracestate: delegatedResult.Tracestate,
	}
}

func (ds *deterministicSampler) Description() string {
	return "HoneycombDeterministicSampler"
}
