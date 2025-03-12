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
	err := client.Moderation().UpdateActivityModerationStatus(ctx, "foo", "moderator_123", "complete", "watch", "mark_safe")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/status/?api_key=key",
		`{"entity_id":"foo", "entity_type":"stream:feeds:v2:activity", "moderator_id": "moderator_123", "latest_moderator_action":"mark_safe", "recommended_action":"watch", "status":"complete"}`,
	)
}

func TestUpdateReactionModerationStatus(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	err := client.Moderation().UpdateReactionModerationStatus(ctx, "foo", "moderator_123", "complete", "watch", "mark_safe")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/status/?api_key=key",
		`{"entity_id":"foo", "entity_type":"stream:feeds:v2:reaction", "moderator_id": "moderator_123", "latest_moderator_action":"mark_safe", "recommended_action":"watch", "status":"complete"}`,
	)
}

func TestInvalidateUserCache(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	err := client.Moderation().InvalidateUserCache(ctx, "foo")
	require.NoError(t, err)
	testRequest(
		t,
		requester.req,
		http.MethodDelete,
		"https://api.stream-io-api.com/api/v1.0/moderation/user/cache/foo/?api_key=key",
		"",
	)

	err = client.Moderation().InvalidateUserCache(ctx, "")
	require.Error(t, err)
	require.Equal(t, "empty userID", err.Error())
}
