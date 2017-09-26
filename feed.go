package stream

import "fmt"

// Feed is a generic Stream feed, exporting the generic functions common to any
// Stream feed.
type Feed interface {
	ID() string
	Slug() string
	UserID() string
	AddActivities(...Activity) (*AddActivitiesResponse, error)
	UpdateActivities(...Activity) error
	RemoveActivityByID(string) error
	RemoveActivityByForeignID(string) error
	Follow(*FlatFeed, ...FollowFeedOption) error
	GetFollowers(...FollowersOption) (*FollowersResponse, error)
	GetFollowing(...FollowingOption) (*FollowingResponse, error) // TODO test filter param
	Unfollow(Feed, ...UnfollowOption) error                      // TODO test heep_history param
	UpdateToTargets(Activity, ...UpdateToTargetsOption) error
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

func (f feed) GetFollowers(opts ...FollowersOption) (*FollowersResponse, error) {
	return f.client.getFollowers(f, opts...)
}

func (f feed) GetFollowing(opts ...FollowingOption) (*FollowingResponse, error) {
	return f.client.getFollowing(f, opts...)
}

func (f feed) Unfollow(target Feed, opts ...UnfollowOption) error {
	return f.client.unfollow(f, target.ID(), opts...)
}

func (f feed) UpdateToTargets(activity Activity, opts ...UpdateToTargetsOption) error {
	return f.client.updateToTargets(f, activity, opts...)
}
