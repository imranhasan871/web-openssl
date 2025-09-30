<script lang="ts">
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';
  import { apiClient } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';

  let stats = {
    totalOperations: 0,
    certificatesGenerated: 0,
    encryptionOperations: 0,
    usageThisMonth: 0
  };
  let recentOperations: any[] = [];
  let loading = true;

  onMount(async () => {
    // Wait for auth store to be initialized
    const unsubscribe = authStore.subscribe((auth) => {
      if (!auth.loading && auth.isAuthenticated) {
        loadDashboardData();
        unsubscribe();
      }
    });
  });

  async function loadDashboardData() {
    loading = true;

    try {
      // Load user stats
      const statsResponse = await apiClient.get('/api/v1/operations/stats');
      if (statsResponse.success) {
        stats = statsResponse.data;
      }

      // Load recent operations
      const operationsResponse = await apiClient.get('/api/v1/operations/?limit=5');
      if (operationsResponse.success) {
        recentOperations = operationsResponse.data?.operations || [];
      }
    } catch (error) {
      notifications.error('Error', 'Failed to load dashboard data');
    } finally {
      loading = false;
    }
  }

  const quickActions = [
    {
      name: 'Generate Certificate',
      description: 'Create a new SSL certificate',
      icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z',
      href: '/dashboard/certificates',
      color: 'bg-blue-500 hover:bg-blue-600'
    },
    {
      name: 'Encrypt Data',
      description: 'Encrypt text or files',
      icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z',
      href: '/dashboard/encryption',
      color: 'bg-green-500 hover:bg-green-600'
    },
    {
      name: 'Generate Key',
      description: 'Create private/public keys',
      icon: 'M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z',
      href: '/dashboard/certificates',
      color: 'bg-purple-500 hover:bg-purple-600'
    },
    {
      name: 'Hash Generator',
      description: 'Generate hashes and HMACs',
      icon: 'M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2h3a1 1 0 011 1v1a1 1 0 01-1 1h-1v9a2 2 0 01-2 2H7a2 2 0 01-2-2V7H4a1 1 0 01-1-1V5a1 1 0 011-1h3z',
      href: '/dashboard/encryption',
      color: 'bg-yellow-500 hover:bg-yellow-600'
    }
  ];
</script>

<svelte:head>
  <title>Dashboard - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <!-- Page header -->
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-gray-900">
      Welcome back, {$authStore.user?.firstName}!
    </h1>
    <p class="mt-1 text-sm text-gray-500">
      Here's what's happening with your OpenSSL operations today.
    </p>
  </div>

  {#if loading}
    <!-- Loading state -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
      {#each Array(4) as _}
        <div class="bg-white overflow-hidden shadow rounded-lg animate-pulse">
          <div class="p-5">
            <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div class="h-8 bg-gray-200 rounded w-1/2"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <!-- Stats overview -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8">
      <div class="bg-white overflow-hidden shadow rounded-lg">
        <div class="p-5">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
              </svg>
            </div>
            <div class="ml-5 w-0 flex-1">
              <dl>
                <dt class="text-sm font-medium text-gray-500 truncate">Total Operations</dt>
                <dd class="text-lg font-medium text-gray-900">{stats.totalOperations}</dd>
              </dl>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white overflow-hidden shadow rounded-lg">
        <div class="p-5">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
              </svg>
            </div>
            <div class="ml-5 w-0 flex-1">
              <dl>
                <dt class="text-sm font-medium text-gray-500 truncate">Certificates</dt>
                <dd class="text-lg font-medium text-gray-900">{stats.certificatesGenerated}</dd>
              </dl>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white overflow-hidden shadow rounded-lg">
        <div class="p-5">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
              </svg>
            </div>
            <div class="ml-5 w-0 flex-1">
              <dl>
                <dt class="text-sm font-medium text-gray-500 truncate">Encryptions</dt>
                <dd class="text-lg font-medium text-gray-900">{stats.encryptionOperations}</dd>
              </dl>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white overflow-hidden shadow rounded-lg">
        <div class="p-5">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-6 w-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
              </svg>
            </div>
            <div class="ml-5 w-0 flex-1">
              <dl>
                <dt class="text-sm font-medium text-gray-500 truncate">This Month</dt>
                <dd class="text-lg font-medium text-gray-900">{stats.usageThisMonth}</dd>
              </dl>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick actions -->
    <div class="mb-8">
      <h2 class="text-lg font-medium text-gray-900 mb-4">Quick Actions</h2>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {#each quickActions as action}
          <a
            href={action.href}
            class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
          >
            <div>
              <span class="rounded-lg inline-flex p-3 {action.color} text-white">
                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={action.icon}/>
                </svg>
              </span>
            </div>
            <div class="mt-4">
              <h3 class="text-lg font-medium text-gray-900">
                <span class="absolute inset-0" aria-hidden="true"></span>
                {action.name}
              </h3>
              <p class="mt-2 text-sm text-gray-500">
                {action.description}
              </p>
            </div>
          </a>
        {/each}
      </div>
    </div>

    <!-- Recent operations -->
    <div class="bg-white shadow rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-medium text-gray-900">Recent Operations</h2>
          <Button href="/dashboard/operations" variant="outline" size="sm">
            View All
          </Button>
        </div>

        {#if recentOperations.length > 0}
          <div class="flow-root">
            <ul class="-mb-8">
              {#each recentOperations as operation, i}
                <li>
                  <div class="relative pb-8">
                    {#if i < recentOperations.length - 1}
                      <span class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true"></span>
                    {/if}
                    <div class="relative flex space-x-3">
                      <div>
                        <span class="h-8 w-8 rounded-full bg-blue-500 flex items-center justify-center ring-8 ring-white">
                          <svg class="h-4 w-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                          </svg>
                        </span>
                      </div>
                      <div class="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
                        <div>
                          <p class="text-sm text-gray-500">
                            <span class="font-medium text-gray-900">{operation.type}</span>
                            {operation.description}
                          </p>
                        </div>
                        <div class="text-right text-sm whitespace-nowrap text-gray-500">
                          <time>{new Date(operation.createdAt).toLocaleDateString()}</time>
                        </div>
                      </div>
                    </div>
                  </div>
                </li>
              {/each}
            </ul>
          </div>
        {:else}
          <div class="text-center py-6">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No operations yet</h3>
            <p class="mt-1 text-sm text-gray-500">Get started by creating your first SSL certificate.</p>
            <div class="mt-6">
              <Button href="/dashboard/certificates">
                Generate Certificate
              </Button>
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>