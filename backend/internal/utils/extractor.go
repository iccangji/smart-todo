package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetString(m bson.M, key string) (string, bool) {
	val, ok := m[key]
	if !ok || val == nil {
		return "", false
	}

	str, ok := val.(string)
	if !ok || str == "" {
		return "", false
	}

	return str, true
}
func GetBool(m bson.M, key string) (bool, bool) {
	val, ok := m[key]
	if !ok || val == nil {
		return false, false
	}

	b, ok := val.(bool)
	if !ok {
		return false, false
	}

	return b, true
}
func GetTime(m bson.M, key string) (time.Time, bool) {
	val, ok := m[key]
	if !ok || val == nil {
		return time.Time{}, false
	}

	t, ok := val.(time.Time)
	if !ok {
		return time.Time{}, false
	}

	return t, true
}
