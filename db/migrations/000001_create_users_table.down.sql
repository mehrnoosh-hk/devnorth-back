-- Drop users table
DROP TABLE IF EXISTS users;

-- Drop user role enum
DROP TYPE IF EXISTS user_role;

-- Note: We don't drop CITEXT extension as it might be used by other tables
