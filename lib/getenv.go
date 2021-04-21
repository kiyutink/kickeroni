package lib

import "os"

func Getenv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic("Missing environment variable " + name)
	}

	return value
}
