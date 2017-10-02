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

func (a authenticator) feedAuthToken(resource resource, action action, feedID string) (string, error) {
	claims := jwt.MapClaims{
		"resource": resource,
		"action":   action,
		"feed_id":  feedID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}

func (a authenticator) feedSignature(feedID string) string {
	token, err := a.feedAuthToken(resFeed, "*", feedID)
	if err != nil {
		return feedID
	}
	return fmt.Sprintf("%s %s", feedID, token)
}

func (a authenticator) feedID(feed Feed) string {
	if feed == nil {
		return "*"
	}
	return fmt.Sprintf("%s%s", feed.Slug(), feed.UserID())
}

func (a authenticator) feedAuth(resource resource, feed Feed) authFunc {
	return func(req *http.Request) error {
		auth, err := a.feedAuthToken(resource, actions[req.Method], a.feedID(feed))
		if err != nil {
			return fmt.Errorf("cannot make auth: %s", err)
		}
		req.Header.Add("Stream-Auth-Type", "jwt")
		req.Header.Add("Authorization", auth)
		return nil
	}
}

func (a authenticator) applicationAuth(key string) authFunc {
	return func(req *http.Request) error {
		req.Header.Set("X-API-Key", key)
		signer, err := httpsig.NewRequestSigner(key, a.secret, "hmac-sha256")
		if err != nil {
			return fmt.Errorf("cannot sign request: %s", err)
		}
		return signer.SignRequest(req, []string{}, nil)
	}
}
