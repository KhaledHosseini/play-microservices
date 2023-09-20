Contents:
- [Summary](#summary)
- [Tools](#tools)
- [Docker dev environment](#docker-dev-environment)
- [Metadata service: Zookeeper](#metadata-service-zookeeper)
- [Zoonavigator service](#zoonavigator-service)
- [Message broker service: Kafka](#message-broker-service-kafka)
- [Kafka-ui service](#kafka-ui-service)
- [SMTP service: OPTIONAL](#smtp-service-optional)
- [Mailjob service: Python](#mailjob-service-python)
- [To DO](#to-do)

---

This is the 4th part of a series of articles under the name **"Play Microservices"**. Links to other parts<br/>
Part 1: [Play Microservices: Bird's eye view](https://dev.to/khaledhosseini/play-microservices-birds-eye-view-3d44)<br/>
Part 2: [Play Microservices: Authentication](https://dev.to/khaledhosseini/play-microservices-authentication-4di3)<br/>
Part 3: [Play Microservices: Scheduler](https://dev.to/khaledhosseini/play-microservices-scheduler-19km)<br/>
Part 4: You are here<br/>
Part 5: [Play Microservices: Report service](https://dev.to/khaledhosseini/play-microservices-report-service-4jcm)<br/>
Part 6: [Play Microservices: Api-gateway service](https://dev.to/khaledhosseini/play-microservices-api-gateway-service-4a9j)<br/>
Part 7: [Play Microservices: Client service](https://dev.to/khaledhosseini/play-microservices-client-service-4jbf)<br/>
Part 8: [Play Microservices: Integration via docker-compose](https://dev.to/khaledhosseini/play-microservices-integration-via-docker-compose-2ddc)<br/>
Part 9: [Play Microservices: Security](https://dev.to/khaledhosseini/play-microservices-security-45e4)<br/>

---

## Summary

In the 3rd part, we developed a Scheduler service. Now, our objective is to create an Email job executer service that consumes events of topic-job-run produced by scheduler service and after running the job, produce topic-job-run-result to the message broker. To achieve this, we need four distinct services: a message broker service, a metadata database service dedicated to supporting the message broker, an smtp-server service and the email service itself. Additionally, in the development environment, we include two extra services specifically for debugging purposes. These services consist of Kafkaui for managing our Kafka service, Zoonavigator for the Zookeeper service. Please keep in mind that we are following service per team pattern, assuming each service in our microservice architecture has its own repository and development team and in the development environment they are completely independent (Technically not logically). The first four parts are a copy from the previous step. You can skip them and got to  **SMTP service**


![Summary](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/i9r4om74wndas4furupa.png)

In the end, the project directory structure will appear as follows:



![Folder structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/il43d4u239n88xtpaebl.PNG)

---

## Tools

The tools required In the host machine:

  - [Docker](https://www.docker.com/): Containerization tool
  - [VSCode](https://code.visualstudio.com/): Code editing tool
  - [Dev containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension for VSCode
  - [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) extension for VSCode
  - [Git](https://git-scm.com/)

The tools and technologies that we will use Inside containers for each service:

 - Messaging service: [Kafka](https://hub.docker.com/r/bitnami/kafka/)
 - Kafka-ui service: [Kafka-ui](https://hub.docker.com/r/provectuslabs/kafka-ui)
 - Metadata service: [Zookeeper](https://hub.docker.com/_/zookeeper)
 - Zoonavigator service: [Zoonavigator](https://hub.docker.com/r/elkozmon/zoonavigator)
 - SMTP service: [Postfix](https://hub.docker.com/r/catatnight/postfix)
 - Email service: 
  - [Python](https://www.python.org/) : programming language
  - [Kafka-python](https://pypi.org/project/kafka-python/): For our message broker communications from go.

---

## Docker dev environment

Development inside Docker containers can provide several benefits such as consistent environments, isolated dependencies, and improved collaboration. By using Docker, development workflows can be containerized and shared with team members, allowing for consistent deployments across different machines and platforms. Developers can easily switch between different versions of dependencies and libraries without worrying about conflicts.

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/id0cijn8hm77ohn0hk2j.png)

When developing inside a Docker container, you only need to install `Docker`, `Visual Studio Code`, and the `Dev Containers` and `Docker `extensions on VS Code. Then you can run a container using Docker and map a host folder to a folder inside the container, then attach VSCode to the running container and start coding, and all changes will be reflected in the host folder. If you remove the images and containers, you can easily start again by recreating the container using the Dockerfile and copying the contents from the host folder to the container folder. However, it's important to note that in this case, any tools required inside the container will need to be downloaded again. Under the hood, When attaching VSCode to a running container, Visual Studio code install and run a special server inside the container which handle the sync of changes between the container and the host machine.

---

## Metadata service: Zookeeper

[ZooKeeper ](https://zookeeper.apache.org/) is a centralized service for maintaining configuration information. we use it as metadata storage for our Kafka messaging service.
> - Inside root directory create a folder with the name zookeeper-service
> - Create a Dockerfile and set content to `FROM bitnami/zookeeper:3.8.1`
> - Create a file named .env and set content to 

```bash
ZOO_SERVER_USERS=admin,user1
# for development environment only
ALLOW_ANONYMOUS_LOGIN="yes"
# if yes, uses SASL
ZOO_ENABLE_AUTH="no" 
```
> - Create a file named server_passwords.properties and set content to `password123,password_for_user1` Please choose your own passwords.
> - Add the following to the .env file of the docker-compose (the .env file at the root directory of the project.)

```bash
ZOOKEEPER_PORT=2181
ZOOKEEPER_ADMIN_CONTAINER_PORT=8078
ZOOKEEPER_ADMIN_PORT=8078
```
> - Add the following to the service part of the docker-compose.yml.

```yaml
  e-zk1:
    build:
      context: ./zookeeper-service
      dockerfile: Dockerfile
    container_name: e-zk1-service
    secrets:
      - zoo-server-pass
    env_file:
      - ./zookeeper-service/.env
    environment:
      ZOO_SERVER_ID: 1
      ZOO_SERVERS: e-zk1-service:${ZOOKEEPER_PORT}:${ZOOKEEPER_PORT}
      ZOO_SERVER_PASSWORDS_FILE: /run/secrets/zoo-server-pass
      ZOO_ENABLE_ADMIN_SERVER: yes
      ZOO_ADMIN_SERVER_PORT_NUMBER: ${ZOOKEEPER_ADMIN_CONTAINER_PORT}
    ports:
      - '${ZOOKEEPER_PORT}:${ZOOKEEPER_PORT}'
      - '${ZOOKEEPER_ADMIN_PORT}:${ZOOKEEPER_ADMIN_CONTAINER_PORT}'
    volumes:
      - "e-zookeeper_data:/bitnami"

volumes:
  e-zookeeper_data:
    driver: local
```
> - Add the following to the secrets part of the docker-compose.yml.

```yaml
  zoo-server-pass:
    file: zookeeper-service/server_passwords.properties
```
> -ZooKeeper is a distributed application that allows us to run multiple servers simultaneously. It enables multiple clients to connect to these servers, facilitating communication between them. ZooKeeper servers collaborate to handle data and respond to requests in a coordinated manner. In this case, our zookeeper consumers (clients) are Kafka servers which is again a distributed event streaming platform. We can run multiple zookeeper services as an ensemble of zookeeper servers and attach them together via `ZOO_SERVERS` environment variable.
> - The Bitnami ZooKeeper Docker image provides a zoo_client entrypoint, which acts as an internal client and allows us to run the zkCli.sh command-line tool to interact with the ZooKeeper server as a client. But we are going to use a GUI client for debugging purposes: Zoonavigator.

---

## Zoonavigator service

This service exists only in the development environment for debugging purposes. We use it to connect to zookeeper-service and manage the data.

> - Inside root directory create a folder with the name zoonavigator-service
> - Create a Dockerfile and set content to `FROM elkozmon/zoonavigator:1.1.2`
> - Add `ZOO_NAVIGATOR_PORT=9000` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  e-zoo-navigator:
    build: 
      context: ./zoonavigator-service
      dockerfile: Dockerfile
    container_name: e-zoo-navigator-service
    ports:
      - '${ZOO_NAVIGATOR_PORT}:${ZOO_NAVIGATOR_PORT}'
    environment:
      - CONNECTION_LOCALZK_NAME = Local-zookeeper
      - CONNECTION_LOCALZK_CONN = localhost:${ZOOKEEPER_PORT}
      - AUTO_CONNECT_CONNECTION_ID = LOCALZK
    depends_on:
      - e-zk1
``` 

> - Now from the terminal run `docker-compose up -d --build`
> - While running go to `http://localhost:9000/`. You will see the following screen:


![zoonavigatoe](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/h151vcjvlc60lp6hr7lq.PNG)

> - Enter the container name of a zookeeper service (here e-zk1).  If everything goes according to plan, you should be able to establish a connection to the ZooKeeper service.

![zoonavigator-zk1](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/aldo17byoo9ef9uppbcc.PNG) 

> - No run docker-compose down. We will return to these tools later.

---

## Message broker service: Kafka

[Apache Kafka](https://kafka.apache.org/)  is an open-source distributed event streaming platform that is well-suited for Microservices architecture. It is an ideal choice for implementing patterns such as event sourcing. Here We use it as an message broker for our scheduler service. 

> - Inside root directory create a folder with the name kafka-service
> - Create a Dockerfile and set content to `FROM bitnami/kafka:3.4.1`
> - Create a .env file beside the Docker file and set the content to: 

```bash
ALLOW_PLAINTEXT_LISTENER=yes
KAFKA_ENABLE_KRAFT=no
```
> - Add `KAFKA1_PORT=9092` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  e-kafka1:
    build: 
      context: ./kafka-service
      dockerfile: Dockerfile
    container_name: e-kafka1-service
    ports:
      - '${KAFKA1_PORT}:${KAFKA1_PORT}'
    volumes:
      - "e-kafka_data:/bitnami"
    env_file:
      - ./kafka-service/.env
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_CFG_ZOOKEEPER_CONNECT: zk1:${ZOOKEEPER_PORT}
      KAFKA_ZOOKEEPER_PROTOCOL: PLAINTEXT #if auth is enabled in zookeeper use one of: SASL, SASL_SSL see https://hub.docker.com/r/bitnami/kafka
      KAFKA_CFG_LISTENERS: PLAINTEXT://:${KAFKA1_PORT}
    depends_on:
      - e-zk1
```
> - To connect to our Kafka brokers for debugging purposes, we run another service. Kafka-ui.


---

## Kafka-ui service

This service exists only in the development environment for debugging purposes. We use it to connect to kafka-service and manage the data.

> - Inside root directory create a folder with the name kafkaui-service
> - Create a Dockerfile and set content to `FROM provectuslabs/kafka-ui:latest`
> - Add `KAFKAUI_PORT=8080` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  e-kafka-ui:
    build: 
      context: ./kafkaui-service
      dockerfile: Dockerfile
    container_name: e-kafka-ui-service
    restart: always
    ports:
      - ${KAFKAUI_PORT}:${KAFKAUI_PORT}
    environment:
     KAFKA_CLUSTERS_0_NAME: local
     KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: e-kafka1:${KAFKA1_PORT}
     DYNAMIC_CONFIG_ENABLED: 'true'
    depends_on:
      - e-kafka1
```

> - Now run `docker-compose run -d --build`. While containers are running, go to `http://localhost:8080/` to open Kafka-ui dashboard. 


![Kafka-ui](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/gcwqxf16jbig0q3xxzn7.PNG)

> - From the interface, you have the ability to view and manage brokers, topics, and consumers. We'll revisit these elements in more detail shortly.
> - Run `docker-compose down`

---

## SMTP service: OPTIONAL

We want to run a Simple Outbound Email Service on our local computer via docker. This step is optional and it is preferred to use 3rd party services like Amazon SES.

> - Inside root directory create a folder with the name postfix-service
> - Create a Dockerfile and set content to `FROM FROM catatnight/postfix:latest`
> - Add `POSTFIX_PORT=25` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
    postfix:
    build: 
      context: ./postfix-service
      dockerfile: Dockerfile
    container_name: postfix
    restart: always
    environment:
      - EMAIL_DOMAIN=yourdomain.com
      - SMTP_USER=username
      - SMTP_PASSWORD=password
    ports:
      - ${POSTFIX_PORT}:${POSTFIX_PORT}
```

---

> - Our required services are ready and running. Now it is time to Prepare development environment for our email job executer service. 

---

## Mailjob service: Python

Our goal is to create a simple python service that consumes kafka events for topic `topic-job-run` (produce by our scheduler service), execute the job (sending an email) and then produce an event for topic `topic-job-run-result`. A piece of cake!

![piece of cake](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/a7mki930114n2a7tby34.png) 


> - Create a folder named `mailjob-service` inside scheduler folder.
> - Create a Dockerfile inside `scheduler-service` and set the contents to

```bash
FROM python:3.11.4

WORKDIR /usr/src/app
```

> - Add the following to the service part of our docker-compose.yml file.

```yaml
  mailjob-service:
    build: 
      context: ./mailjob-service
      dockerfile: Dockerfile
    container_name: mailjob-service
    command: sleep infinity
    ports:
      - ${EMAIL_SERVICE_PORT}:${EMAIL_SERVICE_PORT}
    environment:
      ENVIRONMENT: development
      KAFKA_BROKERS: e-kafka1-service:${KAFKA1_PORT}
      # TOPICS_FILE: ''
      MAIL_SERVER_HOST: postfix
      MAIL_SERVER_PORT: 25
      EMAIL_DOMAIN: yourdomain.com
      SMTP_USER: username
      SMTP_PASSWORD: password

    volumes:
      - ./mailjob-service:/usr/src/app
    depends_on:
      - e-kafka1
      - postfix
```

> - We are going to do all the development inside a docker container without installing Python in our host machine. To do so, we run the containers and then attach VSCode to the mailjob-service container. As you may noticed, the Dockerfile for mailjob-service has no entry-point therefore we set the command value of mailjob-service to `sleep infinity` to keep the container awake.
> - Now run `docker-compose up -d --build`
> - While running, attach to the mailjob service by clicking bottom-left icon and then select `attach to running container`. Select mailjob-service and wait for a new instance of VSCode to start. At the beginning the VScode asks us to open a folder inside the container. We have selected  `WORKDIR /usr/src/app` inside our Dockerfile, so we will open this folder inside the container. This folder is mounted to mailjob-service folder inside the host machine using docker compose volume, therefor whatever change we made will be synced to the host folder too.
> - After opening the folder `/usr/src/app`, create a file named requirements.txt and set the contents to:

```bash
kafka-python==2.0.2
python-dotenv==1.0.0
```
> - Open a new terminal and run `pip install -r requirements.txt`
> - Create a file named .env.topics and set the contents to:

```bash
TOPIC_JOB_RUN="topic-job-run"
TOPIC_JOB_RUN_CONSUMER_GROUP_ID="topic-job-run-consumer-emailexecutor"
TOPIC_JOB_RUN_CONSUMER_WORKER_COUNT=1
TOPIC_JOB_RUN_RESULT="topic-job-run-result"
TOPIC_JOB_RUN_RESULT_CREATE_PARTITIONS=1
TOPIC_JOB_RUN_RESULT_CREATE_REPLICAS=1
```
> - Create a file named config.py and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/email/mailjob-service/config.py).
> - create a folder named app then a file named \_\_init\_\_.py for the app module. set the contents to

```python
from config import Config

def create_app():
    config = Config()
```

> - Create a file named main.py in the root directory and set the contents to

```python
from app import create_app
import logging

def main():
    create_app()

if __name__ == '__main__':
    logging.basicConfig(
        format='%(asctime)s.%(msecs)s:%(name)s:%(thread)d:' + '%(levelname)s:%(process)d:%(message)s',
        level=logging.INFO
    )

    main()
```
> - Create a folder named email inside app, then a file named \_\_init\_\_.py and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/email/mailjob-service/app/email/email_sender.py). This class is responsible for executing the email job. 
> - Our goal is to consume kafka messages with the following structure. 

```json
{
    "Id": "649f07e619fca8aa63d842f6",
    "Name": "job1",
    "ScheduleTime": "2024-01-01T00:00:00Z",
    "CreatedAt": "2023-06-30T16:50:46.3042083Z",
    "UpdatedAt": "2023-06-30T16:50:46.3042086Z",
    "Status": 2,
    "JobData": {
        "SourceAddress": "example@example.com",
        "DestinationAddress": "example@example.com",
        "Subject": "subject ",
        "Message": "message"
    }
}
```

> - Create a folder named models inside app, then create \_\_init\_\_.py for the module. Set the contents to:

```python
import json

class JsonObject:
    def toJsonData(self):
        return self.toJsonStr().encode('utf-8')
    
    def toJsonStr(self):
        return json.dumps(self.to_dict())
    
    def to_dict(self):
        return vars(self)

from app.models._email_job import EmailJob
from app.models._email_job import Email
from app.models._email_job import JobStatusEnum
```
> - Create a file named _email_job.py and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/email/mailjob-service/app/models/_email_job.py).

> - Create this folder tree `email_job/message_broker` inside models folder. Then Create a file named \_\_init\_\_.py and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/email/mailjob-service/app/models/email_job/message_broker/__init__.py). This file contains our kafka server and worker classes which consumes `TOPIC_JOB_RUN`, then runs the job and finally produce the `TOPIC_JOB_RUN_RESULT` topic. `MessageBrokerService` class accepts config and email executer in the constructor then Creates `TOPIC_JOB_RUN_RESULT` if it doesnot exists and finally runs kafka consumer workers.

```python
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
```

> - The logic for running kafka consumer for `TOPIC_JOB_RUN` and producing `TOPIC_JOB_RUN_RESULT` is inside `JobConsumerWorker` class. We receive config and email sender class in the constructor and then run the consumer via following function:

```python
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
```
> Create a folder inside app folder named server and a file named \_\_init\_\_.py. set the content to:

```python
from config import Config
from app.email import EmailSender
from app.models.email_job.message_broker import MessageBrokerService

class Server:
    def run(self, cfg: Config):
        es = EmailSender(cfg)
        ecs = MessageBrokerService(cfg,es)
        ecs.run()
```

> - Now return to \_\_init\_\_.py of app folder and change the code to

```python
import logging
from config import Config
from app.server import Server

def create_app():
    logging.info("Creating job-executor app")
    config = Config()
    s = Server()
    s.run(config)
```
> - Now run `python run main.py` If everything goes according to plan, our app starts and waiting for messages from kafka.
> - go to `http://localhost:8080/`. Under topics select `topic-job-run`. From top right, click on `Produce message` button and set the value to the following json and click Produce message.

```json
{
    "Id": "649f07e619fca8aa63d842f6",
    "Name": "job1",
    "ScheduleTime": "2024-01-01T00:00:00Z",
    "CreatedAt": "2023-06-30T16:50:46.3042083Z",
    "UpdatedAt": "2023-06-30T16:50:46.3042086Z",
    "Status": 2,
    "JobData": {
        "SourceAddress": "example@example.com",
        "DestinationAddress": "example@example.com",
        "Subject": "subject ",
        "Message": "message"
    }
}

```


![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/rr33q1mz3g4ougv0zg0o.PNG)

> - The app consumes the message, Then try to send the email and then produce another message for the topic `topic-job-run-result` that in our microservice app will be consumed by the scheduler-service. You can go to `http://localhost:8080/` and under topics you can see that another topic has been created and it has one message. (Note: We have activated automatic topic creation in kafka in development environment. In production environment, it is common to have a separate service for topic creation and management for the whole microservice application).

---

## To DO

> - Add tests
> - Add tracing using Jaeger
> - Add monitoring and analysis using grafana
> - Refactoring

---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).