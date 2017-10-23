package stream_test

import (
	"fmt"
	"net/http"
	"testing"

	stream "github.com/GetStream/stream-go2"
	"github.com/stretchr/testify/assert"
)

func TestFlatFeedGetActivities(t *testing.T) {
	client, requester := newClient(t)
	flat := newFlatFeedWithUserID(client, "123")
	testCases := []struct {
		opts []stream.GetActivitiesOption
		url  string
	}{
		{
			url: "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key",
		},
		{
			opts: []stream.GetActivitiesOption{stream.WithActivitiesLimit(42)},
			url:  "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&limit=42",
		},
		{
			opts: []stream.GetActivitiesOption{
				stream.WithActivitiesLimit(42),
				stream.WithActivitiesOffset(11),
				stream.WithActivitiesIDGT("aabbcc"),
				stream.WithActivitiesIDGTE("ccddee"),
				stream.WithActivitiesIDLT("ffgghh"),
				stream.WithActivitiesIDLTE("iijjkk"),
			},
			url: "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&limit=42&offset=11&id_gt=aabbcc&id_gte=ccddee&id_lt=ffgghh&id_lte=iijjkk",
		},
	}

	for _, tc := range testCases {
		_, err := flat.GetActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)

		_, err = flat.GetActivitiesWithRanking("popularity", tc.opts...)
		testRequest(t, requester.req, http.MethodGet, fmt.Sprintf("%s&ranking=popularity", tc.url), "")
		assert.NoError(t, err)
	}
}
