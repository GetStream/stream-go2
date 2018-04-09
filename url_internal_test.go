package stream

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_URLString(t *testing.T) {
	testCases := []struct {
		urlBuilder apiURLBuilder
		expected   string
	}{
		{
			urlBuilder: apiURLBuilder{},
			expected:   fmt.Sprintf("https://api.%s/api/v1.0/", domain),
		},
		{
			urlBuilder: newAPIURLBuilder("us-east", "2.0"),
			expected:   fmt.Sprintf("https://us-east-api.%s/api/v2.0/", domain),
		},
		{
			urlBuilder: newAPIURLBuilder("eu-west", "2.0"),
			expected:   fmt.Sprintf("https://eu-west-api.%s/api/v2.0/", domain),
		},
		{
			urlBuilder: newAPIURLBuilder("singapore", "2.0"),
			expected:   fmt.Sprintf("https://singapore-api.%s/api/v2.0/", domain),
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tc.urlBuilder.url())
	}
}
