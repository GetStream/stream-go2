package stream

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	ModerationActivity = "stream:feeds:v2:activity"
	ModerationReaction = "stream:feeds:v2:reaction"
	ModerationUser     = "stream:user"
)

type ModerationClient struct {
	client *Client
}

type flagRequest struct {
	// the user doing the reporting
	UserID     string `json:"user_id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Reason     string `json:"reason"`
}

func (c *ModerationClient) FlagActivity(ctx context.Context, userID, activityID, reason string) error {
	r := flagRequest{
		UserID:     userID,
		EntityType: ModerationActivity,
		EntityID:   activityID,
		Reason:     reason,
	}
	return c.flagContent(ctx, r)
}

func (c *ModerationClient) FlagReaction(ctx context.Context, userID, reactionID, reason string) error {
	r := flagRequest{
		UserID:     userID,
		EntityType: ModerationReaction,
		EntityID:   reactionID,
		Reason:     reason,
	}
	return c.flagContent(ctx, r)
}

func (c *ModerationClient) FlagUser(ctx context.Context, userID, targetUserID, reason string) error {
	r := flagRequest{
		UserID:     userID,
		EntityType: ModerationUser,
		EntityID:   targetUserID,
		Reason:     reason,
	}
	return c.flagContent(ctx, r)
}

func (c *ModerationClient) flagContent(ctx context.Context, r flagRequest) error {
	endpoint := c.client.makeEndpoint("moderation/flag/")

	_, err := c.client.post(ctx, endpoint, r, c.client.authenticator.moderationAuth)
	return err
}

type updateStatusRequest struct {
	EntityType            string `json:"entity_type"`
	EntityID              string `json:"entity_id"`
	ModeratorID           string `json:"moderator_id"`
	Status                string `json:"status"`
	RecommendedAction     string `json:"recommended_action"`
	LatestModeratorAction string `json:"latest_moderator_action"`
}

func (c *ModerationClient) UpdateActivityModerationStatus(ctx context.Context, activityID, modID, status, recAction, modAction string) error {
	r := updateStatusRequest{
		EntityType:            ModerationActivity,
		EntityID:              activityID,
		ModeratorID:           modID,
		Status:                status,
		RecommendedAction:     recAction,
		LatestModeratorAction: modAction,
	}
	return c.updateStatus(ctx, r)
}

func (c *ModerationClient) UpdateReactionModerationStatus(ctx context.Context, reactionID, modID, status, recAction, modAction string) error {
	r := updateStatusRequest{
		EntityType:            ModerationReaction,
		EntityID:              reactionID,
		ModeratorID:           modID,
		Status:                status,
		RecommendedAction:     recAction,
		LatestModeratorAction: modAction,
	}
	return c.updateStatus(ctx, r)
}

func (c *ModerationClient) updateStatus(ctx context.Context, r updateStatusRequest) error {
	endpoint := c.client.makeEndpoint("moderation/status/")

	_, err := c.client.post(ctx, endpoint, r, c.client.authenticator.moderationAuth)
	return err
}

type UpdateStatusBatchRequest struct {
	ModeratorID string                  `json:"moderator_id"`
	Updates     []UpdateStatusBatchItem `json:"updates"`
}

type UpdateStatusBatchItem struct {
	EntityType            string `json:"entity_type"`
	EntityID              string `json:"entity_id"`
	Status                string `json:"status"`
	RecommendedAction     string `json:"recommended_action"`
	LatestModeratorAction string `json:"latest_moderator_action"`
}

type UpdateStatusBatchResponse struct {
	BaseResponse
	Results []UpdateStatusBatchResult `json:"results"`
}

type UpdateStatusBatchResult struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

func (c *ModerationClient) UpdateStatusBatch(ctx context.Context, req UpdateStatusBatchRequest) (*UpdateStatusBatchResponse, error) {
	// Validate the request
	if err := validateUpdateStatusBatchRequest(req); err != nil {
		return nil, err
	}

	endpoint := c.client.makeEndpoint("moderation/status/batch/")

	resp, err := c.client.post(ctx, endpoint, req, c.client.authenticator.moderationAuth)
	if err != nil {
		return nil, err
	}

	var result UpdateStatusBatchResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func validateUpdateStatusBatchRequest(req UpdateStatusBatchRequest) error {
	// Check moderator ID
	if req.ModeratorID == "" {
		return fmt.Errorf("moderator_id is required")
	}

	// Check maximum limit
	if len(req.Updates) > 100 {
		return fmt.Errorf("maximum of 100 updates allowed, got %d", len(req.Updates))
	}

	// Check for empty updates
	if len(req.Updates) == 0 {
		return fmt.Errorf("at least one update is required")
	}

	// Check for duplicate entity type/ID combinations
	seen := make(map[string]bool)
	for i, update := range req.Updates {
		if update.EntityType == "" {
			return fmt.Errorf("entity_type is required for update at index %d", i)
		}
		if update.EntityID == "" {
			return fmt.Errorf("entity_id is required for update at index %d", i)
		}

		key := update.EntityType + ":" + update.EntityID
		if seen[key] {
			return fmt.Errorf("duplicate entity found: entity_type=%s, entity_id=%s", update.EntityType, update.EntityID)
		}
		seen[key] = true
	}

	return nil
}

func (c *ModerationClient) InvalidateUserCache(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("empty userID")
	}

	endpoint := c.client.makeEndpoint("moderation/user/cache/%s/", userID)

	_, err := c.client.delete(ctx, endpoint, nil, c.client.authenticator.moderationAuth)
	return err
}
