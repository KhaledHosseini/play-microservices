use crate::Config;
use log::{info, warn, error};
//grpc models
use tonic::{Request, Response, Status};
use crate::proto::{
    user_service_server::UserService, CreateUserResponse, CreateUserRequest,
    LoginUserRequest, LoginUserResponse,
    RefreshTokenRequest,RefreshTokenResponse, LogOutRequest, LogOutResponse
};

use crate::models::user::db::{UserDBInterface, UserCacheInterface, User , NewUser, Validate};

//others
use crate::utils::jwtoken;

use argon2::{
    password_hash::{rand_core::OsRng, PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
    Argon2,
};

// #[derive(Debug, Default)]
pub struct MyUserService {
    pub db: Box<dyn UserDBInterface + Send + Sync>,
    pub cache: Box<dyn UserCacheInterface + Send + Sync>,
    pub env: Config,
}


// Implement UserService trait for Our MyUserService
#[tonic::async_trait]
impl UserService for MyUserService {
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

    async fn login_user(
        &self,
        request: Request<LoginUserRequest>,
    ) -> Result<Response<LoginUserResponse>, Status> {

        info!("MyUserService:login_user a request: {:#?}", &request);

        let login_user_request: LoginUserRequest = request.into_inner();
        match login_user_request.validate() {
            Ok(_)=>{
                info!("MyUserService:login_user validation passed for LoginUserRequest");
            }
            Err(err)=> {
                error!("MyUserService:login_user error validating of  LoginUserRequest. {:#?}", &err);
                return Err(Status::invalid_argument("invalid argument"));
            }
       }
        let user: User = match self.db.get_user_by_email(&login_user_request.email).await {
            Ok(user)=> {
                info!("Optional user: {:#?}", &user);
                if let Some(user) = user {
                    user
                }else {
                    return Err(Status::not_found("User not found."))
                }
            },
            Err(err)=> {
                error!("MyUserService:login_user user not found error in db. {:#?}", &err);
                return Err(Status::not_found("User not found."))
            }
        };
        

        let is_password_correct = PasswordHash::new(&user.password)
        .and_then(|parsed_hash| {
            Argon2::default().verify_password(login_user_request.password.as_bytes(), &parsed_hash)
        })
        .map_or(false, |_| true);

        if !is_password_correct {
            warn!("MyUserService:login_user invalid password.");
            return Err(Status::permission_denied("Invalid password or username"))
        }

        let access_token_details = match jwtoken::generate_jwt_token(
            user.id,
            format!("{:?}", user.role),
            self.env.access_token_max_age,
            self.env.access_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                error!("MyUserService:login_user Token generation failed with error:{:#?}", &e);
                return Err(Status::internal("Token generation failed."))
            }
        };
    
        let refresh_token_details = match jwtoken::generate_jwt_token(
            0,
            "".into(),
            self.env.refresh_token_max_age,
            self.env.refresh_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                error!("MyUserService:login_user Token generation failed with error:{:#?}", &e);
                return Err(Status::internal("Token generation failed."))
            }
        };


        match self.cache.set_expiration(&refresh_token_details.token_uuid.to_string(), &user.id.to_string(), (self.env.refresh_token_max_age * 60) as usize).await {
            Ok(_)=> {},
            Err(err)=> {
                error!("MyUserService:login_user Cache saving failed with error:{:#?}", &err);
                return Err(Status::internal("Error with caching."))
            }
        }

        let reply = LoginUserResponse {
            access_token : access_token_details.token.unwrap(),
            access_token_age : self.env.access_token_max_age * 60,
            refresh_token : refresh_token_details.token.unwrap(),
            refresh_token_age : self.env.refresh_token_max_age * 60,
        };

        Ok(Response::new(reply))
    }

    async fn refresh_access_token(
        &self,
        request: Request<RefreshTokenRequest>,
    ) -> Result<Response<RefreshTokenResponse>, Status> {
        info!("MyUserService:refresh_access_token Got a request: {:#?}", &request);
        
        let refresh_token_request: RefreshTokenRequest = request.into_inner();
        match refresh_token_request.validate() {
            Ok(_)=>{
                info!("MyUserService:refresh_access_token validation passed for RefreshTokenRequest");
            }
            Err(err)=> {
                error!("MyUserService:refresh_access_token error validating of  RefreshTokenRequest. {:#?}", &err);
                return Err(Status::invalid_argument("invalid argument"));
            }
       }

        let refresh_token_details =
        match jwtoken::verify_jwt_token(self.env.refresh_token_public_key.to_owned(), &refresh_token_request.refresh_token)
        {
            Ok(token_details) => token_details,
            Err(e) => {
                error!("MyUserService:refresh_access_token error validating of  RefreshTokenRequest. {:#?}", &e);
                return Err(Status::unauthenticated("Token validation failed"));
            }
        };

        let user_id_str = match self.cache.get_value(&refresh_token_details.token_uuid.to_string()).await {
            Ok(user_id)=> match user_id {
                Some(value) => {
                    value
                },
                None => {
                    error!("MyUserService:refresh_access_token UserId for this token is not exists in cache: {:#?}", &refresh_token_details.token_uuid.to_string());
                    return Err(Status::unauthenticated("UserId for this token does not exist"));
                }
            },
            Err(err)=> {
                error!("MyUserService:refresh_access_token User id retrival failed: {:#?}", &err);
                return Err(Status::unauthenticated("User not found"));
            }
        };

        let user_id = user_id_str.parse::<i32>().unwrap();
        let user = match self.db.get_user_by_id(&user_id).await {
            Ok(user) => user,
            Err(err) => {
                error!("MyUserService:refresh_access_token User not found: {:#?}", &err);
                return Err(Status::unauthenticated("the user belonging to this token no logger exists."));
            }
        };

        let access_token_details = match jwtoken::generate_jwt_token(
            user.id,
            format!("{:?}", user.role),
            self.env.access_token_max_age,
            self.env.access_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(err) => {
                error!("MyUserService:refresh_access_token Token generation failed with error: {:#?}", &err);
                return Err(Status::unauthenticated("Token generation failed"))
            }
        };

        let reply = RefreshTokenResponse {
            access_token : access_token_details.token.clone().unwrap(),
            access_token_age : self.env.access_token_max_age * 60,
        };

        Ok(Response::new(reply))
    }

    async fn log_out_user(
        &self,
        request: Request<LogOutRequest>,
    ) -> Result<Response<LogOutResponse>, Status> {
        info!("MyUserService:log_out_user Got a request: {:#?}", &request);
        
        let log_out_request: LogOutRequest  = request.into_inner();
        match log_out_request.validate() {
            Ok(_)=>{
                info!("MyUserService:log_out_user validation passed for LogOutRequest");
            }
            Err(err)=> {
                error!("MyUserService:log_out_user error validating of  LogOutRequest. {:#?}", &err);
                return Err(Status::invalid_argument("invalid argument"));
            }
       }
       let refresh_token_details =
       match jwtoken::verify_jwt_token(self.env.refresh_token_public_key.to_owned(), &log_out_request.refresh_token)
       {
           Ok(token_details) => token_details,
           Err(err) => {
            error!("MyUserService:log_out_user error invalid token. {:#?}", &err);
               return Err(Status::unauthenticated("Token validation failed"));
           }
       };

        match self.cache.delete_value_for_key(&refresh_token_details.token_uuid.to_string()).await {
            Ok(_)=> {
                info!("MyUserService:log_out_user tokens deleted from cache");
            },
            Err(err)=> {
                error!("MyUserService:log_out_user Logout failed with error: {:#?}", &err);
                return Err(Status::internal("Logout failed with error in caching"));
            }
        }

        let reply: LogOutResponse = "Logout successfull".into();

        Ok(Response::new(reply))
    }
}