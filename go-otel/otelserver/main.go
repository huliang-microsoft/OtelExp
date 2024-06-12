package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"rai-go-otel/otelsetup"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := otelsetup.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":12345",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// update the token every 5 mins.
	// This is to test if token is correctly set/removed.
	go func() {
		auth := true
		for {
			otelsetup.UpdateMetricExporterAuthToken(auth, context.Background())
			auth = !auth
			println("Sleep for 300s")
			time.Sleep(300 * time.Second)
			println("Sleep done")
		}

	}()

	// For testing purpose, run a infinite loop to emit metrics
	go func() {
		testCnt, _ := meter.Int64Counter("rai.experiment.annotate",
		metric.WithDescription("The number of annotate call."),
		metric.WithUnit("{annotate}"))
		for {
			attr1 := attribute.String("annotate.result.value", "123")
			attr2 := attribute.Int("annotate.result.remote.tms.latency", 100)
			//span.SetAttributes(annotateResultValueAttr)
			testCnt.Add(context.Background(), 1, metric.WithAttributes(attr1, attr2))
			time.Sleep(3 * time.Second)
		}
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// Register handlers.
	handleFunc("/annotate", annotate)
	handleFunc("/api/health", liveness)
	handleFunc("/readiness", readiness)

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}