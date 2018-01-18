package stream

import "encoding/json"

// FlatFeed is a Stream flat feed.
type FlatFeed struct {
	feed
}

// GetActivities returns the activities for the given FlatFeed, filtering
// results with the provided GetActivitiesOption parameters.
func (f *FlatFeed) GetActivities(opts ...GetActivitiesOption) (*FlatFeedResponse, error) {
	body, err := f.client.getActivities(f, opts...)
	if err != nil {
		return nil, err
	}
	var resp FlatFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNextPageActivities returns the activities for the given FlatFeed at the "next" page
// of a previous *FlatFeedResponse response, if any.
func (f *FlatFeed) GetNextPageActivities(resp *FlatFeedResponse) (*FlatFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetActivities(opts...)
}

// GetActivitiesWithRanking returns the activities (filtered) for the given FlatFeed,
// using the provided ranking method.
func (f *FlatFeed) GetActivitiesWithRanking(ranking string, opts ...GetActivitiesOption) (*FlatFeedResponse, error) {
	return f.GetActivities(append(opts, withActivitiesRanking(ranking))...)
}

// GetFollowers returns the feeds following the given FlatFeed.
func (f *FlatFeed) GetFollowers(opts ...FollowersOption) (*FollowersResponse, error) {
	return f.client.getFollowers(f, opts...)
}
