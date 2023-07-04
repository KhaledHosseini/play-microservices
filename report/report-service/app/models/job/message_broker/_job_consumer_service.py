from threading import Thread
import logging

from config import Config
from app.models import ReportDB

from app.models.job.message_broker._job_consumer_worker import JobConsumerWorker


class ConsumerService:
  def __init__(self,cfg: Config, reportDb: ReportDB):
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