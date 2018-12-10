package stream

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CollectionsClient is a specialized client used to interact with the Collection endpoints.
type CollectionsClient struct {
	client *Client
}

// Upsert creates new or updates existing objects for the given collection's name.
func (c *CollectionsClient) Upsert(collection string, objects ...CollectionObject) error {
	if collection == "" {
		return fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/")
	data := map[string]interface{}{
		"data": map[string][]CollectionObject{
			collection: objects,
		},
	}
	_, err := c.client.post(endpoint, data, c.client.authenticator.collectionsAuth)
	return err
}

// Select returns a list of CollectionObjects for the given collection name
// having the given IDs.
func (c *CollectionsClient) Select(collection string, ids ...string) ([]GetCollectionResponseObject, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection name required")
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
	var selectResp getCollectionResponseWrap
	err = json.Unmarshal(resp, &selectResp)
	if err != nil {
		return nil, err
	}
	return selectResp.Response.Data, nil
}

// DeleteMany removes from a collection the objects having the given IDs.
func (c *CollectionsClient) DeleteMany(collection string, ids ...string) error {
	if collection == "" {
		return fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/")
	endpoint.addQueryParam(makeRequestOption("collection_name", collection))
	endpoint.addQueryParam(makeRequestOption("ids", strings.Join(ids, ",")))
	_, err := c.client.delete(endpoint, nil, c.client.authenticator.collectionsAuth)
	return err
}

//Add adds a single object to a collection.
func (c *CollectionsClient) Add(collection string, object CollectionObject, opts ...AddObjectOption) (*CollectionObject, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/", collection)

	req := addCollectionRequest{}

	for _, opt := range opts {
		opt(&req)
	}

	req.ID = object.ID
	req.Data = object.Data

	resp, err := c.client.post(endpoint, req, c.client.authenticator.collectionsAuth)
	if err != nil {
		return nil, err
	}
	result := &CollectionObject{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//Get retrives a collection obejct having the given ID.
func (c *CollectionsClient) Get(collection string, id string) (*CollectionObject, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/%s/", collection, id)

	resp, err := c.client.get(endpoint, nil, c.client.authenticator.collectionsAuth)
	if err != nil {
		return nil, err
	}
	result := &CollectionObject{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//Update updates the given user's data.
func (c *CollectionsClient) Update(collection string, id string, data map[string]interface{}) (*CollectionObject, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/%s/", collection, id)
	reqData := map[string]interface{}{
		"data": data,
	}

	resp, err := c.client.put(endpoint, reqData, c.client.authenticator.collectionsAuth)
	if err != nil {
		return nil, err
	}
	result := &CollectionObject{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete removes from a collection the object having the given ID.
func (c *CollectionsClient) Delete(collection string, id string) error {
	if collection == "" {
		return fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("collections/%s/%s/", collection, id)

	_, err := c.client.delete(endpoint, nil, c.client.authenticator.collectionsAuth)
	return err
}

// CreateReference returns a new reference string in the form SO:<collection>:<id>.
func (c *CollectionsClient) CreateReference(collection, id string) string {
	return fmt.Sprintf("SO:%s:%s", collection, id)
}
