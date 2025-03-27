package env

import (
	"os"
	"strconv"
)

/*
The main purpose of these functions is to
retrieve environment variables with fallback values
if the variables are not set or cannot be converted to the desired type.
*/

func GetString(key, fallback string) string {
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
		return fallback
	}
	return valAsInt

}
