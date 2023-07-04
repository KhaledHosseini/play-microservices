from datetime import datetime
import logging
from typing import List

from app.message_engines import KafkaConsumerWorker

from app.models import ReportDB
from app.models import Job
from app.models import Report

class JobConsumerWorker(KafkaConsumerWorker):
    def __init__(self, topic: str, consumerGroupId: str, kafkaBrokers: List[str], reportDB: ReportDB):
        super().__init__(topic = topic,consumerGroupId = consumerGroupId,kafkaBrokers=kafkaBrokers)
        self.reportDB = reportDB
    def messageRecieved(self,topic: str, message: any):
        try:
            job = Job.from_dict(message)
            report = Report(
                type=0,
                topic=topic,
                created_time=datetime.now(),
                report_data=job.toJsonStr()
            )
            self.reportDB.create(report=report)
        except Exception as e:
            logging.error("Unable to load json for job ", e)