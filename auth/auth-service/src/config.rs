use dotenv::dotenv;
use std::fs;

fn get_env_var(var_name: &str) -> String {
    std::env::var(var_name).unwrap_or_else(|_| panic!("{} must be set", var_name))
}

fn get_file_contents(file_path: &str)-> String {
    fs::read_to_string(file_path)
        .unwrap_or_else(|_| panic!("Cannot read the file {}",file_path))
}

#[derive(Debug, Clone)]
pub struct Config {
    pub database_host: String,
    pub database_port: String,
    pub database_name: String,
    pub database_user: String,
    pub database_pass: String,
    pub database_url: String,

    pub redis_host_name: String,
    pub redis_password: String,
    pub redis_url_scheme: String,
    pub redis_port: String,
    pub redis_url: String,

    pub server_port: String,
    
    pub access_token_private_key: String,
    pub access_token_public_key: String,
    pub access_token_expires_in: String,
    pub access_token_max_age: i64,

    pub refresh_token_private_key: String,
    pub refresh_token_public_key: String,
    pub refresh_token_expires_in: String,
    pub refresh_token_max_age: i64,
}

impl Config {
    pub fn init() -> Config {
        dotenv().ok();

        let database_host = get_env_var("POSTGRES_HOST");
        let database_port = get_env_var("POSTGRES_PORT");

        let database_name_file = get_env_var("POSTGRES_DB_FILE");
        let database_name = get_file_contents(&database_name_file);
        let database_user_file = get_env_var("POSTGRES_USER_FILE");
        let database_user = get_file_contents(&database_user_file);
        let database_pass_file = get_env_var("POSTGRES_PASSWORD_FILE");
        let database_pass = get_file_contents(&database_pass_file);

        let database_url = format!("postgresql://{}:{}@{}:{}/{}",
        database_user,
        database_pass,
        database_host,
        database_port,
        database_name);

        let redis_url_scheme = get_env_var("REDIS_URI_SCHEME");
        let redis_host_name = get_env_var("REDIS_HOST");
        let redis_port = get_env_var("REDIS_PORT");
        let redis_password_file = get_env_var("REDIS_PASS_FILE");
        let redis_password = get_file_contents(&redis_password_file);
        
        let redis_url = format!("{}://:{}@{}:{}", 
        redis_url_scheme, 
        redis_password, 
        redis_host_name,
        redis_port);

        let server_port = get_env_var("SERVER_PORT");

        let access_token_private_key = get_env_var("ACCESS_TOKEN_PRIVATE_KEY");
        let access_token_public_key = get_env_var("ACCESS_TOKEN_PUBLIC_KEY");
        let access_token_expires_in = get_env_var("ACCESS_TOKEN_EXPIRED_IN");
        let access_token_max_age = get_env_var("ACCESS_TOKEN_MAXAGE");

        let refresh_token_private_key = get_env_var("REFRESH_TOKEN_PRIVATE_KEY");
        let refresh_token_public_key = get_env_var("REFRESH_TOKEN_PUBLIC_KEY");
        let refresh_token_expires_in = get_env_var("REFRESH_TOKEN_EXPIRED_IN");
        let refresh_token_max_age = get_env_var("REFRESH_TOKEN_MAXAGE");

        Config {
            database_host,
            database_port,
            database_name,
            database_user,
            database_pass,
            database_url,
            
            redis_host_name,
            redis_password,
            redis_url_scheme,
            redis_port,
            redis_url,

            server_port,

            access_token_private_key,
            access_token_public_key,
            refresh_token_private_key,
            refresh_token_public_key,
            access_token_expires_in,
            refresh_token_expires_in,
            access_token_max_age: access_token_max_age.parse::<i64>().unwrap(),
            refresh_token_max_age: refresh_token_max_age.parse::<i64>().unwrap(),
        }
    }
}
