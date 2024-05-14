package database

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func CreateNewChannel(t *testing.T, channelName string) Channel {

	name := channelName

	if name == "" {
		name = faker.Name()
	}

	channel, err := testQueries.CreateChannel(context.Background(), name)

	require.NoError(t, err)

	require.Equal(t, channel.Name, name)

	return channel
}

func TestCreateChannel(t *testing.T) {
	name := faker.Name()
	channel := CreateNewChannel(t, name)
	require.NotEmpty(t, channel)
	require.Equal(t, name, channel.Name)
	require.WithinRange(t, channel.CreatedAt.Time, time.Now().Add(-5*time.Second), time.Now())
}

func TestListChannels(t *testing.T) {
	for i := 0; i < 5; i++ {
		name := faker.Name()
		CreateNewChannel(t, name)
	}

	channels, err := testQueries.ListChannels(context.Background(), ListChannelsParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, channels, 5)
}

func TestGetChannelByID(t *testing.T) {
	name := faker.Name()
	channel := CreateNewChannel(t, name)

	dbChannel, err := testQueries.GetChannelByID(context.Background(), channel.ID)
	require.NoError(t, err)
	require.Equal(t, dbChannel.ID, channel.ID)
}
