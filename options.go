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

func withLimit(limit int) RequestOption {
	return makeRequestOption("limit", limit)
}

func withOffset(offset int) RequestOption {
	return makeRequestOption("offset", offset)
}

// GetActivitiesOption is an option usable by GetActivities methods for flat and aggregated feeds.
type GetActivitiesOption struct {
	RequestOption
}

func withRanking(ranking string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("ranking", ranking)}
}

// WithLimit adds the limit parameter to API calls which support it, limiting
// the number of results in the response to the provided limit threshold.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func WithLimit(limit int) GetActivitiesOption {
	return GetActivitiesOption{withLimit(limit)}
}

// WithOffset adds the offset parameter to API calls which support it, getting
// results starting from the provided offset index.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func WithOffset(offset int) GetActivitiesOption {
	return GetActivitiesOption{withOffset(offset)}
}

// WithIDGTE adds the id_gte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater or
// equal than the provided id.
func WithIDGTE(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_gte", id)}
}

// WithIDGT adds the id_gt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater than
// the provided id.
func WithIDGT(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_gt", id)}
}

// WithIDLTE adds the id_lte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser or equal
// than the provided id.
func WithIDLTE(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_lte", id)}
}

// WithIDLT adds the id_lt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser than the
// provided id.
func WithIDLT(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_lt", id)}
}

// FollowingOption is an option usable by following feed methods.
type FollowingOption struct {
	RequestOption
}

// WithFilter adds the filter parameter to API calls, used when retrieving
// following feeds, allowing the check whether certain feeds are being followed.
func WithFilter(ids ...string) FollowingOption {
	return FollowingOption{makeRequestOption("filter", strings.Join(ids, ","))}
}

// WithFollowingLimit limits the number of followings in the response to the provided limit.
func WithFollowingLimit(limit int) FollowingOption {
	return FollowingOption{withLimit(limit)}
}

// WithFollowingOffset returns followings starting from the given offset.
func WithFollowingOffset(offset int) FollowingOption {
	return FollowingOption{withOffset(offset)}
}

// WithActivityCopyLimitQuery sets the activity copy threshold for Follow Many
// API calls.
func WithActivityCopyLimitQuery(limit int) FollowingOption {
	return FollowingOption{makeRequestOption("activity_copy_limit", limit)}
}

// FollowersOption is an option usable by followers feed methods.
type FollowersOption struct {
	RequestOption
}

// WithFollowersLimit limits the number of followers in the response to the provided limit.
func WithFollowersLimit(limit int) FollowersOption {
	return FollowersOption{withLimit(limit)}
}

// WithFollowersOffset returns followers starting from the given offset.
func WithFollowersOffset(offset int) FollowersOption {
	return FollowersOption{withOffset(offset)}
}

// UnfollowOption is an option usable with the Unfollow feed method.
type UnfollowOption struct {
	RequestOption
}

// WithKeepHistory adds the `keep_history` parameter to API calls, used to keep
// history when unfollowing feeds, rather than purging it (default behavior).
// If the keepHistory parameter is false, nothing happens.
func WithKeepHistory(keepHistory bool) UnfollowOption {
	if !keepHistory {
		return UnfollowOption{}
	}
	return UnfollowOption{makeRequestOption("keep_history", 1)}
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

// UpdateToTargetsOption is used for configuring an UpdateToTargets API call.
type UpdateToTargetsOption func(*UpdateToTargetsRequest)

// NewToTargets sets the new to targets for an UpdateToTargets request. They will replace
// all previous to targets.
func NewToTargets(targets []Feed) UpdateToTargetsOption {
	return func(r *UpdateToTargetsRequest) {
		r.New = make([]string, len(targets))
		for i := range targets {
			r.New[i] = targets[i].ID()
		}
	}
}

// AddToTargets sets the to targets additions for an UpdateToTargets request. They
// will be added to the existing to targets.
func AddToTargets(targets []Feed) UpdateToTargetsOption {
	return func(r *UpdateToTargetsRequest) {
		r.Adds = make([]string, len(targets))
		for i := range targets {
			r.Adds[i] = targets[i].ID()
		}
	}
}

// RemoveToTargets sets the targets to be removed for an UpdateToTargets request.
func RemoveToTargets(targets []Feed) UpdateToTargetsOption {
	return func(r *UpdateToTargetsRequest) {
		r.Removes = make([]string, len(targets))
		for i := range targets {
			r.Removes[i] = targets[i].ID()
		}
	}
}
