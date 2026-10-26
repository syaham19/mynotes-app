package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_vibe/config"
	"go_vibe/models"

	"github.com/gin-gonic/gin"
)

// ---- Struct untuk input ----

// NoteInput adalah data yang dikirim saat membuat atau mengedit catatan.
type NoteInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	Deadline string `json:"deadline"` // Format: "2006-01-02T15:04:05Z" (ISO 8601), kosong jika tidak ada
}

// ---- Helper: ambil user_id dari context ----

// getUserID mengambil user_id dari Gin context yang sudah diset oleh middleware auth.
func getUserID(c *gin.Context) uint {
	userID, _ := c.Get("user_id")
	return userID.(uint)
}

// ---- Handler Functions ----

// GetNotes mengembalikan semua catatan milik pengguna yang sedang login.
// Endpoint: GET /api/notes
func GetNotes(c *gin.Context) {
	userID := getUserID(c)

	var notes []models.Note
	// Ambil catatan milik user, urutkan dari yang terbaru
	config.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&notes)

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

// GetSharedNotes mengembalikan catatan yang dibagikan ke pengguna yang sedang login.
// Endpoint: GET /api/notes/shared
func GetSharedNotes(c *gin.Context) {
	userID := getUserID(c)

	// Cari catatan yang dibagikan ke user ini
	var sharedNotes []models.SharedNote
	config.DB.Where("shared_with_user_id = ?", userID).Find(&sharedNotes)

	// Ambil detail catatan untuk setiap shared note
	var noteIDs []uint
	for _, sn := range sharedNotes {
		noteIDs = append(noteIDs, sn.NoteID)
	}

	var notes []models.Note
	if len(noteIDs) > 0 {
		// Preload User agar nama pemilik catatan ikut dimuat
		config.DB.Preload("User").Where("id IN ?", noteIDs).Order("created_at DESC").Find(&notes)
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

// GetNoteByID mengembalikan detail satu catatan.
// Bisa diakses oleh pemilik ATAU pengguna yang dibagikan akses.
// Endpoint: GET /api/notes/:id
func GetNoteByID(c *gin.Context) {
	userID := getUserID(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID catatan tidak valid."})
		return
	}

	var note models.Note
	if err := config.DB.Preload("User").First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan."})
		return
	}

	// Cek akses: pemilik ATAU pengguna yang dibagikan
	if note.UserID != userID {
		// Cek apakah catatan dibagikan ke user ini
		var shared models.SharedNote
		result := config.DB.Where("note_id = ? AND shared_with_user_id = ?", noteID, userID).First(&shared)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke catatan ini."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"note":     note,
		"is_owner": note.UserID == userID,
	})
}

// CreateNote membuat catatan baru milik pengguna yang sedang login.
// Endpoint: POST /api/notes
func CreateNote(c *gin.Context) {
	userID := getUserID(c)

	var input NoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid. Pastikan judul terisi."})
		return
	}

	note := models.Note{
		UserID:  userID,
		Title:   input.Title,
		Content: input.Content,
	}

	// Parse deadline jika ada
	if input.Deadline != "" {
		deadline, err := parseDeadline(input.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format tenggat waktu tidak valid. Gunakan format: 2006-01-02T15:04:05Z"})
			return
		}
		note.Deadline = &deadline
	}

	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat catatan. Silakan coba lagi."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Catatan berhasil dibuat!",
		"note":    note,
	})
}

// UpdateNote mengedit catatan. Hanya pemilik yang bisa mengedit.
// Endpoint: PUT /api/notes/:id
func UpdateNote(c *gin.Context) {
	userID := getUserID(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID catatan tidak valid."})
		return
	}

	// Cari catatan dan pastikan milik user yang sedang login
	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan."})
		return
	}

	// Hanya pemilik yang bisa edit
	if note.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk mengedit catatan ini."})
		return
	}

	var input NoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid. Pastikan judul terisi."})
		return
	}

	// Update field
	note.Title = input.Title
	note.Content = input.Content

	// Update deadline
	if input.Deadline != "" {
		deadline, err := parseDeadline(input.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format tenggat waktu tidak valid."})
			return
		}
		note.Deadline = &deadline
	} else {
		note.Deadline = nil // Hapus deadline jika dikosongkan
	}

	if err := config.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate catatan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Catatan berhasil diupdate!",
		"note":    note,
	})
}

// DeleteNote menghapus catatan. Hanya pemilik yang bisa menghapus.
// Endpoint: DELETE /api/notes/:id
func DeleteNote(c *gin.Context) {
	userID := getUserID(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID catatan tidak valid."})
		return
	}

	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan."})
		return
	}

	// Hanya pemilik yang bisa hapus
	if note.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk menghapus catatan ini."})
		return
	}

	// Hapus juga data shared_notes dan notifications terkait
	config.DB.Where("note_id = ?", noteID).Delete(&models.SharedNote{})
	config.DB.Where("related_note_id = ?", noteID).Delete(&models.Notification{})

	// Hapus catatan
	config.DB.Delete(&note)

	c.JSON(http.StatusOK, gin.H{
		"message": "Catatan berhasil dihapus!",
	})
}

// parseDeadline mengubah string deadline ke time.Time.
// Mendukung beberapa format umum.
func parseDeadline(s string) (time.Time, error) {
	// Coba beberapa format yang umum dipakai
	formats := []string{
		"2006-01-02T15:04:05Z",       // ISO 8601 UTC
		"2006-01-02T15:04:05+07:00",  // ISO 8601 dengan timezone
		"2006-01-02T15:04:05-07:00",  // ISO 8601 dengan timezone negatif
		"2006-01-02T15:04",           // Tanpa detik
		"2006-01-02",                 // Tanggal saja
	}

	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("format tidak dikenali: %s", s)
}
