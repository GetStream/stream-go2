package stream_test

import (
	"testing"

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
	assert.Len(t, resp.Results[0].Activities, size/2)
	assert.Len(t, resp.Results[1].Activities, size/2)
	// TODO test read options
}
