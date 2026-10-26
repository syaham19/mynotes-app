// Package handlers berisi fungsi-fungsi yang menangani request HTTP.
package handlers

import (
	"net/http"

	"go_vibe/config"
	"go_vibe/models"
	"go_vibe/utils"

	"github.com/gin-gonic/gin"
)

// ---- Struct untuk menerima data dari request body ----

// RegisterInput adalah data yang dikirim saat registrasi.
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput adalah data yang dikirim saat login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ---- Handler Functions ----

// Register menangani pendaftaran pengguna baru.
// Endpoint: POST /api/auth/register
func Register(c *gin.Context) {
	var input RegisterInput

	// Validasi input dari request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid. Pastikan nama, email, dan password (min 6 karakter) terisi."})
		return
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	result := config.DB.Where("email = ?", input.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar. Silakan gunakan email lain."})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memproses password."})
		return
	}

	// Buat user baru
	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	// Simpan ke database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun. Silakan coba lagi."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil! Silakan login.",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Login menangani proses login pengguna.
// Jika berhasil, mengembalikan JWT token.
// Endpoint: POST /api/auth/login
func Login(c *gin.Context) {
	var input LoginInput

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid. Pastikan email dan password terisi."})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		// Pesan error umum agar tidak membocorkan apakah email terdaftar atau tidak
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah."})
		return
	}

	// Cek password
	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah."})
		return
	}

	// Buat JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token. Silakan coba lagi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// GetProfile mengembalikan data profil pengguna yang sedang login.
// Endpoint: GET /api/auth/profile
func GetProfile(c *gin.Context) {
	// Ambil user_id dari context (diset oleh middleware)
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

