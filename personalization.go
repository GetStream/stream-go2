package stream

import (
	"encoding/json"
	"errors"
	"fmt"
)

// PersonalizationClient is a specialized client for personalization features.
type PersonalizationClient struct {
	client *Client
}

// Get obtains a PersonalizationResponse for the given resource and params.
func (c *PersonalizationClient) Get(resource string, params map[string]interface{}) (*PersonalizationResponse, error) {
	if resource == "" {
		return nil, errors.New("missing resource")
	}
	endpoint := c.client.makeEndpoint("%s/", resource)
	for k, v := range params {
		endpoint.addQueryParam(makeRequestOption(k, v))
	}
	resp, err := c.client.get(endpoint, nil, c.client.authenticator.personalizationAuth)
	if err != nil {
		return nil, err
	}
	var personalizationResp PersonalizationResponse
	err = json.Unmarshal(resp, &personalizationResp)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal resp: %w", err)
	}
	return &personalizationResp, nil
}

// Post sends data to the given resource, adding the given params to the request.
func (c *PersonalizationClient) Post(resource string, params, data map[string]interface{}) error {
	if resource == "" {
		return errors.New("missing resource")
	}
	endpoint := c.client.makeEndpoint("%s/", resource)
	for k, v := range params {
		endpoint.addQueryParam(makeRequestOption(k, v))
	}
	if data != nil {
		data = map[string]interface{}{
			"data": data,
		}
	}
	_, err := c.client.post(endpoint, data, c.client.authenticator.personalizationAuth)
	return err
}

// Delete removes data from the given resource, adding the given params to the request.
func (c *PersonalizationClient) Delete(resource string, params map[string]interface{}) error {
	if resource == "" {
		return errors.New("missing resource")
	}
	endpoint := c.client.makeEndpoint("%s/", resource)
	for k, v := range params {
		endpoint.addQueryParam(makeRequestOption(k, v))
	}
	_, err := c.client.delete(endpoint, nil, c.client.authenticator.personalizationAuth)
	return err
}
