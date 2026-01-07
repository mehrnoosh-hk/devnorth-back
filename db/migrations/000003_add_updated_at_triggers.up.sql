-- Create a reusable function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add trigger for users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add trigger for competencies table
CREATE TRIGGER update_competencies_updated_at
    BEFORE UPDATE ON competencies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
