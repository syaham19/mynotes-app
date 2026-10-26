package models

import "time"

// Note merepresentasikan tabel "notes" di database.
// Setiap catatan dimiliki oleh satu pengguna (user_id).
// Deadline bersifat opsional (pointer, bisa nil).
type Note struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`       // Pemilik catatan
	Title     string     `json:"title" gorm:"not null"`               // Judul catatan
	Content   string     `json:"content" gorm:"type:text"`            // Isi catatan
	Deadline  *time.Time `json:"deadline,omitempty"`                   // Tenggat waktu (opsional)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relasi: satu catatan bisa dimiliki oleh satu user
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

