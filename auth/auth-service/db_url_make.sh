#!/bin/bash

db_name_file="${DATABSE_DB_FILE}"
db_name=$(cat "$db_name_file")
db_user_file="${DATABSE_USER_FILE}"
db_user=$(cat "$db_user_file")
db_pass_file="${DATABSE_PASSWORD_FILE}"
db_pass=$(cat "$db_pass_file")

db_scheme="${DATABSE_SCHEME}"
db_domain="${DATABSE_DOMAIN}"
db_port="${DATABSE_PORT}"

db_url="${db_scheme}://${db_user}:${db_pass}@${db_domain}:${db_port}/${db_name}"
echo "database url is ${db_url}"
export DATABASE_URL=${db_url}