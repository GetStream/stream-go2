package stream

import (
	"encoding/json"
	"errors"
)

// UsersClient is a specialized client used to interact with the Reactions endpoints.
type ReactionsClient struct {
	client *Client
}

//Add adds a reaction.
func (c *ReactionsClient) Add(r AddReactionRequestObject) (*Reaction, error) {
	endpoint := c.client.makeEndpoint("reaction/")
	if r.ParentID != "" {
		return nil, errors.New("`Parent` not empty. For adding child reactions use `AddChild`")
	}

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

//AddChild adds a child reaction to the provided parent.
func (c *ReactionsClient) AddChild(parentID string, r AddReactionRequestObject) (*Reaction, error) {
	endpoint := c.client.makeEndpoint("reaction/")
	r.ParentID = parentID

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
	endpoint := c.client.makeEndpoint("user/%s/", id)

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

//Get retrieves a reaction having the given id.
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

//Delete deletes a reaction having the given id.
func (c *ReactionsClient) Delete(id string) error {
	endpoint := c.client.makeEndpoint("reaction/%s/", id)

	_, err := c.client.delete(endpoint, nil, c.client.authenticator.reactionsAuth)
	return err
}

func (c *ReactionsClient) Filter() {}
