package stream_test

import (
	"testing"

	"github.com/reifcode/stream-go2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregatedFeedGetActivities(t *testing.T) {
	var (
		client     = newClient(t)
		aggregated = client.AggregatedFeed("aggregated", randString(10))
		size       = 10
	)
	prepareFeedForTestGetActivities(t, aggregated, size)
	resp, err := aggregated.GetActivities()
	require.NoError(t, err)
	assert.Len(t, resp.Results, 2)

	resp, err = aggregated.GetActivities(stream.GetActivitiesWithLimit(1))
	require.NoError(t, err)
	assert.Len(t, resp.Results, 1)

	_, err = aggregated.AddActivities(
		stream.Activity{
			Actor:  "test",
			Verb:   randString(10),
			Object: randString(10)},
	)
	require.NoError(t, err)

	resp, err = aggregated.GetActivities(stream.GetActivitiesWithOffset(0))
	require.NoError(t, err)
	assert.Len(t, resp.Results, 3)

	resp, err = aggregated.GetActivities(
		stream.GetActivitiesWithOffset(0),
		stream.GetActivitiesWithLimit(2),
	)
	require.NoError(t, err)
	assert.Len(t, resp.Results, 2)
}
