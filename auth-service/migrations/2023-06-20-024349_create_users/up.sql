-- Your SQL goes here
CREATE TABLE users (
        id int NOT NULL PRIMARY KEY DEFAULT 1,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL UNIQUE,
        verified BOOLEAN NOT NULL DEFAULT FALSE,
        password VARCHAR(100) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'user'
    );
CREATE INDEX users_email_idx ON users (email);