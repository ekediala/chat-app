-- name: CreateMessage :one
INSERT INTO messages (user_id, channel_id, message) VALUES (?, ?, ?) RETURNING *;