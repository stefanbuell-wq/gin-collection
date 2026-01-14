import { create } from 'zustand';
import type { Gin, GinStats, SearchParams } from '../types';
import { ginAPI } from '../api/services';

interface GinState {
  gins: Gin[];
  currentGin: Gin | null;
  stats: GinStats | null;
  total: number;
  page: number;
  limit: number;
  isLoading: boolean;
  error: string | null;

  // Actions
  fetchGins: (params?: SearchParams) => Promise<void>;
  fetchGin: (id: number) => Promise<void>;
  fetchStats: () => Promise<void>;
  createGin: (data: Partial<Gin>) => Promise<Gin>;
  updateGin: (id: number, data: Partial<Gin>) => Promise<void>;
  deleteGin: (id: number) => Promise<void>;
  searchGins: (query: string) => Promise<void>;
  setCurrentGin: (gin: Gin | null) => void;
  clearError: () => void;
}

export const useGinStore = create<GinState>((set, get) => ({
  gins: [],
  currentGin: null,
  stats: null,
  total: 0,
  page: 1,
  limit: 20,
  isLoading: false,
  error: null,

  fetchGins: async (params?: SearchParams) => {
    set({ isLoading: true, error: null });
    try {
      const response = await ginAPI.list(params);
      const { gins, total, page, limit } = response.data;

      set({
        gins,
        total,
        page,
        limit,
        isLoading: false,
      });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to fetch gins',
        isLoading: false,
      });
      throw error;
    }
  },

  fetchGin: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await ginAPI.get(id);
      set({
        currentGin: response.data,
        isLoading: false,
      });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to fetch gin',
        isLoading: false,
      });
      throw error;
    }
  },

  fetchStats: async () => {
    try {
      const response = await ginAPI.stats();
      set({ stats: response.data });
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    }
  },

  createGin: async (data: Partial<Gin>) => {
    set({ isLoading: true, error: null });
    try {
      const response = await ginAPI.create(data as any);
      const newGin = response.data;

      set((state) => ({
        gins: [newGin, ...state.gins],
        total: state.total + 1,
        isLoading: false,
      }));

      return newGin;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to create gin',
        isLoading: false,
      });
      throw error;
    }
  },

  updateGin: async (id: number, data: Partial<Gin>) => {
    set({ isLoading: true, error: null });
    try {
      const response = await ginAPI.update(id, data);
      const updatedGin = response.data;

      set((state) => ({
        gins: state.gins.map((g) => (g.id === id ? updatedGin : g)),
        currentGin: state.currentGin?.id === id ? updatedGin : state.currentGin,
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to update gin',
        isLoading: false,
      });
      throw error;
    }
  },

  deleteGin: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      await ginAPI.delete(id);

      set((state) => ({
        gins: state.gins.filter((g) => g.id !== id),
        total: state.total - 1,
        currentGin: state.currentGin?.id === id ? null : state.currentGin,
        isLoading: false,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to delete gin',
        isLoading: false,
      });
      throw error;
    }
  },

  searchGins: async (query: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await ginAPI.search(query);
      const { gins, total } = response.data;

      set({
        gins,
        total,
        isLoading: false,
      });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to search gins',
        isLoading: false,
      });
      throw error;
    }
  },

  setCurrentGin: (gin: Gin | null) => set({ currentGin: gin }),
  clearError: () => set({ error: null }),
}));
