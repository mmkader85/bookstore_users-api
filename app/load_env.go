package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// loadEnv reads .env when present. In containers the variables are usually
// injected directly, so a missing file is not an error.
func loadEnv() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Unable to load environment config file. \n%s", err.Error())
	}
}
