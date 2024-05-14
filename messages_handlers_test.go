package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	database "github.com/ekediala/chat-app/database/sqlc"
	"github.com/ekediala/chat-app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestMessagesCreate(t *testing.T) {
	server := NewServer()
	url := utils.ComposeMessageRoute(CREATE)
	server.router.POST(url, server.RequiresAuth, server.CreateMessageHandler)
	token, _ := CreateAndLoginNewTestUser(t)

	channel, err := server.Db.CreateChannel(context.Background(), faker.Username())
	require.NoError(t, err)
	requestBody := CreateMessageRequestBody{
		ChannelID: channel.ID,
		Message:   faker.Sentence(),
	}

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
	server.AddHeaders(req, token)
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestListMessagesByChannelIDHandler(t *testing.T) {
	n := 5

	server := NewServer()

	channelName := faker.Username()
	channel, err := server.Db.CreateChannel(context.Background(), channelName)
	require.NoError(t, err)

	token, user := CreateAndLoginNewTestUser(t)

	userMap, ok := user.(map[string]interface{})
	require.True(t, ok)

	userId, ok := userMap["id"].(float64)
	require.True(t, ok)

	for range n {
		_, err := server.Db.CreateMessage(context.Background(), database.CreateMessageParams{
			UserID:    int64(userId),
			ChannelID: channel.ID,
			Message:   faker.Sentence(),
		})
		require.NoError(t, err)
	}

	urlValues := url.Values{}
	urlValues.Set("limit", fmt.Sprintf("%d", n))
	urlValues.Set("page", "1")
	urlValues.Set("channel_id", fmt.Sprintf("%d", channel.ID))
	params := urlValues.Encode()

	target := utils.ComposeMessageRoute(LIST)

	server.router.GET(target, server.RequiresAuth, server.ListMessagesByChannelIDHandler)

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	server.AddHeaders(req, token)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var res utils.ResponsePayload
	responseBody, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	err = json.Unmarshal(responseBody, &res)
	require.NoError(t, err)
	require.Len(t, res.Data, n)
}

func TestListMessagesInvalidLimit(t *testing.T) {
	server := NewServer()

	token, _ := CreateAndLoginNewTestUser(t)

	urlValues := url.Values{}
	urlValues.Set("limit", "0")
	urlValues.Set("page", "1")
	urlValues.Set("channel_id", fmt.Sprintf("%d", 3))
	params := urlValues.Encode()

	target := utils.ComposeMessageRoute(LIST)

	server.router.GET(target, server.RequiresAuth, server.ListMessagesByChannelIDHandler)

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	server.AddHeaders(req, token)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestListMessagesInvalidPage(t *testing.T) {
	server := NewServer()

	token, _ := CreateAndLoginNewTestUser(t)

	urlValues := url.Values{}
	urlValues.Set("limit", "1")
	urlValues.Set("page", "-1")
	urlValues.Set("channel_id", fmt.Sprintf("%d", 5))
	params := urlValues.Encode()

	target := utils.ComposeMessageRoute(LIST)

	server.router.GET(target, server.RequiresAuth, server.ListMessagesByChannelIDHandler)

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	server.AddHeaders(req, token)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestListMessagesInvalidChannelID(t *testing.T) {
	server := NewServer()

	token, _ := CreateAndLoginNewTestUser(t)

	urlValues := url.Values{}
	urlValues.Set("limit", "1")
	urlValues.Set("page", "1")
	urlValues.Set("channel_id", "channel")
	params := urlValues.Encode()

	target := utils.ComposeMessageRoute(LIST)

	server.router.GET(target, server.RequiresAuth, server.ListMessagesByChannelIDHandler)

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	server.AddHeaders(req, token)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestListMessagesEmptyChannel(t *testing.T) {
	server := NewServer()
	token, _ := CreateAndLoginNewTestUser(t)

	channel, err := server.Db.CreateChannel(context.Background(), faker.Username())
	require.NoError(t, err)

	urlValues := url.Values{}
	urlValues.Set("limit", "1")
	urlValues.Set("page", "1")
	urlValues.Set("channel_id", fmt.Sprintf("%d", channel.ID))
	params := urlValues.Encode()

	target := utils.ComposeMessageRoute(LIST)

	server.router.GET(target, server.RequiresAuth, server.ListMessagesByChannelIDHandler)

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	server.AddHeaders(req, token)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	var res utils.ResponsePayload
	responseBody, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	err = json.Unmarshal(responseBody, &res)
	require.NoError(t, err)
	require.Len(t, res.Data, 0)

}
