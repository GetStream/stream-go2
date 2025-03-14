package stream_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v8"
)

func TestQueryAuditLogs(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	// Test with all filters and pager options
	filters := stream.QueryAuditLogsFilters{
		EntityType: "feed",
		EntityID:   "123",
		UserID:     "user-42",
	}
	pager := stream.QueryAuditLogsPager{
		Next:  "next-token",
		Prev:  "prev-token",
		Limit: 25,
	}

	// Set mock response
	now := time.Now()
	mockResp := struct {
		AuditLogs []stream.AuditLog `json:"audit_logs"`
		Next      string            `json:"next"`
		Prev      string            `json:"prev"`
	}{
		AuditLogs: []stream.AuditLog{
			{
				EntityType: "feed",
				EntityID:   "123",
				Action:     "create",
				UserID:     "user-42",
				Custom:     map[string]any{"key": "value"},
				CreatedAt:  now,
			},
		},
		Next: "next-page-token",
		Prev: "prev-page-token",
	}
	respBytes, err := json.Marshal(mockResp)
	require.NoError(t, err)
	requester.resp = string(respBytes)

	// Call the function
	resp, err := client.AuditLogs().QueryAuditLogs(ctx, filters, pager)
	require.NoError(t, err)

	// Verify request
	testRequest(
		t,
		requester.req,
		http.MethodGet,
		"https://api.stream-io-api.com/api/v1.0/audit_logs/?api_key=key&entity_id=123&entity_type=feed&limit=25&next=next-token&prev=prev-token&user_id=user-42",
		"",
	)

	// Verify response
	assert.Len(t, resp.AuditLogs, 1)
	assert.Equal(t, "feed", resp.AuditLogs[0].EntityType)
	assert.Equal(t, "123", resp.AuditLogs[0].EntityID)
	assert.Equal(t, "create", resp.AuditLogs[0].Action)
	assert.Equal(t, "user-42", resp.AuditLogs[0].UserID)
	assert.Equal(t, "value", resp.AuditLogs[0].Custom["key"])
	assert.Equal(t, now.Truncate(time.Second).UTC(), resp.AuditLogs[0].CreatedAt.Truncate(time.Second).UTC())
	assert.Equal(t, "next-page-token", resp.Next)
	assert.Equal(t, "prev-page-token", resp.Prev)
}

func TestQueryAuditLogsWithMinimalParams(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	// Test with minimal filters and pager options
	filters := stream.QueryAuditLogsFilters{}
	pager := stream.QueryAuditLogsPager{}

	// Set mock response
	mockResp := struct {
		AuditLogs []stream.AuditLog `json:"audit_logs"`
		Next      string            `json:"next"`
		Prev      string            `json:"prev"`
	}{
		AuditLogs: []stream.AuditLog{},
		Next:      "",
		Prev:      "",
	}
	respBytes, err := json.Marshal(mockResp)
	require.NoError(t, err)
	requester.resp = string(respBytes)

	// Call the function
	resp, err := client.AuditLogs().QueryAuditLogs(ctx, filters, pager)
	require.NoError(t, err)

	// Verify request
	testRequest(
		t,
		requester.req,
		http.MethodGet,
		"https://api.stream-io-api.com/api/v1.0/audit_logs/?api_key=key",
		"",
	)

	// Verify response
	assert.Empty(t, resp.AuditLogs)
	assert.Empty(t, resp.Next)
	assert.Empty(t, resp.Prev)
}
