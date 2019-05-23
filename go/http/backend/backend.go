package main

import (
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
func callGoogle() string {
	resp, err := http.Get("https://www.google.com")
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
	ctx := r.Context()
	_, span := trace.StartSpan(ctx, "call Google")
	defer span.End()
	returnCode := callGoogle()
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

	log.Fatal(http.ListenAndServe(":8080", handler))
}
