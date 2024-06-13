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

var personalizationOverrides = map[string]string{
	"eu-west": "dublin",
	"dublin":  "dublin",
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
	addr string
	regionalURLBuilder
}

func newAPIURLBuilder(addr, region, version string) apiURLBuilder {
	return apiURLBuilder{addr, newRegionalURLBuilder(region, version)}
}

func (u apiURLBuilder) url() string {
	if u.addr != "" {
		return fmt.Sprintf("%s/api/v%s/", u.addr, u.makeVersion())
	}
	return fmt.Sprintf("%s/api/v%s/", u.makeHost("api"), u.makeVersion())
}

type personalizationURLBuilder struct {
	region string
}

func newPersonalizationURLBuilder(region string) personalizationURLBuilder {
	return personalizationURLBuilder{
		region: region,
	}
}

func (b personalizationURLBuilder) url() string {
	if envHost := os.Getenv("STREAM_URL"); envHost != "" {
		return envHost
	}
	defaultPath := fmt.Sprintf("personalization.%s/personalization/v1.0/", domain)
	if override, ok := personalizationOverrides[b.region]; ok {
		return fmt.Sprintf("https://%s-%s", override, defaultPath)
	}

	return fmt.Sprintf("https://%s", defaultPath)
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
