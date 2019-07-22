package util

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"reflect"
	"time"
)

func UnmarshalConfig(ptr interface{}) error {
	if err := viper.Unmarshal(ptr, func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(config.DecodeHook, ToTimeHookFunc())
	}); err != nil {
		return err
	}
	return nil
}

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			s := data.(string)
			if s == "" {
				return time.Time{}, nil
			}
			return cast.ToTimeE(data)
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}
