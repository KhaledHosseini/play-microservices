import logging

from typing import Any
from dataclasses import dataclass
from app.models import JsonObject
from enum import IntEnum


class JobStatusEnum(IntEnum):
    JOB_STATUS_UNKNOWN = 0
    JOB_STATUS_PENDING = 1
    JOB_STATUS_SCHEDULED = 2
    JOB_STATUS_RUNNING = 3
    JOB_STATUS_SUCCEEDED = 4
    JOB_STATUS_FAILED = 5


@dataclass
class EmailJob(JsonObject):
    Id: str
    Status: JobStatusEnum
    JobData: str

    @staticmethod
    def from_dict(obj: Any) -> 'EmailJob':
        try:
            values = {
                'JobID': obj.get("Id"),
                'Status': JobStatusEnum(obj.get("Status")),
                'JobData': obj.get("JobData"),
            }
            logging.info(values.values())
            if any(value is None for value in values.values()):
                logging.info("None is found")
                raise ValueError("Invalid EmailJob json")
            return EmailJob(**values)
        except Exception as e:
            raise e
    def get_email(self):
        return Email.from_dict(self.JobData)

@dataclass
class Email(JsonObject):
    SourceAddress: str
    DestinationAddress: str
    Subject: str
    Message: str

    @staticmethod
    def from_dict(obj: Any) -> 'Email':
        try:
            values = {
                'SourceAddress': obj.get("SourceAddress"),
                'DestinationAddress': obj.get("DestinationAddress"),
                'Subject': obj.get("Subject"),
                'Message': obj.get("Message"),
            }
            logging.info("Email values are:")
            logging.info(values)
            if any(value is None for value in values.values()):
                logging.info("None is found")
                raise ValueError("Invalid Email json")
            return Email(**values)
        except Exception as e:
            raise e
