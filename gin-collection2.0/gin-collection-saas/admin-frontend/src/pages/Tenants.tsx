import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { adminApi } from '../api';
import {
  Building2,
  Search,
  ChevronLeft,
  ChevronRight,
  Loader2,
  AlertTriangle,
  RefreshCw,
  Users,
  Wine,
  Calendar,
  MoreVertical,
  Play,
  Pause
} from 'lucide-react';
import './Tenants.css';

interface Tenant {
  tenant: {
    id: number;
    name: string;
    subdomain: string;
    tier: string;
    status: string;
    created_at: string;
  };
  user_count: number;
  gin_count: number;
}

export default function Tenants() {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [activeDropdown, setActiveDropdown] = useState<number | null>(null);
  const [actionLoading, setActionLoading] = useState<number | null>(null);

  useEffect(() => {
    loadTenants();
  }, [page]);

  const loadTenants = async () => {
    try {
      setLoading(true);
      const response = await adminApi.getTenants(page, 20);
      setTenants(response.data.tenants || []);
      setTotal(response.data.total);
    } catch (err) {
      setError('Tenants konnten nicht geladen werden');
    } finally {
      setLoading(false);
    }
  };

  const handleSuspend = async (id: number) => {
    if (!confirm('Möchtest du diesen Tenant wirklich suspendieren?')) return;
    setActionLoading(id);
    setActiveDropdown(null);
    try {
      await adminApi.suspendTenant(id);
      loadTenants();
    } catch (err) {
      alert('Tenant konnte nicht suspendiert werden');
    } finally {
      setActionLoading(null);
    }
  };

  const handleActivate = async (id: number) => {
    setActionLoading(id);
    setActiveDropdown(null);
    try {
      await adminApi.activateTenant(id);
      loadTenants();
    } catch (err) {
      alert('Tenant konnte nicht aktiviert werden');
    } finally {
      setActionLoading(null);
    }
  };

  const handleTierChange = async (id: number, tier: string) => {
    setActionLoading(id);
    try {
      await adminApi.updateTenantTier(id, tier);
      loadTenants();
    } catch (err) {
      alert('Tier konnte nicht geändert werden');
    } finally {
      setActionLoading(null);
    }
  };

  const getTierBadgeClass = (tier: string) => {
    switch (tier) {
      case 'free': return 'admin-tenants-badge--free';
      case 'basic': return 'admin-tenants-badge--basic';
      case 'pro': return 'admin-tenants-badge--pro';
      case 'enterprise': return 'admin-tenants-badge--enterprise';
      default: return 'admin-tenants-badge--free';
    }
  };

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'active': return 'admin-tenants-badge--active';
      case 'suspended': return 'admin-tenants-badge--suspended';
      case 'cancelled': return 'admin-tenants-badge--cancelled';
      default: return 'admin-tenants-badge--active';
    }
  };

  const filteredTenants = tenants.filter((item) =>
    item.tenant.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.tenant.subdomain.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { staggerChildren: 0.05 }
    }
  };

  const rowVariants = {
    hidden: { opacity: 0, x: -20 },
    visible: {
      opacity: 1,
      x: 0,
      transition: { type: 'spring', stiffness: 100, damping: 15 }
    }
  };

  if (loading && tenants.length === 0) {
    return (
      <div className="admin-tenants">
        <div className="admin-tenants-loader">
          <Loader2 className="admin-tenants-loader__icon" />
          <span>Lade Tenants...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-tenants">
      {/* Header */}
      <motion.div
        className="admin-tenants-header"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <div className="admin-tenants-header__content">
          <h1 className="admin-tenants-title">
            <Building2 />
            Tenants
          </h1>
          <span className="admin-tenants-count">{total} Organisationen</span>
        </div>

        <div className="admin-tenants-actions">
          <div className="admin-tenants-search">
            <Search />
            <input
              type="text"
              placeholder="Tenant suchen..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="admin-tenants-search__input"
            />
          </div>
          <motion.button
            className="admin-tenants-refresh"
            onClick={loadTenants}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <RefreshCw />
          </motion.button>
        </div>
      </motion.div>

      {error && (
        <motion.div
          className="admin-tenants-error"
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <AlertTriangle />
          <span>{error}</span>
        </motion.div>
      )}

      {/* Table */}
      <motion.div
        className="admin-tenants-table-wrapper"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <table className="admin-tenants-table">
          <thead>
            <tr>
              <th>Organisation</th>
              <th>Tier</th>
              <th>Status</th>
              <th>Nutzer / Gins</th>
              <th>Erstellt</th>
              <th></th>
            </tr>
          </thead>
          <motion.tbody variants={containerVariants} initial="hidden" animate="visible">
            {filteredTenants.map((item) => (
              <motion.tr
                key={item.tenant.id}
                variants={rowVariants}
                className={actionLoading === item.tenant.id ? 'admin-tenants-row--loading' : ''}
              >
                <td>
                  <div className="admin-tenants-org">
                    <div className="admin-tenants-org__avatar">
                      <Building2 />
                    </div>
                    <div className="admin-tenants-org__info">
                      <span className="admin-tenants-org__name">{item.tenant.name}</span>
                      <span className="admin-tenants-org__subdomain">{item.tenant.subdomain}</span>
                    </div>
                  </div>
                </td>
                <td>
                  <select
                    value={item.tenant.tier}
                    onChange={(e) => handleTierChange(item.tenant.id, e.target.value)}
                    className={`admin-tenants-tier-select ${getTierBadgeClass(item.tenant.tier)}`}
                    disabled={actionLoading === item.tenant.id}
                  >
                    <option value="free">Free</option>
                    <option value="basic">Basic</option>
                    <option value="pro">Pro</option>
                    <option value="enterprise">Enterprise</option>
                  </select>
                </td>
                <td>
                  <span className={`admin-tenants-badge ${getStatusBadgeClass(item.tenant.status)}`}>
                    {item.tenant.status === 'active' && 'Aktiv'}
                    {item.tenant.status === 'suspended' && 'Suspendiert'}
                    {item.tenant.status === 'cancelled' && 'Gekündigt'}
                  </span>
                </td>
                <td>
                  <div className="admin-tenants-stats">
                    <span className="admin-tenants-stats__item">
                      <Users />
                      {item.user_count}
                    </span>
                    <span className="admin-tenants-stats__divider">/</span>
                    <span className="admin-tenants-stats__item">
                      <Wine />
                      {item.gin_count}
                    </span>
                  </div>
                </td>
                <td>
                  <div className="admin-tenants-date">
                    <Calendar />
                    {new Date(item.tenant.created_at).toLocaleDateString('de-DE')}
                  </div>
                </td>
                <td>
                  <div className="admin-tenants-dropdown">
                    <motion.button
                      className="admin-tenants-dropdown__trigger"
                      onClick={() => setActiveDropdown(activeDropdown === item.tenant.id ? null : item.tenant.id)}
                      whileHover={{ scale: 1.1 }}
                      whileTap={{ scale: 0.9 }}
                    >
                      {actionLoading === item.tenant.id ? (
                        <Loader2 className="admin-tenants-dropdown__spinner" />
                      ) : (
                        <MoreVertical />
                      )}
                    </motion.button>

                    <AnimatePresence>
                      {activeDropdown === item.tenant.id && (
                        <motion.div
                          className="admin-tenants-dropdown__menu"
                          initial={{ opacity: 0, scale: 0.95, y: -10 }}
                          animate={{ opacity: 1, scale: 1, y: 0 }}
                          exit={{ opacity: 0, scale: 0.95, y: -10 }}
                        >
                          {item.tenant.status === 'active' ? (
                            <button
                              className="admin-tenants-dropdown__item admin-tenants-dropdown__item--danger"
                              onClick={() => handleSuspend(item.tenant.id)}
                            >
                              <Pause />
                              Suspendieren
                            </button>
                          ) : (
                            <button
                              className="admin-tenants-dropdown__item admin-tenants-dropdown__item--success"
                              onClick={() => handleActivate(item.tenant.id)}
                            >
                              <Play />
                              Aktivieren
                            </button>
                          )}
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </div>
                </td>
              </motion.tr>
            ))}
          </motion.tbody>
        </table>

        {filteredTenants.length === 0 && (
          <div className="admin-tenants-empty">
            <Building2 />
            <span>Keine Tenants gefunden</span>
          </div>
        )}
      </motion.div>

      {/* Pagination */}
      {total > 20 && (
        <motion.div
          className="admin-tenants-pagination"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          <button
            className="admin-tenants-pagination__btn"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
          >
            <ChevronLeft />
            Zurück
          </button>
          <span className="admin-tenants-pagination__info">
            Seite {page} von {Math.ceil(total / 20)}
          </span>
          <button
            className="admin-tenants-pagination__btn"
            onClick={() => setPage((p) => p + 1)}
            disabled={page * 20 >= total}
          >
            Weiter
            <ChevronRight />
          </button>
        </motion.div>
      )}

      {/* Click outside handler */}
      {activeDropdown && (
        <div
          className="admin-tenants-overlay"
          onClick={() => setActiveDropdown(null)}
        />
      )}
    </div>
  );
}
