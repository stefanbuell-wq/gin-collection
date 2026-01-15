import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import { subscriptionAPI, tenantAPI } from '../api/services';
import type { TenantTier, BillingCycle, SubscriptionPlan, Subscription as SubscriptionType } from '../types';
import {
  Crown,
  Check,
  X,
  CreditCard,
  Calendar,
  Zap,
  Building2,
  Users,
  HardDrive,
  Image,
  Star,
  AlertTriangle,
  ExternalLink,
  Sparkles
} from 'lucide-react';
import './Subscription.css';

// Plan definitions
const PLANS: SubscriptionPlan[] = [
  {
    id: 'free',
    name: 'Free',
    tier: 'free',
    description: 'Perfekt zum Ausprobieren',
    price_monthly: 0,
    price_yearly: 0,
    features: [
      'Bis zu 25 Gins',
      '3 Fotos pro Gin',
      'Basis-Statistiken',
      'Community Support'
    ],
    limits: {
      max_gins: 25,
      max_photos_per_gin: 3,
      features: ['basic_stats']
    }
  },
  {
    id: 'basic',
    name: 'Basic',
    tier: 'basic',
    description: 'Fur Hobby-Sammler',
    price_monthly: 4.99,
    price_yearly: 49.99,
    features: [
      'Bis zu 100 Gins',
      '10 Fotos pro Gin',
      'Erweiterte Statistiken',
      'Export-Funktion',
      'E-Mail Support'
    ],
    limits: {
      max_gins: 100,
      max_photos_per_gin: 10,
      features: ['basic_stats', 'advanced_stats', 'export']
    }
  },
  {
    id: 'pro',
    name: 'Pro',
    tier: 'pro',
    description: 'Fur ernsthafte Sammler',
    price_monthly: 9.99,
    price_yearly: 99.99,
    features: [
      'Bis zu 500 Gins',
      '25 Fotos pro Gin',
      'Alle Statistiken',
      'API-Zugang',
      'Priority Support',
      'Barcode-Scanner'
    ],
    limits: {
      max_gins: 500,
      max_photos_per_gin: 25,
      features: ['basic_stats', 'advanced_stats', 'export', 'api_access', 'barcode_scanner']
    }
  },
  {
    id: 'enterprise',
    name: 'Enterprise',
    tier: 'enterprise',
    description: 'Fur Bars & Handler',
    price_monthly: 29.99,
    price_yearly: 299.99,
    features: [
      'Unbegrenzte Gins',
      'Unbegrenzte Fotos',
      'Team-Verwaltung',
      'White-Label Option',
      'Dedizierter Support',
      'Custom Integrations',
      'SLA-Garantie'
    ],
    limits: {
      max_gins: -1,
      max_photos_per_gin: -1,
      features: ['basic_stats', 'advanced_stats', 'export', 'api_access', 'barcode_scanner', 'team_management', 'white_label']
    }
  }
];

const TIER_ICONS: Record<TenantTier, React.ReactNode> = {
  free: <Zap size={20} />,
  basic: <Star size={20} />,
  pro: <Crown size={20} />,
  enterprise: <Building2 size={20} />
};

const Subscription = () => {
  const { tenant, setTenant } = useAuthStore();
  const [billingCycle, setBillingCycle] = useState<BillingCycle>('monthly');
  const [selectedPlan, setSelectedPlan] = useState<string | null>(null);
  const [isUpgrading, setIsUpgrading] = useState(false);
  const [isCancelling, setIsCancelling] = useState(false);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [currentSubscription, setCurrentSubscription] = useState<SubscriptionType | null>(null);
  const [upgradeError, setUpgradeError] = useState('');
  const [cancelReason, setCancelReason] = useState('');

  const currentTier = tenant?.tier || 'free';
  const currentPlan = PLANS.find(p => p.tier === currentTier);

  useEffect(() => {
    const loadSubscription = async () => {
      try {
        const response = await subscriptionAPI.getCurrent();
        const apiResponse = response.data as unknown as { success: boolean; data: { subscription: SubscriptionType | null } };
        if (apiResponse.success && apiResponse.data?.subscription) {
          setCurrentSubscription(apiResponse.data.subscription);
        }
      } catch (err) {
        console.error('Failed to load subscription:', err);
      }
    };
    loadSubscription();
  }, []);

  const handleUpgrade = async (planId: string) => {
    setSelectedPlan(planId);
    setUpgradeError('');
    setShowConfirmModal(true);
  };

  const confirmUpgrade = async () => {
    if (!selectedPlan) return;

    setIsUpgrading(true);
    setUpgradeError('');
    try {
      const backendPlanId = mapPlanIdToBackend(selectedPlan, billingCycle);
      const response = await subscriptionAPI.upgrade(backendPlanId, billingCycle);
      const apiResponse = response.data as unknown as {
        success: boolean;
        data: {
          paypal_approval_url: string;
          subscription_id: number;
        }
      };

      if (apiResponse.success && apiResponse.data?.paypal_approval_url) {
        window.location.href = apiResponse.data.paypal_approval_url;
      } else {
        setUpgradeError('Upgrade konnte nicht gestartet werden. Bitte versuche es erneut.');
      }
    } catch (err: unknown) {
      console.error('Upgrade failed:', err);
      const error = err as { response?: { data?: { error?: string } } };
      setUpgradeError(error.response?.data?.error || 'Upgrade fehlgeschlagen. Bitte versuche es erneut.');
    } finally {
      setIsUpgrading(false);
    }
  };

  const handleCancelSubscription = async () => {
    setIsCancelling(true);
    try {
      await subscriptionAPI.cancel();
      setShowCancelModal(false);

      const tenantResponse = await tenantAPI.getCurrent();
      const tenantData = tenantResponse.data as unknown as { success: boolean; data: { tenant: typeof tenant } };
      if (tenantData.success && tenantData.data?.tenant) {
        setTenant(tenantData.data.tenant);
      }

      setCurrentSubscription(null);
    } catch (err) {
      console.error('Cancel failed:', err);
    } finally {
      setIsCancelling(false);
    }
  };

  const mapPlanIdToBackend = (planId: string, cycle: BillingCycle): string => {
    const mapping: Record<string, Record<BillingCycle, string>> = {
      basic: { monthly: 'PLAN_BASIC_MONTHLY', yearly: 'PLAN_BASIC_YEARLY' },
      pro: { monthly: 'PLAN_PRO_MONTHLY', yearly: 'PLAN_PRO_YEARLY' },
      enterprise: { monthly: 'PLAN_ENTERPRISE', yearly: 'PLAN_ENTERPRISE' },
    };
    return mapping[planId]?.[cycle] || planId;
  };

  const getYearlySavings = (plan: SubscriptionPlan) => {
    const monthlyTotal = plan.price_monthly * 12;
    const savings = monthlyTotal - plan.price_yearly;
    return savings > 0 ? Math.round((savings / monthlyTotal) * 100) : 0;
  };

  const getPlanPrice = (plan: SubscriptionPlan) => {
    return billingCycle === 'yearly' ? plan.price_yearly : plan.price_monthly;
  };

  const isPlanUpgrade = (plan: SubscriptionPlan) => {
    const tierOrder: TenantTier[] = ['free', 'basic', 'pro', 'enterprise'];
    return tierOrder.indexOf(plan.tier) > tierOrder.indexOf(currentTier);
  };

  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 }
  };

  return (
    <div className="subscription-page">
      {/* Ambient Background */}
      <div className="subscription-page__ambient">
        <div className="subscription-ambient-orb subscription-ambient-orb--1" />
        <div className="subscription-ambient-orb subscription-ambient-orb--2" />
        <div className="subscription-ambient-orb subscription-ambient-orb--3" />
      </div>

      <div className="subscription-container">
        {/* Header */}
        <motion.div
          className="subscription-header"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <div className="subscription-header__content">
            <div className="subscription-header__icon">
              <CreditCard size={28} />
            </div>
            <div className="subscription-header__text">
              <h1>Abonnement</h1>
              <p>Verwalte dein Abonnement und wahle den passenden Plan</p>
            </div>
          </div>
          {currentTier !== 'enterprise' && (
            <motion.button
              className="subscription-upgrade-btn"
              onClick={() => document.getElementById('plans')?.scrollIntoView({ behavior: 'smooth' })}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Zap size={18} />
              Upgrade
            </motion.button>
          )}
        </motion.div>

        {/* Current Plan */}
        <motion.div
          className="subscription-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <div className="subscription-card__header">
            <div className="subscription-card__title">
              <div className="subscription-card__title-icon">
                <Crown size={20} />
              </div>
              <h2>Aktueller Plan</h2>
            </div>
          </div>

          <div className="subscription-current">
            <div className="subscription-current__info">
              <div className="subscription-current__tier">
                <span className={`subscription-tier-badge subscription-tier-badge--${currentTier}`}>
                  {TIER_ICONS[currentTier]}
                  {currentPlan?.name}
                </span>
              </div>
              <span className="subscription-current__price">
                {currentPlan?.price_monthly === 0 ? 'Kostenlos' : `${currentPlan?.price_monthly} EUR/Monat`}
              </span>
            </div>
          </div>

          <div className="subscription-stats-grid">
            <div className="subscription-stat-item">
              <div className="subscription-stat-item__icon">
                <HardDrive size={20} />
              </div>
              <div className="subscription-stat-item__content">
                <span className="subscription-stat-item__label">Max. Gins</span>
                <span className="subscription-stat-item__value">
                  {currentPlan?.limits.max_gins === -1 ? 'Unbegrenzt' : currentPlan?.limits.max_gins}
                </span>
              </div>
            </div>
            <div className="subscription-stat-item">
              <div className="subscription-stat-item__icon">
                <Image size={20} />
              </div>
              <div className="subscription-stat-item__content">
                <span className="subscription-stat-item__label">Fotos pro Gin</span>
                <span className="subscription-stat-item__value">
                  {currentPlan?.limits.max_photos_per_gin === -1 ? 'Unbegrenzt' : currentPlan?.limits.max_photos_per_gin}
                </span>
              </div>
            </div>
            <div className="subscription-stat-item">
              <div className="subscription-stat-item__icon">
                <Users size={20} />
              </div>
              <div className="subscription-stat-item__content">
                <span className="subscription-stat-item__label">Team-Mitglieder</span>
                <span className="subscription-stat-item__value">
                  {currentTier === 'enterprise' ? 'Unbegrenzt' : currentTier === 'pro' ? '5' : '1'}
                </span>
              </div>
            </div>
          </div>
        </motion.div>

        {/* Billing Information */}
        <motion.div
          className="subscription-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <div className="subscription-card__header">
            <div className="subscription-card__title">
              <div className="subscription-card__title-icon">
                <CreditCard size={20} />
              </div>
              <h2>Zahlungsinformationen</h2>
            </div>
          </div>

          <div className="subscription-billing-grid">
            <div className="subscription-billing-item">
              <span className="subscription-billing-item__label">Zahlungsmethode</span>
              {currentTier === 'free' ? (
                <div className="subscription-billing-box subscription-billing-box--empty">
                  <CreditCard size={32} />
                  <p>Keine Zahlungsmethode erforderlich</p>
                  <span>Upgrade um Zahlungsmethode hinzuzufugen</span>
                </div>
              ) : (
                <div className="subscription-billing-box">
                  <div className="subscription-billing-content">
                    <div className="subscription-billing-content__icon">
                      <CreditCard size={20} />
                    </div>
                    <div className="subscription-billing-content__info">
                      <span className="subscription-billing-content__title">PayPal</span>
                      <span className="subscription-billing-content__subtitle">Verbunden</span>
                    </div>
                    <button className="subscription-billing-content__action">Andern</button>
                  </div>
                </div>
              )}
            </div>

            <div className="subscription-billing-item">
              <span className="subscription-billing-item__label">Abrechnungszeitraum</span>
              {currentTier === 'free' ? (
                <div className="subscription-billing-box subscription-billing-box--empty">
                  <Calendar size={32} />
                  <p>Kein Abrechnungszeitraum</p>
                  <span>Kostenloser Plan ohne Abrechnung</span>
                </div>
              ) : (
                <div className="subscription-billing-box">
                  <div className="subscription-billing-content">
                    <div className="subscription-billing-content__icon">
                      <Calendar size={20} />
                    </div>
                    <div className="subscription-billing-content__info">
                      <span className="subscription-billing-content__title">
                        {currentSubscription?.billing_cycle === 'yearly' ? 'Jahrlich' : 'Monatlich'}
                      </span>
                      <span className="subscription-billing-content__subtitle">
                        {currentSubscription?.next_billing_date
                          ? `Nachste Abrechnung: ${new Date(currentSubscription.next_billing_date).toLocaleDateString('de-DE')}`
                          : 'Aktives Abonnement'
                        }
                      </span>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Invoice History */}
          {currentTier !== 'free' && (
            <div className="subscription-invoices">
              <span className="subscription-invoices__title">Rechnungshistorie</span>
              <table className="subscription-invoices-table">
                <thead>
                  <tr>
                    <th>Datum</th>
                    <th>Beschreibung</th>
                    <th>Betrag</th>
                    <th>Status</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>15. Jan 2026</td>
                    <td style={{ color: 'var(--text-secondary)' }}>{currentPlan?.name} Plan - Monatlich</td>
                    <td>{currentPlan?.price_monthly} EUR</td>
                    <td>
                      <span className="subscription-invoice-status subscription-invoice-status--paid">
                        Bezahlt
                      </span>
                    </td>
                    <td>
                      <button className="subscription-invoice-pdf">PDF</button>
                    </td>
                  </tr>
                  <tr>
                    <td>15. Dez 2025</td>
                    <td style={{ color: 'var(--text-secondary)' }}>{currentPlan?.name} Plan - Monatlich</td>
                    <td>{currentPlan?.price_monthly} EUR</td>
                    <td>
                      <span className="subscription-invoice-status subscription-invoice-status--paid">
                        Bezahlt
                      </span>
                    </td>
                    <td>
                      <button className="subscription-invoice-pdf">PDF</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          )}
        </motion.div>

        {/* Available Plans */}
        <motion.div
          id="plans"
          className="subscription-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.5, delay: 0.3 }}
        >
          <div className="subscription-plans-header">
            <div className="subscription-plans-header__info">
              <h2>Verfugbare Plane</h2>
              <p>Wahle den Plan, der am besten zu dir passt</p>
            </div>

            <div className="subscription-billing-toggle">
              <button
                onClick={() => setBillingCycle('monthly')}
                className={`subscription-billing-toggle__btn ${billingCycle === 'monthly' ? 'subscription-billing-toggle__btn--active' : ''}`}
              >
                Monatlich
              </button>
              <button
                onClick={() => setBillingCycle('yearly')}
                className={`subscription-billing-toggle__btn ${billingCycle === 'yearly' ? 'subscription-billing-toggle__btn--active' : ''}`}
              >
                Jahrlich
                <span className="subscription-billing-toggle__badge">-17%</span>
              </button>
            </div>
          </div>

          <div className="subscription-plans-grid">
            {PLANS.map((plan, index) => {
              const isCurrentPlan = plan.tier === currentTier;
              const isUpgrade = isPlanUpgrade(plan);
              const savings = getYearlySavings(plan);

              return (
                <motion.div
                  key={plan.id}
                  className={`subscription-plan-card ${isCurrentPlan ? 'subscription-plan-card--current' : ''} ${plan.tier === 'pro' && !isCurrentPlan ? 'subscription-plan-card--popular' : ''}`}
                  initial={{ opacity: 0, y: 30 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.4, delay: 0.4 + index * 0.1 }}
                >
                  {/* Badges */}
                  {plan.tier === 'pro' && !isCurrentPlan && (
                    <span className="subscription-plan-badge subscription-plan-badge--popular">
                      Beliebt
                    </span>
                  )}
                  {isCurrentPlan && (
                    <span className="subscription-plan-badge subscription-plan-badge--current">
                      Aktueller Plan
                    </span>
                  )}

                  {/* Plan Header */}
                  <div className="subscription-plan-header">
                    <div className={`subscription-plan-icon subscription-plan-icon--${plan.tier}`}>
                      {TIER_ICONS[plan.tier]}
                    </div>
                    <h3 className="subscription-plan-name">{plan.name}</h3>
                    <p className="subscription-plan-description">{plan.description}</p>
                  </div>

                  {/* Price */}
                  <div className="subscription-plan-price">
                    <div className="subscription-plan-price__amount">
                      <span className="subscription-plan-price__value">
                        {getPlanPrice(plan)} EUR
                      </span>
                      {plan.price_monthly > 0 && (
                        <span className="subscription-plan-price__period">
                          /{billingCycle === 'yearly' ? 'Jahr' : 'Monat'}
                        </span>
                      )}
                    </div>
                    {billingCycle === 'yearly' && savings > 0 && (
                      <p className="subscription-plan-price__savings">
                        Spare {savings}% gegenuber monatlich
                      </p>
                    )}
                  </div>

                  {/* Features */}
                  <ul className="subscription-plan-features">
                    {plan.features.map((feature, idx) => (
                      <li key={idx} className="subscription-plan-feature">
                        <Check size={16} />
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>

                  {/* Action Button */}
                  {isCurrentPlan ? (
                    <button className="subscription-plan-btn subscription-plan-btn--current" disabled>
                      Aktueller Plan
                    </button>
                  ) : isUpgrade ? (
                    <motion.button
                      onClick={() => handleUpgrade(plan.id)}
                      className={`subscription-plan-btn ${plan.tier === 'pro' ? 'subscription-plan-btn--popular' : 'subscription-plan-btn--upgrade'}`}
                      whileHover={{ scale: 1.02 }}
                      whileTap={{ scale: 0.98 }}
                    >
                      <Zap size={16} />
                      Upgrade
                    </motion.button>
                  ) : (
                    <button
                      onClick={() => handleUpgrade(plan.id)}
                      className="subscription-plan-btn subscription-plan-btn--downgrade"
                    >
                      Downgrade
                    </button>
                  )}
                </motion.div>
              );
            })}
          </div>
        </motion.div>

        {/* Feature Comparison */}
        <motion.div
          className="subscription-card"
          initial="hidden"
          animate="visible"
          variants={cardVariants}
          transition={{ duration: 0.5, delay: 0.5 }}
        >
          <div className="subscription-card__header">
            <div className="subscription-card__title">
              <div className="subscription-card__title-icon">
                <Sparkles size={20} />
              </div>
              <h2>Feature-Vergleich</h2>
            </div>
          </div>

          <div className="subscription-comparison">
            <table className="subscription-comparison-table">
              <thead>
                <tr>
                  <th>Feature</th>
                  {PLANS.map((plan) => (
                    <th
                      key={plan.id}
                      className={plan.tier === currentTier ? 'subscription-comparison-current' : ''}
                    >
                      {plan.name}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Max. Gins</td>
                  <td>25</td>
                  <td>100</td>
                  <td>500</td>
                  <td>Unbegrenzt</td>
                </tr>
                <tr>
                  <td>Fotos pro Gin</td>
                  <td>3</td>
                  <td>10</td>
                  <td>25</td>
                  <td>Unbegrenzt</td>
                </tr>
                <tr>
                  <td>Team-Mitglieder</td>
                  <td>1</td>
                  <td>1</td>
                  <td>5</td>
                  <td>Unbegrenzt</td>
                </tr>
                <tr>
                  <td>Statistiken</td>
                  <td>Basis</td>
                  <td>Erweitert</td>
                  <td>Alle</td>
                  <td>Alle + Custom</td>
                </tr>
                <tr>
                  <td>Export</td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                </tr>
                <tr>
                  <td>API-Zugang</td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                </tr>
                <tr>
                  <td>Barcode-Scanner</td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                </tr>
                <tr>
                  <td>Team-Verwaltung</td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                </tr>
                <tr>
                  <td>White-Label</td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><X size={16} className="subscription-comparison-x" /></td>
                  <td><Check size={16} className="subscription-comparison-check" /></td>
                </tr>
                <tr>
                  <td>Support</td>
                  <td>Community</td>
                  <td>E-Mail</td>
                  <td>Priority</td>
                  <td>Dediziert</td>
                </tr>
              </tbody>
            </table>
          </div>
        </motion.div>

        {/* Cancel Subscription */}
        {currentTier !== 'free' && (
          <motion.div
            className="subscription-card subscription-cancel-card"
            initial="hidden"
            animate="visible"
            variants={cardVariants}
            transition={{ duration: 0.5, delay: 0.6 }}
          >
            <div className="subscription-cancel-content">
              <div className="subscription-cancel-icon">
                <AlertTriangle size={22} />
              </div>
              <div className="subscription-cancel-info">
                <h3>Abonnement kundigen</h3>
                <p>
                  Wenn du dein Abonnement kundigst, wird es am Ende des aktuellen Abrechnungszeitraums beendet.
                  Du behaltst den Zugang bis dahin.
                </p>
                <button
                  className="subscription-cancel-btn"
                  onClick={() => setShowCancelModal(true)}
                >
                  Abonnement kundigen
                </button>
              </div>
            </div>
          </motion.div>
        )}
      </div>

      {/* Cancel Modal */}
      <AnimatePresence>
        {showCancelModal && (
          <motion.div
            className="subscription-modal-overlay"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <motion.div
              className="subscription-modal"
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
            >
              <div className="subscription-modal__header">
                <div className="subscription-modal__header-icon subscription-modal__header-icon--danger">
                  <AlertTriangle size={22} />
                </div>
                <h3>Abonnement kundigen</h3>
              </div>

              <div className="subscription-modal__body">
                <p>
                  Bist du sicher, dass du dein <strong>{currentPlan?.name}</strong> Abonnement kundigen mochtest?
                </p>

                <label style={{ display: 'block', marginBottom: '8px', fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                  Grund fur die Kundigung (optional)
                </label>
                <textarea
                  value={cancelReason}
                  onChange={(e) => setCancelReason(e.target.value)}
                  placeholder="Hilf uns, besser zu werden..."
                  className="subscription-modal__textarea"
                />

                <div className="subscription-modal__alert subscription-modal__alert--warning">
                  <AlertTriangle size={16} />
                  <span>
                    Nach der Kundigung wirst du zum kostenlosen Plan herabgestuft und verlierst Zugang zu Premium-Funktionen.
                  </span>
                </div>
              </div>

              <div className="subscription-modal__actions">
                <button
                  onClick={() => setShowCancelModal(false)}
                  className="subscription-modal__btn subscription-modal__btn--secondary"
                  disabled={isCancelling}
                >
                  Behalten
                </button>
                <button
                  onClick={handleCancelSubscription}
                  className="subscription-modal__btn subscription-modal__btn--danger"
                  disabled={isCancelling}
                >
                  {isCancelling ? (
                    <>
                      <span className="subscription-spinner" />
                      Wird gekundigt...
                    </>
                  ) : (
                    'Jetzt kundigen'
                  )}
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Upgrade Modal */}
      <AnimatePresence>
        {showConfirmModal && selectedPlan && (
          <motion.div
            className="subscription-modal-overlay"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            <motion.div
              className="subscription-modal"
              initial={{ opacity: 0, scale: 0.95, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 20 }}
            >
              <div className="subscription-modal__header">
                <div className="subscription-modal__header-icon subscription-modal__header-icon--info">
                  <Zap size={22} />
                </div>
                <h3>Plan-Wechsel bestatigen</h3>
              </div>

              <div className="subscription-modal__body">
                {upgradeError && (
                  <div className="subscription-error">
                    <AlertTriangle size={16} />
                    {upgradeError}
                  </div>
                )}

                {(() => {
                  const newPlan = PLANS.find(p => p.id === selectedPlan);
                  const isUpgrade = newPlan && isPlanUpgrade(newPlan);

                  return (
                    <>
                      <p>
                        Du wechselst von <strong>{currentPlan?.name}</strong> zu{' '}
                        <strong>{newPlan?.name}</strong>.
                      </p>

                      <div className="subscription-modal__price-box">
                        <div className="subscription-modal__price-row">
                          <span className="subscription-modal__price-label">Neuer Preis:</span>
                          <span className="subscription-modal__price-value">
                            {getPlanPrice(newPlan!)} EUR/{billingCycle === 'yearly' ? 'Jahr' : 'Monat'}
                          </span>
                        </div>
                        {isUpgrade && (
                          <p className="subscription-modal__price-note">
                            Die Differenz wird anteilig fur den aktuellen Zeitraum berechnet.
                          </p>
                        )}
                      </div>

                      <div className="subscription-modal__alert subscription-modal__alert--info">
                        <ExternalLink size={16} />
                        <span>Du wirst zu PayPal weitergeleitet, um die Zahlung abzuschliessen.</span>
                      </div>
                    </>
                  );
                })()}
              </div>

              <div className="subscription-modal__actions">
                <button
                  onClick={() => {
                    setShowConfirmModal(false);
                    setUpgradeError('');
                  }}
                  className="subscription-modal__btn subscription-modal__btn--secondary"
                  disabled={isUpgrading}
                >
                  Abbrechen
                </button>
                <button
                  onClick={confirmUpgrade}
                  className="subscription-modal__btn subscription-modal__btn--primary"
                  disabled={isUpgrading}
                >
                  {isUpgrading ? (
                    <>
                      <span className="subscription-spinner" />
                      Weiterleitung...
                    </>
                  ) : (
                    <>
                      <ExternalLink size={16} />
                      Zu PayPal
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

export default Subscription;
