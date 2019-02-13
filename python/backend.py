# service = grpcBackend
 # method = returnValue
 # message = grpcMessage
 # field = payload

import grpc
from concurrent import futures
import time
import random
import datetime
import logging

# import the generated classes
import grpcDemo_pb2
import grpcDemo_pb2_grpc

# import classes for tracing
from opencensus.trace.samplers import always_on
from opencensus.trace.tracer import Tracer
from opencensus.trace.ext.grpc import server_interceptor

# import classes for exporting traces to Stackdriver
import os
import opencensus.common.transports
from opencensus.common.transports import async_
from opencensus.trace.exporters import stackdriver_exporter

# create a class to define the server functions
class grpcDemoServicer(grpcDemo_pb2_grpc.grpcBackendServicer):

    def returnValue(self, request, context):
        # create trace span
        tracer = Tracer(sampler=always_on.AlwaysOnSampler())
        with tracer.span(name='return value') as span:
            # delay for 0-5 seconds
            timeToBlock = random.randint(0,5)
            startTime = datetime.datetime.now()
            endTime = startTime + datetime.timedelta(seconds = timeToBlock)
            result = 0
            response = grpcDemo_pb2.grpcMessage()
            response.payload = request.payload
            response.payload = 'Delayed for ' + str(timeToBlock) + ' seconds'
            span.add_annotation('returning data', len=len(request.payload))
            while True: 
                result += random.random() * random.random()
                if (datetime.datetime.now() > endTime):
                    logging.info("exiting loop")
                    return response

# Setup the gRPC integration/interceptor
tracer_interceptor = server_interceptor.OpenCensusServerInterceptor(
        always_on.AlwaysOnSampler())

# set up the Stackdriver export
exporter = stackdriver_exporter.StackdriverExporter(
    # project_id=os.environ.get('YOUR_GOOGLE_PROJECT_ID_HERE'), # TODO = get this from env
    project_id='ygrinshteyn-sandbox',
    transport=async_.AsyncTransport)

# create a gRPC server with the interceptor
server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
            interceptors=(tracer_interceptor,))

# add defined class to server
grpcDemo_pb2_grpc.add_grpcBackendServicer_to_server(
    grpcDemoServicer(), server)

# start the Stackdriver exporter
tracer_interceptor = server_interceptor.OpenCensusServerInterceptor(
    always_on.AlwaysOnSampler())

# listen on port 50051
print('Starting server. Listening on port 50051.')
server.add_insecure_port('[::]:50051')
server.start()

# since server.start() will not block,
# a sleep-loop is added to keep alive
try:
    while True:
        time.sleep(86400)
except KeyboardInterrupt:
    server.stop(0)