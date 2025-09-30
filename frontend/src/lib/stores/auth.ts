import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';

export interface User {
  id: number;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
  plan: string;
  isActive: boolean;
  usageCount: number;
  apiKey?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
  loading: boolean;
}

const initialState: AuthState = {
  isAuthenticated: false,
  user: null,
  token: null,
  loading: true
};

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>(initialState);

  return {
    subscribe,

    // Initialize auth state from localStorage
    init: () => {
      if (browser) {
        const token = localStorage.getItem('auth_token');
        const userStr = localStorage.getItem('auth_user');

        if (token && userStr) {
          try {
            const user = JSON.parse(userStr);
            update(state => ({
              ...state,
              isAuthenticated: true,
              user,
              token,
              loading: false
            }));
          } catch (e) {
            localStorage.removeItem('auth_token');
            localStorage.removeItem('auth_user');
            update(state => ({ ...state, loading: false }));
          }
        } else {
          update(state => ({ ...state, loading: false }));
        }
      }
    },

    // Login user
    login: (user: User, token: string) => {
      if (browser) {
        localStorage.setItem('auth_token', token);
        localStorage.setItem('auth_user', JSON.stringify(user));
      }

      update(state => ({
        ...state,
        isAuthenticated: true,
        user,
        token,
        loading: false
      }));
    },

    // Logout user
    logout: () => {
      if (browser) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
      }

      set({
        isAuthenticated: false,
        user: null,
        token: null,
        loading: false
      });

      goto('/login');
    },

    // Update user data
    updateUser: (user: User) => {
      if (browser) {
        localStorage.setItem('auth_user', JSON.stringify(user));
      }

      update(state => ({
        ...state,
        user
      }));
    },

    // Set loading state
    setLoading: (loading: boolean) => {
      update(state => ({ ...state, loading }));
    }
  };
}

export const authStore = createAuthStore();