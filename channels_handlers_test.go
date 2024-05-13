package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekediala/chat-app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateChannel(t *testing.T) {
	server := NewServer()
	server.router.POST(utils.ComposeChannelRoute(utils.CREATE), server.RequiresAuth, server.CreateChannel)
	token, _ := CreateAndLoginNewTestUser(t)

	w := httptest.NewRecorder()

	user := CreateChannelPayload{
		Name: faker.Name(),
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", utils.ComposeChannelRoute(utils.CREATE), bytes.NewBuffer(body))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}
