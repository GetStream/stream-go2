package stream_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v4"
)

func TestGetReaction(t *testing.T) {
	client, requester := newClient(t)
	_, err := client.Reactions().Get("id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/id1/?api_key=key", "")
}

func TestDeleteReaction(t *testing.T) {
	client, requester := newClient(t)
	err := client.Reactions().Delete("id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.stream-io-api.com/api/v1.0/reaction/id1/?api_key=key", "")
}

func TestAddReaction(t *testing.T) {
	client, requester := newClient(t)

	testCases := []struct {
		input        stream.AddReactionRequestObject
		expectedURL  string
		expectedBody string
	}{
		{
			input: stream.AddReactionRequestObject{
				Kind:       "like",
				ActivityID: "some-act-id",
				UserID:     "user-id",
				Data: map[string]interface{}{
					"field": "value",
				},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/reaction/?api_key=key",
			expectedBody: `{"kind":"like","activity_id":"some-act-id","user_id":"user-id","data":{"field":"value"}}`,
		},
		{
			input: stream.AddReactionRequestObject{
				Kind:        "like",
				ActivityID:  "some-act-id",
				UserID:      "user-id",
				TargetFeeds: []string{"user:bob"},
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/reaction/?api_key=key",
			expectedBody: `{"kind":"like","activity_id":"some-act-id","user_id":"user-id","target_feeds":["user:bob"]}`,
		},
	}
	for _, tc := range testCases {
		_, err := client.Reactions().Add(tc.input)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestAddChildReaction(t *testing.T) {
	client, requester := newClient(t)

	reaction := stream.AddReactionRequestObject{
		Kind:       "like",
		ActivityID: "some-act-id",
		UserID:     "user-id",
		Data: map[string]interface{}{
			"field": "value",
		},
	}
	expectedURL := "https://api.stream-io-api.com/api/v1.0/reaction/?api_key=key"
	expectedBody := `{"kind":"like","activity_id":"some-act-id","user_id":"user-id","data":{"field":"value"},"parent":"pid"}`

	_, err := client.Reactions().AddChild("pid", reaction)
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodPost, expectedURL, expectedBody)
}

func TestUpdateReaction(t *testing.T) {
	client, requester := newClient(t)

	testCases := []struct {
		id           string
		data         map[string]interface{}
		targetFeeds  []string
		expectedURL  string
		expectedBody string
	}{
		{
			id: "r-id",
			data: map[string]interface{}{
				"field": "value",
			},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/reaction/r-id/?api_key=key",
			expectedBody: `{"data":{"field":"value"}}`,
		},
		{
			id:           "r-id2",
			targetFeeds:  []string{"user:bob"},
			expectedURL:  "https://api.stream-io-api.com/api/v1.0/reaction/r-id2/?api_key=key",
			expectedBody: `{"target_feeds":["user:bob"]}`,
		},
	}
	for _, tc := range testCases {
		_, err := client.Reactions().Update(tc.id, tc.data, tc.targetFeeds)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPut, tc.expectedURL, tc.expectedBody)
	}
}

func TestFilterReactions(t *testing.T) {
	client, requester := newClient(t)
	testCases := []struct {
		attr        stream.FilterReactionsAttribute
		opts        []stream.FilterReactionsOption
		expectedURL string
	}{
		{
			attr:        stream.ByActivityID("aid"),
			expectedURL: "https://api.stream-io-api.com/api/v1.0/reaction/activity_id/aid/?api_key=key",
		},
		{
			attr:        stream.ByReactionID("rid"),
			expectedURL: "https://api.stream-io-api.com/api/v1.0/reaction/reaction_id/rid/?api_key=key",
		},
		{
			attr:        stream.ByUserID("uid"),
			expectedURL: "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/?api_key=key",
		},
		{
			attr:        stream.ByUserID("uid").ByKind("upvote"),
			expectedURL: "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/upvote/?api_key=key",
		},
		{
			attr:        stream.ByUserID("uid").ByKind("upvote"),
			opts:        []stream.FilterReactionsOption{stream.WithLimit(100)},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&limit=100",
		},
		{
			attr:        stream.ByUserID("uid").ByKind("upvote"),
			opts:        []stream.FilterReactionsOption{stream.WithLimit(100), stream.WithActivityData(), stream.WithIDGTE("uid1")},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gte=uid1&limit=100&with_activity_data=true",
		},
	}

	for _, tc := range testCases {
		_, err := client.Reactions().Filter(tc.attr, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expectedURL, "")
	}
}

func TestGetNextPageReactions(t *testing.T) {
	client, requester := newClient(t)

	requester.resp = `{"next":"/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gt=uid1&limit=100&with_activity_data=true"}`
	resp, err := client.Reactions().Filter(stream.ByUserID("uid").ByKind("like"), stream.WithLimit(10), stream.WithActivityData(), stream.WithIDGT("id1"))
	require.NoError(t, err)

	_, err = client.Reactions().GetNextPageFilteredReactions(resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/like/?api_key=key&id_gt=uid1&limit=100&with_activity_data=true", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gt=uid1&limit=100&with_own_children=true"}`
	resp, err = client.Reactions().Filter(stream.ByUserID("uid").ByKind("like"), stream.WithLimit(10), stream.WithActivityData(), stream.WithIDGT("id1"))
	require.NoError(t, err)

	_, err = client.Reactions().GetNextPageFilteredReactions(resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/like/?api_key=key&id_gt=uid1&limit=100&with_own_children=true", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gt=uid1&limit=100&with_activity_data=false"}`
	resp, err = client.Reactions().Filter(stream.ByUserID("uid").ByKind("like"), stream.WithLimit(10), stream.WithActivityData(), stream.WithIDGT("id1"))
	require.NoError(t, err)

	_, err = client.Reactions().GetNextPageFilteredReactions(resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/like/?api_key=key&id_gt=uid1&limit=100", "")
	require.NoError(t, err)

	requester.resp = `{"next":"123"}`
	resp, err = client.Reactions().Filter(stream.ByActivityID("aid"))
	require.NoError(t, err)
	_, err = client.Reactions().GetNextPageFilteredReactions(resp)
	require.Error(t, err)

	requester.resp = `{"next":"?q=a%"}`
	resp, err = client.Reactions().Filter(stream.ByActivityID("aid"))
	require.NoError(t, err)
	_, err = client.Reactions().GetNextPageFilteredReactions(resp)
	require.Error(t, err)
}
