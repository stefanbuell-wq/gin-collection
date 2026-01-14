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
  gin_type?: string;
  abv?: number;
  fill_level?: string;
  price?: number;
  barcode?: string;
  purchase_date?: string;
  purchase_location?: string;
  rating?: number;
  nose_notes?: string;
  taste_notes?: string;
  finish_notes?: string;
  overall_notes?: string;
  serving_suggestion?: string;
  primary_photo_url?: string;
  is_favorite: boolean;
  is_available: boolean;
  created_at: string;
  updated_at: string;
}

export interface GinCreateRequest {
  name: string;
  brand?: string;
  country?: string;
  gin_type?: string;
  abv?: number;
  price?: number;
  barcode?: string;
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
