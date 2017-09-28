# stream-go2

stream-go2 is a Go client for [Stream](https://getstream.io) API.

You can sign up for a Stream account at [getstream.io/get_started](https://getstream.io/get_started/).

[![Build Status](https://travis-ci.org/GetStream/stream-go2.svg?branch=master)](https://travis-ci.org/GetStream/stream-go2)
[![godoc](https://godoc.org/github.com/GetStream/stream-go2?status.svg)](https://godoc.org/github.com/GetStream/stream-go2)
[![codecov](https://codecov.io/gh/GetStream/stream-go2/branch/master/graph/badge.svg)](https://codecov.io/gh/GetStream/stream-go2)


## Usage

Get the client:

```
$ go get github.com/GetStream/stream-go2
```

stream-go2 uses [dep](https://github.com/golang/dep) for managing dependencies (see `Gopkg.toml` and `Gopkg.lock`).
You can get required dependencies simply by running:
```
$ dep ensure
```

Even better: at Stream we have developed [vg](https://github.com/GetStream/vg), a powerful workspace manager for Go based on `dep` itself.
If you use vg (and you should!) you can just:
```
$ vg init && vg ensure
```

### Creating a Client

```go
key := "YOUR_API_KEY"
secret := "YOUR_API_SECRET"

client, err := stream.NewClient(key, secret)
if err != nil {
    // ...
}
```

You can pass additional options when creating a client using the available `ClientOption` functions:

```go
client, err := stream.NewClient(key, secret, 
    ClientWithRegion("us-east"),
    ClientWithVersion("1.0"),
    ...,
)
```

You can also create a client using environment variables:
```go
client, err := stream.NewClientFromEnv()
```

Available environment variables:
* `STREAM_API_KEY`
* `STREAM_API_SECRET`
* `STREAM_API_REGION`
* `STREAM_API_VERSION`

### Creating a Feed

Create a flat feed from slug and user ID:
```go
flat := client.FlatFeed("user", "123")
```

Create an aggregated feed from slug and user ID:
```go
aggr := client.AggregatedFeed("aggregated", "123")
```

Create a notification feed from slug and user ID:
```go
notif := client.NotificationFeed("notification", "123")
```

Flat, aggregated, and notification feeds implement the `Feed` interface methods.

In the snippets below, `feed` indicates any kind of feed, while `flat`, `aggregated`, and `notification` are used
to indicate that only that kind of feed has certain methods or can perform certain operations.

### Retrieving activities

#### Flat feeds
```go
resp, err := flat.GetActivities()
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Next:", resp.Next)
fmt.Println("Activities:")
for _, activity := range resp.Results {
    fmt.Println(activity)
}
```

You can retrieve flat feeds with [custom ranking](https://getstream.io/docs/#custom_ranking), using the dedicated method:
```go
resp, err := flat.GetActivitiesWithRanking("popularity")
if err != nil {
    // ...
}
```

#### Aggregated feeds
```go
resp, err := aggregated.GetActivities()
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Next:", resp.Next)
fmt.Println("Groups:")
for _, group := range resp.Results {
    fmt.Println("Group:", group.Name, "ID:", group.ID, "Verb:", group.Verb)
    fmt.Println("Activities:", group.ActivityCount, "Actors:", group.ActorCount)
    for _, activity := range group.Activities {
        // ...
    }
}
```

#### Notification feeds
```go
resp, err := notification.GetActivities()
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Next:", resp.Next)
fmt.Println("Unseen:", resp.Unseen, "Unread:", resp.Unread)
fmt.Println("Groups:")
for _, group := range resp.Results {
    fmt.Println("Group:", group.Group, "ID:", group.ID, "Verb:", group.Verb)
    fmt.Println("Seen:", group.IsSeen, "Read:", group.IsRead)
    fmt.Println("Activities:", group.ActivityCount, "Actors:", group.ActorCount)
    for _, activity := range group.Activities {
        // ...
    }
}
```


#### Options
You can pass supported options and filters when retrieving activities. For example:
```go
resp, err := flat.GetActivities(
    stream.GetActivitiesWithIDGTE("f505b3fb-a212-11e7-..."),
    stream.GetActivitiesWithLimit(5), 
)
```

### Adding activities
Add a single activity:
```go
resp, err := feed.AddActivity(stream.Activity{Actor: "bob", ...})
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Activity:", resp.Activity) // resp wraps the stream.Activity type
```

Add multiple activities:
```go
a1 := stream.Activity{Actor: "bob", ...}
a2 := stream.Activity{Actor: "john", ...}
a3 := stream.Activity{Actor: "alice", ...}

resp, err := feed.AddActivities(a1, a2, a3)
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Activities:")
for _, activity := range resp.Activities {
    fmt.Println(activity)
}
```

### Updating activities
```go
err := feed.UpdateActivities(a1, a2, ...)
if err != nil {
    // ...
}
```

### Removing activities
You can either remove activities by ID or ForeignID:
```go
err := feed.RemoveActivityByID("f505b3fb-a212-11e7-...")
if err != nil {
    // ...
}

err := feed.RemoveActivityByForeignID("bob:123")
if err != nil {
    // ...
}
```

### Following another feed
```go
err := feed.Follow(anotherFeed)
if err != nil {
    // ...
}
```
Beware that it's possible to follow only flat feeds.

#### Options
You can pass options to the `Follow` method. For example:
```go
err := feed.Follow(anotherFeed, 
    stream.FollowWithActivityCopyLimit(15), 
    ...,
)
```

### Retrieving followers and followings
#### Followings
Get the feeds that a feed is following:
```go
resp, err := feed.GetFollowings()
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
for _, followed := range resp.Results {
    fmt.Println(followed.FeedID, followed.TargetID)
}
```

You can pass options to `GetFollowings`:
```go
resp, err := feed.GetFollowings(
    stream.FollowingWithLimit(5),
    ...,
)
```

#### Followers
```go
resp, err := flat.GetFollowers()
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
for _, follower := range resp.Results {
    fmt.Println(follower.FeedID, follower.TargetID)
}
```
Note: this is only possible for `FlatFeed` types.

### Unfollowing a feed
```go
err := flat.Unfollow(anotherFeed)
if err != nil {
    // ...
}
```

You can pass options to `Unfollow`:
```go
err := flat.Unfollow(anotherFeed,
    stream.UnfollowWithKeepHistory(true),
    ...,
)
```

### Updating an activity's `to` targets
Remove all old targets and set new ones (replace):
```go
newTargets := []stream.Feed{f1, f2}

err := feed.ReplaceToTargets(activity, newTargets)
if err != nil {
    // ...
}
```

Add some targets and remove some others:
```go
add := []stream.Feed{target1, target2}
remove := []stream.Feed{oldTarget1, oldTarget2}

err := feed.UpdateToTargets(activity, add, remove)
if err != nil {
    // ...
}
```

### Batch adding activities
You can add the same activities to multiple feeds at once with the `(*Client).AddToMany` method ([docs](https://getstream.io/docs_rest/#add_to_many)):
```go
err := client.AddToMany(activity,
    feed1, feed2, ...,
)
if err != nil {
    // ...
}
```

### Batch creating follows
You can create multiple follow relationships at once with the `(*Client).FollowMany` method ([docs](https://getstream.io/docs_rest/#follow_many)):
```go
relationships := []stream.FollowRelationship{
    stream.NewFollowRelationship(source, target),
    ...,
}

err := client.FollowMany(relationships)
if err != nil {
    // ...
}
```

### Realtime tokens
You can get a token suitable for client-side [real-time feed updates](https://getstream.io/docs/go/#realtime) as:
```go
// Read+Write token
token := feed.Token(false)

// Read-only token
readonlyToken := feed.Token(true)
```

## License
stream-go2 is licensed under the [GNU General Public License v3.0](LICENSE).

Permissions of this strong copyleft license are conditioned on making available complete source code of licensed works and modifications, which include larger works using a licensed work, under the same license. Copyright and license notices must be preserved. Contributors provide an express grant of patent rights.

See the [LICENSE](LICENSE) file.

## We're hiring!
Would you like to work on cool projects like this?

We are currently hiring for talented Gophers in Amsterdam and Boulder, get in touch with us if you are interested! tommaso@getstream.io