-- +goose Up
-- Add password column to users table
ALTER TABLE users ADD COLUMN hashed_password TEXT NOT NULL;

-- +goose Down
-- Remove password column from users table
ALTER TABLE users DROP COLUMN hashed_password;
