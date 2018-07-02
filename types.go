package stream

import (
	"encoding/json"
	"fmt"
	"net/url"
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
	var tmp interface{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	switch v := tmp.(type) {
	case string:
		*d, err = durationFromString(v)
	case float64:
		*d, err = durationFromString(fmt.Sprintf("%fs", v))
	default:
		err = fmt.Errorf("invalid duration")
	}
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
	var err error
	for _, layout := range timeLayouts {
		var tt time.Time
		tt, err = time.Parse(layout, s)
		if err == nil {
			return Time{tt}, nil
		}
	}
	return Time{}, err
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

var (
	// ErrMissingNextPage is returned when trying to read the next page of a response
	// which has an empty "next" field.
	ErrMissingNextPage = fmt.Errorf("request missing next page")
	// ErrInvalidNextPage is returned when trying to read the next page of a response
	// which has an invalid "next" field.
	ErrInvalidNextPage = fmt.Errorf("invalid format for Next field")
)

func (r readResponse) parseNext() ([]GetActivitiesOption, error) {
	if r.Next == "" {
		return nil, ErrMissingNextPage
	}

	urlParts := strings.Split(r.Next, "?")
	if len(urlParts) != 2 {
		return nil, ErrInvalidNextPage
	}
	values, err := url.ParseQuery(urlParts[1])
	if err != nil {
		return nil, ErrInvalidNextPage
	}

	var opts []GetActivitiesOption

	limit, ok, err := parseIntValue(values, "limit")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, WithActivitiesLimit(limit))
	}

	offset, ok, err := parseIntValue(values, "offset")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, WithActivitiesOffset(offset))
	}

	if idLT := values.Get("id_lt"); idLT != "" {
		opts = append(opts, WithActivitiesIDLT(idLT))
	}

	if ranking := values.Get("ranking"); ranking != "" {
		opts = append(opts, withActivitiesRanking(ranking))
	}

	return opts, nil
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

// AddToManyRequest is the API request body for adding an activity to multiple
// feeds at once.
type AddToManyRequest struct {
	Activity Activity `json:"activity,omitempty"`
	FeedIDs  []string `json:"feeds,omitempty"`
}

// FollowRelationship represents a follow relationship between a source
// ("follower") and a target ("following"), used for FollowMany requests.
type FollowRelationship struct {
	Source            string `json:"source,omitempty"`
	Target            string `json:"target,omitempty"`
	ActivityCopyLimit *int   `json:"activity_copy_limit,omitempty"`
}

// NewFollowRelationship is a helper for creating a FollowRelationship from the
// source ("follower") and target ("following") feeds.
func NewFollowRelationship(source, target Feed, opts ...FollowRelationshipOption) FollowRelationship {
	r := FollowRelationship{
		Source: source.ID(),
		Target: target.ID(),
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

// FollowRelationshipOption customizes a FollowRelationship.
type FollowRelationshipOption func(r *FollowRelationship)

// WithFollowRelationshipActivityCopyLimit sets the ActivityCopyLimit field for a given FollowRelationship.
func WithFollowRelationshipActivityCopyLimit(activityCopyLimit int) FollowRelationshipOption {
	return func(r *FollowRelationship) {
		r.ActivityCopyLimit = &activityCopyLimit
	}
}

type updateToTargetsRequest struct {
	ForeignID string   `json:"foreign_id,omitempty"`
	Time      string   `json:"time,omitempty"`
	New       []string `json:"new_targets,omitempty"`
	Adds      []string `json:"added_targets,omitempty"`
	Removes   []string `json:"removed_targets,omitempty"`
}

// UnfollowRelationship represents a single follow relationship to remove, used for
// UnfollowMany requests.
type UnfollowRelationship struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	KeepHistory bool   `json:"keep_history"`
}

// CollectionObject is a collection's object.
type CollectionObject struct {
	ID   string
	Name string
	Data map[string]interface{}
}

// MarshalJSON marshals the CollectionObject to a flat JSON object.
func (o CollectionObject) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":   o.ID,
		"name": o.Name,
	}
	for k, v := range o.Data {
		m[k] = v
	}
	return json.Marshal(m)
}

type selectCollectionResponse struct {
	Response selectCollectionResponseData `json:"response"`
}

type selectCollectionResponseData struct {
	Data []selectCollectionResponseObject `json:"data"`
}

type selectCollectionResponseObject struct {
	ForeignID string                 `json:"foreign_id"`
	Data      map[string]interface{} `json:"data"`
}

func (o selectCollectionResponseObject) toCollectionObject() CollectionObject {
	parts := strings.SplitN(o.ForeignID, ":", 2)
	id, name := parts[0], parts[1]
	return CollectionObject{
		ID:   id,
		Name: name,
		Data: o.Data,
	}
}

// PersonalizationResponse is a generic response from the personalization endpoints
// obtained after a PersonalizationClient.Get call.
// Common JSON fields are directly available as struct fields, while non-standard
// JSON fields can be retrieved using the Extra() method.
type PersonalizationResponse struct {
	AppID    int                      `json:"app_id"`
	Duration Duration                 `json:"duration"`
	Limit    int                      `json:"limit"`
	Offset   int                      `json:"offset"`
	Version  string                   `json:"version"`
	Results  []map[string]interface{} `json:"results"`
	extra    map[string]interface{}
}

// Extra returns the non-common response fields as a map[string]interface{}.
func (r *PersonalizationResponse) Extra() map[string]interface{} {
	return r.extra
}

// UnmarshalJSON for PersonalizationResponse is required because of the incoming duration string, and
// for storing non-standard fields without losing their values, so they can be retrieved
// later on with the Extra() function.
func (r *PersonalizationResponse) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	meta, err := decodeData(m, r)
	if err != nil {
		return err
	}
	r.extra = make(map[string]interface{})
	for _, k := range meta.Unused {
		r.extra[k] = m[k]
	}
	return nil
}

// GetActivitiesResponse contains a slice of Activity returned by GetActivitiesByID
// and GetActivitiesByForeignID requests.
type GetActivitiesResponse struct {
	response
	Results []Activity `json:"results"`
}

// ForeignIDTimePair couples an activity's foreignID and timestamp.
type ForeignIDTimePair struct {
	ForeignID string
	Timestamp Time
}

// NewForeignIDTimePair creates a new ForeignIDTimePair with the given foreign ID and timestamp.
func NewForeignIDTimePair(foreignID string, timestamp Time) ForeignIDTimePair {
	return ForeignIDTimePair{
		ForeignID: foreignID,
		Timestamp: timestamp,
	}
}

type updateActivityRequest struct {
	ID        *string                `json:"id,omitempty"`
	ForeignID *string                `json:"foreign_id,omitempty"`
	Time      *Time                  `json:"time,omitempty"`
	Set       map[string]interface{} `json:"set,omitempty"`
	Unset     []string               `json:"unset,omitempty"`
}

// UpdateActivityResponse is the response returned by the UpdateActivityByID and
// UpdateActivityByForeignID methods.
type UpdateActivityResponse struct {
	Activity
}
