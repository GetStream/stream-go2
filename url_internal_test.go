package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_URLString(t *testing.T) {
	testCases := []struct {
		url      *apiURL
		expected string
	}{
		{
			url:      &apiURL{},
			expected: "https://api.getstream.io/api/v1.0/",
		},
		{
			url:      &apiURL{region: "eu-central", version: "2.0"},
			expected: "https://eu-central-api.getstream.io/api/v2.0/",
		},
		{
			url:      &apiURL{region: "localhost"},
			expected: "http://localhost:8000/api/v1.0/",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tc.url.String())
	}
}
