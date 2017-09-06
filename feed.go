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
	GetFollowers(...RequestOption) (*FollowersResponse, error)
	GetFollowing(...RequestOption) (*FollowingResponse, error) // TODO test filter param
	Unfollow(Feed, ...RequestOption) error                     // TODO test heep_history param
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
	return f.client.addActivities(f.Slug(), f.UserID(), activities...)
}

func (f feed) UpdateActivities(activities ...Activity) error {
	return f.client.updateActivities(activities...)
}

func (f feed) RemoveActivityByID(id string) error {
	return f.client.removeActivityByID(f.Slug(), f.UserID(), id)
}

func (f feed) RemoveActivityByForeignID(foreignID string) error {
	return f.client.removeActivityByForeignID(f.Slug(), f.UserID(), foreignID)
}

func (f feed) Follow(feed *FlatFeed, opts ...FollowFeedOption) error {
	followOptions := &followFeedOptions{
		Target:            fmt.Sprintf("%s:%s", feed.Slug(), f.UserID()),
		ActivityCopyLimit: defaultActivityCopyLimit,
	}
	for _, opt := range opts {
		opt(followOptions)
	}
	return f.client.follow(f.Slug(), f.UserID(), followOptions)
}

func (f feed) GetFollowers(opts ...RequestOption) (*FollowersResponse, error) {
	return f.client.getFollowers(f.Slug(), f.UserID(), opts...)
}

func (f feed) GetFollowing(opts ...RequestOption) (*FollowingResponse, error) {
	return f.client.getFollowing(f.Slug(), f.UserID(), opts...)
}

func (f feed) Unfollow(target Feed, opts ...RequestOption) error {
	return f.client.unfollow(f.Slug(), f.UserID(), target.ID(), opts...)
}
