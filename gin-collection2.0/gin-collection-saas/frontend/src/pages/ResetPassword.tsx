import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import apiClient, { getErrorMessage } from '../api/client';
import {
  ShieldCheck,
  Lock,
  Eye,
  EyeOff,
  AlertCircle,
  ArrowRight,
  CheckCircle,
  XCircle,
  Loader2
} from 'lucide-react';
import './Login.css';

const ResetPassword = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');

  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [isValidating, setIsValidating] = useState(true);
  const [isTokenValid, setIsTokenValid] = useState(false);

  // Validate token on mount
  useEffect(() => {
    const validateToken = async () => {
      if (!token) {
        setIsValidating(false);
        setIsTokenValid(false);
        return;
      }

      try {
        const response = await apiClient.get(`/auth/validate-reset-token?token=${token}`);
        setIsTokenValid(response.data.valid);
      } catch (err) {
        setIsTokenValid(false);
      } finally {
        setIsValidating(false);
      }
    };

    validateToken();
  }, [token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validate passwords match
    if (password !== confirmPassword) {
      setError('Die Passwörter stimmen nicht überein.');
      return;
    }

    // Validate password length
    if (password.length < 8) {
      setError('Das Passwort muss mindestens 8 Zeichen lang sein.');
      return;
    }

    setIsLoading(true);

    try {
      await apiClient.post('/auth/reset-password', {
        token,
        new_password: password,
      });
      setIsSuccess(true);

      // Redirect to login after 3 seconds
      setTimeout(() => {
        navigate('/login');
      }, 3000);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  // Loading state while validating token
  if (isValidating) {
    return (
      <div className="login-page">
        <div className="login-page__ambient">
          <div className="login-ambient-orb login-ambient-orb--1" />
          <div className="login-ambient-orb login-ambient-orb--2" />
        </div>
        <motion.div
          className="login-card"
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          style={{ textAlign: 'center', padding: '60px 40px' }}
        >
          <Loader2 size={48} className="login-spinner" style={{ margin: '0 auto 20px', width: '48px', height: '48px' }} />
          <p style={{ color: 'var(--text-muted)' }}>Link wird überprüft...</p>
        </motion.div>
      </div>
    );
  }

  // Invalid or expired token
  if (!isTokenValid) {
    return (
      <div className="login-page">
        <div className="login-page__ambient">
          <div className="login-ambient-orb login-ambient-orb--1" />
          <div className="login-ambient-orb login-ambient-orb--2" />
        </div>
        <div className="login-page__decor">
          <div className="decor-line decor-line--1" />
          <div className="decor-line decor-line--2" />
        </div>
        <motion.div
          className="login-card"
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <div style={{ textAlign: 'center' }}>
            <div
              style={{
                width: '80px',
                height: '80px',
                margin: '0 auto 24px',
                borderRadius: '50%',
                background: 'rgba(220, 38, 38, 0.15)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <XCircle size={40} style={{ color: '#F87171' }} />
            </div>
            <h1 className="login-title" style={{ fontSize: '1.5rem', marginBottom: '12px' }}>
              Link ungültig
            </h1>
            <p style={{ color: 'var(--text-muted)', marginBottom: '32px' }}>
              Dieser Link zum Zurücksetzen des Passworts ist ungültig oder abgelaufen.
              Bitte fordere einen neuen Link an.
            </p>
            <Link to="/forgot-password">
              <motion.button
                className="login-submit"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <span className="login-submit__content">
                  <span>Neuen Link anfordern</span>
                  <ArrowRight size={18} />
                </span>
              </motion.button>
            </Link>
            <div style={{ marginTop: '24px' }}>
              <Link to="/login" className="forgot-password" style={{ textAlign: 'center' }}>
                Zurück zum Login
              </Link>
            </div>
          </div>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="login-page">
      {/* Ambient Background */}
      <div className="login-page__ambient">
        <div className="login-ambient-orb login-ambient-orb--1" />
        <div className="login-ambient-orb login-ambient-orb--2" />
        <div className="login-ambient-orb login-ambient-orb--3" />
      </div>

      {/* Decorative Lines */}
      <div className="login-page__decor">
        <div className="decor-line decor-line--1" />
        <div className="decor-line decor-line--2" />
        <div className="decor-line decor-line--3" />
        <div className="decor-line decor-line--4" />
      </div>

      {/* Card */}
      <motion.div
        className="login-card"
        initial={{ opacity: 0, y: 30, scale: 0.95 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
      >
        {/* Header */}
        <motion.div
          className="login-header"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.1 }}
        >
          <motion.div
            className="login-logo"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{
              type: 'spring',
              stiffness: 200,
              damping: 15,
              delay: 0.2
            }}
          >
            <ShieldCheck size={36} />
          </motion.div>
          <h1 className="login-title">
            {isSuccess ? 'Geschafft!' : 'Neues Passwort'}
          </h1>
          <p className="login-subtitle">
            {isSuccess
              ? 'Dein Passwort wurde geändert'
              : 'Wähle ein neues sicheres Passwort'}
          </p>
        </motion.div>

        {isSuccess ? (
          /* Success State */
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.4 }}
          >
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                gap: '20px',
                padding: '20px 0',
              }}
            >
              <div
                style={{
                  width: '64px',
                  height: '64px',
                  borderRadius: '50%',
                  background: 'rgba(126, 205, 160, 0.15)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                <CheckCircle size={32} style={{ color: 'var(--mint)' }} />
              </div>
              <div style={{ textAlign: 'center' }}>
                <p style={{ color: 'var(--text-primary)', marginBottom: '8px' }}>
                  Passwort erfolgreich geändert!
                </p>
                <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>
                  Du wirst in wenigen Sekunden zur Anmeldung weitergeleitet...
                </p>
              </div>
            </div>

            <Link to="/login">
              <motion.button
                type="button"
                className="login-submit"
                style={{ marginTop: '24px' }}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <span className="login-submit__content">
                  <span>Jetzt anmelden</span>
                  <ArrowRight size={18} />
                </span>
              </motion.button>
            </Link>
          </motion.div>
        ) : (
          /* Form State */
          <>
            {/* Error Message */}
            {error && (
              <motion.div
                className="login-error"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.3 }}
              >
                <AlertCircle size={18} />
                <span>{error}</span>
              </motion.div>
            )}

            {/* Form */}
            <motion.form
              className="login-form"
              onSubmit={handleSubmit}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              {/* New Password Field */}
              <div className="form-group">
                <label className="form-label" htmlFor="password">
                  Neues Passwort
                </label>
                <div className="form-input-wrapper">
                  <input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="form-input"
                    placeholder="Mindestens 8 Zeichen"
                    required
                    disabled={isLoading}
                    autoComplete="new-password"
                    minLength={8}
                    style={{ paddingRight: '48px' }}
                    autoFocus
                  />
                  <Lock size={18} className="form-input-icon" />
                  <button
                    type="button"
                    className="password-toggle"
                    onClick={() => setShowPassword(!showPassword)}
                    tabIndex={-1}
                  >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>

              {/* Confirm Password Field */}
              <div className="form-group">
                <label className="form-label" htmlFor="confirmPassword">
                  Passwort bestätigen
                </label>
                <div className="form-input-wrapper">
                  <input
                    id="confirmPassword"
                    type={showConfirmPassword ? 'text' : 'password'}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    className="form-input"
                    placeholder="Passwort wiederholen"
                    required
                    disabled={isLoading}
                    autoComplete="new-password"
                    minLength={8}
                    style={{ paddingRight: '48px' }}
                  />
                  <Lock size={18} className="form-input-icon" />
                  <button
                    type="button"
                    className="password-toggle"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    tabIndex={-1}
                  >
                    {showConfirmPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>

              {/* Password Requirements */}
              <div
                style={{
                  fontSize: '0.8rem',
                  color: 'var(--text-muted)',
                  display: 'flex',
                  flexDirection: 'column',
                  gap: '4px',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <div
                    style={{
                      width: '6px',
                      height: '6px',
                      borderRadius: '50%',
                      background: password.length >= 8 ? 'var(--mint)' : 'var(--text-muted)',
                    }}
                  />
                  <span style={{ color: password.length >= 8 ? 'var(--mint)' : 'inherit' }}>
                    Mindestens 8 Zeichen
                  </span>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <div
                    style={{
                      width: '6px',
                      height: '6px',
                      borderRadius: '50%',
                      background:
                        password && password === confirmPassword ? 'var(--mint)' : 'var(--text-muted)',
                    }}
                  />
                  <span
                    style={{
                      color: password && password === confirmPassword ? 'var(--mint)' : 'inherit',
                    }}
                  >
                    Passwörter stimmen überein
                  </span>
                </div>
              </div>

              {/* Submit Button */}
              <motion.button
                type="submit"
                className="login-submit"
                disabled={isLoading || password.length < 8 || password !== confirmPassword}
                whileHover={{ scale: isLoading ? 1 : 1.02 }}
                whileTap={{ scale: isLoading ? 1 : 0.98 }}
              >
                <span className="login-submit__content">
                  {isLoading ? (
                    <>
                      <span className="login-spinner" />
                      <span>Speichern...</span>
                    </>
                  ) : (
                    <>
                      <span>Passwort speichern</span>
                      <ArrowRight size={18} />
                    </>
                  )}
                </span>
              </motion.button>
            </motion.form>
          </>
        )}
      </motion.div>
    </div>
  );
};

export default ResetPassword;
