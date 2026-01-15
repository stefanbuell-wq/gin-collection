import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { adminApi } from '../api';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from 'recharts';
import {
  Building2,
  Users,
  Wine,
  TrendingUp,
  Crown,
  Clock,
  AlertTriangle,
  CheckCircle2,
  Loader2,
  RefreshCw
} from 'lucide-react';
import './Dashboard.css';

interface Stats {
  total_tenants: number;
  active_tenants: number;
  suspended_tenants: number;
  cancelled_tenants: number;
  total_users: number;
  total_gins: number;
  tenants_by_tier: Record<string, number>;
  new_tenants_last_7d: number;
  new_tenants_last_30d: number;
}

export default function Dashboard() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await adminApi.getStats();
      setStats(response.data);
    } catch (err) {
      setError('Statistiken konnten nicht geladen werden');
    } finally {
      setLoading(false);
    }
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.08
      }
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

  const statCards = [
    {
      label: 'Gesamt Tenants',
      value: stats?.total_tenants || 0,
      icon: Building2,
      color: 'gold',
      description: 'Registrierte Organisationen'
    },
    {
      label: 'Aktive Tenants',
      value: stats?.active_tenants || 0,
      icon: CheckCircle2,
      color: 'mint',
      description: 'Aktive Abonnements'
    },
    {
      label: 'Nutzer',
      value: stats?.total_users || 0,
      icon: Users,
      color: 'purple',
      description: 'Registrierte Benutzer'
    },
    {
      label: 'Gins',
      value: stats?.total_gins || 0,
      icon: Wine,
      color: 'copper',
      description: 'Erfasste Flaschen'
    },
    {
      label: 'Neu (7 Tage)',
      value: stats?.new_tenants_last_7d || 0,
      icon: TrendingUp,
      color: 'green',
      description: 'Neue Registrierungen'
    },
    {
      label: 'Suspendiert',
      value: stats?.suspended_tenants || 0,
      icon: AlertTriangle,
      color: 'red',
      description: 'Pausierte Konten'
    }
  ];

  // Tier chart data
  const tierData = stats?.tenants_by_tier
    ? Object.entries(stats.tenants_by_tier).map(([name, value]) => ({
        name: name.charAt(0).toUpperCase() + name.slice(1),
        value
      }))
    : [];

  const TIER_COLORS: Record<string, string> = {
    Free: '#6B6B63',
    Basic: '#4ECDC4',
    Pro: '#A855F7',
    Enterprise: '#D4A857'
  };

  // Mock growth data (would come from API in real app)
  const growthData = [
    { name: 'Mo', tenants: 12, users: 45 },
    { name: 'Di', tenants: 15, users: 52 },
    { name: 'Mi', tenants: 18, users: 58 },
    { name: 'Do', tenants: 14, users: 62 },
    { name: 'Fr', tenants: 22, users: 78 },
    { name: 'Sa', tenants: 8, users: 45 },
    { name: 'So', tenants: 5, users: 38 }
  ];

  if (loading) {
    return (
      <div className="admin-dashboard">
        <div className="admin-dashboard-loader">
          <Loader2 className="admin-dashboard-loader__icon" />
          <span>Lade Statistiken...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="admin-dashboard">
        <div className="admin-dashboard-error">
          <AlertTriangle />
          <span>{error}</span>
          <button onClick={loadStats} className="admin-dashboard-error__retry">
            <RefreshCw />
            Erneut versuchen
          </button>
        </div>
      </div>
    );
  }

  return (
    <motion.div
      className="admin-dashboard"
      variants={containerVariants}
      initial="hidden"
      animate="visible"
    >
      {/* Header */}
      <motion.div className="admin-dashboard-header" variants={itemVariants}>
        <div className="admin-dashboard-header__content">
          <h1 className="admin-dashboard-title">Dashboard</h1>
          <p className="admin-dashboard-subtitle">
            Platform-Übersicht und Statistiken
          </p>
        </div>
        <motion.button
          className="admin-dashboard-refresh"
          onClick={loadStats}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          <RefreshCw />
          Aktualisieren
        </motion.button>
      </motion.div>

      {/* Stats Grid */}
      <motion.div className="admin-dashboard-stats" variants={itemVariants}>
        {statCards.map((stat) => {
          const Icon = stat.icon;
          return (
            <motion.div
              key={stat.label}
              className={`admin-dashboard-stat admin-dashboard-stat--${stat.color}`}
              variants={itemVariants}
              whileHover={{ y: -4, transition: { duration: 0.2 } }}
            >
              <div className="admin-dashboard-stat__icon">
                <Icon />
              </div>
              <div className="admin-dashboard-stat__content">
                <span className="admin-dashboard-stat__value">{stat.value.toLocaleString()}</span>
                <span className="admin-dashboard-stat__label">{stat.label}</span>
                <span className="admin-dashboard-stat__desc">{stat.description}</span>
              </div>
            </motion.div>
          );
        })}
      </motion.div>

      {/* Charts Row */}
      <div className="admin-dashboard-charts">
        {/* Growth Chart */}
        <motion.div className="admin-dashboard-chart" variants={itemVariants}>
          <div className="admin-dashboard-chart__header">
            <h2 className="admin-dashboard-chart__title">Wachstum (7 Tage)</h2>
            <div className="admin-dashboard-chart__legend">
              <span className="admin-dashboard-chart__legend-item admin-dashboard-chart__legend-item--gold">
                <span />
                Tenants
              </span>
              <span className="admin-dashboard-chart__legend-item admin-dashboard-chart__legend-item--mint">
                <span />
                Users
              </span>
            </div>
          </div>
          <div className="admin-dashboard-chart__container">
            <ResponsiveContainer width="100%" height={280}>
              <AreaChart data={growthData}>
                <defs>
                  <linearGradient id="colorTenants" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#D4A857" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#D4A857" stopOpacity={0} />
                  </linearGradient>
                  <linearGradient id="colorUsers" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#4ECDC4" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#4ECDC4" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(212, 168, 87, 0.1)" />
                <XAxis
                  dataKey="name"
                  stroke="#6B6B63"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                />
                <YAxis
                  stroke="#6B6B63"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#0a0f0d',
                    border: '1px solid rgba(212, 168, 87, 0.2)',
                    borderRadius: '12px',
                    boxShadow: '0 8px 30px rgba(0, 0, 0, 0.3)'
                  }}
                  labelStyle={{ color: '#F5F5F0', marginBottom: '8px' }}
                  itemStyle={{ color: '#A8A8A0' }}
                />
                <Area
                  type="monotone"
                  dataKey="tenants"
                  stroke="#D4A857"
                  strokeWidth={2}
                  fillOpacity={1}
                  fill="url(#colorTenants)"
                />
                <Area
                  type="monotone"
                  dataKey="users"
                  stroke="#4ECDC4"
                  strokeWidth={2}
                  fillOpacity={1}
                  fill="url(#colorUsers)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </motion.div>

        {/* Tier Distribution */}
        <motion.div className="admin-dashboard-chart admin-dashboard-chart--small" variants={itemVariants}>
          <div className="admin-dashboard-chart__header">
            <h2 className="admin-dashboard-chart__title">Tier-Verteilung</h2>
          </div>
          <div className="admin-dashboard-chart__container admin-dashboard-chart__container--pie">
            <ResponsiveContainer width="100%" height={200}>
              <PieChart>
                <Pie
                  data={tierData}
                  cx="50%"
                  cy="50%"
                  innerRadius={50}
                  outerRadius={80}
                  paddingAngle={4}
                  dataKey="value"
                >
                  {tierData.map((entry, idx) => (
                    <Cell
                      key={`cell-${idx}`}
                      fill={TIER_COLORS[entry.name] || '#6B6B63'}
                      stroke="none"
                    />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#0a0f0d',
                    border: '1px solid rgba(212, 168, 87, 0.2)',
                    borderRadius: '12px'
                  }}
                  labelStyle={{ color: '#F5F5F0' }}
                />
              </PieChart>
            </ResponsiveContainer>
            <div className="admin-dashboard-tier-legend">
              {tierData.map((tier) => (
                <div key={tier.name} className="admin-dashboard-tier-legend__item">
                  <span
                    className="admin-dashboard-tier-legend__color"
                    style={{ backgroundColor: TIER_COLORS[tier.name] || '#6B6B63' }}
                  />
                  <span className="admin-dashboard-tier-legend__name">{tier.name}</span>
                  <span className="admin-dashboard-tier-legend__value">{tier.value}</span>
                </div>
              ))}
            </div>
          </div>
        </motion.div>
      </div>

      {/* Quick Stats */}
      <motion.div className="admin-dashboard-quick" variants={itemVariants}>
        <div className="admin-dashboard-quick__item">
          <Crown className="admin-dashboard-quick__icon" />
          <div className="admin-dashboard-quick__content">
            <span className="admin-dashboard-quick__value">
              {stats?.tenants_by_tier?.enterprise || 0}
            </span>
            <span className="admin-dashboard-quick__label">Enterprise Kunden</span>
          </div>
        </div>
        <div className="admin-dashboard-quick__divider" />
        <div className="admin-dashboard-quick__item">
          <Clock className="admin-dashboard-quick__icon" />
          <div className="admin-dashboard-quick__content">
            <span className="admin-dashboard-quick__value">
              {stats?.new_tenants_last_30d || 0}
            </span>
            <span className="admin-dashboard-quick__label">Neu (30 Tage)</span>
          </div>
        </div>
        <div className="admin-dashboard-quick__divider" />
        <div className="admin-dashboard-quick__item">
          <Wine className="admin-dashboard-quick__icon" />
          <div className="admin-dashboard-quick__content">
            <span className="admin-dashboard-quick__value">
              {stats?.total_users ? Math.round(stats.total_gins / stats.total_users) : 0}
            </span>
            <span className="admin-dashboard-quick__label">Gins/User</span>
          </div>
        </div>
      </motion.div>
    </motion.div>
  );
}
