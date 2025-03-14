package stream

import (
	"context"
	"encoding/json"
	"time"
)

type AuditLogsClient struct {
	client *Client
}

type AuditLog struct {
	EntityType string         `json:"entity_type"`
	EntityID   string         `json:"entity_id"`
	Action     string         `json:"action"`
	UserID     string         `json:"user_id"`
	Custom     map[string]any `json:"custom"`
	CreatedAt  time.Time      `json:"created_at"`
}

type QueryAuditLogsResponse struct {
	AuditLogs []AuditLog `json:"audit_logs"`
	Next      string     `json:"next"`
	Prev      string     `json:"prev"`
	response
}

type QueryAuditLogsFilters struct {
	EntityType string
	EntityID   string
	UserID     string
}

type QueryAuditLogsPager struct {
	Next  string
	Prev  string
	Limit int
}

func (c *AuditLogsClient) QueryAuditLogs(ctx context.Context, filters QueryAuditLogsFilters, pager QueryAuditLogsPager) (*QueryAuditLogsResponse, error) {
	endpoint := c.client.makeEndpoint("audit_logs/")
	if filters.EntityType != "" && filters.EntityID != "" {
		endpoint.addQueryParam(makeRequestOption("entity_type", filters.EntityType))
		endpoint.addQueryParam(makeRequestOption("entity_id", filters.EntityID))
	}
	if filters.UserID != "" {
		endpoint.addQueryParam(makeRequestOption("user_id", filters.UserID))
	}
	if pager.Next != "" {
		endpoint.addQueryParam(makeRequestOption("next", pager.Next))
	}
	if pager.Prev != "" {
		endpoint.addQueryParam(makeRequestOption("prev", pager.Prev))
	}
	if pager.Limit > 0 {
		endpoint.addQueryParam(makeRequestOption("limit", pager.Limit))
	}
	body, err := c.client.get(ctx, endpoint, nil, c.client.authenticator.auditLogsAuth)
	if err != nil {
		return nil, err
	}

	var resp QueryAuditLogsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
