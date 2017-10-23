package stream_test

import (
	"net/http"
	"strconv"
	"testing"

	stream "github.com/GetStream/stream-go2"
	"github.com/stretchr/testify/require"
)

func TestAddToMany(t *testing.T) {
	var (
		client, requester = newClient(t)
		activity          = stream.Activity{Actor: "bob", Verb: "like", Object: "cake"}
		flat              = newFlatFeedWithUserID(client, "123")
		aggregated        = newAggregatedFeedWithUserID(client, "123")
	)

	err := client.AddToMany(activity, flat, aggregated)
	require.NoError(t, err)
	body := `{"activity":{"actor":"bob","object":"cake","verb":"like"},"feeds":["flat:123","aggregated:123"]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/feed/add_to_many/?api_key=key", body)
}

func TestFollowMany(t *testing.T) {
	var (
		client, requester = newClient(t)
		relationships     = make([]stream.FollowRelationship, 3)
		flat              = newFlatFeedWithUserID(client, "123")
	)

	for i := range relationships {
		other := newAggregatedFeedWithUserID(client, strconv.Itoa(i))
		relationships[i] = stream.NewFollowRelationship(other, flat)
	}

	err := client.FollowMany(relationships)
	require.NoError(t, err)
	body := `[{"source":"aggregated:0","target":"flat:123"},{"source":"aggregated:1","target":"flat:123"},{"source":"aggregated:2","target":"flat:123"}]`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/follow_many/?api_key=key", body)
}
