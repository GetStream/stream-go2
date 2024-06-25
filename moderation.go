package stream

import (
	"context"
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
	Status                string `json:"status"`
	RecommendedAction     string `json:"recommended_action"`
	LatestModeratorAction string `json:"latest_moderator_action"`
}

func (c *ModerationClient) UpdateActivityModerationStatus(ctx context.Context, activityID, status, recAction, modAction string) error {
	r := updateStatusRequest{
		EntityType:            ModerationActivity,
		EntityID:              activityID,
		Status:                status,
		RecommendedAction:     recAction,
		LatestModeratorAction: modAction,
	}
	return c.updateStatus(ctx, r)
}

func (c *ModerationClient) UpdateReactionModerationStatus(ctx context.Context, reactionID, status, recAction, modAction string) error {
	r := updateStatusRequest{
		EntityType:            ModerationReaction,
		EntityID:              reactionID,
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

func (c *ModerationClient) InvalidateUserCache(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("empty userID")
	}

	endpoint := c.client.makeEndpoint("moderation/user/cache/%s/", userID)

	_, err := c.client.delete(ctx, endpoint, nil, c.client.authenticator.moderationAuth)
	return err
}
