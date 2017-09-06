package stream

import (
	"fmt"
	"net/http"

	httpsig "gopkg.in/LeisureLink/httpsig.v1"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

type authFunc func(*http.Request) error

type authenticator struct {
	secret string
}

type resource string

const (
	followerResource   resource = "follower"
	activitiesResource resource = "activities"
	feedResource       resource = "feed"
)

const (
	readAction   = "read"
	writeAction  = "write"
	deleteAction = "delete"
)

type action string

var actions = map[string]action{
	http.MethodGet:     readAction,
	http.MethodOptions: readAction,
	http.MethodHead:    readAction,
	http.MethodPost:    writeAction,
	http.MethodPut:     writeAction,
	http.MethodPatch:   writeAction,
	http.MethodDelete:  deleteAction,
}

func newAuthenticator(secret string) authenticator {
	return authenticator{
		secret: secret,
	}
}

func (a authenticator) feedAuth(resource resource, method string) authFunc {
	return func(req *http.Request) error {
		claims := jwt.MapClaims{
			"resource": resource,
			"action":   actions[method],
			"feed_id":  "*",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		auth, err := token.SignedString([]byte(a.secret))
		if err != nil {
			return fmt.Errorf("cannot make auth: %s", err)
		}
		req.Header.Add("stream-auth-type", "jwt")
		req.Header.Add("authorization", auth)
		return nil
	}
}

func (a authenticator) applicationAuth(key string) authFunc {
	return func(req *http.Request) error {
		signer, err := httpsig.NewRequestSigner(key, a.secret, "hmac-sha256")
		if err != nil {
			return fmt.Errorf("cannot sign request: %s", err)
		}
		return signer.SignRequest(req, []string{}, nil)
	}
}
