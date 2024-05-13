package utils

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type ResponsePayload struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func RespondWithError(c *gin.Context, code int, message string) {
	Logger.Error("Error", "message", message)
	RespondWithJSON(c, code, ResponsePayload{Message: message})
}

func RespondWithJSON(c *gin.Context, code int, payload ResponsePayload) {
	if payload.Message == "" {
		payload.Message = http.StatusText(code)
	}
	c.JSON(code, payload)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ComposeUserRoute(path string) string {
	return fmt.Sprintf("%s/%s", ROOT_USER_ROUTE, path)
}

func ComposeChannelRoute(path string) string {
	return fmt.Sprintf("%s/%s", ROOT_CHANNEL_ROUTE, path)
}
