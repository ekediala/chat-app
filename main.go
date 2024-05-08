package main

import (
	"log"
	"log/slog"

	"github.com/ekediala/chat-app/utils"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	slog.SetDefault(utils.Logger)

	server := NewServer()

	err := server.Start()

	if err != nil {
		log.Fatal(err)
	}
}
