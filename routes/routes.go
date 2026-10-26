// Package routes mendaftarkan semua endpoint API ke router Gin.
package routes

import (
	"go_vibe/handlers"
	"go_vibe/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mendaftarkan semua route API.
// Route dibagi menjadi 2 kelompok:
//   - Public: bisa diakses tanpa login (register, login)
//   - Protected: harus login dulu (notes, share, notifications)
func SetupRoutes(router *gin.Engine) {

	// ---- Route API ----
	api := router.Group("/api")
	{
		// -- Auth (Public) --
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register) // Daftar akun baru
			auth.POST("/login", handlers.Login)        // Login
		}

		// -- Protected Routes (butuh JWT token) --
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth - profil
			protected.GET("/auth/profile", handlers.GetProfile) // Lihat profil sendiri

			// Notes - CRUD catatan
			protected.GET("/notes", handlers.GetNotes)           // Daftar catatan milik sendiri
			protected.GET("/notes/shared", handlers.GetSharedNotes) // Daftar catatan yang dibagikan ke saya
			protected.POST("/notes", handlers.CreateNote)         // Buat catatan baru
			protected.GET("/notes/:id", handlers.GetNoteByID)     // Lihat detail catatan
			protected.PUT("/notes/:id", handlers.UpdateNote)      // Edit catatan
			protected.DELETE("/notes/:id", handlers.DeleteNote)   // Hapus catatan

			// Share - berbagi catatan
			protected.GET("/notes/:id/share", handlers.GetSharedUsers)       // Daftar user yang dibagikan
			protected.POST("/notes/:id/share", handlers.ShareNote)           // Bagikan catatan ke user lain
			protected.DELETE("/notes/:id/share/:userId", handlers.UnshareNote) // Hapus akses user

			// Notifications
			protected.GET("/notifications", handlers.GetNotifications)              // Daftar notifikasi
			protected.PUT("/notifications/:id/read", handlers.MarkNotificationRead) // Tandai satu sebagai dibaca
			protected.PUT("/notifications/read-all", handlers.MarkAllNotificationsRead) // Tandai semua sebagai dibaca
		}
	}
}

