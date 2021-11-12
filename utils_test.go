package stream_test

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v6"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newClient(t *testing.T) (*stream.Client, *mockRequester) {
	requester := &mockRequester{}
	client, err := stream.New("key", "secret", stream.WithHTTPRequester(requester))
	require.NoError(t, err)
	return client, requester
}

type mockRequester struct {
	req  *http.Request
	resp string
}

func (m *mockRequester) Do(req *http.Request) (*http.Response, error) {
	m.req = req
	body := "{}"
	if m.resp != "" {
		body = m.resp
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func testRequest(t *testing.T, req *http.Request, method, url, body string) {
	assert.Equal(t, url, req.URL.String())
	assert.Equal(t, method, req.Method)
	if req.Method == http.MethodPost {
		reqBody, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		assert.JSONEq(t, body, string(reqBody))
	}
	headers := req.Header
	if headers.Get("X-API-Key") == "" {
		assert.NotEmpty(t, headers.Get("Stream-Auth-Type"))
		assert.NotEmpty(t, headers.Get("Authorization"))
	}
}

func getTime(t time.Time) stream.Time {
	st, _ := time.Parse(stream.TimeLayout, t.Truncate(time.Second).Format(stream.TimeLayout))
	return stream.Time{Time: st}
}

func newFlatFeedWithUserID(c *stream.Client, userID string) (*stream.FlatFeed, error) {
	return c.FlatFeed("flat", userID)
}

func newAggregatedFeedWithUserID(c *stream.Client, userID string) (*stream.AggregatedFeed, error) {
	return c.AggregatedFeed("aggregated", userID)
}

func newNotificationFeedWithUserID(c *stream.Client, userID string) (*stream.NotificationFeed, error) {
	return c.NotificationFeed("notification", userID)
}
