import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import type { AuthResponse, APIError } from '../types';

// Create axios instance with default config
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Helper to get tenant subdomain from auth storage
function getTenantSubdomain(): string | null {
  try {
    const authStorage = localStorage.getItem('auth-storage');
    if (authStorage) {
      const parsed = JSON.parse(authStorage);
      return parsed?.state?.tenant?.subdomain || null;
    }
  } catch {
    // Ignore parse errors
  }
  return null;
}

// Request interceptor to add auth token and tenant header
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Add JWT token from localStorage
    const token = localStorage.getItem('auth_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Add tenant subdomain header
    const hostname = window.location.hostname;
    const hostnameParts = hostname.split('.');

    // Check if we're on a tenant subdomain (e.g., basic.ginvault.cloud has 3+ parts)
    // Main domain (ginvault.cloud, www.ginvault.cloud) should use stored subdomain
    const isMainDomain = hostnameParts.length <= 2 ||
                         hostnameParts[0] === 'www' ||
                         hostnameParts[0] === 'ginvault' ||
                         hostnameParts[0] === 'localhost' ||
                         hostnameParts[0] === '127';

    if (config.headers) {
      // Always prefer stored subdomain from auth (set after login)
      const storedSubdomain = getTenantSubdomain();
      if (storedSubdomain) {
        config.headers['X-Tenant-Subdomain'] = storedSubdomain;
      } else if (!isMainDomain && hostnameParts[0]) {
        // Fallback: extract from hostname subdomain (e.g., basic.ginvault.cloud)
        config.headers['X-Tenant-Subdomain'] = hostnameParts[0];
      }
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors and token refresh
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error: AxiosError<APIError>) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Handle 401 Unauthorized (try to refresh token)
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
          const response = await axios.post<AuthResponse>('/api/v1/auth/refresh', {
            refresh_token: refreshToken,
          });

          const { token } = response.data;
          localStorage.setItem('auth_token', token);

          // Retry original request with new token
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${token}`;
          }
          return apiClient(originalRequest);
        }
      } catch (refreshError) {
        // Refresh failed, logout user
        localStorage.removeItem('auth_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    // Handle 403 Forbidden (feature not available in tier)
    if (error.response?.status === 403 && error.response.data?.upgrade_required) {
      // Could trigger upgrade modal here
      console.warn('Upgrade required:', error.response.data.error);
    }

    return Promise.reject(error);
  }
);

export default apiClient;

// Helper function to handle API errors
export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const apiError = error.response?.data as APIError;
    return apiError?.error || error.message || 'An unexpected error occurred';
  }
  return 'An unexpected error occurred';
}

// Helper to check if error requires upgrade
export function isUpgradeRequired(error: unknown): boolean {
  if (axios.isAxiosError(error)) {
    const apiError = error.response?.data as APIError;
    return apiError?.upgrade_required === true;
  }
  return false;
}
