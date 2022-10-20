package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

// EnrichedActivity is an enriched Stream activity entity.
type EnrichedActivity struct {
	ID              string                         `json:"id,omitempty"`
	Actor           Data                           `json:"actor,omitempty"`
	Verb            string                         `json:"verb,omitempty"`
	Object          Data                           `json:"object,omitempty"`
	ForeignID       string                         `json:"foreign_id,omitempty"`
	Target          Data                           `json:"target,omitempty"`
	Time            Time                           `json:"time,omitempty"`
	Origin          Data                           `json:"origin,omitempty"`
	To              []string                       `json:"to,omitempty"`
	Score           float64                        `json:"score,omitempty"`
	ReactionCounts  map[string]int                 `json:"reaction_counts,omitempty"`
	OwnReactions    map[string][]*EnrichedReaction `json:"own_reactions,omitempty"`
	LatestReactions map[string][]*EnrichedReaction `json:"latest_reactions,omitempty"`
	Extra           map[string]any                 `json:"-"`
}

// EnrichedReaction is an enriched Stream reaction entity.
type EnrichedReaction struct {
	ID                string                         `json:"id,omitempty"`
	Kind              string                         `json:"kind"`
	ActivityID        string                         `json:"activity_id"`
	UserID            string                         `json:"user_id"`
	Data              map[string]any                 `json:"data,omitempty"`
	TargetFeeds       []string                       `json:"target_feeds,omitempty"`
	ParentID          string                         `json:"parent,omitempty"`
	ChildrenReactions map[string][]*EnrichedReaction `json:"latest_children,omitempty"`
	OwnChildren       map[string][]*EnrichedReaction `json:"own_children,omitempty"`
	ChildrenCounters  map[string]int                 `json:"children_counts,omitempty"`
	User              Data                           `json:"user,omitempty"`
	CreatedAt         Time                           `json:"created_at,omitempty"`
	UpdatedAt         Time                           `json:"updated_at,omitempty"`
}

// UnmarshalJSON decodes the provided JSON payload into the EnrichedActivity. It's required
// because of the custom JSON fields and time formats.
func (a *EnrichedActivity) UnmarshalJSON(b []byte) error {
	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if _, ok := data["to"]; ok {
		tos := data["to"].([]any)
		simpleTos := make([]string, len(tos))
		for i := range tos {
			switch to := tos[i].(type) {
			case string:
				simpleTos[i] = to
			case []any:
				tos, ok := to[0].(string)
				if !ok {
					return errors.New("invalid format for to targets")
				}
				simpleTos[i] = tos
			}
		}
		data["to"] = simpleTos
	}

	if v, ok := data["foreign_id"]; ok { // handle activity reference in foreign id
		if val, ok := v.(map[string]any); ok {
			id, ok := val["id"].(string)
			if !ok {
				return fmt.Errorf("invalid format for enriched referenced activity id: %v", val["id"])
			}
			data["foreign_id_ref"] = data["foreign_id"]
			data["foreign_id"] = "SA:" + id
		}
	}

	return a.decode(data)
}

// MarshalJSON encodes the EnrichedActivity to a valid JSON bytes slice. It's required because of
// the custom JSON fields and time formats.
func (a EnrichedActivity) MarshalJSON() ([]byte, error) {
	s := structs.New(a)
	fields := s.Fields()
	data := s.Map()
	for _, f := range fields {
		tag := f.Tag("json")
		root := strings.TrimSuffix(tag, ",omitempty")

		if f.Kind() != reflect.Struct ||
			(strings.HasSuffix(tag, ",omitempty") &&
				structs.IsZero(f.Value())) {
			continue
		}

		data[root] = f.Value()
	}
	for k, v := range a.Extra {
		data[k] = v
	}

	if _, ok := data["time"]; ok {
		data["time"] = a.Time.Format(TimeLayout)
	}
	return json.Marshal(data)
}

func (a *EnrichedActivity) decode(data map[string]any) error {
	meta, err := decodeData(data, a)
	if err != nil {
		return err
	}
	if len(meta.Unused) > 0 {
		a.Extra = make(map[string]any)
		for _, k := range meta.Unused {
			a.Extra[k] = data[k]
		}
	}
	return nil
}

// EnrichedActivityGroup is a group of enriched Activities obtained from aggregated feeds.
type EnrichedActivityGroup struct {
	baseActivityGroup
	Activities []EnrichedActivity `json:"activities,omitempty"`
	Score      float64            `json:"score,omitempty"`
}

// EnrichedNotificationFeedResult is a notification-feed specific response, containing
// the list of enriched activities in the group, plus the extra fields about the group read+seen status.
type EnrichedNotificationFeedResult struct {
	baseNotificationFeedResult
	Activities []EnrichedActivity `json:"activities"`
}
