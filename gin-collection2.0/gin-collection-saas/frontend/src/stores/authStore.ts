import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User, Tenant } from '../types';
import { authAPI } from '../api/services';
import { fetchCSRFToken, clearCSRFToken } from '../api/client';

interface AuthState {
  user: User | null;
  tenant: Tenant | null;
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
  initializeAuth: () => Promise<void>;
  checkAuth: () => Promise<boolean>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      tenant: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (email: string, password: string) => {
        set({ isLoading: true });
        try {
          const response = await authAPI.login(email, password);
          // API returns { success: true, data: { user, tenant } }
          // Tokens are now stored in HttpOnly cookies by the server
          const apiResponse = response.data as unknown as { success: boolean; data: { user: User; tenant: Tenant } };
          const { user, tenant } = apiResponse.data;

          set({
            user,
            tenant,
            isAuthenticated: true,
            isLoading: false,
          });

          // Fetch CSRF token after successful login
          fetchCSRFToken();
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      logout: async () => {
        try {
          // Server clears HttpOnly cookies
          await authAPI.logout();
        } catch (error) {
          console.error('Logout error:', error);
        } finally {
          clearCSRFToken();

          set({
            user: null,
            tenant: null,
            isAuthenticated: false,
          });
        }
      },

      register: async (data) => {
        set({ isLoading: true });
        try {
          const response = await authAPI.register(data);
          // API returns { success: true, data: { user, tenant } }
          // Tokens are now stored in HttpOnly cookies by the server
          const apiResponse = response.data as unknown as { success: boolean; data: { user: User; tenant: Tenant } };
          const { user, tenant } = apiResponse.data;

          set({
            user,
            tenant,
            isAuthenticated: true,
            isLoading: false,
          });

          // Fetch CSRF token after successful registration
          fetchCSRFToken();
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      setUser: (user: User) => set({ user }),
      setTenant: (tenant: Tenant) => set({ tenant }),

      // Check if user is authenticated by calling /auth/me
      checkAuth: async () => {
        try {
          const response = await authAPI.me();
          const apiResponse = response.data as unknown as { success: boolean; data: User };
          if (apiResponse.success && apiResponse.data) {
            set({ user: apiResponse.data, isAuthenticated: true });
            return true;
          }
          return false;
        } catch {
          set({ user: null, tenant: null, isAuthenticated: false });
          return false;
        }
      },

      // Initialize auth on app load - verify auth state with server
      initializeAuth: async () => {
        const state = get();
        // If we think we're authenticated, verify with the server
        if (state.isAuthenticated) {
          await state.checkAuth();
          // Fetch CSRF token if still authenticated
          if (get().isAuthenticated) {
            fetchCSRFToken();
          }
        }
      },
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
