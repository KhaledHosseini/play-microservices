// @generated automatically by Diesel CLI.

pub mod sql_types {
    #[derive(diesel::query_builder::QueryId, diesel::sql_types::SqlType)]
    #[diesel(postgres_type(name = "role_type"))]
    pub struct RoleType;
}

diesel::table! {
    use diesel::sql_types::*;
    use super::sql_types::RoleType;

    users (id) {
        id -> Int4,
        #[max_length = 100]
        name -> Varchar,
        #[max_length = 100]
        email -> Varchar,
        verified -> Bool,
        #[max_length = 100]
        password -> Varchar,
        role -> RoleType,
    }
}
