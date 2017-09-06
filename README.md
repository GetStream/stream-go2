# stream-go2
GetStream.io Go client

This is a Go client for [Stream](https://getstream.io).

## Usage


### Create new client
```go
import "github.com/reifcode/stream-go2"

// Create a new client
client, err := stream.NewClient("your_api_key", "your_api_secret")
```

### Create feeds
```go
flat := client.FlatFeed("user", "123")
aggregated := client.AggregatedFeed("aggregated", "123")
```

### Activities
```go
// Retrieve activities
activities, err := flat.GetActivities()

// Retrieve activities with options
activities, err := flat.GetActivities(stream.WithLimit(10), stream.WithOffset(42))

// Add activities to feeds
johnActivity := stream.Activity{Actor: "john", Verb: "like", Object: "apples"}
aliceActivity := stream.Activity{Actor: "alice", Verb: "like", Object: "pears"}

resp, err := flat.AddActivities(johnActivity, aliceActivity)

// Remove activities
err := flat.RemoveActivityByID(johnActivity.ID)
err := flat.RemoveActivityByForeignID(johnActivity.ForeignID)

// Update activities
johnActivity.Extra["popularity"] = 9000
err := flat.UpdateActivities(johnActivity)
```

### Follow
```go
// Follow a feed
aggregated.Follow(flat)
aggregated.Follow(flat, stream.WithActivityCopyLimit(42))

// Unfollow a feed
aggregated.Unfollow(flat)
aggregated.Unfollow(flat, stream.WithKeepHistory(true))

// Get followers
followers, err := flat.GetFollowers()

// Get following
following, err := aggregated.GetFollowing()
```

### Batch operations
```go
// Add an activity to many feeds at once
activity := stream.Activity{Actor: "bob", Verb: "post", Object: "status updates"}
err := client.AddToMany(activity, flat, aggregated)

// Create multiple follows at once
relationships := []stream.FollowRelationship{
    stream.NewFollowRelationship(aggregated, flat),
}
err := client.FollowMany(relationships)
err := client.FollowMany(relationships, stream.WithActivityCopyLimitQuery(42))
```