package stream

import "encoding/json"

// FlatFeed is a Stream flat feed.
type FlatFeed struct {
	feed
}

// GetActivities returns the activities for the given FlatFeed, filtering
// results with the provided RequestOption parameters.
func (f *FlatFeed) GetActivities(opts ...RequestOption) (*FlatFeedResponse, error) {
	body, err := f.client.getActivities(f.Slug(), f.UserID(), opts...)
	if err != nil {
		return nil, err
	}
	var resp FlatFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
