<script lang="ts">
  import { apiClient } from '$lib/api/client';
  import { notifications } from '$lib/stores/notifications';
  import Button from '$lib/components/ui/Button.svelte';

  let formData = {
    commonName: '',
    organization: '',
    organizationalUnit: '',
    city: '',
    state: '',
    country: '',
    email: '',
    keySize: 2048
  };

  let loading = false;
  let result = '';

  async function generateCSR() {
    if (!formData.commonName) {
      notifications.error('Validation Error', 'Common Name is required');
      return;
    }

    loading = true;
    try {
      const response = await apiClient.post('/api/v1/openssl/certificates/csr', formData);

      if (response.success && response.data) {
        result = response.data.csr || '';
        notifications.success('Success', 'CSR generated successfully');
      } else {
        notifications.error('Error', response.error || 'Failed to generate CSR');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function copyToClipboard() {
    navigator.clipboard.writeText(result);
    notifications.success('Copied', 'CSR copied to clipboard');
  }

  function downloadCSR() {
    const blob = new Blob([result], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${formData.commonName.replace(/[^a-z0-9]/gi, '_')}.csr`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    notifications.success('Downloaded', 'CSR downloaded');
  }
</script>

<svelte:head>
  <title>Generate CSR - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <Button href="/dashboard/certificates" variant="outline" size="sm">
      ← Back to Certificates
    </Button>
    <h1 class="text-2xl font-bold text-gray-900 mt-4">Generate Certificate Signing Request</h1>
    <p class="mt-1 text-sm text-gray-500">
      Create a CSR to submit to a Certificate Authority
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">CSR Details</h2>

      <form on:submit|preventDefault={generateCSR} class="space-y-4">
        <div>
          <label for="commonName" class="block text-sm font-medium text-gray-700">Common Name (CN) *</label>
          <input
            type="text"
            id="commonName"
            bind:value={formData.commonName}
            placeholder="example.com"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            required
          />
          <p class="mt-1 text-xs text-gray-500">Domain name or server hostname</p>
        </div>

        <div>
          <label for="organization" class="block text-sm font-medium text-gray-700">Organization (O)</label>
          <input
            type="text"
            id="organization"
            bind:value={formData.organization}
            placeholder="Acme Corporation"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div>
          <label for="organizationalUnit" class="block text-sm font-medium text-gray-700">Organizational Unit (OU)</label>
          <input
            type="text"
            id="organizationalUnit"
            bind:value={formData.organizationalUnit}
            placeholder="IT Department"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="city" class="block text-sm font-medium text-gray-700">City (L)</label>
            <input
              type="text"
              id="city"
              bind:value={formData.city}
              placeholder="San Francisco"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="state" class="block text-sm font-medium text-gray-700">State (ST)</label>
            <input
              type="text"
              id="state"
              bind:value={formData.state}
              placeholder="California"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>
        </div>

        <div>
          <label for="country" class="block text-sm font-medium text-gray-700">Country (C)</label>
          <input
            type="text"
            id="country"
            bind:value={formData.country}
            placeholder="US"
            maxlength="2"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          />
          <p class="mt-1 text-xs text-gray-500">Two-letter country code</p>
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">Email Address</label>
          <input
            type="email"
            id="email"
            bind:value={formData.email}
            placeholder="admin@example.com"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          />
        </div>

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

        <Button type="submit" disabled={loading} class="w-full">
          {loading ? 'Generating...' : 'Generate CSR'}
        </Button>
      </form>
    </div>

    {#if result}
      <div class="bg-white shadow rounded-lg p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Generated CSR</h2>

        <div class="mb-4">
          <div class="flex justify-between items-center mb-2">
            <label class="block text-sm font-medium text-gray-700">Certificate Signing Request</label>
            <div class="flex gap-2">
              <button
                on:click={copyToClipboard}
                class="text-sm text-blue-600 hover:text-blue-800"
              >
                Copy
              </button>
              <button
                on:click={downloadCSR}
                class="text-sm text-blue-600 hover:text-blue-800"
              >
                Download
              </button>
            </div>
          </div>
          <textarea
            readonly
            value={result}
            rows="20"
            class="block w-full rounded-md border-gray-300 shadow-sm font-mono text-xs bg-gray-50"
          ></textarea>
        </div>

        <div class="bg-blue-50 border border-blue-200 rounded-md p-4">
          <h3 class="text-sm font-medium text-blue-900 mb-2">Next Steps</h3>
          <ol class="list-decimal list-inside text-sm text-blue-700 space-y-1">
            <li>Submit this CSR to your Certificate Authority</li>
            <li>Wait for the CA to validate your request</li>
            <li>Download the signed certificate from the CA</li>
            <li>Install the certificate on your server</li>
          </ol>
        </div>
      </div>
    {/if}
  </div>
</div>