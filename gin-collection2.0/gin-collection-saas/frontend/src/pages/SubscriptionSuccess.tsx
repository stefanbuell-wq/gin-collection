import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import { subscriptionAPI, tenantAPI } from '../api/services';
import {
  CheckCircle,
  Loader2,
  AlertTriangle,
  ArrowRight,
  Wine,
  Sparkles,
  Shield,
  Zap,
  HelpCircle,
  Check
} from 'lucide-react';
import './SubscriptionStatus.css';

const SubscriptionSuccess = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { setTenant } = useAuthStore();
  const [status, setStatus] = useState<'activating' | 'success' | 'error'>('activating');
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const activateSubscription = async () => {
      const subscriptionId = searchParams.get('subscription_id');

      if (!subscriptionId) {
        setStatus('error');
        setErrorMessage('Keine Abonnement-ID gefunden. Bitte kontaktiere den Support.');
        return;
      }

      try {
        await subscriptionAPI.activate(subscriptionId);

        const tenantResponse = await tenantAPI.getCurrent();
        const tenantData = tenantResponse.data as unknown as { success: boolean; data: { tenant: ReturnType<typeof useAuthStore.getState>['tenant'] } };
        if (tenantData.success && tenantData.data?.tenant) {
          setTenant(tenantData.data.tenant);
        }

        setStatus('success');
      } catch (err) {
        console.error('Failed to activate subscription:', err);
        setStatus('error');
        setErrorMessage('Das Abonnement konnte nicht aktiviert werden. Bitte kontaktiere den Support.');
      }
    };

    activateSubscription();
  }, [searchParams, setTenant]);

  const cardVariants = {
    hidden: { opacity: 0, y: 20, scale: 0.95 },
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

  const iconVariants = {
    hidden: { scale: 0 },
    visible: {
      scale: 1,
      transition: {
        type: 'spring',
        stiffness: 200,
        damping: 15,
        delay: 0.2
      }
    }
  };

  return (
    <div className="subscription-status-page">
      {/* Ambient Background */}
      <div className="subscription-status-ambient">
        <div className="subscription-status-orb subscription-status-orb--gold" />
        <div className="subscription-status-orb subscription-status-orb--mint" />
        <div className="subscription-status-orb subscription-status-orb--green" />
      </div>

      <motion.div
        className="subscription-status-card"
        variants={cardVariants}
        initial="hidden"
        animate="visible"
      >
        {/* Logo */}
        <div className="subscription-status-logo">
          <div className="subscription-status-logo__icon">
            <Wine />
          </div>
          <span className="subscription-status-logo__text">GinVault</span>
        </div>

        {/* Activating State */}
        {status === 'activating' && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3 }}
          >
            <motion.div
              className="subscription-status-icon subscription-status-icon--activating"
              variants={iconVariants}
              initial="hidden"
              animate="visible"
            >
              <Loader2 className="subscription-status-spinner" />
            </motion.div>

            <h1 className="subscription-status-title">
              Abonnement wird aktiviert
              <span className="subscription-status-dots">
                <span className="subscription-status-dot" />
                <span className="subscription-status-dot" />
                <span className="subscription-status-dot" />
              </span>
            </h1>

            <p className="subscription-status-description">
              Bitte warte einen Moment, während wir dein Abonnement aktivieren.
            </p>
          </motion.div>
        )}

        {/* Success State */}
        {status === 'success' && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3 }}
          >
            <motion.div
              className="subscription-status-icon subscription-status-icon--success subscription-status-pulse"
              variants={iconVariants}
              initial="hidden"
              animate="visible"
            >
              <CheckCircle />
            </motion.div>

            <div className="subscription-status-badge">
              <span className="subscription-status-badge__dot" />
              <span className="subscription-status-badge__text">Erfolgreich aktiviert</span>
            </div>

            <h1 className="subscription-status-title">Willkommen im Premium!</h1>

            <p className="subscription-status-description">
              Dein Abonnement wurde erfolgreich aktiviert. Du hast jetzt Zugang zu allen Premium-Funktionen.
            </p>

            <div className="subscription-status-features">
              <p className="subscription-status-features__title">Jetzt freigeschaltet</p>
              <ul className="subscription-status-features__list">
                <li className="subscription-status-features__item">
                  <Check />
                  Unbegrenzte Gin-Einträge
                </li>
                <li className="subscription-status-features__item">
                  <Check />
                  Erweiterte Statistiken & Analysen
                </li>
                <li className="subscription-status-features__item">
                  <Check />
                  Premium Support
                </li>
              </ul>
            </div>

            <div className="subscription-status-actions">
              <motion.button
                onClick={() => navigate('/subscription')}
                className="subscription-status-btn subscription-status-btn--primary"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                Zum Abonnement
                <ArrowRight />
              </motion.button>

              <motion.button
                onClick={() => navigate('/gins')}
                className="subscription-status-btn subscription-status-btn--secondary"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <Sparkles />
                Sammlung erkunden
              </motion.button>
            </div>
          </motion.div>
        )}

        {/* Error State */}
        {status === 'error' && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3 }}
          >
            <motion.div
              className="subscription-status-icon subscription-status-icon--error"
              variants={iconVariants}
              initial="hidden"
              animate="visible"
            >
              <AlertTriangle />
            </motion.div>

            <h1 className="subscription-status-title">Aktivierung fehlgeschlagen</h1>

            <p className="subscription-status-description">
              {errorMessage}
            </p>

            <div className="subscription-status-actions">
              <motion.button
                onClick={() => navigate('/subscription')}
                className="subscription-status-btn subscription-status-btn--primary"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                Zurück zum Abonnement
              </motion.button>
            </div>

            <div className="subscription-status-divider">
              <div className="subscription-status-support">
                <p className="subscription-status-support__text">
                  Das Problem besteht weiterhin?
                </p>
                <a
                  href="mailto:support@gin-collection.local"
                  className="subscription-status-support__link"
                >
                  <HelpCircle />
                  Support kontaktieren
                </a>
              </div>
            </div>
          </motion.div>
        )}
      </motion.div>
    </div>
  );
};

export default SubscriptionSuccess;
