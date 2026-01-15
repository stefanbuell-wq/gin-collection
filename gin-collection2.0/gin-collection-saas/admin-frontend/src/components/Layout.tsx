import { Link, useLocation, useNavigate } from 'react-router-dom';
import { ReactNode, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuthStore } from '../stores/authStore';
import {
  Wine,
  LayoutDashboard,
  Building2,
  Users,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Shield,
  Menu,
  X
} from 'lucide-react';
import './Layout.css';

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { admin, logout } = useAuthStore();
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isMobileOpen, setIsMobileOpen] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const navItems = [
    { path: '/', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/tenants', label: 'Tenants', icon: Building2 },
    { path: '/users', label: 'Users', icon: Users },
  ];

  const sidebarVariants = {
    expanded: { width: 280 },
    collapsed: { width: 80 }
  };

  const labelVariants = {
    show: { opacity: 1, x: 0, display: 'block' },
    hide: { opacity: 0, x: -10, transitionEnd: { display: 'none' } }
  };

  return (
    <div className="admin-layout">
      {/* Mobile Overlay */}
      <AnimatePresence>
        {isMobileOpen && (
          <motion.div
            className="admin-layout-overlay"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setIsMobileOpen(false)}
          />
        )}
      </AnimatePresence>

      {/* Sidebar */}
      <motion.aside
        className={`admin-sidebar ${isMobileOpen ? 'admin-sidebar--open' : ''}`}
        variants={sidebarVariants}
        animate={isCollapsed ? 'collapsed' : 'expanded'}
        transition={{ type: 'spring', stiffness: 300, damping: 30 }}
      >
        {/* Logo */}
        <div className="admin-sidebar-header">
          <Link to="/" className="admin-sidebar-logo">
            <div className="admin-sidebar-logo__icon">
              <Wine />
            </div>
            <AnimatePresence>
              {!isCollapsed && (
                <motion.div
                  className="admin-sidebar-logo__text"
                  variants={labelVariants}
                  initial="hide"
                  animate="show"
                  exit="hide"
                >
                  <span className="admin-sidebar-logo__title">GinVault</span>
                  <span className="admin-sidebar-logo__badge">Admin</span>
                </motion.div>
              )}
            </AnimatePresence>
          </Link>

          {/* Mobile Close Button */}
          <button
            className="admin-sidebar-close"
            onClick={() => setIsMobileOpen(false)}
          >
            <X />
          </button>
        </div>

        {/* Navigation */}
        <nav className="admin-sidebar-nav">
          <ul className="admin-sidebar-menu">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;

              return (
                <li key={item.path}>
                  <Link
                    to={item.path}
                    className={`admin-sidebar-link ${isActive ? 'admin-sidebar-link--active' : ''}`}
                    onClick={() => setIsMobileOpen(false)}
                  >
                    <div className="admin-sidebar-link__icon">
                      <Icon />
                    </div>
                    <AnimatePresence>
                      {!isCollapsed && (
                        <motion.span
                          className="admin-sidebar-link__label"
                          variants={labelVariants}
                          initial="hide"
                          animate="show"
                          exit="hide"
                        >
                          {item.label}
                        </motion.span>
                      )}
                    </AnimatePresence>
                    {isActive && (
                      <motion.div
                        className="admin-sidebar-link__indicator"
                        layoutId="activeIndicator"
                        transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                      />
                    )}
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>

        {/* Collapse Toggle */}
        <button
          className="admin-sidebar-toggle"
          onClick={() => setIsCollapsed(!isCollapsed)}
        >
          {isCollapsed ? <ChevronRight /> : <ChevronLeft />}
        </button>

        {/* User Section */}
        <div className="admin-sidebar-footer">
          <div className={`admin-sidebar-user ${isCollapsed ? 'admin-sidebar-user--collapsed' : ''}`}>
            <div className="admin-sidebar-user__avatar">
              <Shield />
            </div>
            <AnimatePresence>
              {!isCollapsed && (
                <motion.div
                  className="admin-sidebar-user__info"
                  variants={labelVariants}
                  initial="hide"
                  animate="show"
                  exit="hide"
                >
                  <span className="admin-sidebar-user__name">{admin?.name || 'Admin'}</span>
                  <span className="admin-sidebar-user__role">Platform Admin</span>
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          <button
            onClick={handleLogout}
            className={`admin-sidebar-logout ${isCollapsed ? 'admin-sidebar-logout--collapsed' : ''}`}
          >
            <LogOut />
            <AnimatePresence>
              {!isCollapsed && (
                <motion.span
                  variants={labelVariants}
                  initial="hide"
                  animate="show"
                  exit="hide"
                >
                  Abmelden
                </motion.span>
              )}
            </AnimatePresence>
          </button>
        </div>
      </motion.aside>

      {/* Main Content Area */}
      <div className={`admin-main ${isCollapsed ? 'admin-main--expanded' : ''}`}>
        {/* Mobile Header */}
        <header className="admin-mobile-header">
          <button
            className="admin-mobile-menu"
            onClick={() => setIsMobileOpen(true)}
          >
            <Menu />
          </button>
          <div className="admin-mobile-logo">
            <Wine />
            <span>GinVault Admin</span>
          </div>
        </header>

        {/* Content */}
        <main className="admin-content">
          {children}
        </main>
      </div>
    </div>
  );
}
