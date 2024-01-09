package stream

import (
	"context"
	"encoding/json"
)

// AggregatedFeed is a Stream aggregated feed, which contains activities grouped
// based on the grouping function defined on the dashboard.
type AggregatedFeed struct {
	feed
}

// GetActivities requests and retrieves the activities and groups for the
// aggregated feed.
func (f *AggregatedFeed) GetActivities(ctx context.Context, opts ...GetActivitiesOption) (*AggregatedFeedResponse, error) {
	body, err := f.client.getActivities(ctx, f, opts...)
	if err != nil {
		return nil, err
	}
	var resp AggregatedFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetActivitiesWithRanking returns the activities and groups for the given AggregatedFeed,
// using the provided ranking method.
func (f *AggregatedFeed) GetActivitiesWithRanking(ctx context.Context, ranking string, opts ...GetActivitiesOption) (*AggregatedFeedResponse, error) {
	return f.GetActivities(ctx, append(opts, WithActivitiesRanking(ranking))...)
}

// GetNextPageActivities returns the activities for the given AggregatedFeed at the "next" page
// of a previous *AggregatedFeedResponse response, if any.
func (f *AggregatedFeed) GetNextPageActivities(ctx context.Context, resp *AggregatedFeedResponse) (*AggregatedFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetActivities(ctx, opts...)
}

// GetEnrichedActivities requests and retrieves the enriched activities and groups for the
// aggregated feed.
func (f *AggregatedFeed) GetEnrichedActivities(ctx context.Context, opts ...GetActivitiesOption) (*EnrichedAggregatedFeedResponse, error) {
	body, err := f.client.getEnrichedActivities(ctx, f, opts...)
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
func (f *AggregatedFeed) GetNextPageEnrichedActivities(ctx context.Context, resp *EnrichedAggregatedFeedResponse) (*EnrichedAggregatedFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetEnrichedActivities(ctx, opts...)
}

// GetEnrichedActivitiesWithRanking returns the enriched activities and groups for the given AggregatedFeed,
// using the provided ranking method.
func (f *AggregatedFeed) GetEnrichedActivitiesWithRanking(ctx context.Context, ranking string, opts ...GetActivitiesOption) (*EnrichedAggregatedFeedResponse, error) {
	return f.GetEnrichedActivities(ctx, append(opts, WithActivitiesRanking(ranking))...)
}
