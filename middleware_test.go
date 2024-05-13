package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekediala/chat-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestEmptyAuthorization(t *testing.T) {
	server := NewServer()
	url := "/middleware_test"
	server.router.POST(url, server.RequiresAuth)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	defer w.Result().Body.Close()

	var res utils.ResponsePayload
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Equal(t, http.StatusText(http.StatusForbidden), res.Message)
}

func TestInvalidAuthorization(t *testing.T) {
	server := NewServer()
	url := "/middleware_test"
	server.router.POST(url, server.RequiresAuth)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	defer w.Result().Body.Close()

	var res utils.ResponsePayload
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Equal(t, http.StatusText(http.StatusUnauthorized), res.Message)
}

func TestValidAuthorization(t *testing.T) {
	server := NewServer()
	url := "/middleware_test"
	token, _ := CreateAndLoginNewTestUser(t)

	server.router.POST(url, server.RequiresAuth, func(ctx *gin.Context) {
		claims := ctx.Request.Context().Value(claimsKey)

		jwtClaims := claims.(jwt.MapClaims)

		user := jwtClaims["user"]

		utils.RespondWithJSON(ctx, http.StatusOK, utils.ResponsePayload{
			Data:    user,
			Message: http.StatusText(http.StatusOK),
		})
	})

	req, err := http.NewRequest("POST", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	defer w.Result().Body.Close()

	var res utils.ResponsePayload
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.NotEmpty(t, res.Data)

	userMap, ok := res.Data.(map[string]interface{})
	require.True(t, ok)
	for _, field := range userMap {
		require.NotEmpty(t, field)
	}
}
