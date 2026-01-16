import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useGinStore } from '../stores/ginStore';
import { useAuthStore } from '../stores/authStore';
import {
  Wine,
  Star,
  TrendingUp,
  Plus,
  Search,
  Bell,
  Settings,
  ChevronRight,
  Sparkles,
  Globe,
  Filter,
  Scan,
  Droplets,
} from 'lucide-react';
import './Dashboard.css';

// Types for the component
interface GinCardData {
  id: number;
  uuid?: string;
  name: string;
  brand: string;
  country: string;
  abv: number;
  fill_level: number;
  rating: number;
  price: number;
  gin_type: string;
  primary_photo_url?: string;
}

// Stat Card Component
const StatCard = ({
  icon: Icon,
  label,
  value,
  suffix = '',
  prefix = '',
  delay = 0,
  accent = false
}: {
  icon: React.ElementType;
  label: string;
  value: string | number;
  suffix?: string;
  prefix?: string;
  delay?: number;
  accent?: boolean;
}) => (
  <motion.div
    className={`stat-card ${accent ? 'stat-card--accent' : ''}`}
    initial={{ opacity: 0, y: 30 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ duration: 0.6, delay, ease: [0.22, 1, 0.36, 1] }}
  >
    <div className="stat-card__icon">
      <Icon size={22} />
    </div>
    <div className="stat-card__content">
      <span className="stat-card__label">{label}</span>
      <span className="stat-card__value">
        {prefix}<span className="stat-card__number">{value}</span>{suffix}
      </span>
    </div>
    <div className="stat-card__glow" />
  </motion.div>
);

// Fill Level Indicator Component
const FillLevel = ({ level }: { level: number }) => {
  const getColor = () => {
    if (level > 70) return 'var(--mint)';
    if (level > 30) return 'var(--gold)';
    return 'var(--copper)';
  };

  return (
    <div className="fill-level">
      <div className="fill-level__track">
        <motion.div
          className="fill-level__bar"
          initial={{ height: 0 }}
          animate={{ height: `${level}%` }}
          transition={{ duration: 1, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
          style={{ backgroundColor: getColor() }}
        />
      </div>
      <span className="fill-level__text" style={{ color: getColor() }}>{level}%</span>
    </div>
  );
};

// Rating Stars Component
const RatingStars = ({ rating }: { rating: number }) => (
  <div className="rating-stars">
    {[...Array(5)].map((_, i) => (
      <Star
        key={i}
        size={12}
        className={i < rating ? 'star--filled' : 'star--empty'}
        fill={i < rating ? 'var(--gold)' : 'none'}
        stroke={i < rating ? 'var(--gold)' : 'var(--glass-border)'}
      />
    ))}
  </div>
);

// Gin Card Component
const GinCard = ({ gin, index }: { gin: GinCardData; index: number }) => {
  const navigate = useNavigate();

  return (
    <motion.article
      className="gin-card"
      initial={{ opacity: 0, y: 40 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        duration: 0.6,
        delay: 0.1 * index + 0.4,
        ease: [0.22, 1, 0.36, 1]
      }}
      whileHover={{ y: -8, transition: { duration: 0.3 } }}
      onClick={() => navigate(`/gins/${gin.id}`)}
    >
      <div className="gin-card__image">
        {gin.primary_photo_url ? (
          <img
            src={gin.primary_photo_url}
            alt={gin.name}
            style={{ width: '100%', height: '100%', objectFit: 'cover' }}
          />
        ) : (
          <div className="gin-card__bottle">
            <Wine size={48} strokeWidth={1} />
          </div>
        )}
        <div className="gin-card__type-badge">{gin.gin_type || 'Gin'}</div>
      </div>

      <div className="gin-card__content">
        <div className="gin-card__header">
          <h3 className="gin-card__name">{gin.name}</h3>
          <span className="gin-card__brand">{gin.brand}</span>
        </div>

        <div className="gin-card__details">
          <div className="gin-card__meta">
            <span className="gin-card__country">
              <Globe size={12} />
              {gin.country || 'Unbekannt'}
            </span>
            <span className="gin-card__abv">{gin.abv}% vol</span>
          </div>
          <RatingStars rating={gin.rating || 0} />
        </div>

        <div className="gin-card__footer">
          <span className="gin-card__price">€{(gin.price || 0).toFixed(2)}</span>
          <FillLevel level={gin.fill_level || 100} />
        </div>
      </div>

      <div className="gin-card__shine" />
    </motion.article>
  );
};

// Scanner Button Component
const ScannerButton = () => (
  <motion.button
    className="scanner-btn"
    whileHover={{ scale: 1.05 }}
    whileTap={{ scale: 0.95 }}
    initial={{ opacity: 0, scale: 0.8 }}
    animate={{ opacity: 1, scale: 1 }}
    transition={{ duration: 0.5, delay: 1 }}
    title="Barcode scannen"
  >
    <span className="scanner-btn__ring" />
    <span className="scanner-btn__ring scanner-btn__ring--delayed" />
    <Scan size={28} />
  </motion.button>
);

// Main Dashboard Component
const Dashboard = () => {
  const navigate = useNavigate();
  const { stats, fetchStats, gins, fetchGins, isLoading } = useGinStore();
  const { user } = useAuthStore();
  const [searchQuery, setSearchQuery] = useState('');
  const [displayGins, setDisplayGins] = useState<GinCardData[]>([]);

  useEffect(() => {
    fetchStats();
    fetchGins({ limit: 12, sort: 'created_at' });
  }, [fetchStats, fetchGins]);

  useEffect(() => {
    // Transform gins to match our card data format
    const transformed = gins.map(gin => ({
      id: gin.id,
      uuid: gin.uuid,
      name: gin.name,
      brand: gin.brand || '',
      country: gin.country || '',
      abv: gin.abv || 0,
      fill_level: gin.fill_level || 100,
      rating: gin.rating || 0,
      price: gin.price || 0,
      gin_type: gin.gin_type || 'Gin',
      primary_photo_url: gin.primary_photo_url,
    }));
    setDisplayGins(transformed);
  }, [gins]);

  // Filter gins based on search
  const filteredGins = displayGins.filter(gin =>
    gin.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    gin.brand.toLowerCase().includes(searchQuery.toLowerCase()) ||
    gin.country.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // Get user initials for avatar
  const getUserInitials = () => {
    if (user?.first_name && user?.last_name) {
      return `${user.first_name[0]}${user.last_name[0]}`.toUpperCase();
    }
    return user?.email?.substring(0, 2).toUpperCase() || 'GV';
  };

  return (
    <div className="dashboard">
      {/* Ambient Background */}
      <div className="dashboard__ambient">
        <div className="ambient-orb ambient-orb--1" />
        <div className="ambient-orb ambient-orb--2" />
        <div className="ambient-orb ambient-orb--3" />
      </div>

      {/* Header */}
      <motion.header
        className="dashboard__header"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        <div className="header__brand">
          <div className="header__logo">
            <Sparkles size={24} />
          </div>
          <div className="header__title">
            <h1>GinVault</h1>
            <span>Premium Collection</span>
          </div>
        </div>

        <div className="header__actions">
          <button className="header__icon-btn">
            <Bell size={20} />
            <span className="header__notification-dot" />
          </button>
          <button
            className="header__icon-btn"
            onClick={() => navigate('/settings')}
          >
            <Settings size={20} />
          </button>
          <div className="header__avatar">
            <span>{getUserInitials()}</span>
          </div>
        </div>
      </motion.header>

      {/* Stats Section */}
      <section className="dashboard__stats">
        <StatCard
          icon={TrendingUp}
          label="Gesamtwert"
          value={stats?.total_value?.toLocaleString('de-DE', { minimumFractionDigits: 2 }) || '0,00'}
          prefix="€"
          delay={0.1}
          accent
        />
        <StatCard
          icon={Wine}
          label="Flaschen"
          value={stats?.total_gins || 0}
          delay={0.2}
        />
        <StatCard
          icon={Star}
          label="Favoriten"
          value={stats?.favorite_count || 0}
          delay={0.3}
        />
        <StatCard
          icon={Globe}
          label="Länder"
          value={stats?.by_country ? Object.keys(stats.by_country).length : 0}
          delay={0.4}
        />
      </section>

      {/* Search & Filter Bar */}
      <motion.section
        className="dashboard__toolbar"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.6, delay: 0.3 }}
      >
        <div className="search-bar">
          <Search size={18} className="search-bar__icon" />
          <input
            type="text"
            placeholder="Suche nach Gin, Marke oder Land..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-bar__input"
          />
        </div>
        <button className="filter-btn">
          <Filter size={18} />
          <span>Filter</span>
        </button>
        <Link to="/gins/new" className="add-btn">
          <Plus size={18} />
          <span>Hinzufügen</span>
        </Link>
      </motion.section>

      {/* Collection Header */}
      <motion.div
        className="dashboard__section-header"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.6, delay: 0.4 }}
      >
        <h2>Meine Sammlung</h2>
        <Link to="/gins" className="view-all-btn">
          Alle anzeigen
          <ChevronRight size={16} />
        </Link>
      </motion.div>

      {/* Gin Grid */}
      <section className="dashboard__collection">
        <AnimatePresence mode="wait">
          {isLoading ? (
            <motion.div
              className="loading-state"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            >
              <div className="loading-spinner" />
              <span>Lade Sammlung...</span>
            </motion.div>
          ) : filteredGins.length > 0 ? (
            <motion.div className="gin-grid">
              {filteredGins.slice(0, 8).map((gin, index) => (
                <GinCard key={gin.id} gin={gin} index={index} />
              ))}
            </motion.div>
          ) : (
            <motion.div
              className="loading-state"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
            >
              <Wine size={48} style={{ opacity: 0.3 }} />
              <span>
                {searchQuery
                  ? 'Keine Gins gefunden'
                  : 'Noch keine Gins in der Sammlung'}
              </span>
              {!searchQuery && (
                <Link to="/gins/new" className="add-btn" style={{ marginTop: '16px' }}>
                  <Plus size={18} />
                  <span>Ersten Gin hinzufügen</span>
                </Link>
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </section>

      {/* Floating Scanner Button */}
      <ScannerButton />

      {/* Bottom Navigation */}
      <motion.nav
        className="dashboard__nav"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.8 }}
      >
        <Link to="/" className="nav-item nav-item--active">
          <Wine size={22} />
          <span>Sammlung</span>
        </Link>
        <Link to="/cocktails" className="nav-item">
          <Droplets size={22} />
          <span>Cocktails</span>
        </Link>
        <Link to="/statistics" className="nav-item">
          <TrendingUp size={22} />
          <span>Statistik</span>
        </Link>
        <Link to="/settings" className="nav-item">
          <Settings size={22} />
          <span>Einstellungen</span>
        </Link>
      </motion.nav>
    </div>
  );
};

export default Dashboard;
