// ===========================================
// notes.js — Logika CRUD catatan + sidebar
// ===========================================

// ---- Muat Catatan Milik Sendiri ----
async function loadMyNotes() {
    const notesList = document.getElementById('notes-list');
    const emptyState = document.getElementById('empty-state');

    try {
        const response = await fetch(`${API_URL}/api/notes`, {
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        const data = await response.json();
        const notes = data.notes || [];

        // Kosongkan daftar
        notesList.innerHTML = '';

        if (notes.length === 0) {
            notesList.innerHTML = '<p class="empty-state">Belum ada catatan. Klik "+ Catatan Baru" untuk mulai.</p>';
            return;
        }

        // Render setiap catatan sebagai card
        notes.forEach(note => {
            const card = createNoteCard(note, false);
            notesList.appendChild(card);
        });

    } catch (err) {
        notesList.innerHTML = '<p class="empty-state">Gagal memuat catatan. Coba refresh halaman.</p>';
    }
}

// ---- Muat Catatan yang Dibagikan ke Saya ----
async function loadSharedNotes() {
    const notesList = document.getElementById('notes-list');

    try {
        const response = await fetch(`${API_URL}/api/notes/shared`, {
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        const data = await response.json();
        const notes = data.notes || [];

        notesList.innerHTML = '';

        if (notes.length === 0) {
            notesList.innerHTML = '<p class="empty-state">Belum ada catatan yang dibagikan ke Anda.</p>';
            return;
        }

        notes.forEach(note => {
            const card = createNoteCard(note, true);
            notesList.appendChild(card);
        });

    } catch (err) {
        notesList.innerHTML = '<p class="empty-state">Gagal memuat catatan.</p>';
    }
}

// ---- Buat Element Card Catatan ----
function createNoteCard(note, isShared) {
    const card = document.createElement('div');
    card.className = `note-card${isShared ? ' shared' : ''}`;
    card.onclick = () => openNoteDetail(note.id, isShared);

    // Preview isi (maks 100 karakter)
    const contentPreview = note.content
        ? (note.content.length > 100 ? note.content.substring(0, 100) + '...' : note.content)
        : 'Tidak ada isi.';

    // Info deadline
    let deadlineHtml = '';
    if (note.deadline) {
        const overdueClass = isOverdue(note.deadline) ? ' overdue' : '';
        deadlineHtml = `<span class="note-card-deadline${overdueClass}">⏰ ${formatDateShort(note.deadline)}</span>`;
    }

    // Info pemilik (hanya untuk catatan yang dibagikan)
    let ownerHtml = '';
    if (isShared && note.user) {
        ownerHtml = `<span class="note-card-owner">oleh ${note.user.name}</span>`;
    }

    card.innerHTML = `
        <div class="note-card-title">${escapeHtml(note.title)}</div>
        <div class="note-card-content">${escapeHtml(contentPreview)}</div>
        <div class="note-card-meta">
            <span>${formatDate(note.created_at)}</span>
            ${deadlineHtml}
            ${ownerHtml}
        </div>
    `;

    return card;
}

// ---- Buka Detail Catatan di Sidebar ----
async function openNoteDetail(noteId, isShared) {
    try {
        const response = await fetch(`${API_URL}/api/notes/${noteId}`, {
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        if (!response.ok) {
            alert('Gagal memuat detail catatan.');
            return;
        }

        const data = await response.json();
        const note = data.note;
        const isOwner = data.is_owner;

        if (isOwner) {
            // ---- Mode Edit (pemilik) ----
            document.getElementById('sidebar-title').textContent = 'Edit Catatan';
            document.getElementById('note-form').style.display = 'block';
            document.getElementById('readonly-view').style.display = 'none';
            document.getElementById('sidebar-actions-owner').style.display = 'flex';

            // Isi form
            document.getElementById('note-id').value = note.id;
            document.getElementById('note-title').value = note.title;
            document.getElementById('note-content').value = note.content || '';

            // Isi deadline
            if (note.deadline) {
                // Format ke datetime-local: "2006-01-02T15:04"
                const d = new Date(note.deadline);
                const local = new Date(d.getTime() - d.getTimezoneOffset() * 60000);
                document.getElementById('note-deadline').value = local.toISOString().slice(0, 16);
            } else {
                document.getElementById('note-deadline').value = '';
            }

            // Tampilkan tombol hapus dan bagian berbagi
            document.getElementById('btn-delete-note').style.display = 'block';
            document.getElementById('share-section').style.display = 'block';

            // Muat daftar user yang dibagikan
            loadSharedUsers(note.id);
        } else {
            // ---- Mode Read-Only (catatan yang dibagikan) ----
            document.getElementById('sidebar-title').textContent = 'Detail Catatan';
            document.getElementById('note-form').style.display = 'none';
            document.getElementById('readonly-view').style.display = 'block';
            document.getElementById('sidebar-actions-owner').style.display = 'none';

            document.getElementById('readonly-title').textContent = note.title;
            document.getElementById('readonly-content').textContent = note.content || 'Tidak ada isi.';
            document.getElementById('readonly-deadline').textContent = note.deadline ? formatDate(note.deadline) : 'Tidak ada';
            document.getElementById('readonly-owner').textContent = note.user ? note.user.name : '-';
        }

        openSidebar();

    } catch (err) {
        alert('Gagal memuat catatan.');
    }
}

// ---- Simpan Catatan (Buat Baru / Update) ----
document.getElementById('note-form').addEventListener('submit', async function(e) {
    e.preventDefault();

    const noteId = document.getElementById('note-id').value;
    const title = document.getElementById('note-title').value;
    const content = document.getElementById('note-content').value;
    const deadline = document.getElementById('note-deadline').value;

    const body = { title, content };

    // Kirim deadline dalam format ISO 8601
    if (deadline) {
        body.deadline = new Date(deadline).toISOString();
    } else {
        body.deadline = '';
    }

    const btnSave = document.getElementById('btn-save-note');
    btnSave.disabled = true;
    btnSave.textContent = 'Menyimpan...';

    try {
        let url, method;
        if (noteId) {
            // Update catatan yang sudah ada
            url = `${API_URL}/api/notes/${noteId}`;
            method = 'PUT';
        } else {
            // Buat catatan baru
            url = `${API_URL}/api/notes`;
            method = 'POST';
        }

        const response = await fetch(url, {
            method: method,
            headers: authHeaders(),
            body: JSON.stringify(body)
        });

        if (handleUnauthorized(response)) return;

        if (response.ok) {
            closeSidebar();
            // Muat ulang daftar catatan
            if (currentTab === 'my-notes') {
                loadMyNotes();
            }
        } else {
            const data = await response.json();
            alert(data.error || 'Gagal menyimpan catatan.');
        }
    } catch (err) {
        alert('Gagal menyimpan catatan.');
    } finally {
        btnSave.disabled = false;
        btnSave.textContent = 'Simpan';
    }
});

// ---- Hapus Catatan ----
document.getElementById('btn-delete-note').addEventListener('click', async function() {
    const noteId = document.getElementById('note-id').value;
    if (!noteId) return;

    if (!confirm('Apakah Anda yakin ingin menghapus catatan ini?')) return;

    try {
        const response = await fetch(`${API_URL}/api/notes/${noteId}`, {
            method: 'DELETE',
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        if (response.ok) {
            closeSidebar();
            loadMyNotes();
        } else {
            const data = await response.json();
            alert(data.error || 'Gagal menghapus catatan.');
        }
    } catch (err) {
        alert('Gagal menghapus catatan.');
    }
});

// ---- Escape HTML (keamanan: mencegah XSS) ----
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

