from threading import Thread
import logging
from datetime import datetime
from typing import List

from config import Config

from app.message_engines import KafkaConsumerWorker

from app.models import ReportDBInterface
from app.models import ReportDBInterface
from app.models import Job
from app.models import Report

class ConsumerService:
  def __init__(self,cfg: Config, reportDb: ReportDBInterface):
      self.cfg = cfg
      self.reportDB = reportDb
      self.threads = []

  def start(self):
      self.runConsumer(topic = self.cfg.TopicJobRun,consumerGroupId = self.cfg.TopicJobRunConsumerGroupID,workerCount=self.cfg.TopicJobRunWorkerCount)
      self.runConsumer(topic=self.cfg.TopicJobRunResult,consumerGroupId=self.cfg.TopicJobRunResultConsumerGroupID,workerCount=self.cfg.TopicJobRunResultConsumerWorkerCount)
      # run other topics

      for t in self.threads:
        t.join()

  def runConsumer(self, topic: str, consumerGroupId: str, workerCount:int):
        logging.info(f"Starting email job consumers with  {workerCount } workers...")
        for i in range(0,workerCount):
            t = Thread(target=self.run_worker,args=(topic,consumerGroupId))
            t.Daemon = True
            self.threads.append(t)
            t.start()
            logging.info(f"Worker {i} started for consuming job events...")
  
  def run_worker(self, topic: str, consumerGroupId: str):
    consumer = JobConsumerWorker(topic = topic,
                           consumerGroupId = consumerGroupId,
                           kafkaBrokers=self.cfg.KafkaBrokers,
                           reportDB= self.reportDB)
    consumer.run()

class JobConsumerWorker(KafkaConsumerWorker):
    def __init__(
        self,
        topic: str,
        consumerGroupId: str,
        kafkaBrokers: List[str],
        reportDB: ReportDBInterface,
    ):
        super().__init__(
            topic=topic, consumerGroupId=consumerGroupId, kafkaBrokers=kafkaBrokers
        )
        self.reportDB = reportDB

    def messageRecieved(self, topic: str, message: any):
        try:
            job = Job.from_dict(message)
            report = Report(
                type=0,
                topic=topic,
                created_time=datetime.now(),
                report_data=job.toJsonStr(),
            )
            self.reportDB.create(report=report)
        except Exception as e:
            logging.error("Unable to load json for job ", e)