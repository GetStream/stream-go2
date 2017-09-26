package stream_test

import (
	"sort"
	"testing"

	"github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddToMany(t *testing.T) {
	var (
		client     = newClient(t)
		activity   = stream.Activity{Actor: "bob", Verb: "like", Object: "cake"}
		flat       = newFlatFeed(client)
		aggregated = newAggregatedFeed(client)
	)

	err := client.AddToMany(activity, flat, aggregated)
	require.Nil(t, err)

	flatActivities, err := flat.GetActivities()
	require.NoError(t, err)
	assert.Len(t, flatActivities.Results, 1)

	aggregatedActivities, err := aggregated.GetActivities()
	require.NoError(t, err)
	assert.Len(t, aggregatedActivities.Results, 1)
}

func TestFollowMany(t *testing.T) {
	var (
		client        = newClient(t)
		relationships = make([]stream.FollowRelationship, 10)
		flat          = client.FlatFeed("flat", randString(10))
	)

	for i := range relationships {
		other := client.AggregatedFeed("aggregated", randString(10))
		relationships[i] = stream.NewFollowRelationship(other, flat)
	}

	err := client.FollowMany(relationships)
	require.NoError(t, err)

	follows, err := flat.GetFollowers()
	require.NoError(t, err)
	require.Len(t, follows.Results, 10)

	expected := make([]string, 10)
	for i := range expected {
		expected[i] = relationships[i].Source
	}
	sort.Strings(expected)

	actual := make([]string, 10)
	for i := range actual {
		actual[i] = follows.Results[i].FeedID
	}
	sort.Strings(actual)

	assert.Equal(t, expected, actual)
}
