package stream

import (
	"context"
	"encoding/json"
)

// FlatFeed is a Stream flat feed.
type FlatFeed struct {
	feed
}

// GetActivities returns the activities for the given FlatFeed, filtering
// results with the provided GetActivitiesOption parameters.
func (f *FlatFeed) GetActivities(ctx context.Context, opts ...GetActivitiesOption) (*FlatFeedResponse, error) {
	body, err := f.client.getActivities(ctx, f, opts...)
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
func (f *FlatFeed) GetNextPageActivities(ctx context.Context, resp *FlatFeedResponse) (*FlatFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetActivities(ctx, opts...)
}

// GetActivitiesWithRanking returns the activities (filtered) for the given FlatFeed,
// using the provided ranking method.
func (f *FlatFeed) GetActivitiesWithRanking(ctx context.Context, ranking string, opts ...GetActivitiesOption) (*FlatFeedResponse, error) {
	return f.GetActivities(ctx, append(opts, withActivitiesRanking(ranking))...)
}

// GetFollowers returns the feeds following the given FlatFeed.
func (f *FlatFeed) GetFollowers(ctx context.Context, opts ...FollowersOption) (*FollowersResponse, error) {
	return f.client.getFollowers(ctx, f, opts...)
}

// GetEnrichedActivities returns the enriched activities for the given FlatFeed, filtering
// results with the provided GetActivitiesOption parameters.
func (f *FlatFeed) GetEnrichedActivities(ctx context.Context, opts ...GetActivitiesOption) (*EnrichedFlatFeedResponse, error) {
	body, err := f.client.getEnrichedActivities(ctx, f, opts...)
	if err != nil {
		return nil, err
	}
	var resp EnrichedFlatFeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNextPageEnrichedActivities returns the enriched activities for the given FlatFeed at the "next" page
// of a previous *EnrichedFlatFeedResponse response, if any.
func (f *FlatFeed) GetNextPageEnrichedActivities(ctx context.Context, resp *EnrichedFlatFeedResponse) (*EnrichedFlatFeedResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return f.GetEnrichedActivities(ctx, opts...)
}

// GetEnrichedActivitiesWithRanking returns the enriched activities (filtered) for the given FlatFeed,
// using the provided ranking method.
func (f *FlatFeed) GetEnrichedActivitiesWithRanking(ctx context.Context, ranking string, opts ...GetActivitiesOption) (*EnrichedFlatFeedResponse, error) {
	return f.GetEnrichedActivities(ctx, append(opts, withActivitiesRanking(ranking))...)
}

// FollowStats returns the follower/following counts of the feed.
// If options are given, counts are filtered for the given slugs.
// Counts will be capped at 10K, if higher counts are needed and contact to support.
func (f *FlatFeed) FollowStats(ctx context.Context, opts ...FollowStatOption) (*FollowStatResponse, error) {
	return f.client.followStats(ctx, f, opts...)
}
