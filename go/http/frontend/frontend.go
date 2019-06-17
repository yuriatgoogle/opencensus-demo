package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"

	"github.com/gorilla/mux"
)

var (
	projectID   = os.Getenv("PROJECT_ID")
	backendAddr = os.Getenv("BACKEND")
	location    = os.Getenv("LOCATION")
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// create root span
	ctx, rootspan := trace.StartSpan(context.Background(), "incoming call")
	defer rootspan.End()

	// create child span for backend call
	ctx, childspan := trace.StartSpan(ctx, "call to backend")
	defer childspan.End()

	// create request for backend call
	req, err := http.NewRequest("GET", backendAddr, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	childCtx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	req = req.WithContext(childCtx)

	// add span context to backend call and make request
	format := &tracecontext.HTTPFormat{}
	format.SpanContextToRequest(childspan.SpanContext(), req)
	//format.SpanContextToRequest(rootspan.SpanContext(), req)
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%v\n", res.StatusCode)
}

func main() {
	// set up Stackdriver exporter
	exporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectID, Location: location})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	// handle root request
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	var handler http.Handler = r
	handler = &ochttp.Handler{
		Handler:     handler,
		Propagation: &tracecontext.HTTPFormat{}}

	log.Fatal(http.ListenAndServe(":8081", handler))
}
