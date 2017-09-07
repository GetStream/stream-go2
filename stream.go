package stream

import (
	"os"

	"github.com/fatih/structs"
)

const (
	// TimeLayout is the default time parse layout for Stream API JSON time fields
	TimeLayout = "2006-01-02T15:04:05.999999"
)

var (
	host = "https://getstream.io/api/v1.0"
)

func init() {
	structs.DefaultTagName = "json"
	envHost := os.Getenv("STREAM_HOST")
	if envHost != "" {
		host = envHost
	}
}
