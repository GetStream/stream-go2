package stream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	timeout       time.Duration
}

// Requester performs HTTP requests.
type Requester interface {
	Do(*http.Request) (*http.Response, error)
}

// New builds a new Client with the provided API key and secret. It can be
// configured further by passing any number of ClientOption parameters.
func New(key, secret string, opts ...ClientOption) (*Client, error) {
	if key == "" || secret == "" {
		return nil, errMissingCredentials
	}
	c := &Client{
		key:           key,
		timeout:       time.Second * 6,
		authenticator: authenticator{secret: secret},
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.requester == nil {
		c.requester = newRequester(c.timeout)
	}
	c.urlBuilder = newAPIURLBuilder(c.region, c.version)
	return c, nil
}

func newRequester(timeout time.Duration) Requester {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: timeout / 2,
			}).DialContext,
		},
	}
}

// NewFromEnv build a new Client using environment variables values, with
// possible values being STREAM_API_KEY, STREAM_API_SECRET, STREAM_API_REGION,
// and STREAM_API_VERSION.
// Additional options can still be provided as parameters.
func NewFromEnv(extraOptions ...ClientOption) (*Client, error) {
	key := os.Getenv("STREAM_API_KEY")
	secret := os.Getenv("STREAM_API_SECRET")
	region := os.Getenv("STREAM_API_REGION")
	version := os.Getenv("STREAM_API_VERSION")
	return New(key, secret, append(extraOptions, WithAPIRegion(region), WithAPIVersion(version))...)
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

// WithTimeout sets the HTTP request timeout
func WithTimeout(dur time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = dur
	}
}

// WithTimeout clones the client with the given timeout.
// If a custom requester was given while initializing, it will be overridden.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	nc := *c
	nc.timeout = timeout
	nc.requester = newRequester(timeout)
	return &nc
}

// FlatFeed returns a new Flat Feed with the provided slug and userID.
func (c *Client) FlatFeed(slug, userID string) (*FlatFeed, error) {
	feed, err := newFeed(slug, userID, c)
	if err != nil {
		return nil, err
	}
	return &FlatFeed{*feed}, nil
}

// AggregatedFeed returns a new Aggregated Feed with the provided slug and
// userID.
func (c *Client) AggregatedFeed(slug, userID string) (*AggregatedFeed, error) {
	feed, err := newFeed(slug, userID, c)
	if err != nil {
		return nil, err
	}
	return &AggregatedFeed{*feed}, nil
}

// NotificationFeed returns a new Notification Feed with the provided slug and
// userID.
func (c *Client) NotificationFeed(slug, userID string) (*NotificationFeed, error) {
	feed, err := newFeed(slug, userID, c)
	if err != nil {
		return nil, err
	}
	return &NotificationFeed{*feed}, nil
}

// GenericFeed returns a standard Feed implementation using the provided target id.
func (c *Client) GenericFeed(targetID string) (Feed, error) {
	parts := strings.Split(targetID, feedSlugIDSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid target id: %q", targetID)
	}

	return newFeed(parts[0], parts[1], c)
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
	_, err := c.post(endpoint, req, c.authenticator.feedAuth(resFeed, nil))
	return err
}

// FollowMany creates multiple follows at once.
func (c *Client) FollowMany(relationships []FollowRelationship, opts ...FollowManyOption) error {
	endpoint := c.makeEndpoint("follow_many/")
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}
	_, err := c.post(endpoint, relationships, c.authenticator.feedAuth(resFollower, nil))
	return err
}

// UnfollowMany removes multiple follow relationships at once.
func (c *Client) UnfollowMany(relationships []UnfollowRelationship) error {
	endpoint := c.makeEndpoint("unfollow_many/")
	_, err := c.post(endpoint, relationships, c.authenticator.feedAuth(resFollower, nil))
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

// Users returns a new UsersClient.
func (c *Client) Users() *UsersClient {
	b := newAPIURLBuilder(c.region, c.version)
	return &UsersClient{client: c.cloneWithURLBuilder(b)}
}

// Reactions returns a new ReactionsClient.
func (c *Client) Reactions() *ReactionsClient {
	b := newAPIURLBuilder(c.region, c.version)
	return &ReactionsClient{client: c.cloneWithURLBuilder(b)}
}

// Personalization returns a new PersonalizationClient.
func (c *Client) Personalization() *PersonalizationClient {
	b := newPersonalizationURLBuilder(c.region)
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

// GetEnrichedActivitiesByID returns enriched activities for the current app having the given IDs.
func (c *Client) GetEnrichedActivitiesByID(ids []string, opts ...GetActivitiesOption) (*GetEnrichedActivitiesResponse, error) {
	options := []GetActivitiesOption{{makeRequestOption("ids", strings.Join(ids, ","))}}
	return c.getAppEnrichedActivities(append(options, opts...)...)
}

// GetEnrichedActivitiesByForeignID returns enriched activities for the current app having the given foreign IDs and timestamps.
func (c *Client) GetEnrichedActivitiesByForeignID(values []ForeignIDTimePair, opts ...GetActivitiesOption) (*GetEnrichedActivitiesResponse, error) {
	foreignIDs := make([]string, len(values))
	timestamps := make([]string, len(values))
	for i, v := range values {
		foreignIDs[i] = v.ForeignID
		timestamps[i] = v.Timestamp.Format(TimeLayout)
	}
	options := []GetActivitiesOption{
		{makeRequestOption("foreign_ids", strings.Join(foreignIDs, ","))},
		{makeRequestOption("timestamps", strings.Join(timestamps, ","))},
	}

	return c.getAppEnrichedActivities(append(options, opts...)...)
}

func (c *Client) getAppEnrichedActivities(options ...GetActivitiesOption) (*GetEnrichedActivitiesResponse, error) {
	endpoint := c.makeEndpoint("enrich/activities/")
	for _, v := range options {
		endpoint.addQueryParam(v.requestOption)
	}
	data, err := c.get(endpoint, nil, c.authenticator.feedAuth(resActivities, nil))
	if err != nil {
		return nil, err
	}
	var resp GetEnrichedActivitiesResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateActivities updates existing activities.
func (c *Client) UpdateActivities(activities ...Activity) (*BaseResponse, error) {
	req := struct {
		Activities []Activity `json:"activities,omitempty"`
	}{
		Activities: activities,
	}
	endpoint := c.makeEndpoint("activities/")
	return decode(c.post(endpoint, req, c.authenticator.feedAuth(resActivities, nil)))
}

// PartialUpdateActivities performs a partial update on multiple activities with the given set and unset operations
// specified by each changeset. This returns the affected activities.
func (c *Client) PartialUpdateActivities(changesets ...UpdateActivityRequest) (*UpdateActivitiesResponse, error) {
	req := struct {
		Activities []UpdateActivityRequest `json:"changes,omitempty"`
	}{
		Activities: changesets,
	}
	endpoint := c.makeEndpoint("activity/")
	data, err := c.post(endpoint, req, c.authenticator.feedAuth(resActivities, nil))
	if err != nil {
		return nil, err
	}
	var resp UpdateActivitiesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, err
}

// UpdateActivityByID performs a partial activity update with the given set and unset operations, returning the
// affected activity, on the activity with the given ID.
func (c *Client) UpdateActivityByID(id string, set map[string]interface{}, unset []string) (*UpdateActivityResponse, error) {
	return c.updateActivity(UpdateActivityRequest{
		ID:    &id,
		Set:   set,
		Unset: unset,
	})
}

// UpdateActivityByForeignID performs a partial activity update with the given set and unset operations, returning the
// affected activity, on the activity with the given foreign ID and timestamp.
func (c *Client) UpdateActivityByForeignID(foreignID string, timestamp Time, set map[string]interface{}, unset []string) (*UpdateActivityResponse, error) {
	return c.updateActivity(UpdateActivityRequest{
		ForeignID: &foreignID,
		Time:      &timestamp,
		Set:       set,
		Unset:     unset,
	})
}

func (c *Client) updateActivity(req UpdateActivityRequest) (*UpdateActivityResponse, error) {
	endpoint := c.makeEndpoint("activity/")
	data, err := c.post(endpoint, req, c.authenticator.feedAuth(resActivities, nil))
	if err != nil {
		return nil, err
	}
	var resp UpdateActivityResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) makeStreamError(statusCode int, rate *Rate, body io.Reader) error {
	if body == nil {
		return errors.New("invalid body")
	}
	errBody, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	var streamErr APIError
	if err := json.Unmarshal(errBody, &streamErr); err != nil {
		return fmt.Errorf("unexpected error (status code %d)", statusCode)
	}
	streamErr.StatusCode = statusCode
	streamErr.Rate = rate
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
			return nil, fmt.Errorf("cannot marshal request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, endpoint.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	c.setBaseHeaders(req)

	if authFn != nil {
		if err := authFn(req); err != nil {
			return nil, err
		}
	}

	resp, err := c.requester.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot perform request: %w", err)
	}
	defer resp.Body.Close()

	rate := NewRate(resp.Header)

	if resp.StatusCode/100 != 2 {
		return nil, c.makeStreamError(resp.StatusCode, rate, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	out := map[string]interface{}{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	out["ratelimit"] = rate

	if body, err = json.Marshal(out); err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	return body, nil
}

func (c *Client) addActivity(feed Feed, activity Activity) (*AddActivityResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/", feed.Slug(), feed.UserID())
	resp, err := c.post(endpoint, activity, c.authenticator.feedAuth(resFeed, feed))
	if err != nil {
		return nil, err
	}
	var out AddActivityResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) addActivities(feed Feed, activities ...Activity) (*AddActivitiesResponse, error) {
	reqBody := struct {
		Activities []Activity `json:"activities,omitempty"`
	}{
		Activities: activities,
	}
	endpoint := c.makeEndpoint("feed/%s/%s/", feed.Slug(), feed.UserID())
	resp, err := c.post(endpoint, reqBody, c.authenticator.feedAuth(resFeed, feed))
	if err != nil {
		return nil, err
	}
	var out AddActivitiesResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	return &out, nil
}

func (c *Client) removeActivityByID(feed Feed, activityID string) (*RemoveActivityResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/%s/", feed.Slug(), feed.UserID(), activityID)
	resp, err := c.delete(endpoint, nil, c.authenticator.feedAuth(resFeed, feed))
	if err != nil {
		return nil, err
	}
	var out RemoveActivityResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) removeActivityByForeignID(feed Feed, foreignID string) (*RemoveActivityResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/%s/", feed.Slug(), feed.UserID(), foreignID)
	endpoint.addQueryParam(makeRequestOption("foreign_id", 1))
	resp, err := c.delete(endpoint, nil, c.authenticator.feedAuth(resFeed, feed))
	if err != nil {
		return nil, err
	}
	var out RemoveActivityResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) getActivities(feed Feed, opts ...GetActivitiesOption) ([]byte, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/", feed.Slug(), feed.UserID())
	return c.getActivitiesInternal(endpoint, feed, opts...)
}

func (c *Client) getEnrichedActivities(feed Feed, opts ...GetActivitiesOption) ([]byte, error) {
	endpoint := c.makeEndpoint("enrich/feed/%s/%s/", feed.Slug(), feed.UserID())
	return c.getActivitiesInternal(endpoint, feed, opts...)
}

func (c *Client) getActivitiesInternal(endpoint endpoint, feed Feed, opts ...GetActivitiesOption) ([]byte, error) {
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}
	return c.get(endpoint, nil, c.authenticator.feedAuth(resFeed, feed))
}

func (c *Client) follow(feed Feed, opts *followFeedOptions) (*BaseResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/follows/", feed.Slug(), feed.UserID())
	return decode(c.post(endpoint, opts, c.authenticator.feedAuth(resFollower, feed)))
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

func (c *Client) unfollow(feed Feed, target string, opts ...UnfollowOption) (*BaseResponse, error) {
	endpoint := c.makeEndpoint("feed/%s/%s/follows/%s/", feed.Slug(), feed.UserID(), target)
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}

	return decode(c.delete(endpoint, nil, c.authenticator.feedAuth(resFollower, feed)))
}

func (c *Client) followStats(feed Feed, opts ...FollowStatOption) (*FollowStatResponse, error) {
	endpoint := c.makeEndpoint("stats/follow/")
	endpoint.addQueryParam(makeRequestOption("followers", feed.ID()))
	endpoint.addQueryParam(makeRequestOption("following", feed.ID()))
	for _, opt := range opts {
		endpoint.addQueryParam(opt)
	}

	resp, err := c.get(endpoint, nil, c.authenticator.feedAuth(resFollower, nil))
	if err != nil {
		return nil, err
	}
	var out FollowStatResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) updateToTargets(feed Feed, activity Activity, opts ...UpdateToTargetsOption) (*UpdateToTargetsResponse, error) {
	endpoint := c.makeEndpoint("feed_targets/%s/%s/activity_to_targets/", feed.Slug(), feed.UserID())

	req := &updateToTargetsRequest{
		ForeignID: activity.ForeignID,
		Time:      activity.Time.Format(TimeLayout),
	}
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.post(endpoint, req, c.authenticator.feedAuth(resFeedTargets, feed))
	if err != nil {
		return nil, err
	}
	var out UpdateToTargetsResponse
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
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
