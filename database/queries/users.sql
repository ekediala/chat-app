-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES (?, ?) RETURNING username, id, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;