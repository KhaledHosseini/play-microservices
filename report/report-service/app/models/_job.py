from typing import Any
from dataclasses import dataclass

from app.models import JsonObject

@dataclass
class Job(JsonObject):
    jobId: str
    name: str
    scheduleTime: str
    createdAt: str
    updatedAt: str
    status: int
    jobData: str

    @staticmethod
    def from_dict(obj: Any) -> 'Job':
        _jobId = str(obj.get("jobId"))
        _name = str(obj.get("name"))
        _scheduleTime = str(obj.get("scheduleTime"))
        _createdAt = str(obj.get("createdAt"))
        _updatedAt = str(obj.get("updatedAt"))
        _status = int(obj.get("status"))
        _jobData = str(obj.get("jobData"))
        return Job(_jobId, _name, _scheduleTime, _createdAt, _updatedAt, _status, _jobData)
