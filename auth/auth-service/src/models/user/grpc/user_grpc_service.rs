use crate::Config;
use log::info;
//grpc models
use tonic::{Request, Response, Status};
use crate::proto::{
    user_service_server::UserService, CreateUserResponse, CreateUserRequest,
    LoginUserRequest, LoginUserResponse,
    RefreshTokenRequest,RefreshTokenResponse, LogOutRequest, LogOutResponse
};

use crate::models::user::db::{UserDBInterface, UserCacheInterface};

//orm models
use crate::models::user::db::{User , NewUser};

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
       
    //    let CreateUserRequest { name, email, password, role } = &request.into_inner();
       let create_user_request: CreateUserRequest = request.into_inner();
       match self.db.is_user_exists_with_email(&create_user_request.email).await {
            Ok(exists) => {
                if exists {
                    //https://docs.rs/tonic/latest/tonic/struct.Status.html
                    return Err(Status::already_exists("User with that email already exists"));
                }
            },
            Err(err) => {
                return Err(Status::already_exists(format!("Problem checking user with error: {}",err)));
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
                let reply = CreateUserResponse { message: format!("Successfull user registration")};
                return Ok(Response::new(reply))
            },
            Err(err)=> {
                return Err(Status::internal(format!("Problem creating user with error: {}",err)));
            }
        };
    }

    async fn login_user(
        &self,
        request: Request<LoginUserRequest>,
    ) -> Result<Response<LoginUserResponse>, Status> {

        info!("MyUserService:login_user a request: {:#?}", &request);

        let LoginUserRequest { email, password } = request.into_inner();

        let user: User = match self.db.get_user_by_email(&email).await {
            Ok(user)=> {
                info!("Optional user: {:#?}", &user);
                if let Some(user) = user {
                    user
                }else {
                    return Err(Status::not_found("User not found."))
                }
            },
            Err(err)=> return Err(Status::not_found(format!("User not found. error: {}",err)))
        };
        

        let is_password_correct = PasswordHash::new(&user.password)
        .and_then(|parsed_hash| {
            Argon2::default().verify_password(password.as_bytes(), &parsed_hash)
        })
        .map_or(false, |_| true);

        if !is_password_correct {
            return Err(Status::internal(format!("Invalid password for: {}", email)))
        }

        let access_token_details = match jwtoken::generate_jwt_token(
            user.id,
            format!("{:?}", user.role),
            self.env.access_token_max_age,
            self.env.access_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };
    
        let refresh_token_details = match jwtoken::generate_jwt_token(
            user.id,
            format!("{:?}", user.role),
            self.env.refresh_token_max_age,
            self.env.refresh_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };

        
        match self.cache.set_expiration(&access_token_details.token_uuid.to_string(), &user.id.to_string(), (self.env.access_token_max_age * 60) as usize).await {
            Ok(_)=> {},
            Err(err)=> {
                return Err(Status::internal(format!("Cache saving failed with error: {}", err)));
            }
        }

        match self.cache.set_expiration(&refresh_token_details.token_uuid.to_string(), &user.id.to_string(), (self.env.refresh_token_max_age * 60) as usize).await {
            Ok(_)=> {},
            Err(err)=> {
                return Err(Status::internal(format!("Cache saving failed with error: {}", err)));
            }
        }

        let reply = LoginUserResponse {
            access_token : access_token_details.token.clone().unwrap(),
            access_token_age : self.env.access_token_max_age * 60,
            refresh_token : refresh_token_details.token.unwrap(),
            refresh_token_age : self.env.access_token_max_age * 60,
        };

        Ok(Response::new(reply))
    }

    async fn refresh_access_token(
        &self,
        request: Request<RefreshTokenRequest>,
    ) -> Result<Response<RefreshTokenResponse>, Status> {
        info!("MyUserService:refresh_access_token Got a request: {:#?}", &request);
        
        let RefreshTokenRequest { refresh_token } = request.into_inner();

        let refresh_token_details =
        match jwtoken::verify_jwt_token(self.env.refresh_token_public_key.to_owned(), &refresh_token)
        {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token validation failed with error: {}", e)));
            }
        };

        let user_id_str = match self.cache.get_value(&refresh_token_details.token_uuid.to_string()).await {
            Ok(user_id)=> user_id,
            Err(err)=> {
                return Err(Status::internal(format!("User id retrival failed: {}", err)));
            }
        };

        let user_id = user_id_str.parse::<i32>().unwrap();
        let user = match self.db.get_user_by_id(&user_id).await {
            Ok(user) => user,
            Err(err) => {
                info!("Error finding user: {}", err);
                return Err(Status::internal(format!("the user belonging to this token no logger exists.")));
            }
        };

        let access_token_details = match jwtoken::generate_jwt_token(
            user.id,
            format!("{:?}", user.role),
            self.env.access_token_max_age,
            self.env.access_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };

        match self.cache.set_expiration(&access_token_details.token_uuid.to_string(), &user.id.to_string(),(self.env.access_token_max_age * 60) as usize).await {
            Ok(_)=> {},
            Err(err)=> {
                return Err(Status::internal(format!("Catch saving failed with error: {}",err)));
            }
        }

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
        
        let LogOutRequest { refresh_token, access_token } = request.into_inner();

        let refresh_token_details =
        match jwtoken::verify_jwt_token(self.env.refresh_token_public_key.to_owned(), &refresh_token)
        {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token validation failed with error: {}", e)));
            }
        };

        match self.cache.delete_value_for_key(vec![
            refresh_token_details.token_uuid.to_string(),
            access_token.to_string(),
        ]).await {
            Ok(_)=> {},
            Err(err)=> {
                return Err(Status::internal(format!("Logout failed with error: {}",err)));
            }
        }

        let reply = LogOutResponse {
            message : "Logout successfull".to_string()
        };

        Ok(Response::new(reply))
    }
}