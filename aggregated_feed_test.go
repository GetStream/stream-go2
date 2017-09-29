package stream_test

import (
	"net/http"
	"testing"

	"github.com/GetStream/stream-go2"

	"github.com/stretchr/testify/assert"
)

func TestAggregatedFeedGetActivities(t *testing.T) {
	client, requester := newClient(t)
	aggregated := newAggregatedFeedWithUserID(client, "123")
	testCases := []struct {
		opts []stream.GetActivitiesOption
		url  string
	}{
		{
			url: "https://api.getstream.io/api/v1.0/feed/aggregated/123/?api_key=key",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithActivitiesLimit(42)},
			url:  "https://api.getstream.io/api/v1.0/feed/aggregated/123/?api_key=key&limit=42",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithActivitiesLimit(42), stream.WithActivitiesOffset(11), stream.WithActivitiesIDGT("aabbcc")},
			url:  "https://api.getstream.io/api/v1.0/feed/aggregated/123/?api_key=key&limit=42&offset=11&id_gt=aabbcc",
		},
	}

	for _, tc := range testCases {
		_, err := aggregated.GetActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)
	}
}
