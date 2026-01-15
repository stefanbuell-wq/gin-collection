import { useEffect, useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useGinStore } from '../stores/ginStore';
import {
  Wine,
  Star,
  Plus,
  Search,
  Heart,
  Droplets,
  TrendingUp,
  Award,
  Package
} from 'lucide-react';
import './GinList.css';

const GinList = () => {
  const { gins, total, fetchGins, isLoading } = useGinStore();
  const [searchQuery, setSearchQuery] = useState('');
  const [filter, setFilter] = useState<'all' | 'available' | 'favorite'>('all');

  useEffect(() => {
    fetchGins({ filter, limit: 50 });
  }, [filter]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      fetchGins({ q: searchQuery, limit: 50 });
    } else {
      fetchGins({ filter, limit: 50 });
    }
  };

  // Calculate collection stats
  const stats = useMemo(() => {
    const totalValue = gins.reduce((sum, gin) => sum + (gin.price || 0), 0);
    const avgRating = gins.length > 0
      ? gins.reduce((sum, gin) => sum + (gin.rating || 0), 0) / gins.filter(g => g.rating).length
      : 0;
    const availableCount = gins.filter(g => (g.fill_level || 0) > 0).length;
    const favoriteCount = gins.filter(g => g.is_favorite).length;

    return { totalValue, avgRating, availableCount, favoriteCount };
  }, [gins]);

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.08
      }
    }
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

  const statVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        type: 'spring',
        stiffness: 100,
        damping: 12
      }
    }
  };

  // Render star rating
  const renderStars = (rating: number | undefined) => {
    const stars = [];
    const ratingValue = rating || 0;
    for (let i = 1; i <= 5; i++) {
      stars.push(
        <Star
          key={i}
          className={`gin-list-card__star ${i <= ratingValue ? 'gin-list-card__star--filled' : ''}`}
        />
      );
    }
    return stars;
  };

  return (
    <div className="gin-list-page">
      {/* Ambient Background */}
      <div className="gin-list-ambient">
        <div className="gin-list-orb gin-list-orb--gold" />
        <div className="gin-list-orb gin-list-orb--mint" />
        <div className="gin-list-orb gin-list-orb--green" />
      </div>

      <div className="gin-list-content">
        {/* Header */}
        <motion.div
          className="gin-list-header"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          <div className="gin-list-header__info">
            <h1 className="gin-list-header__title">
              <div className="gin-list-header__icon">
                <Wine />
              </div>
              Meine Sammlung
            </h1>
            <p className="gin-list-header__subtitle">
              <span className="gin-list-header__count">{total}</span> Gins in deinem Tresor
            </p>
          </div>

          <motion.div
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <Link to="/gins/new" className="gin-list-add-btn">
              <Plus />
              Gin hinzufügen
            </Link>
          </motion.div>
        </motion.div>

        {/* Stats Bar */}
        {gins.length > 0 && (
          <motion.div
            className="gin-list-stats"
            initial="hidden"
            animate="visible"
            variants={containerVariants}
          >
            <motion.div className="gin-list-stat" variants={statVariants}>
              <div className="gin-list-stat__icon">
                <Package />
              </div>
              <div className="gin-list-stat__info">
                <span className="gin-list-stat__value">{total}</span>
                <span className="gin-list-stat__label">Flaschen</span>
              </div>
            </motion.div>

            <motion.div className="gin-list-stat" variants={statVariants}>
              <div className="gin-list-stat__icon">
                <TrendingUp />
              </div>
              <div className="gin-list-stat__info">
                <span className="gin-list-stat__value">{stats.totalValue.toFixed(0)}€</span>
                <span className="gin-list-stat__label">Gesamtwert</span>
              </div>
            </motion.div>

            <motion.div className="gin-list-stat" variants={statVariants}>
              <div className="gin-list-stat__icon gin-list-stat__icon--mint">
                <Droplets />
              </div>
              <div className="gin-list-stat__info">
                <span className="gin-list-stat__value">{stats.availableCount}</span>
                <span className="gin-list-stat__label">Verfügbar</span>
              </div>
            </motion.div>

            <motion.div className="gin-list-stat" variants={statVariants}>
              <div className="gin-list-stat__icon">
                <Award />
              </div>
              <div className="gin-list-stat__info">
                <span className="gin-list-stat__value">{stats.avgRating ? stats.avgRating.toFixed(1) : '-'}</span>
                <span className="gin-list-stat__label">Ø Bewertung</span>
              </div>
            </motion.div>
          </motion.div>
        )}

        {/* Search & Filters */}
        <motion.div
          className="gin-list-filters"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <div className="gin-list-filters__row">
            <form onSubmit={handleSearch} className="gin-list-search">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Gins durchsuchen..."
                className="gin-list-search__input"
              />
              <div className="gin-list-search__icon">
                <Search />
              </div>
            </form>

            <div className="gin-list-filter-buttons">
              <button
                onClick={() => setFilter('all')}
                className={`gin-list-filter-btn ${filter === 'all' ? 'gin-list-filter-btn--active' : ''}`}
              >
                Alle
              </button>
              <button
                onClick={() => setFilter('available')}
                className={`gin-list-filter-btn ${filter === 'available' ? 'gin-list-filter-btn--active' : ''}`}
              >
                Verfügbar
              </button>
              <button
                onClick={() => setFilter('favorite')}
                className={`gin-list-filter-btn ${filter === 'favorite' ? 'gin-list-filter-btn--active' : ''}`}
              >
                Favoriten
              </button>
            </div>
          </div>
        </motion.div>

        {/* Gin Grid */}
        {isLoading ? (
          <div className="gin-list-loading">
            <div className="gin-list-loading__spinner" />
            <p className="gin-list-loading__text">Lade deine Sammlung...</p>
          </div>
        ) : gins.length > 0 ? (
          <>
            <p className="gin-list-results">
              <span className="gin-list-results__count">{gins.length}</span> Ergebnisse
            </p>
            <motion.div
              className="gin-list-grid"
              initial="hidden"
              animate="visible"
              variants={containerVariants}
            >
              <AnimatePresence mode="popLayout">
                {gins.map((gin, index) => (
                  <motion.div
                    key={gin.id}
                    variants={cardVariants}
                    layout
                    whileHover={{ y: -8 }}
                    transition={{ type: 'spring', stiffness: 300, damping: 20 }}
                  >
                    <Link to={`/gins/${gin.id}`} className="gin-list-card">
                      <div className="gin-list-card__image-container">
                        {gin.primary_photo_url ? (
                          <img
                            src={gin.primary_photo_url}
                            alt={gin.name}
                            className="gin-list-card__image"
                          />
                        ) : (
                          <div className="gin-list-card__placeholder">
                            <Wine />
                          </div>
                        )}

                        {/* Fill Level Badge */}
                        {gin.fill_level !== undefined && gin.fill_level > 0 && (
                          <div className="gin-list-card__fill-badge">
                            <Droplets />
                            {gin.fill_level}%
                          </div>
                        )}

                        {/* Favorite Badge */}
                        {gin.is_favorite && (
                          <div className="gin-list-card__favorite-badge">
                            <Heart />
                          </div>
                        )}
                      </div>

                      <div className="gin-list-card__content">
                        <h3 className="gin-list-card__name">{gin.name}</h3>
                        <p className="gin-list-card__brand">
                          {gin.brand}
                          {gin.country && (
                            <>
                              <span className="gin-list-card__brand-separator" />
                              {gin.country}
                            </>
                          )}
                        </p>

                        <div className="gin-list-card__footer">
                          <div className="gin-list-card__rating">
                            {gin.rating ? (
                              <>
                                <div className="gin-list-card__stars">
                                  {renderStars(gin.rating)}
                                </div>
                                <span className="gin-list-card__rating-value">{gin.rating}/5</span>
                              </>
                            ) : (
                              <span className="gin-list-card__rating-value" style={{ color: 'var(--text-muted)' }}>
                                Keine Bewertung
                              </span>
                            )}
                          </div>

                          {gin.abv && (
                            <span className="gin-list-card__abv">{gin.abv}%</span>
                          )}

                          {!gin.abv && gin.price && (
                            <span className="gin-list-card__price">{gin.price}€</span>
                          )}
                        </div>
                      </div>
                    </Link>
                  </motion.div>
                ))}
              </AnimatePresence>
            </motion.div>
          </>
        ) : (
          <motion.div
            className="gin-list-empty"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5 }}
          >
            <div className="gin-list-empty__icon">
              <Wine />
            </div>
            <h3 className="gin-list-empty__title">Noch keine Gins</h3>
            <p className="gin-list-empty__text">
              Starte deine Sammlung und füge deinen ersten Gin hinzu
            </p>
            <motion.div
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Link to="/gins/new" className="gin-list-empty__btn">
                <Plus />
                Ersten Gin hinzufügen
              </Link>
            </motion.div>
          </motion.div>
        )}
      </div>
    </div>
  );
};

export default GinList;
