package utils

import (
	"fmt"
	"os"
)

func EnvOrDefault(key, v string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return v
	}
	return envVar
}

func EnvOrPanic(key string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		panic(fmt.Sprintf("must provide environment variable %s", key))
	}
	return envVar
}