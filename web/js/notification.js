// ===========================================
// notification.js — Logika panel notifikasi
// ===========================================

let notifPanelOpen = false;

function formatDateTime(dateStr) {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleDateString('id-ID', {
        day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit'
    });
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function loadNotifications() {
    try {
        const response = await fetch(`${API_URL}/api/notifications`, {
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        const data = await response.json();
        const notifications = data.notifications || [];
        const unreadCount = data.unread_count || 0;

        updateNotifBadge(unreadCount);
        renderNotifications(notifications);

    } catch (err) {
        console.error('Gagal memuat notifikasi:', err);
        const notifList = document.getElementById('notif-list');
        if (notifList) {
            notifList.innerHTML = '<p class="empty-state">Gagal memuat notifikasi.</p>';
        }
    }
}

function updateNotifBadge(count) {
    const badge = document.getElementById('notif-badge');
    if (!badge) return;

    if (count > 0) {
        badge.textContent = count > 99 ? '99+' : count;
        badge.style.display = 'flex';
    } else {
        badge.style.display = 'none';
    }
}

function renderNotifications(notifications) {
    const notifList = document.getElementById('notif-list');
    if (!notifList) return;

    notifList.innerHTML = '';

    if (notifications.length === 0) {
        notifList.innerHTML = '<p class="empty-state">Belum ada notifikasi.</p>';
        return;
    }

    notifications.forEach(notif => {
        const item = document.createElement('div');
        item.className = `notif-item${!notif.is_read ? ' unread' : ''}`;
        item.dataset.id = notif.id;

        let icon = '📝';
        if (notif.type === 'shared') icon = '📤';
        else if (notif.type === 'deadline') icon = '⏰';

        item.innerHTML = `
            <span class="notif-icon">${icon}</span>
            <div class="notif-content">
                <p class="notif-message">${escapeHtml(notif.message)}</p>
                <span class="notif-time">${formatDateTime(notif.created_at)}</span>
            </div>
            ${!notif.is_read ? '<span class="notif-dot"></span>' : ''}
        `;

        item.addEventListener('click', async () => {
            if (!notif.is_read) {
                await markNotificationRead(notif.id);
                item.classList.remove('unread');
                const dot = item.querySelector('.notif-dot');
                if (dot) dot.remove();
                const currentCount = parseInt(document.getElementById('notif-badge').textContent) || 0;
                updateNotifBadge(Math.max(0, currentCount - 1));
            }

            if (notif.related_note_id) {
                closeNotifPanel();
                openNoteDetail(notif.related_note_id, false);
            }
        });

        notifList.appendChild(item);
    });
}

async function markNotificationRead(notifId) {
    try {
        const response = await fetch(`${API_URL}/api/notifications/${notifId}/read`, {
            method: 'PUT',
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        if (!response.ok) {
            console.error('Gagal menandai notifikasi dibaca');
        }
    } catch (err) {
        console.error('Error marking notification read:', err);
    }
}

async function markAllNotificationsRead() {
    try {
        const response = await fetch(`${API_URL}/api/notifications/read-all`, {
            method: 'PUT',
            headers: authHeaders()
        });

        if (handleUnauthorized(response)) return;

        if (response.ok) {
            document.querySelectorAll('.notif-item.unread').forEach(item => {
                item.classList.remove('unread');
                const dot = item.querySelector('.notif-dot');
                if (dot) dot.remove();
            });
            updateNotifBadge(0);
        } else {
            const data = await response.json();
            alert(data.error || 'Gagal menandai semua dibaca.');
        }
    } catch (err) {
        console.error('Error marking all read:', err);
        alert('Gagal menandai semua dibaca.');
    }
}

function toggleNotifPanel() {
    const panel = document.getElementById('notif-panel');
    if (!panel) return;

    if (notifPanelOpen) {
        closeNotifPanel();
    } else {
        openNotifPanel();
    }
}

function openNotifPanel() {
    const panel = document.getElementById('notif-panel');
    if (!panel) return;

    panel.style.display = 'block';
    notifPanelOpen = true;

    loadNotifications();

    document.addEventListener('click', closeNotifPanelOutside);
}

function closeNotifPanel() {
    const panel = document.getElementById('notif-panel');
    if (!panel) return;

    panel.style.display = 'none';
    notifPanelOpen = false;

    document.removeEventListener('click', closeNotifPanelOutside);
}

function closeNotifPanelOutside(event) {
    const panel = document.getElementById('notif-panel');
    const btnNotif = document.getElementById('btn-notif');

    if (panel && btnNotif && !panel.contains(event.target) && !btnNotif.contains(event.target)) {
        closeNotifPanel();
    }
}

document.addEventListener('DOMContentLoaded', function() {
    const btnNotif = document.getElementById('btn-notif');
    if (btnNotif) {
        btnNotif.addEventListener('click', function(e) {
            e.stopPropagation();
            toggleNotifPanel();
        });
    }

    const btnReadAll = document.getElementById('btn-read-all');
    if (btnReadAll) {
        btnReadAll.addEventListener('click', function(e) {
            e.stopPropagation();
            markAllNotificationsRead();
        });
    }
});
