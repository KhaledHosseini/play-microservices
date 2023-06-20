db_name_file="${POSTGRES_DB_FILE}"
db_name=$(cat "$db_name_file")
db_user_file="${POSTGRES_USER_FILE}"
db_user=$(cat "$db_user_file")
db_pass_file="${POSTGRES_PASSWORD_FILE}"
db_pass=$(cat "$db_pass_file")

db_host="${POSTGRES_HOST}"
db_port="${POSTGRES_PORT}"
db_url="postgresql://${db_user}:${db_pass}@${db_host}:${db_port}/${db_name}"
echo "database url is ${db_url}"
export DATABASE_URL=${db_url}