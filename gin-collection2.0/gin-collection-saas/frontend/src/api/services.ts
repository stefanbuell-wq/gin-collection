import apiClient from './client';
import type {
  AuthResponse,
  User,
  Tenant,
  Gin,
  GinCreateRequest,
  GinListResponse,
  GinStats,
  GinPhoto,
  Subscription,
  SubscriptionPlan,
  Botanical,
  GinBotanical,
  Cocktail,
  SearchParams,
  GinReference,
  GinReferenceFilters,
} from '../types';

// ============================================================================
// Authentication API
// ============================================================================

export const authAPI = {
  login: (email: string, password: string) =>
    apiClient.post<AuthResponse>('/auth/login', { email, password }),

  register: (data: {
    tenant_name: string;
    subdomain: string;
    email: string;
    password: string;
    first_name?: string;
    last_name?: string;
  }) => apiClient.post<AuthResponse>('/auth/register', data),

  refreshToken: (refreshToken: string) =>
    apiClient.post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken }),

  logout: () => apiClient.post('/auth/logout'),

  // Profile management
  getMe: () => apiClient.get<User>('/auth/me'),

  updateProfile: (data: { first_name?: string; last_name?: string }) =>
    apiClient.put<User>('/auth/profile', data),

  changePassword: (currentPassword: string, newPassword: string) =>
    apiClient.post('/auth/change-password', {
      current_password: currentPassword,
      new_password: newPassword,
    }),
};

// ============================================================================
// Tenant API
// ============================================================================

export const tenantAPI = {
  getCurrent: () => apiClient.get<Tenant>('/tenants/current'),

  updateCurrent: (data: Partial<Tenant>) =>
    apiClient.put<Tenant>('/tenants/current', data),

  getUsage: () =>
    apiClient.get<{
      gin_count: number;
      storage_mb: number;
      limits: {
        max_gins: number;
        max_photos_per_gin: number;
        storage_limit_mb?: number;
      };
    }>('/tenants/usage'),
};

// ============================================================================
// Gin API
// ============================================================================

export const ginAPI = {
  list: (params?: SearchParams) =>
    apiClient.get<GinListResponse>('/gins', { params }),

  get: (id: number) => apiClient.get<Gin>(`/gins/${id}`),

  create: (data: GinCreateRequest) => apiClient.post<Gin>('/gins', data),

  update: (id: number, data: Partial<Gin>) =>
    apiClient.put<Gin>(`/gins/${id}`, data),

  delete: (id: number) => apiClient.delete(`/gins/${id}`),

  search: (query: string) =>
    apiClient.get<{ gins: Gin[]; total: number }>('/gins/search', {
      params: { q: query },
    }),

  stats: () => apiClient.get<GinStats>('/gins/stats'),

  export: (format: 'json' | 'csv') =>
    apiClient.post('/gins/export', { format }, { responseType: 'blob' }),

  import: (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return apiClient.post('/gins/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },

  suggestions: (id: number, limit = 5) =>
    apiClient.get<{ suggestions: Gin[]; count: number }>(`/gins/${id}/suggestions`, {
      params: { limit },
    }),
};

// ============================================================================
// Photo API
// ============================================================================

export const photoAPI = {
  getPhotos: (ginId: number) =>
    apiClient.get<{ photos: GinPhoto[]; count: number }>(`/gins/${ginId}/photos`),

  upload: (ginId: number, file: File, photoType: string, caption?: string) => {
    const formData = new FormData();
    formData.append('photo', file);
    formData.append('photo_type', photoType);
    if (caption) formData.append('caption', caption);

    return apiClient.post<GinPhoto>(`/gins/${ginId}/photos`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },

  delete: (ginId: number, photoId: number) =>
    apiClient.delete(`/gins/${ginId}/photos/${photoId}`),

  setPrimary: (ginId: number, photoId: number) =>
    apiClient.put(`/gins/${ginId}/photos/${photoId}/primary`),
};

// ============================================================================
// Subscription API
// ============================================================================

export const subscriptionAPI = {
  getCurrent: () => apiClient.get<Subscription>('/subscriptions/current'),

  getPlans: () => apiClient.get<{ plans: SubscriptionPlan[] }>('/subscriptions/plans'),

  upgrade: (planId: string, billingCycle: 'monthly' | 'yearly') =>
    apiClient.post<{ approval_url: string; subscription_id: string }>(
      '/subscriptions/upgrade',
      { plan_id: planId, billing_cycle: billingCycle }
    ),

  activate: (subscriptionId: string) =>
    apiClient.post('/subscriptions/activate', { subscription_id: subscriptionId }),

  cancel: () => apiClient.post('/subscriptions/cancel'),
};

// ============================================================================
// Botanical API
// ============================================================================

export const botanicalAPI = {
  getAll: () => apiClient.get<{ botanicals: Botanical[] }>('/botanicals'),

  getGinBotanicals: (ginId: number) =>
    apiClient.get<{ botanicals: GinBotanical[] }>(`/gins/${ginId}/botanicals`),

  updateGinBotanicals: (
    ginId: number,
    botanicals: { botanical_id: number; prominence: string }[]
  ) => apiClient.put(`/gins/${ginId}/botanicals`, { botanicals }),
};

// ============================================================================
// Cocktail API
// ============================================================================

export const cocktailAPI = {
  getAll: () => apiClient.get<{ cocktails: Cocktail[] }>('/cocktails'),

  getById: (id: number) => apiClient.get<Cocktail>(`/cocktails/${id}`),

  getGinCocktails: (ginId: number) =>
    apiClient.get<{ cocktails: Cocktail[] }>(`/gins/${ginId}/cocktails`),

  linkCocktail: (ginId: number, cocktailId: number) =>
    apiClient.post(`/gins/${ginId}/cocktails/${cocktailId}`),

  unlinkCocktail: (ginId: number, cocktailId: number) =>
    apiClient.delete(`/gins/${ginId}/cocktails/${cocktailId}`),
};

// ============================================================================
// User API (Enterprise only)
// ============================================================================

export const userAPI = {
  list: () => apiClient.get<{ users: User[]; count: number }>('/users'),

  invite: (data: {
    email: string;
    first_name?: string;
    last_name?: string;
    role: string;
  }) => apiClient.post<User>('/users/invite', data),

  update: (
    id: number,
    data: {
      email: string;
      first_name?: string;
      last_name?: string;
      role: string;
      is_active: boolean;
    }
  ) => apiClient.put<User>(`/users/${id}`, data),

  delete: (id: number) => apiClient.delete(`/users/${id}`),

  generateAPIKey: (id: number) =>
    apiClient.post<{ api_key: string; message: string }>(`/users/${id}/api-key`),

  revokeAPIKey: (id: number) => apiClient.delete(`/users/${id}/api-key`),
};

// ============================================================================
// Gin Reference API (catalog for quick add)
// ============================================================================

export const ginReferenceAPI = {
  search: (params?: { q?: string; country?: string; type?: string; limit?: number; offset?: number }) =>
    apiClient.get<{ data: { gins: GinReference[]; total: number; limit: number; offset: number } }>(
      '/gin-references',
      { params }
    ),

  getById: (id: number) =>
    apiClient.get<{ data: GinReference }>(`/gin-references/${id}`),

  getFilters: () =>
    apiClient.get<{ data: GinReferenceFilters }>('/gin-references/filters'),
};
