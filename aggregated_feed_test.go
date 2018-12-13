package stream_test

import (
	"net/http"
	"testing"

	stream "github.com/GetStream/stream-go2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregatedFeedGetActivities(t *testing.T) {
	client, requester := newClient(t)
	aggregated := newAggregatedFeedWithUserID(client, "123")
	testCases := []struct {
		opts        []stream.GetActivitiesOption
		url         string
		enrichedURL string
	}{
		{
			url:         "https://api.stream-io-api.com/api/v1.0/feed/aggregated/123/?api_key=key",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/aggregated/123/?api_key=key",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithActivitiesLimit(42)},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/aggregated/123/?api_key=key&limit=42",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/aggregated/123/?api_key=key&limit=42",
		},
		{
			opts:        []stream.GetActivitiesOption{stream.WithActivitiesLimit(42), stream.WithActivitiesOffset(11), stream.WithActivitiesIDGT("aabbcc")},
			url:         "https://api.stream-io-api.com/api/v1.0/feed/aggregated/123/?api_key=key&id_gt=aabbcc&limit=42&offset=11",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/aggregated/123/?api_key=key&id_gt=aabbcc&limit=42&offset=11",
		},
	}

	for _, tc := range testCases {
		_, err := aggregated.GetActivities(tc.opts...)
		assert.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		_, err = aggregated.GetEnrichedActivities(tc.opts...)
		assert.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.enrichedURL, "")
	}
}

func TestAggregatedFeedGetNextPageActivities(t *testing.T) {
	client, requester := newClient(t)
	aggregated := newAggregatedFeedWithUserID(client, "123")

	requester.resp = `{"next":"/api/v1.0/feed/aggregated/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	resp, err := aggregated.GetActivities()
	require.NoError(t, err)
	_, err = aggregated.GetNextPageActivities(resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/feed/aggregated/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/enrich/feed/aggregated/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	enrichedResp, err := aggregated.GetEnrichedActivities()
	require.NoError(t, err)
	_, err = aggregated.GetNextPageEnrichedActivities(enrichedResp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/enrich/feed/aggregated/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":123}`
	_, err = aggregated.GetActivities()
	require.Error(t, err)

	requester.resp = `{"next":"123"}`
	resp, err = aggregated.GetActivities()
	require.NoError(t, err)
	_, err = aggregated.GetNextPageActivities(resp)
	require.Error(t, err)

	requester.resp = `{"next":"?q=a%"}`
	resp, err = aggregated.GetActivities()
	require.NoError(t, err)
	_, err = aggregated.GetNextPageActivities(resp)
	require.Error(t, err)
}
