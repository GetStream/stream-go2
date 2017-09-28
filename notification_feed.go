package stream

import "encoding/json"

type NotificationFeed struct {
	feed
}

// GetActivities returns the activities for the given NotificationFeed, filtering
// results with the provided GetActivitiesOption parameters.
func (f *NotificationFeed) GetActivities(opts ...GetActivitiesOption) (*NotificationFeedResponse, error) {
	body, err := f.client.getActivities(f, opts...)
	if err != nil {
		return nil, err
	}
	var resp NotificationFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
