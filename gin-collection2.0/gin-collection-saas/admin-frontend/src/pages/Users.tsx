import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { adminApi } from '../api';
import {
  Users as UsersIcon,
  Search,
  ChevronLeft,
  ChevronRight,
  Loader2,
  AlertTriangle,
  RefreshCw,
  Calendar,
  Clock,
  Crown,
  Shield,
  User,
  Eye,
  Hash
} from 'lucide-react';
import './Users.css';

interface User {
  id: number;
  tenant_id: number;
  email: string;
  first_name?: string;
  last_name?: string;
  role: string;
  is_active: boolean;
  created_at: string;
  last_login_at?: string;
}

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadUsers();
  }, [page]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      const response = await adminApi.getUsers(page, 20);
      setUsers(response.data.users || []);
      setTotal(response.data.total);
    } catch (err) {
      setError('Benutzer konnten nicht geladen werden');
    } finally {
      setLoading(false);
    }
  };

  const getRoleIcon = (role: string) => {
    switch (role) {
      case 'owner': return Crown;
      case 'admin': return Shield;
      case 'member': return User;
      case 'viewer': return Eye;
      default: return User;
    }
  };

  const getRoleBadgeClass = (role: string) => {
    switch (role) {
      case 'owner': return 'admin-users-badge--owner';
      case 'admin': return 'admin-users-badge--admin';
      case 'member': return 'admin-users-badge--member';
      case 'viewer': return 'admin-users-badge--viewer';
      default: return 'admin-users-badge--member';
    }
  };

  const getRoleLabel = (role: string) => {
    switch (role) {
      case 'owner': return 'Owner';
      case 'admin': return 'Admin';
      case 'member': return 'Mitglied';
      case 'viewer': return 'Viewer';
      default: return role;
    }
  };

  const filteredUsers = users.filter((user) =>
    user.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
    (user.first_name?.toLowerCase().includes(searchQuery.toLowerCase())) ||
    (user.last_name?.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { staggerChildren: 0.05 }
    }
  };

  const rowVariants = {
    hidden: { opacity: 0, x: -20 },
    visible: {
      opacity: 1,
      x: 0,
      transition: { type: 'spring', stiffness: 100, damping: 15 }
    }
  };

  if (loading && users.length === 0) {
    return (
      <div className="admin-users">
        <div className="admin-users-loader">
          <Loader2 className="admin-users-loader__icon" />
          <span>Lade Benutzer...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-users">
      {/* Header */}
      <motion.div
        className="admin-users-header"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <div className="admin-users-header__content">
          <h1 className="admin-users-title">
            <UsersIcon />
            Benutzer
          </h1>
          <span className="admin-users-count">{total} Accounts</span>
        </div>

        <div className="admin-users-actions">
          <div className="admin-users-search">
            <Search />
            <input
              type="text"
              placeholder="Benutzer suchen..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="admin-users-search__input"
            />
          </div>
          <motion.button
            className="admin-users-refresh"
            onClick={loadUsers}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <RefreshCw />
          </motion.button>
        </div>
      </motion.div>

      {error && (
        <motion.div
          className="admin-users-error"
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <AlertTriangle />
          <span>{error}</span>
        </motion.div>
      )}

      {/* Table */}
      <motion.div
        className="admin-users-table-wrapper"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <table className="admin-users-table">
          <thead>
            <tr>
              <th>Benutzer</th>
              <th>Tenant</th>
              <th>Rolle</th>
              <th>Status</th>
              <th>Letzter Login</th>
              <th>Erstellt</th>
            </tr>
          </thead>
          <motion.tbody variants={containerVariants} initial="hidden" animate="visible">
            {filteredUsers.map((user) => {
              const RoleIcon = getRoleIcon(user.role);
              return (
                <motion.tr key={user.id} variants={rowVariants}>
                  <td>
                    <div className="admin-users-user">
                      <div className="admin-users-user__avatar">
                        {user.first_name?.[0] || user.email[0].toUpperCase()}
                      </div>
                      <div className="admin-users-user__info">
                        <span className="admin-users-user__email">{user.email}</span>
                        {(user.first_name || user.last_name) && (
                          <span className="admin-users-user__name">
                            {user.first_name} {user.last_name}
                          </span>
                        )}
                      </div>
                    </div>
                  </td>
                  <td>
                    <div className="admin-users-tenant">
                      <Hash />
                      {user.tenant_id}
                    </div>
                  </td>
                  <td>
                    <span className={`admin-users-badge ${getRoleBadgeClass(user.role)}`}>
                      <RoleIcon />
                      {getRoleLabel(user.role)}
                    </span>
                  </td>
                  <td>
                    <span className={`admin-users-status ${user.is_active ? 'admin-users-status--active' : 'admin-users-status--inactive'}`}>
                      <span className="admin-users-status__dot" />
                      {user.is_active ? 'Aktiv' : 'Inaktiv'}
                    </span>
                  </td>
                  <td>
                    <div className="admin-users-date">
                      <Clock />
                      {user.last_login_at
                        ? new Date(user.last_login_at).toLocaleDateString('de-DE')
                        : 'Nie'}
                    </div>
                  </td>
                  <td>
                    <div className="admin-users-date">
                      <Calendar />
                      {new Date(user.created_at).toLocaleDateString('de-DE')}
                    </div>
                  </td>
                </motion.tr>
              );
            })}
          </motion.tbody>
        </table>

        {filteredUsers.length === 0 && (
          <div className="admin-users-empty">
            <UsersIcon />
            <span>Keine Benutzer gefunden</span>
          </div>
        )}
      </motion.div>

      {/* Pagination */}
      {total > 20 && (
        <motion.div
          className="admin-users-pagination"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          <button
            className="admin-users-pagination__btn"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
          >
            <ChevronLeft />
            Zurück
          </button>
          <span className="admin-users-pagination__info">
            Seite {page} von {Math.ceil(total / 20)}
          </span>
          <button
            className="admin-users-pagination__btn"
            onClick={() => setPage((p) => p + 1)}
            disabled={page * 20 >= total}
          >
            Weiter
            <ChevronRight />
          </button>
        </motion.div>
      )}
    </div>
  );
}
