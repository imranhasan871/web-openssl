<script lang="ts">
  import { goto } from '$app/navigation';
  import { authStore } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';
  import { authAPI } from '$lib/api/auth';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let email = '';
  let password = '';
  let confirmPassword = '';
  let firstName = '';
  let lastName = '';
  let loading = false;
  let acceptedTerms = false;

  async function handleRegister() {
    // Validation
    if (!email || !password || !firstName || !lastName) {
      notifications.error('Validation Error', 'Please fill in all required fields');
      return;
    }

    if (password !== confirmPassword) {
      notifications.error('Validation Error', 'Passwords do not match');
      return;
    }

    if (password.length < 6) {
      notifications.error('Validation Error', 'Password must be at least 6 characters long');
      return;
    }

    if (!acceptedTerms) {
      notifications.error('Terms Required', 'Please accept the terms and conditions');
      return;
    }

    loading = true;

    try {
      const response = await authAPI.register({
        email,
        password,
        firstName,
        lastName
      });

      if (response.success && response.data) {
        authStore.login(response.data.user, response.data.accessToken);
        notifications.success('Welcome!', 'Account created successfully');
        goto('/dashboard');
      } else {
        notifications.error('Registration Failed', response.error || 'Failed to create account');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Sign Up - OpenSSL UI</title>
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
      <h1 class="text-2xl font-semibold mb-2">Create your account</h1>
      <p class="text-gray-600">Start managing SSL certificates today</p>
    </div>

    <!-- Registration Form -->
    <div class="bg-white p-6 rounded-lg shadow-lg">
      <h2 class="text-xl font-semibold mb-6">Sign Up</h2>

      <form on:submit|preventDefault={handleRegister} class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <Input
            id="firstName"
            label="First Name"
            bind:value={firstName}
            placeholder="John"
            required
            disabled={loading}
          />

          <Input
            id="lastName"
            label="Last Name"
            bind:value={lastName}
            placeholder="Doe"
            required
            disabled={loading}
          />
        </div>

        <Input
          id="email"
          type="email"
          label="Email"
          bind:value={email}
          placeholder="john@company.com"
          required
          disabled={loading}
        />

        <Input
          id="password"
          type="password"
          label="Password"
          bind:value={password}
          placeholder="At least 6 characters"
          hint="Must be at least 6 characters long"
          required
          disabled={loading}
        />

        <Input
          id="confirmPassword"
          type="password"
          label="Confirm Password"
          bind:value={confirmPassword}
          placeholder="Repeat your password"
          required
          disabled={loading}
        />

        <div class="flex items-start">
          <input
            id="terms"
            type="checkbox"
            bind:checked={acceptedTerms}
            disabled={loading}
            class="mt-1 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label for="terms" class="ml-2 text-sm text-gray-600">
            I accept the
            <a href="/terms" class="text-blue-600 hover:underline">Terms and Conditions</a>
            and
            <a href="/privacy" class="text-blue-600 hover:underline">Privacy Policy</a>
          </label>
        </div>

        <Button
          type="submit"
          fullWidth
          {loading}
          disabled={loading}
        >
          Create Account
        </Button>
      </form>

      <div class="mt-6 text-center text-sm text-gray-600">
        Already have an account?
        <a href="/login" class="text-blue-600 hover:underline font-medium">
          Sign in
        </a>
      </div>
    </div>

    <!-- Features Preview -->
    <div class="mt-6 bg-white p-4 rounded-lg shadow">
      <h3 class="text-sm font-medium text-gray-900 mb-2">What you get:</h3>
      <ul class="text-sm text-gray-600 space-y-1">
        <li class="flex items-center">
          <svg class="h-4 w-4 text-green-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          SSL Certificate Generation
        </li>
        <li class="flex items-center">
          <svg class="h-4 w-4 text-green-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          Encryption & Security Tools
        </li>
        <li class="flex items-center">
          <svg class="h-4 w-4 text-green-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          API Access & Integration
        </li>
      </ul>
    </div>
  </div>
</div>