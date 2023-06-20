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
    pub role: String,
}

#[derive(Debug, Insertable, PartialEq)]
#[diesel(table_name = crate::schema::users)]
pub struct NewUser {
    pub name: String,
    pub email: String,
    pub password: String,
}
