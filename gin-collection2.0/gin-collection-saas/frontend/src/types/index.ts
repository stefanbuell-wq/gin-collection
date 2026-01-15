// User & Authentication Types
export interface User {
  id: number;
  tenant_id: number;
  uuid: string;
  email: string;
  first_name?: string;
  last_name?: string;
  role: UserRole;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export type UserRole = 'owner' | 'admin' | 'member' | 'viewer';

export interface AuthResponse {
  token: string;
  refresh_token?: string;
  user: User;
  tenant: Tenant;
}

// Tenant Types
export interface Tenant {
  id: number;
  uuid: string;
  name: string;
  subdomain: string;
  tier: TenantTier;
  status: TenantStatus;
  branding?: TenantBranding;
  created_at: string;
  updated_at: string;
}

export type TenantTier = 'free' | 'basic' | 'pro' | 'enterprise';
export type TenantStatus = 'active' | 'suspended' | 'cancelled';

export interface TenantBranding {
  logo_url?: string;
  primary_color?: string;
  secondary_color?: string;
}

export interface TenantLimits {
  max_gins: number;
  max_photos_per_gin: number;
  storage_limit_mb?: number;
  features: string[];
}

// Gin Types
export interface Gin {
  id: number;
  tenant_id: number;
  uuid: string;
  name: string;
  brand?: string;
  country?: string;
  region?: string;
  gin_type?: string;
  abv?: number;
  bottle_size?: number;
  fill_level?: number;
  price?: number;
  current_market_value?: number;
  barcode?: string;
  purchase_date?: string;
  purchase_location?: string;
  rating?: number;
  nose_notes?: string;
  palate_notes?: string;
  finish_notes?: string;
  general_notes?: string;
  description?: string;
  recommended_tonic?: string;
  recommended_garnish?: string;
  photo_url?: string;
  primary_photo_url?: string;
  is_finished: boolean;
  is_favorite: boolean;
  is_available: boolean;
  created_at: string;
  updated_at: string;
}

export interface GinCreateRequest {
  name: string;
  brand?: string;
  country?: string;
  region?: string;
  gin_type?: string;
  abv?: number;
  bottle_size?: number;
  fill_level?: number;
  price?: number;
  current_market_value?: number;
  barcode?: string;
  purchase_date?: string;
  purchase_location?: string;
  rating?: number;
  nose_notes?: string;
  palate_notes?: string;
  finish_notes?: string;
  general_notes?: string;
  description?: string;
  recommended_tonic?: string;
  recommended_garnish?: string;
  is_finished?: boolean;
}

export interface GinListResponse {
  gins: Gin[];
  total: number;
  page: number;
  limit: number;
}

export interface GinStats {
  total_gins: number;
  available_gins: number;
  favorite_count: number;
  avg_rating: number;
  total_value: number;
  by_country: Record<string, number>;
  by_type: Record<string, number>;
}

// Photo Types
export interface GinPhoto {
  id: number;
  tenant_id: number;
  gin_id: number;
  photo_url: string;
  photo_type: PhotoType;
  caption?: string;
  is_primary: boolean;
  created_at: string;
}

export type PhotoType = 'bottle' | 'label' | 'moment' | 'tasting';

// Subscription Types
export interface Subscription {
  id: number;
  tenant_id: number;
  plan_id: string;
  status: SubscriptionStatus;
  billing_cycle: BillingCycle;
  amount: number;
  currency: string;
  current_period_start?: string;
  current_period_end?: string;
  next_billing_date?: string;
  created_at: string;
  updated_at: string;
}

export type SubscriptionStatus = 'active' | 'pending' | 'cancelled' | 'suspended' | 'expired';
export type BillingCycle = 'monthly' | 'yearly';

export interface SubscriptionPlan {
  id: string;
  name: string;
  tier: TenantTier;
  description: string;
  price_monthly: number;
  price_yearly: number;
  features: string[];
  limits: TenantLimits;
}

// Botanical Types
export interface Botanical {
  id: number;
  name: string;
  category: string;
  description?: string;
}

export interface GinBotanical {
  gin_id: number;
  botanical_id: number;
  botanical: Botanical;
  prominence: BotanicalProminence;
}

export type BotanicalProminence = 'dominant' | 'notable' | 'subtle';

// Cocktail Types
export interface Cocktail {
  id: number;
  name: string;
  description?: string;
  instructions?: string;
  difficulty?: string;
  ingredients: CocktailIngredient[];
}

export interface CocktailIngredient {
  ingredient: string;
  amount: string;
  is_gin: boolean;
}

// API Error Types
export interface APIError {
  error: string;
  upgrade_required?: boolean;
  code?: string;
}

// Search & Filter Types
export interface SearchParams {
  q?: string;
  filter?: 'all' | 'available' | 'favorite';
  sort?: 'name' | 'rating' | 'price' | 'created_at';
  type?: string;
  country?: string;
  page?: number;
  limit?: number;
}
