package stream

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	httpsig "gopkg.in/LeisureLink/httpsig.v1"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

type authFunc func(*http.Request) error

type resource string

const (
	resFollower          resource = "follower"
	resActivities        resource = "activities"
	resFeed              resource = "feed"
	resFeedTargets       resource = "feed_targets"
	resCollections       resource = "collections"
	resPersonalization   resource = "personalization"
	resAnalytics         resource = "analytics"
	resAnalyticsRedirect resource = "redirect_and_track"
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

func (a authenticator) feedSignature(feedID string) string {
	return fmt.Sprintf("%s %s", feedID, a.feedToken(feedID))
}

func (a authenticator) feedToken(feedID string) string {
	id := strings.Replace(feedID, ":", "", -1)
	hash := sha1.New()
	hash.Write([]byte(a.secret))
	mac := hmac.New(sha1.New, hash.Sum(nil))
	mac.Write([]byte(id))
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return a.urlSafe(digest)
}

func (a authenticator) feedID(feed Feed) string {
	if feed == nil {
		return "*"
	}
	return fmt.Sprintf("%s%s", feed.Slug(), feed.UserID())
}

func (a authenticator) feedAuth(resource resource, feed Feed) authFunc {
	return func(req *http.Request) error {
		var feedID string
		if feed != nil {
			feedID = a.feedID(feed)
		} else {
			feedID = "*"
		}
		return a.jwtSignRequest(req, a.jwtFeedClaims(resource, actions[req.Method], feedID))
	}
}

func (a authenticator) collectionsAuth(req *http.Request) error {
	claims := jwt.MapClaims{
		"action":   "*",
		"user_id":  "*",
		"feed_id":  "*",
		"resource": resCollections,
	}
	return a.jwtSignRequest(req, claims)
}

func (a authenticator) personalizationAuth(req *http.Request) error {
	claims := jwt.MapClaims{
		"action":   "*",
		"user_id":  "*",
		"feed_id":  "*",
		"resource": resPersonalization,
	}
	return a.jwtSignRequest(req, claims)
}

func (a authenticator) analyticsAuth(req *http.Request) error {
	claims := jwt.MapClaims{
		"action":   "*",
		"user_id":  "*",
		"resource": resAnalytics,
	}
	return a.jwtSignRequest(req, claims)
}

func (a authenticator) signAnalyticsRedirectEndpoint(endpoint *endpoint) error {
	claims := jwt.MapClaims{
		"action":   "*",
		"user_id":  "*",
		"resource": resAnalyticsRedirect,
	}
	signature, err := a.jwtSignatureFromClaims(claims)
	if err != nil {
		return err
	}
	endpoint.addQueryParam(makeRequestOption("stream-auth-type", "jwt"))
	endpoint.addQueryParam(makeRequestOption("authorization", signature))
	return nil
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

func (a authenticator) jwtSignatureFromClaims(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}

func (a authenticator) jwtFeedClaims(resource resource, action action, feedID string) jwt.MapClaims {
	return jwt.MapClaims{
		"resource": resource,
		"action":   action,
		"feed_id":  feedID,
	}
}

func (a authenticator) jwtSignRequest(req *http.Request, claims jwt.MapClaims) error {
	auth, err := a.jwtSignatureFromClaims(claims)
	if err != nil {
		return fmt.Errorf("cannot make auth: %s", err)
	}
	req.Header.Add("Stream-Auth-Type", "jwt")
	req.Header.Add("Authorization", auth)
	return nil
}

func (a authenticator) urlSafe(src string) string {
	src = strings.Replace(src, "+", "-", -1)
	src = strings.Replace(src, "/", "_", -1)
	src = strings.Trim(src, "=")
	return src
}
