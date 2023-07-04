# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import proto.report_pb2 as report__pb2


class ReportServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.ListReports = channel.unary_unary(
                '/reportService.ReportService/ListReports',
                request_serializer=report__pb2.ListReportsRequest.SerializeToString,
                response_deserializer=report__pb2.ListReportResponse.FromString,
                )


class ReportServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def ListReports(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ReportServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'ListReports': grpc.unary_unary_rpc_method_handler(
                    servicer.ListReports,
                    request_deserializer=report__pb2.ListReportsRequest.FromString,
                    response_serializer=report__pb2.ListReportResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'reportService.ReportService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ReportService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def ListReports(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/reportService.ReportService/ListReports',
            report__pb2.ListReportsRequest.SerializeToString,
            report__pb2.ListReportResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
