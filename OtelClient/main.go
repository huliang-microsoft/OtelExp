package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func main() {  
	// Create a ticker that triggers every 10 seconds  
	ticker := time.NewTicker(10 * time.Second)

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(context.Background())
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
				annotateResultValueAttr := attribute.String("annotate.result.value", "NothingSpecial")
				annotateFakeLatency := attribute.Int("annotate.result.remote.tms.latency", rand.Intn(1000))
				annotateCnt.Add(context.Background(), 1, metric.WithAttributes(annotateResultValueAttr, annotateFakeLatency))
  
				// Print the response status code  
				fmt.Printf("Response status: %s\n", resp.Status)  
				resp.Body.Close()  
			}  
		}  
	}()

	go delayedExecution()
  
	// Block the main goroutine to keep the program running  
	select {}  
}

func delayedExecution() {
	i := 0
	for {  
		// Call your function here
		time.Sleep(24 * time.Hour)  
		meterProvider, err := newMeterProviderWithoutAuth()
		fmt.Println("Created " + strconv.Itoa(i) + "th  meter provider.")
		if err == nil {
			otel.SetMeterProvider(meterProvider)
			println("Set " +  strconv.Itoa(i) + "th meter provider.")
		}
		i = i + 1
  
		// Sleep for 24 hours  

	} 
}  