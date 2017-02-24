package main

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

func StringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Now()) {
			return data, nil
		}
		return time.Parse(time.RFC3339, data.(string))
	}
}
