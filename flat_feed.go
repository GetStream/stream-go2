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

// GetActivitiesWithRanking returns the activities (filtered) for the given FlatFeed,
// using the provided ranking method.
func (f *FlatFeed) GetActivitiesWithRanking(ranking string, opts ...GetActivitiesOption) (*FlatFeedResponse, error) {
	return f.GetActivities(append(opts, withRanking(ranking))...)
}
