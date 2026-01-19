import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import type { APIError } from '../types';

// CSRF token storage
let csrfToken: string | null = null;

// Create axios instance with default config
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Required for CSRF cookies
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

// Fetch CSRF token from server
export async function fetchCSRFToken(): Promise<string | null> {
  try {
    const response = await axios.get('/api/v1/csrf-token', {
      withCredentials: true,
    });
    csrfToken = response.data?.csrf_token || response.headers['x-csrf-token'];
    return csrfToken;
  } catch (error) {
    console.error('Failed to fetch CSRF token:', error);
    return null;
  }
}

// Get current CSRF token (fetch if not available)
export async function getCSRFToken(): Promise<string | null> {
  if (!csrfToken) {
    return fetchCSRFToken();
  }
  return csrfToken;
}

// Clear CSRF token (call on logout)
export function clearCSRFToken(): void {
  csrfToken = null;
}

// Request interceptor to add tenant header and CSRF token
// Note: JWT is now stored in HttpOnly cookies (sent automatically by browser)
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Add CSRF token for state-changing requests
    const method = config.method?.toUpperCase();
    if (csrfToken && config.headers && method && ['POST', 'PUT', 'DELETE', 'PATCH'].includes(method)) {
      config.headers['X-CSRF-Token'] = csrfToken;
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
    // JWT tokens are now in HttpOnly cookies - refresh token is sent automatically
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Call refresh endpoint - HttpOnly cookie is sent automatically by browser
        await axios.post('/api/v1/auth/refresh', {}, {
          withCredentials: true,
        });

        // Retry original request - new access token cookie is now set
        return apiClient(originalRequest);
      } catch (refreshError) {
        // Refresh failed, redirect to login
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    // Handle 403 Forbidden (feature not available in tier)
    if (error.response?.status === 403 && error.response.data?.upgrade_required) {
      // Could trigger upgrade modal here
      console.warn('Upgrade required:', error.response.data.error);
    }

    // Handle CSRF token errors - refetch token and retry
    if (error.response?.status === 403 && error.response.data?.code?.startsWith('CSRF_')) {
      console.warn('CSRF token error, refetching token...');
      const newToken = await fetchCSRFToken();
      if (newToken && originalRequest.headers) {
        originalRequest.headers['X-CSRF-Token'] = newToken;
        return apiClient(originalRequest);
      }
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
