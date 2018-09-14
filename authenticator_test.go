package stream

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTAuth(t *testing.T) {
	a := authenticator{secret: "something very secret"}
	req, err := http.NewRequest(http.MethodPost, "", nil)
	require.NoError(t, err)

	err = a.jwtAuth(resFeed, "*")(req)
	assert.NoError(t, err)
	assert.Equal(t, "jwt", req.Header.Get("stream-auth-type"))
	expectedAuth := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3Rpb24iOiJ3cml0ZSIsImZlZWRfaWQiOiIqIiwicmVzb3VyY2UiOiJmZWVkIn0.LnWdqnKryMuXEX3p8HepCBRVGfvhdzINmA2jU1j3TUA"
	assert.Equal(t, expectedAuth, req.Header.Get("authorization"))
}

func TestFeedToken(t *testing.T) {
	a := authenticator{secret: "test"}
	token := a.feedSignature("flat:123")
	assert.Equal(t, "flat:123 u1jO3CcR8PeAD22yGNtVobcHrxI", token)
}
