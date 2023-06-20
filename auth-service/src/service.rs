//grpc models
use tonic::{Request, Response, Status};
use crate::proto::{
    user_server::User, CreateUserReply, CreateUserRequest,
    UserReply, UserRequest, LoginUserRequest, LoginUserReply,
    RefreshTokenRequest,RefreshTokenReply, LogOutRequest, LogOutReply
};

//orm models
use crate::models:: {
    User as orm_user, NewUser
};
// users table
use crate::schema::{users as users_table};

use diesel::{pg::PgConnection, prelude::*};
use diesel::r2d2::{ Pool, PooledConnection, ConnectionManager, PoolError };
use diesel::dsl::select;
use diesel::dsl::exists;

//cache
use redis::{Client, AsyncCommands};

//others
use crate::token;
use crate::config::Config;
use argon2::{
    password_hash::{rand_core::OsRng, PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
    Argon2,
};


type PgPool = Pool<ConnectionManager<PgConnection>>;
type PgPooledConnection = PooledConnection<ConnectionManager<PgConnection>>;

// #[derive(Debug, Default)]
pub struct UserService {
    redis_client: Client,
    db_pool: PgPool,
    env: Config,
}

impl UserService {
    pub fn new(db_url: &str, redis_url: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let redis_client = match redis::Client::open(redis_url) {
            Ok(client) => {
                println!("✅Connection to the redis is successful!");
                client
            }
            Err(e) => {
                println!("🔥 Error connecting to Redis: {}", e);
                return Err(Box::new(e))
            }
        };

        let manager = ConnectionManager::<PgConnection>::new(db_url);
        // let db_pool = Pool::builder().build(manager)?;
        let db_pool = match Pool::builder().build(manager) {
            Ok(pool) => {
                println!("✅ Successfully created connection pool");
                pool
            }
            Err(e) => {
                println!("🔥 Error creating connection pool: {}", e);
                return Err(Box::new(e))
            }
        };

        let env = Config::init();
        Ok(Self {
            redis_client,
            db_pool,
            env,
        })
    }

     // helper function to get a connection from the pool
     fn get_conn(&self) -> Result<PgPooledConnection, PoolError> {
        let _pool = self.db_pool.get().unwrap();
        Ok(_pool)
    }
}


// Implement User trait for Our UserService
#[tonic::async_trait]
impl User for UserService {
    async fn get_user(&self, request: Request<UserRequest>) -> Result<Response<UserReply>, Status> {
        println!("Got a request: {:#?}", &request);

        //retrieve grpc model
        let UserRequest { id } = request.into_inner();

        // get a connection from the pool
        let conn = &mut self.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;
        // run a querry
        let user_result = users_table::table.find(id).get_result::<orm_user>(conn);

        match user_result {
            Ok(user) => {
                //convert from orm model to grpc model
                let u: UserReply = user.into();
                Ok(Response::new(u))
            }
            Err(err) => {
                eprintln!("Error finding user: {}", err);
                return Err(Status::not_found("User not found."))
            }
        }
    }

    async fn create_user(
        &self,
        request: Request<CreateUserRequest>
    ) -> Result<Response<CreateUserReply>, Status> {
       println!("Got a request: {:#?}", &request);
       
       let CreateUserRequest { name, email, password } = &request.into_inner();
       
       let conn = &mut self.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;

       let email_exists = select(exists(users_table::table.filter(users_table::dsl::email.eq(&email))))
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
        .hash_password(password.as_bytes(), &salt)
        .expect("Error while hashing password")
        .to_string();
        
        let new_user = NewUser {
            name: name.to_string(),
            email: email.to_string(),
            password: hashed_password
        };

        let insert_response = diesel::insert_into(users_table::dsl::users)
        .values(&new_user)
        .execute(conn);

        match insert_response {
            Ok(_)=> {
                let reply = CreateUserReply { message: format!("Successfull user registration")};
                Ok(Response::new(reply))
            }
            Err (e) => {
                return Err(Status::already_exists(format!("Problem creating user with error: {}",e)));
            }
        }
       
    }

    async fn login_user(
        &self,
        request: Request<LoginUserRequest>,
    ) -> Result<Response<LoginUserReply>, Status> {

        println!("Got a request: {:#?}", &request);

        let LoginUserRequest { email, password } = request.into_inner();

        let conn = &mut self.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;
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
            self.env.access_token_max_age,
            self.env.access_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };
    
        let refresh_token_details = match token::generate_jwt_token(
            user.id,
            self.env.refresh_token_max_age,
            self.env.refresh_token_private_key.to_owned(),
        ) {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token generation failed with error: {}", e)))
            }
        };

        let mut redis_client = match self.redis_client.get_async_connection().await {
            Ok(redis_client) => redis_client,
            Err(e) => {
                return Err(Status::internal(format!("Catching service retrival failed with error: {}", e)))
            }
        };
    
        let access_result: redis::RedisResult<()> = redis_client
            .set_ex(
                access_token_details.token_uuid.to_string(),
                user.id.to_string(),
                (self.env.access_token_max_age * 60) as usize,
            )
            .await;
    
        if let Err(e) = access_result {
            return Err(Status::internal(format!("Catch saving failed with error: {}", e)));
        }
    
        let refresh_result: redis::RedisResult<()> = redis_client
            .set_ex(
                refresh_token_details.token_uuid.to_string(),
                user.id.to_string(),
                (self.env.refresh_token_max_age * 60) as usize,
            )
            .await;
    
        if let Err(e) = refresh_result {
            return Err(Status::internal(format!("Catch saving failed with error: {}", e)));
        }

        let reply = LoginUserReply {
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
    ) -> Result<Response<RefreshTokenReply>, Status> {
        println!("Got a request: {:#?}", &request);
        
        let RefreshTokenRequest { refresh_token } = request.into_inner();

        let refresh_token_details =
        match token::verify_jwt_token(self.env.refresh_token_public_key.to_owned(), &refresh_token)
        {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token validation failed with error: {}", e)));
            }
        };

        let result = self.redis_client.get_async_connection().await;
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
        let conn = &mut self.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;
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
            self.env.access_token_max_age,
            self.env.access_token_private_key.to_owned(),
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
                (self.env.access_token_max_age * 60) as usize,
            )
            .await;

        if redis_result.is_err() {
            return Err(Status::internal(format!("Catch saving failed")));
        }

        let reply = RefreshTokenReply {
            access_token : access_token_details.token.clone().unwrap(),
            access_token_age : self.env.access_token_max_age * 60,
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
        match token::verify_jwt_token(self.env.refresh_token_public_key.to_owned(), &refresh_token)
        {
            Ok(token_details) => token_details,
            Err(e) => {
                return Err(Status::internal(format!("Token validation failed with error: {}", e)));
            }
        };

        let result = self.redis_client.get_async_connection().await;
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