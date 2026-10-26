package models

import "time"

// NotificationType mendefinisikan jenis-jenis notifikasi.
const (
	// NotifShared: Notifikasi saat pengguna ditambahkan ke catatan orang lain.
	NotifShared = "shared"

	// NotifDeadline: Notifikasi pengingat H-1 sebelum tenggat waktu catatan.
	NotifDeadline = "deadline"
)

// Notification merepresentasikan tabel "notifications" di database.
// Menyimpan notifikasi untuk setiap pengguna.
type Notification struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;index"` // Pengguna yang menerima notifikasi
	Message       string    `json:"message" gorm:"not null"`       // Isi pesan notifikasi
	Type          string    `json:"type" gorm:"not null"`          // Jenis: "shared" atau "deadline"
	IsRead        bool      `json:"is_read" gorm:"default:false"`  // Sudah dibaca atau belum
	RelatedNoteID *uint     `json:"related_note_id,omitempty"`     // Catatan terkait (opsional)
	CreatedAt     time.Time `json:"created_at"`

	// Relasi
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

