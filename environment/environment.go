package environment

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

// GoDotEnvVariable to load/read the .env file and return the value of the key
func GoDotEnvVariable(key string) string {
	// load .env file
	dir, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		log.Fatal(err)
	}
	environmentPath := filepath.Join(dir, ".env")
	envs, err := godotenv.Read(environmentPath)

	if err != nil {
		log.Fatalf("Error reading .env file")
	}
	return envs[key]
}
