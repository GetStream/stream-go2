package stream

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedAuth(t *testing.T) {
	a := authenticator{secret: "something very secret"}
	req, err := http.NewRequest(http.MethodPost, "", nil)
	require.NoError(t, err)

	err = a.feedAuth(resFeed, nil)(req)
	assert.NoError(t, err)
	assert.Equal(t, "jwt", req.Header.Get("stream-auth-type"))
	expectedAuth := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb24iOiJ3cml0ZSIsImZlZWRfaWQiOiIqIiwicmVzb3VyY2UiOiJmZWVkIn0.LnWdqnKryMuXEX3p8HepCBRVGfvhdzINmA2jU1j3TUA"
	assert.Equal(t, expectedAuth, req.Header.Get("authorization"))
}

func TestApplicationAuth(t *testing.T) {
	a := authenticator{secret: "something very secret"}
	req, err := http.NewRequest(http.MethodPost, "", nil)
	require.NoError(t, err)

	err = a.applicationAuth("my key")(req)
	assert.NoError(t, err)
	assert.Equal(t, "my key", req.Header.Get("X-API-Key"))
	expectedAuthRe := regexp.MustCompile(`Signature keyId="my key",algorithm="hmac-sha256",headers="date",signature="[0-9a-zA-Z/+]{43}="`)
	assert.Regexp(t, expectedAuthRe, req.Header.Get("authorization"))
}
