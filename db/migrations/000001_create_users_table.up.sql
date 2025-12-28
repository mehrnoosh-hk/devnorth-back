-- Enable CITEXT extension for case-insensitive text
CREATE EXTENSION IF NOT EXISTS citext;

-- Create user role enum
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email CITEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'USER',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);
