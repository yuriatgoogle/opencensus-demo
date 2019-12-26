/*
Copyright 2019 Google LLC
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

// [START imports]
import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver"
	logs "github.com/GoogleCloudPlatform/opencensus-spanner-demo/applog"
	trace "go.opencensus.io/trace"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"

	"github.com/gorilla/mux"
)

// [END imports]

var (
	projectID = os.Getenv("PROJECT_ID")
	location  = os.Getenv("LOCATION")
	client    *logging.Client
)

// [START mainHandler]
func mainHandler(w http.ResponseWriter, r *http.Request) {
	// get context from incoming request
	ctx := r.Context()
	// get span context from incoming request
	HTTPFormat := &tracecontext.HTTPFormat{}
	if spanContext, ok := HTTPFormat.SpanContextFromRequest(r); ok {
		// create new span
		_, span := trace.StartSpanWithRemoteParent(ctx, "execute backend logic", spanContext)
		defer span.End()

		// generate a random 0-10 int and sleep for that many seconds
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(10) // n will be between 0 and 10
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Printf("Slept for %d seconds...\n", n)      // to console
		fmt.Fprintf(w, strconv.FormatInt(int64(n), 10)) // to client/browser

		// create a new context using the span and request context
		c := trace.NewContext(ctx, span)

		// create log entry with trace ID
		logs.Printf(c, "The backend process took %d seconds\n", n)
	}
}

// [END mainHandler]
// [START mainFunction]
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

	// set up logging
	logs.Initialize(projectID)
	defer logs.Close()

	// handle incoming request
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	var handler http.Handler = r

	handler = &ochttp.Handler{
		Handler:     handler,
		Propagation: &tracecontext.HTTPFormat{}}

	log.Fatal(http.ListenAndServe(":8080", handler))
}

// [END mainFunction]
