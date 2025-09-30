/**
 * useUser Composable
 * Single Responsibility: Encapsulate user profile management
 */

import { writable } from 'svelte/store';
import { userService } from '$lib/services';
import { notifications } from '$lib/stores/notifications';
import { authStore } from '$lib/stores/auth';
import type { UpdateProfileRequest } from '$lib/services/interfaces';
import type { User } from '$lib/stores/auth';

export function useUser() {
  const user = writable<User | null>(null);
  const loading = writable(false);

  async function loadProfile() {
    loading.set(true);

    try {
      const response = await userService.getProfile();

      if (response.success && response.data) {
        user.set(response.data);
        authStore.updateUser(response.data);
      } else {
        notifications.error('Error', 'Failed to load profile');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading.set(false);
    }
  }

  async function updateProfile(data: UpdateProfileRequest) {
    loading.set(true);

    try {
      const response = await userService.updateProfile(data);

      if (response.success && response.data) {
        user.set(response.data);
        authStore.updateUser(response.data);
        notifications.success('Success', 'Profile updated successfully');
        return { success: true };
      } else {
        notifications.error('Error', response.error || 'Failed to update profile');
        return { success: false };
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
      return { success: false };
    } finally {
      loading.set(false);
    }
  }

  async function deleteAccount() {
    if (!confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
      return { success: false };
    }

    loading.set(true);

    try {
      const response = await userService.deleteAccount();

      if (response.success) {
        notifications.success('Success', 'Account deleted successfully');
        authStore.logout();
        return { success: true };
      } else {
        notifications.error('Error', 'Failed to delete account');
        return { success: false };
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
      return { success: false };
    } finally {
      loading.set(false);
    }
  }

  async function generateAPIKey() {
    loading.set(true);

    try {
      const response = await userService.generateAPIKey();

      if (response.success && response.data) {
        notifications.success('Success', 'New API key generated');
        return { success: true, apiKey: response.data.apiKey };
      } else {
        notifications.error('Error', 'Failed to generate API key');
        return { success: false };
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
      return { success: false };
    } finally {
      loading.set(false);
    }
  }

  return {
    user,
    loading,
    loadProfile,
    updateProfile,
    deleteAccount,
    generateAPIKey
  };
}