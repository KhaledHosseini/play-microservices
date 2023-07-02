from threading import Thread
from config import Config
from app.models.message_broker.consumer import JobConsumer
from app.email.email_sender import EmailSender

def worker(config: Config):
  email_sender = EmailSender(config)
  job_consumer = JobConsumer(config, email_sender)
  job_consumer.run()

def run_job_consumers(config: Config):
      
      if config.TopicJobRunWorkerCount == 1:
        worker(config)
      else:
        worker_threads = []
        for i in range(0,config.TopicJobRunWorkerCount):
          print("running worker ", i)
          t = Thread(target=worker, args=(config))
          t.Daemon = True
          worker_threads.append(t)
        for t in worker_threads:
          t.join()