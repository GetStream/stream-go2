# stream-go2

stream-go2 is a Go client for [Stream](https://getstream.io) API.

You can sign up for a Stream account at [getstream.io/get_started](https://getstream.io/get_started/).

[![Build Status](https://travis-ci.org/GetStream/stream-go2.svg?branch=master)](https://travis-ci.org/GetStream/stream-go2)
[![godoc](https://godoc.org/github.com/GetStream/stream-go2?status.svg)](https://godoc.org/github.com/GetStream/stream-go2)
[![codecov](https://codecov.io/gh/GetStream/stream-go2/branch/master/graph/badge.svg)](https://codecov.io/gh/GetStream/stream-go2)
[![Go Report Card](https://goreportcard.com/badge/github.com/GetStream/stream-go2)](https://goreportcard.com/report/github.com/GetStream/stream-go2)

## Contents

* [Getting started](#usage)
* [Creating a Client](#creating-a-client)
* [Creating a Feed](#creating-a-feed)
* [Retrieving Activities](#retrieving-activities)
  * [Flat feeds](#flat-feeds)
  * [Aggregated feeds](#aggregated-feeds)
  * [Notification feeds](#notification-feeds)
  * [Options](#options)
* [Adding activities](#adding-activities)
* [Updating activities](#updating-activities)
* [Removing activities](#removing-activities)
* [Retrieving follows](#retrieving-followers-and-followings)
  * [Following](#following)
  * [Followers](#followers)
* [Unfollow](#unfollowing-a-feed)
* [Update `to` targets](#updating-an-activitys-to-targets)
* [Batch activities](#batch-adding-activities)
* [Batch follows](#batch-creating-follows)
* [Realtime tokens](#realtime-tokens)
* [Analytics](#analytics)
  * [Tracking engagement](#tracking-engagement)
  * [Tracking impressions](#tracking-impressions)
  * [Email tracking](#email-tracking)
* [Personalization](#personalization)
* [Collections](#collections)
* [License](#license)

## Usage

Get the client:

```
$ go get gopkg.in/GetStream/stream-go2.v1
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
    stream.WithAPIRegion("us-east"),
    stream.WithAPIVersion("1.0"),
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
You can pass supported options and filters when retrieving activities:
```go
resp, err := flat.GetActivities(
    stream.WithActivitiesIDGTE("f505b3fb-a212-11e7-..."),
    stream.WithActivitiesLimit(5), 
    ...,
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
    stream.WithFollowFeedActivityCopyLimit(15), 
    ...,
)
```

### Retrieving followers and followings
#### Following
Get the feeds that a feed is following:
```go
resp, err := feed.GetFollowing()
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
for _, followed := range resp.Results {
    fmt.Println(followed.FeedID, followed.TargetID)
}
```

You can pass options to `GetFollowing`:
```go
resp, err := feed.GetFollowing(
    stream.WithFollowingLimit(5),
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

You can pass options to `GetFollowers`:
```go
resp, err := feed.GetFollowing(
    stream.WithFollowersLimit(5),
    ...,
)
```

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
    stream.WithUnfollowKeepHistory(true),
    ...,
)
```

### Updating an activity's `to` targets
Remove all old targets and set new ones (replace):
```go
newTargets := []stream.Feed{f1, f2}

err := feed.UpdateToTargets(activity, stream.WithToTargetsNew(newTargets...))
if err != nil {
    // ...
}
```

Add some targets and remove some others:
```go
add := []stream.Feed{target1, target2}
remove := []stream.Feed{oldTarget1, oldTarget2}

err := feed.UpdateToTargets(
    activity, 
    stream.WithToTargetsAdd(add),
    stream.WithToTargetsRemove(remove),
)
if err != nil {
    // ...
}
```

Note: you can't mix `stream.WithToTargetsNew` with `stream.WithToTargetsAdd` or `stream.WithToTargetsRemove`.


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
token := feed.RealtimeToken(false)

// Read-only token
readonlyToken := feed.RealtimeToken(true)
```

## Analytics

If your app is enabled for analytics collection you can use the Go client to track events. The main documentation for the analytics features is available [in our Docs page](https://getstream.io/docs/#analytics_setup).

### Obtaining an Analytics client

You can obtain a specialized Analytics client (`*stream.AnalyticsClient`) from a regular client, which you can use to track events:

```go
// Create the client
analytics := client.Analytics()
```

### Tracking engagement

Engagement events can be tracked with the `TrackEngagement` method of `AnalyticsClient`. It accepts any number of `EngagementEvent`s.

Events' syntax is not checked by the client, so be sure to follow our [documentation](https://getstream.io/docs/#analytics_engagements) about it.

Events are simple maps, but the `stream` package offers handy helpers to populate such events easily.

```go
// Create the event
event := stream.EngagementEvent{}.
    WithLabel("click").
    WithForeignID("event:1234").
    WithUserData(stream.NewUserData().String("john")).
    WithFeatures(
        stream.NewEventFeature("color", "blue"),
        stream.NewEventFeature("shape", "rectangle"),
    ).
    WithLocation("homepage")

// Track the event(s)
err := analytics.TrackEngagement(event)
if err != nil {
    // ...
}
```

### Tracking impressions

Impression events can be tracked with the `TrackImpression` method of `AnalyticsClient` ([syntax docs](https://getstream.io/docs/#analytics_impressions)):

```go
// Create the impression events
imp := stream.ImpressionEventData{}.
    WithForeignIDs("product:1", "product:2", "product:3").
    WithUserData(stream.NewUserData().String("john")).
    WithLocation("storepage")

// Track the events
err := analytics.TrackImpression(imp)
if err != nil {
    // ...
}
```

### Email tracking

You can generate URLs to track events and redirect to a specific URL with the `RedirectAndTrack` method of `AnalyticsClient` ([syntax docs](https://getstream.io/docs/#analytics_email)). It accepts any number of engagement and impression events:

```go
// Create the events
engagement := stream.EngagementEvent{}.
    WithLabel("click").
    WithForeignID("event:1234").
    WithUserData(stream.NewUserData().String("john")).
    WithFeatures(
        stream.NewEventFeature("color", "blue"),
        stream.NewEventFeature("shape", "rectangle"),
    ).
    WithLocation("homepage")

impressions := stream.ImpressionEventData{}.
    WithForeignIDs("product:1", "product:2", "product:3").
    WithUserData(stream.NewUserData().String("john")).
    WithLocation("storepage")

// Generate the tracking and redirect URL, which once followed
// will redirect the user to the targetURL.
targetURL := "https://google.com"
url, err := analytics.RedirectAndTrack(targetURL, engagement, impression)
if err != nil {
    // ...
}

// Display the obtained url where needed.
```

## Personalization

[Personalization endpoints](https://getstream.io/personalization) for enabled apps can be reached using a `PersonalizationClient`, a specialized client obtained with the `Personalization()` function of a regular `Client`.

```go
personalization := client.Personalization()
```

The `PersonalizationClient` exposes three functions that you can use to retrieve and manipulate data: `Get`, `Post`, and `Delete`.

For example, to retrieve follow recommendations:

```go
// Get follow recommendations
data := map[string]interface{}{
    "user_id":          123,
    "source_feed_slug": "timeline",
    "target_feed_slug": "user",
}
resp, err = personalization.Get("follow_recommendations", data)
if err != nil {
    // ...
}
fmt.Println(resp)
```

See the complete [docs and examples](https://getstream.io/docs/#personalization_introduction) about personalization features on Stream's documentation pages.

## Collections

[Collections](https://getstream.io/docs/#collections) endpoints for enabled apps can be reached using a specialized `CollectionsClient` which, like `PersonalizationClient`, can be obtained from a regular `Client`:

```go
collections := client.Collections()
```

`CollectionsClient` exposes three functions, `Upsert`, `Get`, and `Delete`:

```go
// Upsert the "user" collection
object := stream.CollectionObject{
    ID:   "123",
    Name: "johndoe",
    Data: map[string]interface{}{
        "favorite_color": "blue",
    },
}
err = collections.Upsert("user", object)
if err != nil {
    // ...
}

// Get the data from the "user" collection for ID "123" and "456"
objects, err := collections.Get("user", "123", "456")
if err != nil {
    // ...
}

// Delete the data from the "user" collection for user with ID "123"
err = collections.Delete("user", "123")
if err != nil {
    // ...
}
```

See the complete [docs and examples](https://getstream.io/docs/#collections) about collections on Stream's documentation pages.

## License
stream-go2 is licensed under the [GNU General Public License v3.0](LICENSE).

Permissions of this strong copyleft license are conditioned on making available complete source code of licensed works and modifications, which include larger works using a licensed work, under the same license. Copyright and license notices must be preserved. Contributors provide an express grant of patent rights.

See the [LICENSE](LICENSE) file.

## We're hiring!
Would you like to work on cool projects like this?

We are currently hiring for talented Gophers in Amsterdam and Boulder, get in touch with us if you are interested! tommaso@getstream.io
