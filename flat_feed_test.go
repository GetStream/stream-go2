package stream_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v8"
)

func TestFlatFeedGetActivities(t *testing.T) {
	ctx := context.Background()
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
		_, err := flat.GetActivities(ctx, tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)

		_, err = flat.GetActivitiesWithRanking(ctx, "popularity", tc.opts...)
		testRequest(t, requester.req, http.MethodGet, fmt.Sprintf("%s&ranking=popularity", tc.url), "")
		assert.NoError(t, err)

		_, err = flat.GetEnrichedActivities(ctx, tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.enrichedURL, "")
		assert.NoError(t, err)

		_, err = flat.GetEnrichedActivitiesWithRanking(ctx, "popularity", tc.opts...)
		testRequest(t, requester.req, http.MethodGet, fmt.Sprintf("%s&ranking=popularity", tc.enrichedURL), "")
		assert.NoError(t, err)
	}
}
func TestFlatFeedGetActivitiesExternalRanking(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	flat, _ := newFlatFeedWithUserID(client, "123")
	testCases := []struct {
		opts        []stream.GetActivitiesOption
		url         string
		enrichedURL string
	}{
		{
			name : "external ranking vars",
			opts: []stream.GetActivitiesOption{
				stream.WithExternalRankingVars(map[string]any{
					"music":   1,
					"sports":  2.1,
					"boolVal": true,
					"string":  "str",
				}),
			},
			//url:         "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&ranking_vars=%7B%22boolVal%22%3Atrue%2C%22music%22%3A1%2C%22sports%22%3A2.1%2C%22string%22%3A%22str%22%7D&ranking=popularity",
			url:         "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&ranking_vars=%7B%22boolVal%22%3Atrue%2C%22music%22%3A1%2C%22sports%22%3A2.1%2C%22string%22%3A%22str%22%7D",
			enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?api_key=key&ranking_vars=%7B%22boolVal%22%3Atrue%2C%22music%22%3A1%2C%22sports%22%3A2.1%2C%22string%22%3A%22str%22%7D",
			//enrichedURL: "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?api_key=key&ranking_vars=%7B%22boolVal%22%3Atrue%2C%22music%22%3A1%2C%22sports%22%3A2.1%2C%22string%22%3A%22str%22%7D&ranking=popularity",
		},
	}

	for _, tc := range testCases {
		_, err := flat.GetActivities(ctx, tc.opts...)
		testRequest(t, requester.req, http.MethodGet, tc.url, "")
		assert.NoError(t, err)
	}
}

func TestFlatFeedGetNextPageActivities(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	flat, _ := newFlatFeedWithUserID(client, "123")

	requester.resp = `{"next":"/api/v1.0/feed/flat/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	resp, err := flat.GetActivities(ctx)
	require.NoError(t, err)

	_, err = flat.GetNextPageActivities(ctx, resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/feed/flat/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/enrich/feed/flat/123/?id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25"}`
	enrichedResp, err := flat.GetEnrichedActivities(ctx)
	require.NoError(t, err)

	_, err = flat.GetNextPageEnrichedActivities(ctx, enrichedResp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/enrich/feed/flat/123/?api_key=key&id_lt=78c1a709-aff2-11e7-b3a7-a45e60be7d3b&limit=25", "")
	require.NoError(t, err)

	requester.resp = `{"next":123}`
	_, err = flat.GetActivities(ctx)
	require.Error(t, err)

	requester.resp = `{"next":"123"}`
	resp, err = flat.GetActivities(ctx)
	require.NoError(t, err)
	_, err = flat.GetNextPageActivities(ctx, resp)
	require.Error(t, err)

	requester.resp = `{"next":"?q=a%"}`
	resp, err = flat.GetActivities(ctx)
	require.NoError(t, err)
	_, err = flat.GetNextPageActivities(ctx, resp)
	require.Error(t, err)
}

func TestFlatFeedFollowStats(t *testing.T) {
	ctx := context.Background()
	endpoint := "https://api.stream-io-api.com/api/v1.0/stats/follow/?api_key=key"

	client, requester := newClient(t)
	flat, _ := newFlatFeedWithUserID(client, "123")

	_, err := flat.FollowStats(ctx)
	testRequest(t, requester.req, http.MethodGet, endpoint+"&followers=flat%3A123&following=flat%3A123", "")
	assert.NoError(t, err)

	_, err = flat.FollowStats(ctx, stream.WithFollowerSlugs("a", "b"))
	testRequest(t, requester.req, http.MethodGet, endpoint+"&followers=flat%3A123&followers_slugs=a%2Cb&following=flat%3A123", "")
	assert.NoError(t, err)

	_, err = flat.FollowStats(ctx, stream.WithFollowingSlugs("c", "d"))
	testRequest(t, requester.req, http.MethodGet, endpoint+"&followers=flat%3A123&following=flat%3A123&following_slugs=c%2Cd", "")
	assert.NoError(t, err)

	_, err = flat.FollowStats(ctx, stream.WithFollowingSlugs("c", "d"), stream.WithFollowerSlugs("a", "b"))
	testRequest(t, requester.req, http.MethodGet, endpoint+"&followers=flat%3A123&followers_slugs=a%2Cb&following=flat%3A123&following_slugs=c%2Cd", "")
	assert.NoError(t, err)
}
