// Package models berisi semua struct yang merepresentasikan tabel di database.
// GORM akan membuat tabel secara otomatis berdasarkan struct di sini.
package models

import "time"

// User merepresentasikan tabel "users" di database.
// Setiap pengguna punya nama, email (unik), dan password yang di-hash.
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"` // json:"-" agar tidak muncul di response API
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

