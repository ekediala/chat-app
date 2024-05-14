package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekediala/chat-app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateAndLoginNewTestUser(t *testing.T) (token string, u interface{}) {
	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(LOGIN), server.login)
	server.router.POST(utils.ComposeUserRoute(CREATE), server.CreateUser)

	w := httptest.NewRecorder()

	user := LoginUserPayload{
		Username: faker.Email(),
		Password: faker.Password(),
	}

	body, _ := json.Marshal(user)

	createUserReq, _ := http.NewRequest("POST", utils.ComposeUserRoute(CREATE), bytes.NewBuffer(body))
	createUserReq.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, createUserReq)

	w = httptest.NewRecorder()

	loginReq, _ := http.NewRequest("POST", utils.ComposeUserRoute(LOGIN), bytes.NewBuffer(body))
	loginReq.Header.Add("Content-Type", "application/json")
	server.router.ServeHTTP(w, loginReq)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)

	defer w.Result().Body.Close()

	var res utils.ResponsePayload
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	data, ok := res.Data.(map[string]interface{})
	require.Equal(t, ok, true)
	require.NotEmpty(t, data["user"])
	require.NotEmpty(t, data["token"])
	return data["token"].(string), data["user"]
}

func TestCreateUserRoute(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(CREATE), server.CreateUser)

	w := httptest.NewRecorder()

	user := LoginUserPayload{
		Username: faker.Email(),
		Password: faker.Password(),
	}

	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(CREATE), bytes.NewBuffer(body))
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

	server.router.POST(utils.ComposeUserRoute(CREATE), server.CreateUser)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"","password":"password"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(CREATE), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginUserRouteInvalidData(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(LOGIN), server.login)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"","password":"password"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(LOGIN), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	server.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLoginUserUserNotExist(t *testing.T) {

	server := NewServer()

	server.router.POST(utils.ComposeUserRoute(LOGIN), server.login)

	w := httptest.NewRecorder()

	body := []byte(`{"username":"abracadabrawhatever","password":"passwordghdgdgh"}`)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(LOGIN), bytes.NewBuffer(body))
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

	server.router.POST(utils.ComposeUserRoute(LOGIN), server.login)

	w := httptest.NewRecorder()

	_, u := CreateAndLoginNewTestUser(t)

	userMap, ok := u.(map[string]interface{})
	require.True(t, ok)

	user := LoginUserPayload{
		Username: userMap["username"].(string),
		Password: faker.Password(),
	}

	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", utils.ComposeUserRoute(LOGIN), bytes.NewBuffer(body))
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
	token, user := CreateAndLoginNewTestUser(t)
	require.NotEmpty(t, token)
	require.NotEmpty(t, user)
}
