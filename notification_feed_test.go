package stream_test

import (
	"net/http"
	"testing"

	stream "github.com/GetStream/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNotificationActivities(t *testing.T) {
	client, requester := newClient(t)
	notification := newNotificationFeedWithUserID(client, "123")
	testCases := []struct {
		opts        []stream.GetActivitiesOption
		url         string
		enrichedURL string
	}{
		{
			url:         "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithActivitiesLimit(42)},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&limit=42",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key&limit=42",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithActivitiesLimit(42), stream.WithActivitiesOffset(11), stream.WithActivitiesIDGT("aabbcc")},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&id_gt=aabbcc&limit=42&offset=11",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key&id_gt=aabbcc&limit=42&offset=11",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithNotificationsMarkRead(false, "f1", "f2", "f3")},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_read=f1%2Cf2%2Cf3",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key&mark_read=f1%2Cf2%2Cf3",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithNotificationsMarkRead(true, "f1", "f2", "f3")},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_read=true",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key&mark_read=true",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithNotificationsMarkSeen(false, "f1", "f2", "f3")},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_seen=f1%2Cf2%2Cf3",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key&mark_seen=f1%2Cf2%2Cf3",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithNotificationsMarkSeen(true, "f1", "f2", "f3")},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_seen=true",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key&mark_seen=true",
		},
	}

	for _, tc := range testCases {
		_, err := notification.GetActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)

		_, err = notification.GetEnrichedActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.enrichedURL, "")
		assert.NoError(t, err)
	}
}

func TestNotificationFeedGetNextPageActivities(t *testing.T) {
	client, requester := newClient(t)
	notification := newNotificationFeedWithUserID(client, "123")

	requester.resp = `{"next":"/api/v1.0/feed/notification/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	resp, err := notification.GetActivities()
	require.NoError(t, err)

	_, err = notification.GetNextPageActivities(resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/enrich/feed/notification/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	enrichedResp, err := notification.GetEnrichedActivities()
	require.NoError(t, err)

	_, err = notification.GetNextPageEnrichedActivities(enrichedResp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/enrich/feed/notification/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":123}`
	_, err = notification.GetActivities()
	require.Error(t, err)

	requester.resp = `{"next":"123"}`
	resp, err = notification.GetActivities()
	require.NoError(t, err)
	_, err = notification.GetNextPageActivities(resp)
	require.Error(t, err)

	requester.resp = `{"next":"?q=a%"}`
	resp, err = notification.GetActivities()
	require.NoError(t, err)
	_, err = notification.GetNextPageActivities(resp)
	require.Error(t, err)
}
