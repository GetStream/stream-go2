package stream_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	stream "github.com/flockfw64/stream-go2/v7"
)

func TestErrorUnmarshal(t *testing.T) {
	data := []byte(`{"code":42,"detail":"the details","duration":"10ms","exception":"boom","status_code":123}`)
	var apiErr stream.APIError
	err := json.Unmarshal(data, &apiErr)
	assert.NoError(t, err)
	expected := stream.APIError{
		Code:       42,
		Detail:     "the details",
		Duration:   stream.Duration{Duration: 10 * time.Millisecond},
		Exception:  "boom",
		StatusCode: 123,
	}
	assert.Equal(t, expected, apiErr)
}

func TestErrorString(t *testing.T) {
	apiErr := stream.APIError{Detail: "boom"}
	assert.Equal(t, "boom", apiErr.Error())
}

func TestToAPIError(t *testing.T) {
	testCases := []struct {
		err   error
		match bool
	}{
		{
			err:   fmt.Errorf("this is an error"),
			match: false,
		},
		{
			err:   stream.APIError{},
			match: true,
		},
	}

	for _, tc := range testCases {
		err, ok := stream.ToAPIError(tc.err)
		assert.Equal(t, tc.match, ok)
		if ok {
			assert.IsType(t, stream.APIError{}, err)
		}
	}
}
