package stream

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadResponse_parseNext(t *testing.T) {
	testCases := []struct {
		next        string
		shouldError bool
		err         error
		expected    []GetActivitiesOption
	}{
		{
			next:        "",
			shouldError: true,
			err:         ErrMissingNextPage,
		},
		{
			next:        "/test",
			shouldError: true,
			err:         ErrInvalidNextPage,
		},
		{
			next:        "/test?k=%",
			shouldError: true,
			err:         ErrInvalidNextPage,
		},
		{
			next:        "/test?limit=a",
			shouldError: true,
			err:         fmt.Errorf(`strconv.Atoi: parsing "a": invalid syntax`),
		},
		{
			next:        "/test?offset=a",
			shouldError: true,
			err:         fmt.Errorf(`strconv.Atoi: parsing "a": invalid syntax`),
		},
		{
			next:        "/test?limit=1&offset=2",
			shouldError: false,
			expected: []GetActivitiesOption{
				WithActivitiesLimit(1),
				WithActivitiesOffset(2),
			},
		},
		{
			next:        "/test?limit=1&offset=2&id_lt=foo",
			shouldError: false,
			expected: []GetActivitiesOption{
				WithActivitiesLimit(1),
				WithActivitiesOffset(2),
				WithActivitiesIDLT("foo"),
			},
		},
		{
			next:        "/test?limit=1&offset=2&id_lt=foo&ranking=bar",
			shouldError: false,
			expected: []GetActivitiesOption{
				WithActivitiesLimit(1),
				WithActivitiesOffset(2),
				WithActivitiesIDLT("foo"),
				withActivitiesRanking("bar"),
			},
		},
	}

	for _, tc := range testCases {
		r := readResponse{Next: tc.next}
		opts, err := r.parseNext()
		if tc.shouldError {
			require.Error(t, err)
			assert.Equal(t, tc.err.Error(), err.Error())
		} else {
			require.NoError(t, err)
			assert.Equal(t, tc.expected, opts)
		}
	}
}
