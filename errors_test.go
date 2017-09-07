package stream_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/assert"
)

func TestErrorUnmarshal(t *testing.T) {
	data := []byte(`{"code":42,"detail":"the details","duration":"10ms","exception":"boom","status_code":123}`)
	var apiErr stream.APIError
	err := json.Unmarshal(data, &apiErr)
	assert.NoError(t, err)
	expected := stream.APIError{
		Code:       42,
		Detail:     "the details",
		Duration:   10 * time.Millisecond,
		Exception:  "boom",
		StatusCode: 123,
	}
	assert.Equal(t, expected, apiErr)
}
