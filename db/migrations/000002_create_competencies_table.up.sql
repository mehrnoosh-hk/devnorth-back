CREATE TABLE competencies (
    id SERIAL PRIMARY KEY,
    name CITEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_competencies_name ON competencies(name);
