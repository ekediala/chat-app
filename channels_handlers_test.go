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

	"github.com/ekediala/chat-app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateChannel(t *testing.T) {
	server := NewServer()
	server.router.POST(utils.ComposeChannelRoute(CREATE), server.RequiresAuth, server.CreateChannel)
	token, _ := CreateAndLoginNewTestUser(t)

	w := httptest.NewRecorder()

	user := CreateChannelPayload{
		Name: faker.Name(),
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", utils.ComposeChannelRoute(CREATE), bytes.NewBuffer(body))
	server.AddHeaders(req, token)

	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestListChannelsHandlerWrongOffset(t *testing.T) {
	server := NewServer()
	target := utils.ComposeChannelRoute(LIST)
	server.router.GET(target, server.RequiresAuth, server.ListChannels)
	token, _ := CreateAndLoginNewTestUser(t)

	data := url.Values{}
	data.Set("limit", "10")
	data.Set("page", "0")

	params := data.Encode()

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	w := httptest.NewRecorder()

	server.AddHeaders(req, token)
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestListChannelsHandlerWrongLimit(t *testing.T) {
	server := NewServer()
	target := utils.ComposeChannelRoute(LIST)
	server.router.GET(target, server.RequiresAuth, server.ListChannels)
	token, _ := CreateAndLoginNewTestUser(t)

	data := url.Values{}
	data.Set("limit", "-1")
	data.Set("page", "1")

	params := data.Encode()

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	w := httptest.NewRecorder()

	server.AddHeaders(req, token)
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestListChannelsHandler(t *testing.T) {
	server := NewServer()
	target := utils.ComposeChannelRoute(LIST)
	server.router.GET(target, server.RequiresAuth, server.ListChannels)
	token, _ := CreateAndLoginNewTestUser(t)
	n := 5

	for range n {
		server.Db.CreateChannel(context.Background(), faker.Username())
	}

	data := url.Values{}
	data.Set("limit", fmt.Sprintf("%d", n))
	data.Set("page", "1")

	params := data.Encode()

	req := httptest.NewRequest("GET", target+fmt.Sprintf("?%s", params), nil)
	w := httptest.NewRecorder()

	server.AddHeaders(req, token)
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	responseBody, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)

	defer w.Result().Body.Close()

	var res utils.ResponsePayload
	err = json.Unmarshal(responseBody, &res)
	require.NoError(t, err)
	require.Len(t, res.Data, n)
}
