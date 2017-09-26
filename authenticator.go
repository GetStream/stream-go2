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
	resFollower    resource = "follower"
	resActivities  resource = "activities"
	resFeed        resource = "feed"
	resFeedTargets resource = "feed_targets"
)

type action string

const (
	actionRead   action = "read"
	actionWrite  action = "write"
	actionDelete action = "delete"
)

var actions = map[string]action{
	http.MethodGet:     actionRead,
	http.MethodOptions: actionRead,
	http.MethodHead:    actionRead,
	http.MethodPost:    actionWrite,
	http.MethodPut:     actionWrite,
	http.MethodPatch:   actionWrite,
	http.MethodDelete:  actionDelete,
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
