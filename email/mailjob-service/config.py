from dotenv import load_dotenv
import os
from pathlib import Path

class Config:
    def __init__(self) -> None:
        
        self.Environment = os.getenv("ENVIRONMENT")

        topicsFile = os.getenv("TOPICS_FILE")
        if topicsFile is None:
            topicsFile = "./.env.topics"
        dotenv_path = Path(topicsFile)
        load_dotenv(dotenv_path=dotenv_path)

        self.GCP_PROJECT_ID = os.getenv('GCP_PROJECT_ID')
        self.TopicJobRun = os.getenv("TOPIC_JOB_RUN")
        self.TopicJobRunConsumerGroupID = os.getenv("TOPIC_JOB_RUN_CONSUMER_GROUP_ID")
        self.TopicJobRunWorkerCount = int(os.getenv("TOPIC_JOB_RUN_CONSUMER_WORKER_COUNT"))
        self.TopicJobRunResult = os.getenv("TOPIC_JOB_RUN_RESULT")
        self.TopicJobRunResultConsumerGroupID = os.getenv("TOPIC_JOB_RUN_RESULT_CONSUMER_GROUP_ID")
        self.TopicJobRunResultConsumerWorkerCount = os.getenv("TOPIC_JOB_RUN_RESULT_CONSUMER_WORKER_COUNT")

        self.KafkaBrokers = os.getenv("KAFKA_BROKERS").split(",")

        self.MailServerHost = os.getenv("MAIL_SERVER_HOST")
        self.MailServerPort = os.getenv("MAIL_SERVER_PORT")
        self.EMailDomain =  os.getenv("EMAIL_DOMAIN")
        self.SmtpUser =  os.getenv("SMTP_USER")
        self.SmtpPassord =  os.getenv("SMTP_PASSWORD")