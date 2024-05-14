package main

import (
	"database/sql"
	"net/http"

	database "github.com/ekediala/chat-app/database/sqlc"
	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
)

type CreateMessageRequestBody struct {
	ChannelID int64  `json:"channel_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type ListChannelRequestParams struct {
	ListRequestParams
	ChannelID int64 `form:"channel_id" binding:"required,min=1"`
}

func (server *Server) CreateMessageHandler(c *gin.Context) {
	var body CreateMessageRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.RespondWithError(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		return
	}

	user, ok := getCurrentUser(c)
	if !ok {
		utils.RespondWithError(c, http.StatusExpectationFailed, http.StatusText(http.StatusExpectationFailed))
		return
	}

	userId, ok := user["id"].(int64)
	if !ok {
		utils.RespondWithError(c, http.StatusExpectationFailed, http.StatusText(http.StatusExpectationFailed))
		return
	}

	channel, err := server.Db.GetChannelByID(c, body.ChannelID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		utils.RespondWithError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	message, err := server.Db.CreateMessage(c, database.CreateMessageParams{
		UserID:    userId,
		ChannelID: channel.ID,
		Message:   body.Message,
	})

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, utils.ResponsePayload{
		Data:    message,
		Message: http.StatusText(http.StatusCreated),
	})
}

func (server *Server) ListMessagesByChannelIDHandler(c *gin.Context) {
	var params ListChannelRequestParams

	if err := c.ShouldBindQuery(&params); err != nil {
		utils.FailedValidationResponse(c, err.Error())
		return
	}

	messages, err := server.Db.ListMessagesByChannelID(c, database.ListMessagesByChannelIDParams{
		Offset:    (params.Page - 1) * params.Limit,
		Limit:     params.Limit,
		ChannelID: params.ChannelID,
	})

	if err != nil {
		utils.InternalServerResponse(c, "")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.ResponsePayload{
		Data: messages,
	})

}
