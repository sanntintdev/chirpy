-- +goose Up
-- Add password column to users table
ALTER TABLE users ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
-- Remove password column from users table
ALTER TABLE users DROP COLUMN is_chirpy_red;
