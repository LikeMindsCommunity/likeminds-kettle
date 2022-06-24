package environment

import (
	"github.com/joho/godotenv"
	"log"
)

// GoDotEnvVariable to load/read the .env file and return the value of the key
func GoDotEnvVariable(key string) string {
	// load .env file
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatalf("Error reading .env file")
	}
	return envs[key]
}
