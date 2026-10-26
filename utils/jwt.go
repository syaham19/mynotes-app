// Package utils berisi fungsi-fungsi pembantu yang dipakai di banyak tempat.
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret menyimpan secret key untuk menandatangani token.
// Diset sekali saat aplikasi dimulai lewat fungsi SetJWTSecret.
var jwtSecret []byte

// SetJWTSecret mengatur secret key yang digunakan untuk JWT.
// Dipanggil sekali dari main.go saat aplikasi dimulai.
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// JWTClaims adalah data yang disimpan di dalam JWT token.
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT token baru untuk pengguna yang berhasil login.
// Token berlaku selama 24 jam.
func GenerateToken(userID uint, email string) (string, error) {
	// Buat claims (data di dalam token)
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Berlaku 24 jam
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan metode signing HMAC SHA-256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token dengan secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi JWT token dan mengembalikan claims jika valid.
// Mengembalikan error jika token kadaluarsa, format salah, atau tidak valid.
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse token dan validasi
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing sesuai (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Ambil claims dari token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}

