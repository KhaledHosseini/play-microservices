extern crate log;
// in the main file, declare all modules of sibling files.
mod config;
mod models;
mod schema;
mod utils;
mod interceptors;
// We have set the our directory for building the .proto inside build script to ./src/proto
// we define a module called proto and include the generated proto.rs file.
mod proto {
    include!("./proto/proto.rs");
    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] = include_bytes!("./proto/user_descriptor.bin");
}

//interceptor
use crate::interceptors::auth::AuthenticatedService;
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
    let profile_service_intercepted = AuthenticatedService {
        inner: profile_service,
        env: config,
    };
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
        .add_service(profile_service_intercepted)
        .add_service(user_service)
        .add_service(reflection_service)
        .serve(addr);
    println!("Server listening on {}", addr);
    svr.await?;
    Ok(())
}