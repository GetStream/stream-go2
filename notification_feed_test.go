package stream_test

import (
	"fmt"
	"testing"

	stream "github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNotificationActivities(t *testing.T) {
	client := newClient(t)
	notification := client.NotificationFeed("notification", randString(10))
	activities := make([]stream.Activity, 5)
	for i := range activities {
		activities[i] = stream.Activity{
			Actor:  fmt.Sprintf("actor-%d", i),
			Verb:   "like",
			Object: randString(10),
		}
	}
	_, err := notification.AddActivities(activities...)
	require.NoError(t, err)
	_, err = notification.GetActivities(stream.GetActivitiesWithLimit(123))
	assert.NoError(t, err)
}
