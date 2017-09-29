package stream

import "fmt"

// Feed is a generic Stream feed, exporting the generic functions common to any
// Stream feed.
type Feed interface {
	ID() string
	Slug() string
	UserID() string
	AddActivity(Activity) (*AddActivityResponse, error)
	AddActivities(...Activity) (*AddActivitiesResponse, error)
	UpdateActivities(...Activity) error
	RemoveActivityByID(string) error
	RemoveActivityByForeignID(string) error
	Follow(*FlatFeed, ...FollowFeedOption) error
	GetFollowing(...FollowingOption) (*FollowingResponse, error)
	Unfollow(Feed, ...UnfollowOption) error
	ReplaceToTargets(Activity, []string) error
	UpdateToTargets(Activity, []string, []string) error
	Token(bool) string
}

type feed struct {
	slug   string
	userID string
	client *Client
}

// ID returns the feed ID, as slug:user_id.
func (f *feed) ID() string {
	return fmt.Sprintf("%s:%s", f.slug, f.userID)
}

// Slug returns the feed's slug.
func (f *feed) Slug() string {
	return f.slug
}

// UserID returns the feed's user_id.
func (f *feed) UserID() string {
	return f.userID
}

func newFeed(slug, userID string, client *Client) feed {
	return feed{userID: userID, slug: slug, client: client}
}

// AddActivity adds a new Activity to the feed.
func (f *feed) AddActivity(activity Activity) (*AddActivityResponse, error) {
	return f.client.addActivity(f, activity)
}

// AddActivities adds multiple activities to the feed.
func (f *feed) AddActivities(activities ...Activity) (*AddActivitiesResponse, error) {
	return f.client.addActivities(f, activities...)
}

// UpdateActivities updates existing activities in the feed.
func (f *feed) UpdateActivities(activities ...Activity) error {
	return f.client.updateActivities(activities...)
}

// RemoveActivityByID removes an activity from the feed (if present), using the provided
// id string argument as the ID field of the activity.
func (f *feed) RemoveActivityByID(id string) error {
	return f.client.removeActivityByID(f, id)
}

// RemoveActivityByID removes an activity from the feed (if present), using the provided
// foreignID string argument as the foreign_id field of the activity.
func (f *feed) RemoveActivityByForeignID(foreignID string) error {
	return f.client.removeActivityByForeignID(f, foreignID)
}

// Follow follows the provided feed (which must be a FlatFeed), applying the provided FollowFeedOptions,
// if any.
func (f *feed) Follow(feed *FlatFeed, opts ...FollowFeedOption) error {
	followOptions := &followFeedOptions{
		Target:            fmt.Sprintf("%s:%s", feed.Slug(), feed.UserID()),
		ActivityCopyLimit: defaultActivityCopyLimit,
	}
	for _, opt := range opts {
		opt(followOptions)
	}
	return f.client.follow(f, followOptions)
}

// GetFollowing returns the list of the feeds following the feed, applying the provided FollowingOptions,
// if any.
func (f *feed) GetFollowing(opts ...FollowingOption) (*FollowingResponse, error) {
	return f.client.getFollowing(f, opts...)
}

// Unfollow unfollows the provided feed, applying the provided UnfollowOptions, if any.
func (f *feed) Unfollow(target Feed, opts ...UnfollowOption) error {
	return f.client.unfollow(f, target.ID(), opts...)
}

// ReplaceToTargets removes all the existing "to" targets from the provided activity, replacing
// them with the ones in the provided new slice of feed IDs.
func (f *feed) ReplaceToTargets(activity Activity, new []string) error {
	return f.client.updateToTargets(f, activity, updateToTargetsWithNew(new...))
}

// UpdateToTargets updates the "to" targets for the provided activity, adding the feeds in the add slice of feed IDs
// on top of the existing ones, and removing the ones in the remove slice of feed IDs.
func (f *feed) UpdateToTargets(activity Activity, add []string, remove []string) error {
	return f.client.updateToTargets(f, activity, updateToTargetsWithAdd(add...), updateToTargetsWithRemove(remove...))
}

// Token returns a token that can be used client-side to listen in real-time to feed changes.
func (f *feed) Token(readonly bool) string {
	var action action
	if readonly {
		action = actionRead
	} else {
		action = actionWrite
	}
	token, err := f.client.authenticator.feedAuthToken(resFeed, action, f)
	if err != nil {
		return ""
	}
	return token
}
