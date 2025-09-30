import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from '../../../routes/$types';
import { get } from 'svelte/store';
import { authStore } from '$lib/stores/auth';
import { browser } from '$app/environment';

export const requireAuth: LayoutLoad = async ({ url }) => {
  if (browser) {
    const auth = get(authStore);

    if (!auth.isAuthenticated && !auth.loading) {
      throw redirect(302, `/login?redirect=${encodeURIComponent(url.pathname)}`);
    }
  }

  return {};
};

export const requireGuest: LayoutLoad = async ({ url }) => {
  if (browser) {
    const auth = get(authStore);

    if (auth.isAuthenticated) {
      throw redirect(302, '/dashboard');
    }
  }

  return {};
};

export const requireRole = (requiredRole: string): LayoutLoad => {
  return async ({ url }) => {
    if (browser) {
      const auth = get(authStore);

      if (!auth.isAuthenticated) {
        throw redirect(302, `/login?redirect=${encodeURIComponent(url.pathname)}`);
      }

      if (auth.user?.role !== requiredRole) {
        throw redirect(302, '/unauthorized');
      }
    }

    return {};
  };
};