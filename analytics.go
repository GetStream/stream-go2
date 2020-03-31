package stream

import (
	"encoding/json"
)

// AnalyticsClient is a specialized client used to send and track
// analytics events for enabled apps.
type AnalyticsClient struct {
	client *Client
}

// TrackEngagement is used to send and track analytics EngagementEvents.
func (c *AnalyticsClient) TrackEngagement(events ...EngagementEvent) (*BaseResponse, error) {
	endpoint := c.client.makeEndpoint("engagement/")
	data := map[string]interface{}{
		"content_list": events,
	}
	return decode(c.client.post(endpoint, data, c.client.authenticator.analyticsAuth))
}

// TrackImpression is used to send and track analytics ImpressionEvents.
func (c *AnalyticsClient) TrackImpression(eventsData ImpressionEventsData) (*BaseResponse, error) {
	endpoint := c.client.makeEndpoint("impression/")
	return decode(c.client.post(endpoint, eventsData, c.client.authenticator.analyticsAuth))
}

// RedirectAndTrack is used to send and track analytics ImpressionEvents. It tracks
// the events data (either EngagementEvents or ImpressionEvents) and redirects to the provided
// URL string.
func (c *AnalyticsClient) RedirectAndTrack(url string, events ...map[string]interface{}) (string, error) {
	endpoint := c.client.makeEndpoint("redirect/")
	eventsData, err := json.Marshal(events)
	if err != nil {
		return "", err
	}
	endpoint.addQueryParam(makeRequestOption("events", string(eventsData)))
	endpoint.addQueryParam(makeRequestOption("url", url))
	err = c.client.authenticator.signAnalyticsRedirectEndpoint(&endpoint)
	return endpoint.String(), err
}
