package stream_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlagActivity(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	err := client.Moderation().FlagActivity(ctx, "jimmy", "foo", "reason1")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/flag/?api_key=key",
		`{"entity_id":"foo", "entity_type":"stream:feeds:v2:activity", "reason":"reason1", "user_id":"jimmy"}`,
	)
}

func TestFlagReaction(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	err := client.Moderation().FlagReaction(ctx, "jimmy", "foo", "reason1")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/flag/?api_key=key",
		`{"entity_id":"foo", "entity_type":"stream:feeds:v2:reaction", "reason":"reason1", "user_id":"jimmy"}`,
	)
}

func TestFlagUser(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	err := client.Moderation().FlagUser(ctx, "jimmy", "foo", "reason1")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/flag/?api_key=key",
		`{"entity_id":"foo", "entity_type":"stream:user", "reason":"reason1", "user_id":"jimmy"}`,
	)
}

func TestUpdateActivityModerationStatus(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	err := client.Moderation().UpdateActivityModerationStatus(ctx, "foo", "complete", "watch", "mark_safe")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/status/?api_key=key",
		`{"entity_id":"foo", "entity_type":"stream:feeds:v2:activity", "latest_moderator_action":"mark_safe", "recommended_action":"watch", "status":"complete"}`,
	)
}

func TestUpdateReactionModerationStatus(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	err := client.Moderation().UpdateReactionModerationStatus(ctx, "foo", "complete", "watch", "mark_safe")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/status/?api_key=key",
		`{"entity_id":"foo", "entity_type":"stream:feeds:v2:reaction", "latest_moderator_action":"mark_safe", "recommended_action":"watch", "status":"complete"}`,
	)
}
