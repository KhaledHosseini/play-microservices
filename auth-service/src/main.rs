// in the main file, declare all modules of sibling files.
mod config;
mod services;
mod convert;
mod models;
mod schema;
mod token;

//proto
use services::user::{UserService};
use services::user_profile::{UserProfileService};
use tonic::{transport::Server};
use proto::user_server::{UserServer};
mod proto {
    tonic::include_proto!("proto");
    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] =
        tonic::include_file_descriptor_set!("user_descriptor");
}

//diesel
use diesel::{pg::PgConnection, prelude::*};
use diesel::r2d2::{ Pool, PooledConnection, ConnectionManager, PoolError };
use diesel_migrations::EmbeddedMigrations;
use diesel_migrations::{embed_migrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("./migrations");

use redis::{Client};
use config::Config;
use std::net::SocketAddr;



type PgPool = Pool<ConnectionManager<PgConnection>>;
type PgPooledConnection = PooledConnection<ConnectionManager<PgConnection>>;
pub struct Context {
    db_pool: PgPool,
    redis_client: Client,
    env: Config,
}

impl Context {
    // helper function to get a connection from the pool
    pub fn get_conn(&self) -> Result<PgPooledConnection, PoolError> {
       let _pool = self.db_pool.get().unwrap();
       Ok(_pool)
   }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    
    // env::set_var("RUST_BACKTRACE", "full");

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
    
    let my_user_service = UserService { context: Context {
            db_pool: db_pool,
            redis_client: redis_client,
            env: config,
        }
    };

    let user_service = UserServer::new(my_user_service);

    // let ups = UserProfileService::new()
    // let profile_service = AuthenticateMiddleware {
    //     inner: domain::gateway::service::users_server::UsersServer::new(UsersServerImpl {}),
    //     lock: cache.clone(),
    // };

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
        .add_service(user_service)
        .add_service(reflection_service)
        .serve(addr);
    println!("Server listening on {}", addr);
    svr.await?;
    Ok(())
}