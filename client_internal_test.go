package stream

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
		c, err := NewClient(tc.key, tc.secret, tc.opts...)
		if tc.shouldError {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedRegion, c.url.region)
		assert.Equal(t, tc.expectedVersion, c.url.version)
	}
}

func Test_makeEndpoint(t *testing.T) {
	prev := os.Getenv("STREAM_URL")
	defer os.Setenv("STREAM_URL", prev)

	testCases := []struct {
		url      *apiURL
		format   string
		env      string
		args     []interface{}
		expected string
	}{
		{
			url:      &apiURL{},
			format:   "test-%d-%s",
			args:     []interface{}{42, "asd"},
			expected: "https://api.stream-io-api.com/api/v1.0/test-42-asd?api_key=test",
		},
		{
			url:      &apiURL{},
			env:      "http://localhost:8000/api/v1.0/",
			format:   "test-%d-%s",
			args:     []interface{}{42, "asd"},
			expected: "http://localhost:8000/api/v1.0/test-42-asd?api_key=test",
		},
	}

	for _, tc := range testCases {
		os.Setenv("STREAM_URL", tc.env)
		c := &Client{url: tc.url, key: "test"}
		assert.Equal(t, tc.expected, c.makeEndpoint(tc.format, tc.args...).String())
	}
}
