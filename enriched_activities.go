package stream

type EnrichedActivity map[string]interface{}

// EnrichedActivityGroup is a group of enriched Activities obtained from aggregated feeds.
type EnrichedActivityGroup struct {
	baseActivityGroup
	Activities []EnrichedActivity `json:"activities,omitempty"`
}

// EnrichedNotificationFeedResult is a notification-feed specific response, containing
// the list of enriched activities in the group, plus the extra fields about the group read+seen status.
type EnrichedNotificationFeedResult struct {
	baseNotificationFeedResult
	Activities []EnrichedActivity `json:"activities"`
}
