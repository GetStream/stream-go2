package stream

import "fmt"

type apiURL struct {
	region  string
	version string
}

func (u *apiURL) String() string {
	region := "api"
	if u.region != "" {
		region = u.region + "-" + region
	}
	version := "1.0"
	if u.version != "" {
		version = u.version
	}
	return fmt.Sprintf("https://%s.getstream.io/api/v%s/", region, version)
}
