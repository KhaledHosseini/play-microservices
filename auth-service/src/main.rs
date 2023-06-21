// in the main file, declare all modules of sibling files.
mod config;
mod service;
mod convert;
mod models;
mod schema;
mod token;

//proto
use service::{UserService};
use tonic::{transport::Server, Request, Response, Status};
use proto::user_server::{UserServer};
mod proto {
    tonic::include_proto!("proto");
    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] =
        tonic::include_file_descriptor_set!("user_descriptor");
}

//diesel
use diesel::{pg::PgConnection, prelude::*};
use diesel_migrations::EmbeddedMigrations;
use diesel_migrations::{embed_migrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("./migrations");

use config::Config;

use std::net::SocketAddr;


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    
    // env::set_var("RUST_BACKTRACE", "full");

    let config = Config::init();
    
    //do migration
    let database_url = config.database_url.clone();
    let mut connection = PgConnection::establish(&database_url)
    .unwrap_or_else(|_| panic!("Error connecting to {}", database_url));
    _ = connection.run_pending_migrations(MIGRATIONS);

    //user service
    let redis_url = config.redis_url.clone();
    let user_service = UserService::new(&database_url, &redis_url)?;
    let user_server = UserServer::with_interceptor(user_service, intercept);

    //reflection service for user
    let reflection_service = tonic_reflection::server::Builder::configure()
        .register_encoded_file_descriptor_set(proto::FILE_DESCRIPTOR_SET)
        .build()
        .unwrap();

    let mut server_addr = "[::0]:".to_string();
    server_addr.push_str(&config.server_port);
    let addr: SocketAddr = server_addr.as_str()
    .parse()
    .expect("Unable to parse socket address");

    println!("Server listening on {}", addr);

    Server::builder()
        .add_service(user_server)
        .add_service(reflection_service)
        .serve(addr)
        .await?;
    Ok(())
}

fn intercept(mut req: Request<()>) -> Result<Request<()>, Status> {
    println!("Intercepting request: {:?}", req);

    // Set an extension that can be retrieved by `say_hello`
    req.extensions_mut().insert(MyExtension {
        some_piece_of_data: "foo".to_string(),
    });

    Ok(req)
}

struct MyExtension {
    some_piece_of_data: String,
}