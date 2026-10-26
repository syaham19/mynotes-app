// Package middleware berisi fungsi-fungsi yang dijalankan sebelum handler.
package middleware

import (
	"net/http"
	"strings"

	"go_vibe/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memeriksa apakah request memiliki JWT token yang valid.
// Jika valid, user_id dan email disimpan di context Gin agar bisa diakses handler.
// Jika tidak valid, request ditolak dengan status 401 Unauthorized.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan. Silakan login terlebih dahulu."})
			c.Abort()
			return
		}

		// Format header: "Bearer <token>"
		// Ambil bagian token saja
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid. Gunakan: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa. Silakan login ulang."})
			c.Abort()
			return
		}

		// Simpan data user di context agar bisa diakses oleh handler
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}

