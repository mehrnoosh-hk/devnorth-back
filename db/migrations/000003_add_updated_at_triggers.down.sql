-- Drop triggers
DROP TRIGGER IF EXISTS update_competencies_updated_at ON competencies;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();
