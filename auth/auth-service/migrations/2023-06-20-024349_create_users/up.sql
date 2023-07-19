CREATE TYPE role_type AS ENUM (
  'admin',
  'customer'
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  password VARCHAR(100) NOT NULL,
  role role_type NOT NULL DEFAULT 'customer'
);

CREATE INDEX users_email_idx ON users (email);