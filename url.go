package stream

import "fmt"

type apiURL struct {
	region  string
	version string
}

const domain = "stream-io-api.com"

func (u *apiURL) String() string {
	return fmt.Sprintf("https://%s.%s/api/v%s/", u.makeRegion(), domain, u.makeVersion())
}

func (u *apiURL) makeVersion() string {
	if u.version != "" {
		return u.version
	}
	return "1.0"
}

// handy rewrites for regions
var regionOverrides = map[string]string{
	"us-east":   "us-east-api",
	"eu-west":   "eu-west-api",
	"singapore": "singapore-api",
}

func (u *apiURL) makeRegion() string {
	if u.region != "" {
		if override, ok := regionOverrides[u.region]; ok {
			return override
		}
		return u.region
	}
	return "api"
}
