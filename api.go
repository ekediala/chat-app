package main

import (
	"database/sql"
	"log"
	"os"

	database "github.com/ekediala/chat-app/database/sqlc"
	"github.com/ekediala/chat-app/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

type Server struct {
	Db     *database.Queries
	router *gin.Engine
}

func NewServer() *Server {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
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

	return &Server{
		Db:     database.New(db),
		router: router,
	}
}

func (server *Server) registerRoutes() {
	server.router.POST(utils.ComposeUserRoute(utils.CREATE_USER), server.createUser)
	server.router.POST(utils.ComposeUserRoute(utils.LOGIN), server.login)

}

func (api *Server) Start() error {
	api.registerRoutes()
	return api.router.Run(":8080")
}
