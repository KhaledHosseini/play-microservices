Welcome, everyone! This is the first article in a series called **"Play Microservices"** where I will be implementing a microservices architecture featuring a variety of services and technologies. Throughout this series, I aim to provide concise explanations without being overly verbose. So without further ado let's begin. 

---

## First:
 
As I am not going to get into the details of Microservices architecture, I suggest you to have a glance at [microservices.io](https://microservices.io/). I aim to introduce the steps involved in developing a typical microservices software; however, please note that the code and steps detailed throughout this series should not be considered complete nor of production-level quality.


---

## Contents:

- Goal application
- DevOps review
- Applying DevOps principles to our application: Development-> plan

> - Architecture selection
> - Microsystems architecture design of our app
> - Microsystem patterns
> - Choreography diagram
> - DevOps diagram
> - Tools and technologies
> - Next

---

**Our goal application:**
A simple task scheduling app with authentication. Users register, log in and schedule notification tasks.   

---
**DevOps review**
Before delving into the specifics, we must first take a broad overview of our product, tools and technologies, and development pipeline. In order to produce a final product, companies must progress through a multi-stage pipeline, each stage operating seamlessly and continuously. Due to the numerous stages involved, we will only briefly cover those that are relevant to full-stack developers. [DevOps](https://en.wikipedia.org/wiki/DevOps)  is a software development approach that emphasizes collaboration and communication between developers and IT operations to speed up software delivery, increase reliability, and improve customer satisfaction. It involves a wide range of tools and practices, including automation, continuous integration and delivery, and agile development. Usually the lifecycle of DevOps is symbolized in the form of an infinity loop.  
![DevOps life-cycle](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/uwlymslug7i10jfdj3ck.png)


In DevOps, we continuously perform the following tasks:


![DevOps phases](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/b7e5gm8otavif18aoknc.png)
To summarize the tasks performed in each phase, we can roughly state the following.

- _**Continuous development:**_

> - Plan
> - Code
> - On each commit integration phase is started.

- _**Continuous integration:**_

> - New code review
> - Integrate new code into existing code
> - Deploy to build server
> - Unit testing
> - Deploy to test environment
> - Feedback 

- _**Continuous testing:**_

> - Integration tests
> - Component tests
> - Contract tests
> - E2E tests
> - UI/UX tests
> - Security tests
> - Feedback

- _**Continuous deployment:**_

> - Deploy to stage environment
> - Configure and deploy to production environment

- _**Continuous feedback**_
- _**Continuous monitoring**_
- _**Continuous Operations**_

---

**Applying DevOps principles to our application: Development-> plan**
Now we want to apply above approach and go through this pipeline step by step.

- **Architecture selection**

The first step in our journey is to choose an architectural pattern for our application. Our goal is to utilize a [microsystems architecture](https://en.wikipedia.org/wiki/Microservices); however, another approach is the [Monolithic architecture](https://en.wikipedia.org/wiki/Monolithic_architecture). Each architecture has its own set of strengths and weaknesses, as well as unique features at its core.

- **Microsystems architecture design of our app**

How to design the microsystem architecture for our app? To design the architecture of our Microservice we follow the instructions in [assemblage](https://microservices.io/post/architecture/2023/02/09/assemblage-architecture-definition-process.html) and [A pattern language for microservices](https://microservices.io/patterns/).

> **_Discovering system operations:_** According to the description of our goal application, our system operations can be listed as:
> - _signUpUser()_
> - _signInUser()_
> - _refreshToken()_
> - _getUser()_
> - _scheduleNotificationTask()_
> - _listScheduledNotificationTasks()_
> - _sendNotification()_
> - _getReportOfSentNotifications()_

> **_Services and their collaborations:_** After conducting further investigations and refactoring, we have compiled a list of our services:
- Authentication service
- Scheduleservice
- Notification service
- API gateway service
- Report service
- Client service

> An initial diagram of our design can be depicted as:


![Initial diagram](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/jpui1ki75qkcq1f32dvm.png)


- **Microsystem patterns**

Our application is based on a microservices architecture, and to simplify the design process, we implemented microservices-specific patterns. The following provides an overview of the patterns that were utilized and the ones that should be taken into consideration. Please note that the goal is simply to demonstrate an implementation of the patterns. The fact that we are using these particular patterns does not necessarily mean they are the best fit for our specific use case.

- Service collaboration patterns:

> - [_Event sourcing_](https://microservices.io/patterns/data/event-sourcing.html): persist aggregates as a sequence of events (instead of saving objects in database)
> - [_CQRS_](https://microservices.io/patterns/data/cqrs.html): We define a view database (a new service) which subscribe to events from other services (on data change) and save a view of the desired aggregate that we want. This view is read only and we only query from that service. Reports service uses this pattern to create reports of notification tasks.
> - [_API composition_](https://microservices.io/patterns/data/api-composition.html): In this pattern, the result for a specific client request is generated from joining multiple queries to multiple services.
> - [_Saga_](https://microservices.io/patterns/data/saga.html): In this pattern we design a business transaction as a sequence of local transactions (which occurs in a single service locally, then publish an event). Reverse transactions can occur on failed events. Here task scheduler and notification service use this pattern to update the state of scheduled notifications.
> - [_Domain event_](https://microservices.io/patterns/data/domain-event.html)


- Other microservice patterns:

> - _Database per service._
> - _Team per service._
> - _Code source per service._

- Individual service patterns / principles:

> - [_TDD (Test driven development)_](https://en.wikipedia.org/wiki/Test-driven_development)
> - [_BDD (Behavior driven development)_](https://en.wikipedia.org/wiki/Behavior-driven_development)
> - [_SOLID principles:_](https://en.wikipedia.org/wiki/SOLID)  Following these patterns ensures adherence to SOLID principles.


![Microsystem patterns](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/tk8cs940gt2o7kyy7pvd.png)


- **Choreography diagram**

When we combine all of these components -- the services, their collaborations and connections, and the data flow -- we can represent the architecture of our app in the form of a diagram like this:


![Choreography](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/a4svuxcgx5lhriskmy6g.png)


- **DevOps diagram**

We can sketch the DevOps pipeline of our app development as follows:


![DevOps diagram](https://github.com/KhaledHosseini/play-microservices/blob/master/plan/DevOps.png)


- **Tools and technologies**

Since our goal is to demonstrate the flexibility of microservices architecture in terms of independence from specific technologies and systems, we will employ various programming languages for different services of the system. The following are the tools and technologies utilized in our application:

- **_Development tools:_**

> - [Docker](https://www.docker.com/): Containerization. 
> - [Visual studio code](https://code.visualstudio.com/): Code editor
> - [Dev containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers): Edit code inside containers

- **_Integration / deployment tools:_**

> - [Jenkins](https://www.jenkins.io/). Alternatives: [circleci](https://circleci.com/), [travis-ci](https://www.travis-ci.com/), [gitlab ci/cd](https://docs.gitlab.com/ee/ci/), [AWS CodePipeline](https://aws.amazon.com/codepipeline/), [Azure Pipelines](https://azure.microsoft.com/en-us/products/devops/pipelines)

- **_service collaboration / communication technologies:_**

> - [Kafka](https://kafka.apache.org/): Message broker. alternatives: [rabbitmq](https://www.rabbitmq.com/), [mosquitto](https://mosquitto.org/), [NSQ](https://nsq.io/)
> - [grpc](https://grpc.io/)
> - [Rest](https://restfulapi.net/)
> - [graphql](https://graphql.org/)

- **_monitoring technologies:_**

> - Metrices: [prometheus](https://prometheus.io/)
> - Tracing: [jaeger](https://www.jaegertracing.io/)
> - Logging: [Elasticsearch](https://www.elastic.co/)

- **_Auth service_**
[Rust](https://www.rust-lang.org/)
[postgres](https://www.postgresql.org/): database
[redis](https://redis.io/): cache

- **_Scheduler service_**
[Go](https://go.dev/)


- **Tools and technologies**

We've covered some of the phases of our DevOps pipeline, with a focus on the development: plan stage. Within the context of our microservices architecture, each service can have its own independent code base and unique technical implementation. However, there are logical connections between services, so collaboration between teams is required.

Once the codebase environments for each service have been prepared, we need to set up a CI/CD pipeline using tools such as Jenkins. Teams write their code and tests, push them to the code base, and trigger the CI/CD pipeline. The CI/CD tool will check out the code base (including any dependent repositories), start the build, test, and deployment pipeline, and automatically provide feedback to the team.

In the next part of this series, I'll implement the auth service.

Links to other parts:
[](url)
[](url)
[](url)

---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).