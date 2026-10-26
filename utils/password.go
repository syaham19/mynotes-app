package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengubah password teks biasa menjadi hash menggunakan bcrypt.
// Hasil hash aman disimpan di database karena tidak bisa dikembalikan ke teks asli.
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10, cukup aman untuk kebanyakan aplikasi
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword membandingkan password teks biasa dengan hash yang tersimpan di database.
// Mengembalikan true jika cocok, false jika tidak.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

