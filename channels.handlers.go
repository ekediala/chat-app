package main

import (
	"net/http"

	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
)

type CreateChannelPayload struct {
	Name string `json:"name" binding:"required"`
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

func (server *Server) ListChannels(c *gin.Context){
	
}