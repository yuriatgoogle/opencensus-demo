# service = grpcBackend
 # method = returnValue
 # message = grpcMessage
 # field = payload

import grpc
from concurrent import futures
import time

# import the generated classes
import grpcDemo_pb2
import grpcDemo_pb2_grpc

# create a class to define the server functions
# derived from calculator_pb2_grpc.CalculatorServicer
class grpcDemoServicer(grpcDemo_pb2_grpc.grpcBackendServicer):

    def returnValue(self, request, context):
        # response = request.payload
        response = grpcDemo_pb2.grpcMessage()
        response.payload = request.payload
        return response

# create a gRPC server
server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

# add defined class to server
grpcDemo_pb2_grpc.add_grpcBackendServicer_to_server(
    grpcDemoServicer(), server)

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