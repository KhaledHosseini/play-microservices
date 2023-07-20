use std::error::Error;
//Proto models (gRPC)
use crate::proto::{ CreateUserRequest, CreateUserResponse, GetUserResponse,LogOutResponse, RoleType as RoleTypeEnumProto};

//Our ORM models.
use diesel::prelude::*;

#[derive(Queryable, Selectable, PartialEq,  Identifiable, Debug,)]
#[diesel(table_name = crate::schema::users)]
#[diesel(check_for_backend(diesel::pg::Pg))]
pub struct User {
    pub id: i32,
    pub name: String,
    pub email: String,
    pub verified: bool,
    pub password: String,
    pub role: RoleType,
}

#[derive(Insertable, Queryable, Debug, PartialEq)]
#[diesel(table_name = crate::schema::users)]
#[diesel(check_for_backend(diesel::pg::Pg))]
pub struct NewUser {
    pub name: String,
    pub email: String,
    pub password: String,
    pub role: RoleType,
}

#[derive(diesel_derive_enum::DbEnum,Debug,PartialEq)]
#[ExistingTypePath = "crate::schema::sql_types::RoleType"]
#[DbValueStyle = "snake_case"]//Important. it should be styled according to .sql: see https://docs.rs/diesel-derive-enum/2.1.0/diesel_derive_enum/derive.DbEnum.html
pub enum RoleType {
    ADMIN,
    CUSTOMER
}

impl std::convert::From<RoleTypeEnumProto> for RoleType {
    fn from(proto_form: RoleTypeEnumProto) -> Self {
        match proto_form {
            RoleTypeEnumProto::RoletypeAdmin => RoleType::ADMIN,
            RoleTypeEnumProto::RoletypeCustomer => RoleType::CUSTOMER,
        }
    }
}

impl std::convert::From<i32> for RoleTypeEnumProto {
    fn from(proto_form: i32) -> Self {
        match proto_form {
            0 => RoleTypeEnumProto::RoletypeAdmin,
            1 => RoleTypeEnumProto::RoletypeCustomer,
            _ => RoleTypeEnumProto::RoletypeCustomer,
        }
    }
}

//Convert from User ORM model to UserReply proto model
impl From<User> for GetUserResponse{
    fn from(rust_form: User) -> Self {
        GetUserResponse{
            id: rust_form.id,
            name: rust_form.name,
            email: rust_form.email
        }
    }
}

//convert from CreateUserRequest proto model to  NewUser ORM model
impl From<CreateUserRequest> for NewUser {
    fn from(proto_form: CreateUserRequest) -> Self {
        //CreateUserRequest.role is i32
        let role: RoleTypeEnumProto = proto_form.role.into();
        NewUser {
            name: proto_form.name,
            email: proto_form.email,
            password: proto_form.password,
            role: role.into(),
        }
    }
}

impl From<&str> for LogOutResponse {
    fn from(str: &str)-> Self {
        LogOutResponse{
            message: str.to_string()
        }
    }
}
impl From<&str> for CreateUserResponse {
    fn from(str: &str)-> Self {
        CreateUserResponse{
            message: str.to_string()
        }
    }
}


#[tonic::async_trait]
pub trait UserDBInterface {
    async fn is_user_exists_with_email(&self,email: &String)-> Result<bool, Box<dyn Error>>;
    async fn insert_new_user(&self,new_user: &NewUser)-> Result<usize, Box<dyn Error>>;
    async fn get_user_by_email(&self,email: &String)-> Result<std::option::Option<User>, Box<dyn Error>>;
    async fn get_user_by_id(&self,id: &i32)-> Result<User, Box<dyn Error>>;
    async fn get_users_list(&self,page: &i64, size: &i64)-> Result<Vec<User>, Box<dyn Error>>;
}

#[tonic::async_trait]
pub trait UserCacheInterface {
    async fn set_expiration(&self,key: &String,value: &String, seconds: usize)-> Result<(), Box<dyn Error>>;
    async fn get_value(&self,key: &String)-> Result<String, Box<dyn Error>>;
    async fn delete_value_for_key(&self,keys: Vec<String>)-> Result<String, Box<dyn Error>>;
}