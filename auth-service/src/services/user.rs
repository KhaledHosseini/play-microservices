use crate::Context;
//grpc models
use tonic::{Request, Response, Status};
use crate::proto::{
    user_server::User, CreateUserReply, CreateUserRequest,
    LoginUserRequest, LoginUserReply,
    RefreshTokenRequest,RefreshTokenReply, LogOutRequest, LogOutReply
};

//orm models
use crate::models:: {
    User as orm_user, NewUser
};
// users table
use crate::schema::{users as users_table};

use diesel::{ prelude::*};
use diesel::dsl::select;
use diesel::dsl::exists;
use diesel::sql_query;

//cache
use redis::{ AsyncCommands};

//others
use crate::token;
use argon2::{
    password_hash::{rand_core::OsRng, PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
    Argon2,
};

// #[derive(Debug, Default)]
pub struct UserService {
   pub context: Context
}


// Implement User trait for Our UserService
#[tonic::async_trait]
impl User for UserService {
    async fn create_user(
        &self,
        request: Request<CreateUserRequest>
    ) -> Result<Response<CreateUserReply>, Status> {
       println!("Got a request: {:#?}", &request);
       
    //    let CreateUserRequest { name, email, password, role } = &request.into_inner();
       let create_user_request: CreateUserRequest = request.into_inner();
       let conn = &mut self.context.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;

       let email_exists = select(exists(users_table::table.filter(users_table::dsl::email.eq(&create_user_request.email))))
       .get_result::<bool>(conn);
       
       match email_exists {
            Ok(exists)=> {
                if exists {
                    //https://docs.rs/tonic/latest/tonic/struct.Status.html
                    return Err(Status::already_exists("User with that email already exists"));
                }
            }
            Err(e) => {
                return Err(Status::already_exists(format!("Problem checking user with error: {}",e)));
            }
        }

        //we do not save the password. but a hash of it.
        let salt = SaltString::generate(&mut OsRng);
        let hashed_password = Argon2::default()
        .hash_password(create_user_request.password.as_bytes(), &salt)
        .expect("Error while hashing password")
        .to_string();
        
        let mut new_user: NewUser = create_user_request.into();
        new_user.password = hashed_password;
        
        println!("Inserting new user: {:#?}", &new_user);

        // I could not come up with a solution for storing enums via diesel native queries.
        // let insert_response = diesel::insert_into(users_table::table)
        // .values(&new_user)
        // .execute(conn);

        let role_enum_string = format!("{:?}", new_user.role);
        let query = format!("INSERT INTO users (name, email, password, role) VALUES ('{}', '{}', '{}', '{}')",
        &new_user.name,&new_user.email,&new_user.password,role_enum_string);
        let insert_response = sql_query(query)
        .execute(conn);

        match insert_response {
            Ok(_)=> {
                let reply = CreateUserReply { message: format!("Successfull user registration")};
                Ok(Response::new(reply))
            }
            Err (e) => {
                return Err(Status::internal(format!("Problem creating user with error: {}",e)));
            }
        }
       
    }

    async fn login_user(
        &self,
        request: Request<LoginUserRequest>,
    ) -> Result<Response<LoginUserReply>, Status> {

        println!("Got a request: {:#?}", &request);

        let LoginUserRequest { email, password } = request.into_inner();

        let conn = &mut self.context.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;
        let user_result = users_table::table.filter(users_table::dsl::email.eq(&email)).first::<orm_user>(conn).optional()
        .map_err(|err| Status::internal(err.to_string()))?;
        
        let user = match user_result {
            Some(user) => user,
            None => {
                return Err(Status::not_found("User not found."))
            }
        };

        let is_password_correct = PasswordHash::new(&user.password)
        .and_then(|parsed_hash| {
            Argon2::default().verify_password(password.as_bytes(), &parsed_hash)
        })
        .map_or(false, |_| true);

        if !is_password_correct {
            return Err(Status::internal(format!("Invalid password for: {}", email)))
        }

        let access_token_details = match token::generate_jwt_token(
            user.id,
            self.context.env.access_token_max_age,
            self.context.env.access_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };
    
        let refresh_token_details = match token::generate_jwt_token(
            user.id,
            self.context.env.refresh_token_max_age,
            self.context.env.refresh_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };

        let mut redis_client = match self.context.redis_client.get_async_connection().await {
            Ok(redis_client) => redis_client,
            Err(e) => {
                return Err(Status::internal(format!("Catching service retrival failed with error: {}", e)))
            }
        };
    
        let access_result: redis::RedisResult<()> = redis_client
            .set_ex(
                access_token_details.token_uuid.to_string(),
                user.id.to_string(),
                (self.context.env.access_token_max_age * 60) as usize,
            )
            .await;
    
        if let Err(e) = access_result {
            return Err(Status::internal(format!("Catch saving failed with error: {}", e)));
        }
    
        let refresh_result: redis::RedisResult<()> = redis_client
            .set_ex(
                refresh_token_details.token_uuid.to_string(),
                user.id.to_string(),
                (self.context.env.refresh_token_max_age * 60) as usize,
            )
            .await;
    
        if let Err(e) = refresh_result {
            return Err(Status::internal(format!("Catch saving failed with error: {}", e)));
        }

        let reply = LoginUserReply {
            access_token : access_token_details.token.clone().unwrap(),
            access_token_age : self.context.env.access_token_max_age * 60,
            refresh_token : refresh_token_details.token.unwrap(),
            refresh_token_age : self.context.env.access_token_max_age * 60,
        };

        Ok(Response::new(reply))
    }

    async fn refresh_access_token(
        &self,
        request: Request<RefreshTokenRequest>,
    ) -> Result<Response<RefreshTokenReply>, Status> {
        println!("Got a request: {:#?}", &request);
        
        let RefreshTokenRequest { refresh_token } = request.into_inner();

        let refresh_token_details =
        match token::verify_jwt_token(self.context.env.refresh_token_public_key.to_owned(), &refresh_token)
        {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token validation failed with error: {}", e)));
            }
        };

        let result = self.context.redis_client.get_async_connection().await;
        let mut redis_client = match result {
            Ok(redis_client) => redis_client,
            Err(e) => {
                return Err(Status::internal(format!("Catching service connection failed with error: {}", e)));
            }
        };
        let redis_result: redis::RedisResult<String> = redis_client
            .get(refresh_token_details.token_uuid.to_string())
            .await;

        let user_id_str = match redis_result {
            Ok(value) => value,
            Err(_) => {
                return Err(Status::internal(format!("User id retrival failed.")));
            }
        };

        let user_id = user_id_str.parse::<i32>().unwrap();
        // get a connection from the pool
        let conn = &mut self.context.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;
        // run a querry
        let user_result = users_table::table.find(user_id).get_result::<orm_user>(conn);
        let user = match user_result {
            Ok(user) => user,
            Err(err) => {
                eprintln!("Error finding user: {}", err);
                return Err(Status::internal(format!("the user belonging to this token no logger exists.")));
            }
        };

        let access_token_details = match token::generate_jwt_token(
            user.id,
            self.context.env.access_token_max_age,
            self.context.env.access_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };

        let redis_result: redis::RedisResult<()> = redis_client
            .set_ex(
                access_token_details.token_uuid.to_string(),
                user.id.to_string(),
                (self.context.env.access_token_max_age * 60) as usize,
            )
            .await;

        if redis_result.is_err() {
            return Err(Status::internal(format!("Catch saving failed")));
        }

        let reply = RefreshTokenReply {
            access_token : access_token_details.token.clone().unwrap(),
            access_token_age : self.context.env.access_token_max_age * 60,
        };

        Ok(Response::new(reply))
    }

    async fn log_out_user(
        &self,
        request: Request<LogOutRequest>,
    ) -> Result<Response<LogOutReply>, Status> {
        println!("Got a request: {:#?}", &request);
        
        let LogOutRequest { refresh_token, access_token_uuid } = request.into_inner();

        let refresh_token_details =
        match token::verify_jwt_token(self.context.env.refresh_token_public_key.to_owned(), &refresh_token)
        {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token validation failed with error: {}", e)));
            }
        };

        let result = self.context.redis_client.get_async_connection().await;
        let mut redis_client = match result {
            Ok(redis_client) => redis_client,
            Err(e) => {
                return Err(Status::internal(format!("Catching service connection failed with error: {}", e)));
            }
        };

        let redis_result: redis::RedisResult<usize> = redis_client
        .del(&[
            refresh_token_details.token_uuid.to_string(),
            access_token_uuid.to_string(),
        ])
        .await;

        if redis_result.is_err() {
            return Err(Status::internal(format!("Logout failed")));
        }

        let reply = LogOutReply {
            message : "Logout successfull".to_string()
        };

        Ok(Response::new(reply))
    }
}