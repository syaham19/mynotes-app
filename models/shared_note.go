package models

import "time"

// SharedNote merepresentasikan tabel "shared_notes" di database.
// Menyimpan data catatan yang dibagikan ke pengguna lain.
// Pengguna yang ditambahkan hanya bisa MELIHAT catatan (read-only).
type SharedNote struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	NoteID           uint      `json:"note_id" gorm:"not null;index"`             // Catatan yang dibagikan
	SharedWithUserID uint      `json:"shared_with_user_id" gorm:"not null;index"` // Pengguna yang menerima akses
	CreatedAt        time.Time `json:"created_at"`

	// Relasi
	Note Note `json:"note,omitempty" gorm:"foreignKey:NoteID"`
	User User `json:"user,omitempty" gorm:"foreignKey:SharedWithUserID"`
}

