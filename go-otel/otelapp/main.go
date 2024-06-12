package main

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"rai-go-otel/otelsetup"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	timestamp := int64(1714203347)  
	t := time.Unix(timestamp, 0)  
  
	formattedTime := t.Format("01-02-2006 15:04:05")  
  
	println(formattedTime) 

	// Set up OpenTelemetry.
	otelShutdown, err := otelsetup.SetupOTelSDK(ctx)
	if err != nil {
		println(err.Error())
		return
	}

	// For testing purpose, run a infinite loop to emit metrics
	go func() {
		// testCnt, _ := meter.Int64Counter("rai.experiment.session",
		// 	metric.WithDescription("The number of annotate call."),
		// 	metric.WithUnit("{session}"))
		println("FC!")
		testCnt, _ := meter.Int64Counter("rai.experiment.eus2",
			metric.WithDescription("The latency of segment."),
			metric.WithUnit("{ms}"))

		// ctx := context.Background()

		for {
			num := rand.Float64() * 99.999 
			attr1 := attribute.String("segmenter", "FixedValue")
			attr2 := attribute.Float64("latency", num)
			//span.SetAttributes(annotateResultValueAttr)
			testCnt.Add(context.Background(), 1, metric.WithAttributes(attr1, attr2))
			//testHst.Record(ctx, num, metric.WithAttributes(attr1))
			time.Sleep(1 * time.Second)

		}
	}()
	
	// update the token every 5 mins.
	// This is to test if token is correctly set/removed.
	go func() {
		auth := true
		for {
			otelsetup.UpdateMetricExporterAuthToken(auth, context.Background())
			auth = !auth
			println("Sleep for 2 hours")
			time.Sleep(7200 * time.Second)
			println("Sleep done")
		}
	}()

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Wait for the context to be done (cancelled)  
	<-ctx.Done()

	// Perform cleanup or other necessary actions  
	stop()
}