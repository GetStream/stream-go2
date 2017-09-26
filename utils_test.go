package stream_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	stream "github.com/reifcode/stream-go2"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newClient(t *testing.T) *stream.Client {
	client, err := stream.NewClient(os.Getenv("STREAM_API_KEY"), os.Getenv("STREAM_API_SECRET"))
	require.NoError(t, err)
	return client
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func getTime(t time.Time) stream.Time {
	st, _ := time.Parse(stream.TimeLayout, t.Truncate(time.Second).Format(stream.TimeLayout))
	return stream.Time{Time: st}
}

func newFlatFeed(c *stream.Client) *stream.FlatFeed {
	return c.FlatFeed("flat", randString(10))
}

func newAggregatedFeed(c *stream.Client) *stream.AggregatedFeed {
	return c.AggregatedFeed("timeline_aggregated", randString(10))
}
