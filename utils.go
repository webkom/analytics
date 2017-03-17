package main

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
	"strconv"
)

func NormalizeTypesHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() == reflect.Float64 && t.Kind() == reflect.String {
			return strconv.FormatFloat(data.(float64), 'f', -1, 64), nil
		}

		if f.Kind() == reflect.Int && t.Kind() == reflect.String {
			return strconv.Itoa(data.(int)), nil

		}

		if t == reflect.TypeOf(time.Now()) {
			return time.Parse(time.RFC3339, data.(string))
		}

		return data, nil
	}
}
