from threading import Thread
import logging

from config import Config
from app.email import EmailSender

from kafka.admin import KafkaAdminClient, NewTopic

class MessageBrokerService:
    def __init__(self,cfg: Config, emailExecuter: EmailSender) -> None:
        self.cfg = cfg
        self.emailExecuter = emailExecuter
        # producers are responsible for creation of topics they produce
        admin_client = KafkaAdminClient(bootstrap_servers=cfg.KafkaBrokers)
        topic_exists = cfg.TopicJobRunResult in admin_client.list_topics()
        if not topic_exists:
            try:
                topic_list = []
                topic_list.append(NewTopic(name=cfg.TopicJobRunResult, num_partitions=cfg.TopicJobRunResultCreatePartitions, replication_factor=cfg.TopicJobRunResultCreateReplicas))
                admin_client.create_topics(new_topics=topic_list, validate_only=False)
                logging.info(f"Topic '{cfg.TopicJobRunResult}' created successfully.")
            except TopicAlreadyExistsError as e:
                logging.info(f"Topic '{cfg.TopicJobRunResult}' already exists. Error: {e}")
        else:
            logging.info(f"Topic '{cfg.TopicJobRunResult}' already exists.")
        

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
        job_consumer = JobConsumerWorker(self.cfg, self.emailExecuter)
        job_consumer.run()


from kafka import KafkaConsumer, KafkaProducer
import sys
import signal
import json
from app.models import EmailJob, JobStatusEnum
class JobConsumerWorker:
    def __init__(self, config: Config, email_sender: EmailSender) -> None:
        self.cfg = config
        self.email_sender = email_sender
        self.producer = KafkaProducer(bootstrap_servers=self.cfg.KafkaBrokers,compression_type='snappy')

    def run_kafka_consumer(self):
        consumer = KafkaConsumer(self.cfg.TopicJobRun,
        group_id=self.cfg.TopicJobRunConsumerGroupID, 
        bootstrap_servers=self.cfg.KafkaBrokers,
        value_deserializer= self.loadJson)

        for msg in consumer:
            if isinstance(msg.value, EmailJob):
                logging.info(f"An email job json recieved. handling the job. {msg.value}")
                self.handleJob(msg.value)
            else:
                logging.error(f"error handling: {msg}")

    def handleJob(self,job: EmailJob):
        try:
            self.email_sender.send(job.get_email())
            logging.info(f"email sending Succeeded:")
            job.status = JobStatusEnum.JOB_STATUS_SUCCEEDED
            message = job.toJsonData()
            self.producer.send(self.cfg.TopicJobRunResult, value=message)
            self.producer.flush()
        except:
            logging.warning(f"email sending failed:")
            job.status = JobStatusEnum.JOB_STATUS_FAILED
            message = job.toJsonData()
            self.producer.send(self.cfg.TopicJobRunResult, value=message)
            self.producer.flush()

    def loadJson(self,value):
        logging.info(f"decoding message: {value}")
        try:
            js = json.loads(value.decode('utf-8'))
            return EmailJob.from_dict(js)
        except Exception as e: 
            logging.error(e)
            return "ERROR"

    def run(self):
        signal.signal(signal.SIGINT, self.handle_shutdown)
        signal.signal(signal.SIGTERM, self.handle_shutdown)
        self.run_kafka_consumer()

    def handle_shutdown(self, signal, frame):
        logging.info("Shutting down Kafka consumer and producer...")
        # Clean up Kafka consumer and producer here if necessary
        sys.exit(0)