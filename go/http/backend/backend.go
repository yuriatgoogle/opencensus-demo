package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"contrib.go.opencensus.io/exporter/stackdriver"
	trace "go.opencensus.io/trace"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"

	"github.com/gorilla/mux"
)

var (
	projectID = os.Getenv("PROJECT_ID")
	destURL   = os.Getenv("DESTINATION_URL")
	location  = os.Getenv("LOCATION")
)

// make an outbound call
func callRemoteEndpoint() string {
	resp, err := http.Get(destURL)
	if err != nil {
		log.Fatal("could not fetch remote endpoint")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("could not read response from Google")
		log.Fatal(body)
	}

	return strconv.Itoa(resp.StatusCode)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// get context from incoming request
	ctx := r.Context()
	// get span context from incoming request
	HTTPFormat := &tracecontext.HTTPFormat{}
	if spanContext, ok := HTTPFormat.SpanContextFromRequest(r); ok {
		_, span := trace.StartSpanWithRemoteParent(ctx, "call remote endpoint", spanContext)
		defer span.End()
		returnCode := callRemoteEndpoint()
		fmt.Fprintf(w, returnCode)
	}
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

	// handle incoming request
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	var handler http.Handler = r
	// handler = &logHandler{log: log, next: handler}

	handler = &ochttp.Handler{
		Handler:     handler,
		Propagation: &tracecontext.HTTPFormat{}}

	log.Fatal(http.ListenAndServe(":8080", handler))
}
