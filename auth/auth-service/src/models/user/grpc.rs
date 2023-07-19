pub mod user_grpc_service;
pub mod user_profile_grpc_service;

pub use crate::models::user::grpc::user_grpc_service::MyUserService;
pub use crate::models::user::grpc::user_profile_grpc_service::MyUserProfileService;