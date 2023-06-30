This is the 3rd part of a series of articles under the name **"Play Microservices"**. Links to other parts:
Part 1: [Play Microservices: Bird's eye view](https://dev.to/khaledhosseini/play-microservices-birds-eye-view-3d44)
Part 2: [Play Microservices: Authentication](https://dev.to/khaledhosseini/play-microservices-authentication-4di3)
Part 3: You are here

The source code for the project can be found [here](https://github.com/KhaledHosseini/play-microservices):

---

## Contents:

 - **Summary**
 - **Tools**
 - **Docker dev environment**
 - **Database service: Mongodb**
 - **Mongo express service**
 - **Kafka metadata service: Zookeeper**
 - **Zoonavigator service**
 - **Message broker service: Kafka**
 - **Kafka-ui service**
 - **Scheduler-grpcui-service**
 - **Scheduler service: Golang**
 - **To do**

---

- **Summary**

In the 2nd part we have developed an authentication service. Here our goal is to create a job scheduler service for our microservices application. To accomplish our goal, we require four separate services: a database service, a message broker service , a metadata database service for the message broker and the scheduler service which is a gRPC api service. In the development environment, we include four additional services for debugging purposes. These services are Mongo express to manage our database service, Kafkaui service to manage our kafka service, Zoonavigator for Zookeeper service and grpcui service to test our gRPC api.


![Summary](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/cqdkr25890v08auqfdiz.png)

At the end, the project directory structure will appear as follows:


![Folder structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/3p0r2lpw5i4mdd19je0k.PNG)
---
- **Tools**

The tools required In the host machine:

  - [Docker](https://www.docker.com/): Containerization tool
  - [VSCode](https://code.visualstudio.com/): Code editing tool
  - [Dev containsers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension for VSCode
  - [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) extension for VSCode
  - [Git](https://git-scm.com/)

The tools and technologies that we will use Inside containers for each service:

 - Database service: [Mongo](https://hub.docker.com/_/mongo)
 - mongo express service: [Mongo express](https://hub.docker.com/_/mongo-express)
 - Messaging service: [Kafka](https://hub.docker.com/r/bitnami/kafka/)
 - Kafka-ui service: [Kafka-ui](https://hub.docker.com/r/provectuslabs/kafka-ui)
 - metadata service: [Zookeeper](https://hub.docker.com/_/zookeeper)
 - Zoonavigator service: [Zoonavigator](https://hub.docker.com/r/elkozmon/zoonavigator)
 - grpcui service: [grpcui](https://hub.docker.com/r/fullstorydev/grpcui)
 - Scheduler api service: 
  - [Golang](https://go.dev/) : programming language
  - [gRPC-GO](https://pkg.go.dev/google.golang.org/grpc): gRPC framework for golang
  - [Mongo driver](https://pkg.go.dev/go.mongodb.org/mongo-driver): Query builder for our database communication.
  - [Kafka go](https://pkg.go.dev/github.com/segmentio/kafka-go) for our message broker communications from go.
  - [Quartz ](github.com/reugn/go-quartz) for scheduling purposes.

---

- **Docker dev environment**

Development inside Docker containers can provide several benefits such as consistent environments, isolated dependencies, and improved collaboration. By using Docker, development workflows can be containerized and shared with team members, allowing for consistent deployments across different machines and platforms. Developers can easily switch between different versions of dependencies and libraries without worrying about conflicts.

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/id0cijn8hm77ohn0hk2j.png)

When developing inside a Docker container, you only need to install `Docker`, `Visual Studio Code`, and the `Dev Containers` and `Docker `extensions on VS Code. Then you can run a container using Docker and map a host folder to a folder inside the container, then attach VSCode to the running container and start coding, and all changes will be reflected in the host folder. If you remove the images and containers, you can easily start again by recreating the container using the Dockerfile and copying the contents from the host folder to the container folder. However, it's important to note that in this case, any tools required inside the container will need to be downloaded again. Under the hood, When attaching VSCode to a running container, Visual Studio code install and run a special server inside the container which handle the sync of changes between the container and the host machine.

---

- **Database service: Mongo**

> - Create a folder for the project and choose a name for it (such as 'microservice'). Then create a folder named `scheduler`. This folder is the root directory of scheduler service. You can then open the root folder in VS Code by right-clicking on the folder and selecting 'Open with Code'.
> - Inside the root directory create a folder with the name scheduler-db-service, then create the following files inside.
> - Create a Dockerfile and set content to `FROM mongo:7.0.0-rc5`
> - Create a file named pass.txt and set content to a `password`
> - Create a file named user.txt and set content to `admin `
> - Create a file named db_name.txt and set content to `jobs_db`
> - Create a file named .env in the root directory and set the content to `MONGODB_PORT=27017`.
> - Inside root directory create a file named docker-compose.yml and add the following content.

```yaml
version: '3'
services:
# database service for scheduler service
  scheduler-db-service:
    build: 
      context: ./scheduler-db-service
      dockerfile: Dockerfile
    container_name: scheduler-db-service
    environment:
      MONGO_INITDB_ROOT_USERNAME_FILE: /run/secrets/scheduler-db-user
      MONGO_INITDB_ROOT_PASSWORD_FILE: /run/secrets/scheduler-db-pass
    env_file:
      - ./scheduler-db-service/.env
    ports:
      - ${MONGODB_PORT}:${MONGODB_PORT}
    secrets:
      - scheduler-db-user
      - scheduler-db-pass
      - scheduler-db-dbname 
    volumes:
      -  scheduler-db-service-VL:/data/db

volumes:
  scheduler-db-service-VL:

secrets:
  scheduler-db-user:
    file: scheduler-db-service/user.txt
  scheduler-db-pass:
    file: scheduler-db-service/pass.txt
  scheduler-db-dbname:
    file: scheduler-db-service/db_name.txt
```

> - In our Docker Compose file, we use secrets to securely share credential data between containers. While we could use an .env file and environment variables, this is not considered safe. When defining secrets in the Compose file, Docker creates a file inside each container (which has the secrets name) under the `/run/secrets/` path, which the containers can then read and use. For example, we will set the path of the Docker Compose secret `scheduler-db-pass` to the `DATABASE_PASS_FILE` environment variable of scheduler service. The service then will go to the path (/run/secrets/scheduler-db-pass) and read the password file. We will be using these secrets in other services later in the project.

---

- **Mongo express service**

The purpose of this service is solely for debugging and management of our running database server in the development environment.

> - Inside root directory create a folder with the name mongo-express-service
> - Create a Dockerfile and set content to `FROM mongo-express:1.0.0-alpha.4`
> - Create a file named .env beside Dockerfile and set the content to 

```
ME_CONFIG_BASICAUTH_USERNAME=admin
ME_CONFIG_BASICAUTH_PASSWORD=password123
ME_CONFIG_MONGODB_ENABLE_ADMIN=true
```

> - Add the following lines to the .env file of the docker-compose (the .env file at the root directory of the project.)

```
MONGO_EXPRESS_PORT=8081
```
> - Add the following to the service part of the docker-compose.yml.

```yaml
    mongo-express:
    build:
      context: ./mongo-express-service
      dockerfile: Dockerfile
    container_name: mongo-express-service
    restart: always
    environment:
      - ME_CONFIG_MONGODB_PORT=${MONGODB_PORT}
      - ME_CONFIG_MONGODB_SERVER=scheduler-db-service
      - ME_CONFIG_MONGODB_ADMINUSERNAME=root
      - ME_CONFIG_MONGODB_ADMINPASSWORD=password123
    env_file:
      - ./mongo-express-service/.env
    ports:
      - ${MONGO_EXPRESS_PORT}:${MONGO_EXPRESS_PORT}
    depends_on:
      - scheduler-db-service
```
> - Here, as mongo express doesn't provide a capability to read database password from the files, we simply pass mongodb credentials using environment variables (Do not forget we are in development environment).

> - Now open a terminal in your project directory and run docker-compose up. Docker Compose will download and cache the required images before starting your containers. For the first run, this may take a couple of minutes. If everything goes according to plan, you can then access the mongo express panel at http://localhost:8081/ and log in using the mongo express credentials from the .env file inside the mongo-express-service container. You should see that it has successfully connected to the scheduler-db-service container.


![Mongo express](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/d687yoms8vxa5h5rlnbu.PNG)

> - Now run docker-compose down

---

- **Metadata service: Zookeeper**

[ZooKeeper ](https://zookeeper.apache.org/) is a centralized service for maintaining configuration information. we use it as metadata storage for our Kafka messaging service.
> - Inside root directory create a folder with the name zookeeper-service
> - Create a Dockerfile and set content to `FROM bitnami/zookeeper:3.8.1`
> - Create a file named .env and set content to 

```
ZOO_SERVER_USERS=admin,user1
# for development environment only
ALLOW_ANONYMOUS_LOGIN="yes"
# if yes, uses SASL
ZOO_ENABLE_AUTH="no" 
```
> - Create a file named server_passwords.properties and set content to `password123,password_for_user1` Please choose your own passwords.
> - Add the following to the .env file of the docker-compose (the .env file at the root directory of the project.)

```
ZOOKEEPER_PORT=2181
ZOOKEEPER_ADMIN_CONTAINER_PORT=8078
ZOOKEEPER_ADMIN_PORT=8078
```
> - Add the following to the service part of the docker-compose.yml.

```yaml
  zk1:
    build:
      context: ./zookeeper-service
      dockerfile: Dockerfile
    container_name: zk1-service
    secrets:
      - zoo-server-pass
    env_file:
      - ./zookeeper-service/.env
    environment:
      ZOO_SERVER_ID: 1
      ZOO_SERVERS: zk1:${ZOOKEEPER_PORT}:${ZOOKEEPER_PORT} #,zk2:{ZOOKEEPER_PORT}:${ZOOKEEPER_PORT}
      ZOO_SERVER_PASSWORDS_FILE: /run/secrets/zoo-server-pass
      ZOO_ENABLE_ADMIN_SERVER: yes
      ZOO_ADMIN_SERVER_PORT_NUMBER: ${ZOOKEEPER_ADMIN_CONTAINER_PORT}
    ports:
      - '${ZOOKEEPER_PORT}:${ZOOKEEPER_PORT}'
      - '${ZOOKEEPER_ADMIN_PORT}:${ZOOKEEPER_ADMIN_CONTAINER_PORT}'
    volumes:
      - "zookeeper_data:/bitnami"
```
> - Add the following to the secrets part of the docker-compose.yml.

```yaml
  zoo-server-pass:
    file: zookeeper-service/server_passwords.properties
```
> -ZooKeeper is a distributed application that allows us to run multiple servers simultaneously. It enables multiple clients to connect to these servers, facilitating communication between them. ZooKeeper servers collaborate to handle data and respond to requests in a coordinated manner. In this case, our zookeeper consumers (clients) are Kafka servers which is again a distributed event streaming platform. We can run multiple zookeeper services as an ensemble of zookeeper servers and attach them together via `ZOO_SERVERS` environment variable.
> - The Bitnami ZooKeeper Docker image provides a zoo_client entrypoint, which acts as an internal client and allows us to run the zkCli.sh command-line tool to interact with the ZooKeeper server as a client. But we are going to use a GUI client for debugging purposes: Zoonavigator.

---

 - **Zoonavigator service**

This service exists only in the development environment for debugging purposes. We use it to connect to zookeeper-service and manage the data.

> - Inside root directory create a folder with the name zoonavigator-service
> - Create a Dockerfile and set content to `FROM elkozmon/zoonavigator:1.1.2`
> - Add `ZOO_NAVIGATOR_PORT=9000` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  zoo-navigator:
    build: 
      context: ./zoonavigator-service
      dockerfile: Dockerfile
    container_name: zoo-navigator-service
    ports:
      - '${ZOO_NAVIGATOR_PORT}:${ZOO_NAVIGATOR_PORT}'
    environment:
      - CONNECTION_LOCALZK_NAME = Local-zookeeper
      - CONNECTION_LOCALZK_CONN = localhost:${ZOOKEEPER_PORT}
      - AUTO_CONNECT_CONNECTION_ID = LOCALZK
    depends_on:
      - zk1
``` 

> - Now from the terminal run `docker-compose up -d --build`
> - While running go to `http://localhost:9000/`. You will see the following screen:


![zoonavigatoe](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/h151vcjvlc60lp6hr7lq.PNG)

> - Enter the container name of a zookeeper service (here zk1). If everything go according to plan you can connect to the zookeeper service.

![zoonavigator-zk1](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/aldo17byoo9ef9uppbcc.PNG) 

> - No run docker-compose down. We will return to these tools later.

---

- **Message broker service: Kafka**

[Apache Kafka](https://kafka.apache.org/)  is an open-source distributed event streaming platform that is well-suited for Microservices architecture. It is an ideal choice for implementing patterns such as event sourcing. Here We use it as an message broker for our scheduler service. 

> - Inside root directory create a folder with the name kafka-service
> - Create a Dockerfile and set content to `FROM bitnami/kafka:3.4.1`
> - Create a .env file beside the Docker file and set the content to: 

```
ALLOW_PLAINTEXT_LISTENER=yes
KAFKA_ENABLE_KRAFT=no
```
> - Add `KAFKA1_PORT=9092` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  kafka1:
    build: 
      context: ./kafka-service
      dockerfile: Dockerfile
    container_name: kafka1-service
    ports:
      - '${KAFKA1_PORT}:${KAFKA1_PORT}'
    volumes:
      - "kafka_data:/bitnami"
    env_file:
      - ./kafka-service/.env
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ADVERTISED_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka1:${KAFKA1_PORT},LISTENER_DOCKER_EXTERNAL://${DOCKER_HOST_IP:-127.0.0.1}:${KAFKA1_PORT}
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: LISTENER_DOCKER_INTERNAL:PLAINTEXT,LISTENER_DOCKER_EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: LISTENER_DOCKER_INTERNAL
      KAFKA_CFG_ZOOKEEPER_CONNECT: zk1:${ZOOKEEPER_PORT}
      KAFKA_ZOOKEEPER_PROTOCOL: PLAINTEXT #if auth is enabled in zookeeper use one of: SASL, SASL_SSL see https://hub.docker.com/r/bitnami/kafka
      KAFKA_CFG_LISTENERS: PLAINTEXT://:${KAFKA1_PORT}
    depends_on:
      - zk1
```
> - In order to connect to our Kafka brokers for debugging purposes, we run another service. Kafka-ui.


---

 - **Kafka-ui service**

This service exists only in the development environment for debugging purposes. We use it to connect to kafka-service and manage the data.

> - Inside root directory create a folder with the name kafkaui-service
> - Create a Dockerfile and set content to `FROM provectuslabs/kafka-ui:latest`
> - Add `KAFKAUI_PORT=8080` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  kafka-ui:
    build: 
      context: ./kafkaui-service
      dockerfile: Dockerfile
    container_name: kafka-ui-service
    restart: always
    ports:
      - ${KAFKAUI_PORT}:${KAFKAUI_PORT}
    environment:
     KAFKA_CLUSTERS_0_NAME: local
     KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka1:${KAFKA1_PORT}
     DYNAMIC_CONFIG_ENABLED: 'true'
    depends_on:
      - kafka1
```

> - Now run `docker-compose run -d --build`. While containers are running, go to `http://localhost:8080/` to open Kafka-ui dashboard. 


![Kafka-ui](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/gcwqxf16jbig0q3xxzn7.PNG)

> - As you can see, you can view and manage brokers, topics and consumers. We will return to these items later. 
> - Run `docker-compose down`
> - Our required services are ready and running. Now it is time to Prepare development environment for our scheduler service. 

![Coding time](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/8hhqfd5dimjlf8d54t3b.jpg)

---

- **Scheduler-grpcui-service**

Before starting our Scheduler service development, lets add another service to our development environment which we will use to interact with our scheduler-service for debugging purposes.
> - Create a folder named grpcui-service inside scheduler folder.
> - Create a Docker file and set contents to `FROM fullstorydev/grpcui:v1.3.1`
> - Add the following to the services part of the docker-compose.yml file. 

```yaml
  grpcui-service:
    build:
      context: ./grpcui-service
      dockerfile: Dockerfile
    container_name: grpcui-service
    command: -port $GRPCUI_PORT -plaintext scheduler-service:${SCHEDULER_PORT}
    restart: always
    ports:
      - ${GRPCUI_PORT}:${GRPCUI_PORT}
    depends_on:
      - scheduler-service
```
> - Add `GRPCUI_PORT=5000` to compose .env file. (.env file beside docker-compose)
> - We will return to this service later.

---

- **Scheduler service: Golang**

Our goal is to develop a gRPC server with Go. The typical pipeline for developing a gRPC server is quite straightforward. You define your gRPC schema inside a .proto file (see [here](https://protobuf.dev/) for more info). Then you compile (Actually you transform) the .proto to your target programming language using a protocol buffer compiler tool and import it to your project. .proto models are are business layer models. Then you use a gRPC framework in your target language to run a gRPC server. Then you can define corresponding database layer models and use a converter to transform between them. You receive gRPC models vis gRPC server, convert them to database models and store them in a database. In case of queries, you query the data from the database, transform them to gRPC models and return them to the user.

Here is a summary of what we are going to do: We first install [protoc](https://github.com/protocolbuffers/protobuf/releases/download/v23.3/protoc-23.3-linux-x86_64.zip) in our development environment. then initial our go project. Define our proto scheme and compile it using the above tool and then running an initial gRPC server. Then we add database layer models and classes. 

> - Create a folder named `scheduler-service` inside scheduler folder.
> - Create a Dockerfile inside `scheduler-service` and set the contents to

```bash
FROM golang:1.19

ENV PROTOC_VERSION=23.3
ENV PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-x86_64.zip
RUN apt-get update && apt-get install -y unzip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \ 
    && rm -f $PROTOC_ZIP
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN export PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /usr/src/app
```
> - Add the following to the service part of our docker-compose.yml file.

```yaml
  scheduler-service:
    build: 
      context: ./scheduler-service
      dockerfile: Dockerfile
    container_name: scheduler-service
    command: sleep infinity
    environment:
      ENVIRONMENT: development
      SERVER_PORT: ${SCHEDULER_PORT}
      DATABASE_USER_FILE: /run/secrets/scheduler-db-user
      DATABASE_PASS_FILE: /run/secrets/scheduler-db-pass
      DATABASE_DB_NAME_FILE: /run/secrets/scheduler-db-dbname
      DATABASE_SCHEMA: mongodb
      DATABASE_HOST_NAME: scheduler-db-service
      DATABASE_PORT: ${MONGODB_PORT}
      KAFKA_BROKERS: kafka1-service:${KAFKA1_PORT}
    ports:
      - ${SCHEDULER_PORT}:${SCHEDULER_PORT}
    volumes:
      - ./scheduler-service:/usr/src/app
    secrets:
      - scheduler-db-user
      - scheduler-db-pass
      - scheduler-db-dbname
```
> - We are going to do all the development inside a docker container without installing Golang in our host machine. To do so, we run the containers and then attach to scheduler-service container using VSCode. As you may noticed, the dockerfile for scheduler-service has no entry-point therefore we set the command value of scheduler-service to `sleep infinity` to stay the container awake.
> - Now run `docker-compose up -d --build`
> - While running, attach to the scheduler service by clicking bottom-left icon and then select `attach to running container `. Select scheduler-service and wait for a new instance of VSCode to start. At the beginning the VScode asks us to open a folder inside the container. We have selected  `WORKDIR /usr/src/app` inside our dockerfile, so we will open this folder inside the container. this folder is mounted to scheduler-service folder inside the host machine using docker compose file, therefor whatever change we made will be synced to the host folder too.
> - After opening the folder `/usr/src/app`, open a new terminal and initialize the go project by running `go mod init github.com/<your_username>/play-microservices/scheduler/scheduler-service`. This command will create a go.mod file.
> - Run `go get -u google.golang.org/grpc`. This is a gRPC framework for running grpc server using Golang.
> - Run `go get -u google.golang.org/grpc/reflection`. We ad reflection to our gRPC server so that our grpcui-service can connect to it and retrieve the endpoints and messages easily for debugging purposes.
> - Now create a folder named proto and create a file named job.proto inside. Set the content from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/proto/job.proto)
> - Run `protoc --go_out=./proto --go-grpc_out=./proto proto/*.proto`. this will compile our .proto file to Golang. two files will be created. job.pb.go and job_grpc.pb.go. The first contains the proto models and the second contains the code for job service interface (We need to create our service and implement this interface).
> - Note: We adopt a Golang project structure that aligns with the recommended guidelines stated [here](https://github.com/golang-standards/project-layout)
> - Create a folder named config and a file named `config.go`. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/config/config.go). Also create a file named .env in the same folder. we will put our internal environment variables here. Set contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/config/.env). Create another file named .env.topics for putting kafka topics. For production environment, we pass the topics file via docker compose secrets and send the position of file via `TOPICS_FILE` environment variable. Then we load the contents of the file. 
> - Create a folder named pkg in the root directory (beside mod.go). We will put general packages here. Inside Create a folder named logger then a file named logger.go and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/pkg/logger/logger.go).
> - Create this folder tree: `internal/models/job/grpc`. inside grpc folder create a folder named job_service.go and set the contents to

```go
package grpc

import (
	context "context"
	pb "github.com/<your_username>/play-microservices/scheduler/scheduler-service/proto"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type JobService struct {
pb.UnimplementedJobServiceServer
}

func NewJobService() *JobService {
	return &JobService{}
}

func (j *JobService) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateJob not implemented")
}

func (JobService) GetJob(context.Context, *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJob not implemented")
}

func (JobService) ListJobs(context.Context, *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListJobs not implemented")
}

func (JobService) UpdateJob(context.Context, *pb.UpdateJobRequest) (*pb.UpdateJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateJob not implemented")
}

func (JobService) DeleteJob(context.Context, *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteJob not implemented")
}

```

> - Create a folder named server inside internal folder. Then a file named server.go. set the content to

```go
package server

import (
	"log"
	"net"

	"github.com/<your_username>/play-microservices/scheduler/scheduler-service/config"
"github.com/<your_username>/play-microservices/scheduler/scheduler-service/pkg/logger"
	MyJobGRPCService "github.com/<your_username>/play-microservices/scheduler/scheduler-service/internal/models/job/grpc"
	JobGRPCServiceProto "github.com/<your_username>/play-microservices/scheduler/scheduler-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
        log       logger.Logger
	cfg       *config.Config
}

// NewServer constructor
func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg}
}

func (s *server) Run() error {
	lis, err := net.Listen("tcp", ":"+s.cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpc_server := grpc.NewServer()

	job_service := MyJobGRPCService.NewJobService()
	JobGRPCServiceProto.RegisterJobServiceServer(grpc_server, job_service)
	reflection.Register(grpc_server)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpc_server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

```
> - Create a folder named cmd and a file named main.go. Set the content to 

```go
package main

import (
	"log"

	"github.com/<your_username>/play-microservices/scheduler/scheduler-service/config"
"github.com/<your_username>/play-microservices/scheduler/scheduler-service/pkg/logger"
	"github.com/<your_username>/play-microservices/scheduler/scheduler-service/internal/server"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting user server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Environment: %s",
		cfg.AppVersion,
		cfg.Logger_Level,
		cfg.Environment,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	s := server.NewServer(appLogger, cfg)
	s.Run()
}
```
> - Run `go mod tidy`
> - Run `go run cmd/main.go`
> - While our server is running, go to docker desktop and restart grpcui-service. Now go to `http://localhost:5000/`. If everything goes on plan, you can connect the server.


![grpcui-start](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/rkff8ngf3rrevkkspbw1.PNG)

> - Invoking any of the methods results in `method xxx not implemented` because we still have not implemented our job service methods.
> - return to VSCode instance that is already attached to our scheduler-service. Stop the service by hitting `ctl + c`
> - Create a file named job.go inside models folder. Inside this file we are going to define database layer models corresponding to proto models. Logic for converting from/to proto models comes here. Also we define an interface for the database of our job model and an interface for the event messaging of our job model. Set the contents of the job.go from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/internal/models/job.go).

```go
// databas interface for Job model
type JobDB interface {
	Create(ctx context.Context, job *Job) (*Job, error)
	Update(ctx context.Context, job *Job) (*Job, error)
	GetByID(ctx context.Context, jobID primitive.ObjectID) (*Job, error)
	DeleteByID(ctx context.Context, jobID primitive.ObjectID) (error)
	GetByScheduledKey(ctx context.Context, jobScheduledKey int) (*Job, error)
	DeleteByScheduledKey(ctx context.Context, jobScheduledKey int) (error)
	ListALL(ctx context.Context, pagination *utils.Pagination) (*JobsList, error)
}

//Message broker interface for Job model
type JobsProducer interface {
	PublishCreate(ctx context.Context, job *Job) error
	PublishUpdate(ctx context.Context, job *Job) error
	PublishRun(ctx context.Context, job *Job) error
}
```
> - We then pass these two interfaces to our JobService model and do our logic there without knowing which database engine or messaging broker we are using. This gives us the flexibility to select whatever database (mongo or postgres) or message broker (kafka or rabitmq) we want.
> - For scheduling we use this [package](https://github.com/reugn/go-quartz). Run `go get github.com/reugn/go-quartz/quartz`
> - Now change the definition of JobService inside job_service.go: Final file is [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/internal/models/job/grpc/job_grpc_service.go)

```go
type JobService struct {
	jobDB         models.JobDB
	jobsProducer  models.JobsProducer
	jobsScheduler scheduler.Scheduler
        pb.UnimplementedJobServiceServer
}

func NewJobService(jobDB models.JobDB, jobsProducer models.JobsProducer, jobsScheduler scheduler.Scheduler) *JobService {
	return &JobService{jobDB: jobDB, jobsProducer: jobsProducer, jobsScheduler: jobsScheduler}
}
```
> - Now it is time to select a database engine and implement `JobDB` interface using it. Run `go get go.mongodb.org/mongo-driver`. Create a folder named `database` inside  `internal/models/job` folder. Then a file named job_db_mongo.go. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/internal/models/job/database/job_db_mongo.go)
> - Create a folder named mongodb inside pkg directory and then a file named mongodb.go. This package is used to initialize our mongo db database. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/pkg/mongodb/mongodb.go).
> - Now we select a message broker framework and implement `JobsProducer ` interface using it. Run `go get github.com/segmentio/kafka-go`. Create a folder named message_broker inside `internal/models/job` folder. Then create a filed named `job_producer_kafka.go` here we implement `JobsProducer ` interface. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/internal/models/job/message_broker/job_producer_kafka.go).
> - Create a folder named kafka inside pkg and then a file named kafka.go. This package is used to initialize kafka connection. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/pkg/kafka/kafka.go).
> - Some notes on kafka architecture: 
>  - Topics: Core abstraction in kafka and represent stream of records.
   - Partitions: Topics which are a chain of records can be divided to partitions to enable parallel processing and scalability. Partitions can be distributed among brokers or reside only in one broker.
   - Brokers: Brokers are the Kafka servers that form the cluster. They store and manage the published records in a distributed manner.
   - Replication: Kafka provides replication of data for fault tolerance. Each partition can have multiple replicas spread across different brokers. Replicas ensure that if a broker fails, another broker can take over and continue serving the data seamlessly.
   - Producers: Producers are responsible for publishing data to Kafka topics. 
   - Consumers: Consumers read data from Kafka topics. Consumers belongs to consumers groups. Each consumer can subscribe to multiple topics but only for **one partition** of a topic.
> - Configuring a kafka environment can be tricky and depends on use cases.

![Kafka](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/vwk64gi7yn9ltyy2y62o.png)

> - Now we reiterate through main.go, server.go and job_service.go files and complete the contents. Change the contents of main.go from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/cmd/main.go) and server.go from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/internal/server/server.go) to include mongo database and kafka brokers. We initialize them inside main and pass them to the server struct. 
> - Add the remaining packages from [pkg](https://github.com/KhaledHosseini/play-microservices/tree/master/scheduler/scheduler-service/pkg) folder. Install the required packages via `go get <packagename>` command.

```
"github.com/pkg/errors"
```

> - Set the contents for job_service.go from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/internal/models/job/grpc/job_grpc_service.go). The code for the CreateJob function is as follows:

```go
func (j *JobService) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {

	job := models.JobFromProto_CreateJobRequest(req)
	job.Status = int32(pb.JobStatus_SCHEDULED)
	jobFingerPrint := fmt.Sprintf("%s:%s:%s:%p", job.Name, job.Description, job.JobData, &job.ScheduleTime)
	job.ScheduledKey = int(fnv1a.HashString64(jobFingerPrint)) //We assume 64 bit systems!
	//cretae job in the database
	created, err := j.jobDB.Create(ctx, job)
	if err != nil {
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	jobID := created.JobID
	functionJob := jobFucntion.NewFunctionJobWithKey(job.ScheduledKey, func(_ context.Context) (int, error) {

		jb, err := j.jobDB.GetByID(ctx, jobID)
		if err != nil {
			return 0, grpcErrors.ErrorResponse(err, err.Error())
		}
		//set the state of the job to running
		jb.Status = int32(pb.JobStatus_RUNNING)
		_, err2 := j.jobDB.Update(ctx, jb)
		if err2 != nil {
			return 0, grpcErrors.ErrorResponse(err2, err2.Error())
		}

		j.jobsProducer.PublishRun(ctx, job)

		return 0, nil
	})

	triggerTime := time.Duration(time.Until(job.ScheduleTime).Seconds())
	j.jobsScheduler.ScheduleJob(ctx, functionJob, scheduler.NewSimpleTrigger(triggerTime))

	j.jobsProducer.PublishCreate(ctx, created)

	return &pb.CreateJobResponse{Id: created.JobID.String()}, nil
}
```
> - We receive createJob in our gRPC server. We first convert .proto model to database layer model in `job := models.JobFromProto_CreateJobRequest(req)`. Then we set the status to SCHEDULED and save it to the database and retrieve the id. We then schedule the job and publish the job created topic to be consumed by other services like reports service and finally we return the CreateJobResponse to the user. Inside the schedule functuin, We retrieve the job from database, then set the jobState to RUNNING and save it again to the database. Then we publish run topic to be consumed by our job runner service. We then listen to run-update event which will be triggered by job runner. If the result was success we change the state of of our job to COMPLETE save it to the database and then publish job update to be recorded by our report service.

> - Run `go mod tidy`
> - Run `go run cmd/main.go`
> - Now go to docker desktop an restart grpcui service. then go to `http://localhost:5000/`. If everything goes according to plan, you can connect to the service. Select CreateJob from method name and fill in the form. For Schedule time if you pass a time before time.Now() the schedule will trigger immediately. For job data depending on the job type we need to send a specific Json string. For email type, we need to send a json with the following structure:

```json
{
	"SourceAddress": "example@example.com",
	"DestinationAddress": "example@example.com",
	"Subject": "Message From example@example.com contact form",
	"Message": "This is a production test!!!!"
}
```
> - Push Invoke button. You will receive created job id.
> - Go to `http://localhost:8081/` and check the database.

![Mongo test](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/azur9yxh83jmj39gpqgb.PNG)

> - Now go to `http://localhost:8080/`. You can see that 3 topics have been created. The number of messages for `topic-job-create` and `topic-job-run` is 1. Because we have published `topic-job-create` once and `topic-job-run` has been published inside our scheduled function.


![kafka topics](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/pcpj95xmmlse9gd8y9qe.PNG)

> - Stop the service by hitting `ctl + c`
> - A note on kafka topics creation: There are generally two common approaches for topic creation in a microservice architecture. 
   - Centralized Topic Creation: a dedicated team or infrastructure administrators are responsible for creating and managing Kafka topics.
   - Self-Service Topic Creation: In this approach, each microservice is responsible for creating and managing its own Kafka topics. Here, In the development environment we follow this approach and the producer of a topic is responsible for topic creation. We define topic names in .env.topic file for development environment and load them from that file. For production environment, we can receive the file path for the topics using TOPICS_FILE environment variables (location of the file from docker compose secrets) and then load the data from it.
> - We need to listen to `topic-job-run` results. To accomplish this we have to subscribe for `topic-job-run-result` topic. To accomplish that, create a file named `job_consumer_kafka.go` in the path internal/models/job/message_broker. Then set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/internal/models/job/message_broker/job_consumer_kafka.go).
> - Add the following to server.go file.

```
jobsConsumer := kafka.NewJobsConsumerGroup(s.log, job_db)
jobsConsumer.Run(ctx, cancel, s.kafkaConn, s.cfg)
```
> - Run `go mod tidy` then run `go run cmd/main.go`
> - Go to `http://localhost:8080/` and under consumers you can see `job-run-result-consumer` which is the group id of our consumer. Now go to Topics-> topic-job-create and in the messages tab click the preview of the value of the message and copy the json structure of the message. Something like this (You need to copy yours because when listening, we search the database for jobId):

```json
{
	"jobId": "649f07e619fca8aa63d842f6",
	"name": "job1",
	"scheduleTime": "1970-01-01T00:00:00Z",
	"createdAt": "2023-06-30T16:50:46.3042083Z",
	"updatedAt": "2023-06-30T16:50:46.3042086Z",
	"status": 2,
	"scheduledKey": 6072375870331110000
}
```
> - Set the value of job status of json string to 4. 
> - Now go to topics -> topic-job-run-result and click produce message from top right corner. Paste the json string into value and click on produce message. If everything goes according to plan, you can see the results from the terminal of VSCode. Also you can go to mongo express at http://localhost:8081/ and see that job status has changed to 4.


---


- **To DO**

> - Add tests
> - Add tracing using Jaeger
> - Add monitoring and analysis using grafana
> - Refactoring

---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).