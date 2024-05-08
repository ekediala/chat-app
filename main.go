package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	slog.SetDefault(utils.Logger)

	slog.Info("Info", "working directory", workingDirectory)

	sqliteLocation := workingDirectory + "/chat.db"

	db, err := sql.Open("sqlite", sqliteLocation)

	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	r := gin.Default()

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
