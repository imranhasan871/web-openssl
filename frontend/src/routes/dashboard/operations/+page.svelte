<script lang="ts">
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import { apiClient } from '$lib/api/client';
  import { notifications } from '$lib/stores/notifications';
  import Button from '$lib/components/ui/Button.svelte';

  let operations: any[] = [];
  let loading = true;

  onMount(() => {
    const unsubscribe = authStore.subscribe((auth) => {
      if (!auth.loading && auth.isAuthenticated) {
        loadOperations();
        unsubscribe();
      }
    });
  });

  async function loadOperations() {
    loading = true;
    try {
      const response = await apiClient.get('/api/v1/operations/');
      if (response.success) {
        operations = response.data?.operations || [];
      } else {
        notifications.error('Error', 'Failed to load operations');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  async function deleteOperation(id: number) {
    if (!confirm('Are you sure you want to delete this operation?')) {
      return;
    }

    try {
      const response = await apiClient.delete(`/api/v1/operations/${id}`);
      if (response.success) {
        notifications.success('Success', 'Operation deleted');
        operations = operations.filter(op => op.id !== id);
      } else {
        notifications.error('Error', 'Failed to delete operation');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    }
  }

  function getOperationIcon(type: string) {
    switch (type) {
      case 'certificate':
        return 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z';
      case 'encryption':
        return 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z';
      case 'hash':
        return 'M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2h3a1 1 0 011 1v1a1 1 0 01-1 1h-1v9a2 2 0 01-2 2H7a2 2 0 01-2-2V7H4a1 1 0 01-1-1V5a1 1 0 011-1h3z';
      default:
        return 'M13 10V3L4 14h7v7l9-11h-7z';
    }
  }
</script>

<svelte:head>
  <title>Operations History - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-gray-900">Operations History</h1>
    <p class="mt-1 text-sm text-gray-500">
      View and manage your OpenSSL operation history
    </p>
  </div>

  {#if loading}
    <div class="bg-white shadow rounded-lg p-6">
      <div class="animate-pulse space-y-4">
        {#each Array(5) as _}
          <div class="flex space-x-4">
            <div class="rounded-full bg-gray-200 h-10 w-10"></div>
            <div class="flex-1 space-y-2 py-1">
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {:else if operations.length === 0}
    <div class="bg-white shadow rounded-lg">
      <div class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No operations yet</h3>
        <p class="mt-1 text-sm text-gray-500">Get started by performing your first operation.</p>
        <div class="mt-6">
          <Button href="/dashboard">Back to Dashboard</Button>
        </div>
      </div>
    </div>
  {:else}
    <div class="bg-white shadow rounded-lg overflow-hidden">
      <ul class="divide-y divide-gray-200">
        {#each operations as operation}
          <li class="px-6 py-4 hover:bg-gray-50">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div class="flex-shrink-0">
                  <span class="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                    <svg class="h-5 w-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getOperationIcon(operation.type)}/>
                    </svg>
                  </span>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-900">
                    {operation.type || 'Operation'}
                  </p>
                  <p class="text-sm text-gray-500">
                    {operation.description || 'No description'}
                  </p>
                  <p class="text-xs text-gray-400 mt-1">
                    {new Date(operation.createdAt).toLocaleString()}
                  </p>
                </div>
              </div>
              <button
                on:click={() => deleteOperation(operation.id)}
                class="text-red-600 hover:text-red-800 text-sm font-medium"
              >
                Delete
              </button>
            </div>
          </li>
        {/each}
      </ul>
    </div>
  {/if}
</div>