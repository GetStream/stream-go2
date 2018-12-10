package stream

import (
	"encoding/json"
	"fmt"
)

// UsersClient is a specialized client used to interact with the Users endpoints.
type UsersClient struct {
	client *Client
}

// Add adds a new user with the specified id and optional extra data.
func (c *UsersClient) Add(user UserObject, getOrCreate bool) (*UserObject, error) {
	endpoint := c.client.makeEndpoint("user/")
	endpoint.addQueryParam(makeRequestOption("get_or_create", getOrCreate))

	resp, err := c.client.post(endpoint, user, c.client.authenticator.usersAuth)
	if err != nil {
		return nil, err
	}

	result := &UserObject{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update updates the user's data.
func (c *UsersClient) Update(id string, data map[string]interface{}) (*UserObject, error) {
	endpoint := c.client.makeEndpoint("user/%s/", id)

	reqData := map[string]interface{}{
		"data": data,
	}
	resp, err := c.client.put(endpoint, reqData, c.client.authenticator.usersAuth)
	if err != nil {
		return nil, err
	}

	result := &UserObject{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//Get retrieves a user having the given id.
func (c *UsersClient) Get(id string) (*UserObject, error) {
	endpoint := c.client.makeEndpoint("user/%s/", id)

	resp, err := c.client.get(endpoint, nil, c.client.authenticator.usersAuth)
	if err != nil {
		return nil, err
	}

	result := &UserObject{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//Delete deletes a user having the given id.
func (c *UsersClient) Delete(id string) error {
	endpoint := c.client.makeEndpoint("user/%s/", id)

	_, err := c.client.delete(endpoint, nil, c.client.authenticator.usersAuth)
	return err
}

// CreateReference returns a new reference string in the form SU:<id>.
func (c *UsersClient) CreateReference(id string) string {
	return fmt.Sprintf("SU:%s", id)
}
