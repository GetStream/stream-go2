package stream_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v7"
)

func TestUserRefHelpers(t *testing.T) {
	client, _ := newClient(t)
	ref := client.Users().CreateReference("bar")
	assert.Equal(t, "SU:bar", ref)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	_, err := client.Users().Get(ctx, "id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/user/id1/?api_key=key", "")
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	_, err := client.Users().Delete(ctx, "id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.stream-io-api.com/api/v1.0/user/id1/?api_key=key", "")
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	testCases := []struct {
		object       stream.User
		getOrCreate  bool
		expectedURL  string
		expectedBody string
	}{
		{
			object: stream.User{
				ID: "user-test",
				Data: map[string]interface{}{
					"is_admin": true,
					"name":     "Johnny",
				},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/user/?api_key=key&get_or_create=false",
			expectedBody: `{"id":"user-test","data":{"is_admin":true,"name":"Johnny"}}`,
		},
		{
			object: stream.User{
				ID: "user-test",
				Data: map[string]interface{}{
					"is_admin": true,
					"name":     "Jane",
				},
			},
			getOrCreate:  true,
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/user/?api_key=key&get_or_create=true",
			expectedBody: `{"id":"user-test","data":{"is_admin":true,"name":"Jane"}}`,
		},
	}

	for _, tc := range testCases {
		_, err := client.Users().Add(ctx, tc.object, tc.getOrCreate)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	data := map[string]interface{}{
		"name": "Jane",
	}
	_, err := client.Users().Update(ctx, "123", data)
	require.NoError(t, err)
	expectedBody := `{"data":{"name":"Jane"}}`
	testRequest(t, requester.req, http.MethodPut, "https://api.stream-io-api.com/api/v1.0/user/123/?api_key=key", expectedBody)
}
