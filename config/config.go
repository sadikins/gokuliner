package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	AppPort    string
}


//LoadConfig membaca variabel lingkungan dari file .env

func LoadConfig() Config{
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, assuming enviroment variables are set.")
	}

	return Config{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBUser: getEnv("DB_USER", "xnfo"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName: getEnv("DB_NAME", "gomakan"),
		DBPort: getEnv("DB_PORT", "5432"),
		AppPort: getEnv("APP_PORT", "8080"),

	}

}


func getEnv(key string, defaultValue string) string{
	if value, exists := os.LookupEnv(key); exists{
		return value
	}

	return defaultValue
}