use crate::models::user::db::UserDBInterface;
use log::info;
use log::error;

//grpc models
use tonic::{Request, Response, Status};
use crate::proto::{
    user_profile_service_server::UserProfileService, GetUserResponse, GetUserRequest, ListUsersRequest, ListUsersResponse
};
use crate::models::user::db::Validate;

pub struct MyUserProfileService {
    pub db: Box<dyn UserDBInterface + Send + Sync>,
}

// Implement UserProfileService trait for Our MyUserProfileService
#[tonic::async_trait]
impl UserProfileService for MyUserProfileService {
    async fn get_user(&self, request: Request<GetUserRequest>) -> Result<Response<GetUserResponse>, Status> {
        info!("MyUserProfileService:get_user Got a request: {:#?}", &request);
        let metadata = request.metadata();
        if let Some(user_id) = metadata.get("user_id") {
            if let Ok(user_id) = user_id.to_str() {
                let user_id = user_id.parse::<i32>().unwrap();
                match self.db.get_user_by_id(&user_id).await {
                    Ok(user) => {
                        //convert from orm model to grpc model
                        let u: GetUserResponse = user.into();
                        return Ok(Response::new(u))
                    }
                    Err(err) => {
                        error!("Error finding user: {:#?}", &err);
                        return Err(Status::internal("User not found."))
                    }
                };
            }
        }
        return Err(Status::permission_denied("Unauthorized."))
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
        let list_users_request: ListUsersRequest = request.into_inner();
        match list_users_request.validate() {
            Ok(_)=>{
                info!("UserProfileService:list_users Input validation passed for ListUsersRequest")
            }
            Err(err)=> {
                error!("UserProfileService:list_users Input validation failed for ListUsersRequest. {:#?}",&err);
                return Err(Status::invalid_argument(format!("invalid argument: {:#?}",&err)));
            }
        }

        match self.db.get_users_list(&list_users_request.page,&list_users_request.size).await {
            Ok(users) => {
                //convert from orm model to grpc model
                let users: Vec<GetUserResponse> = users.into_iter().map(|user| user.into()).collect();
                let result = ListUsersResponse {
                    total_count:0,
                    total_pages: 1,
                    page:list_users_request.page,
                    size:list_users_request.size,
                    has_more: true,
                    users: users
                };
                return Ok(Response::new(result))
            }
            Err(err) => {
                error!("Error in database listing users: {:#?}", &err);
                return Err(Status::internal("errorlisting users."))
            }
        };
    }
}