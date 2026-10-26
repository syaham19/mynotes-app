// My Personal Notes — Entry Point
//
// File ini adalah titik awal aplikasi. Tugasnya:
// 1. Memuat konfigurasi dari .env
// 2. Menghubungkan ke database PostgreSQL
// 3. Membuat tabel otomatis (auto-migrate)
// 4. Menjalankan server Gin
package main

import (
	"fmt"
	"log"

	"go_vibe/config"
	"go_vibe/models"
	"go_vibe/routes"
	"go_vibe/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Muat konfigurasi dari .env
	cfg := config.LoadConfig()
	log.Printf("Konfigurasi dimuat. Server akan berjalan di port %s", cfg.Port)

	// 2. Set JWT secret
	utils.SetJWTSecret(cfg.JWTSecret)

	// 3. Hubungkan ke database
	config.ConnectDatabase(cfg.DatabaseURL)

	// 4. Auto-migrate: buat tabel otomatis dari struct model
	// GORM akan membuat tabel jika belum ada, dan menambahkan kolom baru
	// jika struct berubah (tidak menghapus kolom yang sudah ada).
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Note{},
		&models.SharedNote{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Gagal melakukan auto-migrate database: %v", err))
	}
	log.Println("Auto-migrate database selesai")

	// 5. Buat router Gin
	router := gin.Default()

	// 6. Sajikan file frontend (HTML, CSS, JS) dari folder web/
	router.Static("/css", "./web/css")
	router.Static("/js", "./web/js")
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/register", "./web/register.html")
	router.StaticFile("/app", "./web/app.html")

	// 7. Daftarkan semua route API
	routes.SetupRoutes(router)

	// 8. Jalankan server
	log.Printf("Server berjalan di http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(fmt.Sprintf("Error: Gagal menjalankan server: %v", err))
	}
}

