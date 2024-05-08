-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES (?, ?) RETURNING *;
