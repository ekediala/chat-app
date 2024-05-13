package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	database "github.com/ekediala/chat-app/database/sqlc"
	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserCreatedSuccessfully = "User created successfully"
	UserNotFound            = "User not found"
	WrongUserNameOrPassword = "Wrong username or password"
	OK                      = "Ok"
)

type CreateUserPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserResponse struct {
	User  utils.FrontendUser `json:"user"`
	Token string             `json:"token"`
}

func (server *Server) createToken(user utils.FrontendUser) (string, error) {

	secretKey := server.config.JWT_SECRET

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Username,                    // Subject (user identifier)
		"iss":  "chat-app",                       // Issuer          // Audience (user role)
		"exp":  time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat":  time.Now().Unix(),                // Issued at
		"user": user,
	})

	return claims.SignedString([]byte(secretKey))

}

func (server *Server) verifyToken(tokenString string) (jwt.Claims, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.config.JWT_SECRET), nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token.Claims, nil
}

func (server *Server) createUser(c *gin.Context) {
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
