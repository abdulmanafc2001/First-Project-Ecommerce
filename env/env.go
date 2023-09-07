package env

import (
	"log"

	"github.com/joho/godotenv"
)

func Loadenv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Panic("Failed to load env file")
	}
}
