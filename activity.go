package stream

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

// Activities is a slice of Activity.
type Activities []Activity

// Activity is a Stream activity entity.
type Activity struct {
	ID        string                 `json:"id,omitempty" structs:"id,omitempty"`
	Actor     string                 `json:"actor,omitempty" structs:"actor,omitempty"`
	Verb      string                 `json:"verb,omitempty" structs:"verb,omitempty"`
	Object    string                 `json:"object,omitempty" structs:"object,omitempty"`
	ForeignID string                 `json:"foreign_id,omitempty" structs:"foreign_id,omitempty"`
	Target    string                 `json:"target,omitempty" structs:"target,omitempty"`
	Time      time.Time              `json:"time,omitempty" structs:"time,omitempty"`
	To        []string               `json:"to,omitempty" structs:"to,omitempty"`
	Score     string                 `json:"score,omitempty" structs:"score,omitempty"`
	Extra     map[string]interface{} `json:"-" structs:"-"`
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

// ActivityGroup is a group of Activity obtained from aggregated feeds.
type ActivityGroup struct {
	Activities    []Activity `json:"activities,omitempty"`
	ActivityCount int        `json:"activity_count,omitempty"`
	ActorCount    int        `json:"actor_count,omitempty"`
	Group         string     `json:"group,omitempty"`
	ID            string     `json:"id,omitempty"`
	Verb          string     `json:"verb,omitempty"`
}

// ActivityGroups is a slice of ActivityGroup.
type ActivityGroups []ActivityGroup
