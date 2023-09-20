Contents:
- [Summary](#summary)
- [Tools](#tools)
- [Docker dev environment](#docker-dev-environment)
- [API mock service: Gripmock](#api-mock-service-gripmock)
- [Api-gateway service: Golang](#api-gateway-service-golang)
- [To DO](#to-do)


---

This is the 6th part of a series of articles under the name **"Play Microservices"**. Links to other parts:<br/>
Part 1: [Play Microservices: Bird's eye view](https://dev.to/khaledhosseini/play-microservices-birds-eye-view-3d44)<br/>
Part 2: [Play Microservices: Authentication](https://dev.to/khaledhosseini/play-microservices-authentication-4di3)<br/>
Part 3: [Play Microservices: Scheduler service](https://dev.to/khaledhosseini/play-microservices-scheduler-19km)<br/>
Part 4: [Play Microservices: Email service](https://dev.to/khaledhosseini/play-microservices-email-service-1kmc)<br/>
Part 5: [Play Microservices: Report service](https://dev.to/khaledhosseini/play-microservices-report-service-4jcm)<br/>
Part 6: You are here<br/>
Part 7: [Play Microservices: Client service](https://dev.to/khaledhosseini/play-microservices-client-service-4jbf)<br/>
Part 8: [Play Microservices: Integration via docker-compose](https://dev.to/khaledhosseini/play-microservices-integration-via-docker-compose-2ddc)<br/>
Part 9: [Play Microservices: Security](https://dev.to/khaledhosseini/play-microservices-security-45e4)<br/>

---
## Summary

In the previous stages, we successfully developed a set of services consisting of the auth, scheduler, email, and report services. Our current objective is to establish a gateway service that acts as a single entry point for clients to access these individual services. To ensure smooth development and testing, we have included an additional service in our development environment solely for debugging purposes. As we independently develop the API gateway, the remaining services are temporarily unavailable. To overcome this limitation, we create mock implementations to simulate the behavior of the unavailable services during the development of the API gateway.


![Summary](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/ymno6abm5m5b07rkucgv.png)


At the end, the project directory structure will appear as follows:



![Folder structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/yijlaztrbed8yw0k2da3.PNG)

---

## Tools

The tools required In the host machine:

  - [Docker](https://www.docker.com/): Containerization tool
  - [VSCode](https://code.visualstudio.com/): Code editing tool
  - [Dev containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension for VSCode
  - [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) extension for VSCode
  - [Git](https://git-scm.com/)

The tools and technologies that we will use Inside containers for each service:

 - Gripmock service: [Gripmock](https://registry.hub.docker.com/r/tkpd/gripmock)
 - Api-gateway service: 
  - [Golang](https://go.dev/) : programming language
  - [gRPC-GO](https://pkg.go.dev/google.golang.org/grpc): gRPC framework for golang
  - [Gin](https://gin-gonic.com/):  Is a web framework written in Golang
  - [gin-swagger](https://github.com/swaggo/gin-swagger) for our rest api documentation.

---

## Docker dev environment

Development inside Docker containers can provide several benefits such as consistent environments, isolated dependencies, and improved collaboration. By using Docker, development workflows can be containerized and shared with team members, allowing for consistent deployments across different machines and platforms. Developers can easily switch between different versions of dependencies and libraries without worrying about conflicts.

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/id0cijn8hm77ohn0hk2j.png)

When developing inside a Docker container, you only need to install `Docker`, `Visual Studio Code`, and the `Dev Containers` and `Docker `extensions on VS Code. Then you can run a container using Docker and map a host folder to a folder inside the container, then attach VSCode to the running container and start coding, and all changes will be reflected in the host folder. If you remove the images and containers, you can easily start again by recreating the container using the Dockerfile and copying the contents from the host folder to the container folder. However, it's important to note that in this case, any tools required inside the container will need to be downloaded again. Under the hood, When attaching VSCode to a running container, Visual Studio code install and run a special server inside the container which handle the sync of changes between the container and the host machine.

---

## API mock service: Gripmock

We have 3 services in our microservice application that the api-gateway service will communicate with them. Authentication service, Job scheduler service and report service. These services are all gRPC services. During the development of api-gateway we do not have access to these services, therefore we mock their behavior. For this purpose we use gripmock. Instructions on how these service works can be found [here](https://github.com/tokopedia/gripmock). For gripmock configuration, we need .proto files for our gRPC services. Then we add stubs to the stub service. For this we can either define our stubs inside a json file (static stubbing) or add stubbing on the fly with a simple REST API (HTTP stub server is running on port :4771). Lets begin.

> - Create a folder for the project and choose a name for it (such as 'microservice'). Then create a folder named `api-gateway`. This folder is the root directory of the current project. You can then open the root folder in VS Code by right-clicking on the folder and selecting 'Open with Code'.
> - Inside the root directory create a folder with the name `gripmock`, then create the following files inside.
> - Create a Dockerfile and set content to `FROM tkpd/gripmock:v1.12.1`
> - Copy all .proto files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/plan/proto) to the gripmock directory.
> - Create a file named .env in the root directory and add the following content:

```
GRIPMOCK_ADMIN_PORT=4771
GRIPMOCK_GRPC_PORT=4770
API_GATEWAY_PORT=5010
```
> - Inside gripmock folder create a folder named stubs. Inside these folder we are going to create json stubs for our services. you can copy the files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/api-gateway/gripmock/stub). For example, for report.json we have the following content. Our grpc service name is `ReportService` and the method name is `ListReports`. For input, we define some `equals` critria, that is when we received a request for service `ReportService` and method `ListReports` and the values for filter, page and size are 1, 1, and 10 respectively, then return the following output. For more information see [here](https://github.com/tokopedia/gripmock).

```json
[
  {
    "service": "ReportService",
    "method": "ListReports",
    "input": {
      "equals": {
        "filter": 1,
        "page": 1,
        "size": 10
      }
    },
    "output": {
      "data": "..."
    }
  }
]
```
> - Inside root directory create a file named docker-compose.yml and add the following content.

```yaml
version: '3'
services:
  gripmock:
    build:
      context: ./gripmock
      dockerfile: Dockerfile
    container_name: gripmock
    ports:
      - ${GRIPMOCK_GRPC_PORT}:${GRIPMOCK_GRPC_PORT}
      - ${GRIPMOCK_ADMIN_PORT}:${GRIPMOCK_ADMIN_PORT}
    volumes:
      - ./gripmock:/mock
    # we use admin to manage grpc servers. we connect to grpc servers using our clinets(from code).
    command: > 
      --admin-listen=0.0.0.0 
      --admin-port=${GRIPMOCK_ADMIN_PORT} 
      --grpc-listen=0.0.0.0 
      --grpc-port=${GRIPMOCK_GRPC_PORT} 
      --stub=/mock/stub /mock/job.proto /mock/report.proto /mock/user.proto
```

> - Run `docker-compose up -d --build`. Now go to `http://localhost:4771/`. You can see the available mock grpc servers. These services acts as a real grpc servers. We connect them from our code and query them.

![Gripmock](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/920e3mavrh9acsm7g0fx.PNG)

> - Run `docker-compose down`

---

## Api-gateway service: Golang

We aim to develop an API gateway service that acts as a REST API server and a gRPC client. This service will receive REST requests from client applications and connect to various gRPC servers to fetch the required results. The API gateway acts as an intermediary that facilitates communication between the client and the gRPC servers. A summary of what we are going to do: First, we will import the .proto files for the gRPC servers and compile them into the Go programming language using [protoc](https://github.com/protocolbuffers/protobuf/releases/download/v23.3/protoc-23.3-linux-x86_64.zip). This step ensures that we have the necessary client libraries to interact with the gRPC services.
Next, we'll prepare gRPC clients for each individual gRPC service. These clients will allow our API gateway to establish connections and communicate with the respective gRPC servers. Finally, we will run our REST API service, which will expose multiple endpoints for client applications to interact with. Within each endpoint, we will utilize a specific gRPC client to handle the requests and retrieve data from the corresponding gRPC server. 

> - Create a folder named `api-gateway-service` inside `api-gateway` folder.
> - Create a Dockerfile inside `api-gateway-service` and set the contents to

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
  api-gateway:
    build: 
      context: ./api-gateway-service
      dockerfile: Dockerfile
    container_name: api-gateway
    environment:
      - SERVER_PORT=${API_GATEWAY_PORT}
      - CLIENT_DOMAIN=localhost # we call the server from this domain (for cookie registration ,...)
      - AUTH_SERVICE_URL=gripmock:${GRIPMOCK_GRPC_PORT}
      - SCHEDULER_SERVICE_URL=gripmock:${GRIPMOCK_GRPC_PORT}
      - REPORT_SERVICE_URL=gripmock:${GRIPMOCK_GRPC_PORT}
      - AUTH_PUBLIC_KEY_FILE=/run/secrets/auth-public-key
    ports:
      - ${API_GATEWAY_PORT}:${API_GATEWAY_PORT}
    command: sleep infinity
    volumes:
      - ./api-gateway-service:/usr/src/app
```
> - We are going to do all the development inside a docker container without installing Golang in our host machine. To do so, we run the containers and then attach VSCode to the api-gateway-service container. As you may noticed, the Dockerfile for api-gateway-service has no entry-point therefore we set the command value of it to `sleep infinity` to keep the container awake.
> - Now run `docker-compose up -d --build`
> - While running, attach to the api-gateway service by clicking bottom-left icon and then select `attach to running container`. Select api-gateway service and wait for a new instance of VSCode to start. Upon starting the attached instance of VSCode, you will be prompted to open a folder within the container. As per our Dockerfile configuration, we have designated the WORKDIR as /usr/src/app. Therefore, we will select this folder inside the container. It is important to note that this designated folder is mounted to the api-gateway-service folder on the host machine using Docker Compose volumes. Consequently, any changes made within the selected folder will be automatically synced to the corresponding folder on the host machine. This synchronization ensures that modifications made during development are reflected in both the container and the host environment.
> - After opening the folder `/usr/src/app`, open a new terminal and initialize the go project by running `go mod init github.com/<your_username>/play-microservices/api-gateway/api-gateway-service`. This command will create a go.mod file.
> - Run `go get -u google.golang.org/grpc`. This is a gRPC framework for running grpc server using Golang.
> - Now create a folder named proto and copy the proto files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/plan/proto).
> - Create file named build_grpc.sh and set the contents to:

```bash
#!/bin/bash
declare -a services=("proto")
for SERVICE in "${services[@]}"; do
   DESTDIR='proto'
   mkdir -p $DESTDIR
   protoc \
       --proto_path=$SERVICE/ \
       --go_out=$DESTDIR \
       --go-grpc_out=$DESTDIR \
       $SERVICE/*.proto
done
```
> - Run `source build_grpc.sh`. This command compile our .proto file to Golang. For each .proto file, Two files will be created. X.pb.go and X_grpc.pb.go. The first contains the proto models and the second contains the code for grpc service interface.
> - Note: We adopt a Golang project structure that aligns with the recommended guidelines stated [here](https://github.com/golang-standards/project-layout)
> - Create a folder named config and a file named `config.go`. Set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/api-gateway/api-gateway-service/config/config.go). Also create a file named .env in the same folder. we will put our internal environment variables here. Set contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/api-gateway/api-gateway-service/config/.env).
> - Create a folder named pkg in the root directory (beside mod.go). We will put general packages here. Inside Create a folder named logger then a file named logger.go and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/scheduler/scheduler-service/pkg/logger/logger.go).
> - Inside pkg folder create a folder named grpc and then a file named GRPC_Client.go. set the content to:

```go
package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPC_Client struct{}

func (jc *GRPC_Client) Connect(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
```
> - Create a folder named cookie inside pkg, then a file named cookie.go. Set the contents from [here](here.com).
> - Create this folder tree: `internal/models/report/grpc`. inside grpc folder create a file named client_service.go and set the contents to

```go
package grpc

import (
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/config"
	gr "github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/pkg/grpc"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
)

type ReportGRPCClient struct {
	log logger.Logger
	cfg *config.Config
	gr.GRPC_Client
}

func NewReportGRPCClient(log logger.Logger, cfg *config.Config) *ReportGRPCClient {
	return &ReportGRPCClient{log: log, cfg: cfg}
}

func (jc *ReportGRPCClient) GRPC_ListReports(c *gin.Context, listReportsRequest *proto.ListReportsRequest) (*proto.ListReportResponse, error) {
	jc.log.Info("ReportGRPCClient.GRPC_ListReports: Connecting to grpc server...")
	conn, err := jc.Connect(jc.cfg.ReportServiceURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := proto.NewReportServiceClient(conn)
	jc.log.Info("ReportGRPCClient.GRPC_ListReports: Conneced to grpc server...")
	jc.log.Infof("ReportGRPCClient.GRPC_ListReports: calling server for ListReports: %v", listReportsRequest)
	return client.ListReports(c, listReportsRequest)
}
```
> - This file contains the logic for creating a client for report gRPC service and then retrieve the list of reports. We will use this class in our rest handler.
> - Run `go install github.com/swaggo/swag/cmd/swag@latest` We will use this tool for generation of rest api documentations. 
> - Run `go get -u github.com/swaggo/gin-swagger` and `go get -u github.com/swaggo/files`
> - Create a folder named handler inside report folder and then a file named handler.go. set the contents to 

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/internal/models"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/internal/models/report/grpc"
	grpcutils "github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/pkg/grpc"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/proto"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	log logger.Logger
	*grpc.ReportGRPCClient
}

func NewReportHandler(log logger.Logger, cfg *config.Config) *ReportHandler {
	grpcClient := grpc.NewReportGRPCClient(log, cfg) // Initialize the embedded type
	return &ReportHandler{log: log, ReportGRPCClient: grpcClient}
}

// @Summary Get the list of reports
// @Description retrieve the reports
// @Tags report
// @Produce json
// @Param   page      query    int     true        "Page"
// @Param   size      query    int     true        "Size"
// @Success 200 {array} models.ListReportResponse
// @Router /report/list [get]
func (rh *ReportHandler) ListReports(c *gin.Context) {
	rh.log.Infof("Request arrived: list reports: %v", c.Request)

	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		rh.log.Error("ReportHandler.ListReports: invalid input. page is not provided in the query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		rh.log.Error("ReportHandler.ListReports: invalid input. size is not provided in the query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	res, err := rh.GRPC_ListReports(c.Request.Context(), &proto.ListReportsRequest{Page: page, Size: size})
	if err != nil {
		rh.log.Errorf("ReportHandler.ListReports: error listing reports: %v",err.Error())
		status := grpcutils.GetHttpStatusCodeFromGrpc(err)
		c.AbortWithStatusJSON(status, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ListReportResponseFromProto(res))
}
```
> - Run `swag init -g ./cmd/main.go -o ./docs`. This command will parse comments and generate required files(docs folder and docs/doc.go) at /docs folder.
> - Create a folder named api inside internal folder and then a file named router.go. set the contents to 

```go
package api

import (
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/docs"
	rh "github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/internal/models/report/handler"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/pkg/logger"
	"github.com/gin-gonic/gin"

	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	log logger.Logger
	cfg *config.Config
}

func NewRouter(log logger.Logger, cfg *config.Config) *Router {
	return &Router{log: log, cfg: cfg}
}
func (r *Router) Setup(router *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/api/v1/ping", r.Pong)

	reportHandler := rh.NewReportHandler(r.log, r.cfg)
	router.GET("/api/v1/report/list", reportHandler.ListReports)
}

// @BasePath /api/v1
// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {string} Pong
// @Router /ping [get]
func (s *Router) Pong(c *gin.Context) {
	c.JSON(http.StatusOK, "Pong")
}
```

> - Create a folder named server inside internal folder. Then a file named server.go. set the content to

```go
package server

import (
	"fmt"

	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/config"
	api "github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/internal/api"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	log logger.Logger
	cfg *config.Config
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{log: log, cfg: cfg}
}

func (s *Server) Run() {
	r := gin.Default()

	router := api.NewRouter(s.log, s.cfg)
	router.Setup(r)

	r.Run(fmt.Sprintf("0.0.0.0:%s", s.cfg.ServerPort))
}
```
> - Create a folder named cmd and a file named main.go. Set the content to 

```go
package main

import (
	"log"

	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/internal/server"
	"github.com/<yourusername>/play-microservices/api-gateway/api-gateway-service/pkg/logger"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Environment: %s",
		cfg.AppVersion,
		cfg.Logger_Level,
		cfg.Environment,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	appLogger.Info("Starting the server")
	s := server.NewServer(appLogger, cfg)
	s.Run()
}

```
> - Run `go mod tidy`
> - Run `go run cmd/main.go`
> - Open a browser and go to `http://localhost:5010/api/v1/report/list?page=1&size=10` If everything goes on plan, you can connect to the server and receive the results. What happens here is that from the browser we send a rest request to our api-gateway by calling `http://localhost:5010/api/v1/report/list?page=1&size=10`. The gateway receive our parameters, and then create a gRPC request to our mock server. The mock server returns the result and then the api-gateway return that result to us. 

![Rest report service](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/dm2nyi88q1fdty8lir8h.PNG)

> - Stop the server by hitting `ctl + c`
> - We repeat the exact procedure for user and job services. For user and job you can get the files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/api-gateway/api-gateway-service/internal/models/user) and [here](https://github.com/KhaledHosseini/play-microservices/tree/master/api-gateway/api-gateway-service/internal/models/report) respectively. Also do not forget to add the endpoints to router.go file from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/api-gateway/api-gateway-service/internal/api/router.go).
> - It is essential to consider the clear separation of REST API models and gRPC models. Just as we adhere to the practice of separating the database layer models from the gRPC layer models in our gRPC services, we follow the same approach here. By separating the models specific to the REST layer from those used in the gRPC layer, we gain valuable flexibility in querying and combining data from multiple gRPC services. This distinction allows us to define and manage models unique to each layer in a cohesive manner, facilitating efficient data retrieval and manipulation within the respective contexts.
> - Some notes on swagger and rest api parameter types. We use [gin-swagger](https://github.com/swaggo/gin-swagger) for creating documentation for our rest api. In rest api we have 5 parameter types:
  - Query parameters:  appended to the URL after a question mark (?). example: `http://localhost:5010/api/v1/report/list?page=1&size=10`
  - Path parameters: are part of the URL path. Example: `http://localhost:5010/api/v1/user/{id}`
  - Headers: are key-value pairs included in the request or response headers. Example: Authorization, Content-Type.
  -  The request body: carries additional data sent with the HTTP request. It is used to send complex or larger data payloads, such as JSON or XML, to the server.
  -  Form data: is used to submit data from an HTML form to the server. It consists of key-value pairs representing form fields and their values. Form data is encoded and sent in the body of an HTTP request with the "Content-Type" header set as "application/x-www-form-urlencoded" or "multipart/form-data".

> - For example the swag declarative comments for GetUser endpoint is as follows. For more details see [here](https://github.com/swaggo/swag/blob/master/README.md)

```go
// @Summary Get user
// @Description Get user
// @Tags user
// @Produce json
// @Success 200 {object} models.GetUserResponse
// @Router /user/get [get]
func (uh *UserHandler) GetUser(c *gin.Context) {

	res, err := uh.GRPC_GetUser(c, &proto.GetUserRequest{})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.GetUserResponseFromProto(res))
}
```
> - After adding swag comments Run `swag init --parseDependency -g ./cmd/main.go -o ./docs` and then run `go run cmd/main.go`. Go to `http://localhost:5010/swagger/index.html` and you can see the swagger documentation for your rest api. based on the stubs you have created for your gRPC servers, you can test the api endpoints.

![Swagger](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/1x1sqw6kjroindydcero.PNG)

> - In the end, we will add an auth interceptor to our api-gateway. Create a folder named interceptors inside `internal/api` and then a file named `auth_interceptor.go`. 

```go
package interceptors

import (
	"net/http"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/cookie"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// We do authentication and authorization in the end services. We just attach the auth headers to grpc requests.
func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := cookie.GetAccessToken(c) // Retrieve the access token from the request header
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": err.Error()})
			return
		}
		
		md := metadata.Pairs("authorization", accessToken)
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		// Update the request context with the modified context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

```
> - In microservices architecture, depending on the design one approach is to do authentication inside api-gateway and authorization inside downstream services. If we use another layer of authentication between microservices using protocols like mTLS, this approach can be considered safe. Another approach is to do both authentication and authorization of users inside downstream services. Here we use the second approach and inside interceptor file we just add the authorization header of http call to the metadata header of grpc context. In gRPC, headers are sent as part of the gRPC message metadata using key-value pairs. The headers are serialized using Protocol Buffers (protobuf) and encoded as binary data while REST API headers are sent as part of the HTTP request or response headers. They are typically represented as plain text in a key-value format, following the HTTP header specifications. 

> - Change the contents of router.go inside api folder to add the interceptor.

```go
func (r *Router) Setup(router *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/api/v1/ping", r.Pong)

	userHandler := uh.NewUserHandler(r.log, r.cfg)
	jobHandler := jh.NewJobHandler(r.log, r.cfg)
	reportHandler := rh.NewReportHandler(r.log, r.cfg)

	router.POST("/api/v1/user/create", userHandler.CreateUser)
	router.POST("/api/v1/user/login", userHandler.LoginUser)
	router.POST("/api/v1/user/refresh_token", userHandler.RefreshAccessToken)
	router.POST("/api/v1/user/logout", userHandler.LogOutUser)

	// Apply the middleware to the routes inside the router.Group function
	api := router.Group("/api")
	api.Use(interceptors.AuthenticateUser()) // Apply the middleware here
	{
		api.GET("/v1/user/get", userHandler.GetUser)
		api.GET("/v1/user/list", userHandler.ListUsers)

		api.POST("/v1/job/create", jobHandler.CreateJob)
		api.POST("/v1/job/update", jobHandler.UpdateJob)
		api.GET("/v1/job/get", jobHandler.GetJob)
		api.GET("/v1/job/list", jobHandler.ListJobs)
		api.POST("/v1/job/delete", jobHandler.DeleteJob)

		api.GET("/v1/report/list", reportHandler.ListReports)
	}
}
```

> - For subsequent requests, you need to add authorization headers to the request. In this case to be able to access protected endpoints, you need to login first. On logging in the cookie for the user will be set on the browser. You can inspect the cookies value by right click on the page and select `inspect`. Then under application tab and storage section you can see the values for `authorization`and `x-refresh-token` values. On subsequent requests, these cookies will be sent along with the request. Api-gateway then add the authorization cookie to the metadata of grpc context to be consumed by downstream services for authentication and authorization purposes.

---


## To DO

> - Add tests
> - Add tracing using Jaeger
> - Add monitoring and analysis using grafana
> - Refactoring

---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).