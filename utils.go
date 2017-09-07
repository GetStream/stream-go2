package stream

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

func decodeJSONStringTimes(f reflect.Type, typ reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	switch typ {
	case reflect.TypeOf(time.Time{}):
		return time.Parse(TimeLayout, data.(string))
	case reflect.TypeOf(time.Millisecond):
		return time.ParseDuration(data.(string))
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

func unmarshalWithDuration(b []byte, e interface{}) error {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	_, err := decodeData(data, e)
	return err
}
