package stream

import (
	"fmt"
	"net/http"

	httpsig "gopkg.in/LeisureLink/httpsig.v1"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

type authFunc func(*http.Request) error

type resource string

const (
	followerResource   resource = "follower"
	activitiesResource resource = "activities"
	feedResource       resource = "feed"
)

type action string

const (
	readAction   action = "read"
	writeAction  action = "write"
	deleteAction action = "delete"
)

var actions = map[string]action{
	http.MethodGet:     readAction,
	http.MethodOptions: readAction,
	http.MethodHead:    readAction,
	http.MethodPost:    writeAction,
	http.MethodPut:     writeAction,
	http.MethodPatch:   writeAction,
	http.MethodDelete:  deleteAction,
}

type authenticator struct {
	secret string
}

func (a authenticator) feedAuth(resource resource) authFunc {
	return func(req *http.Request) error {
		claims := jwt.MapClaims{
			"resource": resource,
			"action":   actions[req.Method],
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
		req.Header.Set("x-api-key", key)
		signer, err := httpsig.NewRequestSigner(key, a.secret, "hmac-sha256")
		if err != nil {
			return fmt.Errorf("cannot sign request: %s", err)
		}
		return signer.SignRequest(req, []string{}, nil)
	}
}
