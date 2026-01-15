import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  XCircle,
  ArrowLeft,
  HelpCircle,
  Wine,
  Info,
  CreditCard,
  Shield
} from 'lucide-react';
import './SubscriptionStatus.css';

const SubscriptionCancel = () => {
  const navigate = useNavigate();

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
    hidden: { scale: 0, rotate: -45 },
    visible: {
      scale: 1,
      rotate: 0,
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

        <motion.div
          className="subscription-status-icon subscription-status-icon--cancel"
          variants={iconVariants}
          initial="hidden"
          animate="visible"
        >
          <XCircle />
        </motion.div>

        <h1 className="subscription-status-title">Zahlung abgebrochen</h1>

        <p className="subscription-status-description">
          Du hast den Zahlungsvorgang abgebrochen. Dein aktueller Plan bleibt unverändert.
        </p>

        <div className="subscription-status-warning">
          <p className="subscription-status-warning__text">
            <Info />
            Es wurden keine Zahlungen vorgenommen und dein Konto wurde nicht belastet.
          </p>
        </div>

        <div className="subscription-status-features">
          <p className="subscription-status-features__title">Bereit für ein Upgrade?</p>
          <ul className="subscription-status-features__list">
            <li className="subscription-status-features__item">
              <CreditCard />
              Sichere Zahlung über PayPal
            </li>
            <li className="subscription-status-features__item">
              <Shield />
              Jederzeit kündbar
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
            <ArrowLeft />
            Zurück zur Planauswahl
          </motion.button>

          <motion.button
            onClick={() => navigate('/')}
            className="subscription-status-btn subscription-status-btn--secondary"
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            Zum Dashboard
          </motion.button>
        </div>

        <div className="subscription-status-divider">
          <div className="subscription-status-support">
            <p className="subscription-status-support__text">
              Hast du Fragen oder Probleme?
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
    </div>
  );
};

export default SubscriptionCancel;
