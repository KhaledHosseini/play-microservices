from threading import Thread
from config import Config
from app.models.message_broker.consumer import JobConsumer
from app.email.email_sender import EmailSender
import logging

def worker(config: Config):
  email_sender = EmailSender(config)
  job_consumer = JobConsumer(config, email_sender)
  job_consumer.run()

def run_job_consumers(config: Config):
      logging.info("Running email job consumers...")
      if config.TopicJobRunWorkerCount == 1:
        worker(config)
        logging.info(f"Worker 1 started for consuming job events...")
      else:
        worker_threads = []
        for i in range(0,config.TopicJobRunWorkerCount):
          t = Thread(target=worker, args=(config))
          t.Daemon = True
          worker_threads.append(t)
          t.start()
          logging.info(f"Worker {i} started for consuming job events...")
        for t in worker_threads:
          t.join()