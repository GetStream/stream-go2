package stream

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ReactionsClient is a specialized client used to interact with the Reactions endpoints.
type ReactionsClient struct {
	client *Client
}

// Add adds a reaction.
func (c *ReactionsClient) Add(r AddReactionRequestObject) (*Reaction, error) {
	if r.ParentID != "" {
		return nil, errors.New("`Parent` not empty. For adding child reactions use `AddChild`")
	}
	return c.addReaction(r)
}

// AddChild adds a child reaction to the provided parent.
func (c *ReactionsClient) AddChild(parentID string, r AddReactionRequestObject) (*Reaction, error) {
	r.ParentID = parentID
	return c.addReaction(r)
}

func (c *ReactionsClient) addReaction(r AddReactionRequestObject) (*Reaction, error) {
	endpoint := c.client.makeEndpoint("reaction/")
	resp, err := c.client.post(endpoint, r, c.client.authenticator.reactionsAuth)
	if err != nil {
		return nil, err
	}

	result := &Reaction{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update updates the reaction's data and/or target feeds.
func (c *ReactionsClient) Update(id string, data map[string]interface{}, targetFeeds []string) (*Reaction, error) {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)

	reqData := map[string]interface{}{
		"data":         data,
		"target_feeds": targetFeeds,
	}
	resp, err := c.client.put(endpoint, reqData, c.client.authenticator.reactionsAuth)
	if err != nil {
		return nil, err
	}

	result := &Reaction{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get retrieves a reaction having the given id.
func (c *ReactionsClient) Get(id string) (*Reaction, error) {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)

	resp, err := c.client.get(endpoint, nil, c.client.authenticator.reactionsAuth)
	if err != nil {
		return nil, err
	}

	result := &Reaction{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete deletes a reaction having the given id.
func (c *ReactionsClient) Delete(id string) error {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)

	_, err := c.client.delete(endpoint, nil, c.client.authenticator.reactionsAuth)
	return err
}

// Filter lists reactions based on the provided criteria and with the specified pagination.
func (c *ReactionsClient) Filter(attr FilterReactionsAttribute, opts ...FilterReactionsOption) (*FilterReactionResponse, error) {
	endpointURI := fmt.Sprintf("reaction/%s/", attr())

	endpoint := c.client.makeEndpoint(endpointURI)
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}

	resp, err := c.client.get(endpoint, nil, c.client.authenticator.reactionsAuth)
	if err != nil {
		return nil, err
	}
	result := &FilterReactionResponse{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	result.meta.attr = attr
	return result, nil
}

// GetNextPageFilteredReactions returns the reactions at the "next" page of a previous *FilterReactionResponse response, if any.
func (c *ReactionsClient) GetNextPageFilteredReactions(resp *FilterReactionResponse) (*FilterReactionResponse, error) {
	opts, err := resp.parseNext()
	if err != nil {
		return nil, err
	}
	return c.Filter(resp.meta.attr, opts...)
}
