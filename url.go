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

func (u *apiURL) makeRegion() string {
	if u.region != "" {
		return fmt.Sprintf("%s-api", u.region)
	}
	return "api"
}
