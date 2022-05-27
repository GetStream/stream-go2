# Official Go SDK for [Stream Feeds](https://getstream.io/activity-feeds/) <!-- omit in toc -->

<p align="center">
    <img src="./assets/logo.svg" width="50%" height="50%">
</p>
<p align="center">
    Official Go API client for Stream Feeds, a service for building applications with activity feeds.
    <br />
    <a href="https://getstream.io/activity-feeds/docs/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/GetStream/stream-go2/issues">Report Bug</a>
    ·
    <a href="https://github.com/GetStream/stream-go2/issues">Request Feature</a>
</p>

---
> ## 🚨 Breaking change in v7.0 <
> From version 7.0.0, all methods' first argument is `context.Context`. The reason is that we internally changed `http.NewRequest()` to `http.NewRequestWithContext()` so that it's easier to handle cancellations, timeouts and deadlines for our users.

## About Stream
stream-go2 is a Go client for [Stream](https://getstream.io) Feeds API.

You can sign up for a Stream account at our [Get Started](https://getstream.io/chat/get_started/) page.

You can use this library to access chat API endpoints server-side.

For the client-side integrations (web and mobile) have a look at the JavaScript, iOS and Android SDK libraries ([docs](https://getstream.io/activity-feeds/)).

> 💡 Note: this is a library for the **Feeds** product. The Chat SDKs can be found [here](https://getstream.io/chat/docs/).

[![build](https://github.com/GetStream/stream-go2/workflows/build/badge.svg)](https://github.com/GetStream/stream-go2/actions)
[![godoc](https://pkg.go.dev/badge/GetStream/stream-go2)](https://pkg.go.dev/github.com/GetStream/stream-go2/v7?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/GetStream/stream-go2)](https://goreportcard.com/report/github.com/GetStream/stream-go2)

## Contents

- [About Stream](#about-stream)
- [Contents](#contents)
- [Installation](#installation)
  - [Creating a Client](#creating-a-client)
  - [Rate Limits](#rate-limits)
  - [Creating a Feed](#creating-a-feed)
  - [Retrieving activities](#retrieving-activities)
    - [Flat feeds](#flat-feeds)
    - [Aggregated feeds](#aggregated-feeds)
    - [Notification feeds](#notification-feeds)
    - [Options](#options)
  - [Adding activities](#adding-activities)
  - [Updating activities](#updating-activities)
  - [Partially updating activities](#partially-updating-activities)
  - [Removing activities](#removing-activities)
  - [Following another feed](#following-another-feed)
    - [Options](#options-1)
  - [Retrieving followers and followings](#retrieving-followers-and-followings)
    - [Following](#following)
    - [Followers](#followers)
  - [Unfollowing a feed](#unfollowing-a-feed)
  - [Updating an activity's `to` targets](#updating-an-activitys-to-targets)
  - [Batch adding activities](#batch-adding-activities)
  - [Batch creating follows](#batch-creating-follows)
  - [Realtime tokens](#realtime-tokens)
- [Analytics](#analytics)
  - [Obtaining an Analytics client](#obtaining-an-analytics-client)
  - [Tracking engagement](#tracking-engagement)
  - [Tracking impressions](#tracking-impressions)
  - [Email tracking](#email-tracking)
- [Personalization](#personalization)
- [Collections](#collections)
- [Users](#users)
- [Reactions](#reactions)
- [Enrichment](#enrichment)
- [License](#license)
- [We are hiring!](#we-are-hiring)

## Installation

Get the client:

```bash
$ go get github.com/GetStream/stream-go2/v7
```

> For `v4`, use `github.com/GetStream/stream-go2/v4` but updating to the last version is highly recommended.

### Creating a Client

```go
import stream "github.com/GetStream/stream-go2/v7"

key := "YOUR_API_KEY"
secret := "YOUR_API_SECRET"

client, err := stream.New(key, secret)
if err != nil {
    // ...
}
```

You can pass additional options when creating a client using the available `ClientOption` functions:

```go
client, err := stream.NewClient(key, secret,
    stream.WithAPIRegion("us-east"),
    stream.WithAPIVersion("1.0"),
    stream.WithTimeout(5 * time.Second),
    ...,
)
```

You can also create a client using environment variables:

```go
client, err := stream.NewFromEnv()
```

Available environment variables:

* `STREAM_API_KEY`
* `STREAM_API_SECRET`
* `STREAM_API_REGION`
* `STREAM_API_VERSION`

### Rate Limits

API has different rate limits for each distinct endpoint and this information is returned to the client in response headers and SDK parses headers into `Rate` type.

This info is provided to the user in 2 ways:
* in responses; `resp, _ := feed.GetActivities; resp.Rate`.
* in errors; if request doesn't succeed then error is a type of `APIError` and headers are accessible via `err.(APIError).Rate`.

### Creating a Feed

Create a flat feed from slug and user ID:

```go
flat, err := client.FlatFeed("user", "123")
```

Create an aggregated feed from slug and user ID:

```go
aggr, err := client.AggregatedFeed("aggregated", "123")
```

Create a notification feed from slug and user ID:

```go
notif, err := client.NotificationFeed("notification", "123")
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
fmt.Println("Rate:", resp.Rate)
fmt.Println("Next:", resp.Next)
fmt.Println("Activities:")
for _, activity := range resp.Results {
    fmt.Println(activity)
}
```

You can retrieve flat feeds with [custom ranking](https://getstream.io/docs/#custom_ranking), using the dedicated method:

```go
resp, err := flat.GetActivitiesWithRanking(ctx, "popularity")
if err != nil {
    // ...
}
```

#### Aggregated feeds

```go
resp, err := aggregated.GetActivities(ctx)
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Rate:", resp.Rate)
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
resp, err := notification.GetActivities(ctx)
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Rate:", resp.Rate)
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
    ctx,
    stream.WithActivitiesIDGTE("f505b3fb-a212-11e7-..."),
    stream.WithActivitiesLimit(5),
    ...,
)
```

### Adding activities

Add a single activity:

```go
resp, err := feed.AddActivity(ctx, stream.Activity{Actor: "bob", ...})
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Rate:", resp.Rate)
fmt.Println("Activity:", resp.Activity) // resp wraps the stream.Activity type
```

Add multiple activities:

```go
a1 := stream.Activity{Actor: "bob", ...}
a2 := stream.Activity{Actor: "john", ...}
a3 := stream.Activity{Actor: "alice", ...}

resp, err := feed.AddActivities(ctx, a1, a2, a3)
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Rate:", resp.Rate)
fmt.Println("Activities:")
for _, activity := range resp.Activities {
    fmt.Println(activity)
}
```

### Updating activities

```go
_, err := feed.UpdateActivities(ctx, a1, a2, ...)
if err != nil {
    // ...
}
```

### Partially updating activities

You can partial update activities identified either by ID:

``` go
changesetA := stream.NewUpdateActivityRequestByID("f505b3fb-a212-11e7-...", map[string]interface{}{"key": "new-value"}, []string{"removed", "keys"})
changesetB := stream.NewUpdateActivityRequestByID("f707b3fb-a212-11e7-...", map[string]interface{}{"key": "new-value"}, []string{"removed", "keys"})
resp, err := client.PartialUpdateActivities(ctx, changesetA, changesetB)
if err != nil {
    // ...
}
```

or by a ForeignID and timestamp pair:
``` go
changesetA := stream.NewUpdateActivityRequestByForeignID("dothings:1", stream.Time{...}, map[string]interface{}{"key": "new-value"}, []string{"removed", "keys"})
changesetB := stream.NewUpdateActivityRequestByForeignID("dothings:2", stream.Time{...}, map[string]interface{}{"key": "new-value"}, []string{"removed", "keys"})
resp, err := client.PartialUpdateActivities(ctx, changesetA, changesetB)
if err != nil {
    // ...
}
```

### Removing activities

You can either remove activities by ID or ForeignID:

```go
_, err := feed.RemoveActivityByID(ctx, "f505b3fb-a212-11e7-...")
if err != nil {
    // ...
}

_, err := feed.RemoveActivityByForeignID(ctx, "bob:123")
if err != nil {
    // ...
}
```

### Following another feed

```go
_, err := feed.Follow(ctx, anotherFeed)
if err != nil {
    // ...
}
```

Beware that it's possible to follow only flat feeds.

#### Options

You can pass options to the `Follow` method. For example:

```go
_, err := feed.Follow(ctx,
    anotherFeed,
    stream.WithFollowFeedActivityCopyLimit(15),
    ...,
)
```

### Retrieving followers and followings

#### Following

Get the feeds that a feed is following:

```go
resp, err := feed.GetFollowing(ctx)
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
    ctx,
    stream.WithFollowingLimit(5),
    ...,
)
```

#### Followers

```go
resp, err := flat.GetFollowers(ctx)
if err != nil {
    // ...
}

fmt.Println("Duration:", resp.Duration)
fmt.Println("Rate:", resp.Rate)
for _, follower := range resp.Results {
    fmt.Println(follower.FeedID, follower.TargetID)
}
```

Note: this is only possible for `FlatFeed` types.

You can pass options to `GetFollowers`:

```go
resp, err := feed.GetFollowing(
    ctx,
    stream.WithFollowersLimit(5),
    ...,
)
```

### Unfollowing a feed

```go
_, err := flat.Unfollow(ctx, anotherFeed)
if err != nil {
    // ...
}
```

You can pass options to `Unfollow`:

```go
_, err := flat.Unfollow(ctx,
    anotherFeed,
    stream.WithUnfollowKeepHistory(true),
    ...,
)
```

### Updating an activity's `to` targets

Remove all old targets and set new ones (replace):

```go
newTargets := []stream.Feed{f1, f2}

_, err := feed.UpdateToTargets(ctx, activity, stream.WithToTargetsNew(newTargets...))
if err != nil {
    // ...
}
```

Add some targets and remove some others:

```go
add := []stream.Feed{target1, target2}
remove := []stream.Feed{oldTarget1, oldTarget2}

_, err := feed.UpdateToTargets(
    ctx,
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
_, err := client.AddToMany(ctx,
    activity, feed1, feed2, ...,
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

_, err := client.FollowMany(ctx, relationships)
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
_, err := analytics.TrackEngagement(ctx, event)
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
_, err := analytics.TrackImpression(ctx, imp)
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
resp, err = personalization.Get(ctx, "follow_recommendations", data)
if err != nil {
    // ...
}
fmt.Println(resp)
```

See the complete [docs and examples](https://getstream.io/docs/#personalization_introduction) about personalization features on Stream's documentation pages.

## Collections

[Collections](https://getstream.io/docs/#collections_introduction) endpoints can be reached using a specialized `CollectionsClient` which, like `PersonalizationClient`, can be obtained from a regular `Client`:

```go
collections := client.Collections()
```

`CollectionsClient` exposes three batch functions, `Upsert`, `Select`, and `DeleteMany` as well as CRUD functions: `Add`, `Get`, `Update`, `Delete`:

```go
// Upsert the "picture" collection
object := stream.CollectionObject{
    ID:   "123",
    Data: map[string]interface{}{
        "name": "Rocky Mountains",
        "location": "North America",
    },
}
_, err = collections.Upsert(ctx, "picture", object)
if err != nil {
    // ...
}

// Get the data from the "picture" collection for ID "123" and "456"
objects, err := collections.Select(ctx, "picture", "123", "456")
if err != nil {
    // ...
}

// Delete the data from the "picture" collection for picture with ID "123"
_, err = collections.Delete(ctx, "picture", "123")
if err != nil {
    // ...
}

// Get a single collection object from the "pictures" collection with ID "123"
_, err = collections.Get(ctx, "pictures", "123")
if err != nil {
    // ...
}
```

See the complete [docs and examples](https://getstream.io/docs/#collections_introduction) about collections on Stream's documentation pages.

## Users

[Users](https://getstream.io/docs/#users_introduction) endpoints can be reached using a specialized `UsersClient` which, like `CollectionsClient`, can be obtained from a regular `Client`:

```go
users := client.Users()
```

`UsersClient` exposes CRUD functions: `Add`, `Get`, `Update`, `Delete`:

```go
user := stream.User{
    ID: "123",
    Data: map[string]interface{}{
        "name": "Bobby Tables",
    },
}

insertedUser, err := users.Add(ctx, user, false)
if err != nil {
    // ...
}

newUserData :=map[string]interface{}{
    "name": "Bobby Tables",
    "age": 7,
}

updatedUser, err := users.Update(ctx, "123", newUserData)
if err != nil {
    // ...
}

_, err = users.Delete(ctx, "123")
if err != nil {
    // ...
}
```

See the complete [docs and examples](https://getstream.io/docs/#users_introduction) about users on Stream's documentation pages.

## Reactions

[Reactions](https://getstream.io/docs/#reactions_introduction) endpoints can be reached using a specialized `Reactions` which, like `CollectionsClient`, can be obtained from a regular `Client`:

```go
reactions := client.Reactions()
```

`Reactions` exposes CRUD functions: `Add`, `Get`, `Update`, `Delete`, as well as two specialized functions `AddChild` and `Filter`:

```go
r := stream.AddReactionRequestObject{
    Kind: "comment",
    UserID: "123",
    ActivityID: "87a9eec0-fd5f-11e8-8080-80013fed2f5b",
    Data: map[string]interface{}{
        "text": "Nice post!!",
    },
    TargetFeeds: []string{"user:bob", "timeline:alice"},
}

comment, err := reactions.Add(ctx, r)
if err != nil {
    // ...
}

like := stream.AddReactionRequestObject{
    Kind: "like",
    UserID: "456",
}

childReaction, err := reactions.AddChild(ctx, comment.ID, like)
if err != nil {
    // ...
}

// If we fetch the "comment" reaction now, it will have the child reaction(s) present.
parent, err := reactions.Get(ctx, comment.ID)
if err != nil {
    // ...
}

for kind, children := range parent.ChildrenReactions {
    // child reactions are grouped by kind
}

// update the target feeds for the `comment` reaction
updatedReaction, err := reactions.Update(ctx, comment.ID, nil, []string{"timeline:jeff"})
if err != nil {
    // ...
}

// get all reactions for the activity "87a9eec0-fd5f-11e8-8080-80013fed2f5b",
// paginated 5 at a time, including the activity data
response, err := reactions.Filter(ctx,
    stream.ByActivityID("87a9eec0-fd5f-11e8-8080-80013fed2f5b"),
    stream.WithLimit(5),
    stream.WithActivityData())
if err != nil {
    // ...
}

// since we requested the activity, it will be present in the response
fmt.Println(response.Activity)

for _, reaction := range response.Results{
    // do something for each reaction
}

//get the next page of reactions
response, err = reactions.GetNextPageFilteredReactions(ctx, response)
if err != nil {
    // ...
}

// get all likes by user "123"
response, err = reactions.Filter(ctx, stream.ByUserID("123").ByKind("like"))
if err != nil {
    // ...
}
```

See the complete [docs and examples](https://getstream.io/docs/#reactions_introduction) about reactions on Stream's documentation pages.

## Enrichment

[Enrichment](https://getstream.io/docs/#enrichment_introduction) is a way of retrieving activities from feeds in which references to `Users` and `Collections` will be replaced with the corresponding objects

`FlatFeed`, `AggregatedFeed` and `NotificationFeed` each have a `GetEnrichedActivities` function to retrieve enriched activities.

```go
u := stream.User{
    ID: "123",
    Data: map[string]interface{}{
        "name": "Bobby Tables",
    },
}

// We add a user
user, err := client.Users().Add(ctx, u, true)
if err != nil {
    // ...
}

c := stream.CollectionObject{
    ID: "123",
    Data: map[string]interface{}{
        "name":     "Rocky Mountains",
        "location": "North America",
    },
}

// We add a collection object
collectionObject, err := client.Collections().Add(ctx, "picture", c)
if err != nil {
    // ...
}

act := stream.Activity{
    Time:      stream.Time{Time: time.Now()},
    Actor:     client.Users().CreateReference(user.ID),
    Verb:      "post",
    Object:    client.Collections().CreateReference("picture", collectionObject.ID),
    ForeignID: "picture:1",
}


// We add the activity to the user's feed
feed, _ := client.FlatFeed("user", "123")
_, err = feed.AddActivity(ctx, act)
if err != nil {
    // ...
}

result, err := feed.GetActivities(ctx)
if err != nil {
    // ...
}
fmt.Println(result.Results[0].Actor) // Will output the user reference
fmt.Println(result.Results[0].Object) // Will output the collection reference

enrichedResult, err := feed.GetEnrichedActivities(ctx)
if err != nil {
    // ...
}
fmt.Println(enrichedResult.Results[0]["actor"].(map[string]interface{})) // Will output the user object
fmt.Println(enrichedResult.Results[0]["object"].(map[string]interface{})) // Will output the collection object
```

See the complete [docs and examples](https://getstream.io/docs/#enrichment_introduction) about enrichment on Stream's documentation pages.

## License

Project is licensed under the [BSD 3-Clause](LICENSE).

## We are hiring!

We've recently closed a [$38 million Series B funding round](https://techcrunch.com/2021/03/04/stream-raises-38m-as-its-chat-and-activity-feed-apis-power-communications-for-1b-users/) and we keep actively growing.
Our APIs are used by more than a billion end-users, and you'll have a chance to make a huge impact on the product within a team of the strongest engineers all over the world.

Check out our current openings and apply via [Stream's website](https://getstream.io/team/#jobs).
