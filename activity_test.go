package stream_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v4"
)

func TestActivityMarshalUnmarshalJSON(t *testing.T) {
	now := getTime(time.Now())
	testCases := []struct {
		activity stream.Activity
		data     []byte
	}{
		{
			activity: stream.Activity{Actor: "actor", Verb: "verb", Object: "object"},
			data:     []byte(`{"actor":"actor","object":"object","verb":"verb"}`),
		},
		{
			activity: stream.Activity{Actor: "actor", Verb: "verb", Object: "object", Time: now},
			data:     []byte(`{"actor":"actor","object":"object","time":"` + now.Format(stream.TimeLayout) + `","verb":"verb"}`),
		},
		{
			activity: stream.Activity{Actor: "actor", Verb: "verb", Object: "object", Time: now, Extra: map[string]interface{}{"popularity": 42.0, "size": map[string]interface{}{"width": 800.0, "height": 600.0}}},
			data:     []byte(`{"actor":"actor","object":"object","popularity":42,"size":{"height":600,"width":800},"time":"` + now.Format(stream.TimeLayout) + `","verb":"verb"}`),
		},
		{
			activity: stream.Activity{Actor: "actor", Verb: "verb", Object: "object", Time: now, Extra: map[string]interface{}{"popularity": 42.0, "size": map[string]interface{}{"width": 800.0, "height": 600.0}}},
			data:     []byte(`{"actor":"actor","object":"object","popularity":42,"size":{"height":600,"width":800},"time":"` + now.Format(stream.TimeLayout) + `","verb":"verb"}`),
		},
		{
			activity: stream.Activity{To: []string{"abcd", "efgh"}},
			data:     []byte(`{"to":["abcd","efgh"]}`),
		},
	}
	for _, tc := range testCases {
		data, err := json.Marshal(tc.activity)
		assert.NoError(t, err)
		assert.Equal(t, tc.data, data)

		var out stream.Activity
		err = json.Unmarshal(tc.data, &out)
		require.NoError(t, err)
		assert.Equal(t, tc.activity, out)
	}
}

func TestActivityMarshalUnmarshalJSON_toTargets(t *testing.T) {
	testCases := []struct {
		activity    stream.Activity
		data        []byte
		shouldError bool
	}{
		{
			activity: stream.Activity{To: []string{"abcd", "efgh"}},
			data:     []byte(`{"to":["abcd","efgh"]}`),
		},
		{
			activity: stream.Activity{To: []string{"abcd", "efgh"}},
			data:     []byte(`{"to":[["abcd", "foo"], ["efgh", "bar"]]}`),
		},
		{
			activity:    stream.Activity{To: []string{"abcd", "efgh"}},
			data:        []byte(`{"to":[[123]]}`),
			shouldError: true,
		},
	}
	for _, tc := range testCases {
		var out stream.Activity
		err := json.Unmarshal(tc.data, &out)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.activity, out)
		}
	}
}
