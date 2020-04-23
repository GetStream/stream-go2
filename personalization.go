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

func (c *PersonalizationClient) decode(resp []byte, err error) (*PersonalizationResponse, error) {
	if err != nil {
		return nil, err
	}
	var result PersonalizationResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("cannot unmarshal resp: %w", err)
	}
	return &result, nil
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
	return c.decode(c.client.get(endpoint, nil, c.client.authenticator.personalizationAuth))
}

// Post sends data to the given resource, adding the given params to the request.
func (c *PersonalizationClient) Post(resource string, params, data map[string]interface{}) (*PersonalizationResponse, error) {
	if resource == "" {
		return nil, errors.New("missing resource")
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
	return c.decode(c.client.post(endpoint, data, c.client.authenticator.personalizationAuth))
}

// Delete removes data from the given resource, adding the given params to the request.
func (c *PersonalizationClient) Delete(resource string, params map[string]interface{}) (*PersonalizationResponse, error) {
	if resource == "" {
		return nil, errors.New("missing resource")
	}
	endpoint := c.client.makeEndpoint("%s/", resource)
	for k, v := range params {
		endpoint.addQueryParam(makeRequestOption(k, v))
	}
	return c.decode(c.client.delete(endpoint, nil, c.client.authenticator.personalizationAuth))
}
