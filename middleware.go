package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
)

const (
	FORBIDDEN = "Invalid token. Please try again later."
)

type ClaimsKey string

var claimsKey ClaimsKey = "claims"

func (server *Server) RequiresAuth(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")

	if bearerToken == "" {
		utils.RespondWithError(c, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		c.Abort()
		return
	}

	bearerToken = strings.TrimSpace(bearerToken)
	token := bearerToken[7:] // Bearer + " " is 6 indices. We do +1 so we start at the beginning of token
	claims, err := server.verifyToken(token)

	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		c.Abort()
		return
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), claimsKey, claims))

	c.Next()
}
