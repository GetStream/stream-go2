package stream_test

import (
	"net/http"
	"testing"

	stream "github.com/GetStream/stream-go2"
	"github.com/stretchr/testify/assert"
)

func TestGetNotificationActivities(t *testing.T) {
	client, requester := newClient(t)
	notification := newNotificationFeedWithUserID(client, "123")
	testCases := []struct {
		opts []stream.GetActivitiesOption
		url  string
	}{
		{
			url: "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithActivitiesLimit(42)},
			url:  "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&limit=42",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithActivitiesLimit(42), stream.WithActivitiesOffset(11), stream.WithActivitiesIDGT("aabbcc")},
			url:  "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&id_gt=aabbcc&limit=42&offset=11",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithNotificationsMarkRead(false, "f1", "f2", "f3")},
			url:  "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_read=f1%2Cf2%2Cf3",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithNotificationsMarkRead(true, "f1", "f2", "f3")},
			url:  "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_read=true",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithNotificationsMarkSeen(false, "f1", "f2", "f3")},
			url:  "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_seen=f1%2Cf2%2Cf3",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithNotificationsMarkSeen(true, "f1", "f2", "f3")},
			url:  "https://api.stream-io-api.com/api/v1.0/feed/notification/123/?api_key=key&mark_seen=true",
		},
	}

	for _, tc := range testCases {
		_, err := notification.GetActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)
	}
}
