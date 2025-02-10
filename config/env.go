package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
	//envPath := "./env"
	err := godotenv.Load("env/appenv")
	if err != nil {
		log.Fatalf("Error loading appenv file")
	}
	// fmt.Println("DB_USER", os.Getenv("DB_USER"))
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
