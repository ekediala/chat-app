-- name: CreateChannel :one
INSERT INTO channels (name) VALUES (?) RETURNING *;

-- name: ListChannels :many
SELECT id, name FROM channels;