package stream

import (
	"github.com/fatih/structs"
)

const (
	// TimeLayout is the default time parse layout for Stream API JSON time fields
	TimeLayout = "2006-01-02T15:04:05.999999"
)

var timeLayouts = []string{
	TimeLayout,
	"2006-01-02 15:04:05.999999-07:00",
}

func init() {
	structs.DefaultTagName = "json"
}
