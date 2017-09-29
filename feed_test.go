package stream_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	stream "github.com/GetStream/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedID(t *testing.T) {
	client, _ := newClient(t)
	flat := client.FlatFeed("flat", "123")
	assert.Equal(t, "flat:123", flat.ID())
	aggregated := client.AggregatedFeed("aggregated", "456")
	assert.Equal(t, "aggregated:456", aggregated.ID())
}

func TestAddActivity(t *testing.T) {
	var (
		client, requester = newClient(t)
		flat              = newFlatFeedWithUserID(client, "123")
		bobActivity       = stream.Activity{Actor: "bob", Verb: "like", Object: "ice-cream"}
	)
	_, err := flat.AddActivity(bobActivity)
	require.NoError(t, err)
	body := `{"actor":"bob","object":"ice-cream","verb":"like"}`
	testRequest(t, requester.req, http.MethodPost, "https://api.getstream.io/api/v1.0/feed/flat/123/?api_key=key", body)

	requester.resp = `{"duration": "something-broken"}`
	_, err = flat.AddActivity(bobActivity)
	require.Error(t, err)

	requester.resp = `{"duration": "1ms"}`
	_, err = flat.AddActivity(bobActivity)
	require.NoError(t, err)
}

func TestAddActivities(t *testing.T) {
	var (
		client, requester = newClient(t)
		flat              = newFlatFeedWithUserID(client, "123")
		bobActivity       = stream.Activity{Actor: "bob", Verb: "like", Object: "ice-cream"}
		aliceActivity     = stream.Activity{Actor: "alice", Verb: "dislike", Object: "ice-cream"}
	)
	_, err := flat.AddActivities(bobActivity, aliceActivity)
	require.NoError(t, err)
	body := `{"activities":[{"actor":"bob","object":"ice-cream","verb":"like"},{"actor":"alice","object":"ice-cream","verb":"dislike"}]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.getstream.io/api/v1.0/feed/flat/123/?api_key=key", body)
}

func TestUpdateActivities(t *testing.T) {
	var (
		client, requester = newClient(t)
		flat              = newFlatFeedWithUserID(client, "123")
		now               = getTime(time.Now())
		bobActivity       = stream.Activity{
			Actor:     "bob",
			Verb:      "like",
			Object:    "ice-cream",
			ForeignID: "bob:123",
			Time:      now,
			Extra:     map[string]interface{}{"influence": 42},
		}
	)
	err := flat.UpdateActivities(bobActivity)
	require.NoError(t, err)

	body := fmt.Sprintf(`{"activities":[{"actor":"bob","foreign_id":"bob:123","influence":42,"object":"ice-cream","time":"%s","verb":"like"}]}`, now.Format(stream.TimeLayout))
	testRequest(t, requester.req, http.MethodPost, "https://api.getstream.io/api/v1.0/activities/?api_key=key", body)
}

func TestFollow(t *testing.T) {
	var (
		client, requester = newClient(t)
		f1, f2            = newFlatFeedWithUserID(client, "f1"), newFlatFeedWithUserID(client, "f2")
	)
	testCases := []struct {
		opts         []stream.FollowFeedOption
		expectedURL  string
		expectedBody string
	}{
		{
			expectedURL:  "https://api.getstream.io/api/v1.0/feed/flat/f1/follows/?api_key=key",
			expectedBody: `{"target":"flat:f2","activity_copy_limit":300}`,
		},
		{
			opts:         []stream.FollowFeedOption{stream.WithFollowFeedActivityCopyLimit(123)},
			expectedURL:  "https://api.getstream.io/api/v1.0/feed/flat/f1/follows/?api_key=key",
			expectedBody: `{"target":"flat:f2","activity_copy_limit":123}`,
		},
	}
	for _, tc := range testCases {
		err := f1.Follow(f2, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestGetFollowing(t *testing.T) {
	var (
		client, requester = newClient(t)
		f1                = newFlatFeedWithUserID(client, "f1")
	)
	testCases := []struct {
		opts     []stream.FollowingOption
		expected string
	}{
		{
			expected: "https://api.getstream.io/api/v1.0/feed/flat/f1/follows/?api_key=key",
		},
		{
			opts:     []stream.FollowingOption{stream.WithFollowingFilter("filter"), stream.WithFollowingLimit(42), stream.WithFollowingOffset(84)},
			expected: "https://api.getstream.io/api/v1.0/feed/flat/f1/follows/?api_key=key&filter=filter&limit=42&offset=84",
		},
	}
	for _, tc := range testCases {
		_, err := f1.GetFollowing(tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expected, "")
	}
}

func TestGetFollowers(t *testing.T) {
	var (
		client, requester = newClient(t)
		f1                = newFlatFeedWithUserID(client, "f1")
	)
	testCases := []struct {
		opts     []stream.FollowersOption
		expected string
	}{
		{
			expected: "https://api.getstream.io/api/v1.0/feed/flat/f1/followers/?api_key=key",
		},
		{
			opts:     []stream.FollowersOption{stream.WithFollowersLimit(42), stream.WithFollowersOffset(84)},
			expected: "https://api.getstream.io/api/v1.0/feed/flat/f1/followers/?api_key=key&limit=42&offset=84",
		},
	}
	for _, tc := range testCases {
		_, err := f1.GetFollowers(tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expected, "")
	}
}

func TestUnfollow(t *testing.T) {
	var (
		client, requester = newClient(t)
		f1, f2            = newFlatFeedWithUserID(client, "f1"), newFlatFeedWithUserID(client, "f2")
	)
	testCases := []struct {
		opts     []stream.UnfollowOption
		expected string
	}{
		{
			expected: "https://api.getstream.io/api/v1.0/feed/flat/f1/follows/flat:f2/?api_key=key",
		},
		{
			opts:     []stream.UnfollowOption{stream.WithKeepHistory(false)},
			expected: "https://api.getstream.io/api/v1.0/feed/flat/f1/follows/flat:f2/?api_key=key",
		},
		{
			opts:     []stream.UnfollowOption{stream.WithKeepHistory(true)},
			expected: "https://api.getstream.io/api/v1.0/feed/flat/f1/follows/flat:f2/?api_key=key&keep_history=1",
		},
	}

	for _, tc := range testCases {
		err := f1.Unfollow(f2, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodDelete, tc.expected, "")
	}
}

func TestRemoveActivities(t *testing.T) {
	client, requester := newClient(t)
	flat := newFlatFeedWithUserID(client, "123")
	err := flat.RemoveActivityByID("id-to-remove")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.getstream.io/api/v1.0/feed/flat/123/id-to-remove/?api_key=key", "")
	err = flat.RemoveActivityByForeignID("bob:123")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.getstream.io/api/v1.0/feed/flat/123/bob:123/?api_key=key&foreign_id=1", "")
}

func TestUpdateToTargets(t *testing.T) {
	var (
		client, requester = newClient(t)
		flat              = newFlatFeedWithUserID(client, "123")
		f1, f2, f3        = newFlatFeedWithUserID(client, "f1"), newFlatFeedWithUserID(client, "f2"), newFlatFeedWithUserID(client, "f3")
		now               = getTime(time.Now())
		activity          = stream.Activity{Time: now, ForeignID: "bob:123", Actor: "bob", Verb: "like", Object: "ice-cream", To: []string{f1.ID()}, Extra: map[string]interface{}{"popularity": 9000}}
	)
	err := flat.UpdateToTargets(activity, stream.WithAddToTargets(f2.ID()), stream.WithRemoveToTargets(f1.ID()))
	require.NoError(t, err)
	body := fmt.Sprintf(`{"foreign_id":"bob:123","time":"%s","added_targets":["flat:f2"],"removed_targets":["flat:f1"]}`, now.Format(stream.TimeLayout))
	testRequest(t, requester.req, http.MethodPost, "https://api.getstream.io/api/v1.0/feed_targets/flat/123/activity_to_targets/?api_key=key", body)
	err = flat.UpdateToTargets(activity, stream.WithNewToTargets(f3.ID()))
	require.NoError(t, err)
	body = fmt.Sprintf(`{"foreign_id":"bob:123","time":"%s","new_targets":["flat:f3"]}`, now.Format(stream.TimeLayout))
	testRequest(t, requester.req, http.MethodPost, "https://api.getstream.io/api/v1.0/feed_targets/flat/123/activity_to_targets/?api_key=key", body)
}

func TestToken(t *testing.T) {
	client, err := stream.NewClient("key", "super secret")
	require.NoError(t, err)
	flat := newFlatFeedWithUserID(client, "sample")
	testCases := []struct {
		readOnly bool
		expected string
	}{
		{
			readOnly: false,
			expected: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb24iOiJ3cml0ZSIsImZlZWRfaWQiOiJmbGF0c2FtcGxlIiwicmVzb3VyY2UiOiJmZWVkIn0._7eLZ3-_6dmOoCKp8MvSoKCp0PA-gAerKnr8tuwut2M",
		},
		{
			readOnly: true,
			expected: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb24iOiJyZWFkIiwiZmVlZF9pZCI6ImZsYXRzYW1wbGUiLCJyZXNvdXJjZSI6ImZlZWQifQ.Ab6NX3dAGbBiXkQrEIWg9Z-WRm1R4710ont2y0OONiE",
		},
	}
	for _, tc := range testCases {
		token := flat.Token(tc.readOnly)
		assert.Equal(t, tc.expected, token)
	}
}
