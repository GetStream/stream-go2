package stream

import (
	"encoding/json"
	"errors"

	"github.com/fatih/structs"
)

// Activity is a Stream activity entity.
type Activity struct {
	ID        string                 `json:"id,omitempty"`
	Actor     string                 `json:"actor,omitempty"`
	Verb      string                 `json:"verb,omitempty"`
	Object    string                 `json:"object,omitempty"`
	ForeignID string                 `json:"foreign_id,omitempty"`
	Target    string                 `json:"target,omitempty"`
	Time      Time                   `json:"time,omitempty"`
	Origin    string                 `json:"origin,omitempty"`
	To        []string               `json:"to,omitempty"`
	Score     float64                `json:"score,omitempty"`
	Extra     map[string]interface{} `json:"-"`
}

// UnmarshalJSON decodes the provided JSON payload into the Activity. It's required
// because of the custom JSON fields and time formats.
func (a *Activity) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if _, ok := data["to"]; ok {
		tos := data["to"].([]interface{})
		simpleTos := make([]string, len(tos))
		for i := range tos {
			switch to := tos[i].(type) {
			case string:
				simpleTos[i] = to
			case []interface{}:
				tos, ok := to[0].(string)
				if !ok {
					return errors.New("invalid format for to targets")
				}
				simpleTos[i] = tos
			}
		}
		data["to"] = simpleTos
	}

	return a.decode(data)
}

// MarshalJSON encodes the Activity to a valid JSON bytes slice. It's required because of
// the custom JSON fields and time formats.
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
	if len(meta.Unused) > 0 {
		a.Extra = make(map[string]interface{})
		for _, k := range meta.Unused {
			a.Extra[k] = data[k]
		}
	}
	return nil
}

// baseActivityGroup is the common part of responses obtained from reading normal or enriched aggregated feeds.
type baseActivityGroup struct {
	ActivityCount int    `json:"activity_count,omitempty"`
	ActorCount    int    `json:"actor_count,omitempty"`
	Group         string `json:"group,omitempty"`
	ID            string `json:"id,omitempty"`
	Verb          string `json:"verb,omitempty"`
	CreatedAt     Time   `json:"created_at,omitempty"`
	UpdatedAt     Time   `json:"updated_at,omitempty"`
}

// ActivityGroup is a group of Activity obtained from aggregated feeds.
type ActivityGroup struct {
	baseActivityGroup
	Activities []Activity `json:"activities,omitempty"`
}
