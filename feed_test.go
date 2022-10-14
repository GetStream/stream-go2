package stream_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v7"
)

func TestFeedID(t *testing.T) {
	client, _ := newClient(t)
	flat, _ := client.FlatFeed("flat", "123")
	assert.Equal(t, "flat:123", flat.ID())
	aggregated, _ := client.AggregatedFeed("aggregated", "456")
	assert.Equal(t, "aggregated:456", aggregated.ID())
}

func TestInvalidFeedUserID(t *testing.T) {
	client, _ := newClient(t)

	_, err := client.FlatFeed("flat", "jones:134")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid userID provided", err.Error())

	_, err = client.AggregatedFeed("aggregated", "jones,kim")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid userID provided", err.Error())
}

func TestAddActivity(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		flat, _           = newFlatFeedWithUserID(client, "123")
		bobActivity       = stream.Activity{Actor: "bob", Verb: "like", Object: "ice-cream", To: []string{"flat:456"}}
	)
	_, err := flat.AddActivity(ctx, bobActivity)
	require.NoError(t, err)
	body := `{"actor":"bob","object":"ice-cream","to":["flat:456"],"verb":"like"}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key", body)

	requester.resp = `{"duration": "1ms"}`
	_, err = flat.AddActivity(ctx, bobActivity)
	require.NoError(t, err)
}

func TestAddActivities(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		flat, _           = newFlatFeedWithUserID(client, "123")
		bobActivity       = stream.Activity{Actor: "bob", Verb: "like", Object: "ice-cream"}
		aliceActivity     = stream.Activity{Actor: "alice", Verb: "dislike", Object: "ice-cream", To: []string{"flat:456"}}
	)
	_, err := flat.AddActivities(ctx, bobActivity, aliceActivity)
	require.NoError(t, err)
	body := `{"activities":[{"actor":"bob","object":"ice-cream","verb":"like"},{"actor":"alice","object":"ice-cream","to":["flat:456"],"verb":"dislike"}]}`
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key", body)
}

func TestUpdateActivities(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		now               = getTime(time.Now())
		bobActivity       = stream.Activity{
			Actor:     "bob",
			Verb:      "like",
			Object:    "ice-cream",
			ForeignID: "bob:123",
			Time:      now,
			Extra:     map[string]any{"influence": 42},
		}
	)
	_, err := client.UpdateActivities(ctx, bobActivity)
	require.NoError(t, err)

	body := fmt.Sprintf(`{"activities":[{"actor":"bob","foreign_id":"bob:123","influence":42,"object":"ice-cream","time":%q,"verb":"like"}]}`, now.Format(stream.TimeLayout))
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/activities/?api_key=key", body)
}

func TestFollow(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		f1, _             = newFlatFeedWithUserID(client, "f1")
		f2, _             = newFlatFeedWithUserID(client, "f2")
	)
	testCases := []struct {
		opts         []stream.FollowFeedOption
		expectedURL  string
		expectedBody string
	}{
		{
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/follows/?api_key=key",
			expectedBody: `{"target":"flat:f2","activity_copy_limit":300}`,
		},
		{
			opts:         []stream.FollowFeedOption{stream.WithFollowFeedActivityCopyLimit(123)},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/follows/?api_key=key",
			expectedBody: `{"target":"flat:f2","activity_copy_limit":123}`,
		},
	}
	for _, tc := range testCases {
		_, err := f1.Follow(ctx, f2, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestGetFollowing(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		f1, _             = newFlatFeedWithUserID(client, "f1")
	)
	testCases := []struct {
		opts     []stream.FollowingOption
		expected string
	}{
		{
			expected: "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/follows/?api_key=key",
		},
		{
			opts:     []stream.FollowingOption{stream.WithFollowingFilter("filter"), stream.WithFollowingLimit(42), stream.WithFollowingOffset(84)},
			expected: "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/follows/?api_key=key&filter=filter&limit=42&offset=84",
		},
	}
	for _, tc := range testCases {
		_, err := f1.GetFollowing(ctx, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expected, "")
	}
}

func TestGetFollowers(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		f1, _             = newFlatFeedWithUserID(client, "f1")
	)
	testCases := []struct {
		opts     []stream.FollowersOption
		expected string
	}{
		{
			expected: "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/followers/?api_key=key",
		},
		{
			opts:     []stream.FollowersOption{stream.WithFollowersLimit(42), stream.WithFollowersOffset(84)},
			expected: "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/followers/?api_key=key&limit=42&offset=84",
		},
	}
	for _, tc := range testCases {
		_, err := f1.GetFollowers(ctx, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expected, "")
	}
}

func TestUnfollow(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		f1, _             = newFlatFeedWithUserID(client, "f1")
		f2, _             = newFlatFeedWithUserID(client, "f2")
	)
	testCases := []struct {
		opts     []stream.UnfollowOption
		expected string
	}{
		{
			expected: "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/follows/flat:f2/?api_key=key",
		},
		{
			opts:     []stream.UnfollowOption{stream.WithUnfollowKeepHistory(false)},
			expected: "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/follows/flat:f2/?api_key=key",
		},
		{
			opts:     []stream.UnfollowOption{stream.WithUnfollowKeepHistory(true)},
			expected: "https://api.stream-io-api.com/api/v1.0/feed/flat/f1/follows/flat:f2/?api_key=key&keep_history=1",
		},
	}

	for _, tc := range testCases {
		_, err := f1.Unfollow(ctx, f2, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodDelete, tc.expected, "")
	}
}

func TestRemoveActivities(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	flat, _ := newFlatFeedWithUserID(client, "123")
	_, err := flat.RemoveActivityByID(ctx, "id-to-remove")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.stream-io-api.com/api/v1.0/feed/flat/123/id-to-remove/?api_key=key", "")
	_, err = flat.RemoveActivityByForeignID(ctx, "bob:123")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.stream-io-api.com/api/v1.0/feed/flat/123/bob:123/?api_key=key&foreign_id=1", "")
}

func TestUpdateToTargets(t *testing.T) {
	var (
		client, requester = newClient(t)
		ctx               = context.Background()
		flat, _           = newFlatFeedWithUserID(client, "123")
		f1, _             = newFlatFeedWithUserID(client, "f1")
		f2, _             = newFlatFeedWithUserID(client, "f2")
		f3, _             = newFlatFeedWithUserID(client, "f3")
		now               = getTime(time.Now())
		activity          = stream.Activity{Time: now, ForeignID: "bob:123", Actor: "bob", Verb: "like", Object: "ice-cream", To: []string{f1.ID()}, Extra: map[string]any{"popularity": 9000}}
	)
	_, err := flat.UpdateToTargets(ctx, activity, stream.WithToTargetsAdd(f2.ID()), stream.WithToTargetsRemove(f1.ID()))
	require.NoError(t, err)
	body := fmt.Sprintf(`{"foreign_id":"bob:123","time":%q,"added_targets":["flat:f2"],"removed_targets":["flat:f1"]}`, now.Format(stream.TimeLayout))
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/feed_targets/flat/123/activity_to_targets/?api_key=key", body)
	_, err = flat.UpdateToTargets(ctx, activity, stream.WithToTargetsNew(f3.ID()))
	require.NoError(t, err)
	body = fmt.Sprintf(`{"foreign_id":"bob:123","time":%q,"new_targets":["flat:f3"]}`, now.Format(stream.TimeLayout))
	testRequest(t, requester.req, http.MethodPost, "https://api.stream-io-api.com/api/v1.0/feed_targets/flat/123/activity_to_targets/?api_key=key", body)
}

func TestRealtimeToken(t *testing.T) {
	client, err := stream.New("key", "super secret")
	require.NoError(t, err)
	flat, _ := newFlatFeedWithUserID(client, "sample")
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
		token := flat.RealtimeToken(tc.readOnly)
		assert.Equal(t, tc.expected, token)
	}
}
