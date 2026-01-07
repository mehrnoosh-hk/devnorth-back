-- Recreate indexes (for rollback purposes)
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_competencies_name ON competencies(name);
