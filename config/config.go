// Package config berisi fungsi untuk memuat konfigurasi aplikasi
// dari environment variable dan membuat koneksi ke database.
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AppConfig menyimpan semua konfigurasi aplikasi.
type AppConfig struct {
	Port        string // Port server (default: 8080)
	DatabaseURL string // URL koneksi PostgreSQL
	JWTSecret   string // Secret key untuk JWT token
}

// DB adalah variabel global untuk koneksi database.
// Digunakan oleh handler dan bagian lain aplikasi.
var DB *gorm.DB

// LoadConfig membaca konfigurasi dari file .env dan environment variable.
// Mengembalikan AppConfig yang berisi semua pengaturan aplikasi.
// Akan menghentikan aplikasi jika konfigurasi wajib tidak ditemukan.
func LoadConfig() AppConfig {
	// Coba muat file .env. Tidak error jika file tidak ada
	// (misalnya di lingkungan produksi yang pakai env var langsung).
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// Baca setiap konfigurasi dari environment variable
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Nilai default jika tidak diset
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Error: DATABASE_URL belum diset. Silakan atur di file .env atau environment variable.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Error: JWT_SECRET belum diset. Silakan atur di file .env atau environment variable.")
	}

	return AppConfig{
		Port:        port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
	}
}

// ConnectDatabase membuat koneksi ke PostgreSQL menggunakan GORM.
// Menyimpan koneksi ke variabel global DB.
// Akan menghentikan aplikasi jika koneksi gagal.
func ConnectDatabase(databaseURL string) {
	var err error

	// Buka koneksi ke PostgreSQL
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Gagal terhubung ke database: %v", err))
	}

	log.Println("Berhasil terhubung ke database PostgreSQL")
}

