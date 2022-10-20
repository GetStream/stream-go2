package stream

import "time"

// EventFeature is a single analytics event feature, a pair of group and
// value strings.
type EventFeature struct {
	Group string `json:"group"`
	Value string `json:"value"`
}

// NewEventFeature returns a new EventFeature from the given group and value
// params.
func NewEventFeature(group, value string) EventFeature {
	return EventFeature{
		Group: group,
		Value: value,
	}
}

// UserData represents an analytics event user data field, which can either
// be a single string/integer representing the user's ID, or a dictionary
// made of an ID (string or integer) and a string alias.
// For example NewUserData().Int(123).Alias("john") will result in a dictionary
// like {"user_data":{"id": 123, "alias": "john"}}, while NewUserData().String("bob") will
// result in {"user_data": "bob"}.
type UserData struct {
	id    any
	alias string
}

// NewUserData initializes an empty UserData type, which must be populated
// using the String, Int, and/or Alias methods.
func NewUserData() *UserData {
	return &UserData{}
}

// String sets the ID as the given string.
func (d *UserData) String(id string) *UserData {
	d.id = id
	return d
}

// Int sets the ID as the given integer.
func (d *UserData) Int(id int) *UserData {
	d.id = id
	return d
}

// Alias sets the alias as the given string.
func (d *UserData) Alias(alias string) *UserData {
	d.alias = alias
	return d
}

func (d *UserData) value() any {
	if d.alias == "" {
		return d.id
	}
	return map[string]any{
		"id":    d.id,
		"alias": d.alias,
	}
}

// EngagementEvent represents an analytics engagement event. It must be populated
// with the available methods, or custom data can be arbitrarily added to it
// manually as key(string),value(any) pairs.
type EngagementEvent map[string]any

// WithLabel sets the event's label field to the given string.
func (e EngagementEvent) WithLabel(label string) EngagementEvent {
	e["label"] = label
	return e
}

// WithUserData sets the event's user_data field to the given UserData's value.
func (e EngagementEvent) WithUserData(userData *UserData) EngagementEvent {
	e["user_data"] = userData.value()
	return e
}

// WithForeignID sets the event's content field to the given foreign ID. If the
// content payload must be an object, use the WithContent method.
func (e EngagementEvent) WithForeignID(foreignID string) EngagementEvent {
	e["content"] = foreignID
	return e
}

// WithContent sets the event's content field to the given content map, and also
// sets the foreign_id field of such object to the given foreign ID string.
// If just the foreign ID is required to be sent, use the WithForeignID method.
func (e EngagementEvent) WithContent(foreignID string, content map[string]any) EngagementEvent {
	if content != nil {
		content["foreign_id"] = foreignID
	}
	e["content"] = content
	return e
}

// WithFeedID sets the event's feed_id field to the given string.
func (e EngagementEvent) WithFeedID(feedID string) EngagementEvent {
	e["feed_id"] = feedID
	return e
}

// WithLocation sets the event's location field to the given string.
func (e EngagementEvent) WithLocation(location string) EngagementEvent {
	e["location"] = location
	return e
}

// WithPosition sets the event's position field to the given int.
func (e EngagementEvent) WithPosition(position int) EngagementEvent {
	e["position"] = position
	return e
}

// WithFeatures sets the event's features field to the given list of EventFeatures.
func (e EngagementEvent) WithFeatures(features ...EventFeature) EngagementEvent {
	e["features"] = features
	return e
}

// WithBoost sets the event's boost field to the given int.
func (e EngagementEvent) WithBoost(boost int) EngagementEvent {
	e["boost"] = boost
	return e
}

// WithTrackedAt sets the event's tracked_at field to the given time.Time.
func (e EngagementEvent) WithTrackedAt(trackedAt time.Time) EngagementEvent {
	e["tracked_at"] = trackedAt.Format(time.RFC3339)
	return e
}

// ImpressionEventsData represents the payload of an arbitrary number of impression events.
// It must be populated with the available methods, or custom data can be arbitrarily
// added to it manually as key(string),value(any) pairs.
type ImpressionEventsData map[string]any

// WithForeignIDs sets the content_list field to the given list of strings.
func (d ImpressionEventsData) WithForeignIDs(foreignIDs ...string) ImpressionEventsData {
	d["content_list"] = foreignIDs
	return d
}

// AddForeignIDs adds the given foreign ID strings to the content_list field, creating
// it if it doesn't exist.
func (d ImpressionEventsData) AddForeignIDs(foreignIDs ...string) ImpressionEventsData {
	list, ok := d["content_list"].([]string)
	if !ok {
		return d.WithForeignIDs(foreignIDs...)
	}
	return d.WithForeignIDs(append(list, foreignIDs...)...)
}

// WithUserData sets the user_data field to the given UserData value.
func (d ImpressionEventsData) WithUserData(userData *UserData) ImpressionEventsData {
	d["user_data"] = userData.value()
	return d
}

// WithFeedID sets the feed_id field to the given string.
func (d ImpressionEventsData) WithFeedID(feedID string) ImpressionEventsData {
	d["feed_id"] = feedID
	return d
}

// WithLocation sets the location field to the given string.
func (d ImpressionEventsData) WithLocation(location string) ImpressionEventsData {
	d["location"] = location
	return d
}

// WithPosition sets the position field to the given int.
func (d ImpressionEventsData) WithPosition(position int) ImpressionEventsData {
	d["position"] = position
	return d
}

// WithFeatures sets the features field to the given list of EventFeatures.
func (d ImpressionEventsData) WithFeatures(features ...EventFeature) ImpressionEventsData {
	d["features"] = features
	return d
}

// WithTrackedAt sets the tracked_at field to the given time.Time.
func (d ImpressionEventsData) WithTrackedAt(trackedAt time.Time) ImpressionEventsData {
	d["tracked_at"] = trackedAt.Format(time.RFC3339)
	return d
}
