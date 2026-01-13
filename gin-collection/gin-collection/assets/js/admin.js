// Admin Dashboard JavaScript

let allUsers = [];
let currentSort = { column: 'created_at', direction: 'desc' };

// Load dashboard data
async function loadDashboard() {
    await loadDashboardStats();
    await loadUsers();
}

// Load dashboard statistics
async function loadDashboardStats() {
    try {
        const response = await fetch('api/?action=admin-dashboard-stats');
        const data = await response.json();

        if (data.stats) {
            document.getElementById('stat-users').textContent = data.stats.total_users;
            document.getElementById('stat-gins').textContent = data.stats.total_gins;
        }
    } catch (error) {
        console.error('Error loading dashboard stats:', error);
    }
}

// Load all users
async function loadUsers() {
    try {
        const response = await fetch('api/?action=admin-users');
        const data = await response.json();

        if (data.users) {
            allUsers = data.users;
            displayUsers(allUsers);
        }
    } catch (error) {
        console.error('Error loading users:', error);
        document.getElementById('users-tbody').innerHTML = `
            <tr><td colspan="7" class="empty-state">Fehler beim Laden der Benutzer</td></tr>
        `;
    }
}

// Display users in table
function displayUsers(users) {
    const tbody = document.getElementById('users-tbody');

    if (users.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="empty-state">Keine Benutzer gefunden</td></tr>';
        return;
    }

    tbody.innerHTML = users.map(user => `
        <tr>
            <td>${user.id}</td>
            <td>
                ${user.username}
                ${user.is_admin ? '<span class="admin-badge-small">👑</span>' : ''}
            </td>
            <td>${user.email}</td>
            <td>${user.full_name || '-'}</td>
            <td>${user.gin_count}</td>
            <td>${formatDate(user.created_at)}</td>
            <td class="user-actions">
                <button class="btn-icon" onclick="showUserDetails(${user.id})" title="Details">👁️</button>
                <button class="btn-icon" onclick="editUser(${user.id})" title="Bearbeiten">✏️</button>
                ${!user.is_admin ? `<button class="btn-icon btn-danger" onclick="deleteUser(${user.id}, '${user.username}')" title="Löschen">🗑️</button>` : ''}
            </td>
        </tr>
    `).join('');
}

// Show user details
async function showUserDetails(userId) {
    try {
        const response = await fetch(`api/?action=admin-user-stats&user_id=${userId}`);
        const data = await response.json();

        if (data.user) {
            const user = data.user;
            document.getElementById('modal-title').textContent = `Benutzer: ${user.username}`;
            document.getElementById('modal-body').innerHTML = `
                <div class="detail-grid">
                    <div class="detail-row">
                        <span class="detail-label">ID:</span>
                        <span class="detail-value">${user.id}</span>
                    </div>
                    <div class="detail-row">
                        <span class="detail-label">Username:</span>
                        <span class="detail-value">${user.username}</span>
                    </div>
                    <div class="detail-row">
                        <span class="detail-label">Email:</span>
                        <span class="detail-value">${user.email}</span>
                    </div>
                    <div class="detail-row">
                        <span class="detail-label">Name:</span>
                        <span class="detail-value">${user.full_name || '-'}</span>
                    </div>
                    <div class="detail-row">
                        <span class="detail-label">Anzahl Gins:</span>
                        <span class="detail-value">${user.gin_count}</span>
                    </div>
                    <div class="detail-row">
                        <span class="detail-label">Registriert:</span>
                        <span class="detail-value">${formatDate(user.created_at)}</span>
                    </div>
                </div>

                ${user.latest_gins && user.latest_gins.length > 0 ? `
                    <div class="detail-section">
                        <h4>Neueste Gins:</h4>
                        <ul class="gin-list">
                            ${user.latest_gins.map(gin => `
                                <li>${gin.name} ${gin.brand ? `(${gin.brand})` : ''}</li>
                            `).join('')}
                        </ul>
                    </div>
                ` : ''}
            `;

            document.getElementById('modal-actions').innerHTML = `
                <button class="btn btn-secondary" onclick="closeUserModal()">Schließen</button>
            `;

            document.getElementById('user-modal').classList.add('active');
        }
    } catch (error) {
        console.error('Error loading user details:', error);
        alert('Fehler beim Laden der Benutzer-Details');
    }
}

// Edit user
function editUser(userId) {
    const user = allUsers.find(u => u.id === userId);
    if (!user) return;

    document.getElementById('edit-user-id').value = user.id;
    document.getElementById('edit-email').value = user.email;
    document.getElementById('edit-full-name').value = user.full_name || '';

    document.getElementById('edit-modal').classList.add('active');
}

// Save user changes
async function saveUser() {
    const userId = document.getElementById('edit-user-id').value;
    const email = document.getElementById('edit-email').value;
    const fullName = document.getElementById('edit-full-name').value;

    try {
        const response = await fetch('api/?action=admin-update-user', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: parseInt(userId),
                email: email,
                full_name: fullName
            })
        });

        const data = await response.json();

        if (data.success) {
            alert('Benutzer erfolgreich aktualisiert');
            closeEditModal();
            loadUsers();
        } else {
            alert('Fehler: ' + (data.error || 'Unbekannter Fehler'));
        }
    } catch (error) {
        console.error('Error updating user:', error);
        alert('Fehler beim Aktualisieren des Benutzers');
    }
}

// Delete user
async function deleteUser(userId, username) {
    if (!confirm(`Möchtest du den Benutzer "${username}" wirklich löschen?\n\nAlle Gins dieses Benutzers werden ebenfalls gelöscht!`)) {
        return;
    }

    try {
        const response = await fetch('api/?action=admin-delete-user', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: userId })
        });

        const data = await response.json();

        if (data.success) {
            alert('Benutzer erfolgreich gelöscht');
            loadUsers();
            loadDashboardStats();
        } else {
            alert('Fehler: ' + (data.error || 'Unbekannter Fehler'));
        }
    } catch (error) {
        console.error('Error deleting user:', error);
        alert('Fehler beim Löschen des Benutzers');
    }
}

// Sort users
function sortUsers(column) {
    if (currentSort.column === column) {
        currentSort.direction = currentSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentSort.column = column;
        currentSort.direction = 'asc';
    }

    const sorted = [...allUsers].sort((a, b) => {
        let aVal = a[column];
        let bVal = b[column];

        // Handle null values
        if (aVal === null) aVal = '';
        if (bVal === null) bVal = '';

        // Convert to numbers if needed
        if (column === 'id' || column === 'gin_count') {
            aVal = parseInt(aVal) || 0;
            bVal = parseInt(bVal) || 0;
        }

        if (currentSort.direction === 'asc') {
            return aVal > bVal ? 1 : -1;
        } else {
            return aVal < bVal ? 1 : -1;
        }
    });

    displayUsers(sorted);
}

// Search users
document.getElementById('search-users')?.addEventListener('input', (e) => {
    const query = e.target.value.toLowerCase();

    if (!query) {
        displayUsers(allUsers);
        return;
    }

    const filtered = allUsers.filter(user =>
        user.username.toLowerCase().includes(query) ||
        user.email.toLowerCase().includes(query) ||
        (user.full_name && user.full_name.toLowerCase().includes(query))
    );

    displayUsers(filtered);
});

// Close modals
function closeUserModal() {
    document.getElementById('user-modal').classList.remove('active');
}

function closeEditModal() {
    document.getElementById('edit-modal').classList.remove('active');
}

// Helper: Format date
function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('de-DE', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}
