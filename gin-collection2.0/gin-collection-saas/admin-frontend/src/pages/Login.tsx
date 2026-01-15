import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import {
  Wine,
  Mail,
  Lock,
  Eye,
  EyeOff,
  AlertCircle,
  Shield,
  Loader2
} from 'lucide-react';
import './Login.css';

export default function Login() {
  const navigate = useNavigate();
  const { login, isAuthenticated, isLoading, error, clearError } = useAuthStore();
  const [email, setEmail] = useState('admin@gin-collection.local');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/');
    }
  }, [isAuthenticated, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    clearError();
    setIsSubmitting(true);

    const success = await login(email, password);
    if (success) {
      navigate('/');
    }

    setIsSubmitting(false);
  };

  const cardVariants = {
    hidden: { opacity: 0, y: 30, scale: 0.95 },
    visible: {
      opacity: 1,
      y: 0,
      scale: 1,
      transition: {
        type: 'spring',
        stiffness: 100,
        damping: 15
      }
    }
  };

  const formVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
        delayChildren: 0.2
      }
    }
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 15 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { type: 'spring', stiffness: 100, damping: 15 }
    }
  };

  if (isLoading) {
    return (
      <div className="admin-login-page">
        <div className="admin-login-loader">
          <div className="spinner" />
        </div>
      </div>
    );
  }

  return (
    <div className="admin-login-page">
      {/* Ambient Background */}
      <div className="admin-login-ambient">
        <div className="admin-login-orb admin-login-orb--gold" />
        <div className="admin-login-orb admin-login-orb--green" />
        <div className="admin-login-orb admin-login-orb--mint" />
      </div>

      {/* Pattern Overlay */}
      <div className="admin-login-pattern" />

      <motion.div
        className="admin-login-card"
        variants={cardVariants}
        initial="hidden"
        animate="visible"
      >
        {/* Logo Section */}
        <div className="admin-login-header">
          <motion.div
            className="admin-login-logo"
            initial={{ scale: 0, rotate: -180 }}
            animate={{ scale: 1, rotate: 0 }}
            transition={{ type: 'spring', stiffness: 200, damping: 15, delay: 0.1 }}
          >
            <Wine />
          </motion.div>
          <motion.h1
            className="admin-login-title"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
          >
            GinVault
          </motion.h1>
          <motion.div
            className="admin-login-badge"
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.3 }}
          >
            <Shield />
            <span>Platform Admin</span>
          </motion.div>
        </div>

        {/* Error Message */}
        {error && (
          <motion.div
            className="admin-login-error"
            initial={{ opacity: 0, y: -10, height: 0 }}
            animate={{ opacity: 1, y: 0, height: 'auto' }}
            exit={{ opacity: 0, y: -10, height: 0 }}
          >
            <AlertCircle />
            <span>{error}</span>
          </motion.div>
        )}

        {/* Login Form */}
        <motion.form
          onSubmit={handleSubmit}
          className="admin-login-form"
          variants={formVariants}
          initial="hidden"
          animate="visible"
        >
          <motion.div className="admin-login-field" variants={itemVariants}>
            <label className="admin-login-label">
              <Mail />
              Email
            </label>
            <div className="admin-login-input-wrapper">
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="admin-login-input"
                placeholder="admin@example.com"
                required
                autoComplete="email"
              />
            </div>
          </motion.div>

          <motion.div className="admin-login-field" variants={itemVariants}>
            <label className="admin-login-label">
              <Lock />
              Passwort
            </label>
            <div className="admin-login-input-wrapper">
              <input
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="admin-login-input admin-login-input--password"
                placeholder="Passwort eingeben"
                required
                autoComplete="current-password"
              />
              <button
                type="button"
                className="admin-login-password-toggle"
                onClick={() => setShowPassword(!showPassword)}
              >
                {showPassword ? <EyeOff /> : <Eye />}
              </button>
            </div>
          </motion.div>

          <motion.button
            type="submit"
            className="admin-login-submit"
            disabled={isSubmitting}
            variants={itemVariants}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            {isSubmitting ? (
              <>
                <Loader2 className="admin-login-submit-spinner" />
                Anmeldung läuft...
              </>
            ) : (
              <>
                <Shield />
                Admin Login
              </>
            )}
          </motion.button>
        </motion.form>

        {/* Footer */}
        <motion.div
          className="admin-login-footer"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
        >
          <p className="admin-login-hint">
            Standard: admin@gin-collection.local / Test123456
          </p>
          <div className="admin-login-divider">
            <span>Sicherer Zugang</span>
          </div>
          <p className="admin-login-security">
            Dieser Bereich ist nur für autorisierte Platform-Administratoren zugänglich.
          </p>
        </motion.div>
      </motion.div>
    </div>
  );
}
