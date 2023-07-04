from proto import ReportGRPC
import logging

from app.models import ReportDBInterface

class MyReportService(ReportGRPC.ReportService):
    def __init__(self, reportDB: ReportDBInterface):
        self.reportDB = reportDB
        super().__init__()
    
    def ListReports(self,request, context):
        logging.info(f"message recieved...{request}")
        filter = request.filter
        page = request.page
        size = request.size
        result = self.reportDB.list(type=filter,page=page,size=size)
        return result.toProto()