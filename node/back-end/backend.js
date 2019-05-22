/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

var PROTO_PATH = __dirname + '/../proto/helloworld.proto'
var grpc = require('grpc');
var hello_proto = grpc.load(PROTO_PATH).helloworld;
const request = require('request');

// set up tracing
var tracing = require('@opencensus/nodejs');
var stackdriver = require('@opencensus/exporter-stackdriver');
var projectID = 'thegrinch-project'; //TODO - convert to env variable

// create and start Stackdriver exporter
var exporter = new stackdriver.StackdriverTraceExporter({projectId: projectID});
tracing.registerExporter(exporter).start();
const tracer = tracing.start({samplingRate: 1}).tracer;
 


// implements an outbound HTTP call
function callOutbound (url) {
  request(url, function (error, response, body) {
    console.log('error:', error); 
    console.log('statusCode:', response && response.statusCode); 
  });
}

// Implements the SayHello RPC method - with tracing
function sayHello(call, callback) {
  tracer.startRootSpan({name: 'backendMain'}, rootSpan => {
    // code to be traced goes in here:
    rootSpan.addAnnotation('main span');
    const childSpan = tracer.startChildSpan('google_call');
    childSpan.start();
    callOutbound('https://www.google.com');
    callback(null, {message: 'Google.com loaded in backend'});
    childSpan.end();
    rootSpan.end();
  });
}


 // Start RPC server
function main() {
  var server = new grpc.Server();
  server.addService(hello_proto.Greeter.service, {sayHello: sayHello});
  server.bind('0.0.0.0:' + process.env.BACKENDPORT, grpc.ServerCredentials.createInsecure());
  server.start();
}

main();