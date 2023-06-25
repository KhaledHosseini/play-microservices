use crate::Context;
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

use diesel::{prelude::*};


// #[derive(Debug, Default)]
pub struct UserProfileService {
    pub context: Context,
}

// Implement User trait for Our UserService
#[tonic::async_trait]
impl UserProfile for UserProfileService {
    async fn get_user(&self, request: Request<UserRequest>) -> Result<Response<UserReply>, Status> {
        println!("Got a request: {:#?}", &request);

        //retrieve grpc model
        let UserRequest { id } = request.into_inner();

        // get a connection from the pool
        let conn = &mut self.context.get_conn().map_err(|e| Status::internal(format!("Error getting connection: {}", e)))?;
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