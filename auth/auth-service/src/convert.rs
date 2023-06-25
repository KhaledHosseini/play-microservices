
//ORM models (database models)
use crate::models::{User, NewUser, RoleType as RoleTypeEnumRust};


//Proto models (gRPC)
use crate::proto::{
    CreateUserRequest, UserReply, RoleType as RoleTypeEnumProto
};


impl std::convert::From<RoleTypeEnumProto> for RoleTypeEnumRust {
    fn from(proto_form: RoleTypeEnumProto) -> Self {
        match proto_form {
            RoleTypeEnumProto::Admin => RoleTypeEnumRust::ADMIN,
            RoleTypeEnumProto::Customer => RoleTypeEnumRust::CUSTOMER,
        }
    }
}

impl std::convert::From<i32> for RoleTypeEnumProto {
    fn from(proto_form: i32) -> Self {
        match proto_form {
            0 => RoleTypeEnumProto::Admin,
            1 => RoleTypeEnumProto::Customer,
            _ => RoleTypeEnumProto::Customer,
        }
    }
}

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