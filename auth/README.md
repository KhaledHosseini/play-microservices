

Contents:
- [Other parts](#other-parts)
- [Summary](#summary)
- [Tools](#tools)
- [Docker dev environment](#docker-dev-environment)
- [Database service: postgres](#database-service-postgres)
- [pgAdmin service](#pgadmin-service)
- [Cache service: Redis](#cache-service-redis)
- [Redis commander service](#redis-commander-service)
- [Auth-grpcui-service](#auth-grpcui-service)
- [Auth service: Rust](#auth-service-rust)
- [To DO](#to-do)


---

## Other parts
This is the 2nd part of a series of articles under the name **"Play Microservices"**. Links to other parts:<br/>
Part 1: [Play Microservices: Bird's eye view](https://dev.to/khaledhosseini/play-microservices-birds-eye-view-3d44)<br/>
Part 2: You are here<br/>
Part 3: [Microservices: Scheduler](https://dev.to/khaledhosseini/play-microservices-scheduler-19km)<br/>
Part 4: [Play Microservices: Email service](https://dev.to/khaledhosseini/play-microservices-email-service-1kmc)<br/>
Part 5: [Play Microservices: Report service](https://dev.to/khaledhosseini/play-microservices-report-service-4jcm)<br/>
Part 6: [Play Microservices: Api-gateway service](https://dev.to/khaledhosseini/play-microservices-api-gateway-service-4a9j)<br/>
Part 7: [Play Microservices: Client service](https://dev.to/khaledhosseini/play-microservices-client-service-4jbf)<br/>
Part 8: [Play Microservices: Integration via docker-compose](https://dev.to/khaledhosseini/play-microservices-integration-via-docker-compose-2ddc)<br/>
Part 9: [Play Microservices: Security](https://dev.to/khaledhosseini/play-microservices-security-45e4)<br/>


---

## Summary

Our goal is to create an authentication service for our microservices application. To achieve authentication, we require three separate services: a database service, a cache service, and the authentication api service. In the development environment, we include three additional services for debugging purposes. These services are pgadmin to manage our database service, redis-commander service to manage our cache service and grpcui service to test our gRPC api.

![Summary](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/9xzxucgkknecdihslh2u.png)

At the end, the project directory structure will appear as follows:

![Directory structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/mu3vktfpri9lnq1f1y1o.PNG)

---

## Tools

The tools required In the host machine:

  - [Docker](https://www.docker.com/): Containerization tool
  - [VSCode](https://code.visualstudio.com/): Code editing tool
  - [Dev containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension for VSCode
  - [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) extension for VSCode
  - [Git](https://git-scm.com/)

The tools and technologies that we will use Inside containers for each service:

 - Database service: [Postgres](https://hub.docker.com/_/postgres)
 - Cache service: [Redis](https://hub.docker.com/_/redis)
 - pgadmin service: [Pgadmin](https://hub.docker.com/r/dpage/pgadmin4/)
 - Redis-commander service: [Rediscommander](https://hub.docker.com/r/rediscommander/redis-commander)
 - grpcui service: [grpcui](https://hub.docker.com/r/fullstorydev/grpcui) for testing our gRPC server.
 - Auth api service: 
  - [Rust](https://www.rust-lang.org/): Programming language
  - [Tonic](https://github.com/hyperium/tonic): gRPC framework for Rust
  - [Diesel](https://diesel.rs/): Query builder and ORM for our database communication.
  - [Redis ](https://docs.rs/redis): For our redis server communications.

---

## Docker dev environment

Development inside Docker containers can provide several benefits such as consistent environments, isolated dependencies, and improved collaboration. By using Docker, development workflows can be containerized and shared with team members, allowing for consistent deployments across different machines and platforms. Developers can easily switch between different versions of dependencies and libraries without worrying about conflicts.

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/id0cijn8hm77ohn0hk2j.png)

When developing inside a Docker container, you only need to install `Docker`, `Visual Studio Code`, and the `Dev Containers` and `Docker `extensions on VS Code. Then you can run a container using Docker and map a host folder to a folder inside the container, then attach VSCode to the running container and start coding, and all changes will be reflected in the host folder. If you remove the images and containers, you only need to create the container using the Dockerfile again and copy the contents of the host folder to the container folder to start again. Under the hood, When attaching VSCode to a running container, Visual Studio code install and run a special server inside the container which handle the sync of changes between the container and the host machine.

---

## Database service: postgres

> - Create a folder for the project and choose a name for it (such as 'microservice'). Then create a folder named `auth`. This folder is the root directory of authentication service. You can then open the root folder in VS Code by right-clicking on the folder and selecting 'Open with Code'.
> - Inside the root directory create a folder with the name auth-db-service, then create the following files inside.
> - Create a Dockerfile and set content to `FROM postgres:15.3`
> - Create a file named pass.txt and set content to a `password`
> - Create a file named user.txt and set content to `admin `
> - Create a file named db.txt and set content to `users_db`
> - Create a file named url.txt and set content to `auth-db-service:5432:users_db:admin:password123`
> - Create a file named .env in the root directory and set the content to `AUTH_POSTGRES_PORT=5432`.
> - Inside root directory create a file named docker-compose.yml and add the following content.

```yaml
version: '3'
services:
# database service for auth service
  auth-db-service:
    build:
      context: ./auth-db-service
      dockerfile: Dockerfile
    container_name: auth-db-service
    secrets:
      - auth-db-user
      - auth-db-pass
      - auth-db-db
      - auth-db-url
    environment:
      POSTGRES_USER_FILE: /run/secrets/auth-db-user
      POSTGRES_PASSWORD_FILE: /run/secrets/auth-db-pass
      POSTGRES_DB_FILE: /run/secrets/auth-db-db
    ports:
      - '${AUTH_POSTGRES_PORT}:${AUTH_POSTGRES_PORT}'
    volumes:
      -  auth-db-service-VL:/var/lib/postgresql/data 

volumes:
  auth-db-service-VL:

secrets:
  auth-db-user:
    file: auth-db-service/user.txt
  auth-db-pass:
    file: auth-db-service/pass.txt
  auth-db-db:
    file: auth-db-service/db.txt
  auth-db-url:
    file: auth-db-service/url.txt
```

> - In our Docker Compose file, we use secrets to securely share credential data between containers. While we could use an .env file and environment variables, this is not considered safe. When defining secrets in the Compose file, Docker creates a file inside each container under the `/run/secrets/` path, which the containers can then read and use. For example, we have set the path of the Docker Compose secret `auth-db-pass` to the `POSTGRES_PASSWORD_FILE `environment variable. We will be using these secrets in other services later in the project.

---

## pgAdmin service

The purpose of this service is solely for debugging and management of our running database server in the development environment.

> - Inside root directory create a folder with the name pgadmin-service
> - Create a Dockerfile and set content to `FROM dpage/pgadmin4:7.3`
> - Create a file named .env file and set the content to 

```
PGADMIN_DEFAULT_EMAIL=admin@admin.com
PGADMIN_DEFAULT_PASSWORD=password123
```
> - Create a file named servers.json and set the content to 

```json
{
    "Servers": {
      "1": {
        "Name": "auth-db-service",
        "Group": "Servers",
        "Host": "auth-db-service",
        "Port": 5432,
        "MaintenanceDB": "users_db",
        "Username": "admin",
        "PassFile": "/run/secrets/auth-db-url",
        "SSLMode": "prefer"
      }
    }
  }
```
> - Add the following lines to the .env file of the docker-compose (the .env file at the root directory of the project.)

```
PGADMIN_CONTAINER_PORT=80
PGADMIN__PORT=5050
```
> - Add the following to the service part of the docker-compose.yml.

```yaml
  pgadmin:
    build:
      context: ./pgadmin-service
      dockerfile: Dockerfile
    container_name: pgadmin-service
# load pgadmin defauld username and password to the environment.
    env_file:
      - ./pgadmin-service/.env
    ports:
      - "${PGADMIN__PORT}:${PGADMIN_CONTAINER_PORT}" 
    secrets:
      - auth-db-url
    volumes:
      - ./pgadmin-service/servers.json:/pgadmin4/servers.json #
    depends_on:
      - auth-db-service
```

> - Now open a terminal in your project directory and run docker-compose up. Docker Compose will download and cache the required images before starting your containers. For the first run, this may take a couple of minutes. If everything goes according to plan, you can then access the pgAdmin panel at http://localhost:5050/ and log in using the pgAdmin credentials from the .env file inside the pgadmin-service container.
Once inside pgAdmin, you should see that it has successfully connected to the auth-db-service container. You can also register any other running PostgreSQL services and connect to them directly. However, in this case, you will need to provide the credentials for the specific database service in order to connect.


![pgadmin](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/7d0g2kcxvsw8ppchkkq1.PNG)

> - Now run docker-compose down

---

## Cache service: Redis

> - Inside root directory create a folder with the name auth-cache-service
> - Create a Dockerfile and set content to `FROM redis:7.0.11-alpine`
> - Create a file named pass.txt and set content to a `password`
> - Create a file named users.acl and set content to `user default on >password123 ~* &* +@all`
> - Create a file named redis.conf and set content to `aclfile /run/secrets/auth-redis-acl`
> - Create a file named .env and set the content to `REDIS_URI_SCHEME=redis #rediss for tls`.
> - Add `AUTH_REDIS_PORT=6379` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  auth-cache-service:
    build:
      context: ./auth-cache-service
      dockerfile: Dockerfile
    container_name: auth-cache-service
    secrets:
      - auth-redis-pass
      - auth-redis-acl
    command: redis-server /usr/local/etc/redis/redis.conf
    ports:
      - '${AUTH_REDIS_PORT}:${AUTH_REDIS_PORT}'
    volumes:
      - auth-cache-serviceVL:/var/vcap/store/redis
      - ./auth-cache-service/redis.conf:/usr/local/etc/redis/redis.conf
```
> - Add the following to the secrets part of the docker-compose.yml.

```yaml
  auth-redis-pass:
    file: auth-cache-service/pass.txt
  auth-redis-acl:
    file: auth-cache-service/users.acl
```
> - Here, we first define users.acl for credentials. Then create a redis.conf file inside which we determine the aclfile path. We are going to pass the .acl file via docker compose secrets, therefore the path would be `/run/secrets/auth-redis-acl`. Then from docker compose file using command `redis-server /usr/local/etc/redis/redis.conf` we apply the config file (which has already mounted via a volume). Now run `docker-compose up` again. If everything goes according to plan, you can connect to redis terminal inside the container and check the user. To do so, click on bottom-left corner (Dev container) and select `Attach to running container`. 

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/kgi9gy44q2fmsp3al5r6.png)

> - To attach to the auth-cache-service, go to the list of running containers and select auth-cache-service. This will launch a new instance of Visual Studio Code attached to that service. Within the newly launched instance, open a new terminal and type redis-cli to connect to the Redis command line interface.
Once connected, running the ping command will result in an error message that says "NOAUTH Authentication required" because our server is password-protected. To address this, run the `AUTH password` command (note: password refers to the password that you defined inside the users.acl and pass.txt files within the auth-cache-service folder). Running ping again after authenticating should return a response of `pong`.

> - run `docker-compose down`

---

## Redis commander service

This service exists only in the development environment for debugging purposes. We use it to connect to auth-cache-service and manage the data.

> - Inside root directory create a folder with the name auth-redis-commander-service
> - Create a Dockerfile and set content to `FROM rediscommander/redis-commander:latest`
> - Add `Auth_REDIS_COMMANDER_PORT=8081` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```yaml
  auth-redis-commander-service:
    build:
      context: ./auth-redis-commander-service
      dockerfile: Dockerfile
    container_name: auth-redis-commander-service
    hostname: redis-commander-service
    # restart: always
    secrets:
      - auth-redis-pass
    ports:
      - ${Auth_REDIS_COMMANDER_PORT}:${Auth_REDIS_COMMANDER_PORT}
    environment:
      REDIS_HOST: auth-cache-service
      REDIS_PORT: ${AUTH_REDIS_PORT}
      REDIS_PASSWORD_FILE: /run/secrets/auth-redis-pass
    volumes:
      - redisCommanderVL:/data
    depends_on:
      - auth-cache-service
```

> - Here we passed the REDIS_PASSWORD_FILE via secrets. 
> - Now run `docker-compose up`. this will runs the containers. Now go to `http://localhost:8081/` and voila! the commander has attached to our redis server.


![redis commander](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/ju22ectn6xvvh2vb4ave.PNG)

> - As you can see, you can view and add data to the redis server.
> - run `docker-compose down`

---


## Auth-grpcui-service

Before we dive into coding our auth-server, let's start by initializing a gRPC UI service. This service provide an easy to use ui to test our grpc service in development environment. 
> - Create a folder named Auth-grpcui-service inside root folder.
> - Create a Docker file and set contents to `FROM fullstorydev/grpcui:v1.3.1`
> - Add the following to the services part of the docker-compose.yml file. 

```yaml
  auth-grpcui-service:
    build:
      context: ./auth-grpcui-service
      dockerfile: Dockerfile
    container_name: auth-grpcui-service
    command: -port $AUTH_GRPCUI_PORT -plaintext auth-service:${AUTH_SERVICE_PORT}
    ports:
      - ${AUTH_GRPCUI_PORT}:${AUTH_GRPCUI_PORT}

```
> - Our grpc service is ready. We will use it later.

## Auth service: Rust

Our goal is to develop a gRPC server with Rust. To achieve this, we begin by defining our protocol buffers (protobuf) schema in a .proto file. We then use a protocol buffer compiler tool to generate model structs in Rust based on this schema. These models are known as contract or business layer models, and they should not necessarily match the exact data structures stored in the database. Therefore, it's recommended to separate the models of the gRPC layer from the models of the database layer, and use a converter to transform between them.
This approach has several benefits. For instance, if we need to modify the database models, there is no need to change the gRPC layer models as they are decoupled. Furthermore, by using object-relational mapping (ORM) tools, we can more easily communicate with the database within our code.

![grpc_orm](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/f6zipw660xdmmp3wf7b9.png)

In summary, to build a gRPC app using Rust with a database, we first define our protocol buffers schema in a .proto file and use tonic to generate the gRPC side code. Once the gRPC side is ready, we use diesel to define database schema and run the migration files against our database. Afterward, we define our database models with Diesel.
The next step is tying these two components(gRPC and Diesel models) together. For that, we create a converter class that can translate our gRPC tonic models into Diesel models. This allows us to move data to and from our database.  We then use our Diesel ORM to communicate with the database from both the server and clients in our gRPC app. lets start.

> - create a folder named auth-service inside root directory.
> - create a Dockerfile and set the contents to the following code (Note: You can install rust at host machine first, then init a starter project and then dockerize your app. But here we do all the things inside a container without installing rust on the host machine):

```bash
FROM rust:1.70.0 AS base
ENV PROTOC_VERSION=23.3
ENV PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-x86_64.zip
RUN apt-get update && apt-get install -y unzip && apt-get install libpq-dev
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \ 
    && rm -f $PROTOC_ZIP
RUN apt install -y netcat
RUN cargo install diesel_cli --no-default-features --features postgres

```
> - We install [protobuf](https://protobuf.dev/) and libpq-dev on our container. The first is used by tonic-build to compile .proto files to rust and the second is required by diesel-cli for database operations. 
> - Add `AUTH_SERVICE_PORT=9001` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Create a file named `.env` inside auth-service and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/.env). We will use this variables later, but we need to have it in this step because we have declared it inside docker-compose file.
> - Create a folder named keys inside auth-service. Inside this folder we put RSA keys for jwt generation and signing. Copy the files from [here](https://github.com/KhaledHosseini/play-microservices/tree/master/auth/auth-service/keys). You can Create your own rsa256 key pairs from [here](https://cryptotools.net/rsagen).
> - Add the following to the service part of the docker-compose.yml.

```yaml
  authentication-service:
    build:
      context: ./auth-service
      dockerfile: Dockerfile
      # target: developement
    # restart: always
    # command: bash -c 'while ! nc -z auth-db-service $AUTH_POSTGRES_PORT; do echo ..; sleep 1; done; echo backend is up;cargo sqlx database create; cargo run --bin migrate; cargo run --bin server;'
    container_name: authentication-service
    command: sleep infinity
    environment:
      DATABSE_SCHEME: postgresql
      DATABSE_DOMAIN: auth-db-service
      DATABSE_PORT: ${AUTH_POSTGRES_PORT}
      DATABSE_USER_FILE: /run/secrets/auth-db-user
      DATABSE_PASSWORD_FILE: /run/secrets/auth-db-pass
      DATABSE_DB_FILE: /run/secrets/auth-db-db
      REDIS_SCHEME: redis # rediss
      REDIS_DOMAIN: auth-cache-service
      REDIS_PORT: ${AUTH_REDIS_PORT}
      REDIS_PASS_FILE: /run/secrets/auth-redis-pass
      SERVER_PORT: ${AUTH_SERVICE_PORT}
      REFRESH_TOKEN_PRIVATE_KEY_FILE: /run/secrets/auth-refresh-private-key
      REFRESH_TOKEN_PUBLIC_KEY_FILE: /run/secrets/auth-refresh-public-key
      ACCESS_TOKEN_PRIVATE_KEY_FILE: /run/secrets/auth-access-private-key
      ACCESS_TOKEN_PUBLIC_KEY_FILE: /run/secrets/auth-access-public-key
    secrets:
      - auth-db-user
      - auth-db-pass
      - auth-db-db
      - auth-redis-pass
      - auth-access-public-key
      - auth-access-private-key
      - auth-refresh-public-key
      - auth-refresh-private-key
    env_file:
      - ./auth-service/.env
      - ./auth-cache-service/.env
    # you can access the database service using browser from: AUTH_POSTGRES_PORT and from within docker: AUTH_POSTGRES_CONTAINER_PORT
    ports:
      - ${AUTH_SERVICE_PORT}:${AUTH_SERVICE_PORT}
    volumes:
      - ./auth-service:/usr/src/app 
    depends_on:
      - auth-db-service
      - auth-cache-service
```
> - Add the following to the secrets section of the docker compose file.

```yaml
  auth-access-public-key:
    file: auth-service/keys/access_token.public.pem
  auth-access-private-key:
    file: auth-service/keys/access_token.private.pem
  auth-refresh-public-key:
    file: auth-service/keys/refresh_token.public.pem
  auth-refresh-private-key:
    file: auth-service/keys/refresh_token.private.pem
```
> - Now run `docker-compose up -d --build` we start our containers in detach mode. we temporarily set the command for the auth-service as `sleep infinity` to keep the container alive. Now we will do all the development inside the container. 
> - Click on bottom-left icon and select `Attach to running container` and then select auth-service. A new instance of VSCode will open. On the left side, it asks you to open a folder. We need to open the `/usr/src/app` inside the container. A folder which is mapped to our `auth-service` folder via docker-compose volume. Now Our VSCode is running inside the `/usr/src/app` folder inside the container. Then whatever change we made on this folder will sync to our auth-service folder inside the host machine.

![Open folder](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/dna0wxdbxuebhdvu2r9w.PNG)

> - Open a new terminal inside new VSCode instance that is attached to auth-service container. run `cargo init` . This will initialize our rust app.
> - run `cargo add tonic --features "tls"`. Tonic is a gRPC server.
> - run `cargo add tonic-reflection`. We add reflection to our tonic server. Reflection provides information about publicly-accessible gRPC services on a server, and assists clients at runtime to construct RPC requests and responses without precompiled service information.
> - Create a folder named proto inside auth-service (VSCode has already attached to the app directory inside container, which is mapped to auth-service folder inside host), then create a file named user.proto and define the proto schema inside it. For the contents of the file see [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/proto/user.proto)
> - We will use tonic on top of Rust to deal with gRPC, so we need to add `tonic-build` dependency in the Cargo.toml. To do that run: `cargo add --build tonic-build` and `cargo add prost`. Tonic build uses prost to compile .proto files. Create an empty folder named proto inside src folder (We will set the tonic output folder to this folder so that tonic create compiled .proto files here) and then create a file named Build.rs inside auth-service folder add paste the following code:

```rust
use std::{env,path::PathBuf};
fn main()->Result<(),Box<dyn std::error::Error>>{
   env::set_var("OUT_DIR", "./src/proto");
   let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());
   tonic_build::configure()
   .build_server(true)
   .out_dir(out_dir.clone())
   .file_descriptor_set_path(out_dir.join("user_descriptor.bin"))
   .compile(&["./proto/user.proto"],&["proto"])//.compile(&["../proto/user.proto"],&["../proto"])//first argument is .proto file path, second argument is the folder.
   .unwrap_or_else(|e| panic!("protobuf compile error: {}", e));
   
   Ok(())
}
```
> - Create `models/user/grpc` folder tree inside src and then create a file inside named `user_grpc_service.rs`. Add the following code:

```rust
use tonic::{Request, Response, Status};
use crate::proto::{
    user_service_server::UserService, CreateUserResponse, CreateUserRequest,
    LoginUserRequest, LoginUserResponse,
    RefreshTokenRequest,RefreshTokenResponse, LogOutRequest, LogOutResponse
};

#[derive(Debug, Default)]
pub struct MyUserService{}

#[tonic::async_trait]
impl UserService for MyUserService{
    async fn create_user(
        &self,
        request: Request<CreateUserRequest>,
    ) -> Result<Response<CreateUserResponse>, Status> {
       println!("Got a request: {:#?}", &request);
       let reply = CreateUserResponse { message: format!("Successfull user registration"),};
       Ok(Response::new(reply))
    }

    async fn login_user(
        &self,
        request: Request<LoginUserRequest>,
    ) -> Result<Response<LoginUserResponse>, Status> {

        println!("Got a request: {:#?}", &request);

        let reply = LoginUserResponse {
            access_token : "".to_string(),
            access_token_age : 1,
            refresh_token : "".to_string(),
            refresh_token_age : 1,
        };

        Ok(Response::new(reply))
    }

    async fn refresh_access_token(
        &self,
        request: Request<RefreshTokenRequest>,
    ) -> Result<Response<RefreshTokenResponse>, Status> {
        println!("Got a request: {:#?}", &request);

        let reply = RefreshTokenResponse {
            access_token : "".to_string(),
            access_token_age : 1,
        };

        Ok(Response::new(reply))
    }

    async fn log_out_user(
        &self,
        request: Request<LogOutRequest>,
    ) -> Result<Response<LogOutResponse>, Status> {
        println!("Got a request: {:#?}", &request);

        let reply = LogOutResponse {
            message : "Logout successfull".to_string()
        };

        Ok(Response::new(reply))
    }
}
```
> - Create another file inside grpc folder named user_profile_grpc_service.rs and set the contents to:

```rust
use tonic::{Request, Response, Status};
use crate::proto::{
    user_profile_service_server::UserProfileService, GetUserResponse, GetUserRequest, ListUsersRequest, ListUsersResponse
};

#[derive(Debug, Default)]
pub struct MyUserProfileService{}

#[tonic::async_trait]
impl UserProfileService for MyUserProfileService{
    async fn get_user(&self, request: Request<GetUserRequest>) -> Result<Response<GetUserResponse>, Status>{
        println!("Got a request: {:#?}", &request);
        let reply = GetUserResponse {
            id: 1,
            name: "name".to_string(),
            email: "email".to_string(),
        };
        Ok(Response::new(reply))
    }
   async fn list_users(&self, request: Request<ListUsersRequest>) -> Result<Response<ListUsersResponse>, Status> {
     return Err(Status::internal("unimplemented."))
  }
}
```
> - Inside user folder (parent of grpc folder) create a file named grpc.rs and set the contents to:

```rust
pub mod user_grpc_service;
pub mod user_profile_grpc_service;

pub use crate::models::user::grpc::user_grpc_service::MyUserService;
pub use crate::models::user::grpc::user_profile_grpc_service::MyUserProfileService;
```
> - inside models folder (parent of src folder) create a file named user.rs and set the contents to:

```rust
pub mod grpc;
```
> - inside src folder (parent of models), create a file named models.rs and set the content to

```rust
pub mod user;
```
> - Some notes on how gRPC is initialized via .proto files: To define a schema for our gRPC service, we create a user.proto file. Inside this file, we define a service called UserService, which serves as an interface with functions that have their own input and output parameters. When we run a protocol buffer compiler (no matter the target programming language), it generates this service as an interface, which Rust implements as a trait.
To complete our gRPC service implementation, we need to define our own service class and implement the generated trait. This involves creating a structure that will serve as a gRPC service class and adding method implementations that correspond to the services defined in the UserService trait. 
> - We have created a `.env` file for auth-service. This file contains our environment variables like private and public keys ages. 
> - Run `cargo add dotenv`. 
> - Create a file named config.rs inside src folder. In this file, we are going to read environment variables along with the secrets that are passed via docker-compose. See file contents [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/config.rs). We will use this module later inside code to manage environment variables.
> - Run `cargo add tokio --features "macros sync rt-multi-thread"` [Tokio](https://tokio.rs/)  is a runtime for writing reliable asynchronous applications with Rust.
> - add the following code to the main.rs

```rust
mod models;
mod config;
mod proto {
    include!("./proto/proto.rs");
    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] = include_bytes!("./proto/user_descriptor.bin");
}

use models::user::grpc::MyUserService;
use models::user::grpc::MyUserProfileService;

use tonic::{transport::Server};
use proto::user_service_server::UserServiceServer;
use proto::user_profile_service_server::UserProfileServiceServer;


use config::Config;
use std::net::SocketAddr;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = Config::init();

    let mut server_addr = "[::0]:".to_string();
    server_addr.push_str(&config.server_port);
    let addr: SocketAddr = server_addr.as_str()
    .parse()
    .expect("Unable to parse socket address");

    let user_service = MyUserService::default();
    let user_profile_service = MyUserProfileService::default();
    let reflection_service = 
    tonic_reflection::server::Builder::configure()
.register_encoded_file_descriptor_set(proto::FILE_DESCRIPTOR_SET)
        .build()
        .unwrap();


    Server::builder()
        .add_service(UserServiceServer::new(user_service))
        .add_service(UserProfileServiceServer::new(user_profile_service ))
        .add_service(reflection_service)
        .serve(addr)
        .await?;
    Ok(())
}
```
> - Now our gRPC part is ready. Run `cargo run`. If everything goes according to plan, our service will start. Go to docker desktop and restart `auth-grpcui-service` service. Now go to `http://localhost:5000/ ` and voila. You can test the server. 

![grpc-ui](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/e6yeamoqeked4u3f2bbt.PNG)

> - Return back to VSCode and stop the server (ctl + c).
> - Next step is to prepare our database models and migrations. For this we need a running database server and the database_url to connect to it. For development environments it is common to pass this type of data via environment variables. Here we are following an approach that is more common in production environment and pass credentials via docker-compose secrets and then pass the location of the secrets via environment variables. If you run `printenv` you can see the list of environment variables. We can make the database_url using the environment variables provided to us. Actually, Inside config.rs, we generate database url using environment variables received via docker-compose.

```rust
        let database_scheme = get_env_var("DATABSE_SCHEME");
        let database_domain = get_env_var("DATABSE_DOMAIN");
        let database_port = get_env_var("DATABSE_PORT");
        let database_name_file = get_env_var("DATABSE_DB_FILE");
        let database_name = get_file_contents(&database_name_file);
        let database_user_file = get_env_var("DATABSE_USER_FILE");
        let database_user = get_file_contents(&database_user_file);
        let database_pass_file = get_env_var("DATABSE_PASSWORD_FILE");
        let database_pass = get_file_contents(&database_pass_file);
        //database url for postgres-> postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
        let database_url = format!("{}://{}:{}@{}:{}/{}",
        database_scheme,
        database_user,
        database_pass,
        database_domain,
        database_port,
        database_name);
```
> - In order to make database_url available in the terminal (to be used by diesel cli), we need to create it from the environment values we received via docker-compose. To do so, create a file named db_url_make.sh in the auth-service directory (beside cargo.toml). Then add the following code:

```shell
#!/bin/bash

db_name_file="${DATABSE_DB_FILE}"
db_name=$(cat "$db_name_file")
db_user_file="${DATABSE_USER_FILE}"
db_user=$(cat "$db_user_file")
db_pass_file="${DATABSE_PASSWORD_FILE}"
db_pass=$(cat "$db_pass_file")

db_scheme="${DATABSE_SCHEME}"
db_domain="${DATABSE_DOMAIN}"
db_port="${DATABSE_PORT}"

db_url="${db_scheme}://${db_user}:${db_pass}@${db_domain}:${db_port}/${db_name}"
echo "database url is ${db_url}"
export DATABASE_URL=${db_url}
```
> - Now run `source db_url_make.sh` then you can check if the DATABASE_URL available via `printenv` command.
> - If you encountered errors like `command not found` it is possibly due to the difference between end of line characters in linux and windows. You can convert the new lines in VSCode from bottom right corner.

![new line characters](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/tmz2q80q9w3dgs7dcmth.PNG) 

> - Run `cargo add diesel --features "postgres r2d2"` Diesel is ORM and query builder. `r2d2` feature gives us connection pool management capabilities.
> - Run `cargo add diesel-derive-enum  --features "postgres"` This package helps us use rust enums directly with diesel ORM. (At the moment diesel doesn't support native rust enums).
> - Run `diesel setup` this will create a folder named migrations and diesel.toml file beside cargo.toml.
> - Run `diesel migration generate create_users` This command creates two files named XXX_create_users/up.sql and XXX_create_users/down.sql. Edit up.sql to:

```sql
CREATE TYPE role_type AS ENUM (
  'admin',
  'customer'
);

CREATE TABLE users (
        id int NOT NULL PRIMARY KEY DEFAULT 1,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL UNIQUE,
        verified BOOLEAN NOT NULL DEFAULT FALSE,
        password VARCHAR(100) NOT NULL,
        role role_type NOT NULL DEFAULT 'customer'
    );
CREATE INDEX users_email_idx ON users (email);
```
> - Edit down.sql to

```sql
DROP TABLE IF EXISTS "users";
DROP TYPE role_type;
```

> - Run `diesel migration run` this will run the migrations against database_url and creates schema.rs inside src folder. In case of error `relation "users" already exists`, you can run  `diesel print-schema > src/schema.rs` to create the schema file for you. Also if you want to revert the migration just run: `diesel migration revert`, This will apply down.sql file to the database. you can run `diesel migration run` again.
> - We did the migrations using cli. For production environments, we can run the migrations on start inside the code. To accomplish this, Run `cargo add diesel_migrations` and then add the following codes to the main.rs

```rust
use diesel::{pg::PgConnection, prelude::*};
use diesel_migrations::EmbeddedMigrations;
use diesel_migrations::{embed_migrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("./migrations");

...

// add this to the start of main function.
 let config = Config::init();
    
        let database_url = config.database_url.clone();
    let mut connection = PgConnection::establish(&database_url)
    .unwrap_or_else(|_| panic!("Error connecting to {}", database_url));

    _ = connection.run_pending_migrations(MIGRATIONS);
```

> - Now run `cargo run` This will runs the migrations against the database on app start. You can go to pgadmin panel and see users table under servers: auth-db-service: databases: users_db: schema: public: tables.
> - Create a folder named db inside `src/models/user/` and then a file named model.rs inside db folder.

```rust
//Our ORM models.
use diesel::prelude::*;

#[derive(Queryable, Selectable, PartialEq,  Identifiable, Debug,)]
#[diesel(table_name = crate::schema::users)]
#[diesel(check_for_backend(diesel::pg::Pg))]
pub struct User {
    pub id: i32,
    pub name: String,
    pub email: String,
    pub verified: bool,
    pub password: String,
    pub role: RoleType,
}

#[derive(Insertable, Queryable, Debug, PartialEq)]
#[diesel(table_name = crate::schema::users)]
#[diesel(check_for_backend(diesel::pg::Pg))]
pub struct NewUser {
    pub name: String,
    pub email: String,
    pub password: String,
    pub role: RoleType,
}

#[derive(diesel_derive_enum::DbEnum,Debug,PartialEq)]
#[ExistingTypePath = "crate::schema::sql_types::RoleType"]
#[DbValueStyle = "snake_case"]//Important. it should be styled according to .sql: see https://docs.rs/diesel-derive-enum/2.1.0/diesel_derive_enum/derive.DbEnum.html
pub enum RoleType {
    ADMIN,
    CUSTOMER
}

```
> - inside db.rs, We define our database models. To clarify a bit more, We have grpc layer models that we have defined in our .proto file. Necessarily we do not save the exact same models in our database. One reason is the difference of data types and another is we may have our logic in the database layer. Inside this file, we will implement the logic to convert to/from proto layer models. Rust made this conversion easy to implement. We use `from `and `into `traits in rust. add the following code to the model.rs file.

```rust
use crate::proto::{ CreateUserRequest, CreateUserResponse, GetUserResponse,LogOutResponse, RoleType as RoleTypeEnumProto};
...
impl std::convert::From<RoleTypeEnumProto> for RoleType {
    fn from(proto_form: RoleTypeEnumProto) -> Self {
        match proto_form {
            RoleTypeEnumProto::RoletypeAdmin => RoleType::ADMIN,
            RoleTypeEnumProto::RoletypeCustomer => RoleType::CUSTOMER,
        }
    }
}

impl std::convert::From<i32> for RoleTypeEnumProto {
    fn from(proto_form: i32) -> Self {
        match proto_form {
            0 => RoleTypeEnumProto::RoletypeAdmin,
            1 => RoleTypeEnumProto::RoletypeCustomer,
            _ => RoleTypeEnumProto::RoletypeCustomer,
        }
    }
}

//Convert from User ORM model to UserReply proto model
impl From<User> for GetUserResponse{
    fn from(rust_form: User) -> Self {
        GetUserResponse{
            id: rust_form.id,
            name: rust_form.name,
            email: rust_form.email
        }
    }
}

//convert from CreateUserRequest proto model to  NewUser ORM model
impl From<CreateUserRequest> for NewUser {
    fn from(proto_form: CreateUserRequest) -> Self {
        //CreateUserRequest.role is i32
        let role: RoleTypeEnumProto = proto_form.role.into();
        NewUser {
            name: proto_form.name,
            email: proto_form.email,
            password: proto_form.password,
            role: role.into(),
        }
    }
}

impl From<&str> for LogOutResponse {
    fn from(str: &str)-> Self {
        LogOutResponse{
            message: str.to_string()
        }
    }
}
impl From<&str> for CreateUserResponse {
    fn from(str: &str)-> Self {
        CreateUserResponse{
            message: str.to_string()
        }
    }
}

```

> - Another logic we may consider inside the model.rs file is to define some interfaces for database and cache operations. This interfaces will be used by our grpc services without linking to a specific database or cache service. add the following code to the model.rs file.

```rust
use std::error::Error;
...
#[tonic::async_trait]
pub trait UserDBInterface {
    async fn is_user_exists_with_email(&self,email: &String)-> Result<bool, Box<dyn Error>>;
    async fn insert_new_user(&self,new_user: &NewUser)-> Result<usize, Box<dyn Error>>;
    async fn get_user_by_email(&self,email: &String)-> Result<std::option::Option<User>, Box<dyn Error>>;
    async fn get_user_by_id(&self,id: &i32)-> Result<User, Box<dyn Error>>;
    async fn get_users_list(&self,page: &i64, size: &i64)-> Result<Vec<User>, Box<dyn Error>>;
}

#[tonic::async_trait]
pub trait UserCacheInterface {
    async fn set_expiration(&self,key: &String,value: &String, seconds: usize)-> Result<(), Box<dyn Error>>;
    async fn get_value(&self,key: &String)-> Result<String, Box<dyn Error>>;
    async fn delete_value_for_key(&self,keys: Vec<String>)-> Result<String, Box<dyn Error>>;
}
```
> - add the following line to db.rs (inside user folder).

```rust
pub use crate::models::user::db::model::{User,NewUser,RoleType,UserDBInterface, UserCacheInterface};
```
> - Another logic to be implemented inside db folder is the code for validation of our input models. Create a file named validation.rs inside db folder and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/models/user/db/validation.rs). Add the following code to db.rs inside user folder.

```rust
pub mod validation;
pub use crate::models::user::db::validation::Validate;
```
> - Now it is time to Prepare our grpc services using `UserDBInterface `and `UserCacheInterface`. But before we need to add some packages for our cache service and manipulation of JWTs.
> - Run `cargo add redis --features "tokio-comp"`. 
> - Run `cargo add argon2` we use it for password hashing.
> - Run `cargo add serde --features derive`  [serde](https://serde.rs/) is a framework for serializing and deserializing Rust data structures.
> - Run `cargo add uuid --features "serde v4"`
> - Run `cargo add chrono --features serde`
> - Run `cargo add jsonwebtoken` [jsonwebtoken](https://docs.rs/jsonwebtoken/) Create and parses JWT.
> - Run `cargo add log env_logger`

> - Most of the above packages are being used for [JWT](https://jwt.io/introduction) generation and validations. JWT is a standard way for securely transmitting information between parties. The schematic pipeline of how typically JWT is used across services can be illustrated as below:


![JWT Pipeline](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/wq15mrh1ymbnb43srtq7.png)

> - Another module which is required here is one for creation and verifying json web tokens (JWT). Create a folder named utils inside src and then a file named jwtoken.rs inside utils folder and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/utils/jwtoken.rs).
> - Create a file named utils.rs inside src and set the contents to `pub mod jwtoken;`

> - Change the `MyUserService` definition to match the following code. 

```rust
use crate::models::user::db::{UserDBInterface, UserCacheInterface};
...
pub struct MyUserService {
    pub db: Box<dyn UserDBInterface + Send + Sync>,
    pub cache: Box<dyn UserCacheInterface + Send + Sync>,
    pub env: Config,
}

```
> - We will using above interface in our grpc methods without knowing which database engine or cache service we are using. for example, for create_user function we have:

```rust
        async fn create_user(
        &self,
        request: Request<CreateUserRequest>
    ) -> Result<Response<CreateUserResponse>, Status> {
        info!("MyUserService: create_user Got a request: {:#?}", &request);
        let create_user_request: CreateUserRequest = request.into_inner();
        match create_user_request.validate() {
            Ok(_)=>{
                info!("MyUserService:create_user Input validation passed for CreateUserRequest")
            }
            Err(err)=> {
                error!("MyUserService:create_user Input validation failed for CreateUserRequest. {:#?}",&err);
                return Err(Status::invalid_argument("invalid argument"));
            }
        }

       match self.db.is_user_exists_with_email(&create_user_request.email).await {
            Ok(exists) => {
                if exists {
                    //https://docs.rs/tonic/latest/tonic/struct.Status.html
                    return Err(Status::already_exists("User with that email already exists"));
                }
            },
            Err(err) => {
                error!("MyUserService:create_user Error checking email existence. check the database. {:#?}",&err);
                return Err(Status::internal("Problem checking user"));
            }
        };

        //we do not save the password. but a hash of it.
        let salt = SaltString::generate(&mut OsRng);
        let hashed_password = Argon2::default()
        .hash_password(create_user_request.password.as_bytes(), &salt)
        .expect("Error while hashing password")
        .to_string();
        
        let mut new_user: NewUser = create_user_request.into();
        new_user.password = hashed_password;
        
        info!("Inserting new user: {:#?}", &new_user);

        match self.db.insert_new_user(&new_user).await {
            Ok(_)=> {
                let reply: CreateUserResponse = "Successfull user registration".into();
                return Ok(Response::new(reply))
            },
            Err(err)=> {
                error!("Problem creating user with error: {:#?}",&err);
                return Err(Status::internal("Error registering user."));
            }
        };
    }

```
> - The complete code for user_grpc_service.rs and user_profile_grpc_service.rs can be found [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/models/user/grpc/user_grpc_service.rs) and [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/models/user/grpc/user_profile_grpc_service.rs) respectively.
> - Now it is time to prepare a database and a cache engine so that we can pass them to the constructor of MyUserService and MyUserProfileService structs.
> - inside models/user/db folder create a file named diesel_postgres.rs and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/models/user/db/diesel_postgres.rs). then another file named redis.rs and then set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/models/user/db/redis.rs).
> - Add the following to the db.rs file inside user folder.

```rust
pub mod diesel_postgres;
pub mod redis;
pub use crate::models::user::db::diesel_postgres::{PostgressDB,PgPool,PgPooledConnection};
pub use crate::models::user::db::redis::RedisCache;
```
> - Change the contents of main.rs to

```rust
extern crate log;
// in the main file, declare all modules of sibling files.
mod config;
mod models;
mod schema;
mod utils;
mod proto {
    include!("./proto/proto.rs");
    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] = include_bytes!("./proto/user_descriptor.bin");
}

//proto
use models::user::grpc::MyUserService;
use models::user::grpc::MyUserProfileService;
use tonic::transport::Server;
use proto::user_service_server::UserServiceServer;
use proto::user_profile_service_server::UserProfileServiceServer;

pub use crate::models::user::db::{PostgressDB,PgPool,PgPooledConnection};
pub use crate::models::user::db::RedisCache;
//diesel
use diesel::{pg::PgConnection, prelude::*};
use diesel::r2d2::{ Pool, ConnectionManager };
use diesel_migrations::EmbeddedMigrations;
use diesel_migrations::{embed_migrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("./migrations");

use config::Config;
use std::net::SocketAddr;


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    
    // env::set_var("RUST_BACKTRACE", "full");
    env_logger::init();

    let config = Config::init();
    let redis_url = config.redis_url.clone();
    let database_url = config.database_url.clone();
    let server_port = config.server_port.clone();
    let redis_client = match redis::Client::open(redis_url) {
        Ok(client) => {
            println!("✅Connection to the redis is successful!");
            client
        }
        Err(e) => {
            println!("🔥 Error connecting to Redis: {}", e);
            std::process::exit(1);
        }
    };

    let manager = ConnectionManager::<PgConnection>::new(&database_url);
    let db_pool = match Pool::builder().build(manager) {
        Ok(pool) => {
            println!("✅ Successfully created connection pool");
            pool
        }
        Err(e) => {
            println!("🔥 Error creating connection pool: {}", e);
            std::process::exit(1);
        }
    };

    //do migration
   
    let mut connection = PgConnection::establish(&database_url)
    .unwrap_or_else(|_| panic!("Error connecting to {}", database_url));
    _ = connection.run_pending_migrations(MIGRATIONS);

    //user service
    let my_user_service = MyUserService { 
        db: Box::new(PostgressDB { db_pool: db_pool.clone() }),
        cache: Box::new(RedisCache{redis_client: redis_client.clone()}),
        env: config.clone()
    };
    let user_service = UserServiceServer::new(my_user_service);

    let my_user_profile_service = MyUserProfileService {
        db: Box::new(PostgressDB { db_pool: db_pool })
    };
    let profile_service = UserProfileServiceServer::new(my_user_profile_service);
   
    //reflection service for user
    let reflection_service = tonic_reflection::server::Builder::configure()
        .register_encoded_file_descriptor_set(proto::FILE_DESCRIPTOR_SET)
        .build()
        .unwrap();

    let mut server_addr = "[::0]:".to_string();
    server_addr.push_str(&server_port);
    let addr: SocketAddr = server_addr.as_str()
    .parse()
    .expect("Unable to parse socket address");

    

    let svr = Server::builder()
        .add_service(profile_service)
        .add_service(user_service)
        .add_service(reflection_service)
        .serve(addr);
    println!("Server listening on {}", addr);
    svr.await?;
    Ok(())
}
```
> - Run `cargo run` If everything goes according to plan, our app will be compiled and starts. Go to docker desktop and restart `auth-grpcui-service` service. Assuming everything proceeds as expected, auth-grpcui-service will start at http://localhost:5000/

![grpcui](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/075vnrmpwxghh7j3114l.PNG)
> - As we have used reflector for our service, grpcui can automatically retrieve the endpoints and parameters automatically.
> - lets test our server. Select CreateUser from method name and fill in the form, then click Invoke. As long as nothing unexpected happens, we receive `"message": "Successfull user registration"` to check, go to pgadmin and inspect user's table.

![create user pgadmin](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/za5ddhl3jdh0kqem8qwy.PNG)

> - Now select loginUser and fill in the form with the data You used for CreateUser. You will receive access and refresh tokens. Now you can go to the redis-commander-service at http://localhost:8081/ and see the cache data.

![redis login](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/1j26w7d5reuwo7rq5dur.PNG)
> - Stop the server by hitting `ctl + c`
> We have defined two separate services in our .proto file. User service and profile service. We need to do authentication for profile service. (user services messages already contain credentials). To do so, add the following packages: 
> - Run `cargo add futures`
> - Run `cargo add tower`
> - Run `cargo add hyper`
> - Inside src, add a folder named interceptors and the a file named auth.rs. Set the content from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/src/interceptors/auth.rs)
> - Create a file named interceptors.rs inside src and set the content to `pub mod auth;`
> - Add he following code to the main function:

```rust
...
mod interceptors;

//interceptor
use crate::interceptors::auth::AuthenticatedService;
...
let profile_service_intercepted = AuthenticatedService {
        inner: profile_service,
        env: config,
};
...
let svr = Server::builder()
        .add_service(profile_service_intercepted)
        .add_service(user_service)
        .add_service(reflection_service)
        .serve(addr);
```
> - A summary on the auth interceptor: Requests for UserProfileService will go through this interceptor first. Inside this file, we will read the access token, then validate it and if validation passed, we attach the user id and the role of the token to the header of of the request. Inside UserProfileService  methods, we read those headers and do the authorization.
> - Run `cargo run`. Go to http://localhost:5000/ and do test the app using the following scenario.
  - Select proto.userService from service name.
  - Select CreateUser from method name: (name:admin,email:admin@test.com,password: password,role: admin)
  - Select LoginUser from method name. (email:admin@test.com,password: password)
  - Copy the value of access_token (Without double quotes)
  - Select UserProfileService from sevice name. 
  - Select ListUsers from method name. (page:1, size:10)
  - Add `authorization `as key and the token you have copied from the login user as value for metadata. click invoke.
> - Our service development environment is ready. you can continue coding and debugging.

> - One important note: for packages, we have used `cargo add package_name`, and this is not safe, as the packages may change in the future. In case of falling into a dependency hell, [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth/auth-service/Cargo.toml) is the package versions inside my cargo.toml file:

---


## To DO

> - Add tests
> - Add tracing using Jaeger
> - Add monitoring and analysis using grafana
> - Refactoring

---

I would love to hear your thoughts. Please comment your opinions. If you found this helpful, let's stay connected on Twitter! [khaled11_33](https://twitter.com/khaled11_33).