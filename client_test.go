package stream_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v5"
)

func TestHeaders(t *testing.T) {
	client, requester := newClient(t)
	feed, err := client.FlatFeed("user", "123")
	require.NoError(t, err)

	_, err = feed.GetActivities()
	require.NoError(t, err)
	assert.Equal(t, "application/json", requester.req.Header.Get("content-type"))
	assert.Regexp(t, "^stream-go2-client-v[0-9\\.]+$", requester.req.Header.Get("x-stream-client"))
}

func TestAddToMany(t *testing.T) {
	var (
		client, requester = newClient(t)
		activity          = stream.Activity{Actor: "bob", Verb: "like", Object: "cake"}
		flat, _           = newFlatFeedWithUserID(client, "123")
		aggregated, _     = newAggregatedFeedWithUserID(client, "123")
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
		flat, _           = newFlatFeedWithUserID(client, "123")
	)

	for i := range relationships {
		other, _ := newAggregatedFeedWithUserID(client, strconv.Itoa(i))
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
		flat, _           = newFlatFeedWithUserID(client, "123")
	)

	for i := range relationships {
		other, _ := newAggregatedFeedWithUserID(client, strconv.Itoa(i))
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
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/activities/?api_key=key&foreign_ids=foo%2Cbar&timestamps=0001-01-01T00%3A00%3A00%2C0001-01-01T00%3A00%3A01", "")
}

func TestGetEnrichedActivities(t *testing.T) {
	client, requester := newClient(t)
	_, err := client.GetEnrichedActivitiesByID("foo", "bar", "baz")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/enrich/activities/?api_key=key&ids=foo%2Cbar%2Cbaz", "")
	_, err = client.GetEnrichedActivitiesByForeignID(
		stream.NewForeignIDTimePair("foo", stream.Time{}),
		stream.NewForeignIDTimePair("bar", stream.Time{Time: time.Time{}.Add(time.Second)}),
	)
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/enrich/activities/?api_key=key&foreign_ids=foo%2Cbar&timestamps=0001-01-01T00%3A00%3A00%2C0001-01-01T00%3A00%3A01", "")
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

func TestPartialUpdateActivities(t *testing.T) {
	client, requester := newClient(t)

	_, err := client.PartialUpdateActivities(
		stream.NewUpdateActivityRequestByID(
			"abcdef",
			map[string]interface{}{"foo.bar": "baz"},
			[]string{"qux", "tty"},
		),
		stream.NewUpdateActivityRequestByID(
			"ghijkl",
			map[string]interface{}{"foo.bar": "baz"},
			[]string{"quux", "ttl"},
		),
	)
	require.NoError(t, err)
	body := `{"changes":[{"id":"abcdef","set":{"foo.bar":"baz"},"unset":["qux","tty"]},{"id":"ghijkl","set":{"foo.bar":"baz"},"unset":["quux","ttl"]}]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)

	tt, _ := time.Parse(stream.TimeLayout, "2006-01-02T15:04:05.999999")
	_, err = client.PartialUpdateActivities(
		stream.NewUpdateActivityRequestByForeignID(
			"abcdef:123",
			stream.Time{Time: tt},
			map[string]interface{}{"foo.bar": "baz"},
			nil,
		),
		stream.NewUpdateActivityRequestByForeignID(
			"ghijkl:1234",
			stream.Time{Time: tt},
			nil,
			[]string{"quux", "ttl"},
		),
	)
	require.NoError(t, err)
	body = `{"changes":[{"foreign_id":"abcdef:123","time":"2006-01-02T15:04:05.999999","set":{"foo.bar":"baz"}},{"foreign_id":"ghijkl:1234","time":"2006-01-02T15:04:05.999999","unset":["quux","ttl"]}]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activity/?api_key=key", body)
}

func TestUpdateActivityByForeignID(t *testing.T) {
	client, requester := newClient(t)

	tt := stream.Time{Time: time.Date(2018, 6, 24, 11, 28, 0, 0, time.UTC)}

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

func TestUserSessionToken(t *testing.T) {
	client, _ := newClient(t)
	tokenString, err := client.GetUserSessionToken("user")
	require.NoError(t, err)
	assert.Equal(t, tokenString, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlciJ9.0Kiui6HUywyU-C-00E68n1iq_3o7Eh0aE5VGSOc3pU4")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return []byte("secret"), nil })
	require.NoError(t, err)
	assert.Equal(t, true, token.Valid)
	assert.Equal(t, token.Claims, jwt.MapClaims{"user_id": "user"})
}

func TestUserSessionTokenWithClaims(t *testing.T) {
	client, _ := newClient(t)
	tokenString, err := client.GetUserSessionTokenWithClaims("user", map[string]interface{}{"client": "go"})
	require.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnQiOiJnbyIsInVzZXJfaWQiOiJ1c2VyIn0.Us6UIuH83dJe1cXQIiudseFz9-1kVMr6-SL6-idzIB0", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return []byte("secret"), nil })
	require.NoError(t, err)
	assert.Equal(t, true, token.Valid)
	assert.Equal(t, token.Claims, jwt.MapClaims{"user_id": "user", "client": "go"})
}
