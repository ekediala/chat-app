package utils

import (
	"database/sql"
)

const (
	ROOT_USER_ROUTE     = "/users"
	ROOT_CHANNEL_ROUTE  = "/channels"
	ROOT_MESSAGES_ROUTE = "/messages"
)

type FrontendUser struct {
	Username  string       `json:"username"`
	ID        int64        `json:"id"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
