package stream_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPersonalizationGet(t *testing.T) {
	client, requester := newClient(t)
	p := client.Personalization()
	params := map[string]interface{}{"answer": 42, "feed": "user:123"}
	_, err := p.Get("", params)
	require.Error(t, err)
	_, err = p.Get("some_resource", params)
	require.NoError(t, err)
	expectedURL := "https://personalization.stream-io-api.com/personalization/v1.0/some_resource/?answer=42&api_key=key&feed=user%3A123"
	testRequest(t, requester.req, http.MethodGet, expectedURL, "")
}

func TestPersonalizationPost(t *testing.T) {
	client, requester := newClient(t)
	p := client.Personalization()
	params := map[string]interface{}{"answer": 42, "feed": "user:123"}
	_, err := p.Post("", params, nil)
	require.Error(t, err)
	data := map[string]interface{}{"foo": "bar", "baz": 42}
	_, err = p.Post("some_resource", params, data)
	require.NoError(t, err)
	expectedURL := "https://personalization.stream-io-api.com/personalization/v1.0/some_resource/?answer=42&api_key=key&feed=user%3A123"
	expectedBody := `{"data":{"baz":42,"foo":"bar"}}`
	testRequest(t, requester.req, http.MethodPost, expectedURL, expectedBody)
}

func TestPersonalizationDelete(t *testing.T) {
	client, requester := newClient(t)
	p := client.Personalization()
	params := map[string]interface{}{"answer": 42, "feed": "user:123"}
	_, err := p.Delete("", params)
	require.Error(t, err)
	_, err = p.Delete("some_resource", params)
	require.NoError(t, err)
	expectedURL := "https://personalization.stream-io-api.com/personalization/v1.0/some_resource/?answer=42&api_key=key&feed=user%3A123"
	testRequest(t, requester.req, http.MethodDelete, expectedURL, "")
}
