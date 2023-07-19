use crate::models::user::db::UserDBInterface;
use log::info;
use log::error;

//grpc models
use tonic::{Request, Response, Status};
use crate::proto::{
    user_profile_service_server::UserProfileService, GetUserResponse, GetUserRequest, ListUsersRequest, ListUsersResponse
};

pub struct MyUserProfileService {
    pub db: Box<dyn UserDBInterface + Send + Sync>,
}

// Implement UserProfileService trait for Our MyUserProfileService
#[tonic::async_trait]
impl UserProfileService for MyUserProfileService {
    async fn get_user(&self, request: Request<GetUserRequest>) -> Result<Response<GetUserResponse>, Status> {
        info!("MyUserProfileService:get_user Got a request: {:#?}", &request);
        let metadata = request.metadata();
        let GetUserRequest { id } = request.get_ref();
        let mut authorized = false;
        if let Some(user_id) = metadata.get("user_id") {
            if let Ok(user_id) = user_id.to_str() {
                let user_id = user_id.parse::<i32>().unwrap();
                if user_id == *id{
                    authorized = true;
                }
            }
        }
        if let Some(user_role) = metadata.get("user_role") {
            if let Ok(user_role) = user_role.to_str() {
                if user_role.to_lowercase() == "admin".to_lowercase(){
                    authorized = true;
                }
            }
        }
        if !authorized {
            return Err(Status::permission_denied("Unauthorized."))
        }
        match self.db.get_user_by_id(&id).await {
            Ok(user) => {
                //convert from orm model to grpc model
                let u: GetUserResponse = user.into();
                return Ok(Response::new(u))
            }
            Err(err) => {
                error!("Error finding user: {}", err);
                return Err(Status::not_found("User not found."))
            }
        };
    }

    async fn list_users(&self, request: Request<ListUsersRequest>) -> Result<Response<ListUsersResponse>, Status> {
        info!("MyUserProfileService:list_users Got a request: {:#?}", &request);
        let metadata = request.metadata();
        let mut authorized = false;
        if let Some(user_role) = metadata.get("user_role") {
            if let Ok(user_role) = user_role.to_str() {
                if user_role.to_lowercase() == "admin".to_lowercase(){
                    authorized = true;
                }
            }
        }
        if !authorized {
            return Err(Status::permission_denied("Unauthorized."))
        }
        let ListUsersRequest { page,size } = request.get_ref();
        match self.db.get_users_list(page,size).await {
            Ok(users) => {
                //convert from orm model to grpc model
                let users: Vec<GetUserResponse> = users.into_iter().map(|user| user.into()).collect();
                let result = ListUsersResponse {
                    total_count:0,
                    total_pages: 1,
                    page:1,
                    size:1,
                    has_more: true,
                    users: users
                };
                return Ok(Response::new(result))
            }
            Err(err) => {
                error!("Error listin users: {}", err);
                return Err(Status::not_found("error listi users."))
            }
        };
    }
}