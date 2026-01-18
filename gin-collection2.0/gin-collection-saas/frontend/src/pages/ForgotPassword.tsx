import { useState } from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import apiClient from '../api/client';
import {
  KeyRound,
  Mail,
  AlertCircle,
  ArrowRight,
  ArrowLeft,
  CheckCircle
} from 'lucide-react';
import './Login.css';

const ForgotPassword = () => {
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await apiClient.post('/auth/forgot-password', { email });
      setIsSuccess(true);
    } catch (err: any) {
      // Always show success to prevent email enumeration
      setIsSuccess(true);
    } finally {
      setIsLoading(false);
    }
  };

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
            <KeyRound size={36} />
          </motion.div>
          <h1 className="login-title">Passwort vergessen?</h1>
          <p className="login-subtitle">
            {isSuccess
              ? 'Prüfe dein E-Mail Postfach'
              : 'Kein Problem, wir helfen dir'}
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
                  E-Mail wurde gesendet!
                </p>
                <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>
                  Wenn ein Konto mit der E-Mail <strong style={{ color: 'var(--gold)' }}>{email}</strong> existiert,
                  erhältst du in Kürze einen Link zum Zurücksetzen deines Passworts.
                </p>
              </div>
              <p style={{ color: 'var(--text-muted)', fontSize: '0.85rem', marginTop: '8px' }}>
                Der Link ist 1 Stunde gültig.
              </p>
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
                  <ArrowLeft size={18} />
                  <span>Zurück zum Login</span>
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
              <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem', marginBottom: '8px' }}>
                Gib deine E-Mail-Adresse ein und wir senden dir einen Link zum Zurücksetzen deines Passworts.
              </p>

              {/* Email Field */}
              <div className="form-group">
                <label className="form-label" htmlFor="email">
                  E-Mail Adresse
                </label>
                <div className="form-input-wrapper">
                  <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="form-input"
                    placeholder="name@beispiel.de"
                    required
                    disabled={isLoading}
                    autoComplete="email"
                    autoFocus
                  />
                  <Mail size={18} className="form-input-icon" />
                </div>
              </div>

              {/* Submit Button */}
              <motion.button
                type="submit"
                className="login-submit"
                disabled={isLoading}
                whileHover={{ scale: isLoading ? 1 : 1.02 }}
                whileTap={{ scale: isLoading ? 1 : 0.98 }}
              >
                <span className="login-submit__content">
                  {isLoading ? (
                    <>
                      <span className="login-spinner" />
                      <span>Senden...</span>
                    </>
                  ) : (
                    <>
                      <span>Link senden</span>
                      <ArrowRight size={18} />
                    </>
                  )}
                </span>
              </motion.button>
            </motion.form>

            {/* Back to Login */}
            <div className="login-divider">
              <div className="login-divider__line" />
              <span className="login-divider__text">oder</span>
              <div className="login-divider__line" />
            </div>

            <motion.div
              className="login-footer"
              style={{ borderTop: 'none', marginTop: 0, paddingTop: 0 }}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.6, delay: 0.3 }}
            >
              <p className="login-footer__text">
                Passwort wieder eingefallen?{' '}
                <Link to="/login" className="login-footer__link">
                  Zurück zum Login
                </Link>
              </p>
            </motion.div>
          </>
        )}
      </motion.div>
    </div>
  );
};

export default ForgotPassword;
