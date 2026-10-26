// ===========================================
// auth.js — Logika halaman Login & Register
// ===========================================

// URL dasar API
const API_URL = '';

// ---- Fungsi Bantu: Tampilkan Alert ----
function showAlert(message, type = 'error') {
    const alert = document.getElementById('alert');
    alert.textContent = message;
    alert.className = `alert alert-${type}`;
    alert.style.display = 'block';

    // Sembunyikan setelah 5 detik
    setTimeout(() => {
        alert.style.display = 'none';
    }, 5000);
}

// ---- Cek apakah sudah login ----
// Jika sudah login, langsung ke halaman utama
(function checkAuth() {
    const token = localStorage.getItem('token');
    if (token) {
        window.location.href = '/app';
    }
})();

// ---- Form Login ----
const loginForm = document.getElementById('login-form');
if (loginForm) {
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const btnLogin = document.getElementById('btn-login');

        // Nonaktifkan tombol sementara
        btnLogin.disabled = true;
        btnLogin.textContent = 'Memproses...';

        try {
            const response = await fetch(`${API_URL}/api/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });

            const data = await response.json();

            if (response.ok) {
                // Simpan token dan data user di localStorage
                localStorage.setItem('token', data.token);
                localStorage.setItem('user', JSON.stringify(data.user));

                // Redirect ke halaman utama
                window.location.href = '/app';
            } else {
                showAlert(data.error || 'Login gagal.');
            }
        } catch (err) {
            showAlert('Tidak bisa terhubung ke server.');
        } finally {
            btnLogin.disabled = false;
            btnLogin.textContent = 'Masuk';
        }
    });
}

// ---- Form Register ----
const registerForm = document.getElementById('register-form');
if (registerForm) {
    registerForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        const btnRegister = document.getElementById('btn-register');

        // Validasi: password harus cocok
        if (password !== confirmPassword) {
            showAlert('Password dan konfirmasi password tidak cocok.');
            return;
        }

        btnRegister.disabled = true;
        btnRegister.textContent = 'Memproses...';

        try {
            const response = await fetch(`${API_URL}/api/auth/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, email, password })
            });

            const data = await response.json();

            if (response.ok) {
                showAlert('Registrasi berhasil! Silakan login.', 'success');
                // Redirect ke halaman login setelah 2 detik
                setTimeout(() => {
                    window.location.href = '/';
                }, 2000);
            } else {
                showAlert(data.error || 'Registrasi gagal.');
            }
        } catch (err) {
            showAlert('Tidak bisa terhubung ke server.');
        } finally {
            btnRegister.disabled = false;
            btnRegister.textContent = 'Daftar';
        }
    });
}

