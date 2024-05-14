package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	database "github.com/ekediala/chat-app/database/sqlc"
	"github.com/ekediala/chat-app/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	JWT_SECRET string
	ENV        string
}

type Server struct {
	Db     *database.Queries
	router *gin.Engine
	config AppConfig
}

const (
	USER_JWT_KEY string = "user"
)

func NewServer() *Server {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sqliteLocation := workingDirectory + "/chat.db"

	db, err := sql.Open("sqlite", sqliteLocation)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	// config.AllowAllOrigins = true

	router.Use(cors.New(config))

	env := os.Getenv("ENV")

	gin.SetMode(gin.ReleaseMode)

	if env == "development" {
		gin.SetMode(gin.DebugMode)
	}

	return &Server{
		Db:     database.New(db),
		router: router,
		config: AppConfig{
			JWT_SECRET: os.Getenv("JWT_SECRET"),
			ENV:        env,
		},
	}
}

func (server *Server) registerRoutes() {
	server.router.POST(utils.ComposeUserRoute(utils.CREATE_USER), server.CreateUser)
	server.router.POST(utils.ComposeUserRoute(utils.LOGIN), server.login)

}

func (api *Server) Start() error {
	api.registerRoutes()
	return api.router.Run(":8080")
}

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
