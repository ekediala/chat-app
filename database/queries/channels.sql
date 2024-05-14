-- name: CreateChannel :one
INSERT INTO channels (name) VALUES (?) RETURNING *;

-- name: ListChannels :many
SELECT id, name FROM channels LIMIT ? OFFSET ?;

-- name: GetChannelByID :one
SELECT id, name FROM channels WHERE id = ?;