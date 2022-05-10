package stream_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v7"
)

func TestAggregatedFeedGetActivities(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	aggregated, _ := newAggregatedFeedWithUserID(client, "123")
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
		_, err := aggregated.GetActivities(ctx, tc.opts...)
		assert.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")

		_, err = aggregated.GetActivitiesWithRanking(ctx, "popularity", tc.opts...)
		testRequest(t, requester.req, http.MethodGet, fmt.Sprintf("%s&ranking=popularity", tc.url), "")
		assert.NoError(t, err)

		_, err = aggregated.GetEnrichedActivities(ctx, tc.opts...)
		assert.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.enrichedURL, "")

		_, err = aggregated.GetEnrichedActivitiesWithRanking(ctx, "popularity", tc.opts...)
		testRequest(t, requester.req, http.MethodGet, fmt.Sprintf("%s&ranking=popularity", tc.enrichedURL), "")
		assert.NoError(t, err)
	}
}

func TestAggregatedFeedGetNextPageActivities(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	aggregated, _ := newAggregatedFeedWithUserID(client, "123")

	requester.resp = `{"next":"/api/v1.0/feed/aggregated/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	resp, err := aggregated.GetActivities(ctx)
	require.NoError(t, err)
	_, err = aggregated.GetNextPageActivities(ctx, resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/feed/aggregated/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/enrich/feed/aggregated/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	enrichedResp, err := aggregated.GetEnrichedActivities(ctx)
	require.NoError(t, err)
	_, err = aggregated.GetNextPageEnrichedActivities(ctx, enrichedResp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/enrich/feed/aggregated/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":123}`
	_, err = aggregated.GetActivities(ctx)
	require.Error(t, err)

	requester.resp = `{"next":"123"}`
	resp, err = aggregated.GetActivities(ctx)
	require.NoError(t, err)
	_, err = aggregated.GetNextPageActivities(ctx, resp)
	require.Error(t, err)

	requester.resp = `{"next":"?q=a%"}`
	resp, err = aggregated.GetActivities(ctx)
	require.NoError(t, err)
	_, err = aggregated.GetNextPageActivities(ctx, resp)
	require.Error(t, err)
}
