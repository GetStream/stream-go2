package stream

import (
	"context"
	"fmt"
	"regexp"
)

var (
	_ Feed = (*FlatFeed)(nil)
	_ Feed = (*AggregatedFeed)(nil)
	_ Feed = (*NotificationFeed)(nil)
)

// Feed is a generic Stream feed, exporting the generic functions common to any
// Stream feed.
type Feed interface {
	ID() string
	Slug() string
	UserID() string
	AddActivity(context.Context, Activity) (*AddActivityResponse, error)
	AddActivities(context.Context, ...Activity) (*AddActivitiesResponse, error)
	RemoveActivityByID(context.Context, string) (*RemoveActivityResponse, error)
	RemoveActivityByForeignID(context.Context, string) (*RemoveActivityResponse, error)
	Follow(context.Context, *FlatFeed, ...FollowFeedOption) (*BaseResponse, error)
	GetFollowing(context.Context, ...FollowingOption) (*FollowingResponse, error)
	Unfollow(context.Context, Feed, ...UnfollowOption) (*BaseResponse, error)
	UpdateToTargets(context.Context, Activity, ...UpdateToTargetsOption) (*UpdateToTargetsResponse, error)
	BatchUpdateToTargets(context.Context, []UpdateToTargetsRequest) (*UpdateToTargetsResponse, error)
	RealtimeToken(bool) string
}

const feedSlugIDSeparator = ":"

var userIDRegex *regexp.Regexp

func init() {
	userIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
}

type feed struct {
	client *Client
	slug   string
	userID string
}

// ID returns the feed ID, as slug:user_id.
func (f *feed) ID() string {
	return fmt.Sprintf("%s%s%s", f.slug, feedSlugIDSeparator, f.userID)
}

// Slug returns the feed's slug.
func (f *feed) Slug() string {
	return f.slug
}

// UserID returns the feed's user_id.
func (f *feed) UserID() string {
	return f.userID
}

func newFeed(slug, userID string, client *Client) (*feed, error) {
	ok := userIDRegex.MatchString(userID)
	if !ok {
		return nil, errInvalidUserID
	}
	return &feed{userID: userID, slug: slug, client: client}, nil
}

// AddActivity adds a new Activity to the feed.
func (f *feed) AddActivity(ctx context.Context, activity Activity) (*AddActivityResponse, error) {
	return f.client.addActivity(ctx, f, activity)
}

// AddActivities adds multiple activities to the feed.
func (f *feed) AddActivities(ctx context.Context, activities ...Activity) (*AddActivitiesResponse, error) {
	return f.client.addActivities(ctx, f, activities...)
}

// RemoveActivityByID removes an activity from the feed (if present), using the provided
// id string argument as the ID field of the activity.
func (f *feed) RemoveActivityByID(ctx context.Context, id string) (*RemoveActivityResponse, error) {
	return f.client.removeActivityByID(ctx, f, id)
}

// RemoveActivityByForeignID removes an activity from the feed (if present), using the provided
// foreignID string argument as the foreign_id field of the activity.
func (f *feed) RemoveActivityByForeignID(ctx context.Context, foreignID string) (*RemoveActivityResponse, error) {
	return f.client.removeActivityByForeignID(ctx, f, foreignID)
}

// Follow follows the provided feed (which must be a FlatFeed), applying the provided FollowFeedOptions,
// if any.
func (f *feed) Follow(ctx context.Context, feed *FlatFeed, opts ...FollowFeedOption) (*BaseResponse, error) {
	followOptions := &followFeedOptions{
		Target:            fmt.Sprintf("%s:%s", feed.Slug(), feed.UserID()),
		ActivityCopyLimit: defaultActivityCopyLimit,
	}
	for _, opt := range opts {
		opt(followOptions)
	}
	return f.client.follow(ctx, f, followOptions)
}

// GetFollowing returns the list of the feeds following the feed, applying the provided FollowingOptions,
// if any.
func (f *feed) GetFollowing(ctx context.Context, opts ...FollowingOption) (*FollowingResponse, error) {
	return f.client.getFollowing(ctx, f, opts...)
}

// Unfollow unfollows the provided feed, applying the provided UnfollowOptions, if any.
func (f *feed) Unfollow(ctx context.Context, target Feed, opts ...UnfollowOption) (*BaseResponse, error) {
	return f.client.unfollow(ctx, f, target.ID(), opts...)
}

// UpdateToTargets updates the "to" targets for the provided activity, with the options passed
// as argument for replacing, adding, or removing to targets.
func (f *feed) UpdateToTargets(ctx context.Context, activity Activity, opts ...UpdateToTargetsOption) (*UpdateToTargetsResponse, error) {
	return f.client.updateToTargets(ctx, f, activity, opts...)
}

// BatchUpdateToTargets updates the "to" targets for up to 100 activities, with the options passed
// as argument for replacing, adding, or removing to targets.
// NOTE: Only the first update is executed synchronously (same response as UpdateToTargets()), the remaining N-1 updates will be put in a worker queue and executed asynchronously.
func (f *feed) BatchUpdateToTargets(ctx context.Context, reqs []UpdateToTargetsRequest) (*UpdateToTargetsResponse, error) {
	return f.client.batchUpdateToTargets(ctx, f, reqs)
}

// RealtimeToken returns a token that can be used client-side to listen in real-time to feed changes.
func (f *feed) RealtimeToken(readonly bool) string {
	var action action
	if readonly {
		action = actionRead
	} else {
		action = actionWrite
	}
	id := f.client.authenticator.feedID(f)
	claims := f.client.authenticator.jwtFeedClaims(resFeed, action, id)
	token, err := f.client.authenticator.jwtSignatureFromClaims(claims)
	if err != nil {
		return ""
	}
	return token
}
