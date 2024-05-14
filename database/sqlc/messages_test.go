package database

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func CreateRandomMessage(t *testing.T) Message {
	channelName := faker.Name()
	channel := CreateNewChannel(t, channelName)
	user := createNewUser(t, CreateUserParams{})
	text := faker.Sentence()

	message, err := testQueries.CreateMessage(context.Background(), CreateMessageParams{
		UserID:    user.ID,
		ChannelID: channel.ID,
		Message:   text,
	})

	require.NoError(t, err)
	return message
}

func TestCreateMessage(t *testing.T) {
	message := CreateRandomMessage(t)
	require.NotEmpty(t, message)
	require.NotEmpty(t, message.CreatedAt.Time)
	require.NotEmpty(t, message.UpdatedAt.Time)
}

func TestListMessagesByChannelID(t *testing.T) {
	channelName := faker.Name()
	channel := CreateNewChannel(t, channelName)
	user := createNewUser(t, CreateUserParams{})
	n := 10

	for range n {
		text := faker.Sentence()
		_, err := testQueries.CreateMessage(context.Background(), CreateMessageParams{
			UserID:    user.ID,
			ChannelID: channel.ID,
			Message:   text,
		})
		require.NoError(t, err)
	}

	messages, err := testQueries.ListMessagesByChannelID(context.Background(), ListMessagesByChannelIDParams{
		Limit:     int64(n),
		Offset:    0,
		ChannelID: channel.ID,
	})

	require.NoError(t, err)
	require.Len(t, messages, n)
}
