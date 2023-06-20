
//ORM models (database models)
use crate::models::{User, NewUser};

//Proto models (gRPC)
use crate::proto::{
    CreateUserRequest, UserReply
};

//Convert from User ORM model to UserReply proto model
impl From<User> for UserReply {
    fn from(rust_form: User) -> Self {
        UserReply{
            id: rust_form.id,
            name: rust_form.name,
            email: rust_form.email
        }
    }
}

//convert from CreateUserRequest proto model to  NewUser ORM model
impl From<CreateUserRequest> for NewUser {
    fn from(proto_form: CreateUserRequest) -> Self {
        NewUser {
            name: proto_form.name,
            email: proto_form.email,
            password: proto_form.password,
        }
    }
}
