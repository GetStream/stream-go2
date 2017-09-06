package stream

import (
	"encoding/json"
	"time"
)

type response struct {
	Duration time.Duration `json:"duration,omitempty"`
	Next     string        `json:"next,omitempty"`
}

// FlatFeedResponse is the API response obtained when retrieving activities from
// a flat feed.
type FlatFeedResponse struct {
	response
	Results Activities `json:"results,omitempty"`
}

// UnmarshalJSON decodes the provided JSON payload into the FlatFeedResponse.
func (r *FlatFeedResponse) UnmarshalJSON(b []byte) error {
	type alias FlatFeedResponse
	aux := &struct {
		Duration string `json:"duration,omitempty"`
		*alias
	}{alias: (*alias)(r)}
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	r.Duration, err = time.ParseDuration(aux.Duration)
	if err != nil {
		return err
	}
	return nil
}

// AggregatedFeedResponse is the API response obtained when retrieving
// activities from an aggregated feed.
type AggregatedFeedResponse struct {
	response
	Results ActivityGroups `json:"results,omitempty"`
}

// UnmarshalJSON decodes the provided JSON payload into the
// AggregatedFeedResponse.
func (r *AggregatedFeedResponse) UnmarshalJSON(b []byte) error {
	type alias AggregatedFeedResponse
	aux := &struct {
		Duration string `json:"duration,omitempty"`
		*alias
	}{alias: (*alias)(r)}
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	r.Duration, err = time.ParseDuration(aux.Duration)
	if err != nil {
		return err
	}
	return nil
}

// AddActivitiesResponse is the API response obtained when adding activities to
// a feed.
type AddActivitiesResponse struct {
	response
	Activities []Activity `json:"activities,omitempty"`
}

// UnmarshalJSON decodes the provided JSON payload into the
// AddActivitiesResponse.
func (r *AddActivitiesResponse) UnmarshalJSON(b []byte) error {
	type alias AddActivitiesResponse
	aux := &struct {
		Duration string `json:"duration,omitempty"`
		*alias
	}{alias: (*alias)(r)}
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	r.Duration, err = time.ParseDuration(aux.Duration)
	if err != nil {
		return err
	}
	return nil
}

// Follower is the representation of a feed following another feed.
type Follower struct {
	FeedID   string `json:"feed_id,omitempty"`
	TargetID string `json:"target_id,omitempty"`
}

type followResponse struct {
	response
	Results []Follower `json:"results,omitempty"`
}

// FollowersResponse is the API response obtained when retrieving followers from
// a feed.
type FollowersResponse struct {
	followResponse
}

// FollowingResponse is the API response obtained when retrieving following
// feeds from a feed.
type FollowingResponse struct {
	followResponse
}

func (r *followResponse) UnmarshalJSON(b []byte) error {
	type alias followResponse
	aux := &struct {
		Duration string `json:"duration,omitempty"`
		*alias
	}{alias: (*alias)(r)}
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	r.Duration, err = time.ParseDuration(aux.Duration)
	if err != nil {
		return err
	}
	return nil
}

// AddToManyRequest is the API request body for adding an activity to multiple
// feeds at once.
type AddToManyRequest struct {
	Activity Activity `json:"activity,omitempty"`
	Feeds    []string `json:"feeds,omitempty"`
}

// FollowRelationship represents a follow relationship between a source
// ("follower") and a target ("following"), used for FollowMany requests.
type FollowRelationship struct {
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
}

// NewFollowRelationship is a helper for creating a FollowRelationship from the
// source ("follower") and target ("following") feeds.
func NewFollowRelationship(source, target Feed) FollowRelationship {
	return FollowRelationship{
		Source: source.ID(),
		Target: target.ID(),
	}
}
