/*
#    Copyright 2017 Google Inc.
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        https://www.apache.org/licenses/LICENSE-2.0

#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
*/

// grpc stuff
var protobuf = require('protobufs');
var PROTO_PATH = __dirname + '/helloworld.proto'
console.log("proto path: " + PROTO_PATH)
var grpc = require('grpc');
var hello_proto = grpc.load(PROTO_PATH).helloworld;
console.log("hello_proto: " + hello_proto)


const express = require('express');
const app = express();

var tracing = require('@opencensus/nodejs');
var stackdriver = require('@opencensus/exporter-stackdriver');
var projectID = 'ygrinshteyn-sandbox';

// create Stackdriver exporter
var exporter = new stackdriver.StackdriverTraceExporter({projectId: projectID});

// start Stackdriver exporter
tracing.registerExporter(exporter).start();

// start tracing and set sampling rate
const tracer = tracing.start({samplingRate: 1}).tracer;


app.get('/', (req, res) => {

    console.log('Inbound request received!');
    // set grpc options
    var backendHost = process.env.BACKENDHOST;
    var backendPort = process.env.BACKENDPORT;
    var client = new hello_proto.Greeter(backendHost + ':' + backendPort, grpc.credentials.createInsecure());
    var user = 'Yuri'

    // create root span
    tracer.startRootSpan({name: 'main'}, rootSpan => {
        // code to be traced goes in here:
        rootSpan.addAnnotation('main span');
        // make grpc call
        const grpcCallSpan = tracer.startChildSpan('grpcCall');
        grpcCallSpan.start();
        client.sayHello({name: user}, function(err, response) {
            if (err){
                console.log("could not get grpc response");
                res.send("could not get grpc response");
                return;
            }
            console.log('Greeting:', response.message);
            grpcCallSpan.end();
            rootSpan.end();
            res.send("grpc response is " + response.message);
        });
    });
}); // end app.get

// Start the server
const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
  console.log("Server running at http://127.0.0.1:8080/");
  console.log('Press Ctrl+C to quit.');
});