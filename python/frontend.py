import grpc

# import the generated classes
import basicmessage_pb2
import basicmessage_pb2_grpc

# open a gRPC channel
channel = grpc.insecure_channel('localhost:55001')

# create a stub (client)
backend = basicmessage_pb2_grpc.backendStub(channel)
message = basicmessage_pb2.commString(payload='test')

response = backend.sendMessage(message)

print(response.value)