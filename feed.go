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
	ReplaceToTargets(Activity, []Feed) error
	UpdateToTargets(Activity, []Feed, []Feed) error
}

type feed struct {
	slug   string
	userID string
	client *Client
}

func (f feed) ID() string     { return fmt.Sprintf("%s:%s", f.slug, f.userID) }
func (f feed) Slug() string   { return f.slug }
func (f feed) UserID() string { return f.userID }

func newFeed(slug, userID string, client *Client) feed {
	return feed{userID: userID, slug: slug, client: client}
}

func (f feed) AddActivity(activity Activity) (*AddActivityResponse, error) {
	return f.client.addActivity(f, activity)
}

func (f feed) AddActivities(activities ...Activity) (*AddActivitiesResponse, error) {
	return f.client.addActivities(f, activities...)
}

func (f feed) UpdateActivities(activities ...Activity) error {
	return f.client.updateActivities(activities...)
}

func (f feed) RemoveActivityByID(id string) error {
	return f.client.removeActivityByID(f, id)
}

func (f feed) RemoveActivityByForeignID(foreignID string) error {
	return f.client.removeActivityByForeignID(f, foreignID)
}

func (f feed) Follow(feed *FlatFeed, opts ...FollowFeedOption) error {
	followOptions := &followFeedOptions{
		Target:            fmt.Sprintf("%s:%s", feed.Slug(), f.UserID()),
		ActivityCopyLimit: defaultActivityCopyLimit,
	}
	for _, opt := range opts {
		opt(followOptions)
	}
	return f.client.follow(f, followOptions)
}

func (f feed) GetFollowing(opts ...FollowingOption) (*FollowingResponse, error) {
	return f.client.getFollowing(f, opts...)
}

func (f feed) Unfollow(target Feed, opts ...UnfollowOption) error {
	return f.client.unfollow(f, target.ID(), opts...)
}

func (f feed) ReplaceToTargets(activity Activity, new []Feed) error {
	return f.client.updateToTargets(f, activity, updateToTargetsWithNew(new...))
}

func (f feed) UpdateToTargets(activity Activity, add []Feed, remove []Feed) error {
	return f.client.updateToTargets(f, activity, updateToTargetsWithAdd(add...), updateToTargetsWithRemove(remove...))
}
