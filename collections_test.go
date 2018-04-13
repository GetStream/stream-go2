package stream_test

import (
	"net/http"
	"net/url"
	"testing"

	stream "github.com/GetStream/stream-go2"
	"github.com/stretchr/testify/require"
)

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
				stream.CollectionObject{
					ID:   "1",
					Name: "Juniper",
					Data: map[string]interface{}{
						"hobbies": []string{"playing", "sleeping", "eating"},
					},
				},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/meta/?api_key=key",
			expectedBody: `{"data":{"test-single":[{"hobbies":["playing","sleeping","eating"],"id":"1","name":"Juniper"}]}}`,
		},
		{
			collection: "test-many",
			objects: []stream.CollectionObject{
				stream.CollectionObject{
					ID:   "1",
					Name: "Juniper",
					Data: map[string]interface{}{
						"hobbies": []string{"playing", "sleeping", "eating"},
					},
				},
				stream.CollectionObject{
					ID:   "2",
					Name: "Ruby",
					Data: map[string]interface{}{
						"interests": []string{"sunbeams", "surprise attacks"},
					},
				},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/meta/?api_key=key",
			expectedBody: `{"data":{"test-many":[{"hobbies":["playing","sleeping","eating"],"id":"1","name":"Juniper"},{"id":"2","interests":["sunbeams","surprise attacks"],"name":"Ruby"}]}}`,
		},
	}
	for _, tc := range testCases {
		err := client.Collections().Upsert(tc.collection, tc.objects...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestGetCollectionObjects(t *testing.T) {
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
			expectedURL: "https://api.stream-io-api.com/api/v1.0/meta/?api_key=key&foreign_ids=" + url.QueryEscape("test-single:one"),
		},
		{
			collection:  "test-multiple",
			ids:         []string{"one", "two", "three"},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/meta/?api_key=key&foreign_ids=" + url.QueryEscape("test-multiple:one,test-multiple:two,test-multiple:three"),
		},
	}
	for _, tc := range testCases {
		_, err := client.Collections().Get(tc.collection, tc.ids...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expectedURL, tc.expectedBody)
	}
}

func TestDeleteCollectionObjects(t *testing.T) {
	client, requester := newClient(t)
	testCases := []struct {
		ids         []string
		collection  string
		expectedURL string
	}{
		{
			collection:  "test-single",
			ids:         []string{"one"},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/meta/?api_key=key&collection_name=test-single&ids=one",
		},
		{
			collection:  "test-many",
			ids:         []string{"one", "two", "three"},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/meta/?api_key=key&collection_name=test-many&ids=one%2Ctwo%2Cthree",
		},
	}
	for _, tc := range testCases {
		err := client.Collections().Delete(tc.collection, tc.ids...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodDelete, tc.expectedURL, "")
	}
}
