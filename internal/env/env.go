package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)

	if err != nil {
		log.Println(err)
		return fallback
	}

	return valAsInt
}
func GetDuration(key string, fallback string) time.Duration {
	value, ok := os.LookupEnv(key)

	if !ok {
		duration, err := time.ParseDuration(fallback)
		if err != nil {
			return 15 * time.Minute
		}
		return duration
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 15 * time.Minute
	}
	return duration

}
func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}
	boolValue, err := strconv.ParseBool(value)

	if err != nil {
		return fallback
	}

	return boolValue
}
