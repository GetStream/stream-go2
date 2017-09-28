package stream_test

import (
	"encoding/json"
	"testing"
	"time"

	stream "github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDurationMarshalUnmarshalJSON(t *testing.T) {
	dur := stream.Duration{Duration: 33 * time.Second}
	data := []byte(`"33s"`)
	marshaled, err := json.Marshal(dur)
	assert.NoError(t, err)
	assert.Equal(t, data, marshaled)
	var out stream.Duration
	err = json.Unmarshal(marshaled, &out)
	assert.NoError(t, err)
	assert.Equal(t, dur, out)
}

func TestTimeMarshalUnmarshalJSON(t *testing.T) {
	tt, _ := time.Parse("2006-Jan-02", "2013-Feb-03")
	st := stream.Time{Time: tt}
	data := []byte(`"2013-02-03T00:00:00"`)
	marshaled, err := json.Marshal(st)
	require.NoError(t, err)
	require.Equal(t, data, marshaled)
	var out stream.Time
	err = json.Unmarshal(marshaled, &out)
	assert.NoError(t, err)
	assert.Equal(t, st, out)
}
