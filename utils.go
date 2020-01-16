package stream

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func decodeJSONHook(f, typ reflect.Type, data interface{}) (interface{}, error) {
	switch typ {
	case reflect.TypeOf(Time{}):
		return timeFromString(data.(string))
	case reflect.TypeOf(Duration{}):
		switch v := data.(type) {
		case string:
			return durationFromString(v)
		case float64:
			return durationFromString(fmt.Sprintf("%fs", v))
		default:
			return nil, fmt.Errorf("invalid duration")
		}
	case reflect.TypeOf(Data{}):
		switch v := data.(type) {
		case string:
			return Data{
				ID: v,
			}, nil
		case map[string]interface{}:
			a := Data{}
			if err := a.decode(v); err != nil {
				return nil, err
			}
			return a, nil
		default:
			return nil, fmt.Errorf("invalid data")
		}
	}
	return data, nil
}

func decodeData(data map[string]interface{}, target interface{}) (*mapstructure.Metadata, error) {
	cfg := &mapstructure.DecoderConfig{
		DecodeHook: decodeJSONHook,
		Result:     target,
		Metadata:   &mapstructure.Metadata{},
		TagName:    "json",
	}
	dec, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return nil, err
	}
	if err := dec.Decode(data); err != nil {
		return nil, err
	}
	return cfg.Metadata, nil
}

func parseIntValue(values url.Values, key string) (val int, exits bool, err error) {
	v := values.Get(key)
	if v == "" {
		return 0, false, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, false, err
	}
	return i, true, nil
}

func parseBool(value string) bool {
	v := strings.ToLower(value)
	return v != "" && v != "false" && v != "f" && v != "0"
}
