package database

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func createNewUser(t *testing.T, user CreateUserParams) CreateUserRow {

	if user == (CreateUserParams{}) {
		user = CreateUserParams{
			Username: faker.Email(),
			Password: faker.Password(),
		}
	}

	createdUser, err := testQueries.CreateUser(context.Background(), user)

	require.NoError(t, err)
	require.Equal(t, user.Username, createdUser.Username)
	require.NotEmpty(t, createdUser)

	return createdUser
}

func TestCreateUser(t *testing.T) {

	user := CreateUserParams{
		Username: faker.Email(),
		Password: faker.Password(),
	}

	result := createNewUser(t, user)
	require.NotEmpty(t, result)
	require.WithinRange(t, result.CreatedAt.Time, time.Now().Add(-5*time.Second), time.Now())
}

func TestGetUser(t *testing.T) {

	user := CreateUserParams{
		Username: faker.Email(),
		Password: faker.Password(),
	}

	result := createNewUser(t, user)

	userFromDb, err := testQueries.GetUserByUsername(context.Background(), result.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userFromDb)
	require.Equal(t, result.ID, userFromDb.ID)

}
