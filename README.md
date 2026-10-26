# 📝 My Personal Notes

> Aplikasi backend full-stack untuk mencatat ide, kegiatan, dan pengeluaran. Membuat catatan bisa dibagi ke pengguna lain untuk dibaca bersama.

## 🎯 Deskripsi

My Personal Notes adalah aplikasi catatan pribadi berbasis web yang memungkinkan pengguna:
- Membuat akun / login
- Membuat, mengedit, menghapus catatan
- Menetapkan deadline pada catatan
- Membagikan catatan ke pengguna lain (read-only)
- Menerima notifikasi saat ada catatan yang dibagikan

## 🚀 Fitur

- ✅ Registrasi & Login akun
- ✅ JWT Authentication + Authorization
- ✅ Manajemen catatan (CRUD)
- ✅ Berbagi catatan (read-only) ke pengguna lain
- ✅ Sistem notifikasi (dibagikan ke saya)
- ✅ Frontend responsif (HTML/CSS/JS native)

## 🔧 Teknologi

| Komponen | Teknologi |
|---|---|
| Backend | Go + Gin |
| Database | PostgreSQL |
| ORM | GORM |
| Auth | JWT |
| Password Hash | bcrypt |
| Frontend | HTML, CSS, JavaScript (native) |
| Env Config | godotenv (.env) |

## 📦 Persiapan Lingkungan

Pastikan sudah terpasang:

- **Go** 1.22+ — cek: `go version`
- **PostgreSQL** 13+ — cek: `psql --version`

## ⬇️ Instalasi

```bash
# 1. Clone repo ini
git clone <repository-url>
cd go_vibe

# 2. Download dependencies Go
go mod tidy

# 3. Copy contoh env
cp .env.example .env
```

## ⚙️ Konfigurasi (.env)

Edit file `.env` yang baru dibuat:

```dotenv
APP_PORT=8080
DATABASE_URL=postgres://USERNAME:PASSWORD@localhost:5432/personal_notes?sslmode=disable
JWT_SECRET=ganti-dengan-secret-acak-yang-panjang
```

> 🔒 File `.env` tidak pernah disimpan ke Git (ada di `.gitignore`).

## ▶️ Cara Menjalankan

```bash
# Jalankan server
go run main.go
```

Server akan berjalan di `http://localhost:8080`.

> 💡 `AutoMigrate` otomatis membuat tabel database sesuai struct model. Pastikan database di `DATABASE_URL` sudah ada (bisa dibuat via `createdb personal_notes` atau pgAdmin).

## 📋 API Reference

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-----------|------|
| `POST` | `/api/auth/register` | Registrasi akun baru | ❌ |
| `POST` | `/api/auth/login` | Login & dapatkan token | ❌ |
| `GET` | `/api/auth/profile` | Lihat profil user | ✅ |
| `GET` | `/api/notes` | Daftar catatan sendiri | ✅ |
| `POST` | `/api/notes` | Buat catatan baru | ✅ |
| `GET` | `/api/notes/:id` | Detail catatan | ✅ |
| `PUT` | `/api/notes/:id` | Edit catatan | ✅ |
| `DELETE` | `/api/notes/:id` | Hapus catatan | ✅ |
| `GET` | `/api/notes/shared` | Catatan yang dibagikan ke saya | ✅ |
| `GET` | `/api/notes/:id/share` | Daftar user yang di-share | ✅ |
| `POST` | `/api/notes/:id/share` | Bagikan ke user lain | ✅ |
| `DELETE` | `/api/notes/:id/share/:userId` | Hapus akses user | ✅ |
| `GET` | `/api/notifications` | Daftar notifikasi | ✅ |
| `PUT` | `/api/notifications/:id/read` | Tandai satu dibaca | ✅ |
| `PUT` | `/api/notifications/read-all` | Tandai semua dibaca | ✅ |

## 🧪 Cara Menguji Aplikasi

### Akses Web
Buka browser: [http://localhost:8080](http://localhost:8080)

Flow:
1. Register akun baru → pilih "Daftar di sini"
2. Login dengan email & password
3. Setelah masuk, Anda akan diarahkan ke halaman utama `/app`
4. Klik `"+ Catatan Baru"` → isi judul, isi, deadline → Simpan
5. Klik kartu catatan untuk lihat detail
6. Share catatan: buka detail → bagian "Bagikan ke Pengguna Lain" → masukkan email → Tambah

### Test API pakai curl

```bash

# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@test.com","password":"123456"}'

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@test.com","password":"123456"}' | jq -r '.token')

# Buat catatan
curl -X POST http://localhost:8080/api/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Catatan Demo","content":"Hello World","deadline":"2026-12-31T23:59:59Z"}'

# Ambil daftar catatan
curl http://localhost:8080/api/notes \
  -H "Authorization: Bearer $TOKEN"
```

## 🐛 Troubleshooting

| Masalah | Solusi |
|---|---|
| `DATABASE_URL belum diset` | Salin `.env.example` → `.env`, isi dengan koneksi PostgreSQL Anda |
| `FATAL: password authentication failed` | Periksa username/password di `DATABASE_URL` |
| `database "personal_notes" does not exist` | Buat database: `createdb personal_notes` |
| `go: command not found` | Instal Go dari [golang.org](https://go.dev/dl/) |
| Frontend tidak load setelah login | Pastikan route `/app` mengarah ke file HTML |