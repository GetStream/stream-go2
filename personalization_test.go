package stream_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPersonalizationGet(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	p := client.Personalization()
	params := map[string]any{"answer": 42, "feed": "user:123"}
	_, err := p.Get(ctx, "", params)
	require.Error(t, err)
	_, err = p.Get(ctx, "some_resource", params)
	require.NoError(t, err)
	expectedURL := "https://personalization.stream-io-api.com/personalization/v1.0/some_resource/?answer=42&api_key=key&feed=user%3A123"
	testRequest(t, requester.req, http.MethodGet, expectedURL, "")
}

func TestPersonalizationPost(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	p := client.Personalization()
	params := map[string]any{"answer": 42, "feed": "user:123"}
	_, err := p.Post(ctx, "", params, nil)
	require.Error(t, err)
	data := map[string]any{"foo": "bar", "baz": 42}
	_, err = p.Post(ctx, "some_resource", params, data)
	require.NoError(t, err)
	expectedURL := "https://personalization.stream-io-api.com/personalization/v1.0/some_resource/?answer=42&api_key=key&feed=user%3A123"
	expectedBody := `{"data":{"baz":42,"foo":"bar"}}`
	testRequest(t, requester.req, http.MethodPost, expectedURL, expectedBody)
}

func TestPersonalizationDelete(t *testing.T) {
	ctx := context.Background()
	client, requester := newClient(t)
	p := client.Personalization()
	params := map[string]any{"answer": 42, "feed": "user:123"}
	_, err := p.Delete(ctx, "", params)
	require.Error(t, err)
	_, err = p.Delete(ctx, "some_resource", params)
	require.NoError(t, err)
	expectedURL := "https://personalization.stream-io-api.com/personalization/v1.0/some_resource/?answer=42&api_key=key&feed=user%3A123"
	testRequest(t, requester.req, http.MethodDelete, expectedURL, "")
}
