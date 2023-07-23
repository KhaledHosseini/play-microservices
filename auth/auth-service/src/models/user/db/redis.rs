use std::error::Error;
use log::info;
use redis::{Client, RedisError, Connection, Commands};
use crate::models::user::db::UserCacheInterface;

pub struct RedisCache {
    pub redis_client: Client,
}

impl RedisCache {
    // helper function to get a connection from the pool
    pub fn get_conn(&self) -> Result<Connection, RedisError> {
        self.redis_client.get_connection()
    }
}

#[tonic::async_trait]
impl UserCacheInterface for RedisCache {
    async fn set_expiration(&self,key: &String,value: &String, seconds: usize)-> Result<(), Box<dyn Error>>{
        info!("RedisCache: set_expiration");
        // Get a Redis connection
        let mut connection = self.get_conn()?;
        let _: () = connection.set_ex(key, value, seconds)?;
        Ok(())
    }

    async fn get_value(&self,key: &String)-> Result<Option<String>, Box<dyn Error>> {
        info!("RedisCache: get_value");
        let mut connection = self.get_conn()?;
        let result = connection.get(key)?;
        match result {
            Some(value) => Ok(Some(value)),
            None => Ok(None),
        }
    }

    async fn delete_value_for_key(&self,key: &String)-> Result<u64, Box<dyn Error>> {
        info!("RedisCache: delete_value_for_key");
        let mut connection = self.get_conn()?;
        let result = connection.del(key)?;
        Ok(result)
    }
}