from kafka import KafkaConsumer, KafkaProducer
import sys
import signal
import json
import logging

from config import Config
from app.models.email_job import EmailJob
from app.email.email_sender import EmailSender

class JobConsumer:
    def __init__(self, config: Config, email_sender: EmailSender) -> None:
        self.cfg = config
        self.email_sender = email_sender
        self.producer = KafkaProducer(bootstrap_servers=self.cfg.KafkaBrokers)

    def run_kafka_consumer(self):
        consumer = KafkaConsumer(self.cfg.TopicJobRun,
        group_id=self.cfg.TopicJobRunConsumerGroupID, 
        bootstrap_servers=self.cfg.KafkaBrokers,
        value_deserializer= self.loadJson)

        for msg in consumer:
            if isinstance(msg.value, EmailJob):
                logging.info("An email job json recieved. doing the job.")
                self.handleJob(msg.value)
            else:
                logging.error(f"error handling: {msg}")

    def handleJob(self,job: EmailJob):
        try:
            self.email_sender.send(job.get_email())
            logging.info(f"email sending Succeeded:")
            #  SUCCEEDED = 4;
            #  FAILED = 5;
            job.status = 4
            message = job.toJsonData()
            self.producer.send(self.cfg.TopicJobRunResult, value=message)
            self.producer.flush()
        except:
            logging.warning(f"email sending failed:")
            job.status = 5
            message = job.toJsonData()
            self.producer.send(self.cfg.TopicJobRunResult, value=message)
            self.producer.flush()

    def loadJson(self,value):
        logging.info(f"decoding message: {value}")
        try:
            js = json.loads(value.decode('utf-8'))
            return EmailJob(js)
        except json.decoder.JSONDecodeError as e: 
            logging.error("invalid email job json:", e)
            return "invalid json"

    def run(self):
        signal.signal(signal.SIGINT, self.handle_shutdown)
        signal.signal(signal.SIGTERM, self.handle_shutdown)
        self.run_kafka_consumer()

    def handle_shutdown(self, signal, frame):
        logging.info("Shutting down Kafka consumer and producer...")
        # Clean up Kafka consumer and producer here if necessary
        sys.exit(0)