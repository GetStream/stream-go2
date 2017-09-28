package stream_test

import (
	"net/http"
	"testing"

	stream "github.com/reifcode/stream-go2"
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
			url: "https://api.getstream.io/api/v1.0/feed/notification/123/?api_key=key",
		},
		{
			opts: []stream.GetActivitiesOption{stream.GetActivitiesWithLimit(42)},
			url:  "https://api.getstream.io/api/v1.0/feed/notification/123/?api_key=key&limit=42",
		},
		{
			opts: []stream.GetActivitiesOption{stream.GetActivitiesWithLimit(42), stream.GetActivitiesWithOffset(11), stream.GetActivitiesWithIDGT("aabbcc")},
			url:  "https://api.getstream.io/api/v1.0/feed/notification/123/?api_key=key&limit=42&offset=11&id_gt=aabbcc",
		},
		{
			opts: []stream.GetActivitiesOption{stream.GetNotificationWithMarkRead(false, "f1", "f2", "f3")},
			url:  "https://api.getstream.io/api/v1.0/feed/notification/123/?api_key=key&mark_read=f1,f2,f3",
		},
		{
			opts: []stream.GetActivitiesOption{stream.GetNotificationWithMarkRead(true, "f1", "f2", "f3")},
			url:  "https://api.getstream.io/api/v1.0/feed/notification/123/?api_key=key&mark_read=true",
		},
		{
			opts: []stream.GetActivitiesOption{stream.GetNotificationWithMarkSeen(false, "f1", "f2", "f3")},
			url:  "https://api.getstream.io/api/v1.0/feed/notification/123/?api_key=key&mark_seen=f1,f2,f3",
		},
		{
			opts: []stream.GetActivitiesOption{stream.GetNotificationWithMarkSeen(true, "f1", "f2", "f3")},
			url:  "https://api.getstream.io/api/v1.0/feed/notification/123/?api_key=key&mark_seen=true",
		},
	}

	for _, tc := range testCases {
		_, err := notification.GetActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)
	}
}
