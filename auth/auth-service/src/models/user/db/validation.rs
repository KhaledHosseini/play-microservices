use std::error::Error;
use std::fmt;
use crate::proto::{CreateUserRequest, LoginUserRequest,RefreshTokenRequest, LogOutRequest};

pub trait Validate {
    fn validate(&self)-> Result<(), ValidationError>;
}

#[derive(Debug)]
pub struct ValidationError {
    message: String,
}
impl ValidationError {
    fn new(message: &str) -> Self {
        Self {
            message: message.to_owned(),
        }
    }
}
impl fmt::Display for ValidationError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.message)
    }
}
impl Error for ValidationError {}

impl Validate for CreateUserRequest {
    fn validate(&self)-> Result<(),  ValidationError> {
        if self.name.is_empty() || self.email.is_empty() || self.password.is_empty(){
            return Err(ValidationError::new("Empty field is not allowed."))
        }

        Ok(())
    }
}

impl Validate for LoginUserRequest {
    fn validate(&self)-> Result<(),  ValidationError> {
        if self.email.is_empty() || self.password.is_empty() {
            return Err(ValidationError::new("Empty field is not allowed."))
        }

        Ok(())
    }
}

impl Validate for RefreshTokenRequest {
    fn validate(&self)-> Result<(),  ValidationError> {
        if self.refresh_token.is_empty() {
            return Err(ValidationError::new("Empty token is not allowed."))
        }

        Ok(())
    }
}

impl Validate for LogOutRequest {
    fn validate(&self)-> Result<(),  ValidationError> {
        if self.refresh_token.is_empty() {
            return Err(ValidationError::new("Empty token is not allowed."))
        }
        Ok(())
    }
}
