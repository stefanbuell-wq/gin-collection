import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import { authAPI } from '../api/services';
import {
  User,
  Building2,
  Lock,
  CreditCard,
  Bell,
  Shield,
  Save,
  Check,
  AlertTriangle,
  Crown,
  Mail,
  Eye,
  EyeOff,
  AlertCircle,
  Settings as SettingsIcon,
  Globe,
  Calendar,
  Sparkles
} from 'lucide-react';
import './Settings.css';

type SettingsTab = 'profile' | 'account' | 'security' | 'notifications';

const Settings = () => {
  const { user, tenant } = useAuthStore();
  const [activeTab, setActiveTab] = useState<SettingsTab>('profile');

  // Profile form state
  const [profileForm, setProfileForm] = useState({
    first_name: user?.first_name || '',
    last_name: user?.last_name || '',
    email: user?.email || ''
  });
  const [profileSaving, setProfileSaving] = useState(false);
  const [profileSaved, setProfileSaved] = useState(false);
  const [profileError, setProfileError] = useState('');

  // Password form state
  const [passwordForm, setPasswordForm] = useState({
    current_password: '',
    new_password: '',
    confirm_password: ''
  });
  const [passwordSaving, setPasswordSaving] = useState(false);
  const [passwordSaved, setPasswordSaved] = useState(false);
  const [passwordError, setPasswordError] = useState('');
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);

  const handleProfileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setProfileForm(prev => ({ ...prev, [name]: value }));
    setProfileSaved(false);
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordForm(prev => ({ ...prev, [name]: value }));
    setPasswordSaved(false);
    setPasswordError('');
  };

  const handleProfileSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setProfileError('');
    setProfileSaving(true);

    try {
      const response = await authAPI.updateProfile({
        first_name: profileForm.first_name || undefined,
        last_name: profileForm.last_name || undefined,
      });
      const apiResponse = response.data as unknown as { success: boolean; data: typeof user };
      if (apiResponse.data) {
        useAuthStore.getState().setUser(apiResponse.data);
      }
      setProfileSaved(true);
      setTimeout(() => setProfileSaved(false), 3000);
    } catch {
      setProfileError('Profil konnte nicht aktualisiert werden. Bitte versuche es erneut.');
    } finally {
      setProfileSaving(false);
    }
  };

  const handlePasswordSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordError('');

    if (passwordForm.new_password !== passwordForm.confirm_password) {
      setPasswordError('Passworter stimmen nicht uberein');
      return;
    }

    if (passwordForm.new_password.length < 8) {
      setPasswordError('Passwort muss mindestens 8 Zeichen haben');
      return;
    }

    setPasswordSaving(true);

    try {
      await authAPI.changePassword(passwordForm.current_password, passwordForm.new_password);
      setPasswordForm({ current_password: '', new_password: '', confirm_password: '' });
      setPasswordSaved(true);
      setTimeout(() => setPasswordSaved(false), 3000);
    } catch {
      setPasswordError('Passwort konnte nicht geandert werden. Bitte prufe dein aktuelles Passwort.');
    } finally {
      setPasswordSaving(false);
    }
  };

  const tabs = [
    { id: 'profile' as SettingsTab, label: 'Profil', icon: User },
    { id: 'account' as SettingsTab, label: 'Konto', icon: Building2 },
    { id: 'security' as SettingsTab, label: 'Sicherheit', icon: Lock },
    { id: 'notifications' as SettingsTab, label: 'Benachrichtigungen', icon: Bell }
  ];

  const getTierBadgeClass = (tier: string) => {
    switch (tier) {
      case 'free': return 'settings-badge--tier-free';
      case 'basic': return 'settings-badge--tier-basic';
      case 'pro': return 'settings-badge--tier-pro';
      case 'enterprise': return 'settings-badge--tier-enterprise';
      default: return 'settings-badge--tier-free';
    }
  };

  const getRoleBadgeClass = (role: string) => {
    switch (role) {
      case 'owner': return 'settings-badge--role-owner';
      case 'admin': return 'settings-badge--role-admin';
      case 'member': return 'settings-badge--role-member';
      case 'viewer': return 'settings-badge--role-viewer';
      default: return 'settings-badge--role-viewer';
    }
  };

  const getUserInitials = () => {
    const first = user?.first_name?.[0] || user?.email?.[0] || '?';
    const last = user?.last_name?.[0] || '';
    return (first + last).toUpperCase();
  };

  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 }
  };

  return (
    <div className="settings-page">
      {/* Ambient Background */}
      <div className="settings-page__ambient">
        <div className="settings-ambient-orb settings-ambient-orb--1" />
        <div className="settings-ambient-orb settings-ambient-orb--2" />
        <div className="settings-ambient-orb settings-ambient-orb--3" />
      </div>

      <div className="settings-container">
        {/* Header */}
        <motion.div
          className="settings-header"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <div className="settings-header__content">
            <div className="settings-header__icon">
              <SettingsIcon size={28} />
            </div>
            <div className="settings-header__text">
              <h1>Einstellungen</h1>
              <p>Verwalte dein Konto und deine Praferenzen</p>
            </div>
          </div>
          <div className="settings-header__user">
            <div className="settings-header__avatar">
              {getUserInitials()}
            </div>
            <div className="settings-header__info">
              <span className="settings-header__name">
                {user?.first_name && user?.last_name
                  ? `${user.first_name} ${user.last_name}`
                  : user?.email?.split('@')[0]}
              </span>
              <span className="settings-header__email">{user?.email}</span>
            </div>
          </div>
        </motion.div>

        {/* Tab Navigation */}
        <motion.div
          className="settings-tabs"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          {tabs.map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`settings-tab ${activeTab === tab.id ? 'settings-tab--active' : ''}`}
            >
              <tab.icon className="settings-tab__icon" />
              {tab.label}
            </button>
          ))}
        </motion.div>

        {/* Tab Content */}
        <AnimatePresence mode="wait">
          {/* Profile Tab */}
          {activeTab === 'profile' && (
            <motion.div
              key="profile"
              initial="hidden"
              animate="visible"
              exit="hidden"
              variants={{ visible: { transition: { staggerChildren: 0.1 } } }}
            >
              <motion.div className="settings-card" variants={cardVariants}>
                <div className="settings-card__header">
                  <div className="settings-card__title">
                    <div className="settings-card__title-icon">
                      <User size={20} />
                    </div>
                    <h2>Personliche Informationen</h2>
                  </div>
                </div>

                {profileError && (
                  <motion.div
                    className="settings-error"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                  >
                    <AlertCircle size={18} />
                    <span>{profileError}</span>
                  </motion.div>
                )}

                <form onSubmit={handleProfileSave} className="settings-form">
                  <div className="settings-form-row">
                    <div className="settings-form-group">
                      <label className="settings-form-label">
                        <User size={14} />
                        Vorname
                      </label>
                      <div className="settings-input-wrapper">
                        <input
                          type="text"
                          name="first_name"
                          value={profileForm.first_name}
                          onChange={handleProfileChange}
                          className="settings-input"
                          placeholder="Max"
                        />
                        <User size={18} className="settings-input-icon" />
                      </div>
                    </div>
                    <div className="settings-form-group">
                      <label className="settings-form-label">
                        <User size={14} />
                        Nachname
                      </label>
                      <div className="settings-input-wrapper">
                        <input
                          type="text"
                          name="last_name"
                          value={profileForm.last_name}
                          onChange={handleProfileChange}
                          className="settings-input"
                          placeholder="Mustermann"
                        />
                        <User size={18} className="settings-input-icon" />
                      </div>
                    </div>
                  </div>

                  <div className="settings-form-group">
                    <label className="settings-form-label">
                      <Mail size={14} />
                      E-Mail Adresse
                    </label>
                    <div className="settings-input-wrapper">
                      <input
                        type="email"
                        name="email"
                        value={profileForm.email}
                        onChange={handleProfileChange}
                        className="settings-input"
                        placeholder="name@beispiel.de"
                        disabled
                      />
                      <Mail size={18} className="settings-input-icon" />
                    </div>
                    <span className="settings-input-hint">E-Mail kann nicht geandert werden</span>
                  </div>

                  <div className="settings-form-actions">
                    <div>
                      {profileSaved && (
                        <motion.span
                          className="settings-success"
                          initial={{ opacity: 0, x: -10 }}
                          animate={{ opacity: 1, x: 0 }}
                        >
                          <Check size={16} />
                          Profil erfolgreich gespeichert
                        </motion.span>
                      )}
                    </div>
                    <motion.button
                      type="submit"
                      className="settings-btn settings-btn--primary"
                      disabled={profileSaving}
                      whileHover={{ scale: 1.02 }}
                      whileTap={{ scale: 0.98 }}
                    >
                      {profileSaving ? (
                        <>
                          <span className="settings-spinner" />
                          <span>Speichern...</span>
                        </>
                      ) : (
                        <>
                          <Save size={18} />
                          <span>Speichern</span>
                        </>
                      )}
                    </motion.button>
                  </div>
                </form>
              </motion.div>

              {/* Role Card */}
              <motion.div
                className="settings-card settings-card--accent"
                variants={cardVariants}
              >
                <div className="settings-role-card">
                  <div className="settings-role-info">
                    <span className="settings-role-info__label">Deine Rolle</span>
                    <div className="settings-role-info__value">
                      <span className={`settings-badge ${getRoleBadgeClass(user?.role || '')}`}>
                        {user?.role || 'Unbekannt'}
                      </span>
                      {user?.role === 'owner' && (
                        <Crown size={16} style={{ color: 'var(--gold)' }} />
                      )}
                    </div>
                  </div>
                  <div className="settings-date-info">
                    <span className="settings-date-info__label">Mitglied seit</span>
                    <span className="settings-date-info__value">
                      {user?.created_at
                        ? new Date(user.created_at).toLocaleDateString('de-DE', {
                            day: '2-digit',
                            month: 'long',
                            year: 'numeric'
                          })
                        : '-'}
                    </span>
                  </div>
                </div>
              </motion.div>
            </motion.div>
          )}

          {/* Account Tab */}
          {activeTab === 'account' && (
            <motion.div
              key="account"
              initial="hidden"
              animate="visible"
              exit="hidden"
              variants={{ visible: { transition: { staggerChildren: 0.1 } } }}
            >
              <motion.div className="settings-card" variants={cardVariants}>
                <div className="settings-card__header">
                  <div className="settings-card__title">
                    <div className="settings-card__title-icon">
                      <Building2 size={20} />
                    </div>
                    <h2>Organisation</h2>
                  </div>
                </div>

                <div className="settings-info-grid">
                  <div className="settings-info-item">
                    <span className="settings-info-label">Organisationsname</span>
                    <span className="settings-info-value">{tenant?.name || '-'}</span>
                  </div>
                  <div className="settings-info-item">
                    <span className="settings-info-label">Subdomain</span>
                    <span className="settings-info-value settings-info-value--mono">
                      {tenant?.subdomain || '-'}.ginvault.app
                    </span>
                  </div>
                  <div className="settings-info-item">
                    <span className="settings-info-label">Status</span>
                    <span className={`settings-badge ${tenant?.status === 'active' ? 'settings-badge--status-active' : 'settings-badge--status-inactive'}`}>
                      {tenant?.status === 'active' ? 'Aktiv' : 'Inaktiv'}
                    </span>
                  </div>
                  <div className="settings-info-item">
                    <span className="settings-info-label">Erstellt am</span>
                    <span className="settings-info-value">
                      {tenant?.created_at
                        ? new Date(tenant.created_at).toLocaleDateString('de-DE')
                        : '-'}
                    </span>
                  </div>
                </div>
              </motion.div>

              {/* Subscription Info */}
              <motion.div className="settings-card" variants={cardVariants}>
                <div className="settings-card__header">
                  <div className="settings-card__title">
                    <div className="settings-card__title-icon">
                      <CreditCard size={20} />
                    </div>
                    <h2>Abonnement</h2>
                  </div>
                  <span className={`settings-badge ${getTierBadgeClass(tenant?.tier || '')}`}>
                    <Sparkles size={12} />
                    {tenant?.tier || 'Free'} Plan
                  </span>
                </div>

                <div className="settings-subscription-stats">
                  <div className="settings-stat">
                    <span className="settings-stat__value">
                      {tenant?.tier === 'free' ? '10' : tenant?.tier === 'basic' ? '50' : String.fromCharCode(8734)}
                    </span>
                    <span className="settings-stat__label">Max Gins</span>
                  </div>
                  <div className="settings-stat">
                    <span className="settings-stat__value">
                      {tenant?.tier === 'free' ? '1' : tenant?.tier === 'basic' ? '3' : tenant?.tier === 'pro' ? '10' : '50'}
                    </span>
                    <span className="settings-stat__label">Fotos/Gin</span>
                  </div>
                  <div className="settings-stat">
                    <span className={`settings-stat__value ${tenant?.tier === 'pro' || tenant?.tier === 'enterprise' ? 'settings-stat__value--check' : 'settings-stat__value--cross'}`}>
                      {tenant?.tier === 'pro' || tenant?.tier === 'enterprise' ? String.fromCharCode(10003) : String.fromCharCode(10007)}
                    </span>
                    <span className="settings-stat__label">Botanicals</span>
                  </div>
                  <div className="settings-stat">
                    <span className={`settings-stat__value ${tenant?.tier === 'enterprise' ? 'settings-stat__value--check' : 'settings-stat__value--cross'}`}>
                      {tenant?.tier === 'enterprise' ? String.fromCharCode(10003) : String.fromCharCode(10007)}
                    </span>
                    <span className="settings-stat__label">Multi-User</span>
                  </div>
                </div>

                <motion.button
                  onClick={() => window.location.href = '/subscription'}
                  className="settings-btn settings-btn--secondary settings-btn--full"
                  whileHover={{ scale: 1.01 }}
                  whileTap={{ scale: 0.99 }}
                >
                  <CreditCard size={18} />
                  Abonnement verwalten
                </motion.button>
              </motion.div>
            </motion.div>
          )}

          {/* Security Tab */}
          {activeTab === 'security' && (
            <motion.div
              key="security"
              initial="hidden"
              animate="visible"
              exit="hidden"
              variants={{ visible: { transition: { staggerChildren: 0.1 } } }}
            >
              <motion.div className="settings-card" variants={cardVariants}>
                <div className="settings-card__header">
                  <div className="settings-card__title">
                    <div className="settings-card__title-icon">
                      <Lock size={20} />
                    </div>
                    <h2>Passwort andern</h2>
                  </div>
                </div>

                {passwordError && (
                  <motion.div
                    className="settings-error"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                  >
                    <AlertCircle size={18} />
                    <span>{passwordError}</span>
                  </motion.div>
                )}

                <form onSubmit={handlePasswordSave} className="settings-form">
                  <div className="settings-form-group">
                    <label className="settings-form-label">
                      <Lock size={14} />
                      Aktuelles Passwort
                    </label>
                    <div className="settings-input-wrapper">
                      <input
                        type={showCurrentPassword ? 'text' : 'password'}
                        name="current_password"
                        value={passwordForm.current_password}
                        onChange={handlePasswordChange}
                        className="settings-input"
                        placeholder="Aktuelles Passwort eingeben"
                        required
                        style={{ paddingRight: '48px' }}
                      />
                      <Lock size={18} className="settings-input-icon" />
                      <button
                        type="button"
                        className="settings-password-toggle"
                        onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                      >
                        {showCurrentPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                      </button>
                    </div>
                  </div>

                  <div className="settings-form-group">
                    <label className="settings-form-label">
                      <Lock size={14} />
                      Neues Passwort
                    </label>
                    <div className="settings-input-wrapper">
                      <input
                        type={showNewPassword ? 'text' : 'password'}
                        name="new_password"
                        value={passwordForm.new_password}
                        onChange={handlePasswordChange}
                        className="settings-input"
                        placeholder="Neues Passwort eingeben"
                        minLength={8}
                        required
                        style={{ paddingRight: '48px' }}
                      />
                      <Lock size={18} className="settings-input-icon" />
                      <button
                        type="button"
                        className="settings-password-toggle"
                        onClick={() => setShowNewPassword(!showNewPassword)}
                      >
                        {showNewPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                      </button>
                    </div>
                    <span className="settings-input-hint">Mindestens 8 Zeichen</span>
                  </div>

                  <div className="settings-form-group">
                    <label className="settings-form-label">
                      <Lock size={14} />
                      Passwort bestatigen
                    </label>
                    <div className="settings-input-wrapper">
                      <input
                        type="password"
                        name="confirm_password"
                        value={passwordForm.confirm_password}
                        onChange={handlePasswordChange}
                        className="settings-input"
                        placeholder="Neues Passwort bestatigen"
                        required
                      />
                      <Lock size={18} className="settings-input-icon" />
                    </div>
                  </div>

                  <div className="settings-form-actions">
                    <div>
                      {passwordSaved && (
                        <motion.span
                          className="settings-success"
                          initial={{ opacity: 0, x: -10 }}
                          animate={{ opacity: 1, x: 0 }}
                        >
                          <Check size={16} />
                          Passwort erfolgreich geandert
                        </motion.span>
                      )}
                    </div>
                    <motion.button
                      type="submit"
                      className="settings-btn settings-btn--primary"
                      disabled={passwordSaving}
                      whileHover={{ scale: 1.02 }}
                      whileTap={{ scale: 0.98 }}
                    >
                      {passwordSaving ? (
                        <>
                          <span className="settings-spinner" />
                          <span>Speichern...</span>
                        </>
                      ) : (
                        <>
                          <Shield size={18} />
                          <span>Passwort andern</span>
                        </>
                      )}
                    </motion.button>
                  </div>
                </form>
              </motion.div>

              {/* Security Info */}
              <motion.div
                className="settings-card settings-card--accent"
                variants={cardVariants}
              >
                <div className="settings-card__header">
                  <div className="settings-card__title">
                    <div className="settings-card__title-icon">
                      <Shield size={20} />
                    </div>
                    <h2>Sicherheitsinformationen</h2>
                  </div>
                </div>
                <div className="settings-security-grid">
                  <div className="settings-security-item">
                    <span className="settings-security-label">Letzter Login</span>
                    <span className="settings-security-value">Heute</span>
                  </div>
                  <div className="settings-security-item">
                    <span className="settings-security-label">Account Status</span>
                    <span className={`settings-badge ${user?.is_active ? 'settings-badge--status-active' : 'settings-badge--status-inactive'}`}>
                      {user?.is_active ? 'Aktiv' : 'Inaktiv'}
                    </span>
                  </div>
                </div>
              </motion.div>

              {/* Danger Zone */}
              <motion.div
                className="settings-card settings-card--danger"
                variants={cardVariants}
              >
                <div className="settings-card__header">
                  <div className="settings-card__title settings-danger-title">
                    <div className="settings-card__title-icon" style={{ background: 'rgba(239, 68, 68, 0.15)', color: '#f87171' }}>
                      <AlertTriangle size={20} />
                    </div>
                    <h2>Gefahrenzone</h2>
                  </div>
                </div>
                <p className="settings-danger-text">
                  Wenn du deinen Account loschst, gibt es kein Zuruck. Bitte sei dir sicher.
                </p>
                <motion.button
                  className="settings-btn settings-btn--danger"
                  onClick={() => alert('Account-Loschung ist noch nicht implementiert. Kontaktiere den Support.')}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <AlertTriangle size={18} />
                  Account loschen
                </motion.button>
              </motion.div>
            </motion.div>
          )}

          {/* Notifications Tab */}
          {activeTab === 'notifications' && (
            <motion.div
              key="notifications"
              initial="hidden"
              animate="visible"
              exit="hidden"
              variants={{ visible: { transition: { staggerChildren: 0.1 } } }}
            >
              <motion.div className="settings-card" variants={cardVariants}>
                <div className="settings-card__header">
                  <div className="settings-card__title">
                    <div className="settings-card__title-icon">
                      <Bell size={20} />
                    </div>
                    <h2>Benachrichtigungseinstellungen</h2>
                  </div>
                </div>

                <div className="settings-toggles">
                  <div className="settings-toggle-item">
                    <div className="settings-toggle-content">
                      <span className="settings-toggle-title">E-Mail Benachrichtigungen</span>
                      <span className="settings-toggle-description">Erhalte Updates uber deine Sammlung</span>
                    </div>
                    <label className="settings-toggle">
                      <input type="checkbox" defaultChecked />
                      <span className="settings-toggle__slider" />
                    </label>
                  </div>

                  <div className="settings-toggle-item">
                    <div className="settings-toggle-content">
                      <span className="settings-toggle-title">Abonnement-Hinweise</span>
                      <span className="settings-toggle-description">Werde uber Abrechnungen und Plan-Anderungen informiert</span>
                    </div>
                    <label className="settings-toggle">
                      <input type="checkbox" defaultChecked />
                      <span className="settings-toggle__slider" />
                    </label>
                  </div>

                  <div className="settings-toggle-item">
                    <div className="settings-toggle-content">
                      <span className="settings-toggle-title">Marketing E-Mails</span>
                      <span className="settings-toggle-description">Erhalte Tipps, Neuigkeiten und Sonderangebote</span>
                    </div>
                    <label className="settings-toggle">
                      <input type="checkbox" />
                      <span className="settings-toggle__slider" />
                    </label>
                  </div>

                  <div className="settings-toggle-item">
                    <div className="settings-toggle-content">
                      <span className="settings-toggle-title">Wochentliche Zusammenfassung</span>
                      <span className="settings-toggle-description">Erhalte eine wochentliche Ubersicht deiner Sammlungs-Statistiken</span>
                    </div>
                    <label className="settings-toggle">
                      <input type="checkbox" />
                      <span className="settings-toggle__slider" />
                    </label>
                  </div>
                </div>

                <div className="settings-notification-footer">
                  <Mail size={16} />
                  Benachrichtigungen werden gesendet an: <span>{user?.email}</span>
                </div>
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
};

export default Settings;
