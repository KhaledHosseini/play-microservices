//grpc models
use tonic::{Request, Response, Status};
use crate::proto::{
    user_profile_server::UserProfile, UserReply, UserRequest
};

//orm models
use crate::models:: {
    User as orm_user
};
// users table
use crate::schema::{users as users_table};

use diesel::{pg::PgConnection, prelude::*};
use diesel::r2d2::{ Pool, PooledConnection, ConnectionManager, PoolError };


//others
use crate::config::Config;

type PgPool = Pool<ConnectionManager<PgConnection>>;
type PgPooledConnection = PooledConnection<ConnectionManager<PgConnection>>;

// #[derive(Debug, Default)]
pub struct UserProfileService {
    db_pool: PgPool,
}

impl UserProfileService {
    pub fn new(db_url: &str) -> Result<Self, Box<dyn std::error::Error>> {

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

        Ok(Self {
            db_pool,
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
impl UserProfile for UserProfileService {
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
}