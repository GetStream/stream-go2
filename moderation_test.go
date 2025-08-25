package stream_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v8"
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

func TestUpdateStatusBatch(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	batchReq := stream.UpdateStatusBatchRequest{
		ModeratorID: "moderator_123",
		Updates: []stream.UpdateStatusBatchItem{
			{
				EntityType:            "stream:feeds:v2:activity",
				EntityID:              "activity_1",
				Status:                "complete",
				RecommendedAction:     "watch",
				LatestModeratorAction: "mark_safe",
			},
			{
				EntityType:            "stream:feeds:v2:reaction",
				EntityID:              "reaction_1",
				Status:                "complete",
				RecommendedAction:     "flag",
				LatestModeratorAction: "mark_harmful",
			},
		},
	}

	// Mock the response
	requester.resp = `{
		"results": [
			{
				"entity_id": "activity_1",
				"entity_type": "stream:feeds:v2:activity",
				"success": true
			},
			{
				"entity_id": "reaction_1",
				"entity_type": "stream:feeds:v2:reaction",
				"success": true
			}
		]
	}`

	resp, err := client.Moderation().UpdateStatusBatch(ctx, batchReq)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Results, 2)
	require.Equal(t, "activity_1", resp.Results[0].EntityID)
	require.Equal(t, "stream:feeds:v2:activity", resp.Results[0].EntityType)
	require.True(t, resp.Results[0].Success)
	require.Equal(t, "reaction_1", resp.Results[1].EntityID)
	require.Equal(t, "stream:feeds:v2:reaction", resp.Results[1].EntityType)
	require.True(t, resp.Results[1].Success)

	testRequest(
		t,
		requester.req,
		http.MethodPost,
		"https://api.stream-io-api.com/api/v1.0/moderation/status/batch/?api_key=key",
		`{"moderator_id":"moderator_123","updates":[{"entity_type":"stream:feeds:v2:activity","entity_id":"activity_1","status":"complete","recommended_action":"watch","latest_moderator_action":"mark_safe"},{"entity_type":"stream:feeds:v2:reaction","entity_id":"reaction_1","status":"complete","recommended_action":"flag","latest_moderator_action":"mark_harmful"}]}`,
	)
}

func TestUpdateStatusBatchValidation(t *testing.T) {
	ctx := context.Background()
	client, _ := newClient(t)

	// Test missing moderator ID
	t.Run("missing moderator id", func(t *testing.T) {
		batchReq := stream.UpdateStatusBatchRequest{
			ModeratorID: "", // Missing
			Updates: []stream.UpdateStatusBatchItem{
				{
					EntityType: "stream:feeds:v2:activity",
					EntityID:   "activity_1",
					Status:     "complete",
				},
			},
		}
		_, err := client.Moderation().UpdateStatusBatch(ctx, batchReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "moderator_id is required")
	})

	// Test empty updates
	t.Run("empty updates", func(t *testing.T) {
		batchReq := stream.UpdateStatusBatchRequest{
			ModeratorID: "moderator_123",
			Updates:     []stream.UpdateStatusBatchItem{},
		}
		_, err := client.Moderation().UpdateStatusBatch(ctx, batchReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "at least one update is required")
	})

	// Test too many updates (over 100)
	t.Run("too many updates", func(t *testing.T) {
		updates := make([]stream.UpdateStatusBatchItem, 101)
		for i := 0; i < 101; i++ {
			updates[i] = stream.UpdateStatusBatchItem{
				EntityType: "stream:feeds:v2:activity",
				EntityID:   fmt.Sprintf("activity_%d", i),
				Status:     "complete",
			}
		}
		batchReq := stream.UpdateStatusBatchRequest{
			ModeratorID: "moderator_123",
			Updates:     updates,
		}
		_, err := client.Moderation().UpdateStatusBatch(ctx, batchReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum of 100 updates allowed, got 101")
	})

	// Test duplicate entity type/ID combination
	t.Run("duplicate entities", func(t *testing.T) {
		batchReq := stream.UpdateStatusBatchRequest{
			ModeratorID: "moderator_123",
			Updates: []stream.UpdateStatusBatchItem{
				{
					EntityType: "stream:feeds:v2:activity",
					EntityID:   "activity_1",
					Status:     "complete",
				},
				{
					EntityType: "stream:feeds:v2:activity",
					EntityID:   "activity_1", // Same as above
					Status:     "pending",
				},
			},
		}
		_, err := client.Moderation().UpdateStatusBatch(ctx, batchReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate entity found: entity_type=stream:feeds:v2:activity, entity_id=activity_1")
	})

	// Test missing entity type
	t.Run("missing entity type", func(t *testing.T) {
		batchReq := stream.UpdateStatusBatchRequest{
			ModeratorID: "moderator_123",
			Updates: []stream.UpdateStatusBatchItem{
				{
					EntityType: "", // Missing
					EntityID:   "activity_1",
					Status:     "complete",
				},
			},
		}
		_, err := client.Moderation().UpdateStatusBatch(ctx, batchReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "entity_type is required for update at index 0")
	})

	// Test missing entity ID
	t.Run("missing entity id", func(t *testing.T) {
		batchReq := stream.UpdateStatusBatchRequest{
			ModeratorID: "moderator_123",
			Updates: []stream.UpdateStatusBatchItem{
				{
					EntityType: "stream:feeds:v2:activity",
					EntityID:   "", // Missing
					Status:     "complete",
				},
			},
		}
		_, err := client.Moderation().UpdateStatusBatch(ctx, batchReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "entity_id is required for update at index 0")
	})
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
