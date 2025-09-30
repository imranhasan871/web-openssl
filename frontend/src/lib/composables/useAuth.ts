/**
 * useAuth Composable
 * Single Responsibility: Encapsulate authentication logic
 * Reusable across components
 */

import { authStore } from '$lib/stores/auth';
import { notifications } from '$lib/stores/notifications';
import { authService } from '$lib/services';
import { goto } from '$app/navigation';
import type { LoginRequest, RegisterRequest } from '$lib/services/interfaces';

export function useAuth() {
  async function login(credentials: LoginRequest) {
    try {
      const response = await authService.login(credentials);

      if (response.success && response.data) {
        authStore.login(response.data.user, response.data.accessToken);
        notifications.success('Success', 'Successfully signed in!');
        await goto('/dashboard', { replaceState: true });
        return { success: true };
      } else {
        notifications.error('Login Failed', response.error || 'Invalid credentials');
        return { success: false, error: response.error };
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
      return { success: false, error: 'Unexpected error' };
    }
  }

  async function register(userData: RegisterRequest) {
    try {
      const response = await authService.register(userData);

      if (response.success && response.data) {
        authStore.login(response.data.user, response.data.accessToken);
        notifications.success('Welcome!', 'Account created successfully');
        await goto('/dashboard', { replaceState: true });
        return { success: true };
      } else {
        notifications.error('Registration Failed', response.error || 'Failed to create account');
        return { success: false, error: response.error };
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
      return { success: false, error: 'Unexpected error' };
    }
  }

  function logout() {
    authStore.logout();
    notifications.info('Signed Out', 'You have been signed out successfully');
  }

  async function refreshToken() {
    try {
      const response = await authService.refreshToken();

      if (response.success && response.data) {
        // Update token in store
        // Note: You'd need to add this method to authStore
        return { success: true };
      }

      return { success: false };
    } catch (error) {
      return { success: false };
    }
  }

  return {
    login,
    register,
    logout,
    refreshToken,
    // Expose store for reactive subscriptions
    store: authStore
  };
}