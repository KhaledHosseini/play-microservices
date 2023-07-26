from typing import Any
from dataclasses import dataclass

from app.models import JsonObject

@dataclass
class Job(JsonObject):
    Id: str
    Name: str
    Description: str
    ScheduleTime: str
    CreatedAt: str
    UpdatedAt: str
    Status: int
    JobType: int
    JobData: str

    @staticmethod
    def from_dict(obj: Any) -> 'Job':
        _Id = str(obj.get("Id"))
        _name = str(obj.get("Name"))
        _description=str(obj.get("Description"))
        _scheduleTime = str(obj.get("ScheduleTime"))
        _createdAt = str(obj.get("CreatedAt"))
        _updatedAt = str(obj.get("UpdatedAt"))
        _status = int(obj.get("Status"))
        _JobType = int(obj.get("JobType"))
        _jobData = str(obj.get("JobData"))
        return Job(_Id, _name,_description, _scheduleTime, _createdAt, _updatedAt, _status,_JobType, _jobData)
