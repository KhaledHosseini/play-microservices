# play-microservices
A simple job scheduler app using microservices architecture. Users can sign up (as admin or normal role).
 - Admins can query the list of all users.
 - Admins can schedule email jobs to be run in the future.
 - Admins can query the reports.

Just run 'docker-compose up' from the root directory. Then go to http://localhost:3000/

![Image1](https://github.com/KhaledHosseini/play-microservices/blob/master/plan/image1.png)

![Image2](https://github.com/KhaledHosseini/play-microservices/blob/master/plan/image2.png)

## Articles:

 > - Plan: [Read](https://dev.to/khaledhosseini/play-microservices-birds-eye-view-3d44)
 > - Authentication service: Rust [Read](https://dev.to/khaledhosseini/play-microservices-authentication-4di3)
 > - Scheduler service: Go [Read](https://dev.to/khaledhosseini/play-microservices-scheduler-19km)
 > - Email service: Python [Read](https://dev.to/khaledhosseini/play-microservices-email-service-1kmc)
 > - Report service: Python [Read](https://dev.to/khaledhosseini/play-microservices-report-service-4jcm)
 > - API gateway service: Go [Read](https://dev.to/khaledhosseini/play-microservices-api-gateway-service-4a9j)
 > - Client service: Typescript [Read](https://dev.to/khaledhosseini/play-microservices-client-service-4jbf)

---

## Choreography

 ![Choreography](https://github.com/KhaledHosseini/play-microservices/blob/master/plan/choreography.svg)

---

## CI / CD Flow diagram

 ![Choreography](https://github.com/KhaledHosseini/play-microservices/blob/master/plan/developement_environment.svg)
 
---

## Technologies per service

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
  - [Typescript](https://www.typescriptlang.org/) : Programming language
  - [Next.js](https://nextjs.org/): Web framework that can be used for developing backend and front-end applications.
  - [Tailwind css](https://tailwindcss.com/):  A utility-first CSS framework 
  - [react-hook-form](https://react-hook-form.com/docs)
  - [TanStack Query](https://tanstack.com/query/latest): An asynchronous state management.

---

# How to add my service?
