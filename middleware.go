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

type claimsKeyType string

var claimsKey claimsKeyType = "claims"

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

func (server *Server) registerRoutes() {
	server.router.POST(utils.ComposeUserRoute(CREATE), server.CreateUser)
	server.router.POST(utils.ComposeUserRoute(LOGIN), server.login)
	server.router.POST(utils.ComposeMessageRoute(CREATE), server.RequiresAuth, server.CreateMessageHandler)
	server.router.POST(utils.ComposeChannelRoute(CREATE), server.RequiresAuth, server.CreateChannel)
	server.router.GET(utils.ComposeChannelRoute(LIST), server.RequiresAuth, server.ListChannels)
	server.router.GET(utils.ComposeMessageRoute(LIST), server.RequiresAuth, server.ListMessagesByChannelIDHandler)
}
