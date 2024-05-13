package database

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func createNewChannel(t *testing.T, channelName string) Channel {

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
	channel := createNewChannel(t, name)
	require.NotEmpty(t, channel)
	require.Equal(t, name, channel.Name)
	require.WithinRange(t, channel.CreatedAt.Time, time.Now().Add(-5*time.Second), time.Now())
}

func TestListChannels(t *testing.T) {

}
