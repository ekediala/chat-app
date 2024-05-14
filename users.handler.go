package main

import (
	"database/sql"
	"net/http"

	database "github.com/ekediala/chat-app/database/sqlc"
	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
)

const (
	UserCreatedSuccessfully = "User created successfully"
	UserNotFound            = "User not found"
	WrongUserNameOrPassword = "Wrong username or password"
	OK                      = "Ok"
)

type CreateUserPayload struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserPayload struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserResponse struct {
	User  utils.FrontendUser `json:"user"`
	Token string             `json:"token"`
}

func (server *Server) CreateUser(c *gin.Context) {
	var user CreateUserPayload

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return

	}

	createdUser, err := server.Db.CreateUser(c, database.CreateUserParams{
		Username: user.Username,
		Password: hashedPassword,
	})

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, utils.ResponsePayload{
		Data:    createdUser,
		Message: UserCreatedSuccessfully,
	})

}

func (server *Server) login(c *gin.Context) {
	var user LoginUserPayload

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	userFromDb, err := server.Db.GetUserByUsername(c, user.Username)

	if err == sql.ErrNoRows {
		utils.RespondWithError(c, http.StatusNotFound, UserNotFound)
		return
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	match := utils.CheckPasswordHash(user.Password, userFromDb.Password)

	if !match {
		utils.RespondWithError(c, http.StatusUnauthorized, WrongUserNameOrPassword)
		return
	}

	userResponse := utils.FrontendUser{
		ID:        userFromDb.ID,
		Username:  user.Username,
		CreatedAt: userFromDb.CreatedAt,
		UpdatedAt: userFromDb.UpdatedAt,
	}

	token, err := server.createToken(userResponse)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Error logging user in. Please try again")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, utils.ResponsePayload{
		Message: http.StatusText(http.StatusOK),
		Data: LoginUserResponse{
			User:  userResponse,
			Token: token,
		},
	})

}
