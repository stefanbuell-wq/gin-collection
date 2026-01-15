import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import { getErrorMessage } from '../api/client';
import {
  Sparkles,
  Mail,
  Lock,
  Eye,
  EyeOff,
  AlertCircle,
  ArrowRight,
  Info
} from 'lucide-react';
import './Login.css';

// Demo accounts for quick testing
const demoAccounts = [
  { email: 'test@test.com', tier: 'Free' },
  { email: 'basic@demo.local', tier: 'Basic' },
  { email: 'pro@demo.local', tier: 'Pro' },
  { email: 'enterprise@demo.local', tier: 'Enterprise' },
];

const Login = () => {
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showDemoAccounts, setShowDemoAccounts] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(email, password);
      navigate('/dashboard');
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  const handleDemoAccountClick = (demoEmail: string) => {
    setEmail(demoEmail);
    setPassword('Test123456');
    setShowDemoAccounts(false);
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

      {/* Login Card */}
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
            <Sparkles size={36} />
          </motion.div>
          <h1 className="login-title">GinVault</h1>
          <p className="login-subtitle">Willkommen zurück</p>
        </motion.div>

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

        {/* Login Form */}
        <motion.form
          className="login-form"
          onSubmit={handleSubmit}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
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
              />
              <Mail size={18} className="form-input-icon" />
            </div>
          </div>

          {/* Password Field */}
          <div className="form-group">
            <label className="form-label" htmlFor="password">
              Passwort
            </label>
            <div className="form-input-wrapper">
              <input
                id="password"
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="form-input"
                placeholder="••••••••••"
                required
                disabled={isLoading}
                autoComplete="current-password"
                style={{ paddingRight: '48px' }}
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

          {/* Forgot Password */}
          <Link to="/forgot-password" className="forgot-password">
            Passwort vergessen?
          </Link>

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
                  <span>Anmelden...</span>
                </>
              ) : (
                <>
                  <span>Anmelden</span>
                  <ArrowRight size={18} />
                </>
              )}
            </span>
          </motion.button>
        </motion.form>

        {/* Demo Accounts Toggle */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.4 }}
        >
          <button
            type="button"
            className="forgot-password"
            style={{ width: '100%', textAlign: 'center', marginTop: '16px' }}
            onClick={() => setShowDemoAccounts(!showDemoAccounts)}
          >
            <Info size={14} style={{ display: 'inline', marginRight: '6px', verticalAlign: 'middle' }} />
            {showDemoAccounts ? 'Demo-Accounts ausblenden' : 'Demo-Accounts anzeigen'}
          </button>

          {/* Demo Accounts List */}
          {showDemoAccounts && (
            <motion.div
              className="demo-accounts"
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.3 }}
            >
              <div className="demo-accounts__title">
                <Info size={14} />
                Test-Accounts (Passwort: Test123456)
              </div>
              <div className="demo-accounts__list">
                {demoAccounts.map((account) => (
                  <div
                    key={account.email}
                    className="demo-account"
                    onClick={() => handleDemoAccountClick(account.email)}
                  >
                    <span className="demo-account__email">{account.email}</span>
                    <span className="demo-account__tier">{account.tier}</span>
                  </div>
                ))}
              </div>
            </motion.div>
          )}
        </motion.div>

        {/* Divider */}
        <div className="login-divider">
          <div className="login-divider__line" />
          <span className="login-divider__text">Neu hier?</span>
          <div className="login-divider__line" />
        </div>

        {/* Register Link */}
        <motion.div
          className="login-footer"
          style={{ borderTop: 'none', marginTop: 0, paddingTop: 0 }}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.3 }}
        >
          <p className="login-footer__text">
            Noch kein Account?{' '}
            <Link to="/register" className="login-footer__link">
              Jetzt registrieren
            </Link>
          </p>
        </motion.div>
      </motion.div>
    </div>
  );
};

export default Login;
