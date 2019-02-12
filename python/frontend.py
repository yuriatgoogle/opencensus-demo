 # service = grpcBackend
 # method = returnValue
 # message = grpcMessage
 # field = payload

import grpc

# import the generated classes
import grpcDemo_pb2
import grpcDemo_pb2_grpc

# open a gRPC channel
channel = grpc.insecure_channel('localhost:50051')

# create a stub (client)
stub = grpcDemo_pb2_grpc.grpcBackendStub(channel)
# stub = calculator_pb2_grpc.CalculatorStub(channel)

# create a valid request message
# message to send = pb2.message(field name = field value)
sendString = grpcDemo_pb2.grpcMessage(payload='test')
# number = calculator_pb2.Number(value=16)

# make the call
# response = stub.SquareRoot(number)
# response = stub.method(message to send)
response = stub.returnValue(sendString)

# return
print(response.payload)