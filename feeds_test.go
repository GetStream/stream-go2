package stream_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	stream "github.com/reifcode/stream-go2"
)

func TestClient(t *testing.T) {
	client, err := stream.NewClient(os.Getenv("STREAM_API_KEY"), os.Getenv("STREAM_API_SECRET"))
	require.Nil(t, err)

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	flatFeed := client.FlatFeed("flat", id)
	fmt.Println("slug:", flatFeed.Slug(), "id:", flatFeed.ID())

	ciccio := stream.Activity{Actor: "ciccio", Verb: "like", Object: "ice cream"}
	resp, err := flatFeed.AddActivities(ciccio)
	require.Nil(t, err)
	ciccio.ID = resp.Activities[0].ID
	fmt.Println("added activities in:", resp.Duration)

	out, err := flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}
	fmt.Println()

	now := time.Now()
	pippo := stream.Activity{Actor: "pippo", Verb: "like", Object: "cats", ForeignID: "pippo:smartass", Time: now}
	_, err = flatFeed.AddActivities(pippo)
	require.Nil(t, err)

	out, err = flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, " - ", a.ForeignID, "]")
	}
	fmt.Println()

	pippo.Object = "dogs"
	err = flatFeed.UpdateActivities(pippo)
	require.Nil(t, err)

	out, err = flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}
	fmt.Println()

	fmt.Println("removing activity", ciccio.ID)
	err = flatFeed.RemoveActivityByID(ciccio.ID)
	require.Nil(t, err)
	out, err = flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}

	fmt.Println("removing activity by foreignID", pippo.ForeignID)
	err = flatFeed.RemoveActivityByForeignID(pippo.ForeignID)
	require.Nil(t, err)
	out, err = flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}

	anotherFeed := client.FlatFeed("ranked", id)
	_, err = anotherFeed.AddActivities(stream.Activity{Actor: "something", Verb: "does", Object: "something"})
	require.Nil(t, err)

	fmt.Println("following", anotherFeed.Slug(), anotherFeed.ID())
	err = flatFeed.Follow(anotherFeed)
	require.Nil(t, err)

	fmt.Println("followers:")
	followersResp, err := anotherFeed.GetFollowers()
	require.Nil(t, err)
	for i, f := range followersResp.Results {
		fmt.Println(i, f.FeedID, f.TargetID)
	}
	fmt.Println()

	fmt.Println("followings:")
	followingResp, err := flatFeed.GetFollowing()
	require.Nil(t, err)
	for i, f := range followingResp.Results {
		fmt.Println(i, f.FeedID, f.TargetID)
	}

	_, err = anotherFeed.AddActivities(stream.Activity{Actor: "daoisj", Verb: "random", Object: "wut"})
	require.Nil(t, err)
	out, err = flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}
	fmt.Println()

	fmt.Println("unfollowing")
	err = flatFeed.Unfollow(anotherFeed, stream.WithKeepHistory(true))
	require.Nil(t, err)
	out, err = flatFeed.GetActivities()
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}
	_, err = anotherFeed.AddActivities(stream.Activity{Actor: "daoisj", Verb: "again", Object: "wut"})
	require.Nil(t, err)
	out, err = flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}
	fmt.Println()

	fmt.Println("add to many")
	err = client.AddToMany(stream.Activity{Actor: "attore", Verb: "multipla", Object: "multiplobj"}, flatFeed, anotherFeed)
	require.Nil(t, err)

	out, err = flatFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}
	fmt.Println()

	out, err = anotherFeed.GetActivities()
	require.Nil(t, err)
	fmt.Println("read activities:", out.Duration, len(out.Results))
	for i, a := range out.Results {
		fmt.Println(i, "->", a.Actor, a.Verb, a.Object, "[", a.ID, "]")
	}
	fmt.Println()

	fmt.Println("follow many")
	err = client.FollowMany([]stream.FollowRelationship{stream.NewFollowRelationship(anotherFeed, flatFeed)}, stream.WithActivityCopyLimitQuery(0))
	require.NoError(t, err)
}
