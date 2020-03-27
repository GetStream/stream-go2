package stream

type ClientInterface interface {
	// FlatFeed returns a new Flat Feed with the provided slug and userID.
	FlatFeed(slug, userID string) (*FlatFeed, error)
	// AggregatedFeed returns a new Aggregated Feed with the provided slug and
	// userID.
	AggregatedFeed(slug, userID string) (*AggregatedFeed, error)
	// NotificationFeed returns a new Notification Feed with the provided slug and
	// userID.
	NotificationFeed(slug, userID string) (*NotificationFeed, error)
	// AddToMany adds an activity to multiple feeds at once.
	AddToMany(activity Activity, feeds ...Feed) error
	// FollowMany creates multiple follows at once.
	FollowMany(relationships []FollowRelationship, opts ...FollowManyOption) error
	// UnfollowMany removes multiple follow relationships at once.
	UnfollowMany(relationships []UnfollowRelationship) error
	// Analytics returns a new AnalyticsClient sharing the base configuration of the original Client.
	Analytics() *AnalyticsClient
	// Collections returns a new CollectionsClient.
	Collections() *CollectionsClient
	// Users returns a new UsersClient.
	Users() *UsersClient
	// Reactions returns a new ReactionsClient.
	Reactions() *ReactionsClient
	// Personalization returns a new PersonalizationClient.
	Personalization() *PersonalizationClient
	// GetActivitiesByID returns activities for the current app having the given IDs.
	GetActivitiesByID(ids ...string) (*GetActivitiesResponse, error)
	// GetEnrichedActivitiesByID returns enriched activities for the current app having the given IDs.
	GetEnrichedActivitiesByID(ids ...string) (*GetEnrichedActivitiesResponse, error)
	// GetActivitiesByForeignID returns activities for the current app having the given foreign IDs and timestamps.
	GetActivitiesByForeignID(values ...ForeignIDTimePair) (*GetActivitiesResponse, error)
	// GetEnrichedActivitiesByForeignID returns enriched activities for the current app having the given foreign IDs and timestamps.
	GetEnrichedActivitiesByForeignID(values ...ForeignIDTimePair) (*GetEnrichedActivitiesResponse, error)
	// UpdateActivities updates existing activities.
	UpdateActivities(activities ...Activity) error
	// PartialUpdateActivities performs a partial update on multiple activities with the given set and unset operations
	// specified by each changeset. This returns the affected activities.
	PartialUpdateActivities(changesets ...UpdateActivityRequest) (*UpdateActivitiesResponse, error)
	// UpdateActivityByID performs a partial activity update with the given set and unset operations, returning the
	// affected activity, on the activity with the given ID.
	UpdateActivityByID(id string, set map[string]interface{}, unset []string) (*UpdateActivityResponse, error)
	// UpdateActivityByForeignID performs a partial activity update with the given set and unset operations, returning the
	// affected activity, on the activity with the given foreign ID and timestamp.
	UpdateActivityByForeignID(foreignID string, timestamp Time, set map[string]interface{}, unset []string) (*UpdateActivityResponse, error)
	GetUserSessionToken(userID string) (string, error)
	GetUserSessionTokenWithClaims(userID string, claims map[string]interface{}) (string, error)
}
