package utils

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func StructToBsonM(input any) bson.M {
	result := bson.M{}

	val := reflect.ValueOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if !field.CanInterface() {
			continue
		}

		tag := structField.Tag.Get("json")

		key := strings.Split(tag, ",")[0]

		if key == "" || key == "-" {
			continue
		}

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		if field.Kind() == reflect.Ptr {
			result[key] = field.Elem().Interface()
		} else {
			result[key] = field.Interface()
		}
	}

	return result
}
