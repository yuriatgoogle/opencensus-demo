 # service = grpcBackend
 # method = returnValue
 # message = grpcMessage
 # field = payload

import grpc
from flask import Flask

# import the generated classes
import grpcDemo_pb2
import grpcDemo_pb2_grpc

# import tracing classes
from opencensus.trace.samplers import always_on
from opencensus.trace.tracer import Tracer
from opencensus.trace.ext.grpc import server_interceptor

app = Flask(__name__)

@app.route('/')
def homePage():

    
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
    return (response.payload)

if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0', port=8080)

