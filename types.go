package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
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
		err = errors.New("invalid duration")
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
	*t, err = timeFromString(strings.ReplaceAll(string(b), `"`, ""))
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

// Data is a representation of an enriched activities enriched object,
// such as the the user or the object
type Data struct {
	ID    string                 `json:"id"`
	Extra map[string]interface{} `json:"-"`
}

func (a *Data) decode(data map[string]interface{}) error {
	// We are not using decodeData here because we do not need the DecodeHook
	// since it leads to a stack overflow
	cfg := &mapstructure.DecoderConfig{
		Result:   a,
		Metadata: &mapstructure.Metadata{},
		TagName:  "json",
	}
	dec, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}
	if err := dec.Decode(data); err != nil {
		return err
	}

	if len(cfg.Metadata.Unused) > 0 {
		a.Extra = make(map[string]interface{})
		for _, k := range cfg.Metadata.Unused {
			a.Extra[k] = data[k]
		}
	}
	return nil
}

// MarshalJSON encodes data into json.
func (a Data) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(a.Extra))
	for k, v := range a.Extra {
		m[k] = v
	}
	result := map[string]interface{}{
		"id":   a.ID,
		"data": m,
	}
	return json.Marshal(result)
}

// Rate is the information to be filled if a rate limited response
type Rate struct {
	// Reset is the time for limit to reset
	Reset Time

	// Limit is the existing limit for the resource
	Limit int

	// Remaining is the remaining possible calls
	Remaining int
}

func NewRate(headers http.Header) *Rate {
	var r Rate
	limit, err := strconv.Atoi(headers.Get(HeaderRateLimit))
	if err == nil {
		r.Limit = limit
	}
	remaining, err := strconv.Atoi(headers.Get(HeaderRateRemaining))
	if err == nil {
		r.Remaining = remaining
	}
	reset, err := strconv.ParseInt(headers.Get(HeaderRateReset), 10, 64)
	if err == nil && reset > 0 {
		r.Reset = Time{Time: time.Unix(reset, 0)}
	}
	return &r
}

// Response is the part of StreamAPI responses common throughout the API.
type response struct {
	Rate     Rate     `json:"ratelimit,omitempty"`
	Duration Duration `json:"duration,omitempty"`
}

type BaseResponse struct {
	response
}

// readResponse is the part of StreamAPI responses common for GetActivities API requests.
type readResponse struct {
	response
	Next string `json:"next,omitempty"`
}

var (
	// ErrMissingNextPage is returned when trying to read the next page of a response
	// which has an empty "next" field.
	ErrMissingNextPage = errors.New("request missing next page")
	// ErrInvalidNextPage is returned when trying to read the next page of a response
	// which has an invalid "next" field.
	ErrInvalidNextPage = errors.New("invalid format for Next field")
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

	if enrichOpt := values.Get("withOwnReactions"); parseBool(enrichOpt) {
		opts = append(opts, WithEnrichOwnReactions())
	}

	if enrichOpt := values.Get("withRecentReactions"); parseBool(enrichOpt) {
		opts = append(opts, WithEnrichRecentReactions())
	}

	if enrichOpt := values.Get("withReactionCounts"); parseBool(enrichOpt) {
		opts = append(opts, WithEnrichReactionCounts())
	}

	if enrichOpt := values.Get("withOwnChildren"); parseBool(enrichOpt) {
		opts = append(opts, WithEnrichOwnChildren())
	}

	reactionsLimit, ok, err := parseIntValue(values, "recentReactionsLimit")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, WithEnrichRecentReactionsLimit(reactionsLimit))
	}

	if enrichOpt := values.Get("reactionKindsFilter"); enrichOpt != "" {
		kinds := strings.Split(enrichOpt, ",")
		opts = append(opts, WithEnrichReactionKindsFilter(kinds...))
	}

	return opts, nil
}

// baseNotificationFeedResponse is the common part of responses obtained from reading normal or enriched notification feeds.
type baseNotificationFeedResponse struct {
	readResponse
	Unseen int `json:"unseen"`
	Unread int `json:"unread"`
}

// baseNotificationFeedResult is the common part of responses obtained from reading normal or enriched notification feeds.
type baseNotificationFeedResult struct {
	ID            string `json:"id"`
	ActivityCount int    `json:"activity_count"`
	ActorCount    int    `json:"actor_count"`
	Group         string `json:"group"`
	IsRead        bool   `json:"is_read"`
	IsSeen        bool   `json:"is_seen"`
	Verb          string `json:"verb"`
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
	baseNotificationFeedResponse
	Results []NotificationFeedResult `json:"results"`
}

// NotificationFeedResult is a notification-feed specific response, containing
// the list of activities in the group, plus the extra fields about the group read+seen status.
type NotificationFeedResult struct {
	baseNotificationFeedResult
	Activities []Activity `json:"activities"`
}

// AddActivityResponse is the API response obtained when adding a single activity
// to a feed.
type AddActivityResponse struct {
	activityResponse
}

type activityResponse struct {
	response
	Activity
}

// UnmarshalJSON is the custom unmarshaler since activity custom unmarshaler
// can take extra values.
func (a *activityResponse) UnmarshalJSON(buf []byte) error {
	var r response
	if err := json.Unmarshal(buf, &r); err != nil {
		return err
	}
	var ac Activity
	if err := json.Unmarshal(buf, &ac); err != nil {
		return err
	}
	delete(ac.Extra, "duration")
	delete(ac.Extra, "ratelimit")
	*a = activityResponse{response: r, Activity: ac}
	return nil
}

// RemoveActivityResponse is the API response obtained when removing an activity
// from a feed.
type RemoveActivityResponse struct {
	response
	Removed string `json:"removed"`
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

// followResponse is the API response obtained when retrieving follow graph
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

// UpdateToTargetsResponse is the response for updating to targets of an activity.
type UpdateToTargetsResponse struct {
	response
	Activity map[string]interface{} `json:"activity"`
	Added    []string               `json:"added"`
	Removed  []string               `json:"removed"`
}

// UnfollowRelationship represents a single follow relationship to remove, used for
// UnfollowMany requests.
type UnfollowRelationship struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	KeepHistory bool   `json:"keep_history"`
}

// NewUnfollowRelationship is a helper for creating an UnfollowRelationship from the
// source ("follower") and target ("following") feeds.
func NewUnfollowRelationship(source, target Feed, opts ...UnfollowRelationshipOption) UnfollowRelationship {
	r := UnfollowRelationship{
		Source: source.ID(),
		Target: target.ID(),
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

// WithUnfollowRelationshipKeepHistory sets the KeepHistory field for a given UnfollowRelationship.
func WithUnfollowRelationshipKeepHistory() UnfollowRelationshipOption {
	return func(r *UnfollowRelationship) {
		r.KeepHistory = true
	}
}

// UnfollowRelationshipOption customizes an UnfollowRelationship.
type UnfollowRelationshipOption func(r *UnfollowRelationship)

// CollectionObject is a collection's object.
type CollectionObject struct {
	ID   string                 `json:"id,omitempty"`
	Data map[string]interface{} `json:"data"`
}

type CollectionObjectResponse struct {
	response
	CollectionObject `json:",inline"`
}

// MarshalJSON marshals the CollectionObject to a flat JSON object.
func (o CollectionObject) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id": o.ID,
	}
	for k, v := range o.Data {
		m[k] = v
	}
	return json.Marshal(m)
}

type addCollectionRequest struct {
	UserID *string `json:"user_id,omitempty"`
	CollectionObject
}

func (r addCollectionRequest) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":   r.ID,
		"data": r.Data,
	}
	if r.UserID != nil {
		m["user_id"] = r.UserID
	}

	return json.Marshal(m)
}

// GetCollectionResponseObject represents a single response coming from a Collection
// Get request after a CollectionsClient.Get call.
type GetCollectionResponseObject struct {
	ForeignID string                 `json:"foreign_id"`
	Data      map[string]interface{} `json:"data"`
}

type getCollectionResponse struct {
	Data []GetCollectionResponseObject `json:"data"`
}

type getCollectionResponseWrap struct {
	response
	Response getCollectionResponse `json:"response"`
}

// GetCollectionResponse represents a single response coming from a Collection Select
// request after a CollectionsClient.Select call.
type GetCollectionResponse struct {
	response
	Objects []GetCollectionResponseObject
}

// User represents a user
type User struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type UserResponse struct {
	response
	User `json:",inline"`
}

// Reaction is a reaction retrieved from the API.
type Reaction struct {
	AddReactionRequestObject
	ChildrenReactions map[string][]*Reaction `json:"latest_children,omitempty"`
	OwnChildren       map[string][]*Reaction `json:"own_children,omitempty"`
	ChildrenCounters  map[string]interface{} `json:"children_counts,omitempty"`
}

type ReactionResponse struct {
	response
	Reaction `json:",inline"`
}

// AddReactionRequestObject is an object used only when calling the Add* reaction endpoints
type AddReactionRequestObject struct {
	ID                   string                 `json:"id,omitempty"`
	Kind                 string                 `json:"kind"`
	ActivityID           string                 `json:"activity_id"`
	UserID               string                 `json:"user_id"`
	Data                 map[string]interface{} `json:"data,omitempty"`
	TargetFeeds          []string               `json:"target_feeds,omitempty"`
	TargetFeedsExtraData map[string]interface{} `json:"target_feeds_extra_data,omitempty"`
	ParentID             string                 `json:"parent,omitempty"`
}

// filterResponse is the part of StreamAPI responses common for FilterReactions API requests.
type filterResponse struct {
	response
	Next string `json:"next,omitempty"`
}

func (r filterResponse) parseNext() ([]FilterReactionsOption, error) {
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

	var opts []FilterReactionsOption

	limit, ok, err := parseIntValue(values, "limit")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, WithLimit(limit))
	}

	if idLT := values.Get("id_lt"); idLT != "" {
		opts = append(opts, WithIDLT(idLT))
	}

	if idLT := values.Get("id_gt"); idLT != "" {
		opts = append(opts, WithIDGT(idLT))
	}

	if withActData := values.Get("with_activity_data"); withActData != "" {
		if val := strings.ToLower(withActData); val == "true" || val == "t" || val == "1" {
			opts = append(opts, WithActivityData())
		}
	}

	if withOwnChildren := values.Get("with_own_children"); withOwnChildren != "" {
		if val := strings.ToLower(withOwnChildren); val == "true" || val == "t" || val == "1" {
			opts = append(opts, WithOwnChildren())
		}
	}

	return opts, nil
}

// FilterReactionResponse is the response received from the ReactionsClient.Filter call.
type FilterReactionResponse struct {
	filterResponse
	Results  []Reaction             `json:"results"`
	Activity map[string]interface{} `json:"activity"`
	meta     filterReactionsRequestMetadata
}

// filterReactionsRequestMetadata holds the initial request metadata used for pagination.
type filterReactionsRequestMetadata struct {
	attr FilterReactionsAttribute
}

// PersonalizationResponse is a generic response from the personalization endpoints
// obtained after a PersonalizationClient.Get call.
// Common JSON fields are directly available as struct fields, while non-standard
// JSON fields can be retrieved using the Extra() method.
type PersonalizationResponse struct {
	AppID    int                      `json:"app_id"`
	Duration Duration                 `json:"duration"`
	Rate     Rate                     `json:"ratelimit"`
	Limit    int                      `json:"limit"`
	Offset   int                      `json:"offset"`
	Version  string                   `json:"version"`
	Next     string                   `json:"next"`
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

// EnrichedFlatFeedResponse is the API response obtained when retrieving enriched activities from
// a flat feed.
type EnrichedFlatFeedResponse struct {
	readResponse
	Results []EnrichedActivity `json:"results,omitempty"`
}

// EnrichedAggregatedFeedResponse is the API response obtained when retrieving
// enriched activities from an aggregated feed.
type EnrichedAggregatedFeedResponse struct {
	readResponse
	Results []EnrichedActivityGroup `json:"results,omitempty"`
}

// EnrichedNotificationFeedResponse is the API response obtained when retrieving enriched activities
// from a notification feed.
type EnrichedNotificationFeedResponse struct {
	baseNotificationFeedResponse
	Results []EnrichedNotificationFeedResult `json:"results"`
}

// GetActivitiesResponse contains a slice of Activity returned by GetActivitiesByID
// and GetActivitiesByForeignID requests.
type GetActivitiesResponse struct {
	response
	Results []Activity `json:"results"`
}

// GetEnrichedActivitiesResponse contains a slice of enriched Activity returned by GetEnrichedActivitiesByID
// and GetEnrichedActivitiesByForeignID requests.
type GetEnrichedActivitiesResponse struct {
	response
	Results []EnrichedActivity `json:"results"`
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

// UpdateActivityRequest is the API request body for partially updating an activity.
type UpdateActivityRequest struct {
	ID        *string                `json:"id,omitempty"`
	ForeignID *string                `json:"foreign_id,omitempty"`
	Time      *Time                  `json:"time,omitempty"`
	Set       map[string]interface{} `json:"set,omitempty"`
	Unset     []string               `json:"unset,omitempty"`
}

// NewUpdateActivityRequestByID creates a new UpdateActivityRequest to be used by PartialUpdateActivities
func NewUpdateActivityRequestByID(id string, set map[string]interface{}, unset []string) UpdateActivityRequest {
	return UpdateActivityRequest{
		ID:    &id,
		Set:   set,
		Unset: unset,
	}
}

// NewUpdateActivityRequestByForeignID creates a new UpdateActivityRequest to be used by PartialUpdateActivities
func NewUpdateActivityRequestByForeignID(foreignID string, timestamp Time, set map[string]interface{}, unset []string) UpdateActivityRequest {
	return UpdateActivityRequest{
		ForeignID: &foreignID,
		Time:      &timestamp,
		Set:       set,
		Unset:     unset,
	}
}

// UpdateActivityResponse is the response returned by the UpdateActivityByID and
// UpdateActivityByForeignID methods.
type UpdateActivityResponse struct {
	activityResponse
}

type UpdateActivitiesResponse struct {
	response
	Activities []*Activity `json:"activities"`
}
