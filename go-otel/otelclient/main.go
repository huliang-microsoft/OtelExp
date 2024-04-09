package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"rai-go-otel/otelsetup"
)

var (
	tracer  = otel.Tracer("annotate")
	meter   = otel.Meter("rai.experiment.annotate")
	annotateCnt metric.Int64Counter
)

func main() {  
	// Create a ticker that triggers every 10 seconds  
	ticker := time.NewTicker(10 * time.Second)

	// Set up OpenTelemetry.
	otelShutdown, err := otelsetup.SetupOTelSDK(context.Background())
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Create an HTTP client with OpenTelemetry instrumentation  
	client := http.DefaultClient  
	client.Transport = otelhttp.NewTransport(client.Transport)  
  
	// Define the URL to send the request to  
	url := "http://localhost:12345/annotate"  
  
	// Start a goroutine to handle the periodic requests  
	go func() {  
		for {  
			select {  
			case <-ticker.C:  
				// Send the HTTP GET request  
				resp, err := client.Get(url)  
				if err != nil {  
					fmt.Printf("Error sending request: %s\n", err)  
					continue  
				}
				annotateResultValueAttr := attribute.String("annotate.result.exp.value", "TestStringValue")
				annotateFakeLatency := attribute.Int("annotate.result.exp.tms.latency", rand.Intn(1000))
				annotateCnt.Add(context.Background(), 1, metric.WithAttributes(annotateResultValueAttr, annotateFakeLatency))
  
				// Print the response status code  
				fmt.Printf("Response status: %s\n", resp.Status)  
				resp.Body.Close()  
			}  
		}  
	}()
  
	// Block the main goroutine to keep the program running  
	select {}  
}