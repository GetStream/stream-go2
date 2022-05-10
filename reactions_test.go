package stream_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	stream "github.com/GetStream/stream-go2/v6"
)

func TestGetReaction(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	_, err := client.Reactions().Get(ctx, "id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/id1/?api_key=key", "")
}

func TestDeleteReaction(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	_, err := client.Reactions().Delete(ctx, "id1")
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodDelete, "https://api.stream-io-api.com/api/v1.0/reaction/id1/?api_key=key", "")
}

func TestAddReaction(t *testing.T) {
	ctx := context.Background()
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
		{
			input: stream.AddReactionRequestObject{
				Kind:                 "like",
				ActivityID:           "some-act-id",
				UserID:               "user-id",
				Data:                 map[string]interface{}{"some_extra": "on reaction"},
				TargetFeeds:          []string{"user:bob"},
				TargetFeedsExtraData: map[string]interface{}{"some_extra": "on activity"},
			},
			expectedURL: "https://api.stream-io-api.com/api/v1.0/reaction/?api_key=key",
			expectedBody: `{
				"kind":"like","activity_id":"some-act-id","user_id":"user-id",
				"data":{"some_extra":"on reaction"},
				"target_feeds":["user:bob"],
				"target_feeds_extra_data":{"some_extra":"on activity"}
			}`,
		},
	}
	for _, tc := range testCases {
		_, err := client.Reactions().Add(ctx, tc.input)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPost, tc.expectedURL, tc.expectedBody)
	}
}

func TestAddChildReaction(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	reaction := stream.AddReactionRequestObject{
		Kind:       "like",
		ActivityID: "some-act-id",
		UserID:     "user-id",
		Data: map[string]interface{}{
			"field": "value",
		},
		TargetFeeds: []string{"stalker:timeline"},
		TargetFeedsExtraData: map[string]interface{}{
			"activity_field": "activity_value",
		},
	}
	expectedURL := "https://api.stream-io-api.com/api/v1.0/reaction/?api_key=key"
	expectedBody := `{
		"kind":"like","activity_id":"some-act-id","user_id":"user-id",
		"data":{"field":"value"},"parent":"pid",
		"target_feeds": ["stalker:timeline"],"target_feeds_extra_data":{"activity_field":"activity_value"}
	}`

	_, err := client.Reactions().AddChild(ctx, "pid", reaction)
	require.NoError(t, err)
	testRequest(t, requester.req, http.MethodPost, expectedURL, expectedBody)
}

func TestUpdateReaction(t *testing.T) {
	ctx := context.Background()
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
		_, err := client.Reactions().Update(ctx, tc.id, tc.data, tc.targetFeeds)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodPut, tc.expectedURL, tc.expectedBody)
	}
}

func TestFilterReactions(t *testing.T) {
	ctx := context.Background()
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
		_, err := client.Reactions().Filter(ctx, tc.attr, tc.opts...)
		require.NoError(t, err)
		testRequest(t, requester.req, http.MethodGet, tc.expectedURL, "")
	}
}

func TestGetNextPageReactions(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)

	requester.resp = `{"next":"/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gt=uid1&limit=100&with_activity_data=true"}`
	resp, err := client.Reactions().Filter(ctx, stream.ByUserID("uid").ByKind("like"), stream.WithLimit(10), stream.WithActivityData(), stream.WithIDGT("id1"))
	require.NoError(t, err)

	_, err = client.Reactions().GetNextPageFilteredReactions(ctx, resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/like/?api_key=key&id_gt=uid1&limit=100&with_activity_data=true", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gt=uid1&limit=100&with_own_children=true"}`
	resp, err = client.Reactions().Filter(
		ctx,
		stream.ByUserID("uid").ByKind("like"),
		stream.WithLimit(10),
		stream.WithOwnChildren(),
		stream.WithIDGT("id1"),
	)
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gt=uid1&limit=100&with_own_children=true&with_own_children_kinds=comment,like&user_id=something&children_user_id=child_user_id"}`
	_, err = client.Reactions().Filter(
		ctx,
		stream.ByUserID("uid").ByKind("like"),
		stream.WithLimit(10),
		stream.WithIDGT("id1"),
		stream.WithOwnChildren(),
		stream.FilterReactionsOption(stream.WithEnrichOwnChildrenKindsFilter("comment", "like")),
	)
	require.NoError(t, err)

	_, err = client.Reactions().GetNextPageFilteredReactions(ctx, resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/like/?api_key=key&id_gt=uid1&limit=100&with_own_children=true", "")
	require.NoError(t, err)

	requester.resp = `{"next":"/api/v1.0/reaction/user_id/uid/upvote/?api_key=key&id_gt=uid1&limit=100&with_activity_data=false"}`
	resp, err = client.Reactions().Filter(ctx, stream.ByUserID("uid").ByKind("like"), stream.WithLimit(10), stream.WithActivityData(), stream.WithIDGT("id1"))
	require.NoError(t, err)

	_, err = client.Reactions().GetNextPageFilteredReactions(ctx, resp)
	testRequest(t, requester.req, http.MethodGet, "https://api.stream-io-api.com/api/v1.0/reaction/user_id/uid/like/?api_key=key&id_gt=uid1&limit=100", "")
	require.NoError(t, err)

	requester.resp = `{"next":"123"}`
	resp, err = client.Reactions().Filter(ctx, stream.ByActivityID("aid"))
	require.NoError(t, err)
	_, err = client.Reactions().GetNextPageFilteredReactions(ctx, resp)
	require.Error(t, err)

	requester.resp = `{"next":"?q=a%"}`
	resp, err = client.Reactions().Filter(ctx, stream.ByActivityID("aid"))
	require.NoError(t, err)
	_, err = client.Reactions().GetNextPageFilteredReactions(ctx, resp)
	require.Error(t, err)
}
