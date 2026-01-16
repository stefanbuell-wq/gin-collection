import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuthStore, useUser, useTenant } from '../stores/authStore';
import {
  LayoutDashboard,
  Wine,
  Settings,
  CreditCard,
  Users,
  LogOut,
  Menu,
  X,
} from 'lucide-react';
import { useState } from 'react';

export const Layout = () => {
  const navigate = useNavigate();
  const logout = useAuthStore((state) => state.logout);
  const user = useUser();
  const tenant = useTenant();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const navigation = [
    { name: 'Dashboard', to: '/dashboard', icon: LayoutDashboard },
    { name: 'Gins', to: '/gins', icon: Wine },
    { name: 'Subscription', to: '/subscription', icon: CreditCard },
    { name: 'Settings', to: '/settings', icon: Settings },
  ];

  // Add Users link for Enterprise tier
  if (tenant?.tier === 'enterprise' && (user?.role === 'owner' || user?.role === 'admin')) {
    navigation.push({ name: 'Users', to: '/users', icon: Users });
  }

  return (
    <div className="min-h-screen" style={{ background: 'var(--vault-deep)' }}>
      {/* Header */}
      <header
        className="shadow-lg border-b"
        style={{
          background: 'var(--vault-dark)',
          borderColor: 'var(--glass-border)'
        }}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* Logo & Tenant Name */}
            <div className="flex items-center gap-3">
              <div className="text-2xl">🍸</div>
              <div>
                <h1
                  className="text-xl font-bold"
                  style={{
                    fontFamily: 'var(--font-display)',
                    color: 'var(--gold)'
                  }}
                >
                  GinVault
                </h1>
                {tenant && (
                  <p
                    className="text-xs"
                    style={{ color: 'var(--text-secondary)' }}
                  >
                    {tenant.name} • {tenant.tier.toUpperCase()}
                  </p>
                )}
              </div>
            </div>

            {/* Desktop Navigation */}
            <nav className="hidden md:flex items-center gap-2">
              {navigation.map((item) => {
                const Icon = item.icon;
                return (
                  <NavLink
                    key={item.to}
                    to={item.to}
                    className={({ isActive }) =>
                      `flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200`
                    }
                    style={({ isActive }) => ({
                      background: isActive ? 'var(--glass-bg)' : 'transparent',
                      color: isActive ? 'var(--gold)' : 'var(--text-secondary)',
                      border: isActive ? '1px solid var(--glass-border-hover)' : '1px solid transparent',
                    })}
                  >
                    <Icon className="w-4 h-4" />
                    {item.name}
                  </NavLink>
                );
              })}
            </nav>

            {/* User Menu */}
            <div className="flex items-center gap-4">
              <div className="hidden md:block text-right">
                <p
                  className="text-sm font-medium"
                  style={{ color: 'var(--text-primary)' }}
                >
                  {user?.email}
                </p>
                <p
                  className="text-xs capitalize"
                  style={{ color: 'var(--text-muted)' }}
                >
                  {user?.role}
                </p>
              </div>
              <button
                onClick={handleLogout}
                className="flex items-center gap-2 px-3 py-2 text-sm rounded-lg transition-all duration-200"
                style={{
                  color: 'var(--text-secondary)',
                  background: 'transparent',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'var(--glass-bg)';
                  e.currentTarget.style.color = 'var(--copper)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'transparent';
                  e.currentTarget.style.color = 'var(--text-secondary)';
                }}
              >
                <LogOut className="w-4 h-4" />
                <span className="hidden md:inline">Logout</span>
              </button>

              {/* Mobile menu button */}
              <button
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                className="md:hidden p-2 rounded-lg transition-colors"
                style={{ color: 'var(--text-primary)' }}
              >
                {mobileMenuOpen ? (
                  <X className="w-6 h-6" />
                ) : (
                  <Menu className="w-6 h-6" />
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <div
            className="md:hidden border-t"
            style={{
              background: 'var(--vault-dark)',
              borderColor: 'var(--glass-border)'
            }}
          >
            <nav className="px-4 py-4 space-y-2">
              {navigation.map((item) => {
                const Icon = item.icon;
                return (
                  <NavLink
                    key={item.to}
                    to={item.to}
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200"
                    style={({ isActive }) => ({
                      background: isActive ? 'var(--glass-bg)' : 'transparent',
                      color: isActive ? 'var(--gold)' : 'var(--text-secondary)',
                    })}
                  >
                    <Icon className="w-5 h-5" />
                    {item.name}
                  </NavLink>
                );
              })}
            </nav>
          </div>
        )}
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>
    </div>
  );
};
