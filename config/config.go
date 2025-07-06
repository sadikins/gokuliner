package config

import (
	"fmt" // <<< Tambahkan import ini
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

// LoadConfig membaca variabel lingkungan dari file .env
// Sekarang mengembalikan (Config, error)
func LoadConfig() (Config, error) { // <<< UBAH TANDA TANGAN FUNGSI INI
	err := godotenv.Load()
	if err != nil {
		// Jika file .env tidak ditemukan, itu bukan error fatal jika variabel env sudah diset
		log.Println("Error loading .env file, assuming environment variables are set: ", err)
		// Kita tetap kembalikan Config, tapi dengan error jika ingin menanganinya
	}

	// Lakukan validasi sederhana di sini jika perlu
	// Contoh: jika DB_HOST kosong, kembalikan error
	if getEnv("DB_HOST", "") == "" { // Asumsi jika DB_HOST kosong, itu error
		return Config{}, fmt.Errorf("DB_HOST environment variable is not set") // <<< KEMBALIKAN ERROR
	}


	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "sandi"),
		DBPassword: getEnv("DB_PASSWORD", "minyakkayuputih"),
		DBName:     getEnv("DB_NAME", "vinhpp_db"),
		DBPort:     getEnv("DB_PORT", "5432"), // 3306 for MySQL
		AppPort:    getEnv("APP_PORT", "8080"),
	}, nil // <<< KEMBALIKAN nil UNTUK ERROR JIKA BERHASIL
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}