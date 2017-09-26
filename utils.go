package stream

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func decodeJSONStringTimes(f reflect.Type, typ reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	switch typ {
	case reflect.TypeOf(Time{}):
		return timeFromString(data.(string))
	case reflect.TypeOf(Duration{}):
		return durationFromString(data.(string))
	}
	return data, nil
}

func decodeData(data map[string]interface{}, target interface{}) (*mapstructure.Metadata, error) {
	cfg := &mapstructure.DecoderConfig{
		DecodeHook: decodeJSONStringTimes,
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
