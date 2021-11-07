# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import multiplayerservice_pb2 as multiplayerservice__pb2


class MultiplayerServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Login = channel.stream_stream(
                '/proto.MultiplayerService/Login',
                request_serializer=multiplayerservice__pb2.ClientMessage.SerializeToString,
                response_deserializer=multiplayerservice__pb2.ServerMessage.FromString,
                )


class MultiplayerServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Login(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_MultiplayerServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Login': grpc.stream_stream_rpc_method_handler(
                    servicer.Login,
                    request_deserializer=multiplayerservice__pb2.ClientMessage.FromString,
                    response_serializer=multiplayerservice__pb2.ServerMessage.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'proto.MultiplayerService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class MultiplayerService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Login(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(request_iterator, target, '/proto.MultiplayerService/Login',
            multiplayerservice__pb2.ClientMessage.SerializeToString,
            multiplayerservice__pb2.ServerMessage.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
