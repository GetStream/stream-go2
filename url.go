package stream

import (
	"fmt"
	"os"
)

const domain = "stream-io-api.com"

type urlBuilder interface {
	url() string
}

// handy rewrites for regions
var regionOverrides = map[string]string{
	"us-east":   "us-east-api",
	"eu-west":   "eu-west-api",
	"singapore": "singapore-api",
}

type regionalURLBuilder struct {
	region  string
	version string
}

func newRegionalURLBuilder(region, version string) regionalURLBuilder {
	return regionalURLBuilder{
		region:  region,
		version: version,
	}
}

func (u regionalURLBuilder) makeHost(subdomain string) string {
	if envHost := os.Getenv("STREAM_URL"); envHost != "" {
		return envHost
	}
	return fmt.Sprintf("https://%s.%s", u.makeRegion(subdomain), domain)
}

func (u regionalURLBuilder) makeVersion() string {
	if u.version != "" {
		return u.version
	}
	return "1.0"
}

func (u regionalURLBuilder) makeRegion(subdomain string) string {
	if u.region != "" {
		if override, ok := regionOverrides[u.region]; ok {
			return override
		}
		return u.region
	}
	return subdomain
}

type apiURLBuilder struct {
	regionalURLBuilder
}

func newAPIURLBuilder(region, version string) apiURLBuilder {
	return apiURLBuilder{newRegionalURLBuilder(region, version)}
}

func (u apiURLBuilder) url() string {
	return fmt.Sprintf("%s/api/v%s/", u.makeHost("api"), u.makeVersion())
}

type personalizationURLBuilder struct{}

func newPersonalizationURLBuilder() personalizationURLBuilder {
	return personalizationURLBuilder{}
}

func (b personalizationURLBuilder) url() string {
	if envHost := os.Getenv("STREAM_URL"); envHost != "" {
		return envHost
	}
	return "https://personalization.stream-io-api.com/personalization/v1.0/"
}

type analyticsURLBuilder struct {
	regionalURLBuilder
}

func newAnalyticsURLBuilder(region, version string) analyticsURLBuilder {
	return analyticsURLBuilder{newRegionalURLBuilder(region, version)}
}

func (u analyticsURLBuilder) url() string {
	return fmt.Sprintf("%s/analytics/v%s/", u.makeHost("analytics"), u.makeVersion())
}
