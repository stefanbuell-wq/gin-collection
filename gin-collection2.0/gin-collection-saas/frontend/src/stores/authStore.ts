import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User, Tenant } from '../types';
import { authAPI } from '../api/services';

interface AuthState {
  user: User | null;
  tenant: Tenant | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  register: (data: {
    tenant_name: string;
    subdomain: string;
    email: string;
    password: string;
    first_name?: string;
    last_name?: string;
  }) => Promise<void>;
  setUser: (user: User) => void;
  setTenant: (tenant: Tenant) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      tenant: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (email: string, password: string) => {
        set({ isLoading: true });
        try {
          const response = await authAPI.login(email, password);
          // API returns { success: true, data: { token, user, tenant, ... } }
          const apiResponse = response.data as unknown as { success: boolean; data: { token: string; refresh_token?: string; user: User; tenant: Tenant } };
          const { token, refresh_token, user, tenant } = apiResponse.data;

          // Store tokens
          localStorage.setItem('auth_token', token);
          if (refresh_token) {
            localStorage.setItem('refresh_token', refresh_token);
          }

          set({
            user,
            tenant,
            token,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      logout: async () => {
        try {
          await authAPI.logout();
        } catch (error) {
          console.error('Logout error:', error);
        } finally {
          // Clear tokens
          localStorage.removeItem('auth_token');
          localStorage.removeItem('refresh_token');

          set({
            user: null,
            tenant: null,
            token: null,
            isAuthenticated: false,
          });
        }
      },

      register: async (data) => {
        set({ isLoading: true });
        try {
          const response = await authAPI.register(data);
          // API returns { success: true, data: { token, user, tenant, ... } }
          const apiResponse = response.data as unknown as { success: boolean; data: { token: string; refresh_token?: string; user: User; tenant: Tenant } };
          const { token, refresh_token, user, tenant } = apiResponse.data;

          // Store tokens
          localStorage.setItem('auth_token', token);
          if (refresh_token) {
            localStorage.setItem('refresh_token', refresh_token);
          }

          set({
            user,
            tenant,
            token,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      setUser: (user: User) => set({ user }),
      setTenant: (tenant: Tenant) => set({ tenant }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        tenant: state.tenant,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

// Helper hooks
export const useUser = () => useAuthStore((state) => state.user);
export const useTenant = () => useAuthStore((state) => state.tenant);
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
