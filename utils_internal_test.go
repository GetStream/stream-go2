package stream

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_decodeJSONStringTimes(t *testing.T) {
	now, _ := time.Parse(TimeLayout, "2006-01-02T15:04:05.999999")
	testCases := []struct {
		f           reflect.Type
		typ         reflect.Type
		data        interface{}
		expected    interface{}
		shouldError bool
	}{
		{
			f:        reflect.TypeOf(123),
			data:     123,
			expected: 123,
		},
		{
			f:        reflect.TypeOf(""),
			typ:      reflect.TypeOf(Duration{}),
			data:     "1m2s",
			expected: Duration{time.Minute + time.Second*2},
		},
		{
			f:           reflect.TypeOf(""),
			typ:         reflect.TypeOf(Duration{}),
			data:        "test",
			shouldError: true,
		},
		{
			f:        reflect.TypeOf(""),
			typ:      reflect.TypeOf(Time{}),
			data:     now.Format(TimeLayout),
			expected: Time{now},
		},
		{
			f:           reflect.TypeOf(""),
			typ:         reflect.TypeOf(Time{}),
			data:        "test",
			shouldError: true,
		},
	}
	for _, tc := range testCases {
		out, err := decodeJSONStringTimes(tc.f, tc.typ, tc.data)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.expected, out)
		}
	}
}
