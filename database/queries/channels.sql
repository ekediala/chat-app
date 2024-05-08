-- name: CreateChannel :one
INSERT INTO channels (name) VALUES (?) RETURNING *;
