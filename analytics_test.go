package stream_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v8"
)

func TestAnalyticsTrackEngagement(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	analytics := client.Analytics()
	event1 := stream.EngagementEvent{}.
		WithLabel("click").
		WithForeignID("abcdef").
		WithUserData(stream.NewUserData().Int(12345).Alias("John Doe")).
		WithFeedID("timeline:123").
		WithLocation("hawaii").
		WithPosition(42).
		WithBoost(10)

	event2 := stream.EngagementEvent{}.
		WithLabel("share").
		WithForeignID("aabbccdd").
		WithUserData(stream.NewUserData().String("bob")).
		WithFeedID("timeline:123").
		WithFeatures(
			stream.NewEventFeature("color", "red"),
			stream.NewEventFeature("size", "xxl"),
		)

	_, err := analytics.TrackEngagement(ctx, event1, event2)
	require.NoError(t, err)
	expectedURL := "https://analytics.stream-io-api.com/analytics/v1.0/engagement/?api_key=key"
	expectedBody := `{"content_list":[{"boost":10,"content":"abcdef","feed_id":"timeline:123","label":"click","location":"hawaii","position":42,"user_data":{"alias":"John Doe","id":12345}},{"content":"aabbccdd","features":[{"group":"color","value":"red"},{"group":"size","value":"xxl"}],"feed_id":"timeline:123","label":"share","user_data":"bob"}]}`
	testRequest(t, requester.req, http.MethodPost, expectedURL, expectedBody)
}

func TestAnalyticsTrackImpression(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	analytics := client.Analytics()
	imp := stream.ImpressionEventsData{}.
		WithForeignIDs("a", "b", "c", "d").
		WithUserData(stream.NewUserData().Int(123)).
		WithFeedID("timeline:123").
		WithFeatures(
			stream.NewEventFeature("color", "red"),
			stream.NewEventFeature("size", "xxl"),
		).
		WithLocation("hawaii").
		WithPosition(42)

	_, err := analytics.TrackImpression(ctx, imp)
	require.NoError(t, err)
	expectedURL := "https://analytics.stream-io-api.com/analytics/v1.0/impression/?api_key=key"
	expectedBody := `{"content_list":["a","b","c","d"],"features":[{"group":"color","value":"red"},{"group":"size","value":"xxl"}],"feed_id":"timeline:123","location":"hawaii","position":42,"user_data":123}`
	testRequest(t, requester.req, http.MethodPost, expectedURL, expectedBody)
}

func TestAnalyticsRedirectAndTrack(t *testing.T) {
	client, _ := newClient(t)
	analytics := client.Analytics()
	event1 := stream.EngagementEvent{}.
		WithLabel("click").
		WithForeignID("abcdef").
		WithUserData(stream.NewUserData().Int(12345).Alias("John Doe")).
		WithFeedID("timeline:123").
		WithLocation("hawaii").
		WithPosition(42).
		WithBoost(10)
	event2 := stream.EngagementEvent{}.
		WithLabel("share").
		WithForeignID("aabbccdd").
		WithUserData(stream.NewUserData().String("bob")).
		WithFeedID("timeline:123").
		WithFeatures(
			stream.NewEventFeature("color", "red"),
			stream.NewEventFeature("size", "xxl"),
		)
	imp := stream.ImpressionEventsData{}.
		WithForeignIDs("a", "b", "c", "d").
		WithUserData(stream.NewUserData().Int(123)).
		WithFeedID("timeline:123").
		WithFeatures(
			stream.NewEventFeature("color", "red"),
			stream.NewEventFeature("size", "xxl"),
		).
		WithLocation("hawaii").
		WithPosition(42)

	link, err := analytics.RedirectAndTrack("foo.bar.baz", event1, event2, imp)
	require.NoError(t, err)
	query, err := url.ParseQuery(`api_key=key&authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb24iOiIqIiwicmVzb3VyY2UiOiJyZWRpcmVjdF9hbmRfdHJhY2siLCJ1c2VyX2lkIjoiKiJ9.A1vy9pFwLw5s6qn0chkhRcoy974A16a0lE-x5Vtxb-o&events=[{"boost":10,"content":"abcdef","feed_id":"timeline:123","label":"click","location":"hawaii","position":42,"user_data":{"alias":"John Doe","id":12345}},{"content":"aabbccdd","features":[{"group":"color","value":"red"},{"group":"size","value":"xxl"}],"feed_id":"timeline:123","label":"share","user_data":"bob"},{"content_list":["a","b","c","d"],"features":[{"group":"color","value":"red"},{"group":"size","value":"xxl"}],"feed_id":"timeline:123","location":"hawaii","position":42,"user_data":123}]&stream-auth-type=jwt&url=foo.bar.baz`)
	require.NoError(t, err)
	expected := "https://analytics.stream-io-api.com/analytics/v1.0/redirect/?" + query.Encode()
	assert.Equal(t, expected, link)
}
