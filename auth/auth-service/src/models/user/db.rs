pub mod model;
pub mod validation;
pub mod diesel_postgres;
pub mod redis;
pub mod tests;

pub use crate::models::user::db::model::{User,NewUser,RoleType,UserDBInterface, UserCacheInterface};
pub use crate::models::user::db::validation::Validate;
pub use crate::models::user::db::diesel_postgres::{PostgressDB,PgPool,PgPooledConnection};
pub use crate::models::user::db::redis::RedisCache;