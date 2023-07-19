use std::error::Error;
use log::info;

use crate::models::user::db::{User,NewUser, RoleType ,UserDBInterface};

use diesel::{pg::PgConnection, prelude::*};
use diesel::r2d2::{ Pool, PooledConnection, ConnectionManager, PoolError };
use diesel::dsl::select;
use diesel::dsl::exists;
use crate::schema::users as users_table;

pub type PgPool = Pool<ConnectionManager<PgConnection>>;
pub type PgPooledConnection = PooledConnection<ConnectionManager<PgConnection>>;
pub struct PostgressDB {
    pub db_pool: PgPool,
}

impl PostgressDB {
    // helper function to get a connection from the pool
    pub fn get_conn(&self) -> Result<PgPooledConnection, PoolError> {
       let _pool = self.db_pool.get().unwrap();
       Ok(_pool)
   }
}


#[tonic::async_trait]
impl UserDBInterface for PostgressDB {
    async fn is_user_exists_with_email(&self, email: &String)-> Result<bool,Box<dyn Error>> {
        info!("PostgressDB:is_user_exists_with_email: start");
        let conn = &mut self.get_conn()?;
        let email_exists = select(exists(users_table::table.filter(users_table::dsl::email.eq(email))))
            .get_result::<bool>(conn);
        
        Ok(email_exists?)
    }
   
    async fn insert_new_user(&self,new_user: &NewUser)-> Result<usize, Box<dyn Error>> {
        info!("PostgressDB:insert_new_user: start");
        let conn = &mut self.get_conn()?;
        let insert_response = diesel::insert_into(users_table::table)
        .values(new_user)
        .execute(conn);
        
        return Ok(insert_response?)
    }
    async fn get_user_by_email(&self,email: &String)-> Result<std::option::Option<User>, Box<dyn Error>> {
        info!("PostgressDB:get_user_by_email: start");
        let conn = &mut self.get_conn()?;
        let user_result = users_table::table.filter(users_table::dsl::email.eq(email)).first::<User>(conn).optional();
        return Ok(user_result?)
    }
    async fn get_user_by_id(&self,id: &i32)-> Result<User, Box<dyn Error>> {
        info!("PostgressDB:get_user_by_id: start");
        let conn = &mut self.get_conn()?;
        let user_result = users_table::table.find(id).get_result::<User>(conn);
        return Ok(user_result?)
    }
    async fn get_users_list(&self,page: &i64, size: &i64)-> Result<Vec<User>, Box<dyn Error>> {
        info!("PostgressDB:get_users_list: start");
        let conn = &mut self.get_conn()?;
        let offset = (*page - 1) * *size;
        let users = users_table::table
            .select((users_table::dsl::id, users_table::dsl::name, users_table::dsl::email, users_table::dsl::verified, users_table::dsl::role))
            .limit((*size).into())
            .offset(offset.into())
            .load::<(i32, String, String, bool, RoleType)>(conn)?;

        let users: Vec<User> = users
        .into_iter()
        .map(|(id, name, email,verified,role)| User { id, name, email,verified, password: String::new(),role }) // Set an empty password
        .collect();

        Ok(users)
    }
}