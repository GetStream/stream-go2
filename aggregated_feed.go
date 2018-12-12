package stream

import (
	"encoding/json"
)

// AggregatedFeed is a Stream aggregated feed, which contains activities grouped
// based on the grouping function defined on the dashboard.
type AggregatedFeed struct {
	feed
}

// GetActivities requests and retrieves the activities and groups for the
// aggregated feed.
func (f *AggregatedFeed) GetActivities(opts ...GetActivitiesOption) (*AggregatedFeedResponse, error) {
	body, err := f.client.getActivities(f, opts...)
	if err != nil {
		return nil, err
	}
	var resp AggregatedFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNextPageActivities returns the activities for the given AggregatedFeed at the "next" page
// of a previous *AggregatedFeedResponse response, if any.
func (f *AggregatedFeed) GetNextPageActivities(resp *AggregatedFeedResponse) (*AggregatedFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetActivities(opts...)
}

// GetEnrichedActivities requests and retrieves the enriched activities and groups for the
// aggregated feed.
func (f *AggregatedFeed) GetEnrichedActivities(opts ...GetActivitiesOption) (*EnrichedAggregatedFeedResponse, error) {
	body, err := f.client.getEnrichedActivities(f, opts...)
	if err != nil {
		return nil, err
	}
	var resp EnrichedAggregatedFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNextPageEnrichedActivities returns the enriched activities for the given AggregatedFeed at the "next" page
// of a previous *EnrichedAggregatedFeedResponse response, if any.
func (f *AggregatedFeed) GetNextPageEnrichedActivities(resp *EnrichedAggregatedFeedResponse) (*EnrichedAggregatedFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetEnrichedActivities(opts...)
}
