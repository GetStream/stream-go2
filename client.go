package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

// Client is a Stream API client used for retrieving feeds and performing API
// calls.
type Client struct {
	key           string
	requester     Requester
	authenticator authenticator
	urlBuilder    urlBuilder
	region        string
	version       string
}

// Requester performs HTTP requests.
type Requester interface {
	Do(*http.Request) (*http.Response, error)
}

// NewClient builds a new Client with the provided API key and secret. It can be
// configured further by passing any number of ClientOption parameters.
func NewClient(key, secret string, opts ...ClientOption) (*Client, error) {
	if key == "" || secret == "" {
		return nil, errMissingCredentials
	}
	c := &Client{
		key: key,
		requester: &http.Client{
			Transport: &http.Transport{},
		},
		authenticator: authenticator{secret: secret},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.urlBuilder = newAPIURLBuilder(c.region, c.version)
	return c, nil
}

// NewClientFromEnv build a new Client using environment variables values, with
// possible values being STREAM_API_KEY, STREAM_API_SECRET, STREAM_API_REGION,
// and STREAM_API_VERSION.
func NewClientFromEnv() (*Client, error) {
	key := os.Getenv("STREAM_API_KEY")
	secret := os.Getenv("STREAM_API_SECRET")
	region := os.Getenv("STREAM_API_REGION")
	version := os.Getenv("STREAM_API_VERSION")
	return NewClient(key, secret, WithAPIRegion(region), WithAPIVersion(version))
}

// ClientOption is a function used for adding specific configuration options to
// a Stream client.
type ClientOption func(*Client)

// WithAPIRegion sets the region for a given Client.
func WithAPIRegion(region string) ClientOption {
	return func(c *Client) {
		c.region = region
	}
}

// WithAPIVersion sets the version for a given Client.
func WithAPIVersion(version string) ClientOption {
	return func(c *Client) {
		c.version = version
	}
}

// WithHTTPRequester sets the HTTP requester for a given client, used mostly for testing.
func WithHTTPRequester(requester Requester) ClientOption {
	return func(c *Client) {
		c.requester = requester
	}
}

// FlatFeed returns a new Flat Feed with the provided slug and userID.
func (c *Client) FlatFeed(slug, userID string) *FlatFeed {
	return &FlatFeed{newFeed(slug, userID, c)}
}

// AggregatedFeed returns a new Aggregated Feed with the provided slug and
// userID.
func (c *Client) AggregatedFeed(slug, userID string) *AggregatedFeed {
	return &AggregatedFeed{newFeed(slug, userID, c)}
}

// NotificationFeed returns a new Notification Feed with the provided slug and
// userID.
func (c *Client) NotificationFeed(slug, userID string) *NotificationFeed {
	return &NotificationFeed{newFeed(slug, userID, c)}
}

// AddToMany adds an activity to multiple feeds at once.
func (c *Client) AddToMany(activity Activity, feeds ...Feed) error {
	endpoint := c.makeEndpoint("feed/add_to_many/")
	ids := make([]string, len(feeds))
	for i := range feeds {
		ids[i] = feeds[i].ID()
	}
	req := AddToManyRequest{
		Activity: activity,
		FeedIDs:  ids,
	}
	_, err := c.post(endpoint, req, c.authenticator.applicationAuth(c.key))
	return err
}

// FollowMany creates multiple follows at once.
func (c *Client) FollowMany(relationships []FollowRelationship, opts ...FollowManyOption) error {
	endpoint := c.makeEndpoint("follow_many/")
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}
	_, err := c.post(endpoint, relationships, c.authenticator.applicationAuth(c.key))
	return err
}

// UnfollowMany removes multiple follow relationships at once.
func (c *Client) UnfollowMany(relationships []UnfollowRelationship) error {
	endpoint := c.makeEndpoint("unfollow_many/")
	_, err := c.post(endpoint, relationships, c.authenticator.applicationAuth(c.key))
	return err
}

func (c *Client) cloneWithURLBuilder(builder urlBuilder) *Client {
	return &Client{
		key:           c.key,
		requester:     c.requester,
		authenticator: c.authenticator,
		urlBuilder:    builder,
	}
}

// Analytics returns a new AnalyticsClient sharing the base configuration of the original Client.
func (c *Client) Analytics() *AnalyticsClient {
	b := newAnalyticsURLBuilder(c.region, c.version)
	return &AnalyticsClient{client: c.cloneWithURLBuilder(b)}
}

// Collections returns a new CollectionsClient.
func (c *Client) Collections() *CollectionsClient {
	b := newAPIURLBuilder(c.region, c.version)
	return &CollectionsClient{client: c.cloneWithURLBuilder(b)}
}

// Personalization returns a new PersonalizationClient.
func (c *Client) Personalization() *PersonalizationClient {
	b := newPersonalizationURLBuilder()
	return &PersonalizationClient{client: c.cloneWithURLBuilder(b)}
}

// GetActivitiesByID returns activities for the current app having the given IDs.
func (c *Client) GetActivitiesByID(ids ...string) (*GetActivitiesResponse, error) {
	return c.getAppActivities(makeRequestOption("ids", strings.Join(ids, ",")))
}

// GetActivitiesByForeignID returns activities for the current app having the given foreign IDs and timestamps.
func (c *Client) GetActivitiesByForeignID(values ...ForeignIDTimePair) (*GetActivitiesResponse, error) {
	foreignIDs := make([]string, len(values))
	timestamps := make([]string, len(values))
	for i, v := range values {
		foreignIDs[i] = v.ForeignID
		timestamps[i] = v.Timestamp.Format(TimeLayout)
	}
	return c.getAppActivities(
		makeRequestOption("foreign_ids", strings.Join(foreignIDs, ",")),
		makeRequestOption("timestamps", strings.Join(timestamps, ",")),
	)
}

func (c *Client) getAppActivities(values ...valuer) (*GetActivitiesResponse, error) {
	endpoint := c.makeEndpoint("activities/")
	for _, v := range values {
		endpoint.addQueryParam(v)
	}
	data, err := c.get(endpoint, nil, c.authenticator.feedAuth(resActivities, nil))
	if err != nil {
		return nil, err
	}
	var resp GetActivitiesResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateActivities updates existing activities.
func (c *Client) UpdateActivities(activities ...Activity) error {
	signedActivities := c.signActivities(activities)
	req := struct {
		Activities []Activity `json:"activities,omitempty"`
	}{
		Activities: signedActivities,
	}
	endpoint := c.makeEndpoint("activities/")
	_, err := c.post(endpoint, req, c.authenticator.feedAuth(resActivities, nil))
	return err
}

// UpdateActivityByID performs a partial activity update with the given set and unset operations, returning the
// affected activity, on the activity with the given ID.
func (c *Client) UpdateActivityByID(id string, set map[string]interface{}, unset []string) (*UpdateActivityResponse, error) {
	return c.updateActivity(updateActivityRequest{
		ID:    &id,
		Set:   set,
		Unset: unset,
	})
}

// UpdateActivityByForeignID performs a partial activity update with the given set and unset operations, returning the
// affected activity, on the activity with the given foreign ID and timestamp.
func (c *Client) UpdateActivityByForeignID(foreignID string, timestamp Time, set map[string]interface{}, unset []string) (*UpdateActivityResponse, error) {
	return c.updateActivity(updateActivityRequest{
		ForeignID: &foreignID,
		Time:      &timestamp,
		Set:       set,
		Unset:     unset,
	})
}

func (c *Client) updateActivity(req updateActivityRequest) (*UpdateActivityResponse, error) {
	endpoint := c.makeEndpoint("activity/")
	data, err := c.post(endpoint, req, c.authenticator.feedAuth(resActivities, nil))
	if err != nil {
		return nil, err
	}
	var resp UpdateActivityResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	_, ok := resp.Extra["duration"].(string)
	if ok {
		delete(resp.Extra, "duration")
	}
	return &resp, nil
}

func (c *Client) makeStreamError(statusCode int, body io.Reader) error {
	if body == nil {
		return fmt.Errorf("invalid body")
	}
	errBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	var streamErr APIError
	if err := json.Unmarshal(errBody, &streamErr); err != nil {
		return fmt.Errorf("unexpected error (status code %d)", statusCode)
	}
	streamErr.StatusCode = statusCode
	return streamErr
}

type endpoint struct {
	url   *url.URL
	query url.Values
}

func (e endpoint) String() string {
	e.url.RawQuery = e.query.Encode()
	return e.url.String()
}

func (e *endpoint) addQueryParam(v valuer) {
	if !v.valid() {
		return
	}
	e.query.Add(v.values())
}

func (c *Client) makeEndpoint(format string, a ...interface{}) endpoint {
	host := c.urlBuilder.url()

	path := fmt.Sprintf(format, a...)
	u, _ := url.Parse(host + path)

	query := make(url.Values)
	query.Set("api_key", c.key)

	return endpoint{
		url:   u,
		query: query,
	}
}

func (c *Client) get(endpoint endpoint, data interface{}, authFn authFunc) ([]byte, error) {
	return c.request(http.MethodGet, endpoint, data, authFn)
}

func (c *Client) post(endpoint endpoint, data interface{}, authFn authFunc) ([]byte, error) {
	return c.request(http.MethodPost, endpoint, data, authFn)
}

func (c *Client) put(endpoint endpoint, data interface{}, authFn authFunc) ([]byte, error) {
	return c.request(http.MethodPut, endpoint, data, authFn)
}

func (c *Client) delete(endpoint endpoint, data interface{}, authFn authFunc) ([]byte, error) {
	return c.request(http.MethodDelete, endpoint, data, authFn)
}

func (c *Client) setBaseHeaders(r *http.Request) {
	r.Header.Set("Content-type", "application/json")
	r.Header.Set("X-Stream-Client", fmt.Sprintf("stream-go2-client-%s", Version))
}

func (c *Client) request(method string, endpoint endpoint, data interface{}, authFn authFunc) ([]byte, error) {
	var reader io.Reader
	if data != nil {
		payload, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal request: %s", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, endpoint.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %s", err)
	}
	c.setBaseHeaders(req)

	if authFn != nil {
		if err := authFn(req); err != nil {
			return nil, err
		}
	}

	resp, err := c.requester.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot perform request: %s", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, c.makeStreamError(resp.StatusCode, resp.Body)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %s", err)
	}

	return body, nil
}

func (c *Client) signActivity(activity Activity) Activity {
	if len(activity.To) == 0 {
		return activity
	}
	tos := make([]string, len(activity.To))
	signed := activity
	for i, id := range activity.To {
		tos[i] = c.authenticator.feedSignature(id)
	}
	signed.To = tos
	return signed
}

func (c *Client) signActivities(activities []Activity) []Activity {
	signed := make([]Activity, len(activities))
	for i := range activities {
		signed[i] = c.signActivity(activities[i])
	}
	return signed
}

func (c *Client) addActivity(feed Feed, activity Activity) (*AddActivityResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/", feed.Slug(), feed.UserID())
	signedActivity := c.signActivity(activity)
	resp, err := c.post(endpoint, signedActivity, c.authenticator.feedAuth(resFeed, feed))
	if err != nil {
		return nil, err
	}
	var out AddActivityResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	_, ok := out.Extra["duration"].(string)
	if ok {
		delete(out.Extra, "duration")
	}
	return &out, nil
}

func (c *Client) addActivities(feed Feed, activities ...Activity) (*AddActivitiesResponse, error) {
	signedActivities := c.signActivities(activities)
	reqBody := struct {
		Activities []Activity `json:"activities,omitempty"`
	}{
		Activities: signedActivities,
	}
	endpoint := c.makeEndpoint("feed/%s/%s/", feed.Slug(), feed.UserID())
	resp, err := c.post(endpoint, reqBody, c.authenticator.feedAuth(resFeed, feed))
	if err != nil {
		return nil, err
	}
	var out AddActivitiesResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %s", err)
	}
	return &out, nil
}

func (c *Client) removeActivityByID(feed Feed, activityID string) error {
	endpoint := c.makeEndpoint("feed/%s/%s/%s/", feed.Slug(), feed.UserID(), activityID)
	_, err := c.delete(endpoint, nil, c.authenticator.feedAuth(resFeed, feed))
	return err
}

func (c *Client) removeActivityByForeignID(feed Feed, foreignID string) error {
	endpoint := c.makeEndpoint("feed/%s/%s/%s/", feed.Slug(), feed.UserID(), foreignID)
	endpoint.addQueryParam(makeRequestOption("foreign_id", 1))
	_, err := c.delete(endpoint, nil, c.authenticator.feedAuth(resFeed, feed))
	return err
}

func (c *Client) getActivities(feed Feed, opts ...GetActivitiesOption) ([]byte, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/", feed.Slug(), feed.UserID())
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}
	return c.get(endpoint, nil, c.authenticator.feedAuth(resFeed, feed))
}

func (c *Client) follow(feed Feed, opts *followFeedOptions) error {
	endpoint := c.makeEndpoint("feed/%s/%s/follows/", feed.Slug(), feed.UserID())
	_, err := c.post(endpoint, opts, c.authenticator.feedAuth(resFollower, feed))
	return err
}

func (c *Client) getFollowers(feed Feed, opts ...FollowersOption) (*FollowersResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/followers/", feed.Slug(), feed.UserID())
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}

	resp, err := c.get(endpoint, nil, c.authenticator.feedAuth(resFollower, feed))
	if err != nil {
		return nil, err
	}
	var out FollowersResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) getFollowing(feed Feed, opts ...FollowingOption) (*FollowingResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/follows/", feed.Slug(), feed.UserID())
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}

	resp, err := c.get(endpoint, nil, c.authenticator.feedAuth(resFollower, feed))
	if err != nil {
		return nil, err
	}
	var out FollowingResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) unfollow(feed Feed, target string, opts ...UnfollowOption) error {
	endpoint := c.makeEndpoint("feed/%s/%s/follows/%s/", feed.Slug(), feed.UserID(), target)
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}

	_, err := c.delete(endpoint, nil, c.authenticator.feedAuth(resFollower, feed))
	return err
}

func (c *Client) updateToTargets(feed Feed, activity Activity, opts ...UpdateToTargetsOption) error {
	endpoint := c.makeEndpoint("feed_targets/%s/%s/activity_to_targets/", feed.Slug(), feed.UserID())

	req := &updateToTargetsRequest{
		ForeignID: activity.ForeignID,
		Time:      activity.Time.Format(TimeLayout),
	}
	for _, opt := range opts {
		opt(req)
	}

	_, err := c.post(endpoint, req, c.authenticator.feedAuth(resFeedTargets, feed))
	return err
}

func (c *Client) GetUserSessionToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
	}
	return c.authenticator.jwtSignatureFromClaims(claims)
}

func (c *Client) GetUserSessionTokenWithClaims(userID string, claims map[string]interface{}) (string, error) {
	claims["user_id"] = userID
	jwtclaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtclaims[k] = v
	}
	return c.authenticator.jwtSignatureFromClaims(jwtclaims)
}
