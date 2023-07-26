from proto import ReportGRPCTypes
from typing import Any
from dataclasses import dataclass
from google.protobuf.timestamp_pb2 import Timestamp
from datetime import datetime
from app.models import JsonObject


@dataclass
class Report(JsonObject):
    Topic: str
    CreatedTime: str
    ReportData: str

    @staticmethod
    def from_dict(obj: Any) -> "Report":
        _topic = str(obj.get("Topic"))
        _created_time = obj.get("CreatedTime")
        _report_data = str(obj.get("ReportData"))
        return Report(_topic, _created_time, _report_data)

    def saveToDB(self, db: "ReportDBInterface"):
        db.create(self)

    def toProto(self):
        ts = Timestamp()
        ts.FromDatetime(self.CreatedTime)
        return ReportGRPCTypes.Report(
            topic=self.Topic,
            created_time=ts,
            report_data=self.ReportData,
        )


class ReportList:
    def __init__(
        self, totalCount=0, totalPages=0, page=0, size=0, hasMore=False, reports=[]
    ) -> None:
        self.TotalCount = totalCount
        self.TotalPages = totalPages
        self.Page = page
        self.Size = size
        self.HasMore = hasMore
        self.Reports = reports

    def toProto(self):
        return ReportGRPCTypes.ListReportResponse(
            total_count=self.TotalCount,
            total_pages=self.TotalPages,
            page=self.Page,
            size=self.Size,
            has_more=self.HasMore,
            reports=[r.toProto() for r in self.Reports],
        )


class ReportDBInterface:
    def create(self, report: Report) -> Report:
        raise Exception("create is not implemented")

    def list(self, page: int, size: int) -> ReportList:
        raise Exception("list is not implemented")
