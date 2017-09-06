package stream

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type response struct {
	Duration time.Duration `json:"duration"`
	Next     string        `json:"next"`
}

// Activities is a slice of Activity.
type Activities []Activity

// Activity is a Stream activity entity.
type Activity struct {
	ID        string                 `json:"id" structs:"id"`
	Actor     string                 `json:"actor" structs:"actor"`
	Verb      string                 `json:"verb" structs:"verb"`
	Object    string                 `json:"object" structs:"object"`
	ForeignID string                 `json:"foreign_id" structs:"foreign_id"`
	Target    string                 `json:"target" structs:"target"`
	Time      time.Time              `json:"time" structs:"time"`
	To        []string               `json:"to" structs:"to"`
	Score     string                 `json:"score" structs:"score"`
	Extra     map[string]interface{} `json:"-"`
}

func (a *Activity) decodeStringToTime(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t != reflect.TypeOf(time.Time{}) {
		return data, nil
	}
	tt, err := time.Parse(timeLayout, data.(string))
	return tt, err
}

func (a *Activity) decode(data map[string]interface{}) error {
	meta := &mapstructure.Metadata{}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: a.decodeStringToTime,
		Result:     a,
		Metadata:   meta,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(data); err != nil {
		return err
	}
	a.Extra = make(map[string]interface{})
	for _, k := range meta.Unused {
		a.Extra[k] = data[k]
	}
	return nil
}

// UnmarshalJSON decodes the provided JSON payload into the Activity.
func (a *Activity) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	if err := a.decode(data); err != nil {
		return err
	}
	return nil
}

// MarshalJSON encodes the Activity to a valid JSON bytes slice.
func (a *Activity) MarshalJSON() ([]byte, error) {
	data := structs.New(a).Map()
	for k, v := range a.Extra {
		data[k] = v
	}
	return json.Marshal(data)
}

// ActivityGroup is a group of Activity obtained from aggregated feeds.
type ActivityGroup struct {
	Activities    []Activity `json:"activities"`
	ActivityCount int        `json:"activity_count"`
	ActorCount    int        `json:"actor_count"`
	Group         string     `json:"group"`
	ID            string     `json:"id"`
	Verb          string     `json:"verb"`
}

// ActivityGroups is a slice of ActivityGroup.
type ActivityGroups []ActivityGroup

// FlatFeedResponse is the API response obtained when retrieving activities from
// a flat feed.
type FlatFeedResponse struct {
	response
	Results Activities `json:"results"`
}

// UnmarshalJSON decodes the provided JSON payload into the FlatFeedResponse.
func (r *FlatFeedResponse) UnmarshalJSON(b []byte) error {
	type alias FlatFeedResponse
	aux := &struct {
		Duration string `json:"duration"`
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
	Results ActivityGroups `json:"results"`
}

// UnmarshalJSON decodes the provided JSON payload into the
// AggregatedFeedResponse.
func (r *AggregatedFeedResponse) UnmarshalJSON(b []byte) error {
	type alias AggregatedFeedResponse
	aux := &struct {
		Duration string `json:"duration"`
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
	Activities []Activity `json:"activities"`
}

// UnmarshalJSON decodes the provided JSON payload into the
// AddActivitiesResponse.
func (r *AddActivitiesResponse) UnmarshalJSON(b []byte) error {
	type alias AddActivitiesResponse
	aux := &struct {
		Duration string `json:"duration"`
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
	FeedID   string `json:"feed_id"`
	TargetID string `json:"target_id"`
}

type followResponse struct {
	response
	Results []Follower `json:"results"`
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
		Duration string `json:"duration"`
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
	Activity Activity `json:"activity"`
	Feeds    []string `json:"feeds"`
}

// FollowRelationship represents a follow relationship between a source
// ("follower") and a target ("following"), used for FollowMany requests.
type FollowRelationship struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// NewFollowRelationship is a helper for creating a FollowRelationship from the
// source ("follower") and target ("following") feeds.
func NewFollowRelationship(source, target Feed) FollowRelationship {
	return FollowRelationship{
		Source: source.ID(),
		Target: target.ID(),
	}
}
