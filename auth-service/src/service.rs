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