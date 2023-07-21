use dotenv::dotenv;
use std::fs;

fn get_env_var(var_name: &str) -> String {
    std::env::var(var_name).unwrap_or_else(|_| panic!("{} must be set", var_name))
}

fn get_file_contents(file_path: &str)-> String {
    fs::read_to_string(file_path)
        .unwrap_or_else(|_| panic!("Cannot read the file {}",file_path))
}

fn get_file_contents_bytes(file_path: &str)-> Vec<u8> {
    fs::read(file_path)
        .unwrap_or_else(|_| panic!("Cannot read the file {}",file_path))
}

#[derive(Debug, Clone)]
pub struct Config {
    //database url for postgres-> postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
    pub database_scheme: String,
    pub database_domain: String,
    pub database_port: String,
    pub database_name: String,
    pub database_user: String,
    pub database_pass: String,
    pub database_url: String,

    // rediss://:password=@redis-server-url:6380/0?ssl_cert_reqs=CERT_REQUIRED'
    // redis://arbitrary_usrname:password@ipaddress:6379/0 -> 0: database index
    pub redis_scheme: String,
    pub redis_domain: String,
    pub redis_port: String,
    pub redis_password: String,
    pub redis_url: String,

    pub server_port: String,
    
    pub access_token_private_key: Vec<u8>,
    pub access_token_public_key: Vec<u8>,
    pub access_token_expires_in: String,
    pub access_token_max_age: i64,

    pub refresh_token_private_key: Vec<u8>,
    pub refresh_token_public_key: Vec<u8>,
    pub refresh_token_expires_in: String,
    pub refresh_token_max_age: i64,
}

impl Config {
    pub fn init() -> Config {
        dotenv().ok();

        let database_scheme = get_env_var("DATABSE_SCHEME");
        let database_domain = get_env_var("DATABSE_DOMAIN");
        let database_port = get_env_var("DATABSE_PORT");
        let database_name_file = get_env_var("DATABSE_DB_FILE");
        let database_name = get_file_contents(&database_name_file);
        let database_user_file = get_env_var("DATABSE_USER_FILE");
        let database_user = get_file_contents(&database_user_file);
        let database_pass_file = get_env_var("DATABSE_PASSWORD_FILE");
        let database_pass = get_file_contents(&database_pass_file);
        //database url for postgres-> postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
        let database_url = format!("{}://{}:{}@{}:{}/{}",
        database_scheme,
        database_user,
        database_pass,
        database_domain,
        database_port,
        database_name);

        
        let redis_scheme = get_env_var("REDIS_SCHEME");
        let redis_domain = get_env_var("REDIS_DOMAIN");
        let redis_port = get_env_var("REDIS_PORT");
        let redis_password_file = get_env_var("REDIS_PASS_FILE");
        let redis_password = get_file_contents(&redis_password_file);
        // rediss://:password=@redis-server-url:6380/0?ssl_cert_reqs=CERT_REQUIRED'
        // redis://arbitrary_usrname:password@ipaddress:6379/0 -> 0: database index
        let redis_url = format!("{}://:{}@{}:{}", 
        redis_scheme, 
        redis_password, 
        redis_domain,
        redis_port);

        let server_port = get_env_var("SERVER_PORT");

        let at_pub_file = get_env_var("ACCESS_TOKEN_PRIVATE_KEY_FILE");
        let access_token_private_key = get_file_contents_bytes(&at_pub_file);
        let ac_prvt_file = get_env_var("ACCESS_TOKEN_PUBLIC_KEY_FILE");
        let access_token_public_key = get_file_contents_bytes(&ac_prvt_file);
        let access_token_expires_in = get_env_var("ACCESS_TOKEN_EXPIRED_IN");
        let access_token_max_age = get_env_var("ACCESS_TOKEN_MAXAGE");

        let rt_pub_file = get_env_var("REFRESH_TOKEN_PRIVATE_KEY_FILE");
        let refresh_token_private_key = get_file_contents_bytes(&rt_pub_file);
        let rt_prvt_file = get_env_var("REFRESH_TOKEN_PUBLIC_KEY_FILE");
        let refresh_token_public_key = get_file_contents_bytes(&rt_prvt_file);
        let refresh_token_expires_in = get_env_var("REFRESH_TOKEN_EXPIRED_IN");
        let refresh_token_max_age = get_env_var("REFRESH_TOKEN_MAXAGE");

        Config {
            database_scheme,
            database_domain,
            database_port,
            database_name,
            database_user,
            database_pass,
            database_url,
            
            redis_scheme,
            redis_domain,
            redis_port,
            redis_password,
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
