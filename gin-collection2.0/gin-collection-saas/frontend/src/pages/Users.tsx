import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import { userAPI } from '../api/services';
import type { User, UserRole } from '../types';
import {
  Users as UsersIcon,
  UserPlus,
  Search,
  MoreVertical,
  Mail,
  Shield,
  ShieldCheck,
  Eye,
  Crown,
  Loader2,
  Check,
  X,
  AlertTriangle,
  Trash2,
  UserX,
  UserCheck,
  Lock,
  RefreshCw
} from 'lucide-react';
import './Users.css';

const ROLE_CONFIG: Record<UserRole, { label: string; className: string; icon: React.ReactNode; description: string }> = {
  owner: {
    label: 'Inhaber',
    className: 'users-role-badge--owner',
    icon: <Crown className="w-3.5 h-3.5" />,
    description: 'Vollzugriff, kann Abonnement verwalten'
  },
  admin: {
    label: 'Admin',
    className: 'users-role-badge--admin',
    icon: <ShieldCheck className="w-3.5 h-3.5" />,
    description: 'Kann Benutzer und Einstellungen verwalten'
  },
  member: {
    label: 'Mitglied',
    className: 'users-role-badge--member',
    icon: <Shield className="w-3.5 h-3.5" />,
    description: 'Kann Gins erstellen und bearbeiten'
  },
  viewer: {
    label: 'Betrachter',
    className: 'users-role-badge--viewer',
    icon: <Eye className="w-3.5 h-3.5" />,
    description: 'Nur Lesezugriff'
  }
};

const Users = () => {
  const { user, tenant } = useAuthStore();
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState<UserRole | 'all'>('all');
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [activeDropdown, setActiveDropdown] = useState<number | null>(null);
  const [isUpdating, setIsUpdating] = useState<number | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  // Invite form state
  const [inviteForm, setInviteForm] = useState({
    email: '',
    role: 'member' as UserRole,
    first_name: '',
    last_name: ''
  });
  const [isInviting, setIsInviting] = useState(false);
  const [inviteError, setInviteError] = useState('');

  const currentTier = tenant?.tier || 'free';
  const isEnterprise = currentTier === 'enterprise';
  const canManageUsers = user?.role === 'owner' || user?.role === 'admin';

  // Load users from API
  const loadUsers = async () => {
    setIsLoading(true);
    setLoadError('');
    try {
      const response = await userAPI.list();
      const apiResponse = response.data as unknown as { success: boolean; data: { users: User[]; count: number } };
      if (apiResponse.success && apiResponse.data) {
        setUsers(apiResponse.data.users || []);
      }
    } catch (err) {
      console.error('Failed to load users:', err);
      setLoadError('Benutzer konnten nicht geladen werden.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (isEnterprise) {
      loadUsers();
    }
  }, [isEnterprise]);

  // Filter users
  const filteredUsers = users.filter(u => {
    const matchesSearch =
      u.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
      `${u.first_name} ${u.last_name}`.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesRole = roleFilter === 'all' || u.role === roleFilter;
    return matchesSearch && matchesRole;
  });

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    setInviteError('');

    if (!inviteForm.email) {
      setInviteError('E-Mail-Adresse ist erforderlich');
      return;
    }

    setIsInviting(true);
    try {
      const response = await userAPI.invite({
        email: inviteForm.email,
        first_name: inviteForm.first_name || undefined,
        last_name: inviteForm.last_name || undefined,
        role: inviteForm.role,
      });
      const apiResponse = response.data as unknown as { success: boolean; data: User };
      if (apiResponse.success && apiResponse.data) {
        setUsers([...users, apiResponse.data]);
      }
      setShowInviteModal(false);
      setInviteForm({ email: '', role: 'member', first_name: '', last_name: '' });
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setInviteError(error.response?.data?.error || 'Einladung fehlgeschlagen. Bitte versuche es erneut.');
    } finally {
      setIsInviting(false);
    }
  };

  const handleChangeRole = async (userId: number, newRole: UserRole) => {
    const targetUser = users.find(u => u.id === userId);
    if (!targetUser) return;

    setIsUpdating(userId);
    try {
      const response = await userAPI.update(userId, {
        email: targetUser.email,
        first_name: targetUser.first_name,
        last_name: targetUser.last_name,
        role: newRole,
        is_active: targetUser.is_active,
      });
      const apiResponse = response.data as unknown as { success: boolean; data: User };
      if (apiResponse.success && apiResponse.data) {
        setUsers(users.map(u => u.id === userId ? apiResponse.data : u));
      }
    } catch (err) {
      console.error('Failed to update user role:', err);
    } finally {
      setIsUpdating(null);
      setActiveDropdown(null);
    }
  };

  const handleToggleStatus = async (userId: number) => {
    const targetUser = users.find(u => u.id === userId);
    if (!targetUser) return;

    setIsUpdating(userId);
    try {
      const response = await userAPI.update(userId, {
        email: targetUser.email,
        first_name: targetUser.first_name,
        last_name: targetUser.last_name,
        role: targetUser.role,
        is_active: !targetUser.is_active,
      });
      const apiResponse = response.data as unknown as { success: boolean; data: User };
      if (apiResponse.success && apiResponse.data) {
        setUsers(users.map(u => u.id === userId ? apiResponse.data : u));
      }
    } catch (err) {
      console.error('Failed to toggle user status:', err);
    } finally {
      setIsUpdating(null);
      setActiveDropdown(null);
    }
  };

  const handleDeleteUser = async () => {
    if (!selectedUser) return;

    setIsDeleting(true);
    try {
      await userAPI.delete(selectedUser.id);
      setUsers(users.filter(u => u.id !== selectedUser.id));
      setShowDeleteModal(false);
      setSelectedUser(null);
    } catch (err) {
      console.error('Failed to delete user:', err);
    } finally {
      setIsDeleting(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('de-DE', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric'
    });
  };

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { staggerChildren: 0.1 }
    }
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { type: 'spring', stiffness: 100, damping: 15 }
    }
  };

  const modalVariants = {
    hidden: { opacity: 0, scale: 0.95 },
    visible: {
      opacity: 1,
      scale: 1,
      transition: { type: 'spring', stiffness: 300, damping: 25 }
    },
    exit: {
      opacity: 0,
      scale: 0.95,
      transition: { duration: 0.15 }
    }
  };

  // Show upgrade message for non-enterprise tiers
  if (!isEnterprise) {
    return (
      <div className="users-page">
        <div className="users-ambient">
          <div className="users-orb users-orb--gold" />
          <div className="users-orb users-orb--purple" />
          <div className="users-orb users-orb--green" />
        </div>

        <div className="users-content">
          <motion.div
            className="users-header"
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            <div className="users-header__info">
              <h1 className="users-header__title">
                <div className="users-header__icon">
                  <UsersIcon />
                </div>
                Team-Verwaltung
              </h1>
              <p className="users-header__subtitle">Verwalte dein Team und weise Rollen zu</p>
            </div>
          </motion.div>

          <motion.div
            className="users-card"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
          >
            <div className="users-upgrade">
              <div className="users-upgrade__icon">
                <Lock />
              </div>
              <h2 className="users-upgrade__title">Enterprise-Feature</h2>
              <p className="users-upgrade__text">
                Die Team-Verwaltung ist nur im Enterprise-Plan verfügbar.
                Upgrade jetzt, um unbegrenzt Team-Mitglieder einzuladen und Rollen zu verwalten.
              </p>
              <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                <Link to="/subscription" className="users-upgrade__btn">
                  <Crown />
                  Zum Enterprise-Plan upgraden
                </Link>
              </motion.div>

              <div className="users-upgrade__features">
                <p className="users-upgrade__features-title">Im Enterprise-Plan enthalten</p>
                <div className="users-upgrade__features-grid">
                  <div className="users-upgrade__feature">
                    <UsersIcon />
                    <p className="users-upgrade__feature-text">Unbegrenzte Mitglieder</p>
                  </div>
                  <div className="users-upgrade__feature">
                    <Shield />
                    <p className="users-upgrade__feature-text">Rollen & Berechtigungen</p>
                  </div>
                  <div className="users-upgrade__feature">
                    <Mail />
                    <p className="users-upgrade__feature-text">E-Mail-Einladungen</p>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    );
  }

  // Show loading state
  if (isLoading) {
    return (
      <div className="users-page">
        <div className="users-ambient">
          <div className="users-orb users-orb--gold" />
          <div className="users-orb users-orb--purple" />
          <div className="users-orb users-orb--green" />
        </div>

        <div className="users-content">
          <motion.div
            className="users-header"
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
          >
            <div className="users-header__info">
              <h1 className="users-header__title">
                <div className="users-header__icon">
                  <UsersIcon />
                </div>
                Team-Verwaltung
              </h1>
              <p className="users-header__subtitle">Lade Team-Mitglieder...</p>
            </div>
          </motion.div>

          <motion.div
            className="users-card"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
          >
            <div className="users-loading">
              <div className="users-loading__spinner" />
              <p className="users-loading__text">Lade Team-Mitglieder...</p>
            </div>
          </motion.div>
        </div>
      </div>
    );
  }

  // Show error state
  if (loadError) {
    return (
      <div className="users-page">
        <div className="users-ambient">
          <div className="users-orb users-orb--gold" />
          <div className="users-orb users-orb--purple" />
          <div className="users-orb users-orb--green" />
        </div>

        <div className="users-content">
          <motion.div
            className="users-header"
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
          >
            <div className="users-header__info">
              <h1 className="users-header__title">
                <div className="users-header__icon">
                  <UsersIcon />
                </div>
                Team-Verwaltung
              </h1>
              <p className="users-header__subtitle">Verwalte dein Team und weise Rollen zu</p>
            </div>
          </motion.div>

          <motion.div
            className="users-card"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
          >
            <div className="users-error">
              <AlertTriangle className="users-error__icon" />
              <p className="users-error__text">{loadError}</p>
              <motion.button
                onClick={loadUsers}
                className="users-error__btn"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <RefreshCw />
                Erneut versuchen
              </motion.button>
            </div>
          </motion.div>
        </div>
      </div>
    );
  }

  return (
    <div className="users-page">
      {/* Ambient Background */}
      <div className="users-ambient">
        <div className="users-orb users-orb--gold" />
        <div className="users-orb users-orb--purple" />
        <div className="users-orb users-orb--green" />
      </div>

      <div className="users-content">
        {/* Header */}
        <motion.div
          className="users-header"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <div className="users-header__info">
            <h1 className="users-header__title">
              <div className="users-header__icon">
                <UsersIcon />
              </div>
              Team-Verwaltung
            </h1>
            <p className="users-header__subtitle">
              <span className="users-header__count">{users.length}</span> Mitglied{users.length !== 1 ? 'er' : ''} in deinem Team
            </p>
          </div>

          <div className="users-header__actions">
            <motion.button
              onClick={loadUsers}
              className="users-refresh-btn"
              whileHover={{ rotate: 180 }}
              whileTap={{ scale: 0.95 }}
              title="Aktualisieren"
            >
              <RefreshCw />
            </motion.button>
            {canManageUsers && (
              <motion.button
                onClick={() => setShowInviteModal(true)}
                className="users-invite-btn"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <UserPlus />
                Mitglied einladen
              </motion.button>
            )}
          </div>
        </motion.div>

        {/* Filters */}
        <motion.div
          className="users-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <div className="users-filters">
            <div className="users-search">
              <input
                type="text"
                placeholder="Nach Name oder E-Mail suchen..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="users-search__input"
              />
              <div className="users-search__icon">
                <Search />
              </div>
            </div>

            <select
              value={roleFilter}
              onChange={(e) => setRoleFilter(e.target.value as UserRole | 'all')}
              className="users-filter-select"
            >
              <option value="all">Alle Rollen</option>
              <option value="owner">Inhaber</option>
              <option value="admin">Admin</option>
              <option value="member">Mitglied</option>
              <option value="viewer">Betrachter</option>
            </select>
          </div>
        </motion.div>

        {/* Users Table */}
        <motion.div
          className="users-table-container"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <table className="users-table">
            <thead>
              <tr>
                <th>Benutzer</th>
                <th>Rolle</th>
                <th>Status</th>
                <th>Beigetreten</th>
                <th>Aktionen</th>
              </tr>
            </thead>
            <tbody>
              {filteredUsers.length === 0 ? (
                <tr>
                  <td colSpan={5}>
                    <div className="users-empty">
                      <div className="users-empty__icon">
                        <UsersIcon />
                      </div>
                      <p className="users-empty__text">Keine Benutzer gefunden</p>
                    </div>
                  </td>
                </tr>
              ) : (
                filteredUsers.map((u) => {
                  const roleConfig = ROLE_CONFIG[u.role];
                  const isCurrentUser = u.id === user?.id;
                  const isOwner = u.role === 'owner';

                  return (
                    <motion.tr
                      key={u.id}
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      transition={{ duration: 0.3 }}
                    >
                      {/* User Info */}
                      <td>
                        <div className="users-user-info">
                          <div className="users-avatar">
                            <span className="users-avatar__initials">
                              {(u.first_name?.[0] || u.email[0]).toUpperCase()}
                              {(u.last_name?.[0] || '').toUpperCase()}
                            </span>
                          </div>
                          <div className="users-user-details">
                            <p className="users-user-name">
                              {u.first_name && u.last_name
                                ? `${u.first_name} ${u.last_name}`
                                : u.email.split('@')[0]
                              }
                              {isCurrentUser && (
                                <span className="users-user-badge">Du</span>
                              )}
                            </p>
                            <p className="users-user-email">{u.email}</p>
                          </div>
                        </div>
                      </td>

                      {/* Role */}
                      <td>
                        <span className={`users-role-badge ${roleConfig.className}`}>
                          {roleConfig.icon}
                          {roleConfig.label}
                        </span>
                      </td>

                      {/* Status */}
                      <td>
                        <span className={`users-status ${u.is_active ? 'users-status--active' : 'users-status--inactive'}`}>
                          <span className="users-status__dot"></span>
                          {u.is_active ? 'Aktiv' : 'Inaktiv'}
                        </span>
                      </td>

                      {/* Joined Date */}
                      <td>
                        <span className="users-date">{formatDate(u.created_at)}</span>
                      </td>

                      {/* Actions */}
                      <td>
                        <div className="users-actions">
                          {isUpdating === u.id ? (
                            <Loader2 className="users-spinner" style={{ width: 20, height: 20, color: 'var(--gold)' }} />
                          ) : canManageUsers && !isOwner && !isCurrentUser ? (
                            <>
                              <button
                                onClick={() => setActiveDropdown(activeDropdown === u.id ? null : u.id)}
                                className="users-action-btn"
                              >
                                <MoreVertical />
                              </button>

                              <AnimatePresence>
                                {activeDropdown === u.id && (
                                  <>
                                    <div
                                      className="users-dropdown-backdrop"
                                      onClick={() => setActiveDropdown(null)}
                                    />
                                    <motion.div
                                      className="users-dropdown"
                                      initial={{ opacity: 0, y: -10 }}
                                      animate={{ opacity: 1, y: 0 }}
                                      exit={{ opacity: 0, y: -10 }}
                                      transition={{ duration: 0.15 }}
                                    >
                                      <div className="users-dropdown__section">
                                        <p className="users-dropdown__label">Rolle ändern</p>
                                        {(['admin', 'member', 'viewer'] as UserRole[]).map((role) => (
                                          <button
                                            key={role}
                                            onClick={() => handleChangeRole(u.id, role)}
                                            className={`users-dropdown__item ${u.role === role ? 'users-dropdown__item--active' : ''}`}
                                          >
                                            {ROLE_CONFIG[role].icon}
                                            {ROLE_CONFIG[role].label}
                                            {u.role === role && <Check className="users-dropdown__check" />}
                                          </button>
                                        ))}
                                      </div>

                                      <div className="users-dropdown__section">
                                        <button
                                          onClick={() => handleToggleStatus(u.id)}
                                          className="users-dropdown__item"
                                        >
                                          {u.is_active ? <UserX /> : <UserCheck />}
                                          {u.is_active ? 'Deaktivieren' : 'Aktivieren'}
                                        </button>

                                        <button
                                          onClick={() => {
                                            setSelectedUser(u);
                                            setShowDeleteModal(true);
                                            setActiveDropdown(null);
                                          }}
                                          className="users-dropdown__item users-dropdown__item--danger"
                                        >
                                          <Trash2 />
                                          Entfernen
                                        </button>
                                      </div>
                                    </motion.div>
                                  </>
                                )}
                              </AnimatePresence>
                            </>
                          ) : (
                            <span className="users-action-text">
                              {isOwner ? 'Inhaber' : isCurrentUser ? 'Du selbst' : '—'}
                            </span>
                          )}
                        </div>
                      </td>
                    </motion.tr>
                  );
                })
              )}
            </tbody>
          </table>
        </motion.div>

        {/* Role Descriptions */}
        <motion.div
          className="users-roles-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <h2 className="users-roles-card__title">Rollen-Übersicht</h2>
          <div className="users-roles-grid">
            {(Object.entries(ROLE_CONFIG) as [UserRole, typeof ROLE_CONFIG[UserRole]][]).map(([role, config]) => (
              <div key={role} className="users-role-item">
                <div className="users-role-item__badge">
                  <span className={`users-role-badge ${config.className}`}>
                    {config.icon}
                    {config.label}
                  </span>
                </div>
                <p className="users-role-item__description">{config.description}</p>
              </div>
            ))}
          </div>
        </motion.div>
      </div>

      {/* Invite Modal */}
      <AnimatePresence>
        {showInviteModal && (
          <motion.div
            className="users-modal-backdrop"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <motion.div
              className="users-modal"
              variants={modalVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
            >
              <div className="users-modal__header">
                <h3 className="users-modal__title">
                  <div className="users-modal__title-icon">
                    <UserPlus />
                  </div>
                  Team-Mitglied einladen
                </h3>
                <button
                  onClick={() => setShowInviteModal(false)}
                  className="users-modal__close"
                >
                  <X />
                </button>
              </div>

              <form onSubmit={handleInvite}>
                <div className="users-modal__body">
                  <div className="users-modal-form">
                    {inviteError && (
                      <div className="users-modal-error">
                        <AlertTriangle />
                        {inviteError}
                      </div>
                    )}

                    <div className="users-modal-group">
                      <label className="users-modal-label">
                        E-Mail-Adresse <span className="users-modal-label__required">*</span>
                      </label>
                      <input
                        type="email"
                        value={inviteForm.email}
                        onChange={(e) => setInviteForm({ ...inviteForm, email: e.target.value })}
                        placeholder="name@example.com"
                        className="users-modal-input"
                        required
                      />
                    </div>

                    <div className="users-modal-grid">
                      <div className="users-modal-group">
                        <label className="users-modal-label">Vorname</label>
                        <input
                          type="text"
                          value={inviteForm.first_name}
                          onChange={(e) => setInviteForm({ ...inviteForm, first_name: e.target.value })}
                          placeholder="Max"
                          className="users-modal-input"
                        />
                      </div>
                      <div className="users-modal-group">
                        <label className="users-modal-label">Nachname</label>
                        <input
                          type="text"
                          value={inviteForm.last_name}
                          onChange={(e) => setInviteForm({ ...inviteForm, last_name: e.target.value })}
                          placeholder="Mustermann"
                          className="users-modal-input"
                        />
                      </div>
                    </div>

                    <div className="users-modal-group">
                      <label className="users-modal-label">Rolle</label>
                      <select
                        value={inviteForm.role}
                        onChange={(e) => setInviteForm({ ...inviteForm, role: e.target.value as UserRole })}
                        className="users-modal-select"
                      >
                        <option value="admin">Admin - Kann Benutzer und Einstellungen verwalten</option>
                        <option value="member">Mitglied - Kann Gins erstellen und bearbeiten</option>
                        <option value="viewer">Betrachter - Nur Lesezugriff</option>
                      </select>
                    </div>

                    <div className="users-modal-info">
                      <p className="users-modal-info__text">
                        Eine Einladungs-E-Mail wird an diese Adresse gesendet. Der Benutzer muss die Einladung akzeptieren, um Zugang zu erhalten.
                      </p>
                    </div>
                  </div>
                </div>

                <div className="users-modal__footer">
                  <button
                    type="button"
                    onClick={() => setShowInviteModal(false)}
                    className="users-modal-btn users-modal-btn--secondary"
                    disabled={isInviting}
                  >
                    Abbrechen
                  </button>
                  <button
                    type="submit"
                    className="users-modal-btn users-modal-btn--primary"
                    disabled={isInviting}
                  >
                    {isInviting ? (
                      <>
                        <Loader2 className="users-spinner" />
                        Wird gesendet...
                      </>
                    ) : (
                      <>
                        <Mail />
                        Einladung senden
                      </>
                    )}
                  </button>
                </div>
              </form>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Delete Confirmation Modal */}
      <AnimatePresence>
        {showDeleteModal && selectedUser && (
          <motion.div
            className="users-modal-backdrop"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <motion.div
              className="users-modal"
              variants={modalVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
            >
              <div className="users-modal__header">
                <h3 className="users-modal__title">
                  <div className="users-modal__title-icon users-modal__title-icon--danger">
                    <AlertTriangle />
                  </div>
                  Benutzer entfernen
                </h3>
                <button
                  onClick={() => {
                    setShowDeleteModal(false);
                    setSelectedUser(null);
                  }}
                  className="users-modal__close"
                >
                  <X />
                </button>
              </div>

              <div className="users-modal__body">
                <p className="users-modal__description">
                  Bist du sicher, dass du <strong>{selectedUser.first_name} {selectedUser.last_name}</strong> ({selectedUser.email}) aus dem Team entfernen möchtest?
                </p>
                <p className="users-modal__warning">
                  Der Benutzer verliert sofort den Zugang zu allen Daten. Diese Aktion kann nicht rückgängig gemacht werden.
                </p>
              </div>

              <div className="users-modal__footer">
                <button
                  onClick={() => {
                    setShowDeleteModal(false);
                    setSelectedUser(null);
                  }}
                  className="users-modal-btn users-modal-btn--secondary"
                  disabled={isDeleting}
                >
                  Abbrechen
                </button>
                <button
                  onClick={handleDeleteUser}
                  className="users-modal-btn users-modal-btn--danger"
                  disabled={isDeleting}
                >
                  {isDeleting ? (
                    <>
                      <Loader2 className="users-spinner" />
                      Wird entfernt...
                    </>
                  ) : (
                    <>
                      <Trash2 />
                      Entfernen
                    </>
                  )}
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default Users;
