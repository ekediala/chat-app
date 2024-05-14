package main

import (
	"net/http"

	database "github.com/ekediala/chat-app/database/sqlc"
	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
)

type CreateChannelPayload struct {
	Name string `json:"name" binding:"required"`
}

type ListChannelsRequestParams struct {
	Limit int64 `form:"limit" binding:"required,min=0"`
	Page  int64 `form:"page" binding:"required,min=1"`
}

func (server *Server) CreateChannel(c *gin.Context) {
	var data CreateChannelPayload

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.RespondWithError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	channel, err := server.Db.CreateChannel(c, data.Name)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, utils.ResponsePayload{
		Data:    channel,
		Message: http.StatusText(http.StatusCreated),
	})
}

func (server *Server) ListChannels(c *gin.Context) {
	var requestParams ListChannelsRequestParams
	if err := c.ShouldBindQuery(&requestParams); err != nil {
		utils.RespondWithError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	channels, err := server.Db.ListChannels(c, database.ListChannelsParams{
		Limit:  requestParams.Limit,
		Offset: (requestParams.Page - 1) * requestParams.Limit,
	})

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.ResponsePayload{
		Data:    channels,
		Message: http.StatusText(http.StatusOK),
	})
}
