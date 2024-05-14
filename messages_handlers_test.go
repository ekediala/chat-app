package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekediala/chat-app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestMessagesCreate(t *testing.T) {
	server := NewServer()
	url := utils.ComposeMessageRoute(utils.CREATE)
	server.router.POST(url, server.RequiresAuth, server.CreateMessageHandler)
	token, _ := CreateAndLoginNewTestUser(t)

	channel, err := server.Db.CreateChannel(context.Background(), faker.Username())
	require.NoError(t, err)
	requestBody := ChannelRequestBody{
		ChannelID: channel.ID,
		Message:   faker.Sentence(),
	}

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}
