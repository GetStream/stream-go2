package stream

import "time"

type EventFeature struct {
	Group string `json:"group"`
	Value string `json:"value"`
}

func NewEventFeature(group, value string) EventFeature {
	return EventFeature{
		Group: group,
		Value: value,
	}
}

type UserData struct {
	id    interface{}
	alias string
}

func NewUserData() *UserData {
	return &UserData{}
}

func (d *UserData) String(id string) *UserData {
	d.id = id
	return d
}

func (d *UserData) Int(id int) *UserData {
	d.id = id
	return d
}

func (d *UserData) Alias(alias string) *UserData {
	d.alias = alias
	return d
}

func (d *UserData) value() interface{} {
	if d.alias == "" {
		return d.id
	}
	return map[string]interface{}{
		"id":    d.id,
		"alias": d.alias,
	}
}

type EngagementEvent map[string]interface{}

func (e EngagementEvent) WithLabel(label string) EngagementEvent {
	e["label"] = label
	return e
}

func (e EngagementEvent) WithUserData(userData *UserData) EngagementEvent {
	e["user_data"] = userData.value()
	return e
}

func (e EngagementEvent) WithForeignID(foreignID string) EngagementEvent {
	e["content"] = foreignID
	return e
}

func (e EngagementEvent) WithContent(foreignID string, content map[string]interface{}) EngagementEvent {
	if content != nil {
		content["foreign_id"] = foreignID
	}
	e["content"] = content
	return e
}

func (e EngagementEvent) WithFeedID(feedID string) EngagementEvent {
	e["feed_id"] = feedID
	return e
}

func (e EngagementEvent) WithLocation(location string) EngagementEvent {
	e["location"] = location
	return e
}

func (e EngagementEvent) WithPosition(position int) EngagementEvent {
	e["position"] = position
	return e
}

func (e EngagementEvent) WithFeatures(features ...EventFeature) EngagementEvent {
	e["features"] = features
	return e
}

func (e EngagementEvent) WithBoost(boost int) EngagementEvent {
	e["boost"] = boost
	return e
}

func (e EngagementEvent) WithTrackedAt(trackedAt time.Time) EngagementEvent {
	e["tracked_at"] = trackedAt.Format(time.RFC3339)
	return e
}

type ImpressionEventsData map[string]interface{}

func (d ImpressionEventsData) WithForeignIDs(foreignIDs ...string) ImpressionEventsData {
	d["content_list"] = foreignIDs
	return d
}

func (d ImpressionEventsData) AddForeignIDs(foreignIDs ...string) ImpressionEventsData {
	list, ok := d["content_list"].([]string)
	if !ok {
		return d.WithForeignIDs(foreignIDs...)
	}
	return d.WithForeignIDs(append(list, foreignIDs...)...)
}

func (d ImpressionEventsData) WithUserData(userData *UserData) ImpressionEventsData {
	d["user_data"] = userData.value()
	return d
}

func (d ImpressionEventsData) WithFeedID(feedID string) ImpressionEventsData {
	d["feed_id"] = feedID
	return d
}

func (d ImpressionEventsData) WithLocation(location string) ImpressionEventsData {
	d["location"] = location
	return d
}

func (d ImpressionEventsData) WithPosition(position int) ImpressionEventsData {
	d["position"] = position
	return d
}

func (d ImpressionEventsData) WithFeatures(features ...EventFeature) ImpressionEventsData {
	d["features"] = features
	return d
}

func (d ImpressionEventsData) WithBoost(boost int) ImpressionEventsData {
	d["boost"] = boost
	return d
}

func (d ImpressionEventsData) WithTrackedAt(trackedAt time.Time) ImpressionEventsData {
	d["tracked_at"] = trackedAt.Format(time.RFC3339)
	return d
}
