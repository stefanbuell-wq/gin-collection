import { useState, useMemo } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import { getErrorMessage } from '../api/client';
import {
  Sparkles,
  Building2,
  Globe,
  Mail,
  Lock,
  User,
  Eye,
  EyeOff,
  AlertCircle,
  ArrowRight,
  Check,
  Wine,
  Shield,
  Zap
} from 'lucide-react';
import './Register.css';

// Password strength calculator
const calculatePasswordStrength = (password: string): {
  strength: 'weak' | 'medium' | 'strong' | 'very-strong';
  label: string;
} => {
  if (!password) return { strength: 'weak', label: '' };

  let score = 0;
  if (password.length >= 8) score++;
  if (password.length >= 12) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/[a-z]/.test(password)) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[^A-Za-z0-9]/.test(password)) score++;

  if (score <= 2) return { strength: 'weak', label: 'Schwach' };
  if (score <= 3) return { strength: 'medium', label: 'Mittel' };
  if (score <= 4) return { strength: 'strong', label: 'Stark' };
  return { strength: 'very-strong', label: 'Sehr stark' };
};

const Register = () => {
  const navigate = useNavigate();
  const register = useAuthStore((state) => state.register);
  const [formData, setFormData] = useState({
    tenant_name: '',
    subdomain: '',
    email: '',
    password: '',
    first_name: '',
    last_name: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // Calculate password strength
  const passwordStrength = useMemo(
    () => calculatePasswordStrength(formData.password),
    [formData.password]
  );

  // Calculate current step based on filled fields
  const currentStep = useMemo(() => {
    if (!formData.tenant_name || !formData.subdomain) return 1;
    if (!formData.email || !formData.password) return 2;
    return 3;
  }, [formData]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    // Auto-format subdomain: lowercase, only alphanumeric and hyphens
    if (name === 'subdomain') {
      const formatted = value.toLowerCase().replace(/[^a-z0-9-]/g, '');
      setFormData({ ...formData, [name]: formatted });
      return;
    }

    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await register(formData);
      navigate('/dashboard');
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="register-page">
      {/* Ambient Background */}
      <div className="register-page__ambient">
        <div className="register-ambient-orb register-ambient-orb--1" />
        <div className="register-ambient-orb register-ambient-orb--2" />
        <div className="register-ambient-orb register-ambient-orb--3" />
      </div>

      {/* Decorative Pattern */}
      <div className="register-page__pattern" />

      {/* Register Card */}
      <motion.div
        className="register-card"
        initial={{ opacity: 0, y: 30, scale: 0.95 }}
        animate={{ opacity: 1, y: 0, scale: 1 }}
        transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
      >
        {/* Header */}
        <motion.div
          className="register-header"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.1 }}
        >
          <motion.div
            className="register-logo"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{
              type: 'spring',
              stiffness: 200,
              damping: 15,
              delay: 0.2
            }}
          >
            <Sparkles size={32} />
          </motion.div>
          <h1 className="register-title">GinVault</h1>
          <p className="register-subtitle">Erstelle deine Premium-Sammlung</p>
        </motion.div>

        {/* Progress Steps */}
        <motion.div
          className="register-steps"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.15 }}
        >
          <div className={`register-step ${currentStep >= 1 ? 'register-step--active' : ''} ${currentStep > 1 ? 'register-step--completed' : ''}`}>
            <div className="register-step__dot">
              {currentStep > 1 ? <Check size={14} /> : '1'}
            </div>
          </div>
          <div className="register-step__line" style={{ background: currentStep > 1 ? 'var(--mint)' : undefined }} />
          <div className={`register-step ${currentStep >= 2 ? 'register-step--active' : ''} ${currentStep > 2 ? 'register-step--completed' : ''}`}>
            <div className="register-step__dot">
              {currentStep > 2 ? <Check size={14} /> : '2'}
            </div>
          </div>
          <div className="register-step__line" style={{ background: currentStep > 2 ? 'var(--mint)' : undefined }} />
          <div className={`register-step ${currentStep >= 3 ? 'register-step--active' : ''}`}>
            <div className="register-step__dot">3</div>
          </div>
        </motion.div>

        {/* Error Message */}
        {error && (
          <motion.div
            className="register-error"
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.3 }}
          >
            <AlertCircle size={18} />
            <span>{error}</span>
          </motion.div>
        )}

        {/* Register Form */}
        <motion.form
          className="register-form"
          onSubmit={handleSubmit}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
          {/* Step 1: Organization */}
          <div className="register-form-group">
            <label className="register-form-label" htmlFor="tenant_name">
              <Building2 size={14} />
              Name deiner Sammlung
            </label>
            <div className="register-input-wrapper">
              <input
                id="tenant_name"
                type="text"
                name="tenant_name"
                value={formData.tenant_name}
                onChange={handleChange}
                className="register-input"
                placeholder="Meine Gin-Sammlung"
                required
                disabled={isLoading}
              />
              <Building2 size={16} className="register-input-icon" />
            </div>
          </div>

          <div className="register-form-group">
            <label className="register-form-label" htmlFor="subdomain">
              <Globe size={14} />
              Subdomain
            </label>
            <div className="register-subdomain-group">
              <div className="register-input-wrapper" style={{ flex: 1 }}>
                <input
                  id="subdomain"
                  type="text"
                  name="subdomain"
                  value={formData.subdomain}
                  onChange={handleChange}
                  className="register-input"
                  placeholder="meine-sammlung"
                  required
                  disabled={isLoading}
                  pattern="[a-z0-9-]+"
                  style={{ paddingRight: '14px' }}
                />
                <Globe size={16} className="register-input-icon" />
              </div>
              <span className="register-subdomain-suffix">.ginvault.app</span>
            </div>
          </div>

          {/* Step 2: Account */}
          <div className="register-form-group">
            <label className="register-form-label" htmlFor="email">
              <Mail size={14} />
              E-Mail Adresse
            </label>
            <div className="register-input-wrapper">
              <input
                id="email"
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className="register-input"
                placeholder="name@beispiel.de"
                required
                disabled={isLoading}
                autoComplete="email"
              />
              <Mail size={16} className="register-input-icon" />
            </div>
          </div>

          <div className="register-form-group">
            <label className="register-form-label" htmlFor="password">
              <Lock size={14} />
              Passwort
            </label>
            <div className="register-input-wrapper">
              <input
                id="password"
                type={showPassword ? 'text' : 'password'}
                name="password"
                value={formData.password}
                onChange={handleChange}
                className="register-input"
                placeholder="Mindestens 8 Zeichen"
                required
                minLength={8}
                disabled={isLoading}
                autoComplete="new-password"
                style={{ paddingRight: '48px' }}
              />
              <Lock size={16} className="register-input-icon" />
              <button
                type="button"
                className="register-password-toggle"
                onClick={() => setShowPassword(!showPassword)}
                tabIndex={-1}
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>
            {/* Password Strength Indicator */}
            {formData.password && (
              <>
                <div className={`password-strength password-strength--${passwordStrength.strength}`}>
                  <div className="password-strength__bar" />
                  <div className="password-strength__bar" />
                  <div className="password-strength__bar" />
                  <div className="password-strength__bar" />
                </div>
                <span className={`password-hint password-hint--${passwordStrength.strength}`}>
                  <Shield size={12} />
                  {passwordStrength.label}
                </span>
              </>
            )}
          </div>

          {/* Step 3: Personal Info (Optional) */}
          <div className="register-form-row">
            <div className="register-form-group">
              <label className="register-form-label" htmlFor="first_name">
                <User size={14} />
                Vorname
                <span className="register-form-label__optional">(optional)</span>
              </label>
              <div className="register-input-wrapper">
                <input
                  id="first_name"
                  type="text"
                  name="first_name"
                  value={formData.first_name}
                  onChange={handleChange}
                  className="register-input"
                  placeholder="Max"
                  disabled={isLoading}
                  autoComplete="given-name"
                />
                <User size={16} className="register-input-icon" />
              </div>
            </div>

            <div className="register-form-group">
              <label className="register-form-label" htmlFor="last_name">
                <User size={14} />
                Nachname
                <span className="register-form-label__optional">(optional)</span>
              </label>
              <div className="register-input-wrapper">
                <input
                  id="last_name"
                  type="text"
                  name="last_name"
                  value={formData.last_name}
                  onChange={handleChange}
                  className="register-input"
                  placeholder="Mustermann"
                  disabled={isLoading}
                  autoComplete="family-name"
                />
                <User size={16} className="register-input-icon" />
              </div>
            </div>
          </div>

          {/* Features List */}
          <div className="register-features">
            <div className="register-feature">
              <Wine size={14} />
              <span>Bis zu 5 Gins kostenlos</span>
            </div>
            <div className="register-feature">
              <Shield size={14} />
              <span>Sichere Datenspeicherung</span>
            </div>
            <div className="register-feature">
              <Zap size={14} />
              <span>Jederzeit upgraden</span>
            </div>
          </div>

          {/* Submit Button */}
          <motion.button
            type="submit"
            className="register-submit"
            disabled={isLoading}
            whileHover={{ scale: isLoading ? 1 : 1.02 }}
            whileTap={{ scale: isLoading ? 1 : 0.98 }}
          >
            <span className="register-submit__content">
              {isLoading ? (
                <>
                  <span className="register-spinner" />
                  <span>Account wird erstellt...</span>
                </>
              ) : (
                <>
                  <span>Kostenlos registrieren</span>
                  <ArrowRight size={18} />
                </>
              )}
            </span>
          </motion.button>

          {/* Terms */}
          <p className="register-terms">
            Mit der Registrierung akzeptierst du unsere{' '}
            <Link to="/terms">Nutzungsbedingungen</Link> und{' '}
            <Link to="/privacy">Datenschutzrichtlinie</Link>.
          </p>
        </motion.form>

        {/* Footer - Login Link */}
        <motion.div
          className="register-footer"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.3 }}
        >
          <p className="register-footer__text">
            Bereits einen Account?{' '}
            <Link to="/login" className="register-footer__link">
              Jetzt anmelden
            </Link>
          </p>
        </motion.div>
      </motion.div>
    </div>
  );
};

export default Register;
