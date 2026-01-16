import axios from 'axios';

const api = axios.create({
  baseURL: '/admin/api/v1',
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('admin_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const adminApi = {
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),

  me: () => api.get('/auth/me'),

  getStats: () => api.get('/stats/overview'),

  getTenants: (page = 1, limit = 20) =>
    api.get(`/tenants?page=${page}&limit=${limit}`),

  suspendTenant: (id: number) =>
    api.post(`/tenants/${id}/suspend`),

  activateTenant: (id: number) =>
    api.post(`/tenants/${id}/activate`),

  updateTenantTier: (id: number, tier: string) =>
    api.put(`/tenants/${id}/tier`, { tier }),

  getUsers: (page = 1, limit = 20) =>
    api.get(`/users?page=${page}&limit=${limit}`),

  getHealth: () => api.get('/health'),

  // Server Management
  getServerStatus: () => api.get('/server/status'),

  getQuickActions: () => api.get('/server/actions'),

  executeAction: (actionId: string) =>
    api.post(`/server/actions/${actionId}`),

  deploy: (services: string[], pull = true, noCache = false) =>
    api.post('/server/deploy', { services, pull, no_cache: noCache }),

  restartService: (service: string) =>
    api.post(`/server/restart/${service}`),

  getServiceLogs: (service: string, lines = 100) =>
    api.get(`/server/logs/${service}?lines=${lines}`),

  reloadNginx: () =>
    api.post('/server/nginx/reload'),
};

export default api;
