-- Remove redundant index on users.email (UNIQUE constraint already creates an index)
DROP INDEX IF EXISTS idx_users_email;

-- Remove redundant index on competencies.name (UNIQUE constraint already creates an index)
DROP INDEX IF EXISTS idx_competencies_name;
