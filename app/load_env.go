package app

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to load environment config file. \n%s", err.Error()))
	}
}
