package stream

import (
	"encoding/json"
	"strings"
	"time"
)

// Duration wraps time.Duration, used because of JSON marshaling and
// unmarshaling.
type Duration struct {
	time.Duration
}

// UnmarshalJSON for Duration is required because of the incoming duration string.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var err error
	*d, err = durationFromString(strings.Replace(string(b), `"`, "", -1))
	return err
}

// MarshalJSON marshals the Duration to a string like "30s".
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func durationFromString(s string) (Duration, error) {
	dd, err := time.ParseDuration(s)
	return Duration{dd}, err
}

// Time wraps time.Time, used because of custom API time format in JSON marshaling
// and unmarshaling.
type Time struct {
	time.Time
}

// UnmarshalJSON for Time is required because of the incoming time string format.
func (t *Time) UnmarshalJSON(b []byte) error {
	var err error
	*t, err = timeFromString(strings.Replace(string(b), `"`, "", -1))
	return err
}

// MarshalJSON marshals Time into a string formatted with the TimeLayout format.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(TimeLayout))
}

func timeFromString(s string) (Time, error) {
	tt, err := time.Parse(TimeLayout, s)
	return Time{tt}, err
}

// Response is the part of StreamAPI responses common throughout the API.
type response struct {
	Duration Duration `json:"duration,omitempty"`
}

// readResponse is the part of StreamAPI responses common for GetActivities API requests.
type readResponse struct {
	response
	Next string `json:"next,omitempty"`
}

// FlatFeedResponse is the API response obtained when retrieving activities from
// a flat feed.
type FlatFeedResponse struct {
	readResponse
	Results []Activity `json:"results,omitempty"`
}

// AggregatedFeedResponse is the API response obtained when retrieving
// activities from an aggregated feed.
type AggregatedFeedResponse struct {
	readResponse
	Results []ActivityGroup `json:"results,omitempty"`
}

// NotificationFeedResponse is the API response obtained when retrieving activities
// from a notification feed.
type NotificationFeedResponse struct {
	readResponse
	Unseen  int                      `json:"unseen"`
	Unread  int                      `json:"unread"`
	Results []NotificationFeedResult `json:"results"`
}

// NotificationFeedResult is a notification-feed specific response, containing
// the list of activities in the group, plus the extra fields about the group read+seen status.
type NotificationFeedResult struct {
	ID            string     `json:"id"`
	Activities    []Activity `json:"activities"`
	ActivityCount int        `json:"activity_count"`
	ActorCount    int        `json:"actor_count"`
	Group         string     `json:"group"`
	IsRead        bool       `json:"is_read"`
	IsSeen        bool       `json:"is_seen"`
	Verb          string     `json:"verb"`
}

// AddActivityResponse is the API response obtained when adding a single activity
// to a feed.
type AddActivityResponse struct {
	response
	Activity
}

// AddActivitiesResponse is the API response obtained when adding activities to
// a feed.
type AddActivitiesResponse struct {
	response
	Activities []Activity `json:"activities,omitempty"`
}

// Follower is the representation of a feed following another feed.
type Follower struct {
	FeedID   FeedID `json:"feed_id,omitempty"`
	TargetID FeedID `json:"target_id,omitempty"`
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

// AddToManyRequest is the API request body for adding an activity to multiple
// feeds at once.
type AddToManyRequest struct {
	Activity Activity `json:"activity,omitempty"`
	FeedIDs  []FeedID `json:"feeds,omitempty"`
}

// FollowRelationship represents a follow relationship between a source
// ("follower") and a target ("following"), used for FollowMany requests.
type FollowRelationship struct {
	Source FeedID `json:"source,omitempty"`
	Target FeedID `json:"target,omitempty"`
}

// NewFollowRelationship is a helper for creating a FollowRelationship from the
// source ("follower") and target ("following") feeds.
func NewFollowRelationship(source, target Feed) FollowRelationship {
	return FollowRelationship{
		Source: source.ID(),
		Target: target.ID(),
	}
}

type updateToTargetsRequest struct {
	ForeignID string   `json:"foreign_id,omitempty"`
	Time      string   `json:"time,omitempty"`
	New       []FeedID `json:"new_targets,omitempty"`
	Adds      []FeedID `json:"added_targets,omitempty"`
	Removes   []FeedID `json:"removed_targets,omitempty"`
}

// FeedID is a string representing a feed's slug:userID id
type FeedID string
