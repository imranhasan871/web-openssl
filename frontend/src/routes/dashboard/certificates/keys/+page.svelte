<script lang="ts">
  import { apiClient } from '$lib/api/client';
  import { notifications } from '$lib/stores/notifications';
  import Button from '$lib/components/ui/Button.svelte';

  let formData = {
    algorithm: 'RSA',
    keySize: 2048,
    curve: 'prime256v1'
  };

  let loading = false;
  let privateKey = '';
  let publicKey = '';

  async function generateKeys() {
    loading = true;
    try {
      const response = await apiClient.post('/api/v1/openssl/keys/generate', formData);

      if (response.success && response.data) {
        privateKey = response.data.privateKey || '';
        publicKey = response.data.publicKey || '';
        notifications.success('Success', 'Keys generated successfully');
      } else {
        notifications.error('Error', response.error || 'Failed to generate keys');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function copyPrivateKey() {
    navigator.clipboard.writeText(privateKey);
    notifications.success('Copied', 'Private key copied to clipboard');
  }

  function copyPublicKey() {
    navigator.clipboard.writeText(publicKey);
    notifications.success('Copied', 'Public key copied to clipboard');
  }

  function downloadPrivateKey() {
    const blob = new Blob([privateKey], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `private_key_${formData.algorithm.toLowerCase()}.pem`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    notifications.success('Downloaded', 'Private key downloaded');
  }

  function downloadPublicKey() {
    const blob = new Blob([publicKey], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `public_key_${formData.algorithm.toLowerCase()}.pem`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    notifications.success('Downloaded', 'Public key downloaded');
  }

  function clearKeys() {
    privateKey = '';
    publicKey = '';
  }

  $: showKeySize = formData.algorithm === 'RSA';
  $: showCurve = formData.algorithm === 'ECDSA';
</script>

<svelte:head>
  <title>Generate Keys - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <Button href="/dashboard/certificates" variant="outline" size="sm">
      ← Back to Certificates
    </Button>
    <h1 class="text-2xl font-bold text-gray-900 mt-4">Generate Cryptographic Keys</h1>
    <p class="mt-1 text-sm text-gray-500">
      Create RSA, ECDSA, or ED25519 key pairs
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">Key Configuration</h2>

      <form on:submit|preventDefault={generateKeys} class="space-y-4">
        <div>
          <label for="algorithm" class="block text-sm font-medium text-gray-700">Algorithm</label>
          <select
            id="algorithm"
            bind:value={formData.algorithm}
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          >
            <option value="RSA">RSA</option>
            <option value="ECDSA">ECDSA (Elliptic Curve)</option>
            <option value="ED25519">ED25519</option>
          </select>
          <p class="mt-1 text-xs text-gray-500">
            {#if formData.algorithm === 'RSA'}
              RSA is widely supported and suitable for most use cases
            {:else if formData.algorithm === 'ECDSA'}
              ECDSA provides strong security with smaller key sizes
            {:else}
              ED25519 offers excellent performance and security
            {/if}
          </p>
        </div>

        {#if showKeySize}
          <div>
            <label for="keySize" class="block text-sm font-medium text-gray-700">Key Size</label>
            <select
              id="keySize"
              bind:value={formData.keySize}
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            >
              <option value={2048}>2048 bits (Standard)</option>
              <option value={3072}>3072 bits (Secure)</option>
              <option value={4096}>4096 bits (Maximum Security)</option>
            </select>
          </div>
        {/if}

        {#if showCurve}
          <div>
            <label for="curve" class="block text-sm font-medium text-gray-700">Elliptic Curve</label>
            <select
              id="curve"
              bind:value={formData.curve}
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            >
              <option value="prime256v1">prime256v1 (P-256)</option>
              <option value="secp384r1">secp384r1 (P-384)</option>
              <option value="secp521r1">secp521r1 (P-521)</option>
            </select>
          </div>
        {/if}

        <div class="flex gap-2">
          <Button type="submit" disabled={loading} class="flex-1">
            {loading ? 'Generating...' : 'Generate Key Pair'}
          </Button>
          {#if privateKey || publicKey}
            <Button type="button" variant="outline" on:click={clearKeys}>
              Clear
            </Button>
          {/if}
        </div>
      </form>

      <div class="mt-6 bg-yellow-50 border border-yellow-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-yellow-900">Security Warning</h3>
            <p class="text-sm text-yellow-700 mt-1">
              Keep your private key secure. Never share it or commit it to version control.
            </p>
          </div>
        </div>
      </div>
    </div>

    {#if privateKey || publicKey}
      <div class="space-y-6">
        {#if privateKey}
          <div class="bg-white shadow rounded-lg p-6">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-medium text-gray-900">Private Key</h2>
              <div class="flex gap-2">
                <button
                  on:click={copyPrivateKey}
                  class="text-sm text-blue-600 hover:text-blue-800"
                >
                  Copy
                </button>
                <button
                  on:click={downloadPrivateKey}
                  class="text-sm text-blue-600 hover:text-blue-800"
                >
                  Download
                </button>
              </div>
            </div>
            <textarea
              readonly
              value={privateKey}
              rows="12"
              class="block w-full rounded-md border-gray-300 shadow-sm font-mono text-xs bg-gray-50"
            ></textarea>
          </div>
        {/if}

        {#if publicKey}
          <div class="bg-white shadow rounded-lg p-6">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-medium text-gray-900">Public Key</h2>
              <div class="flex gap-2">
                <button
                  on:click={copyPublicKey}
                  class="text-sm text-blue-600 hover:text-blue-800"
                >
                  Copy
                </button>
                <button
                  on:click={downloadPublicKey}
                  class="text-sm text-blue-600 hover:text-blue-800"
                >
                  Download
                </button>
              </div>
            </div>
            <textarea
              readonly
              value={publicKey}
              rows="8"
              class="block w-full rounded-md border-gray-300 shadow-sm font-mono text-xs bg-gray-50"
            ></textarea>
          </div>
        {/if}

        <div class="bg-blue-50 border border-blue-200 rounded-md p-4">
          <h3 class="text-sm font-medium text-blue-900 mb-2">Key Usage</h3>
          <ul class="list-disc list-inside text-sm text-blue-700 space-y-1">
            <li>Use the private key for signing and decryption</li>
            <li>Share the public key for verification and encryption</li>
            <li>Store the private key in a secure location</li>
            <li>Use password protection for private keys in production</li>
          </ul>
        </div>
      </div>
    {/if}
  </div>
</div>