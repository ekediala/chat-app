package utils

import (
	"database/sql"
)

const (
	ROOT_USER_ROUTE    = "/users"
	CREATE_USER        = "create"
	LOGIN              = "login"
	ROOT_CHANNEL_ROUTE = "/channels"
	CREATE             = "create"
)

type FrontendUser struct {
	Username  string       `json:"username"`
	ID        int64        `json:"id"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
