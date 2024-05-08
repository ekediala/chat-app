package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekediala/chat-app/utils"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserRoute(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(utils.CREATE_USER), server.createUser)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"testuser4","password":"password"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(utils.CREATE_USER), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)

	body, err := io.ReadAll(w.Result().Body)

	defer w.Result().Body.Close()

	require.NoError(t, err)

	var res utils.ResponsePayload

	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, UserCreatedSuccessfully, res.Message)
}

func TestCreateUserRouteInvalidData(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(utils.CREATE_USER), server.createUser)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"","password":"password"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(utils.CREATE_USER), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginUserRouteInvalidData(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(utils.LOGIN), server.login)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"","password":"password"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(utils.LOGIN), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginUserUserNotExist(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(utils.LOGIN), server.login)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"abracadabrawhatever","password":"passwordghdgdgh"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(utils.LOGIN), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)

	body, err := io.ReadAll(w.Result().Body)

	defer w.Result().Body.Close()

	require.NoError(t, err)

	var res utils.ResponsePayload

	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, res.Message, UserNotFound)
}

func TestLoginUserUserIncorrectPassword(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(utils.LOGIN), server.login)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"testuser4","password":"password12"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(utils.LOGIN), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)

	body, err := io.ReadAll(w.Result().Body)

	defer w.Result().Body.Close()

	require.NoError(t, err)

	var res utils.ResponsePayload

	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, res.Message, WrongUserNameOrPassword)
}

func TestLoginUserUserCorrectPassword(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(utils.LOGIN), server.login)
	server.router.POST(utils.ComposeUserRoute(utils.CREATE_USER), server.createUser)

	w := httptest.NewRecorder()

	fake := faker.New()

	user := LoginUserPayload{
		Username: fake.Internet().Email(),
		Password: fake.Lorem().Text(10),
	}

	body, _ := json.Marshal(user)

	createUserReq, _ := http.NewRequest("POST", utils.ComposeUserRoute(utils.CREATE_USER), bytes.NewBuffer(body))
	createUserReq.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, createUserReq)

	w = httptest.NewRecorder()

	loginReq, _ := http.NewRequest("POST", utils.ComposeUserRoute(utils.LOGIN), bytes.NewBuffer(body))
	loginReq.Header.Add("Content-Type", "application/json")
	server.router.ServeHTTP(w, loginReq)

	assert.Equal(t, http.StatusOK, w.Code)
}
