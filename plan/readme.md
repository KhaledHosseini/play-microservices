Welcome, everyone! This is the first article in a series called **"Play Microservices"** where I will be implementing a microservices architecture featuring a variety of services and technologies. Throughout this series, I aim to provide concise explanations without being overly verbose. So without further ado let's begin. 

---

## Contents:

- [Contents:](#contents)
- [First:](#first)
  - [Our goal application:](#our-goal-application)
  - [DevOps review](#devops-review)
    - [Continuous development](#continuous-development)
    - [Continuous integration](#continuous-integration)
    - [Continuous testing](#continuous-testing)
    - [Continuous deployment](#continuous-deployment)
    - [Continuous feedback](#continuous-feedback)
  - [Applying DevOps principles to our application: Development-\> plan](#applying-devops-principles-to-our-application-development--plan)
    - [Architecture selection](#architecture-selection)
    - [Microservices architecture design of our app](#microservices-architecture-design-of-our-app)
      - [Discovering system operations](#discovering-system-operations)
      - [Services and their collaborations](#services-and-their-collaborations)
      - [C4 Model](#c4-model)
    - [Microservices patterns](#microservices-patterns)
    - [Choreography diagram](#choreography-diagram)
    - [DevOps diagram](#devops-diagram)
    - [Tools and technologies](#tools-and-technologies)
    - [Next](#next)

---

## First:
 
As I am not going to get into the details of Microservices architecture, I suggest you to have a glance at [microservices.io](https://microservices.io/). I aim to introduce the steps involved in developing a typical microservices software; however, please note that the code and steps detailed throughout this series should not be considered complete nor of production-level quality.

---

### Our goal application:
A simple task scheduling app with authentication. Users register (Admin or user) and log in to the app. Admins can get the list of all users and schedule email jobs to them. Also admins can query for reports.

---

### DevOps review
Before delving into the specifics, we must first take a broad overview of our product, tools and technologies, and development pipeline. In order to produce a final product, companies must progress through a multi-stage pipeline, each stage operating seamlessly and continuously. Due to the numerous stages involved, we will only briefly cover those that are relevant to full-stack developers. [DevOps](https://en.wikipedia.org/wiki/DevOps)  is a software development approach that emphasizes collaboration and communication between developers and IT operations to speed up software delivery, increase reliability, and improve customer satisfaction. It involves a wide range of tools and practices, including automation, continuous integration and delivery, and agile development. Usually the lifecycle of DevOps is symbolized in the form of an infinity loop.  
![DevOps life-cycle](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/uwlymslug7i10jfdj3ck.png)


In DevOps, we continuously perform the following tasks:


![DevOps phases](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/b7e5gm8otavif18aoknc.png)
To summarize the tasks performed in each phase, we can roughly state the following.

#### Continuous development

 - Plan
 - Code
 - On each commit integration phase is started.

#### Continuous integration

 - New code review
 - Integrate new code into existing code
 - Deploy to build server
 - Unit testing
 - Deploy to test environment
 - Feedback 

#### Continuous testing

 - Integration tests
 - Component tests
 - Contract tests
 - E2E tests
 - UI/UX tests
 - Security tests
 - Feedback

#### Continuous deployment

 - Deploy to stage environment
 - Configure and deploy to production environment

#### Continuous feedback
- _**Continuous monitoring**_
- _**Continuous Operations**_

---

### Applying DevOps principles to our application: Development-> plan
Now we want to apply above approach and go through this pipeline step by step.

#### Architecture selection

The first step in our journey is to choose an architectural pattern for our application. Our goal is to utilize a [microsystems architecture](https://en.wikipedia.org/wiki/Microservices); however, another approach is the [Monolithic architecture](https://en.wikipedia.org/wiki/Monolithic_architecture). Each architecture has its own set of strengths and weaknesses, as well as unique features at its core.

#### Microservices architecture design of our app

How to design the microservices architecture for our app? To design the architecture of our Microservice we follow the instructions in [assemblage](https://microservices.io/post/architecture/2023/02/09/assemblage-architecture-definition-process.html) and [A pattern language for microservices](https://microservices.io/patterns/).

 
##### Discovering system operations
According to the description of our goal application, our system operations can be listed as:
- _signUpUser()_
- _signInUser()_
- _refreshToken()_
- _getUser()_
- _getAllUsers()_
- _scheduleJob()_
- _getJob()_
- _deleteJob()_
- _updateJob()_
- _listJobs()_
- _ExecuteJob(JobTypeEmail)_
- _createReport()_
- _listReport()_

##### Services and their collaborations
After conducting further investigations and refactoring, we have compiled a list of our services:
- Authentication service
- Job-Scheduler service
- Email-job Executor service
- Report service
- API gateway service
- Client service

##### C4 Model
We use the [C4 model](https://c4model.com/) for visualizing our software architecture. According to this model we have 4 layers of abstraction. If you imagine your software as a map, the 1st layer has the top highest view and then you zoom in and in each layer you provide more details.

![context diagram](https://c4model.com/img/c4-overview.png)

In case of our app, we can visualize these layers as follows:
  - System context diagram 

![context diagram](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/233xhoanz6wui7apyx4t.png)
  - Container diagram 

![Containers](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/lzmsg66qtyc9093af7xt.png)

  - Component diagram : We have 6 microservices each have it's own components. 
  - Code diagram: We have 6 microservices each have it's own components. Each component has its own code diagram.

---

#### Microservices patterns

Our application is based on a microservices architecture, and to simplify the design process, we implemented microservices-specific patterns. The following provides an overview of the patterns that were utilized and the ones that should be taken into consideration. Please note that the goal is simply to demonstrate an implementation of the patterns. The fact that we are using these particular patterns does not necessarily mean they are the best fit for our specific use case.

- Service collaboration patterns:

  - [_Event sourcing_](https://microservices.io/patterns/data/event-sourcing.html): persist aggregates as a sequence of events (instead of saving objects in database)
  - [_CQRS_](https://microservices.io/patterns/data/cqrs.html): We define a view database (a new service) which subscribe to events from other services (on data change) and save a view of the desired aggregate that we want. This view is read only and we only query from that service. Reports service uses this pattern to create reports of notification tasks.
  - [_API composition_](https://microservices.io/patterns/data/api-composition.html): In this pattern, the result for a specific client request is generated from joining multiple queries to multiple services.
  - [_Saga_](https://microservices.io/patterns/data/saga.html): In this pattern we design a business transaction as a sequence of local transactions (which occurs in a single service locally, then publish an event). Reverse transactions can occur on failed events. Here task scheduler and notification service use this pattern to update the state of scheduled notifications.
  - [_Domain event_](https://microservices.io/patterns/data/domain-event.html)


- Other microservice patterns:

  - _Database per service._
  - _Team per service._
  - _Code source per service._

- Individual service patterns / principles:

  - [_TDD (Test driven development)_](https://en.wikipedia.org/wiki/Test-driven_development)
  - [_BDD (Behavior driven development)_](https://en.wikipedia.org/wiki/Behavior-driven_development)
  - [_SOLID principles:_](https://en.wikipedia.org/wiki/SOLID)  Following these patterns ensures adherence to SOLID principles.


![Microservice patterns](https://res.cloudinary.com/dql9iembg/image/upload/v1693371946/microservice_patterns.drawio_cjl1jh.svg)


#### Choreography diagram

When we combine all of these components -- the services, their collaborations and connections, and the data flow -- we can represent the architecture of our app in the form of a diagram like this:


![Choreography](https://res.cloudinary.com/dql9iembg/image/upload/v1693371945/choreography_n2imsz.svg)


#### DevOps diagram

We can sketch the DevOps pipeline of our app development as follows:


![DevOps dev environment](https://res.cloudinary.com/dql9iembg/image/upload/v1693371945/developement_environment_vt5cwc.svg)

#### Tools and technologies

Since our goal is to demonstrate the flexibility of microservices architecture in terms of independence from specific technologies and systems, we will employ various programming languages for different services of the system. The following are the tools and technologies utilized in our application:

- **_Development tools:_**

  - [Docker](https://www.docker.com/): Containerization. 
  - [Visual studio code](https://code.visualstudio.com/): Code editor
  - [Dev containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers): Edit code inside containers

- **_Integration / deployment tools:_**

  - [Jenkins](https://www.jenkins.io/). Alternatives: [circleci](https://circleci.com/), [travis-ci](https://www.travis-ci.com/), [gitlab ci/cd](https://docs.gitlab.com/ee/ci/), [AWS CodePipeline](https://aws.amazon.com/codepipeline/), [Azure Pipelines](https://azure.microsoft.com/en-us/products/devops/pipelines)

- **_service collaboration / communication technologies:_**

  - [Kafka](https://kafka.apache.org/): Message broker. alternatives: [rabbitmq](https://www.rabbitmq.com/), [mosquitto](https://mosquitto.org/), [NSQ](https://nsq.io/)
  - [grpc](https://grpc.io/)
  - [Rest](https://restfulapi.net/)
  - [graphql](https://graphql.org/)

- **_monitoring technologies:_**

  - Metrices: [prometheus](https://prometheus.io/)
  - Tracing: [jaeger](https://www.jaegertracing.io/)
  - Logging: [Elasticsearch](https://www.elastic.co/)

- **_Auth service_**

  - [Rust](https://www.rust-lang.org/): Programming language
  - [postgres](https://www.postgresql.org/): database
  - [redis](https://redis.io/): cache
  - [Tonic](https://github.com/hyperium/tonic): gRPC framework for Rust
  - [Diesel](https://diesel.rs/): Query builder and ORM for our database communication.
  - [Redis ](https://docs.rs/redis): For our redis server communications in rust.

- **_Scheduler service_**

  - [Golang](https://go.dev/) : programming language
  - [gRPC-GO](https://pkg.go.dev/google.golang.org/grpc): gRPC framework for golang
  - [Mongo driver](https://pkg.go.dev/go.mongodb.org/mongo-driver): Query builder for our database communication.
  - [Kafka go](https://pkg.go.dev/github.com/segmentio/kafka-go) for our message broker communications from go.
  - [Quartz ](github.com/reugn/go-quartz) for scheduling purposes.

- **_Email job executor service_**

  - [Python](https://www.python.org/) : programming language
  - [Kafka-python](https://pypi.org/project/kafka-python/): For our message broker communications from go.

- **_Reports service_**

  - [Python](https://python.org/) : programming language
  - [grpc](https://pypi.org/project/grpc/): gRPC framework for python
  - [Pymongo](https://pypi.org/project/pymongo/): Query builder for our database communication.
  - [Kafka-Python](https://pypi.org/project/kafka-python/) for our message broker communications from python.
  - [grpcio-tools](https://pypi.org/project/grpcio-tools/) For compiling .proto files to python.

- **_API gateway service_**

  - [Golang](https://go.dev/) : programming language
  - [gRPC-GO](https://pkg.go.dev/google.golang.org/grpc): gRPC framework for golang
  - [Gin](https://gin-gonic.com/):  Is a web framework written in Golang
  - [gin-swagger](https://github.com/swaggo/gin-swagger) for our rest api documentation.


- **_Client service_**

  - [Typescript](https://www.typescriptlang.org/) : programming language
  - [Next.js](https://nextjs.org/): Web framework that can be used for developing backend and front-end applications.
  - [Tailwind css](https://tailwindcss.com/):  A utility-first CSS framework 
  - [react-hook-form](https://react-hook-form.com/docs)
  - [TanStack Query](https://tanstack.com/query/latest): A asynchronous state management.


We've covered some of the phases of our DevOps pipeline, with a focus on the development: plan stage. Within the context of our microservices architecture, each service can have its own independent code base and unique technical implementation. However, there are logical connections between services, so collaboration between teams is required.

Once the codebase environments for each service have been prepared, we need to set up a CI/CD pipeline using tools such as Jenkins. Teams write their code and tests, push them to the code base, and trigger the CI/CD pipeline. The CI/CD tool will check out the code base (including any dependent repositories), start the build, test, and deployment pipeline, and automatically provide feedback to the team.

#### Next

In the next part of this series, I'll implement the authentication service.

Links to other parts:
Part 1: You are here<br/>
Part 2: [Play Microservices: Authentication service](https://dev.to/khaledhosseini/play-microservices-authentication-4di3)<br/>
Part 3: [Microservices: Scheduler](https://dev.to/khaledhosseini/play-microservices-scheduler-19km)<br/>
Part 4: [Play Microservices: Email service](https://dev.to/khaledhosseini/play-microservices-email-service-1kmc)<br/>
Part 5: [Play Microservices: Report service](https://dev.to/khaledhosseini/play-microservices-report-service-4jcm)<br/>
Part 6: [Play Microservices: Api-gateway service](https://dev.to/khaledhosseini/play-microservices-api-gateway-service-4a9j)<br/>
Part 7: [Play Microservices: Client service](https://dev.to/khaledhosseini/play-microservices-client-service-4jbf)<br/>
Part 8: [Play Microservices: Integration via docker-compose](https://dev.to/khaledhosseini/play-microservices-integration-via-docker-compose-2ddc)<br/>
Part 9: [Play Microservices: Security](https://dev.to/khaledhosseini/play-microservices-security-45e4)<br/>


---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).