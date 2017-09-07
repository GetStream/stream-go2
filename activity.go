package stream

import (
	"encoding/json"
	"time"

	"github.com/fatih/structs"
)

// Activities is a slice of Activity.
type Activities []Activity

// Activity is a Stream activity entity.
type Activity struct {
	ID        string                 `json:"id,omitempty"`
	Actor     string                 `json:"actor,omitempty"`
	Verb      string                 `json:"verb,omitempty"`
	Object    string                 `json:"object,omitempty"`
	ForeignID string                 `json:"foreign_id,omitempty"`
	Target    string                 `json:"target,omitempty"`
	Time      time.Time              `json:"time,omitempty"`
	To        []string               `json:"to,omitempty"`
	Score     string                 `json:"score,omitempty"`
	Extra     map[string]interface{} `json:"-"`
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
func (a Activity) MarshalJSON() ([]byte, error) {
	data := structs.New(a).Map()
	for k, v := range a.Extra {
		data[k] = v
	}
	if _, ok := data["time"]; ok {
		data["time"] = a.Time.Format(TimeLayout)
	}
	return json.Marshal(data)
}

func (a *Activity) decode(data map[string]interface{}) error {
	meta, err := decodeData(data, a)
	if err != nil {
		return err
	}

	// this block must become an anonymous function parameter so to generalize
	// the whole process
	if len(meta.Unused) > 0 {
		a.Extra = make(map[string]interface{})
		for _, k := range meta.Unused {
			a.Extra[k] = data[k]
		}
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
