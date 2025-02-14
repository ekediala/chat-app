// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: channels.sql

package database

import (
	"context"
)

const createChannel = `-- name: CreateChannel :one
INSERT INTO channels (name) VALUES (?) RETURNING id, name, created_at
`

func (q *Queries) CreateChannel(ctx context.Context, name string) (Channel, error) {
	row := q.db.QueryRowContext(ctx, createChannel, name)
	var i Channel
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}

const getChannelByID = `-- name: GetChannelByID :one
SELECT id, name FROM channels WHERE id = ?
`

type GetChannelByIDRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) GetChannelByID(ctx context.Context, id int64) (GetChannelByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getChannelByID, id)
	var i GetChannelByIDRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listChannels = `-- name: ListChannels :many
SELECT id, name FROM channels ORDER BY created_at ASC LIMIT ? OFFSET ?
`

type ListChannelsParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type ListChannelsRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) ListChannels(ctx context.Context, arg ListChannelsParams) ([]ListChannelsRow, error) {
	rows, err := q.db.QueryContext(ctx, listChannels, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListChannelsRow{}
	for rows.Next() {
		var i ListChannelsRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
