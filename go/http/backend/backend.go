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

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID, Location: "us-west1-a"})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	// trace.SetDefaultSampler(trace.AlwaysSample())

	// // Automatically add a Stackdriver trace header to outgoing requests:
	// client := &http.Client{
	// 	Transport: &ochttp.Transport{
	// 		Propagation: &propagation.HTTPFormat{},
	// 	},
	// }
	// _ = client // use client

	_, span := trace.StartSpan(context.Background(), "main")
	defer span.End()

	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
