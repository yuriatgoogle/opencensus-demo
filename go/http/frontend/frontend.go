package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"contrib.go.opencensus.io/exporter/stackdriver"
	trace "go.opencensus.io/trace"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"

	"github.com/gorilla/mux"
)

var (
	projectID = "thegrinch-project"
)

// make an outbound call and
func callBackend(ctx context.Context) string {
	_, span := trace.StartSpan(ctx, "backendCall")
	defer span.End()
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		log.Fatal("could not fetch Google")
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
	ctx, span := trace.StartSpan(context.Background(), "root call")
	defer span.End()

	returnCode := callBackend(ctx)
	fmt.Fprintf(w, returnCode)
}

func main() {

	exporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectID, Location: "us-west1-a"})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	var handler http.Handler = r
	// handler = &logHandler{log: log, next: handler}

	handler = &ochttp.Handler{
		Handler:     handler,
		Propagation: &b3.HTTPFormat{}}

	log.Fatal(http.ListenAndServe(":8081", handler))
}
