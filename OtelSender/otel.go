package main

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

var AZURE_TENANT_ID = "33e01921-4d64-4f8c-a055-5bdaffd5e33d"

var AZURE_CLIENT_ID = "9c7ae59d-9323-4423-a0da-38ddce774875"

// var AZURE_CLIENT_SECRET Cannot check in secret to Github repo

// The public otel collector scope.
var SCOPE = "fa65c9d1-e75e-4ac1-b7c1-608189fd7969/.default"

// The RAIUAI identity ID in AME
var UAI_CLIENT_ID = "1baa67a6-59c1-4c0f-a675-ee2682793b42"



func newTraceExporter(ctx context.Context) (trace.SpanExporter, error) {
	token, err := getToken(ctx, SCOPE, UAI_CLIENT_ID)
	if err != nil {
		return nil, err
	}

	kv := make(map[string]string)
	kv["Authorization"] = "Bearer " + token

	return otlptracehttp.New(ctx, otlptracehttp.WithHeaders(kv))
}

func newMetricsExporter(ctx context.Context) (metric.Exporter, error) {
	token, err := getToken(ctx, SCOPE, UAI_CLIENT_ID)
	if err != nil {
		return nil, err
	}

	kv := make(map[string]string)
	kv["Authorization"] = "Bearer " + token

	return otlpmetrichttp.New(ctx, otlpmetrichttp.WithHeaders(kv))
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider() (*trace.TracerProvider, error) {
	//traceExporter, err := stdouttrace.New(
	//	stdouttrace.WithPrettyPrint())
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("StubRAIService"),
		semconv.ServiceNamespaceKey.String("StubRAIOTrace"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("StubRAIOInstance"),
	)

	traceExporter, err := newTraceExporter(context.Background())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("StubRAIService"),
		semconv.ServiceNamespaceKey.String("StubRAIOMetrics"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("StubRAIOInstance"),
	)

	metricExporter, err := newMetricsExporter(context.Background())
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
		metric.WithResource(res),
	)
	return meterProvider, nil
}
