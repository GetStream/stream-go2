package stream_test

import (
	"testing"

	stream "github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareFeedForTestGetActivities(t *testing.T, feed stream.Feed, size int) {
	activities := make(stream.Activities, size)
	for i := range activities {
		var verb string
		if i%2 == 0 {
			verb = "even"
		} else {
			verb = "odd"
		}
		activities[i] = stream.Activity{Actor: "test", Verb: verb, Object: randString(10)}
	}
	_, err := feed.AddActivities(activities...)
	require.NoError(t, err)
}

func TestFlatFeedGetActivities(t *testing.T) {
	var (
		client = newClient(t)
		flat   = client.FlatFeed("flat", randString(10))
		size   = 15
	)
	prepareFeedForTestGetActivities(t, flat, size)
	resp, err := flat.GetActivities()
	require.NoError(t, err)
	assert.Len(t, resp.Results, size)

	limit := 2
	resp, err = flat.GetActivities(stream.GetActivitiesWithLimit(limit))
	require.NoError(t, err)
	assert.Len(t, resp.Results, limit)
}
