<script lang="ts">
  import { notifications } from '$lib/stores/notifications';
  import { apiClient } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let loading = false;
  let mode: 'hash' | 'hmac' = 'hash';

  // Form fields
  let data = '';
  let key = '';
  let algorithm = 'sha256';
  let result = '';

  async function handleSubmit() {
    if (!data) {
      notifications.error('Validation Error', 'Data is required');
      return;
    }

    if (mode === 'hmac' && !key) {
      notifications.error('Validation Error', 'Key is required for HMAC');
      return;
    }

    loading = true;
    result = '';

    try {
      const endpoint = mode === 'hash'
        ? '/api/v1/openssl/hash/generate'
        : '/api/v1/openssl/hash/hmac';

      const payload = mode === 'hash'
        ? { data, algorithm }
        : { data, key, algorithm };

      const response = await apiClient.post(endpoint, payload);

      if (response.success && response.data) {
        result = response.data.hash || response.data.hmac || JSON.stringify(response.data, null, 2);
        notifications.success('Success', `${mode === 'hash' ? 'Hash' : 'HMAC'} generated successfully!`);
      } else {
        notifications.error('Error', response.error || 'Failed to generate');
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
    key = '';
    result = '';
  }
</script>

<svelte:head>
  <title>Hash Generator - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-gray-900">Hash Generator</h1>
    <p class="mt-1 text-sm text-gray-500">
      Generate cryptographic hashes (MD5, SHA-256, etc.) and HMACs
    </p>
  </div>

  <!-- Mode Toggle -->
  <div class="bg-white shadow rounded-lg p-4 mb-6">
    <div class="flex space-x-4">
      <button
        on:click={() => mode = 'hash'}
        class="px-4 py-2 rounded-md {mode === 'hash' ? 'bg-yellow-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
      >
        Hash
      </button>
      <button
        on:click={() => mode = 'hmac'}
        class="px-4 py-2 rounded-md {mode === 'hmac' ? 'bg-yellow-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
      >
        HMAC
      </button>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Form -->
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-semibold mb-6">
        Generate {mode === 'hash' ? 'Hash' : 'HMAC'}
      </h2>

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <div>
          <label for="data" class="block text-sm font-medium text-gray-700 mb-1">
            Data *
          </label>
          <textarea
            id="data"
            bind:value={data}
            rows="6"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-yellow-500"
            placeholder="Enter text to hash..."
            required
            disabled={loading}
          />
        </div>

        {#if mode === 'hmac'}
          <Input
            id="key"
            label="Key *"
            placeholder="Enter secret key"
            bind:value={key}
            required
            disabled={loading}
          />
        {/if}

        <div>
          <label for="algorithm" class="block text-sm font-medium text-gray-700 mb-1">
            Algorithm
          </label>
          <select
            id="algorithm"
            bind:value={algorithm}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-yellow-500"
            disabled={loading}
          >
            <option value="md5">MD5</option>
            <option value="sha1">SHA-1</option>
            <option value="sha256">SHA-256</option>
            <option value="sha384">SHA-384</option>
            <option value="sha512">SHA-512</option>
          </select>
        </div>

        <div class="flex space-x-3 pt-4">
          <Button type="submit" {loading} disabled={loading} fullWidth>
            {loading ? 'Generating...' : 'Generate'}
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
            class="text-sm text-yellow-600 hover:text-yellow-700"
          >
            Copy
          </button>
        {/if}
      </div>

      {#if result}
        <div class="bg-gray-50 rounded-md p-4 overflow-auto">
          <pre class="text-xs font-mono break-all whitespace-pre-wrap">{result}</pre>
        </div>

        <div class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
          <p class="text-xs text-yellow-800">
            <strong>Length:</strong> {result.length} characters
          </p>
        </div>
      {:else}
        <div class="bg-gray-50 rounded-md p-8 text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2h3a1 1 0 011 1v1a1 1 0 01-1 1h-1v9a2 2 0 01-2 2H7a2 2 0 01-2-2V7H4a1 1 0 01-1-1V5a1 1 0 011-1h3z"/>
          </svg>
          <p class="mt-2 text-sm text-gray-500">
            Generated hash will appear here
          </p>
        </div>
      {/if}
    </div>
  </div>
</div>