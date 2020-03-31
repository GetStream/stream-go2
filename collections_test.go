package stream_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v4"
)

func TestCollectionRefHelpers(t *testing.T) {
	client, _ := newClient(t)
	ref := client.Collections().CreateReference("foo", "bar")
	assert.Equal(t, "SO:foo:bar", ref)
}

func TestUpsertCollectionObjects(t *testing.T) {
	client, requester := newClient(t)
	testCases := []struct {
		objects      []stream.CollectionObject
		collection   string
		expectedURL  string
		expectedBody string
	}{
		{
			collection: "test-single",
			objects: []stream.CollectionObject{
				{
					ID: "1",
					Data: map[string]interface{}{
						"name":    "Juniper",
						"hobbies": []string{"playing", "sleeping", "eating"},
					},
				},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/collections/?api_key=key",
			expectedBody: `{"data":{"test-single":[{"hobbies":["playing","sleeping","eating"],"id":"1","name":"Juniper"}]}}`,
		},
		{
			collection: "test-many",
			objects: []stream.CollectionObject{
				{
					ID: "1",
					Data: map[string]interface{}{
						"name":    "Juniper",
						"hobbies": []string{"playing", "sleeping", "eating"},
					},
				},
				{
					ID: "2",
					Data: map[string]interface{}{
						"name":      "Ruby",
						"interests": []string{"sunbeams", "surprise attacks"},
					},
				},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/collections/?api_key=key",
			expectedBody: `{"data":{"test-many":[{"hobbies":["playing","sleeping","eating"],"id":"1","name":"Juniper"},{"id":"2","interests":["sunbeams","surprise attacks"],"name":"Ruby"}]}}`,
		},
	}
	for _, tc := range testCases {
		_, err := client.Collections().Upsert(tc.collection, tc.objects...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestSelectCollectionObjects(t *testing.T) {
	client, requester := newClient(t)
	testCases := []struct {
		ids          []string
		collection   string
		expectedURL  string
		expectedBody string
	}{
		{
			collection:  "test-single",
			ids:         []string{"one"},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/collections/?api_key=key&foreign_ids=" + url.QueryEscape("test-single:one"),
		},
		{
			collection:  "test-multiple",
			ids:         []string{"one", "two", "three"},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/collections/?api_key=key&foreign_ids=" + url.QueryEscape("test-multiple:one,test-multiple:two,test-multiple:three"),
		},
	}
	for _, tc := range testCases {
		_, err := client.Collections().Select(tc.collection, tc.ids...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expectedURL, tc.expectedBody)
	}
}

func TestDeleteManyCollectionObjects(t *testing.T) {
	client, requester := newClient(t)
	testCases := []struct {
		ids         []string
		collection  string
		expectedURL string
	}{
		{
			collection:  "test-single",
			ids:         []string{"one"},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/collections/?api_key=key&collection_name=test-single&ids=one",
		},
		{
			collection:  "test-many",
			ids:         []string{"one", "two", "three"},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/collections/?api_key=key&collection_name=test-many&ids=one%2Ctwo%2Cthree",
		},
	}
	for _, tc := range testCases {
		_, err := client.Collections().DeleteMany(tc.collection, tc.ids...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodDelete, tc.expectedURL, "")
	}
}

func TestGetCollectionObject(t *testing.T) {
	client, requester := newClient(t)

	_, err := client.Collections().Get("test-get-one", "id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/collections/test-get-one/id1/?api_key=key", "")
}

func TestDeleteCollectionObject(t *testing.T) {
	client, requester := newClient(t)

	_, err := client.Collections().Delete("test-get-one", "id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.stream-io-api.com/api/v1.0/collections/test-get-one/id1/?api_key=key", "")
}

func TestAddCollectionObject(t *testing.T) {
	client, requester := newClient(t)
	testCases := []struct {
		object       stream.CollectionObject
		collection   string
		opts         []stream.AddObjectOption
		expectedURL  string
		expectedBody string
	}{
		{
			collection: "no_user_id",
			object: stream.CollectionObject{
				ID: "1",
				Data: map[string]interface{}{
					"name":    "Juniper",
					"hobbies": []string{"playing", "sleeping", "eating"},
				},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/collections/no_user_id/?api_key=key",
			expectedBody: `{"data":{"hobbies":["playing","sleeping","eating"],"name":"Juniper"},"id":"1"}`,
		},
		{
			collection: "with_user_id",
			object: stream.CollectionObject{
				ID: "1",
				Data: map[string]interface{}{
					"name":    "Juniper",
					"hobbies": []string{"playing", "sleeping", "eating"},
				},
			},
			opts:         []stream.AddObjectOption{stream.WithUserID("user1")},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/collections/with_user_id/?api_key=key",
			expectedBody: `{"data":{"hobbies":["playing","sleeping","eating"],"name":"Juniper"},"id":"1","user_id":"user1"}`,
		},
	}
	for _, tc := range testCases {
		_, err := client.Collections().Add(tc.collection, tc.object, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestUpdateCollectionObject(t *testing.T) {
	client, requester := newClient(t)

	data := map[string]interface{}{
		"name": "Jane",
	}
	_, err := client.Collections().Update("test-collection", "123", data)
	require.NoError(t, err)
	expectedBody := `{"data":{"name":"Jane"}}`
	testRequest(t, requester.req, http.MethodPut, "https://api.stream-io-api.com/api/v1.0/collections/test-collection/123/?api_key=key", expectedBody)
}
