from app.proto import ReportGRPCTypes
from typing import Any
from dataclasses import dataclass
from google.protobuf.timestamp_pb2 import Timestamp
from datetime import datetime
from app.models import JsonObject


@dataclass
class Report(JsonObject):
    type: int
    topic: str
    created_time: str
    report_data: str

    def __init__(self, type: int, topic: str, created_time: datetime, report_data: str):
        self.type = type
        self.topic = topic
        self.created_time = created_time
        self.report_data = report_data

    @staticmethod
    def from_dict(obj: Any) -> "Report":
        _type = int(obj.get("type"))
        _topic = str(obj.get("topic"))
        _created_time = obj.get("created_time")
        _report_data = str(obj.get("report_data"))
        return Report(_type, _topic, _created_time, _report_data)

    def saveToDB(self, db: "ReportDB"):
        db.create(self)

    def toProto(self):
        ts = Timestamp()
        ts.FromDatetime(self.created_time)
        return ReportGRPCTypes.Report(
            type=self.type,
            topic=self.topic,
            created_time=ts,
            report_data=self.report_data,
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
            totalCount=self.TotalCount,
            totalPages=self.TotalPages,
            page=self.Page,
            size=self.Size,
            hasMore=self.HasMore,
            reports=[r.toProto() for r in self.Reports],
        )


class ReportDB:
    def create(self, report: Report) -> Report:
        raise Exception("create is not implemented")

    def list(self, type: int, page: int, size: int) -> ReportList:
        raise Exception("list is not implemented")
