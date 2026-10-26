package handlers

import (
	"net/http"
	"strconv"

	"go_vibe/config"
	"go_vibe/models"

	"github.com/gin-gonic/gin"
)

// GetNotifications mengembalikan daftar notifikasi pengguna yang sedang login.
// Endpoint: GET /api/notifications
func GetNotifications(c *gin.Context) {
	userID := getUserID(c)

	var notifications []models.Notification
	// Ambil notifikasi, urutkan dari yang terbaru
	config.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications)

	// Hitung jumlah yang belum dibaca
	var unreadCount int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

// MarkNotificationRead menandai satu notifikasi sebagai sudah dibaca.
// Endpoint: PUT /api/notifications/:id/read
func MarkNotificationRead(c *gin.Context) {
	userID := getUserID(c)
	notifID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID notifikasi tidak valid."})
		return
	}

	var notification models.Notification
	if err := config.DB.First(&notification, notifID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan."})
		return
	}

	// Pastikan notifikasi milik user yang sedang login
	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke notifikasi ini."})
		return
	}

	// Tandai sebagai sudah dibaca
	config.DB.Model(&notification).Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Notifikasi ditandai sebagai sudah dibaca.",
	})
}

// MarkAllNotificationsRead menandai semua notifikasi pengguna sebagai sudah dibaca.
// Endpoint: PUT /api/notifications/read-all
func MarkAllNotificationsRead(c *gin.Context) {
	userID := getUserID(c)

	config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Semua notifikasi ditandai sebagai sudah dibaca.",
	})
}

