package stream_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v4"
)

func TestFlatFeedGetActivities(t *testing.T) {
	client, requester := newClient(t)
	flat, _ := newFlatFeedWithUserID(client, "123")
	testCases := []struct {
		opts        []stream.GetActivitiesOption
		url         string
		enrichedURL string
	}{
		{
			url:         "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?api_key=key",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithActivitiesLimit(42)},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&limit=42",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?api_key=key&limit=42",
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
			url:         "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&id_gt=aabbcc&id_gte=ccddee&id_lt=ffgghh&id_lte=iijjkk&limit=42&offset=11",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?api_key=key&id_gt=aabbcc&id_gte=ccddee&id_lt=ffgghh&id_lte=iijjkk&limit=42&offset=11",
		},
		{
			opts: []stream.GetActivitiesOption{
				stream.WithCustomParam("aaa", "bbb"),
			},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?aaa=bbb&api_key=key",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?aaa=bbb&api_key=key",
		},
	}

	for _, tc := range testCases {
		_, err := flat.GetActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)

		_, err = flat.GetActivitiesWithRanking("popularity", tc.opts...)
		testRequest(t, requester.req, http.MethodGet, fmt.Sprintf("%s&ranking=popularity", tc.url), "")
		assert.NoError(t, err)

		_, err = flat.GetEnrichedActivities(tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.enrichedURL, "")
		assert.NoError(t, err)

		_, err = flat.GetEnrichedActivitiesWithRanking("popularity", tc.opts...)
		testRequest(t, requester.req, http.MethodGet, fmt.Sprintf("%s&ranking=popularity", tc.enrichedURL), "")
		assert.NoError(t, err)
	}
}

func TestFlatFeedGetNextPageActivities(t *testing.T) {
	client, requester := newClient(t)
	flat, _ := newFlatFeedWithUserID(client, "123")

	requester.resp = `{"next":"/api/v1.0/feed/flat/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	resp, err := flat.GetActivities()
	require.NoError(t, err)

	_, err = flat.GetNextPageActivities(resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/enrich/feed/flat/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	enrichedResp, err := flat.GetEnrichedActivities()
	require.NoError(t, err)

	_, err = flat.GetNextPageEnrichedActivities(enrichedResp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":123}`
	_, err = flat.GetActivities()
	require.Error(t, err)

	requester.resp = `{"next":"123"}`
	resp, err = flat.GetActivities()
	require.NoError(t, err)
	_, err = flat.GetNextPageActivities(resp)
	require.Error(t, err)

	requester.resp = `{"next":"?q=a%"}`
	resp, err = flat.GetActivities()
	require.NoError(t, err)
	_, err = flat.GetNextPageActivities(resp)
	require.Error(t, err)
}
