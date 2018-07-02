package stream_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	stream "github.com/GetStream/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	client, requester := newClient(t)
	_, err := client.FlatFeed("user", "123").GetActivities()
	require.NoError(t, err)
	assert.Equal(t, "application/json", requester.req.Header.Get("content-type"))
	assert.Regexp(t, "^stream-go2-client-v[0-9\\.]+$", requester.req.Header.Get("x-stream-client"))
}

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

	err = client.FollowMany(relationships, stream.WithFollowManyActivityCopyLimit(500))
	require.NoError(t, err)
	body = `[{"source":"aggregated:0","target":"flat:123"},{"source":"aggregated:1","target":"flat:123"},{"source":"aggregated:2","target":"flat:123"}]`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/follow_many/?activity_copy_limit=500&api_key=key", body)
}

func TestFollowManyActivityCopyLimit(t *testing.T) {
	var (
		client, requester = newClient(t)
		relationships     = make([]stream.FollowRelationship, 3)
		flat              = newFlatFeedWithUserID(client, "123")
	)

	for i := range relationships {
		other := newAggregatedFeedWithUserID(client, strconv.Itoa(i))
		relationships[i] = stream.NewFollowRelationship(other, flat, stream.WithFollowRelationshipActivityCopyLimit(i))
	}

	err := client.FollowMany(relationships)
	require.NoError(t, err)
	body := `[{"source":"aggregated:0","target":"flat:123","activity_copy_limit":0},{"source":"aggregated:1","target":"flat:123","activity_copy_limit":1},{"source":"aggregated:2","target":"flat:123","activity_copy_limit":2}]`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/follow_many/?api_key=key", body)

	err = client.FollowMany(relationships, stream.WithFollowManyActivityCopyLimit(123))
	require.NoError(t, err)
	body = `[{"source":"aggregated:0","target":"flat:123","activity_copy_limit":0},{"source":"aggregated:1","target":"flat:123","activity_copy_limit":1},{"source":"aggregated:2","target":"flat:123","activity_copy_limit":2}]`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/follow_many/?activity_copy_limit=123&api_key=key", body)
}

func TestUnfollowMany(t *testing.T) {
	var (
		client, requester = newClient(t)
		relationships     = make([]stream.UnfollowRelationship, 3)
	)
	for i := range relationships {
		relationships[i] = stream.UnfollowRelationship{
			Source:      fmt.Sprintf("src:%d", i),
			Target:      fmt.Sprintf("tgt:%d", i),
			KeepHistory: i%2 == 0,
		}
	}
	err := client.UnfollowMany(relationships)
	require.NoError(t, err)
	body := `[{"source":"src:0","target":"tgt:0","keep_history":true},{"source":"src:1","target":"tgt:1","keep_history":false},{"source":"src:2","target":"tgt:2","keep_history":true}]`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/unfollow_many/?api_key=key", body)
}

func TestGetActivities(t *testing.T) {
	client, requester := newClient(t)
	_, err := client.GetActivitiesByID("foo", "bar", "baz")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/activities/?api_key=key&ids=foo%2Cbar%2Cbaz", "")
	_, err = client.GetActivitiesByForeignID(
		stream.NewForeignIDTimePair("foo", stream.Time{}),
		stream.NewForeignIDTimePair("bar", stream.Time{Time: time.Time{}.Add(time.Second)}),
	)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/activities/?api_key=key&foreign_ids=foo%2Cbar&timestamps=0001-01-01T00%3A00%3A00%2C0001-01-01T00%3A00%3A01", "")
}

func TestUpdateActivityByID(t *testing.T) {
	client, requester := newClient(t)

	_, err := client.UpdateActivityByID("abcdef", map[string]interface{}{"foo.bar": "baz", "popularity": 42, "color": map[string]interface{}{"hex": "FF0000", "rgb": "255,0,0"}}, []string{"a", "b", "c"})
	require.NoError(t, err)
	body := `{"id":"abcdef","set":{"color":{"hex":"FF0000","rgb":"255,0,0"},"foo.bar":"baz","popularity":42},"unset":["a","b","c"]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)

	_, err = client.UpdateActivityByID("abcdef", map[string]interface{}{"foo.bar": "baz", "popularity": 42, "color": map[string]interface{}{"hex": "FF0000", "rgb": "255,0,0"}}, nil)
	require.NoError(t, err)
	body = `{"id":"abcdef","set":{"color":{"hex":"FF0000","rgb":"255,0,0"},"foo.bar":"baz","popularity":42}}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)

	_, err = client.UpdateActivityByID("abcdef", nil, []string{"a", "b", "c"})
	require.NoError(t, err)
	body = `{"id":"abcdef","unset":["a","b","c"]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)
}

func TestUpdateActivityByForeignID(t *testing.T) {
	client, requester := newClient(t)

	tt := stream.Time{Time: time.Date(2018, 06, 24, 11, 28, 0, 0, time.UTC)}

	_, err := client.UpdateActivityByForeignID("fid:123", tt, map[string]interface{}{"foo.bar": "baz", "popularity": 42, "color": map[string]interface{}{"hex": "FF0000", "rgb": "255,0,0"}}, []string{"a", "b", "c"})
	require.NoError(t, err)
	body := `{"foreign_id":"fid:123","time":"2018-06-24T11:28:00","set":{"color":{"hex":"FF0000","rgb":"255,0,0"},"foo.bar":"baz","popularity":42},"unset":["a","b","c"]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)

	_, err = client.UpdateActivityByForeignID("fid:123", tt, map[string]interface{}{"foo.bar": "baz", "popularity": 42, "color": map[string]interface{}{"hex": "FF0000", "rgb": "255,0,0"}}, nil)
	require.NoError(t, err)
	body = `{"foreign_id":"fid:123","time":"2018-06-24T11:28:00","set":{"color":{"hex":"FF0000","rgb":"255,0,0"},"foo.bar":"baz","popularity":42}}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)

	_, err = client.UpdateActivityByForeignID("fid:123", tt, nil, []string{"a", "b", "c"})
	require.NoError(t, err)
	body = `{"foreign_id":"fid:123","time":"2018-06-24T11:28:00","unset":["a","b","c"]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)
}
