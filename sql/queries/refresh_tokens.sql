-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, created_at, updated_at, expired_at)
VALUES ($1, $2, NOW(), NOW(), NOW() + INTERVAL '60 day')
RETURNING user_id, token, created_at, updated_at, expired_at;

-- name: GetRefreshToken :one
SELECT user_id, token, created_at, updated_at, expired_at, revoked_at
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW() 
WHERE token = $1;
