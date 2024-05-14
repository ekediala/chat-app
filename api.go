package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-jwt/jwt/v5"
)

const (
	USER_JWT_KEY string = "user"
	LOGIN        string = "login"
	CREATE       string = "create"
	LIST         string = "list"
)

func (server *Server) createToken(user utils.FrontendUser) (string, error) {

	secretKey := server.config.JWT_SECRET

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        user.Username,                    // Subject (user identifier)
		"iss":        "chat-app",                       // Issuer          // Audience (user role)
		"exp":        time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat":        time.Now().Unix(),                // Issued at
		USER_JWT_KEY: user,
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

func getCurrentUser(c *gin.Context) (userMap map[string]interface{}, ok bool) {
	claims := c.Request.Context().Value(claimsKey)

	jwtClaims, ok := claims.(jwt.MapClaims)

	if !ok {
		return userMap, ok
	}

	user, ok := jwtClaims[USER_JWT_KEY]
	if !ok {
		return userMap, ok
	}

	userMap, ok = user.(map[string]interface{})
	if !ok {
		return userMap, ok
	}

	userId, ok := userMap["id"].(float64)
	if !ok {
		return userMap, ok
	}

	userMap["id"] = int64(userId)

	return userMap, ok
}

func (server *Server) AddHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
}
