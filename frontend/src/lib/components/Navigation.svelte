<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';
  import Button from './ui/Button.svelte';

  let mobileMenuOpen = false;

  function handleLogout() {
    authStore.logout();
    notifications.info('Signed Out', 'You have been signed out successfully');
  }

  function toggleMobileMenu() {
    mobileMenuOpen = !mobileMenuOpen;
  }
</script>

<nav class="bg-white shadow">
  <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
    <div class="flex h-16 justify-between">
      <div class="flex">
        <!-- Logo -->
        <div class="flex flex-shrink-0 items-center">
          <a href="/" class="flex items-center space-x-2">
            <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
            </svg>
            <span class="text-xl font-bold text-gray-900">OpenSSL UI</span>
          </a>
        </div>

        <!-- Desktop Navigation -->
        {#if $authStore.isAuthenticated}
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <a
              href="/dashboard"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 hover:text-blue-600"
            >
              Dashboard
            </a>
            <a
              href="/dashboard/certificates"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 hover:text-gray-700"
            >
              Certificates
            </a>
            <a
              href="/dashboard/encryption"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 hover:text-gray-700"
            >
              Encryption
            </a>
            <a
              href="/dashboard/operations"
              class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 hover:text-gray-700"
            >
              History
            </a>
            {#if $authStore.user?.role === 'admin'}
              <a
                href="/admin"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 hover:text-gray-700"
              >
                Admin
              </a>
            {/if}
          </div>
        {/if}
      </div>

      <!-- Right side -->
      <div class="flex items-center">
        {#if $authStore.isAuthenticated}
          <!-- User dropdown -->
          <div class="relative ml-3">
            <div class="flex items-center space-x-4">
              <span class="text-sm text-gray-700">
                {$authStore.user?.firstName} {$authStore.user?.lastName}
              </span>
              <div class="flex space-x-2">
                <a
                  href="/profile"
                  class="text-gray-400 hover:text-gray-500"
                  title="Profile"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                </a>
                <button
                  on:click={handleLogout}
                  class="text-gray-400 hover:text-gray-500"
                  title="Sign Out"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        {:else}
          <!-- Guest buttons -->
          <div class="flex space-x-4">
            <a
              href="/login"
              class="text-gray-500 hover:text-gray-700 px-3 py-2 text-sm font-medium"
            >
              Sign In
            </a>
            <Button href="/register" size="sm">
              Sign Up
            </Button>
          </div>
        {/if}

        <!-- Mobile menu button -->
        <div class="flex items-center sm:hidden">
          <button
            type="button"
            class="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500"
            on:click={toggleMobileMenu}
          >
            <span class="sr-only">Open main menu</span>
            {#if !mobileMenuOpen}
              <svg class="block h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            {:else}
              <svg class="block h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Mobile menu -->
  {#if mobileMenuOpen}
    <div class="sm:hidden">
      <div class="space-y-1 pb-3 pt-2">
        {#if $authStore.isAuthenticated}
          <a
            href="/dashboard"
            class="block py-2 pl-3 pr-4 text-base font-medium text-gray-700 hover:bg-gray-50"
          >
            Dashboard
          </a>
          <a
            href="/dashboard/certificates"
            class="block py-2 pl-3 pr-4 text-base font-medium text-gray-700 hover:bg-gray-50"
          >
            Certificates
          </a>
          <a
            href="/dashboard/encryption"
            class="block py-2 pl-3 pr-4 text-base font-medium text-gray-700 hover:bg-gray-50"
          >
            Encryption
          </a>
          <a
            href="/dashboard/operations"
            class="block py-2 pl-3 pr-4 text-base font-medium text-gray-700 hover:bg-gray-50"
          >
            History
          </a>
          {#if $authStore.user?.role === 'admin'}
            <a
              href="/admin"
              class="block py-2 pl-3 pr-4 text-base font-medium text-gray-700 hover:bg-gray-50"
            >
              Admin
            </a>
          {/if}
          <button
            on:click={handleLogout}
            class="block w-full text-left py-2 pl-3 pr-4 text-base font-medium text-gray-700 hover:bg-gray-50"
          >
            Sign Out
          </button>
        {:else}
          <a
            href="/login"
            class="block py-2 pl-3 pr-4 text-base font-medium text-gray-700 hover:bg-gray-50"
          >
            Sign In
          </a>
          <a
            href="/register"
            class="block py-2 pl-3 pr-4 text-base font-medium text-blue-600 hover:bg-gray-50"
          >
            Sign Up
          </a>
        {/if}
      </div>
    </div>
  {/if}
</nav>