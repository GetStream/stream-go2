package stream

import (
	"fmt"
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
			expected: fmt.Sprintf("https://api.%s/api/v1.0/", domain),
		},
		{
			url:      &apiURL{region: "eu-central", version: "2.0"},
			expected: fmt.Sprintf("https://eu-central-api.%s/api/v2.0/", domain),
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tc.url.String())
	}
}
