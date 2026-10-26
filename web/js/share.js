// ===========================================
// share.js — Logika berbagi catatan ke pengguna lain
// ===========================================

// ---- Muat Daftar User yang Dibagikan ----
async function loadSharedUsers(noteId) {
    const list = document.getElementById('shared-users-list');
    list.innerHTML = '';

    try {
        const response = await fetch(`${API_URL}/api/notes/${noteId}/share`, {
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        if (!response.ok) return;

        const data = await response.json();
        const users = data.shared_users || [];

        if (users.length === 0) {
            list.innerHTML = '<li style="color: #aaa; font-size: 13px;">Belum dibagikan ke siapa pun.</li>';
            return;
        }

        users.forEach(user => {
            const li = document.createElement('li');
            li.className = 'shared-user-item';
            li.innerHTML = `
                <span>${escapeHtml(user.name)} (${escapeHtml(user.email)})</span>
                <button class="btn-remove" onclick="removeSharedUser(${noteId}, ${user.id})" title="Hapus akses">✕</button>
            `;
            list.appendChild(li);
        });

    } catch (err) {
        list.innerHTML = '<li style="color: #e74c3c; font-size: 13px;">Gagal memuat daftar.</li>';
    }
}

// ---- Tambah User ke Catatan ----
document.getElementById('btn-share').addEventListener('click', async function() {
    const noteId = document.getElementById('note-id').value;
    const emailInput = document.getElementById('share-email');
    const email = emailInput.value.trim();

    if (!noteId || !email) {
        alert('Masukkan email pengguna yang ingin ditambahkan.');
        return;
    }

    try {
        const response = await fetch(`${API_URL}/api/notes/${noteId}/share`, {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify({ email })
        });

        if (handleUnauthorized(response)) return;

        const data = await response.json();

        if (response.ok) {
            emailInput.value = '';
            // Muat ulang daftar
            loadSharedUsers(noteId);
        } else {
            alert(data.error || 'Gagal menambahkan pengguna.');
        }
    } catch (err) {
        alert('Gagal menambahkan pengguna.');
    }
});

// ---- Hapus Akses User dari Catatan ----
async function removeSharedUser(noteId, userId) {
    if (!confirm('Hapus akses pengguna ini dari catatan?')) return;

    try {
        const response = await fetch(`${API_URL}/api/notes/${noteId}/share/${userId}`, {
            method: 'DELETE',
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        if (response.ok) {
            loadSharedUsers(noteId);
        } else {
            const data = await response.json();
            alert(data.error || 'Gagal menghapus akses.');
        }
    } catch (err) {
        alert('Gagal menghapus akses.');
    }
}

