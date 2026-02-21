-- name: GetUserByUsername :one
SELECT * FROM user
WHERE username = sqlc.arg('username')
LIMIT 1;

-- name: InsertUser :one
INSERT INTO user (username, password_hash, email)
VALUES (sqlc.arg('username'), sqlc.arg('password_hash'), sqlc.arg('email'))
RETURNING *;

-- name: UpdateLastLogin :exec
UPDATE user
SET last_login = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id');
