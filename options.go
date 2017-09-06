package stream

import (
	"fmt"
	"strings"
)

const (
	defaultActivityCopyLimit = 300
)

// RequestOption is an interface representing API request optional filters and
// parameters.
type RequestOption interface {
	String() string
}

type requestOption struct {
	key   string
	value interface{}
}

func makeRequestOption(key string, value interface{}) RequestOption {
	return requestOption{key: key, value: value}
}

func (o requestOption) String() string {
	if o.key == "" {
		return ""
	}
	return fmt.Sprintf("&%s=%v", o.key, o.value)
}

// WithLimit adds the limit parameter to API calls which support it, limiting
// the number of results in the response to the provided limit threshold.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func WithLimit(limit int) RequestOption {
	return makeRequestOption("limit", limit)
}

// WithOffset adds the offset parameter to API calls which support it, getting
// results starting from the provided offset index.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func WithOffset(offset int) RequestOption {
	return makeRequestOption("offset", offset)
}

// WithRanking adds the ranking parameter to API calls, used when retrieving
// flat feed activities with a specified ranking method.
func WithRanking(ranking string) RequestOption {
	return makeRequestOption("ranking", ranking)
}

// WithIDGTE adds the id_gte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater or
// equal than the provided id.
func WithIDGTE(id string) RequestOption {
	return makeRequestOption("id_gte", id)
}

// WithIDGT adds the id_gt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater than
// the provided id.
func WithIDGT(id string) RequestOption {
	return makeRequestOption("id_gt", id)
}

// WithIDLTE adds the id_lte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser or equal
// than the provided id.
func WithIDLTE(id string) RequestOption {
	return makeRequestOption("id_lte", id)
}

// WithIDLT adds the id_lt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser than the
// provided id.
func WithIDLT(id string) RequestOption {
	return makeRequestOption("id_lt", id)
}

// WithFilter adds the filter parameter to API calls, used when retrieving
// following feeds, allowing the check whether certain feeds are being followed.
func WithFilter(ids ...string) RequestOption {
	return makeRequestOption("filter", strings.Join(ids, ","))
}

// WithKeepHistory adds the `keep_history` parameter to API calls, used to keep
// history when unfollowing feeds, rather than purging it (default behavior).
// If the keepHistory parameter is false, nothing happens.
func WithKeepHistory(keepHistory bool) RequestOption {
	if !keepHistory {
		return requestOption{}
	}
	return makeRequestOption("keep_history", 1)
}

// WithActivityCopyLimitQuery sets the activity copy threshold for Follow Many
// API calls.
func WithActivityCopyLimitQuery(limit int) RequestOption {
	return makeRequestOption("activity_copy_limit", limit)
}

type followFeedOptions struct {
	Target            string `json:"target,omitempty"`
	ActivityCopyLimit int    `json:"activity_copy_limit,omitempty"`
}

// FollowFeedOption is a function used to customize FollowFeed API calls.
type FollowFeedOption func(*followFeedOptions)

// WithActivityCopyLimit sets the activity copy threshold for Follow Feed API
// calls.
func WithActivityCopyLimit(activityCopyLimit int) FollowFeedOption {
	return func(o *followFeedOptions) {
		o.ActivityCopyLimit = activityCopyLimit
	}
}

// TODO fix inheritance for request options on different types (with limit/offset <=> specific stuff)
