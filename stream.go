package stream

import (
	"os"
)

const (
	timeLayout  = "2006-01-02T15:04:05.999999"
	defaultHost = "https://getstream.io/api/v1.0"
)

var host = defaultHost

func init() {
	envHost := os.Getenv("STREAM_HOST")
	if envHost != "" {
		host = envHost
	}
}
