package stream_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"

	stream "github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newClient(t *testing.T) (*stream.Client, *mockRequester) {
	requester := &mockRequester{}
	client, err := stream.NewClient("key", "secret", stream.WithHTTPRequester(requester))
	require.NoError(t, err)
	return client, requester
}

type mockRequester struct {
	req  *http.Request
	resp string
}

func (m *mockRequester) Do(req *http.Request) (*http.Response, error) {
	m.req = req
	var body string
	if m.resp != "" {
		body = m.resp
	} else {
		body = "{}"
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
	return resp, nil
}

func testRequest(t *testing.T, req *http.Request, method, url, body string) {
	assert.Equal(t, url, req.URL.String())
	assert.Equal(t, method, req.Method)
	if req.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)
		assert.Equal(t, body, string(reqBody))
	}
	headers := req.Header
	if headers.Get("X-API-Key") == "" {
		assert.NotEmpty(t, headers.Get("Stream-Auth-Type"))
		assert.NotEmpty(t, headers.Get("Authorization"))
	}
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func getTime(t time.Time) stream.Time {
	st, _ := time.Parse(stream.TimeLayout, t.Truncate(time.Second).Format(stream.TimeLayout))
	return stream.Time{Time: st}
}

func newFlatFeed(c *stream.Client) *stream.FlatFeed {
	return newFlatFeedWithUserID(c, randString(10))
}

func newFlatFeedWithUserID(c *stream.Client, userID string) *stream.FlatFeed {
	return c.FlatFeed("flat", userID)
}

func newAggregatedFeed(c *stream.Client) *stream.AggregatedFeed {
	return newAggregatedFeedWithUserID(c, randString(10))
}

func newAggregatedFeedWithUserID(c *stream.Client, userID string) *stream.AggregatedFeed {
	return c.AggregatedFeed("aggregated", userID)
}

func newNotificationFeed(c *stream.Client) *stream.NotificationFeed {
	return newNotificationFeedWithUserID(c, randString(10))
}

func newNotificationFeedWithUserID(c *stream.Client, userID string) *stream.NotificationFeed {
	return c.NotificationFeed("notification", userID)
}
