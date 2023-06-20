This is the 2nd part of a series of articles under the name **"Play Microservices"**. Links to other parts:
Part 1: Play Microservices: Bird's eye view
Part 2: You are here
Part 3: Coming soon

The source code for the project can be found [here](https://github.com/KhaledHosseini/play-microservices):

---

## Contents:

- **Summary**
- **Tools**
- **Docker dev environment**
- **Database service: postgres**
- **pgAdmin service**
- **Cache service: Redis**
- **Redis commander service**
- **Auth service: rust**
- **Auth-grpcui-service**
- **Conclusion**

---

- **Summary**

Our goal is to create an authentication service for our microservices application. To achieve authentication, we require three separate services: a database service, a cache service, and the authentication api service. In the development environment, we include three additional services for debugging purposes. These services are pgadmin to manage our database service. redis-commander service to manage our cache service and grpcui service to test our gRPC api.

![Summary](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/9xzxucgkknecdihslh2u.png)

At the end, the project directory structure will appear as follows:

![Directory structure](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/icnex5qasbpdq0pfytci.PNG)
---
- **Tools**

The tools required In the host machine:

> - [Docker](https://www.docker.com/): Containerization tool
> - [VSCode](https://code.visualstudio.com/): Code editing tool
> - [Dev containsers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension for vs code
> - [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) extension for VSCode
> - [Git](https://git-scm.com/)

The tools and technologies that we will use Inside containers for each service:

> Database service: [Postgres](https://hub.docker.com/_/postgres)
> cache service: [redis](https://hub.docker.com/_/redis)
> pgadmin service: [pgadmin](https://hub.docker.com/r/dpage/pgadmin4/)
> redis-commander service: [rediscommander](https://hub.docker.com/r/rediscommander/redis-commander)
> grpcui service: [grpcui](https://hub.docker.com/r/fullstorydev/grpcui)
> Auth api service: 
> - [rust](https://www.rust-lang.org/) : programming language
> - [Tonic](https://github.com/hyperium/tonic): gRPC framework for rust
> - [Diesel](https://diesel.rs/): Query builder and ORM for our database communication.
> - [redis ](https://docs.rs/redis) redis for our redis server communications.

---

- **Docker dev environment**

Development inside Docker containers can provide several benefits such as consistent environments, isolated dependencies, and improved collaboration. By using Docker, development workflows can be containerized and shared with team members, allowing for consistent deployments across different machines and platforms. Developers can easily switch between different versions of dependencies and libraries without worrying about conflicts.

![dev container](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/id0cijn8hm77ohn0hk2j.png)

When developing inside a Docker container, you only need to install Docker, Visual Studio Code, and the Dev Containers and Docker extensions on VS Code. Then you can run a container using Docker and map a host folder to a folder inside the container. You can start coding, and all changes will be reflected in the host folder, so if you remove the images and containers, you only need to create the container using the Dockerfile again and copy the contents of the host folder to the container folder to start again.

---

- **Database service: postgres**

> - To create the root folder for the project, choose a name for it (such as 'microservice') and create a new folder to hold all project files. This folder is the root directory of your project. You can then open the root folder in VS Code by right-clicking on the folder and selecting 'Open with Code'.
> - Inside the root directory create a folder with the name auth-db-service, then create the following files inside.
> - Create a Dockerfile and set content to `FROM postgres:15.3`
> - Create a file named pass.txt and set content to a `password`
> - Create a file named user.txt and set content to `admin `
> - Create a file named db.txt and set content to `users_db`
> - Create a file named url.txt and set content to `auth-db-service:5432:users_db:admin:password123`
> - Create a file named .env in the root directory and set the content to `AUTH_POSTGRES_PORT=5432`.
> - Inside root directory create a file named docker-compose.yml and add the following content.

```
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

> - In our Docker Compose file, we use secrets to securely share credential data between containers. While we could use an .env file and environment variables, this is not considered safe. When defining secrets in the Compose file, Docker creates a file inside each container under the `/run/secrets/` path, which the containers can then read and use. 

For example, we have set the path of the Docker Compose secret `auth-db-path` to the `POSTGRES_PASSWORD_FILE `environment variable. We will be using these secrets in other services later in the project.

---

- **pgAdmin service**

The purpose of this service is solely for debugging and management of our running database server in the development environment.

> - Inside root directory create a folder with the name pgadmin-service
> - Create a Dockerfile and set content to `FROM dpage/pgadmin4:7.3`
> - Create a file named .env file and set the content to 

```
PGADMIN_DEFAULT_EMAIL=admin@admin.com
PGADMIN_DEFAULT_PASSWORD=password123
```
> - Create a file named server.json and set the content to 

```
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

```
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
Once inside pgAdmin, you should see that it has successfully connected to the auth-db-service container. You can also register any other running PostgreSQL services and connect to them directly. However, in this case, you will need to provide the credentials for the specific database service in order to connect


![pgadmin](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/7d0g2kcxvsw8ppchkkq1.PNG)

> - Now run docker-compose down

---

- **Cache service: Redis**

> - Inside root directory create a folder with the name auth-cache-service
> - Create a Dockerfile and set content to `FROM redis:7.0.11-alpine`
> - Create a file named pass.txt and set content to a `password`
> - Create a file named users.acl and set content to `user default on >password123 ~* &* +@all`
> - Create a file named redis.conf and set content to `aclfile /run/secrets/auth-redis-acl`
> - Create a file named .env and set the content to `REDIS_URI_SCHEME=redis #rediss for tls`.
> - Add `AUTH_REDIS_PORT=6379` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```
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

```
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

- **Redis commander service**

This service exists only in the development environment for debugging purposes. We use it to connect to auth-cache-serve and manage the data.

> - Inside root directory create a folder with the name auth-redis-commander-service
> - Create a Dockerfile and set content to `FROM rediscommander/redis-commander:latest`
> - Add `Auth_REDIS_COMMANDER_PORT=8081` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Add the following to the service part of the docker-compose.yml.

```
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

- **Auth service: Rust**

Our goal is to develop a gRPC server with Rust. To achieve this, we begin by defining our protocol buffers (protobuf) schema in a .proto file. We then use a protocol buffer compiler tool to generate model structs in Rust based on this schema. These models are known as contract or business layer models, and they should not necessarily match the exact data structures stored in the database. Therefore, it's recommended to separate the models of the gRPC layer from the models of the database layer, and use a converter to transform between them.
This approach has several benefits. For instance, if we need to modify the database models, there is no need to change the gRPC layer models as they are decoupled. Furthermore, by using object-relational mapping (ORM) tools, we can more easily communicate with the database within our code.

![grpc_orm](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/f6zipw660xdmmp3wf7b9.png)

In summary, to build a gRPC app using Rust with a database, we first define our protocol buffers schema in a .proto file and use tonic to generate the gRPC side code. Once the gRPC side is ready, we use diesel to define database schema and run the migration files against our database. Afterward, we define our database models with Diesel.
The next step is tying these two components(gRPC and Diesel models) together. For that, we create a converter class that can translate our gRPC tonic models into Diesel models. This allows us to move data to and from our database.  We then use our Diesel ORM to communicate with the database from both the server and clients in our gRPC app. lets start.

> - create a folder named auth-service inside root directory.
> - create a Dockerfile and set the contents to the following code (Note: You can install rust at host machine first, then init a starter project and then dockerize your app. But here we do all the things inside a container without installing rust on the host machine):

```
FROM rust:1.70.0 AS base
ENV PROTOC_ZIP=protoc-3.13.0-linux-x86_64.zip
RUN apt-get update && apt-get install -y unzip && apt-get install libpq-dev
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \ 
    && rm -f $PROTOC_ZIP
RUN apt install -y netcat
RUN cargo install diesel_cli --no-default-features --features postgres

```
> - We install [protobuf](https://protobuf.dev/) and libpq-dev on our container. The first is used by tonic-build to compile .proto files to rust and the second is required by diesel-cli for database schema creation. 
> - Add `AUTH_SERVICE_PORT=9001` to the .env file of the docker-compose (the .env file at the root directory of the project.)
> - Create a file named `.env` inside auth-service and set the contents from [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth-service/.env). We will use this variables later, but we need to have it in this step because we have declared it inside docker-compose file.
> - Add the following to the service part of the docker-compose.yml.

```
  auth-service:
    build:
      context: ./auth-service
      dockerfile: Dockerfile
    container_name: auth-service
    command: sleep infinity
    environment:
      POSTGRES_HOST: auth-db-service
      POSTGRES_PORT: ${AUTH_POSTGRES_PORT}
      POSTGRES_USER_FILE: /run/secrets/auth-db-user
      POSTGRES_PASSWORD_FILE: /run/secrets/auth-db-pass
      POSTGRES_DB_FILE: /run/secrets/auth-db-db
      REDIS_HOST: auth-cache-service
      REDIS_PORT: ${AUTH_REDIS_PORT}
      REDIS_PASS_FILE: /run/secrets/auth-redis-pass
      SERVER_PORT: ${AUTH_SERVICE_PORT}
    secrets:
      - auth-db-user
      - auth-db-pass
      - auth-db-db
      - auth-redis-pass
    env_file:
      - ./auth-service/.env
      - ./auth-cache-service/.env
    ports:
      - ${AUTH_SERVICE_PORT}:${AUTH_SERVICE_PORT}
    volumes:
      - ./auth-service:/usr/src/app # mount a local path (auth-service) to container path. each time we change the code and then restart this service, the build starts again. see dockerfile for auth service.
    depends_on:
      - auth-db-service
      - auth-cache-service
```
> - Now run `docker-compose up -d --build` we start our containers in detach mode. we temporarily set the command for the auth-service as `sleep infinity` to keep the container alive. Now we will do all the development inside the container. 
> - Click on bottom-left icon and select `Attach to running container` and then select auth-service. A new instance of VSCode will open. On the left side, it asks you to open a folder. We need to open the `/usr/src/app` inside the container. A folder which is mapped to our `auth-service` folder via docker-compose volume. Now Our VSCode is running inside the `/usr/src/app` folder inside the container. Then whatever change we made on this folder will sync to our auth-service folder inside the host machine.
> - Open a new terminal inside new VSCode instance that is attached to auth-service container. run `cargo init` . This will initialize our rust app.
> - run `cargo add tonic --features "tls"`. tonic is a gRPC server.
> - run `cargo add tonic-reflection`. We add reflection to our tonic server. Reflection provides information about publicly-accessible gRPC services on a server, and assists clients at runtime to construct RPC requests and responses without precompiled service information.
> - Create a folder named proto inside auth-service (VSCode has already attached to the app directory inside container, which is mapped to auth-service folder inside host), then create a file named user.proto and define the proto schema inside it. For the contents of the file see [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth-service/proto/user.proto)
> - We will use tonic on top of Rust to deal with gRPC, so we need to add `tonic-build` dependency in the Cargo.toml. To do that run: `cargo add --build tonic-build` and `cargo add prost`. tonic build uses prost to compile .proto files.
Now create a file named Build.rs inside auth-service folder add paste the following code:

```
use std::{env, path::PathBuf};
fn main()->Result<(),Box<dyn std::error::Error>>{

   let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());
   tonic_build::configure()
   .build_server(true)
   //.out_dir(out_dir)
   .file_descriptor_set_path(out_dir.join("user_descriptor.bin"))
   .compile(&["./proto/user.proto"],&["proto"])//first argument is .proto file path, second argument is the folder.
   .unwrap_or_else(|e| panic!("protobuf compile error: {}", e));
   Ok(())
}
```
> - Create a file inside src folder named `service.rs`. Add the following code:

```
use tonic::{Request, Response, Status};
use crate::proto::{
    user_server::User, CreateUserReply, CreateUserRequest,
    UserReply, UserRequest, LoginUserRequest, LoginUserReply,
    RefreshTokenRequest,RefreshTokenReply, LogOutRequest, LogOutReply
};

#[derive(Debug, Default)]
pub struct UserService {}

#[tonic::async_trait]
impl User for UserService {
    async fn get_user(&self, request: Request<UserRequest>) -> Result<Response<UserReply>, Status> {
        println!("Got a request: {:#?}", &request);
        let reply = UserReply {
            id: 1,
            name: "name".to_string(),
            email: "email".to_string(),
        };
        Ok(Response::new(reply))
    }

    async fn create_user(
        &self,
        request: Request<CreateUserRequest>,
    ) -> Result<Response<CreateUserReply>, Status> {
       println!("Got a request: {:#?}", &request);
       let reply = CreateUserReply { message: format!("Successfull user registration"),};
       Ok(Response::new(reply))
    }

    async fn login_user(
        &self,
        request: Request<LoginUserRequest>,
    ) -> Result<Response<LoginUserReply>, Status> {

        println!("Got a request: {:#?}", &request);

        let reply = LoginUserReply {
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
    ) -> Result<Response<RefreshTokenReply>, Status> {
        println!("Got a request: {:#?}", &request);

        let reply = RefreshTokenReply {
            access_token : "".to_string(),
            access_token_age : 1,
        };

        Ok(Response::new(reply))
    }

    async fn log_out_user(
        &self,
        request: Request<LogOutRequest>,
    ) -> Result<Response<LogOutReply>, Status> {
        println!("Got a request: {:#?}", &request);

        let reply = LogOutReply {
            message : "Logout successfull".to_string()
        };

        Ok(Response::new(reply))
    }
}
```
> - Some notes on how gRPC is initialized via .proto files: To define a schema for our gRPC service, we create a user.proto file. Inside this file, we define a service called User, which serves as an interface with functions that have their own input and output parameters. When we run a protocol buffer compiler (no matter the target programming language), it generates this service as an interface, which Rust implements as a trait.
To complete our gRPC service implementation, we need to define our own service class and implement the generated trait. This involves creating a structure that will serve as a gRPC service class and adding method implementations that correspond to the services defined in the User trait. 
> - We have created a `.env` file for auth-service. This file contains our environment variables like private and public keys for token signature. You can generate your own public/private key pair [here](https://cryptotools.net/rsagen).
> - Create a file named `access.public` and set the value to the public key of access token (Only public key). We will pass this value to the services that may need to verify signature of access tokens signed by auth service.
> - Run `cargo add dotenv`. 
> - Create a file named config.rs inside src folder. In this file, we are going to read environment variables along with the secrets that are passed via docker-compose. see file contents [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth-service/src/config.rs). We will use this module later inside code to manage environment variables.
> - Run `cargo add tokio --features "macros sync rt-multi-thread"` [Tokio](https://tokio.rs/)  is a runtime for writing reliable asynchronous applications with Rust.
> - add the following code to the main.rs

```
mod service;
mod config;


use service::{UserService};
use tonic::{transport::Server};
use proto::user_server::{UserServer};
mod proto {
    tonic::include_proto!("proto");
    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] =
        tonic::include_file_descriptor_set!("user_descriptor");
}


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

    let user_service = UserService::default();
    let reflection_service = 
    tonic_reflection::server::Builder::configure()
.register_encoded_file_descriptor_set(proto::FILE_DESCRIPTOR_SET)
        .build()
        .unwrap();


    Server::builder()
        .add_service(UserServer::new(user_service))
        .add_service(reflection_service)
        .serve(addr)
        .await?;
    Ok(())
}
```
> - Now our gRPC part is ready. Run `cargo run`. If everything goes according to plan, our service will start. Then stop the server (ctl + c).
> - Next step is to prepare our database models and migrations. For this we need a running database server and the database_url to connect to it. For development environments it is common to pass this type of data via environment variables. Here we are following an approach that is more common in production environment and pass credentials via docker-compose secrets and then pass the location of the secrets via environment variables. If you run `printenv` you can see the list of environment variables. We can make the database_url using the environment variables provided to us. Actually, Inside config.rs, we generate database url using environment variables received via docker-compose.

```
let database_name_file = get_env_var("POSTGRES_DB_FILE");
        let database_name = get_file_contents(&database_name_file);
        let database_user_file = get_env_var("POSTGRES_USER_FILE");
        let database_user = get_file_contents(&database_user_file);
        let database_pass_file = get_env_var("POSTGRES_PASSWORD_FILE");
        let database_pass = get_file_contents(&database_pass_file);

        let database_url = format!("postgresql://{}:{}@{}:{}/{}",
        database_user,
        database_pass,
        database_host,
        database_port,
        database_name);
```
> - In order to make database_url available in the terminal (to be used by diesel cli, we need to create it from the environment values we received via docker-compose. To do so, create a file named db_url_make.sh in the auth-service directory (beside cargo.toml). Then add the following code:

```
db_name_file="${POSTGRES_DB_FILE}"
db_name=$(cat "$db_name_file")
db_user_file="${POSTGRES_USER_FILE}"
db_user=$(cat "$db_user_file")
db_pass_file="${POSTGRES_PASSWORD_FILE}"
db_pass=$(cat "$db_pass_file")

db_host="${POSTGRES_HOST}"
db_port="${POSTGRES_PORT}"
db_url="postgresql://${db_user}:${db_pass}@${db_host}:${db_port}/${db_name}"
echo "database url is ${db_url}"
export DATABASE_URL=${db_url}
```
> - Now run `source db_url_make.sh` then you can check if the DATABASE_URL available via `printenv` command.
> - Run `cargo add diesel --features "postgres r2d2"` Diesel is ORM and query builder. `r2d2` feature gives us connection pool management capabilities.
> - Run `diesel setup` this will create a folder named migrations and diesel.toml file beside cargo.toml.
> - Run `diesel migration generate create_users` This command creates two files named XXX_create_users/up.sql and XXX_create_users/down.sql. Edit up.sql to:

```
CREATE TABLE users (
        id int NOT NULL PRIMARY KEY DEFAULT 1,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL UNIQUE,
        verified BOOLEAN NOT NULL DEFAULT FALSE,
        password VARCHAR(100) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'user'
    );
CREATE INDEX users_email_idx ON users (email);
```
> - Edit down.sql to

```
DROP TABLE IF EXISTS "users";
```

> - Run `diesel migration run` this will run the migrations against database_url and creates schema.rs inside src folder. In case of error `relation "users" already exists`, you can run  `diesel print-schema > src/schema.rs` to create the schema file for you.
> - We did the migrations using cli. for production environments, we can run the migrations on start inside the code. To accomplish this, Run `cargo add diesel_migrations` and then add the following codes to the main.rs

```
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
> - Create a file named models.rs 

```
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
    pub role: String,
}

#[derive(Debug, Insertable, PartialEq)]
#[diesel(table_name = crate::schema::users)]
pub struct NewUser {
    pub name: String,
    pub email: String,
    pub password: String,
}
```

> - We have prepared proto side and ORM side of our models. Now it is time to create the converter. Create a file named convert.rs inside src directory.

```

//ORM models (database models)
use crate::models::{User, NewUser};

//Proto models (gRPC)
use crate::proto::{
    CreateUserRequest, UserReply
};

//Convert from User ORM model to UserReply proto model
impl From<User> for UserReply {
    fn from(rust_form: User) -> Self {
        UserReply{
            id: rust_form.id,
            name: rust_form.name,
            email: rust_form.email
        }
    }
}

//convert from CreateUserRequest proto model to  NewUser ORM model
impl From<CreateUserRequest> for NewUser {
    fn from(proto_form: CreateUserRequest) -> Self {
        NewUser {
            name: proto_form.name,
            email: proto_form.email,
            password: proto_form.password,
        }
    }
}
```
> - whenever you create a new module, do not forget to include it at the start of your entry point file of your binary (main.rs) to make them accessible to the app.
> - Now it is time to put all together. Change the UserService declaration to:

```
type PgPool = Pool<ConnectionManager<PgConnection>>;
type PgPooledConnection = PooledConnection<ConnectionManager<PgConnection>>;

// #[derive(Debug, Default)]
pub struct UserService {
    redis_client: Client,
    db_pool: PgPool,
}

impl UserService {
    pub fn new(db_url: &str, redis_url: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let redis_client = match redis::Client::open(redis_url) {
            Ok(client) => {
                println!("✅Connection to the redis is successful!");
                client
            }
            Err(e) => {
                println!("🔥 Error connecting to Redis: {}", e);
                return Err(Box::new(e))
            }
        };

        let manager = ConnectionManager::<PgConnection>::new(db_url);
        // let db_pool = Pool::builder().build(manager)?;
        let db_pool = match Pool::builder().build(manager) {
            Ok(pool) => {
                println!("✅ Successfully created connection pool");
                pool
            }
            Err(e) => {
                println!("🔥 Error creating connection pool: {}", e);
                return Err(Box::new(e))
            }
        };

        Ok(Self {
            redis_client,
            db_pool,
        })
    }

     // helper function to get a connection from the pool
     fn get_conn(&self) -> Result<PgPooledConnection, PoolError> {
        let _pool = self.db_pool.get().unwrap();
        Ok(_pool)
    }
}
```
> - For the rest, we will need the following packages:
> - Run `cargo add redis --features "tokio-comp"`. 
> - Run `cargo add argon2` we use it for password hashing.
> - Run `cargo add base64`
> - Run `cargo add serde --features derive`  [serde](https://serde.rs/) is a framework for serializing and deserializing Rust data structures.
> - Run `cargo add uuid --features "serde v4"`
> - Run `cargo add chrono --features serde`
> - Run `cargo add jsonwebtoken` [jsonwebtoken](https://docs.rs/jsonwebtoken/) Create and parses JWT.

> - Most of the above packages are being used for [JWT](https://jwt.io/introduction) generation and validations. JWT is a standard way for securely transmitting information between parties. The schematic pipeline of how typically JWT is used across services can be illustrated as below:


![JWT Pipeline](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/wq15mrh1ymbnb43srtq7.png)
> - Another module which is required here is one for manipulation of tokens. Create a file named token.rs in src folder and set the contents to [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth-service/src/token.rs).
> - The source code for the final service module can be found [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth-service/src/service.rs). We pass database and cache URLs to the service constructor. The service which is a Tonic gRPC service, receive Its function parameters as gRPC models. Then we convert the gRPC models to orm models (if necessary) and query our database. Then convert back the database models to gRPC models and return the results. 
> - the final main.rs file can be found [here](https://github.com/KhaledHosseini/play-microservices/blob/master/auth-service/src/main.rs)
> - Run `cargo run ` If everything goes on plan, our service will start at port 9001. We can access our Tonic grpc server. but how? it is an gRPC server. To test it we need tools like [Postman](https://www.postman.com/) or [Insomnia](https://insomnia.rest/). But we use a docker service for that. [grpcui](https://hub.docker.com/r/fullstorydev/grpcui)
> - before going to the next step, close the VSCode instance attached to auth-service and then run docker-compose stop from the terminal of local VSCode instance.

---

- **Auth-grpcui-service**

In order to connect to our auth-service we create a new service called Auth-grpcui-service. 
> - Create a folder named Auth-grpcui-service inside root folder.
> - Create a Docker file and set contents to `FROM fullstorydev/grpcui:v1.3.1`
> - Add the following to the services part of the docker-compose.yml file. 

```
  auth-grpcui-service:
    build:
      context: ./auth-grpcui-service
      dockerfile: Dockerfile
    container_name: auth-grpcui-service
    command: -port $AUTH_GRPCUI_PORT -plaintext auth-service:${AUTH_SERVICE_PORT}
    ports:
      - ${AUTH_GRPCUI_PORT}:${AUTH_GRPCUI_PORT}
    depends_on:
      - auth-service
```
> - run `docker-compose up -d`
> - Attach again to the auth service using VSCode. Click buttom-left icon -> attach to running container -> auth-service
> - Open a new terminal from the attached VSCode and then run `cargo run`
> - After our gRPC server started running, Go to docker desktop and restart  auth-grpcui-service (which is closed because at start it cannot find auth-service). Assuming everything proceeds as expected, auth-grpcui-service will start at http://localhost:5000/
> - open the http://localhost:5000/ and viola! 

![grpcui](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/075vnrmpwxghh7j3114l.PNG)
> - As we have used reflector for our service, grpcui can automatically retrieve the endpoints and parameters automatically.
> - lets test our server. Select CreateUser from method name and fill in the form, then click Invoke. As long as nothing unexpected happens, we receive `"message": "Successfull user registration"` to check, go to pgadmin and inspect user's table.

![create user pgadmin](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/za5ddhl3jdh0kqem8qwy.PNG)

> - Now select loginUser and fill in the form with the data You used for CreateUser. You will receive access and refresh tokens. Now you can go to the redis-commander-service at http://localhost:8081/ and see the cache data.

![redis login](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/1j26w7d5reuwo7rq5dur.PNG)

> - Our service development environment is ready. you can continue coding and debugging.

> - One important note: for packages, we have used `cargo add package_name`, and this is not safe, as the packages may change in the future. In case of falling into dependency hell, here is the package versions inside my cargo.toml file:

```
[package]
name = "auth_service"
version = "0.1.0"
edition = "2021"
authors = ["khaled hosseini <khaled.hosseini77@gmail.com>"]

# See more keys and their definitions at https://doc.rust-lang.org/cargo/reference/manifest.html

[dependencies]
argon2 = "0.5.0"
base64 = "0.21.2"
chrono = { version = "0.4.26", features = ["serde"] }
diesel = { version = "2.1.0", features = ["postgres", "r2d2"] }
diesel_migrations = "2.1.0"
dotenv = "0.15.0"
jsonwebtoken = "8.3.0"
prost = "0.11.9"
redis = { version = "0.23.0", features = ["tokio-comp"] }
serde = { version = "1.0.164", features = ["derive"] }
tokio = { version = "1.28.2", features = ["macros", "sync", "rt-multi-thread"] }
tonic = { version = "0.9.2", features = ["tls"] }
tonic-reflection = "0.9.2"
uuid = { version = "1.3.4", features = ["serde", "v4"] }

[build-dependencies]
tonic-build = "0.9.2"

# server binary
[[bin]]
name = "server"
path = "src/main.rs"

```


---

- **Conclusion**



 