package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// CollectionsClient is a specialized client used to interact with the Collection endpoints.
type CollectionsClient struct {
	client *Client
}

// Upsert creates new or updates existing objects for the given collection's name.
func (c *CollectionsClient) Upsert(collection string, objects ...CollectionObject) (*BaseResponse, error) {
	if collection == "" {
		return nil, errors.New("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/")
	data := map[string]interface{}{
		"data": map[string][]CollectionObject{
			collection: objects,
		},
	}
	return decode(c.client.post(endpoint, data, c.client.authenticator.collectionsAuth))
}

// Select returns a list of CollectionObjects for the given collection name
// having the given IDs.
func (c *CollectionsClient) Select(collection string, ids ...string) (*GetCollectionResponse, error) {
	if collection == "" {
		return nil, errors.New("collection name required")
	}
	foreignIDs := make([]string, len(ids))
	for i := range ids {
		foreignIDs[i] = fmt.Sprintf("%s:%s", collection, ids[i])
	}
	endpoint := c.client.makeEndpoint("collections/")
	endpoint.addQueryParam(makeRequestOption("foreign_ids", strings.Join(foreignIDs, ",")))
	resp, err := c.client.get(endpoint, nil, c.client.authenticator.collectionsAuth)
	if err != nil {
		return nil, err
	}
	var result GetCollectionResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMany removes from a collection the objects having the given IDs.
func (c *CollectionsClient) DeleteMany(collection string, ids ...string) (*BaseResponse, error) {
	if collection == "" {
		return nil, errors.New("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/")
	endpoint.addQueryParam(makeRequestOption("collection_name", collection))
	endpoint.addQueryParam(makeRequestOption("ids", strings.Join(ids, ",")))
	return decode(c.client.delete(endpoint, nil, c.client.authenticator.collectionsAuth))
}

func (c *CollectionsClient) decodeObject(resp []byte, err error) (*CollectionObjectResponse, error) {
	if err != nil {
		return nil, err
	}
	var result CollectionObjectResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Add adds a single object to a collection.
func (c *CollectionsClient) Add(collection string, object CollectionObject, opts ...AddObjectOption) (*CollectionObjectResponse, error) {
	if collection == "" {
		return nil, errors.New("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/", collection)

	req := addCollectionRequest{}

	for _, opt := range opts {
		opt(&req)
	}

	req.ID = object.ID
	req.Data = object.Data

	return c.decodeObject(c.client.post(endpoint, req, c.client.authenticator.collectionsAuth))
}

// Get retrieves a collection object having the given ID.
func (c *CollectionsClient) Get(collection, id string) (*CollectionObjectResponse, error) {
	if collection == "" {
		return nil, errors.New("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/%s/", collection, id)

	return c.decodeObject(c.client.get(endpoint, nil, c.client.authenticator.collectionsAuth))
}

// Update updates the given collection object's data.
func (c *CollectionsClient) Update(collection, id string, data map[string]interface{}) (*CollectionObjectResponse, error) {
	if collection == "" {
		return nil, errors.New("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/%s/", collection, id)
	reqData := map[string]interface{}{
		"data": data,
	}

	return c.decodeObject(c.client.put(endpoint, reqData, c.client.authenticator.collectionsAuth))
}

// Delete removes from a collection the object having the given ID.
func (c *CollectionsClient) Delete(collection, id string) (*BaseResponse, error) {
	if collection == "" {
		return nil, errors.New("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/%s/", collection, id)

	return decode(c.client.delete(endpoint, nil, c.client.authenticator.collectionsAuth))
}

// CreateReference returns a new reference string in the form SO:<collection>:<id>.
func (c *CollectionsClient) CreateReference(collection, id string) string {
	return fmt.Sprintf("SO:%s:%s", collection, id)
}

// CreateCollectionReference is a convenience helper not to require a client.
var CreateCollectionReference = (&CollectionsClient{}).CreateReference
