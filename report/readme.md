Contents:
- [Summary](#summary)
- [Tools](#tools)
- [Docker dev environment](#docker-dev-environment)
- [Development services](#development-services)
- [Report service: Python](#report-service-python)
- [To DO](#to-do)


---

This is the 5th part of a series of articles under the name **"Play Microservices"**. Links to other parts:<br/>
Part 1: [Play Microservices: Bird's eye view](https://dev.to/khaledhosseini/play-microservices-birds-eye-view-3d44)<br/>
Part 2: [Play Microservices: Authentication](https://dev.to/khaledhosseini/play-microservices-authentication-4di3)<br/>
Part 3: [Play Microservices: Scheduler service](https://dev.to/khaledhosseini/play-microservices-scheduler-19km)<br/>
Part 4: [Play Microservices: Email service](https://dev.to/khaledhosseini/play-microservices-email-service-1kmc)<br/>
Part 5: You are here<br/>
Part 6: [Play Microservices: Api-gateway service](https://dev.to/khaledhosseini/play-microservices-api-gateway-service-4a9j)<br/>
Part 7: [Play Microservices: Client service](https://dev.to/khaledhosseini/play-microservices-client-service-4jbf)<br/>
Part 8: [Play Microservices: Integration via docker-compose](https://dev.to/khaledhosseini/play-microservices-integration-via-docker-compose-2ddc)<br/>
Part 9: [Play Microservices: Security](https://dev.to/khaledhosseini/play-microservices-security-45e4)<br/>

---

## Summary

Up to now we have developed 4 services. Now, our objective is to create a report service for our microservices application. This service is responsible for saving all events in our application. It listens to all kafka events and save them in a database. To develop this service , we need four distinct services: a database service, a message broker service, a metadata database service dedicated to supporting the message broker, and the report service itself, which is a gRPC API service. Additionally, in the development environment, we include four extra services specifically for debugging purposes. These services consist of Mongo Express, used to manage our database service, Kafkaui for managing our Kafka service, Zoonavigator for the Zookeeper service, and grpcui for testing our gRPC API.



![Summary](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/h5bbowu5qqtprxhkwitt.png)

At the end, the project directory structure will appear as follows:



![Folder structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/jwec96d6k79a2k6zqa7s.PNG)

---

## Tools

The tools required In the host machine:

  - [Docker](https://www.docker.com/): Containerization tool
  - [VSCode](https://code.visualstudio.com/): Code editing tool
  - [Dev containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension for VSCode
  - [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) extension for VSCode
  - [Git](https://git-scm.com/)

The tools and technologies that we will use Inside containers for each service:

 - Database service: [Mongo](https://hub.docker.com/_/mongo)
 - Mongo express service: [Mongo express](https://hub.docker.com/_/mongo-express)
 - Messaging service: [Kafka](https://hub.docker.com/r/bitnami/kafka/)
 - Kafka-ui service: [Kafka-ui](https://hub.docker.com/r/provectuslabs/kafka-ui)
 - Metadata service: [Zookeeper](https://hub.docker.com/_/zookeeper)
 - Zoonavigator service: [Zoonavigator](https://hub.docker.com/r/elkozmon/zoonavigator)
 - grpcui service: [grpcui](https://hub.docker.com/r/fullstorydev/grpcui)
 - report service: 
  - [Python](https://python.org/) : programming language
  - [grpc](https://pypi.org/project/grpc/): gRPC framework for python
  - [Pymongo](https://pypi.org/project/pymongo/): Query builder for our database communication.
  - [Kafka-Python](https://pypi.org/project/kafka-python/) for our message broker communications from python.
  - [grpcio-tools](https://pypi.org/project/grpcio-tools/) For compiling .proto files to python.

---

## Docker dev environment

Development inside Docker containers can provide several benefits such as consistent environments, isolated dependencies, and improved collaboration. By using Docker, development workflows can be containerized and shared with team members, allowing for consistent deployments across different machines and platforms. Developers can easily switch between different versions of dependencies and libraries without worrying about conflicts.

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/id0cijn8hm77ohn0hk2j.png)

When developing inside a Docker container, you only need to install `Docker`, `Visual Studio Code`, and the `Dev Containers` and `Docker `extensions on VS Code. Then you can run a container using Docker and map a host folder to a folder inside the container, then attach VSCode to the running container and start coding, and all changes will be reflected in the host folder. If you remove the images and containers, you can easily start again by recreating the container using the Dockerfile and copying the contents from the host folder to the container folder. However, it's important to note that in this case, any tools required inside the container will need to be downloaded again. Under the hood, When attaching VSCode to a running container, Visual Studio code install and run a special server inside the container which handle the sync of changes between the container and the host machine.

---

## Development services

> - Create a folder for the project and choose a name for it (such as 'microservice'). Then create a folder named `report`. This folder is the root directory of current project. You can then open the root folder in VS Code by right-clicking on the folder and selecting 'Open with Code'.

> - Our approach is to develop all services of our microservice application independently, and as development services are Identical to the part 3 you can copy the folders for development services from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/report). Do not forget to also copy docker-compose.yml file.

---

## Report service: Python

Our goal is to create a python service that consumes kafka events for all topics in our application (produced by other services) and then store them in a database. The admin then can query the reports. To respond to the queries, we run an grpc service and for listening to kafka events we subscribe to kafka events using kafka consumers. lets begin.


> - Create a folder named `report-service` inside report folder.
> - Create a Dockerfile inside `report-service` and set the contents to

```bash
FROM python:3.11.4

RUN apt-get update \
  && apt-get install -y --no-install-recommends graphviz libsnappy-dev \
  && rm -rf /var/lib/apt/lists/* \
  && pip install --no-cache-dir pyparsing pydot

WORKDIR /usr/src/app
```
> - Inside `report-service` folder create a folder named keys and then go to this [link](https://cryptotools.net/rsagen), generate a key pair with length 2048. Create a file named `access_token.public.pem` inside the keys folder and then paste the public key. Keep the private key somewhere else. We will use it later. This public key acts as the public key of authentication service and we are going to verify the JWTs with it.
> - Add the following to the service part of our docker-compose.yml file.(if you have copied the file, It already has this data).

```yaml
  report-service:
    build: 
      context: ./report-service
      dockerfile: Dockerfile
    container_name: report-service
    command: sleep infinity
    ports:
      - ${REPORT_SERVICE_PORT}:${REPORT_SERVICE_PORT}
    environment:
      ENVIRONMENT: development
      SERVER_PORT: ${REPORT_SERVICE_PORT}
      DATABASE_USER_FILE: /run/secrets/report-db-user
      DATABASE_PASS_FILE: /run/secrets/report-db-pass
      DATABASE_DB_NAME_FILE: /run/secrets/report-db-dbname
      DATABASE_SCHEMA: mongodb
      DATABASE_HOST_NAME: report-db-service
      DATABASE_PORT: ${MONGODB_PORT}
      KAFKA_BROKERS: r-kafka1-service:${KAFKA1_PORT}
      AUTH_PUBLIC_KEY_FILE: /run/secrets/auth-public-key
      # TOPICS_FILE: ''
    volumes:
      - ./report-service:/usr/src/app
    depends_on:
      - r-kafka1
    secrets:
      - report-db-user
      - report-db-pass
      - report-db-dbname
      - auth-public-key
```

> - We are going to do all the development inside a docker container without installing Python in our host machine. To do so, we run the containers and then attach VSCode to the report-service container. As you may noticed, the Dockerfile for report-service has no entry-point therefore we set the command value of report-service to `sleep infinity` to keep the container awake.
> - Now run `docker-compose up -d --build`
> - While running, attach to the report-service by clicking bottom-left icon and then select `attach to running container`. Select report-service and wait for a new instance of VSCode to start. At the beginning the VScode asks us to open a folder inside the container. We have selected  `WORKDIR /usr/src/app` inside our Dockerfile, so we will open this folder inside the container. This folder is mounted to report-service folder inside the host machine using docker compose volume (see volumes of report-service inside docker-compose), therefor whatever change we made will be synced to the host folder too.
> - After opening the folder `/usr/src/app`, create a file named requirements.txt and set the contents to:

```bash
kafka-python==2.0.2
python-snappy==0.6.1
crc32c==2.3.post0
python-dotenv==1.0.0
pymongo==4.4.0
grpcio-tools==1.56.0
grpcio-reflection==1.56.0
python_jwt==4.0.0
```
> - Open a new terminal and run `pip install -r requirements.txt`
> - Create a file named .env.topics and set the contents to:

```bash
TOPIC_JOB_RUN="topic-job-run"
TOPIC_JOB_RUN_CONSUMER_GROUP_ID="topic-job-run-consumer-report"
TOPIC_JOB_RUN_CONSUMER_WORKER_COUNT=1
TOPIC_JOB_RUN_RESULT="topic-job-run-result"
TOPIC_JOB_RUN_RESULT_CONSUMER_GROUP_ID="job-run-result-consumer-report"
TOPIC_JOB_RUN_RESULT_CONSUMER_WORKER_COUNT=1
TOPIC_JOB_CREATE="topic-job-create"
TOPIC_JOB_CREATE_CONSUMER_GROUP_ID="topic-job-create-consumer-report"
TOPIC_JOB_CREATE_CONSUMER_WORKER_COUNT= 1
TOPIC_JOB_UPDATE="topic-job-update"
TOPIC_JOB_UPDATE_CONSUMER_GROUP_ID="topic-job-update-consumer-report"
TOPIC_JOB_UPDATE_CONSUMER_WORKER_COUNT=1
```
> - Create a file named config.py and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/config.py).
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
> - A summary of what we are going to do: First we define a .proto file for our report. Then we compile the .proto file to python. As it is common for separating protocol models from database models, we define database models and implement the logic for converting to/from protocol models. Then we create a report service and implement the grpc protocol interfaces. This report service read the reports from the database and send it to the end user. Next part is to prepare a database for our report service for reading and writing the reports and finally we prepare our kafka consumer classes to listen to kafka events and then save them to reports database. The folder structure for our service is as follows:


![report folder structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/u9s484lxbkpwt5bw6tuw.PNG)


> - Create a folder named proto and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/proto/report.proto). Now create a file named build_grpc.sh and set the content to:

```bash
#!/bin/bash

declare -a services=("proto")
for SERVICE in "${services[@]}"; do
    DESTDIR='proto'
    mkdir -p $DESTDIR
    python -m grpc_tools.protoc \
        --proto_path=$SERVICE/ \
        --python_out=$DESTDIR \
        --grpc_python_out=$DESTDIR \
        $SERVICE/*.proto
done
```
> - Now run `source build_grpc.sh`. This will compile .proto files and store them in proto directory(As this is a linux command, be careful about the line endings to be LF and not CRLF). Add a file named \_\_init\_\_.py file and set the contents to the following. We do this to avoid repeating ugly file names when importing them in other modules.

```
import proto.report_pb2_grpc as ReportGRPC
import proto.report_pb2 as ReportGRPCTypes
```
> - Also change the `import report_pb2 as report__pb2` inside generated `report_pb2_grpc.py` to `import proto.report_pb2 as report__pb2` to avoid import errors later.
> - Create a folder named models inside app folder. We have two models here. One is Report which we transfer via gRPC and the other is Job which we receive via message events. Here we develop our logic on top of our models. For the Report model, we need a database and a grpc service. We define the Report model inside model module and the database and grpc services inside another module called report. For the job model we need a consumer service to read it from kafka events and as the same logic for report we define the job model itself inside models module and the consumer service inside another module called job. 
> - Create two files named \_report.py and _job.py and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/app/models/_report.py) and [here ](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/app/models/_job.py)respectively. We Add \_ to the files to indicate them as internal files. We expose them via \_\_init\_\_.py file. Then set the contents of the \_\_init\_\_.py file inside models folder to

```python
import json

class JsonObject(object):
    def toJsonData(self):
        return self.toJsonStr().encode('utf-8')
    
    def toJsonStr(self):
        return json.dumps(self.to_dict())
    
    def to_dict(self):
        return vars(self)

from app.models._job import Job
from app.models._report import Report
from app.models._report import ReportList
from app.models._report import ReportDBInterface
```
> - Create a folder named `report `and then two folders named `db ` and `grpc`. As we are going to use Mongo for our database, first create another folder inside app called database_engines and then a file named \_mongo.py. set the content from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/app/database_engines/_mongo.py). Then expose it from \_\_init\_\_.py file using `from app.database_engines._mongo import MongoDB` Inside db folder create a folder named \_report_db_mongo.py and set the content from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/app/models/report/db/_report_db_mongo.py). This file contains ReportDBMongo which has inherited from two classes. One is MongoDB and the other is ReportDBInterface.

```python
class ReportDBMongo(MongoDB, ReportDBInterface):
    def __init__(self, cfg: Config):
        super().__init__(cfg)
        ReportDBInterface.__init__(self)
        self.dbName = cfg.DatabaseDBName
        self.collection = "ReportsCollection"
        self.db = self.client[self.dbName]
```
> - Inside grpc folder, create a file named \_report_grpc_service.py. This file contains our grpc service and implements our grpc protocol interface. Set the contents to:

```python
from proto import ReportGRPC
import logging

from app.models import ReportDBInterface

class MyReportService(ReportGRPC.ReportService):
    def __init__(self, reportDB: ReportDBInterface):
        self.reportDB = reportDB
        super().__init__()
    
    def ListReports(self,request, context):
        logging.info(f"message recieved...{request}")
        filter = request.filter
        page = request.page
        size = request.size
        result = self.reportDB.list(type=filter,page=page,size=size)
        return result.toProto()
```
> - Now it is time to run our job consumer service. As we are going to use kafka-python for our messaging service, first create a folder named message_engines and then a file named \_kafka_consumer_worker.py. set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/app/message_engines/_kafka_consumer_worker.py). This class contains the logic to create and run a kafka consumer for a specific topic. It receives topic name, consumer group id and the list of kafka broker addresses and then subscribes to the topic. whenever a new message arrives `messageRecieved` function will be called. We inherit this class and override `messageRecieved` function and do our logic.
```python
class KafkaConsumerWorker:
    def __init__(self, topic: str, consumerGroupId: str, kafkaBrokers: List[str]):
        self.Topic = topic
        self.ConsumerGroupId = consumerGroupId
        self.KafkaBrokers = kafkaBrokers
    def run_kafka_consumer(self):
        consumer = KafkaConsumer(self.Topic,
        group_id=self.ConsumerGroupId, 
        bootstrap_servers=self.KafkaBrokers,
        enable_auto_commit=False,
        value_deserializer= self.loadJson)

        logging.info("consumer is listening....")
        try:
            for message in consumer:
                json = self.loadJson(message.value)
                self.messageRecieved(self.Topic, json)
                #committing message manually after reading from the topic
                consumer.commit()
        except Exception as e:
            logging.error(f"consumer listener stoped with error: {e}")
        finally:
            consumer.close()
```
 
> - Inside models/job folder create another folder named message_broker and then a file named \_job_consumer_service.py. Inside this file we run consumers for all topics that their messages are a job model. Set the content from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/app/models/job/message_broker/_job_consumer_service.py). One important class in this file is `JobConsumerWorker `that has inherited from `KafkaConsumerWorker`. We have overridden the `messageRecieved` function and when a new job arrives, we create a new report in the ReportDB for it.

```python
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
```

> - Expose ConsumerService inside \_\_init\_\_.py file using `from app.models.job.message_broker._job_consumer_service import ConsumerService`

> - Our services are ready. Now lets run them. Create a folder named server inside app and the a file named \_\_init\_\_.py. Set the contents to 

```python
from concurrent import futures
import logging

import grpc
from grpc_reflection.v1alpha import reflection

from proto import ReportGRPC
from proto import ReportGRPCTypes

from config import Config

from app.models.report import MyReportService
from app.models.report import ReportDBMongo
from app.models.job import ConsumerService

class Server:
    def run(self, cfg: Config):
        logging.info("Running services...")
        db = ReportDBMongo(cfg)
        my_report_grpc_service = MyReportService(reportDB=db)
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        ReportGRPC.add_ReportServiceServicer_to_server(my_report_grpc_service, server)
        SERVICE_NAMES = (
            ReportGRPCTypes.DESCRIPTOR.services_by_name["ReportService"].full_name,
            reflection.SERVICE_NAME,
        )
        reflection.enable_server_reflection(SERVICE_NAMES, server)
        port = cfg.ServerPort
        server.add_insecure_port("[::]:" + port)
        server.start()

        consumer = ConsumerService(cfg,db)
        consumer.start()


        server.wait_for_termination()
```

> - Now go to \_\_init\_\_.py of app folder and change the contents to

```python
from config import Config
from app.server import Server

import logging
def create_app():
    logging.info("creating app...")
    config = Config()
    s = Server()
    s.run(config)
```

> - Finally it is time to run and test our server! run `python main.py`. If everything goes according to plan, both our grpc server and our kafka consumers are running. Go to `http://localhost:8080/` and under the topics, create a new topic called `topic-job-run` (we have defined this named in environment variables). Then click on the topic and produce a message for it by clicking on the top right corner button and pasting the following json content to the value of the message and then clicking produce message.

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

![produce kafka message](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/hzy8jrl3rda59vskzui0.PNG)

> - This message will be consumed by our app and a new report item will be created for it. To confirm, go to `http://localhost:8081/`. in the Mongo Express panel click reports_db and then view. Our report item has been created!

![mongo express](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/k40vqabf84kbxofxvtvb.PNG)

> - You can repeat and reproduce some more messages the same way.
> - Now lets test our grpc server and read our reports. First go to docker desktop and be sure that our grpcui service is running and attached to the service. (restart the service if necessary). Go to `http://localhost:5000/`. Fill in the form by selecting job for filter, 0 for page and 10 for page size and click invoke! if everything goes according to plan you will receive the results!


![grpcui result](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/wy947gjoguutwrrzzi2a.PNG)


> - The last step is to add authentication and authorization to our report service. We first check the authenticity of the JWT to see if it is issued by our auth service or not. then we read the role of inside JWT and if it is admin, we allow access to the reports.
> - Create a folder named `grpc_interceptors `inside app folder. then a file named `_token_validation_interceptor.py` inside `grpc_interceptors`  folder. Set the contents of file from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/report/report-service/app/grpc_interceptors/_token_validation_interceptor.py). This file contains an iterceptor that our grpc calls go through before reaching grpc server methods. We read the authorization header (which is a JWt). then we verify it and check for admin role.

```python
def intercept_service(self, continuation, handler_call_details):
        # Exclude reflection service from authentication so that it can retrieve the service info
        if handler_call_details.method == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo":  
            return continuation(handler_call_details)
            
        # Extract the token from the gRPC request headers
        metadata = dict(handler_call_details.invocation_metadata)
        token = metadata.get('authorization', '')
        if token:
            logging.info(f"Token found: {token}")
            public_key = jwk.JWK.from_pem(self.cfg.AuthPublicKey)
            header, claims = jwt.verify_jwt(token, public_key, ['RS256'], checks_optional=True,ignore_not_implemented=True)
            roleStr = claims['role'].lower()
            if roleStr == "admin":
                logging.info(" role is Admin: Authorized")
                return continuation(handler_call_details)
            return self._notAuthorized
        else:
            logging.info("No token found")
            return self._noAuthHeader
```

> - Create a file named \_\_init\_\_.py inside `grpc_interceptors` and set the contents to 

```python
from app.grpc_interceptors._token_validation_interceptor import TokenValidationInterceptor
```
> - Change the server definition inside server init file to match the following:

```python
        from app.grpc_interceptors import TokenValidationInterceptor
        ...

        server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
            interceptors = (TokenValidationInterceptor(cfg),)
        )
```

> - Run `python main.py`. Now an auth interceptor in active and protect our grpc server. If you invoke any method, the response would be `Missing Authorization header`
> - go to [this site](https://dinochiesa.github.io/jwt/) and put your generated rsa key-pairs in the corresponding boxes. Then define a simple jwt like this:

```json
{
  "iss": "microservice-scheduler.com",
  "sub": "sheniqua",
  "aud": "maxine",
  "iat": 1689904636,
  "exp": 1689910603,
  "role": "admin"
}
```
> - Set other configs as shown in the following image and then click the left arrow. Copy the encoded token to the clipboard.


![JWT generation](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/p6utfo7nic4fyvv3aqa8.PNG)

> - Now go to docker desktop and restart `grpcui-service` then go to http://localhost:5000/. For each request you make, set a metadata with name= authorization and the value= `<token_copied>`
> - This time when you invoke the methods, authentication will pass. you can change the role to `user` and regenerate the token and use the new token. This time the response would be `Not authorized` happy debugging :)

---


## To DO

> - Add tests
> - Add tracing using Jaeger
> - Add monitoring and analysis using grafana
> - Refactoring

---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).