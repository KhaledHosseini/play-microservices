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

        self.ServerPort = os.getenv("SERVER_PORT")
        self.TopicJobRun = os.getenv("TOPIC_JOB_RUN")
        self.TopicJobRunConsumerGroupID = os.getenv("TOPIC_JOB_RUN_CONSUMER_GROUP_ID")
        self.TopicJobRunWorkerCount = int(os.getenv("TOPIC_JOB_RUN_CONSUMER_WORKER_COUNT"))
        self.TopicJobRunResult = os.getenv("TOPIC_JOB_RUN_RESULT")
        self.TopicJobRunResultConsumerGroupID = os.getenv("TOPIC_JOB_RUN_RESULT_CONSUMER_GROUP_ID")
        self.TopicJobRunResultConsumerWorkerCount = int(os.getenv("TOPIC_JOB_RUN_RESULT_CONSUMER_WORKER_COUNT"))
        self.TopicJobCreate = os.getenv("TOPIC_JOB_CREATE")
        self.TopicJobCreateConsumerGroupID = os.getenv("TOPIC_JOB_CREATE_CONSUMER_GROUP_ID")
        self.TopicJobCreateConsumerWorkerCount = int(os.getenv("TOPIC_JOB_CREATE_CONSUMER_WORKER_COUNT"))
        self.TopicJobUpdate = os.getenv("TOPIC_JOB_UPDATE")
        self.TopicJobUpdateConsumerGroupID = os.getenv("TOPIC_JOB_UPDATE_CONSUMER_GROUP_ID")
        self.TopicJobUpdateConsumerWorkerCount = int(os.getenv("TOPIC_JOB_UPDATE_CONSUMER_WORKER_COUNT"))

        self.KafkaBrokers = os.getenv("KAFKA_BROKERS").split(",")

        db_user_file = os.getenv("DATABASE_USER_FILE")
        self.DatabaseUser = self.fileContent(db_user_file)
        db_pass_file = os.getenv("DATABASE_PASS_FILE")
        self.DatabasePass = self.fileContent(db_pass_file)
        db_name_file = os.getenv("DATABASE_DB_NAME_FILE")
        self.DatabaseDBName = self.fileContent(db_name_file)
        self.DatabaseSchema = os.getenv("DATABASE_SCHEMA")
        self.DatabaseHostName = os.getenv("DATABASE_HOST_NAME")
        self.DatabasePort = os.getenv("DATABASE_PORT")

        auth_publickey_file = os.getenv("AUTH_PUBLIC_KEY_FILE")
        self.AuthPublicKey = self.fileContent(auth_publickey_file).encode()
    
    def fileContent(self,path:str):
        with open(path, 'r') as file:
          return file.read()