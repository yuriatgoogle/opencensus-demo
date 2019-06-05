package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	trace "go.opencensus.io/trace"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"

	"github.com/gorilla/mux"
)

var (
	projectID = "thegrinch-project"
)

// make an outbound call with context
func callBackend(ctx context.Context) string {
	_, span := trace.StartSpan(ctx, "call to backend")
	defer span.End()

	// make backend call with context
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1000*time.Millisecond)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%v\n", res.StatusCode)

	return strconv.Itoa(res.StatusCode)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(context.Background(), "incoming call")
	defer span.End()

	returnCode := callBackend(ctx)
	fmt.Fprintf(w, returnCode)
}

func main() {
	// set up Stackdriver exporter
	exporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectID, Location: "us-west1-a"})
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
		Propagation: &b3.HTTPFormat{}}

	log.Fatal(http.ListenAndServe(":8081", handler))
}
