<script lang="ts">
  import { notifications } from '$lib/stores/notifications';
  import { apiClient } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let loading = false;
  let mode: 'encrypt' | 'decrypt' = 'encrypt';

  // Form fields
  let data = '';
  let password = '';
  let algorithm = 'aes-256-cbc';
  let result = '';

  async function handleSubmit() {
    if (!data || !password) {
      notifications.error('Validation Error', 'All fields are required');
      return;
    }

    loading = true;
    result = '';

    try {
      const endpoint = mode === 'encrypt'
        ? '/api/v1/openssl/encrypt/symmetric'
        : '/api/v1/openssl/encrypt/decrypt';

      const response = await apiClient.post(endpoint, {
        data,
        password,
        algorithm: mode === 'encrypt' ? algorithm : undefined
      });

      if (response.success && response.data) {
        result = response.data.result || JSON.stringify(response.data, null, 2);
        notifications.success('Success', `Data ${mode}ed successfully!`);
      } else {
        notifications.error('Error', response.error || `Failed to ${mode} data`);
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function copyResult() {
    navigator.clipboard.writeText(result);
    notifications.success('Copied', 'Result copied to clipboard');
  }

  function reset() {
    data = '';
    password = '';
    result = '';
  }
</script>

<svelte:head>
  <title>Symmetric Encryption - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-gray-900">Symmetric Encryption</h1>
    <p class="mt-1 text-sm text-gray-500">
      Encrypt or decrypt data using symmetric algorithms (AES, DES, etc.)
    </p>
  </div>

  <!-- Mode Toggle -->
  <div class="bg-white shadow rounded-lg p-4 mb-6">
    <div class="flex space-x-4">
      <button
        on:click={() => mode = 'encrypt'}
        class="px-4 py-2 rounded-md {mode === 'encrypt' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
      >
        Encrypt
      </button>
      <button
        on:click={() => mode = 'decrypt'}
        class="px-4 py-2 rounded-md {mode === 'decrypt' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
      >
        Decrypt
      </button>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Form -->
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-semibold mb-6">
        {mode === 'encrypt' ? 'Encrypt' : 'Decrypt'} Data
      </h2>

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <div>
          <label for="data" class="block text-sm font-medium text-gray-700 mb-1">
            {mode === 'encrypt' ? 'Plain Text' : 'Encrypted Data'} *
          </label>
          <textarea
            id="data"
            bind:value={data}
            rows="6"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder={mode === 'encrypt' ? 'Enter text to encrypt...' : 'Enter encrypted data to decrypt...'}
            required
            disabled={loading}
          />
        </div>

        <Input
          id="password"
          type="password"
          label="Password *"
          placeholder="Enter password"
          bind:value={password}
          required
          disabled={loading}
        />

        {#if mode === 'encrypt'}
          <div>
            <label for="algorithm" class="block text-sm font-medium text-gray-700 mb-1">
              Algorithm
            </label>
            <select
              id="algorithm"
              bind:value={algorithm}
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            >
              <option value="aes-256-cbc">AES-256-CBC</option>
              <option value="aes-192-cbc">AES-192-CBC</option>
              <option value="aes-128-cbc">AES-128-CBC</option>
              <option value="des3">3DES</option>
              <option value="des">DES</option>
            </select>
          </div>
        {/if}

        <div class="flex space-x-3 pt-4">
          <Button type="submit" {loading} disabled={loading} fullWidth>
            {loading ? 'Processing...' : mode === 'encrypt' ? 'Encrypt' : 'Decrypt'}
          </Button>
          <Button type="button" variant="outline" on:click={reset} disabled={loading}>
            Reset
          </Button>
        </div>
      </form>
    </div>

    <!-- Result -->
    <div class="bg-white shadow rounded-lg p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold">Result</h2>
        {#if result}
          <button
            on:click={copyResult}
            class="text-sm text-blue-600 hover:text-blue-700"
          >
            Copy
          </button>
        {/if}
      </div>

      {#if result}
        <div class="bg-gray-50 rounded-md p-4 overflow-auto max-h-[600px]">
          <pre class="text-xs font-mono whitespace-pre-wrap">{result}</pre>
        </div>
      {:else}
        <div class="bg-gray-50 rounded-md p-8 text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          <p class="mt-2 text-sm text-gray-500">
            {mode === 'encrypt' ? 'Encrypted' : 'Decrypted'} data will appear here
          </p>
        </div>
      {/if}
    </div>
  </div>
</div>