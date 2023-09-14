package utils

import "os"

func EnvOrDefault(key, v string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return v
	}
	return envVar
}