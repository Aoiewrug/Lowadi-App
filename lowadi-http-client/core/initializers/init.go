package initializers

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVar() {
	err := godotenv.Load(filepath.Join("C:/Users/3/Desktop/oldLowadi/Lowadi-App/lowadi-http-client/", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
