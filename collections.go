package stream

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CollectionsClient struct {
	client *Client
}

// Upsert creates new or updates existing objects for the given collection's name.
func (c *CollectionsClient) Upsert(collection string, objects ...CollectionObject) error {
	if collection == "" {
		return fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("meta/")
	data := map[string]interface{}{
		"data": map[string][]CollectionObject{
			collection: objects,
		},
	}
	_, err := c.client.post(endpoint, data, c.client.authenticator.collectionsAuth)
	return err
}

// Get returns a list of CollectionObjects for the given collection name
// having the given IDs.
func (c *CollectionsClient) Get(collection string, ids ...string) ([]CollectionObject, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection name required")
	}
	foreignIDs := make([]string, len(ids))
	for i := range ids {
		foreignIDs[i] = fmt.Sprintf("%s:%s", collection, ids[i])
	}
	endpoint := c.client.makeEndpoint("meta/")
	endpoint.addQueryParam(makeRequestOption("foreign_ids", strings.Join(foreignIDs, ",")))
	resp, err := c.client.get(endpoint, nil, c.client.authenticator.collectionsAuth)
	if err != nil {
		return nil, err
	}
	var selectResp selectCollectionResponse
	err = json.Unmarshal(resp, &selectResp)
	if err != nil {
		return nil, err
	}
	objects := make([]CollectionObject, len(selectResp.Response.Data))
	for i, obj := range selectResp.Response.Data {
		objects[i] = obj.toCollectionObject()
	}
	return objects, nil
}

// Delete removes from a collection the objects having the given IDs.
func (c *CollectionsClient) Delete(collection string, ids ...string) error {
	if collection == "" {
		return fmt.Errorf("collection name required")
	}
	endpoint := c.client.makeEndpoint("meta/")
	endpoint.addQueryParam(makeRequestOption("collection_name", collection))
	endpoint.addQueryParam(makeRequestOption("ids", strings.Join(ids, ",")))
	_, err := c.client.delete(endpoint, nil, c.client.authenticator.collectionsAuth)
	return err
}
