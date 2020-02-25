package utils

import "os"

func GetEnvOrDie(key string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}
	panic("invalid env parameter: " + key)
}
