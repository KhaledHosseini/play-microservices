import json
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
    job_id: str
    status: JobStatusEnum
    job_data: str

    @staticmethod
    def from_dict(obj: Any) -> 'EmailJob':
        try:
            values = {
                'job_id': obj.get("job_id"),
                'status': JobStatusEnum(obj.get("status")),
                'job_data': obj.get("job_data"),
            }
            logging.info(values.values())
            if any(value is None for value in values.values()):
                logging.info("Non is found")
                raise ValueError("Invalid EmailJob json")
            return EmailJob(**values)
        except Exception as e:
            raise e
    def get_email(self):
        return Email.from_dict(self.job_data)

@dataclass
class Email(JsonObject):
    source_address: str
    destination_address: str
    subject: str
    message: str

    @staticmethod
    def from_dict(obj: Any) -> 'Email':
        try:
            values = {
                'source_address': obj.get("source_address"),
                'destination_address': obj.get("destination_address"),
                'subject': obj.get("subject"),
                'message': obj.get("message"),
            }
            logging.info("Email values are:")
            logging.info(values)
            if any(value is None for value in values.values()):
                logging.info("None is found")
                raise ValueError("Invalid Email json")
            return Email(**values)
        except Exception as e:
            raise e


