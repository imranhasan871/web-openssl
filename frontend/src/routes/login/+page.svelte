<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { authStore } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';
  import { authAPI } from '$lib/api/auth';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let email = '';
  let password = '';
  let loading = false;

  // Redirect to dashboard if already logged in
  onMount(() => {
    const unsubscribe = authStore.subscribe((auth) => {
      if (!auth.loading && auth.isAuthenticated) {
        goto('/dashboard', { replaceState: true });
      }
    });
    return unsubscribe;
  });

  async function handleLogin() {
    if (!email || !password) {
      notifications.error('Validation Error', 'Please enter both email and password');
      return;
    }

    loading = true;

    try {
      const response = await authAPI.login({ email, password });

      if (response.success && response.data) {
        authStore.login(response.data.user, response.data.accessToken);
        notifications.success('Success', 'Successfully signed in!');
        // Add small delay to ensure store is updated
        await new Promise(resolve => setTimeout(resolve, 100));
        goto('/dashboard');
      } else {
        notifications.error('Login Failed', response.error || 'Invalid credentials');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<svelte:head>
  <title>Sign In - OpenSSL UI</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50 p-4">
  <div class="w-full max-w-md">
    <!-- Header -->
    <div class="text-center mb-8">
      <div class="flex items-center justify-center space-x-2 mb-4">
        <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
        </svg>
        <span class="text-2xl font-bold">OpenSSL UI</span>
      </div>
      <h1 class="text-2xl font-semibold mb-2">Welcome back</h1>
      <p class="text-gray-600">Sign in to your account to continue</p>
    </div>

    <!-- Login Form -->
    <div class="bg-white p-6 rounded-lg shadow-lg">
      <h2 class="text-xl font-semibold mb-6">Sign In</h2>

      <form on:submit|preventDefault={handleLogin} class="space-y-4">
        <Input
          id="email"
          type="email"
          label="Email"
          bind:value={email}
          placeholder="Enter your email"
          required
          disabled={loading}
          on:keydown={handleKeyDown}
        />

        <Input
          id="password"
          type="password"
          label="Password"
          bind:value={password}
          placeholder="Enter your password"
          required
          disabled={loading}
          on:keydown={handleKeyDown}
        />

        <div class="flex items-center justify-between">
          <a
            href="/forgot-password"
            class="text-sm text-blue-600 hover:underline"
          >
            Forgot password?
          </a>
        </div>

        <Button
          type="submit"
          fullWidth
          {loading}
          disabled={loading}
        >
          Sign In
        </Button>
      </form>

      <div class="mt-6 text-center text-sm text-gray-600">
        Don't have an account?
        <a href="/register" class="text-blue-600 hover:underline font-medium">
          Sign up
        </a>
      </div>

      <!-- Demo Credentials -->
      <div class="mt-4 p-3 bg-gray-50 rounded-md text-center">
        <p class="text-sm text-gray-600 mb-1">Demo credentials:</p>
        <p class="text-xs font-mono">demo@opensslui.com / demo123</p>
        <button
          type="button"
          class="mt-2 text-xs text-blue-600 hover:underline"
          on:click={() => {
            email = 'demo@opensslui.com';
            password = 'demo123';
          }}
        >
          Use demo credentials
        </button>
      </div>
    </div>
  </div>
</div>