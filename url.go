package stream

import "fmt"

type apiURL struct {
	region  string
	version string
}

func (u *apiURL) String() string {
	if u.region == "localhost" {
		return "http://localhost:8000/api/v1.0/"
	}
	return fmt.Sprintf("https://%s.getstream.io/api/v%s/", u.makeRegion(), u.makeVersion())
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
