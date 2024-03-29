package stream

import (
	"context"
	"encoding/json"
)

// NotificationFeed is a Stream notification feed.
type NotificationFeed struct {
	feed
}

// GetActivities returns the activities for the given NotificationFeed, filtering
// results with the provided GetActivitiesOption parameters.
func (f *NotificationFeed) GetActivities(ctx context.Context, opts ...GetActivitiesOption) (*NotificationFeedResponse, error) {
	body, err := f.client.getActivities(ctx, f, opts...)
	if err != nil {
		return nil, err
	}
	var resp NotificationFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNextPageActivities returns the activities for the given NotificationFeed at the "next" page
// of a previous *NotificationFeedResponse response, if any.
func (f *NotificationFeed) GetNextPageActivities(ctx context.Context, resp *NotificationFeedResponse) (*NotificationFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetActivities(ctx, opts...)
}

// GetEnrichedActivities returns the enriched activities for the given NotificationFeed, filtering
// results with the provided GetActivitiesOption parameters.
func (f *NotificationFeed) GetEnrichedActivities(ctx context.Context, opts ...GetActivitiesOption) (*EnrichedNotificationFeedResponse, error) {
	body, err := f.client.getEnrichedActivities(ctx, f, opts...)
	if err != nil {
		return nil, err
	}
	var resp EnrichedNotificationFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNextPageEnrichedActivities returns the enriched activities for the given NotificationFeed at the "next" page
// of a previous *EnrichedNotificationFeedResponse response, if any.
func (f *NotificationFeed) GetNextPageEnrichedActivities(ctx context.Context, resp *EnrichedNotificationFeedResponse) (*EnrichedNotificationFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetEnrichedActivities(ctx, opts...)
}
