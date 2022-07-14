package stream_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/flockfw64/stream-go2/v7"
)

func TestEnrichedActivityUnmarshalJSON(t *testing.T) {
	now := getTime(time.Now())
	testCases := []struct {
		activity stream.EnrichedActivity
		data     []byte
	}{
		{
			activity: stream.EnrichedActivity{Actor: stream.Data{ID: "actor"}, Verb: "verb", Object: stream.Data{ID: "object"}},
			data:     []byte(`{"actor":"actor","object":"object","verb":"verb"}`),
		},
		{
			activity: stream.EnrichedActivity{Actor: stream.Data{ID: "actor"}, Verb: "verb", Object: stream.Data{ID: "object"}, Time: now},
			data:     []byte(`{"actor":"actor","object":"object","time":"` + now.Format(stream.TimeLayout) + `","verb":"verb"}`),
		},
		{
			activity: stream.EnrichedActivity{Actor: stream.Data{ID: "actor"}, Verb: "verb", Object: stream.Data{ID: "object"}, Time: now, Extra: map[string]interface{}{"popularity": 42.0, "size": map[string]interface{}{"width": 800.0, "height": 600.0}}},
			data:     []byte(`{"actor":"actor","object":"object","popularity":42,"size":{"height":600,"width":800},"time":"` + now.Format(stream.TimeLayout) + `","verb":"verb"}`),
		},
		{
			activity: stream.EnrichedActivity{Actor: stream.Data{ID: "actor"}, Verb: "verb", Object: stream.Data{ID: "object"}, Time: now, Extra: map[string]interface{}{"popularity": 42.0, "size": map[string]interface{}{"width": 800.0, "height": 600.0}}},
			data:     []byte(`{"actor":"actor","object":"object","popularity":42,"size":{"height":600,"width":800},"time":"` + now.Format(stream.TimeLayout) + `","verb":"verb"}`),
		},
		{
			activity: stream.EnrichedActivity{To: []string{"abcd", "efgh"}},
			data:     []byte(`{"to":["abcd","efgh"]}`),
		},
		{
			activity: stream.EnrichedActivity{
				ForeignID: "SA:123",
				Extra: map[string]interface{}{
					"foreign_id_ref": map[string]interface{}{
						"id":         "123",
						"extra_prop": true,
					},
				},
			},
			data: []byte(`{"foreign_id":{"id":"123","extra_prop":true}}`),
		},
	}
	for _, tc := range testCases {
		var out stream.EnrichedActivity
		err := json.Unmarshal(tc.data, &out)
		require.NoError(t, err)
		assert.Equal(t, tc.activity, out)
	}
}

func TestEnrichedActivityUnmarshalJSON_toTargets(t *testing.T) {
	testCases := []struct {
		activity    stream.EnrichedActivity
		data        []byte
		shouldError bool
	}{
		{
			activity: stream.EnrichedActivity{To: []string{"abcd", "efgh"}},
			data:     []byte(`{"to":["abcd","efgh"]}`),
		},
		{
			activity: stream.EnrichedActivity{To: []string{"abcd", "efgh"}},
			data:     []byte(`{"to":[["abcd", "foo"], ["efgh", "bar"]]}`),
		},
		{
			activity:    stream.EnrichedActivity{To: []string{"abcd", "efgh"}},
			data:        []byte(`{"to":[[123]]}`),
			shouldError: true,
		},
	}
	for _, tc := range testCases {
		var out stream.EnrichedActivity
		err := json.Unmarshal(tc.data, &out)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.activity, out)
		}
	}
}
