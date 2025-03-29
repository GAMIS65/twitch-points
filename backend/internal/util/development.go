package util

import "os"

func IsDev() bool {
	environment := os.Getenv("ENVIRONMENT")

	return environment == "DEVELOPMENT"

}
