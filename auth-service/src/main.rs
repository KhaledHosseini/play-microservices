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