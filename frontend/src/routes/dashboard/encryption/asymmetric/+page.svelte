<script lang="ts">
  import { apiClient } from '$lib/api/client';
  import { notifications } from '$lib/stores/notifications';
  import Button from '$lib/components/ui/Button.svelte';

  let mode: 'encrypt' | 'decrypt' = 'encrypt';
  let formData = {
    data: '',
    key: '',
    keyType: 'public' as 'public' | 'private'
  };

  let loading = false;
  let result = '';

  async function processData() {
    if (!formData.data.trim()) {
      notifications.error('Validation Error', 'Please enter data to process');
      return;
    }

    if (!formData.key.trim()) {
      notifications.error('Validation Error', 'Please provide a key');
      return;
    }

    loading = true;
    try {
      const endpoint = mode === 'encrypt'
        ? '/api/v1/openssl/encryption/asymmetric/encrypt'
        : '/api/v1/openssl/encryption/asymmetric/decrypt';

      const response = await apiClient.post(endpoint, {
        data: formData.data,
        key: formData.key
      });

      if (response.success && response.data) {
        result = response.data.result || '';
        notifications.success('Success', `Data ${mode}ed successfully`);
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

  function clearForm() {
    formData.data = '';
    result = '';
  }

  function switchMode() {
    mode = mode === 'encrypt' ? 'decrypt' : 'encrypt';
    formData.keyType = mode === 'encrypt' ? 'public' : 'private';
    clearForm();
  }

  $: formData.keyType = mode === 'encrypt' ? 'public' : 'private';
</script>

<svelte:head>
  <title>Asymmetric Encryption - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <Button href="/dashboard/encryption" variant="outline" size="sm">
      ← Back to Encryption
    </Button>
    <h1 class="text-2xl font-bold text-gray-900 mt-4">Asymmetric Encryption</h1>
    <p class="mt-1 text-sm text-gray-500">
      Encrypt and decrypt data using RSA public/private key pairs
    </p>
  </div>

  <div class="mb-6">
    <div class="flex space-x-2">
      <button
        on:click={() => mode = 'encrypt'}
        class={`px-4 py-2 rounded-md font-medium text-sm transition-colors ${
          mode === 'encrypt'
            ? 'bg-blue-600 text-white'
            : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
        }`}
      >
        Encrypt
      </button>
      <button
        on:click={() => mode = 'decrypt'}
        class={`px-4 py-2 rounded-md font-medium text-sm transition-colors ${
          mode === 'decrypt'
            ? 'bg-blue-600 text-white'
            : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
        }`}
      >
        Decrypt
      </button>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">
        {mode === 'encrypt' ? 'Encrypt' : 'Decrypt'} Data
      </h2>

      <form on:submit|preventDefault={processData} class="space-y-4">
        <div>
          <label for="data" class="block text-sm font-medium text-gray-700 mb-2">
            {mode === 'encrypt' ? 'Data to Encrypt' : 'Encrypted Data'}
          </label>
          <textarea
            id="data"
            bind:value={formData.data}
            rows="6"
            placeholder={mode === 'encrypt' ? 'Enter text to encrypt' : 'Enter base64 encoded encrypted data'}
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            required
          ></textarea>
        </div>

        <div>
          <label for="key" class="block text-sm font-medium text-gray-700 mb-2">
            {mode === 'encrypt' ? 'Public Key' : 'Private Key'} (PEM format)
          </label>
          <textarea
            id="key"
            bind:value={formData.key}
            rows="10"
            placeholder={mode === 'encrypt'
              ? '-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----'
              : '-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----'}
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm font-mono"
            required
          ></textarea>
          <p class="mt-1 text-xs text-gray-500">
            {mode === 'encrypt'
              ? 'Paste the public key (can be shared publicly)'
              : 'Paste the private key (keep this secure)'}
          </p>
        </div>

        <div class="flex gap-2">
          <Button type="submit" disabled={loading} class="flex-1">
            {loading ? 'Processing...' : mode === 'encrypt' ? 'Encrypt' : 'Decrypt'}
          </Button>
          <Button type="button" variant="outline" on:click={clearForm}>
            Clear
          </Button>
        </div>
      </form>

      <div class="mt-6 bg-blue-50 border border-blue-200 rounded-md p-4">
        <h3 class="text-sm font-medium text-blue-900 mb-2">How it works</h3>
        <ul class="list-disc list-inside text-sm text-blue-700 space-y-1">
          <li>Public key: Used for encryption (can be shared)</li>
          <li>Private key: Used for decryption (must be kept secret)</li>
          <li>RSA encryption is suitable for small amounts of data</li>
          <li>For large data, use symmetric encryption with RSA for key exchange</li>
        </ul>
      </div>
    </div>

    {#if result}
      <div class="bg-white shadow rounded-lg p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Result</h2>

        <div class="mb-4">
          <div class="flex justify-between items-center mb-2">
            <label class="block text-sm font-medium text-gray-700">
              {mode === 'encrypt' ? 'Encrypted Data' : 'Decrypted Data'}
            </label>
            <button
              on:click={copyResult}
              class="text-sm text-blue-600 hover:text-blue-800"
            >
              Copy
            </button>
          </div>
          <textarea
            readonly
            value={result}
            rows="10"
            class="block w-full rounded-md border-gray-300 shadow-sm font-mono text-sm bg-gray-50"
          ></textarea>
          {#if mode === 'encrypt'}
            <p class="mt-2 text-xs text-gray-500">
              The encrypted data is base64 encoded. Share this with someone who has the private key to decrypt it.
            </p>
          {/if}
        </div>

        <div class="bg-green-50 border border-green-200 rounded-md p-4">
          <div class="flex">
            <svg class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-green-900">
                {mode === 'encrypt' ? 'Encryption' : 'Decryption'} Complete
              </h3>
              <p class="text-sm text-green-700 mt-1">
                {mode === 'encrypt'
                  ? 'Your data has been encrypted and can only be decrypted with the corresponding private key.'
                  : 'Your data has been successfully decrypted.'}
              </p>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>