from threading import Thread
import logging

from config import Config
from app.email.email_sender import EmailSender
from app.models.email_job.message_broker.emailjob_consumer_worker import JobConsumer


class EmailjobConsumerService:
    def __init__(self,cfg: Config, emailExecuter: EmailSender) -> None:
        self.cfg = cfg
        self.emailExecuter = emailExecuter

    def run(self):
        logging.info(f"Starting email job consumers with  {self.cfg.TopicJobRunWorkerCount } workers...")
        if self.cfg.TopicJobRunWorkerCount == 1:
            self.run_worker()
            logging.info(f"Worker 1 started for consuming job events...")
        else:
            worker_threads = []
            for i in range(0,self.cfg.TopicJobRunWorkerCount):
                t = Thread(target=self.worker)
                t.Daemon = True
                worker_threads.append(t)
                t.start()
                logging.info(f"Worker {i} started for consuming job events...")

            for t in worker_threads:
                t.join()
    def run_worker(self):
        job_consumer = JobConsumer(self.cfg, self.emailExecuter)
        job_consumer.run()