package stream

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_decodeJSONHook(t *testing.T) {
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
		{
			f:           reflect.TypeOf(float64(0)),
			typ:         reflect.TypeOf(Duration{}),
			data:        float64(42),
			shouldError: false,
			expected:    Duration{time.Second * 42},
		},
		{
			f:           reflect.TypeOf(float64(0)),
			typ:         reflect.TypeOf(Duration{}),
			data:        struct{}{},
			shouldError: true,
		},
		{
			f:           reflect.TypeOf(string("")),
			typ:         reflect.TypeOf(Data{}),
			data:        "test",
			shouldError: false,
			expected:    Data{ID: "test"},
		},
		{
			f:   reflect.TypeOf(map[string]interface{}{}),
			typ: reflect.TypeOf(Data{}),
			data: map[string]interface{}{
				"id":    "test",
				"extra": "data",
			},
			shouldError: false,
			expected:    Data{ID: "test", Extra: map[string]interface{}{"extra": "data"}},
		},
		{
			f:           reflect.TypeOf(float64(0)),
			typ:         reflect.TypeOf(Data{}),
			data:        struct{}{},
			shouldError: true,
		},
	}
	for _, tc := range testCases {
		out, err := decodeJSONHook(tc.f, tc.typ, tc.data)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.expected, out)
		}
	}
}

func Test_parseIntValue(t *testing.T) {
	testCases := []struct {
		values       url.Values
		shouldError  bool
		expected     int
		expectedFlag bool
	}{
		{
			values:       url.Values{},
			shouldError:  false,
			expected:     0,
			expectedFlag: false,
		},
		{
			values:       url.Values{"test": []string{"a"}},
			shouldError:  true,
			expected:     0,
			expectedFlag: false,
		},
		{
			values:       url.Values{"test": []string{"123"}},
			shouldError:  false,
			expected:     123,
			expectedFlag: true,
		},
		{
			values:       url.Values{"test": []string{"123.5"}},
			shouldError:  true,
			expected:     0,
			expectedFlag: false,
		},
	}
	for _, tc := range testCases {
		v, ok, err := parseIntValue(tc.values, "test")
		if tc.shouldError {
			require.Error(t, err)
			assert.False(t, ok)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFlag, ok)
			assert.Equal(t, tc.expected, v)
		}
	}
}
