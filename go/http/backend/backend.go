package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
)

var (
	projectID = "thegrinch=project"
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
	returnCode := callGoogle()
	fmt.Fprintf(w, returnCode)
}

func main() {

	handler = &ochttp.Handler{
		Handler:     handler,
		Propagation: &b3.HTTPFormat{}}

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:            projectID,
		BundleDelayThreshold: time.Second / 10,
		BundleCountThreshold: 10})
	if err != nil {
		log.Println(err)
	}
	trace.RegisterExporter(exporter)
	trace.SetDefaultSampler(trace.AlwaysSample())

	_, span := trace.StartSpan(context.Background(), "main")
	defer span.End()

	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
