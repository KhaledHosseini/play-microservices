from concurrent import futures
import logging

import grpc
from grpc_reflection.v1alpha import reflection

from proto import ReportGRPC
from proto import ReportGRPCTypes

from app.grpc_interceptors import TokenValidationInterceptor
from config import Config

from app.models.report import MyReportService
from app.models.report import ReportDBMongo
from app.models.job import ConsumerService

class Server:
    def run(self, cfg: Config):
        logging.info("Running services...")

        db = ReportDBMongo(cfg)
       
        consumer = ConsumerService(cfg,db)
        consumer.start()
        
        my_report_grpc_service = MyReportService(reportDB=db)
        server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
            interceptors = (TokenValidationInterceptor(cfg),)
        )
        ReportGRPC.add_ReportServiceServicer_to_server(my_report_grpc_service, server)
        SERVICE_NAMES = (
            ReportGRPCTypes.DESCRIPTOR.services_by_name["ReportService"].full_name,
            reflection.SERVICE_NAME,
        )

        reflection.enable_server_reflection(SERVICE_NAMES, server)

        port = cfg.ServerPort
        server.add_insecure_port("[::]:" + port)
        server.start()
        server.wait_for_termination()