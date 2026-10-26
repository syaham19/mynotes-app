// ===========================================
// app.js — Logika utama halaman catatan (app.html)
// Mengelola: auth check, navigasi tab, sidebar, logout
// ===========================================

// ---- Variabel Global ----
const API_URL = '';
let currentTab = 'my-notes'; // Tab aktif saat ini

// ---- Fungsi Bantu: Header dengan token ----
function authHeaders() {
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
    };
}

// ---- Fungsi Bantu: Cek apakah response 401 (token kadaluarsa) ----
function handleUnauthorized(response) {
    if (response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/';
        return true;
    }
    return false;
}

// ---- Cek Auth: Redirect ke login jika belum login ----
(function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

    // Tampilkan nama user
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const userNameEl = document.getElementById('user-name');
    if (userNameEl && user.name) {
        userNameEl.textContent = `Halo, ${user.name}!`;
    }
})();

// ---- Logout ----
document.getElementById('btn-logout').addEventListener('click', function() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/';
});

// ---- Tab Navigation ----
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', function() {
        // Update tab aktif
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        this.classList.add('active');

        currentTab = this.dataset.tab;

        // Tampilkan/sembunyikan toolbar (tombol tambah hanya di "Catatan Saya")
        const toolbar = document.getElementById('toolbar-my-notes');
        toolbar.style.display = currentTab === 'my-notes' ? 'flex' : 'none';

        // Muat catatan sesuai tab
        if (currentTab === 'my-notes') {
            loadMyNotes();
        } else {
            loadSharedNotes();
        }
    });
});

// ---- Sidebar: Buka & Tutup ----
function openSidebar() {
    document.getElementById('sidebar').classList.add('active');
    document.getElementById('sidebar-overlay').classList.add('active');
}

function closeSidebar() {
    document.getElementById('sidebar').classList.remove('active');
    document.getElementById('sidebar-overlay').classList.remove('active');
    // Reset form
    document.getElementById('note-form').reset();
    document.getElementById('note-id').value = '';
    document.getElementById('share-section').style.display = 'none';
    document.getElementById('btn-delete-note').style.display = 'none';
    document.getElementById('readonly-view').style.display = 'none';
    document.getElementById('note-form').style.display = 'block';
    document.getElementById('sidebar-actions-owner').style.display = 'flex';
}

// Tombol tutup sidebar
document.getElementById('btn-close-sidebar').addEventListener('click', closeSidebar);
document.getElementById('sidebar-overlay').addEventListener('click', closeSidebar);

// Tombol tambah catatan baru
document.getElementById('btn-add-note').addEventListener('click', function() {
    document.getElementById('sidebar-title').textContent = 'Catatan Baru';
    document.getElementById('note-id').value = '';
    document.getElementById('note-form').reset();
    document.getElementById('share-section').style.display = 'none';
    document.getElementById('btn-delete-note').style.display = 'none';
    document.getElementById('readonly-view').style.display = 'none';
    document.getElementById('note-form').style.display = 'block';
    document.getElementById('sidebar-actions-owner').style.display = 'flex';
    openSidebar();
});

// ---- Format Tanggal ----
function formatDate(dateStr) {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleDateString('id-ID', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function formatDateShort(dateStr) {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });
}

// ---- Cek apakah deadline sudah lewat ----
function isOverdue(deadlineStr) {
    if (!deadlineStr) return false;
    return new Date(deadlineStr) < new Date();
}

// ---- Muat catatan saat halaman dibuka ----
document.addEventListener('DOMContentLoaded', function() {
    loadMyNotes();
    loadNotifications();
});

