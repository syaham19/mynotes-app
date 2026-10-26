package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"go_vibe/config"
	"go_vibe/models"

	"github.com/gin-gonic/gin"
)

// ---- Struct untuk input ----

// ShareInput adalah data yang dikirim saat berbagi catatan.
type ShareInput struct {
	Email string `json:"email" binding:"required,email"` // Email pengguna yang akan ditambahkan
}

// ---- Handler Functions ----

// ShareNote menambahkan pengguna lain ke catatan (read-only).
// Hanya pemilik catatan yang bisa berbagi.
// Endpoint: POST /api/notes/:id/share
func ShareNote(c *gin.Context) {
	userID := getUserID(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID catatan tidak valid."})
		return
	}

	// Pastikan catatan ada dan milik user yang sedang login
	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan."})
		return
	}
	if note.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya pemilik catatan yang bisa berbagi."})
		return
	}

	var input ShareInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak valid."})
		return
	}

	// Cari user yang akan ditambahkan
	var targetUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna dengan email tersebut tidak ditemukan."})
		return
	}

	// Tidak boleh berbagi ke diri sendiri
	if targetUser.ID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak bisa berbagi catatan ke diri sendiri."})
		return
	}

	// Cek apakah sudah dibagikan sebelumnya
	var existingShare models.SharedNote
	result := config.DB.Where("note_id = ? AND shared_with_user_id = ?", noteID, targetUser.ID).First(&existingShare)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Catatan sudah dibagikan ke pengguna ini."})
		return
	}

	// Buat record shared note
	sharedNote := models.SharedNote{
		NoteID:           uint(noteID),
		SharedWithUserID: targetUser.ID,
	}

	if err := config.DB.Create(&sharedNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal berbagi catatan."})
		return
	}

	// Ambil nama pemilik catatan untuk pesan notifikasi
	var owner models.User
	config.DB.First(&owner, userID)

	// Buat notifikasi untuk pengguna yang ditambahkan
	notification := models.Notification{
		UserID:        targetUser.ID,
		Message:       fmt.Sprintf("%s membagikan catatan \"%s\" kepada Anda.", owner.Name, note.Title),
		Type:          models.NotifShared,
		RelatedNoteID: &note.ID,
	}
	config.DB.Create(&notification)

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Catatan berhasil dibagikan ke %s!", targetUser.Name),
	})
}

// UnshareNote menghapus akses pengguna dari catatan.
// Hanya pemilik catatan yang bisa menghapus akses.
// Endpoint: DELETE /api/notes/:id/share/:userId
func UnshareNote(c *gin.Context) {
	userID := getUserID(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID catatan tidak valid."})
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid."})
		return
	}

	// Pastikan catatan milik user yang sedang login
	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan."})
		return
	}
	if note.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya pemilik catatan yang bisa menghapus akses."})
		return
	}

	// Hapus shared note
	result := config.DB.Where("note_id = ? AND shared_with_user_id = ?", noteID, targetUserID).Delete(&models.SharedNote{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak memiliki akses ke catatan ini."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Akses pengguna berhasil dihapus.",
	})
}

// GetSharedUsers mengembalikan daftar pengguna yang memiliki akses ke catatan.
// Endpoint: GET /api/notes/:id/share
func GetSharedUsers(c *gin.Context) {
	userID := getUserID(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID catatan tidak valid."})
		return
	}

	// Pastikan catatan milik user yang sedang login
	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan."})
		return
	}
	if note.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya pemilik catatan yang bisa melihat daftar akses."})
		return
	}

	// Ambil daftar user yang dibagikan
	var sharedNotes []models.SharedNote
	config.DB.Where("note_id = ?", noteID).Find(&sharedNotes)

	var userIDs []uint
	for _, sn := range sharedNotes {
		userIDs = append(userIDs, sn.SharedWithUserID)
	}

	var users []models.User
	if len(userIDs) > 0 {
		config.DB.Where("id IN ?", userIDs).Find(&users)
	}

	// Kembalikan data user tanpa password hash
	var result []gin.H
	for _, u := range users {
		result = append(result, gin.H{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"shared_users": result,
	})
}

