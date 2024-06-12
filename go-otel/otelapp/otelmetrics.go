package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter   = otel.Meter("rai.experiment.session")
	annotateCnt metric.Int64Counter
)

func init() {
	var err error
	annotateCnt, err = meter.Int64Counter("rai.experiment.session",
		metric.WithDescription("The number of annotate call."),
		metric.WithUnit("{annotate}"))
	if err != nil {
		panic(err)
	}
}