package stream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// ReactionsClient is a specialized client used to interact with the Reactions endpoints.
type ReactionsClient struct {
	client *Client
}

// Add adds a reaction.
func (c *ReactionsClient) Add(ctx context.Context, r AddReactionRequestObject) (*ReactionResponse, error) {
	if r.ParentID != "" {
		return nil, errors.New("`Parent` not empty. For adding child reactions use `AddChild`")
	}
	return c.addReaction(ctx, r)
}

// AddChild adds a child reaction to the provided parent.
func (c *ReactionsClient) AddChild(ctx context.Context, parentID string, r AddReactionRequestObject) (*ReactionResponse, error) {
	r.ParentID = parentID
	return c.addReaction(ctx, r)
}

func (c *ReactionsClient) addReaction(ctx context.Context, r AddReactionRequestObject) (*ReactionResponse, error) {
	endpoint := c.client.makeEndpoint("reaction/")
	return c.decode(c.client.post(ctx, endpoint, r, c.client.authenticator.reactionsAuth))
}

func (c *ReactionsClient) decode(resp []byte, err error) (*ReactionResponse, error) {
	if err != nil {
		return nil, err
	}

	var result ReactionResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates the reaction's data and/or target feeds.
func (c *ReactionsClient) Update(ctx context.Context, id string, data map[string]any, targetFeeds []string) (*ReactionResponse, error) {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)

	reqData := map[string]any{
		"data":         data,
		"target_feeds": targetFeeds,
	}
	return c.decode(c.client.put(ctx, endpoint, reqData, c.client.authenticator.reactionsAuth))
}

// Get retrieves a reaction having the given id.
func (c *ReactionsClient) Get(ctx context.Context, id string) (*ReactionResponse, error) {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)

	return c.decode(c.client.get(ctx, endpoint, nil, c.client.authenticator.reactionsAuth))
}

// Delete deletes a reaction having the given id.
func (c *ReactionsClient) Delete(ctx context.Context, id string) (*ReactionResponse, error) {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)

	return c.decode(c.client.delete(ctx, endpoint, nil, c.client.authenticator.reactionsAuth))
}

// SoftDelete soft-deletes a reaction having the given id. It is possible to restore this reaction using ReactionsClient.Restore.
func (c *ReactionsClient) SoftDelete(ctx context.Context, id string) error {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)
	endpoint.addQueryParam(makeRequestOption("soft", true))

	_, err := c.client.delete(ctx, endpoint, nil, c.client.authenticator.reactionsAuth)
	return err
}

// Restore restores a soft deleted reaction having the given id.
func (c *ReactionsClient) Restore(ctx context.Context, id string) error {
	endpoint := c.client.makeEndpoint("reaction/%s/restore/", id)

	_, err := c.client.put(ctx, endpoint, nil, c.client.authenticator.reactionsAuth)
	return err
}

// Filter lists reactions based on the provided criteria and with the specified pagination.
func (c *ReactionsClient) Filter(ctx context.Context, attr FilterReactionsAttribute, opts ...FilterReactionsOption) (*FilterReactionResponse, error) {
	endpointURI := fmt.Sprintf("reaction/%s/", attr())

	endpoint := c.client.makeEndpoint(endpointURI)
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}

	resp, err := c.client.get(ctx, endpoint, nil, c.client.authenticator.reactionsAuth)
	if err != nil {
		return nil, err
	}
	var result FilterReactionResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	result.meta.attr = attr
	return &result, nil
}

// GetNextPageFilteredReactions returns the reactions at the "next" page of a previous *FilterReactionResponse response, if any.
func (c *ReactionsClient) GetNextPageFilteredReactions(ctx context.Context, resp *FilterReactionResponse) (*FilterReactionResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return c.Filter(ctx, resp.meta.attr, opts...)
}
