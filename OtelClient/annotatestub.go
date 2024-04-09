package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer  = otel.Tracer("annotate")
	meter   = otel.Meter("rai.experiment.annotate")
	annotateCnt metric.Int64Counter
)

func init() {
	var err error
	annotateCnt, err = meter.Int64Counter("rai.annotate",
		metric.WithDescription("The number of annotate call."),
		metric.WithUnit("{annotate}"))
	if err != nil {
		panic(err)
	}
}

func annotate(w http.ResponseWriter, r *http.Request) {
	// Read the request body  
	body, err := io.ReadAll(r.Body)  
	if err != nil {  
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)  
		return  
	}  
	
	// Get the request body  
	bodyStr := string(body)

	ctx, span := tracer.Start(r.Context(), "annotate")
	defer span.End()

	// random annotate result
	annotateResult := bodyStr + ": with random number " + strconv.Itoa(rand.Intn(100))

	// Add the custom attribute to the span and counter.
	annotateResultValueAttr := attribute.String("annotate.result.value", annotateResult)
	annotateFakeLatency := attribute.Int("annotate.result.remote.tms.latency", rand.Intn(1000))
	span.SetAttributes(annotateResultValueAttr)
	annotateCnt.Add(ctx, 1, metric.WithAttributes(annotateResultValueAttr, annotateFakeLatency))

	resp := annotateResult + "\n"
	print(annotateResult)
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

// liveness endpoint just returns 200
func liveness(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("OK"))
}

// readiness endpoint just returns 200
func readiness(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("OK"))
}
