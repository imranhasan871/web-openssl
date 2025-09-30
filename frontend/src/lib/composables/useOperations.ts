/**
 * useOperations Composable
 * Single Responsibility: Encapsulate operation management logic
 */

import { writable } from 'svelte/store';
import { operationService } from '$lib/services';
import { notifications } from '$lib/stores/notifications';

export function useOperations() {
  const operations = writable<any[]>([]);
  const stats = writable<any>({
    totalOperations: 0,
    certificatesGenerated: 0,
    encryptionOperations: 0,
    usageThisMonth: 0
  });
  const loading = writable(false);

  async function loadOperations(limit: number = 50) {
    loading.set(true);

    try {
      const response = await operationService.getOperations(limit);

      if (response.success && response.data) {
        operations.set(response.data);
      } else {
        notifications.error('Error', 'Failed to load operations');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading.set(false);
    }
  }

  async function loadStats() {
    try {
      const response = await operationService.getStats();

      if (response.success && response.data) {
        stats.set(response.data);
      } else {
        notifications.error('Error', 'Failed to load stats');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    }
  }

  async function deleteOperation(id: number) {
    try {
      const response = await operationService.deleteOperation(id);

      if (response.success) {
        notifications.success('Success', 'Operation deleted');
        // Remove from local state
        operations.update(ops => ops.filter(op => op.id !== id));
        return { success: true };
      } else {
        notifications.error('Error', 'Failed to delete operation');
        return { success: false };
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
      return { success: false };
    }
  }

  return {
    operations,
    stats,
    loading,
    loadOperations,
    loadStats,
    deleteOperation
  };
}