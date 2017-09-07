package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Client is a Stream API client used for retrieving feeds and performing API
// calls.
type Client struct {
	key           string
	appID         string
	cl            *http.Client
	authenticator authenticator
}

// NewClient builds a new Client with the provided API key and secret. It can be
// configured further by passing any number of ClientOption parameters.
func NewClient(key, secret string, opts ...ClientOption) (*Client, error) {
	if key == "" || secret == "" {
		return nil, errMissingCredentials
	}
	c := &Client{
		key:           key,
		cl:            &http.Client{},
		authenticator: authenticator{secret: secret},
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// ClientOption is a function used for adding specific configuration options to
// a Stream client.
type ClientOption func(*Client) error

// FlatFeed returns a new Flat Feed with the provided slug and userID.
func (c *Client) FlatFeed(slug, userID string) *FlatFeed {
	return &FlatFeed{newFeed(slug, userID, c)}
}

// AggregatedFeed returns a new Aggregated Feed with the provided slug and
// userID.
func (c *Client) AggregatedFeed(slug, userID string) *AggregatedFeed {
	return &AggregatedFeed{newFeed(slug, userID, c)}
}

// AddToMany adds an activity to multiple feeds at once.
func (c *Client) AddToMany(activity Activity, feeds ...Feed) error {
	endpoint := c.makeEndpoint("/feed/add_to_many/")
	ids := make([]string, len(feeds))
	for i := range feeds {
		ids[i] = feeds[i].ID()
	}
	req := AddToManyRequest{
		Activity: activity,
		Feeds:    ids,
	}
	_, err := c.request(http.MethodPost, endpoint, req, c.authenticator.applicationAuth(c.key))
	return err
}

// FollowMany creates multiple follows at once.
func (c *Client) FollowMany(relationships []FollowRelationship, opts ...RequestOption) error { // TODO test activity_copy_limit
	endpoint := c.makeEndpoint("/follow_many/")
	for _, opt := range opts {
		endpoint += opt.String()
	}
	_, err := c.request(http.MethodPost, endpoint, relationships, c.authenticator.applicationAuth(c.key))
	return err
}

func (c *Client) makeStreamError(body io.ReadCloser) error {
	errBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	var streamErr APIError
	if err := json.Unmarshal(errBody, &streamErr); err != nil {
		return err
	}
	return streamErr
}

func (c *Client) makeEndpoint(f string, a ...interface{}) string {
	format := fmt.Sprintf("%s%s?api_key=%s", host, f, c.key)
	k := fmt.Sprintf(format, a...)
	return k
}

func (c *Client) request(method, endpoint string, data interface{}, authFn authFunc) ([]byte, error) {
	var reader io.Reader
	if data != nil {
		payload, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal request: %s", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, endpoint, reader)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %s", err)
	}
	req.Header.Set("Content-type", "application/json")

	if err := authFn(req); err != nil {
		return nil, err
	}

	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot perform request: %s", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, c.makeStreamError(resp.Body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %s", err)
	}
	return body, nil
}

func (c *Client) addActivities(slug, userID string, activities ...Activity) (*AddActivitiesResponse, error) {
	reqBody := struct {
		Activities []Activity `json:"activities,omitempty"`
	}{
		Activities: activities,
	}
	endpoint := c.makeEndpoint("/feed/%s/%s/", slug, userID)
	resp, err := c.request(http.MethodPost, endpoint, reqBody, c.authenticator.feedAuth(feedResource))
	if err != nil {
		return nil, err
	}
	var out AddActivitiesResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %s", err)
	}
	return &out, nil
}

func (c *Client) updateActivities(activities ...Activity) error {
	req := struct {
		Activities []Activity `json:"activities,omitempty"`
	}{
		Activities: activities,
	}
	endpoint := c.makeEndpoint("/activities/")
	_, err := c.request(http.MethodPost, endpoint, req, c.authenticator.feedAuth(activitiesResource))
	return err
}

func (c *Client) removeActivityByID(slug, userID, activityID string) error {
	endpoint := c.makeEndpoint("/feed/%s/%s/%s/", slug, userID, activityID)
	_, err := c.request(http.MethodDelete, endpoint, nil, c.authenticator.feedAuth(feedResource))
	return err
}

func (c *Client) removeActivityByForeignID(slug, userID, foreignID string) error {
	endpoint := c.makeEndpoint("/feed/%s/%s/%s/", slug, userID, foreignID)
	endpoint += "&foreign_id=1"
	_, err := c.request(http.MethodDelete, endpoint, nil, c.authenticator.feedAuth(feedResource))
	return err
}

func (c *Client) getActivities(slug, userID string, opts ...RequestOption) ([]byte, error) {
	endpoint := c.makeEndpoint("/feed/%s/%s/", slug, userID)
	for _, opt := range opts {
		endpoint += opt.String()
	}
	return c.request(http.MethodGet, endpoint, nil, c.authenticator.feedAuth(feedResource))
}

func (c *Client) follow(slug, userID string, opts *followFeedOptions) error {
	endpoint := c.makeEndpoint("/feed/%s/%s/follows/", slug, userID)
	_, err := c.request(http.MethodPost, endpoint, opts, c.authenticator.feedAuth(followerResource))
	return err
}

func (c *Client) getFollowers(slug, userID string, opts ...RequestOption) (*FollowersResponse, error) {
	endpoint := c.makeEndpoint("/feed/%s/%s/followers/", slug, userID)
	for _, opt := range opts {
		endpoint += opt.String()
	}
	resp, err := c.request(http.MethodGet, endpoint, nil, c.authenticator.feedAuth(followerResource))
	if err != nil {
		return nil, err
	}
	var out FollowersResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) getFollowing(slug, userID string, opts ...RequestOption) (*FollowingResponse, error) {
	endpoint := c.makeEndpoint("/feed/%s/%s/follows/", slug, userID)
	for _, opt := range opts {
		endpoint += opt.String()
	}
	resp, err := c.request(http.MethodGet, endpoint, nil, c.authenticator.feedAuth(followerResource))
	if err != nil {
		return nil, err
	}
	var out FollowingResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) unfollow(slug, userID, target string, opts ...RequestOption) error {
	endpoint := c.makeEndpoint("/feed/%s/%s/follows/%s/", slug, userID, target)
	for _, opt := range opts {
		endpoint += opt.String()
	}
	_, err := c.request(http.MethodDelete, endpoint, nil, c.authenticator.feedAuth(followerResource))
	return err
}
