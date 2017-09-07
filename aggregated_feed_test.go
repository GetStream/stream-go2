package stream_test

import (
	"testing"

	stream "github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetActivities(t *testing.T) {
	client := newClient(t)
	aggregated := client.AggregatedFeed("aggregated", randString(10))

	activities := make(stream.Activities, 10)
	for i := range activities {
		var verb string
		if i%2 == 0 {
			verb = "even"
		} else {
			verb = "odd"
		}
		activities[i] = stream.Activity{Actor: "test", Verb: verb, Object: randString(10)}
	}
	_, err := aggregated.AddActivities(activities...)
	require.NoError(t, err)

	resp, err := aggregated.GetActivities()
	require.NoError(t, err)
	assert.Len(t, resp.Results, 2)

	// TODO test read options
}
