-- name: CreateMessage :one
INSERT INTO messages (user_id, channel_id, message) VALUES (?, ?, ?) RETURNING *;

-- name: ListMessagesByChannelID :many
SELECT message.id, message.message, message.created_at, message.updated_at, user.id AS user_id, user.username AS user_name, channel.id AS channel_id, channel.name AS channel_name FROM messages message JOIN users user ON message.user_id = user.id JOIN channels channel ON message.channel_id = channel.id WHERE channel_id = ? ORDER BY message.created_at DESC LIMIT ? OFFSET ?;