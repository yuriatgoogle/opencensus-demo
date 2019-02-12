import grpc
import time
import hashlib
import basicmessage_pb2
import basicmessage_pb2_grpc
from concurrent import futures
 
class backendServer(basicmessage_pb2_grpc.backendServicer):
    """
    gRPC server 
    """
    def __init__(self, *args, **kwargs):
        self.server_port = 55001
 
    def sendMessage(self, request, context):
        """
        Implementation of the rpc sendMessage declared in the proto
        file above.
        """ 
        return basicmessage_pb2.backendReply(request.backendRequest)
 
    def start_server(self):
        """
        Function which actually starts the gRPC server, and preps
        it for serving incoming connections
        """
        # declare a server object with desired number
        # of thread pool workers.
        backendServer = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
 
        # bind the server to the port defined above
        backendServer.add_insecure_port('[::]:{}'.format(self.server_port))
 
        # start the server
        backendServer.start()
        print ('Server running ...')
 
        try:
            # need an infinite loop since the above
            # code is non blocking, and if I don't do this
            # the program will exit
            while True:
                time.sleep(60*60*60)
        except KeyboardInterrupt:
            backendServer.stop(0)
            print('Digestor Server Stopped ...')
 
curr_server = backendServer()
curr_server.start_server()