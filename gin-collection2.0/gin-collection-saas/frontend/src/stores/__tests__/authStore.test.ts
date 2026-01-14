import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useAuthStore } from '../authStore';
import { authAPI } from '../../api/services';

// Mock API
vi.mock('../../api/services', () => ({
  authAPI: {
    login: vi.fn(),
    register: vi.fn(),
    logout: vi.fn(),
  },
}));

describe('authStore', () => {
  beforeEach(() => {
    // Reset store
    useAuthStore.setState({
      user: null,
      tenant: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
    });

    // Clear mocks
    vi.clearAllMocks();
    localStorage.clear();
  });

  describe('login', () => {
    it('should login successfully and set user state', async () => {
      const mockResponse = {
        data: {
          token: 'test-token',
          refresh_token: 'refresh-token',
          user: {
            id: 1,
            email: 'test@example.com',
            role: 'owner',
          },
          tenant: {
            id: 1,
            name: 'Test Tenant',
            tier: 'free',
          },
        },
      };

      vi.mocked(authAPI.login).mockResolvedValue(mockResponse as any);

      const { login } = useAuthStore.getState();
      await login('test@example.com', 'password');

      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(true);
      expect(state.user).toEqual(mockResponse.data.user);
      expect(state.tenant).toEqual(mockResponse.data.tenant);
      expect(state.token).toBe('test-token');
      expect(localStorage.setItem).toHaveBeenCalledWith('auth_token', 'test-token');
      expect(localStorage.setItem).toHaveBeenCalledWith('refresh_token', 'refresh-token');
    });

    it('should handle login errors', async () => {
      vi.mocked(authAPI.login).mockRejectedValue(new Error('Invalid credentials'));

      const { login } = useAuthStore.getState();

      await expect(login('test@example.com', 'wrong')).rejects.toThrow();

      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBeNull();
    });

    it('should set loading state during login', async () => {
      let resolveLogin: any;
      const loginPromise = new Promise((resolve) => {
        resolveLogin = resolve;
      });

      vi.mocked(authAPI.login).mockReturnValue(loginPromise as any);

      const { login } = useAuthStore.getState();
      const loginCall = login('test@example.com', 'password');

      // Check loading state is true
      expect(useAuthStore.getState().isLoading).toBe(true);

      resolveLogin({
        data: {
          token: 'token',
          user: { id: 1 },
          tenant: { id: 1 },
        },
      });

      await loginCall;

      // Check loading state is false after completion
      expect(useAuthStore.getState().isLoading).toBe(false);
    });
  });

  describe('logout', () => {
    it('should logout and clear state', async () => {
      // Set initial authenticated state
      useAuthStore.setState({
        user: { id: 1 } as any,
        tenant: { id: 1 } as any,
        token: 'token',
        isAuthenticated: true,
      });

      localStorage.setItem('auth_token', 'token');
      localStorage.setItem('refresh_token', 'refresh');

      vi.mocked(authAPI.logout).mockResolvedValue({} as any);

      const { logout } = useAuthStore.getState();
      await logout();

      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBeNull();
      expect(state.tenant).toBeNull();
      expect(state.token).toBeNull();
      expect(localStorage.removeItem).toHaveBeenCalledWith('auth_token');
      expect(localStorage.removeItem).toHaveBeenCalledWith('refresh_token');
    });

    it('should clear state even if API call fails', async () => {
      useAuthStore.setState({
        user: { id: 1 } as any,
        isAuthenticated: true,
      });

      vi.mocked(authAPI.logout).mockRejectedValue(new Error('Network error'));

      const { logout } = useAuthStore.getState();
      await logout();

      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBeNull();
    });
  });

  describe('register', () => {
    it('should register successfully and set user state', async () => {
      const mockResponse = {
        data: {
          token: 'test-token',
          user: {
            id: 1,
            email: 'new@example.com',
            role: 'owner',
          },
          tenant: {
            id: 1,
            name: 'New Tenant',
            tier: 'free',
          },
        },
      };

      vi.mocked(authAPI.register).mockResolvedValue(mockResponse as any);

      const { register } = useAuthStore.getState();
      await register({
        tenant_name: 'New Tenant',
        subdomain: 'newtenant',
        email: 'new@example.com',
        password: 'password123',
      });

      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(true);
      expect(state.user?.email).toBe('new@example.com');
      expect(state.tenant?.name).toBe('New Tenant');
    });
  });

  describe('setUser', () => {
    it('should update user state', () => {
      const newUser = {
        id: 2,
        email: 'updated@example.com',
      } as any;

      const { setUser } = useAuthStore.getState();
      setUser(newUser);

      expect(useAuthStore.getState().user).toEqual(newUser);
    });
  });

  describe('setTenant', () => {
    it('should update tenant state', () => {
      const newTenant = {
        id: 2,
        name: 'Updated Tenant',
        tier: 'pro',
      } as any;

      const { setTenant } = useAuthStore.getState();
      setTenant(newTenant);

      expect(useAuthStore.getState().tenant).toEqual(newTenant);
    });
  });
});
