package stream

import (
	"encoding/json"
	"fmt"
)

// UsersClient is a specialized client used to interact with the Users endpoints.
type UsersClient struct {
	client *Client
}

func (c *UsersClient) decode(resp []byte, err error) (*UserResponse, error) {
	if err != nil {
		return nil, err
	}
	var result UserResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Add adds a new user with the specified id and optional extra data.
func (c *UsersClient) Add(user User, getOrCreate bool) (*UserResponse, error) {
	endpoint := c.client.makeEndpoint("user/")
	endpoint.addQueryParam(makeRequestOption("get_or_create", getOrCreate))

	return c.decode(c.client.post(endpoint, user, c.client.authenticator.usersAuth))
}

// Update updates the user's data.
func (c *UsersClient) Update(id string, data map[string]interface{}) (*UserResponse, error) {
	endpoint := c.client.makeEndpoint("user/%s/", id)

	reqData := map[string]interface{}{
		"data": data,
	}
	return c.decode(c.client.put(endpoint, reqData, c.client.authenticator.usersAuth))
}

// Get retrieves a user having the given id.
func (c *UsersClient) Get(id string) (*UserResponse, error) {
	endpoint := c.client.makeEndpoint("user/%s/", id)

	return c.decode(c.client.get(endpoint, nil, c.client.authenticator.usersAuth))
}

// Delete deletes a user having the given id.
func (c *UsersClient) Delete(id string) (*BaseResponse, error) {
	endpoint := c.client.makeEndpoint("user/%s/", id)

	return decode(c.client.delete(endpoint, nil, c.client.authenticator.usersAuth))
}

// CreateReference returns a new reference string in the form SU:<id>.
func (c *UsersClient) CreateReference(id string) string {
	return fmt.Sprintf("SU:%s", id)
}

// CreateUserReference is a convenience helper not to require a client.
var CreateUserReference = (&UsersClient{}).CreateReference
