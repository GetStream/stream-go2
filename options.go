package stream

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	defaultActivityCopyLimit = 300
)

// requestOption is an interface representing API request optional filters and
// parameters.
type requestOption interface {
	valuer
}

type valuer interface {
	values() (string, string)
	valid() bool
}

type baseRequestOption struct {
	key   string
	value string
}

func makeRequestOption(key string, value any) requestOption {
	return baseRequestOption{
		key:   key,
		value: fmt.Sprintf("%v", value),
	}
}

func (o baseRequestOption) values() (key, value string) {
	return o.key, o.value
}

func (o baseRequestOption) valid() bool {
	return o.value != ""
}

func withLimit(limit int) requestOption {
	return makeRequestOption("limit", limit)
}

func withOffset(offset int) requestOption {
	return makeRequestOption("offset", offset)
}

// GetActivitiesOption is an option usable by GetActivities methods for flat and aggregated feeds.
type GetActivitiesOption struct {
	requestOption
}

// WithActivitiesLimit adds the limit parameter to API calls which support it, limiting
// the number of results in the response to the provided limit threshold.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func WithActivitiesLimit(limit int) GetActivitiesOption {
	return GetActivitiesOption{withLimit(limit)}
}

// WithActivitiesOffset adds the offset parameter to API calls which support it, getting
// results starting from the provided offset index.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func WithActivitiesOffset(offset int) GetActivitiesOption {
	return GetActivitiesOption{withOffset(offset)}
}

// WithActivitiesIDGTE adds the id_gte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater or
// equal than the provided id.
func WithActivitiesIDGTE(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_gte", id)}
}

// WithActivitiesIDGT adds the id_gt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater than
// the provided id.
func WithActivitiesIDGT(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_gt", id)}
}

// WithActivitiesIDLTE adds the id_lte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser or equal
// than the provided id.
func WithActivitiesIDLTE(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_lte", id)}
}

// WithActivitiesIDLT adds the id_lt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser than the
// provided id.
func WithActivitiesIDLT(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_lt", id)}
}

func withActivitiesRanking(ranking string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("ranking", ranking)}
}

func WithRankingScoreVars() GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("withScoreVars", true)}
}

// Added a private struct `jsonString` to ensure WithExternalRankingVars accepts only encoded json string created by `MakeExternalVarJson` function
type jsonString struct {
	value string
}

func MakeExternalVarJson(externalRankingVars map[string]any) (jsonString, error) {
	str, err := json.Marshal(externalRankingVars)
	return jsonString{string(str)}, err
}

func WithExternalRankingVars(externalVarJson jsonString) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("ranking_vars", externalVarJson.value)}
}

// WithNotificationsMarkSeen marks as seen the given activity ids in a notification
// feed. If the all parameter is true, every activity in the feed is marked as seen.
func WithNotificationsMarkSeen(all bool, activityIDs ...string) GetActivitiesOption {
	if all {
		return GetActivitiesOption{makeRequestOption("mark_seen", true)}
	}
	return GetActivitiesOption{makeRequestOption("mark_seen", strings.Join(activityIDs, ","))}
}

// WithNotificationsMarkRead marks as read the given activity ids in a notification
// feed. If the all parameter is true, every activity in the feed is marked as read.
func WithNotificationsMarkRead(all bool, activityIDs ...string) GetActivitiesOption {
	if all {
		return GetActivitiesOption{makeRequestOption("mark_read", true)}
	}
	return GetActivitiesOption{makeRequestOption("mark_read", strings.Join(activityIDs, ","))}
}

// WithCustomParam adds a custom parameter to the read request.
func WithCustomParam(name, value string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption(name, value)}
}

// WithEnrichOwnReactions enriches the activities with the reactions to them.
func WithEnrichOwnReactions() GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("withOwnReactions", true)}
}

func WithEnrichUserReactions(userID string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("user_id", userID)}
}

// WithEnrichRecentReactions enriches the activities with the first reactions to them.
func WithEnrichFirstReactions() GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("withFirstReactions", true)}
}

// WithEnrichRecentReactions enriches the activities with the recent reactions to them.
func WithEnrichRecentReactions() GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("withRecentReactions", true)}
}

// WithEnrichReactionCounts enriches the activities with the reaction counts.
func WithEnrichReactionCounts() GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("withReactionCounts", true)}
}

// WithEnrichOwnChildren enriches the activities with the children reactions.
func WithEnrichOwnChildren() GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("withOwnChildren", true)}
}

// WithEnrichRecentReactionsLimit specifies how many recent reactions to include in the enrichment.
func WithEnrichRecentReactionsLimit(limit int) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("recentReactionsLimit", limit)}
}

// WithEnrichReactionsLimit specifies how many reactions to include in the enrichment.
func WithEnrichReactionsLimit(limit int) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("reaction_limit", limit)}
}

// WithEnrichReactionKindsFilter filters the reactions by the specified kinds
func WithEnrichReactionKindsFilter(kinds ...string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("reactionKindsFilter", strings.Join(kinds, ","))}
}

// WithEnrichOwnChildrenKindsFilter filters the reactions by the specified kinds for own children
func WithEnrichOwnChildrenKindsFilter(kinds ...string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("withOwnChildrenKinds", strings.Join(kinds, ","))}
}

// FollowingOption is an option usable by following feed methods.
type FollowingOption struct {
	requestOption
}

// WithFollowingFilter adds the filter parameter to API calls, used when retrieving
// following feeds, allowing the check whether certain feeds are being followed.
func WithFollowingFilter(ids ...string) FollowingOption {
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

// FollowersOption is an option usable by followers feed methods.
type FollowersOption struct {
	requestOption
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
	requestOption
}

// WithUnfollowKeepHistory adds the keep_history parameter to API calls, used to keep
// history when unfollowing feeds, rather than purging it (default behavior).
// If the keepHistory parameter is false, nothing happens.
func WithUnfollowKeepHistory(keepHistory bool) UnfollowOption {
	if !keepHistory {
		return UnfollowOption{nop{}}
	}
	return UnfollowOption{makeRequestOption("keep_history", 1)}
}

type followFeedOptions struct {
	Target            string `json:"target,omitempty"`
	ActivityCopyLimit int    `json:"activity_copy_limit"`
}

// FollowManyOption is an option to customize behavior of Follow Many calls.
type FollowManyOption struct {
	requestOption
}

// WithFollowManyActivityCopyLimit sets how many activities should be copied from the target feed.
func WithFollowManyActivityCopyLimit(activityCopyLimit int) FollowManyOption {
	return FollowManyOption{makeRequestOption("activity_copy_limit", activityCopyLimit)}
}

// FollowFeedOption is a function used to customize FollowFeed API calls.
type FollowFeedOption func(*followFeedOptions)

// WithFollowFeedActivityCopyLimit sets the activity copy threshold for Follow Feed API
// calls.
func WithFollowFeedActivityCopyLimit(activityCopyLimit int) FollowFeedOption {
	return func(o *followFeedOptions) {
		o.ActivityCopyLimit = activityCopyLimit
	}
}

// FollowStatOption is an option used to customize FollowStats API calls.
type FollowStatOption struct {
	requestOption
}

// WithFollowerSlugs sets the follower feed slugs for filtering in counting.
func WithFollowerSlugs(slugs ...string) FollowStatOption {
	return FollowStatOption{makeRequestOption("followers_slugs", strings.Join(slugs, ","))}
}

// WithFollowerSlugs sets the following feed slugs for filtering in counting.
func WithFollowingSlugs(slugs ...string) FollowStatOption {
	return FollowStatOption{makeRequestOption("following_slugs", strings.Join(slugs, ","))}
}

// UpdateToTargetsOption determines what operations perform during an UpdateToTargets API call.
type UpdateToTargetsOption func(*updateToTargetsRequest)

// WithToTargetsNew sets the new to targets, replacing all the existing ones. It cannot be used in combination with any other UpdateToTargetsOption.
func WithToTargetsNew(targets ...string) UpdateToTargetsOption {
	return func(r *updateToTargetsRequest) {
		r.New = targets
	}
}

// WithToTargetsAdd sets the add to targets, adding them to the activity's existing ones.
func WithToTargetsAdd(targets ...string) UpdateToTargetsOption {
	return func(r *updateToTargetsRequest) {
		r.Adds = targets
	}
}

// WithToTargetsRemove sets the remove to targets, removing them from activity's the existing ones.
func WithToTargetsRemove(targets ...string) UpdateToTargetsOption {
	return func(r *updateToTargetsRequest) {
		r.Removes = targets
	}
}

// AddObjectOption is an option usable by the Collections.Add method.
type AddObjectOption func(*addCollectionRequest)

// WithUserID adds the user id to the Collections.Add request object.
func WithUserID(userID string) AddObjectOption {
	return func(req *addCollectionRequest) {
		req.UserID = &userID
	}
}

// FilterReactionsOption is an option used by Reactions.Filter() to support pagination.
type FilterReactionsOption struct {
	requestOption
}

// WithLimit adds the limit parameter to the Reactions.Filter() call.
func WithLimit(limit int) FilterReactionsOption {
	return FilterReactionsOption{withLimit(limit)}
}

// WithIDGTE adds the id_gte parameter to API calls, used when retrieving
// paginated reactions, returning activities with ID greater or
// equal than the provided id.
func WithIDGTE(id string) FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("id_gte", id)}
}

// WithIDGT adds the id_gt parameter to API calls, used when retrieving
// paginated reactions.
func WithIDGT(id string) FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("id_gt", id)}
}

// WithIDLTE adds the id_lte parameter to API calls, used when retrieving
// paginated reactions.
func WithIDLTE(id string) FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("id_lte", id)}
}

// WithIDLT adds the id_lt parameter to API calls, used when retrieving
// paginated reactions.
func WithIDLT(id string) FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("id_lt", id)}
}

// WithActivityData will enable returning the activity data when filtering
// reactions by activity_id.
func WithActivityData() FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("with_activity_data", true)}
}

// WithOwnChildren will enable returning the children reactions when filtering
// reactions by parent ID.
func WithOwnChildren() FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("with_own_children", true)}
}

// WithOwnUserID will enable further filtering by the given user id.
// It's similar to FilterReactionsAttribute user id.
func WithOwnUserID(userID string) FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("user_id", userID)}
}

// WithChildrenUserID will enable further filtering own children by the given user id.
// It's different than FilterReactionsAttribute user id.
func WithChildrenUserID(userID string) FilterReactionsOption {
	return FilterReactionsOption{makeRequestOption("children_user_id", userID)}
}

// FilterReactionsAttribute specifies the filtering method of Reactions.Filter()
type FilterReactionsAttribute func() string

// ByKind filters reactions by kind, after the initial desired filtering method was applied.
func (a FilterReactionsAttribute) ByKind(kind string) FilterReactionsAttribute {
	return func() string {
		base := a()
		return fmt.Sprintf("%s/%s", base, kind)
	}
}

// ByActivityID will filter reactions based on the specified activity id.
func ByActivityID(activityID string) FilterReactionsAttribute {
	return func() string {
		return fmt.Sprintf("activity_id/%s", activityID)
	}
}

// ByReactionID will filter reactions based on the specified parent reaction id.
func ByReactionID(reactionID string) FilterReactionsAttribute {
	return func() string {
		return fmt.Sprintf("reaction_id/%s", reactionID)
	}
}

// ByUserID will filter reactions based on the specified user id.
func ByUserID(userID string) FilterReactionsAttribute {
	return func() string {
		return fmt.Sprintf("user_id/%s", userID)
	}
}

type nop struct{}

func (nop) values() (key, value string) {
	return "", ""
}

func (nop) valid() bool {
	return false
}
