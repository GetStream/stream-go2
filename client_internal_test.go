package stream

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		key             string
		secret          string
		shouldError     bool
		opts            []ClientOption
		expectedRegion  string
		expectedVersion string
	}{
		{
			shouldError: true,
		},
		{
			key: "k", secret: "s",
			expectedRegion:  "",
			expectedVersion: "",
		},
		{
			key: "k", secret: "s",
			opts:            []ClientOption{WithAPIRegion("test")},
			expectedRegion:  "test",
			expectedVersion: "",
		},
		{
			key: "k", secret: "s",
			opts:            []ClientOption{WithAPIVersion("test")},
			expectedRegion:  "",
			expectedVersion: "test",
		},
		{
			key: "k", secret: "s",
			opts:            []ClientOption{WithAPIRegion("test"), WithAPIVersion("more")},
			expectedRegion:  "test",
			expectedVersion: "more",
		},
	}
	for _, tc := range testCases {
		c, err := New(tc.key, tc.secret, tc.opts...)
		if tc.shouldError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedRegion, c.urlBuilder.(apiURLBuilder).region)
		assert.Equal(t, tc.expectedVersion, c.urlBuilder.(apiURLBuilder).version)
	}
}

func Test_makeEndpoint(t *testing.T) {
	prev := os.Getenv("STREAM_URL")
	defer os.Setenv("STREAM_URL", prev)

	testCases := []struct {
		urlBuilder apiURLBuilder
		format     string
		env        string
		args       []interface{}
		expected   string
	}{
		{
			urlBuilder: apiURLBuilder{},
			format:     "test-%d-%s",
			args:       []interface{}{42, "asd"},
			expected:   "https://api.stream-io-api.com/api/v1.0/test-42-asd?api_key=test",
		},
		{
			urlBuilder: apiURLBuilder{},
			env:        "http://localhost:8000",
			format:     "test-%d-%s",
			args:       []interface{}{42, "asd"},
			expected:   "http://localhost:8000/api/v1.0/test-42-asd?api_key=test",
		},
	}

	for _, tc := range testCases {
		os.Setenv("STREAM_URL", tc.env)
		c := &Client{urlBuilder: tc.urlBuilder, key: "test"}
		assert.Equal(t, tc.expected, c.makeEndpoint(tc.format, tc.args...).String())
	}
}

func TestNewFromEnv(t *testing.T) {
	defer func() {
		os.Setenv("STREAM_API_KEY", "")
		os.Setenv("STREAM_API_SECRET", "")
		os.Setenv("STREAM_API_REGION", "")
		os.Setenv("STREAM_API_VERSION", "")
	}()

	_, err := NewFromEnv()
	require.Error(t, err)

	os.Setenv("STREAM_API_KEY", "foo")
	os.Setenv("STREAM_API_SECRET", "bar")

	client, err := NewFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "foo", client.key)
	assert.Equal(t, "bar", client.authenticator.secret)

	os.Setenv("STREAM_API_REGION", "baz")
	client, err = NewFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "baz", client.urlBuilder.(apiURLBuilder).region)

	os.Setenv("STREAM_API_VERSION", "qux")
	client, err = NewFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "qux", client.urlBuilder.(apiURLBuilder).version)
}

type badReader struct{}

func (badReader) Read([]byte) (int, error) { return 0, fmt.Errorf("boom") }

func Test_makeStreamError(t *testing.T) {
	testCases := []struct {
		body     io.Reader
		expected error
		apiErr   APIError
	}{
		{
			body:     nil,
			expected: fmt.Errorf("invalid body"),
		},
		{
			body:     badReader{},
			expected: fmt.Errorf("boom"),
		},
		{
			body:     strings.NewReader(`{{`),
			expected: fmt.Errorf("unexpected error (status code 123)"),
		},
		{
			body:     strings.NewReader(`{"code":"A"}`),
			expected: fmt.Errorf("unexpected error (status code 123)"),
		},
		{
			body:     strings.NewReader(`{"code":1, "detail":"test", "duration": "1m2s", "exception": "boom", "status_code": 456, "exception_fields": {"foo":["bar"]}}`),
			expected: fmt.Errorf("test"),
			apiErr: APIError{
				Code:       1,
				Detail:     "test",
				Duration:   Duration{time.Minute + time.Second*2},
				Exception:  "boom",
				StatusCode: 123,
				ExceptionFields: map[string][]interface{}{
					"foo": {"bar"},
				},
			},
		},
	}
	for _, tc := range testCases {
		err := (&Client{}).makeStreamError(123, nil, tc.body)
		assert.Equal(t, tc.expected.Error(), err.Error())
		if tc.apiErr.Code != 0 {
			assert.Equal(t, tc.apiErr, err)
		}
	}
}

type requester struct {
	code int
	body io.ReadCloser
	err  error
}

func (r requester) Do(*http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: r.code,
		Body:       r.body,
	}
	return resp, r.err
}

func Test_requestErrors(t *testing.T) {
	testCases := []struct {
		data      interface{}
		method    string
		authFn    authFunc
		expected  error
		requester Requester
	}{
		{
			data:     make(chan int),
			expected: fmt.Errorf("cannot marshal request: json: unsupported type: chan int"),
		},
		{
			data:     42,
			authFn:   func(*http.Request) error { return fmt.Errorf("boom") },
			expected: fmt.Errorf("boom"),
		},
		{
			data:     42,
			method:   "Ω",
			expected: fmt.Errorf(`cannot create request: net/http: invalid method "Ω"`),
		},
		{
			data:      42,
			authFn:    func(*http.Request) error { return nil },
			requester: &requester{err: fmt.Errorf("boom")},
			expected:  fmt.Errorf("cannot perform request: boom"),
		},
		{
			data:      42,
			authFn:    func(*http.Request) error { return nil },
			requester: &requester{code: 400, body: ioutil.NopCloser(strings.NewReader(`{"detail":"boom"}`))},
			expected:  fmt.Errorf("boom"),
		},
		{
			data:      42,
			authFn:    func(*http.Request) error { return nil },
			requester: &requester{code: 200, body: ioutil.NopCloser(badReader{})},
			expected:  fmt.Errorf("cannot read response: boom"),
		},
	}

	for _, tc := range testCases {
		c := &Client{requester: tc.requester}
		_, err := c.request(tc.method, endpoint{url: &url.URL{}, query: url.Values{}}, tc.data, tc.authFn)
		require.Error(t, err)
		assert.Equal(t, tc.expected.Error(), err.Error())
	}
}
