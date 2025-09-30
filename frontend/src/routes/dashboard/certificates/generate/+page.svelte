<script lang="ts">
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import { notifications } from '$lib/stores/notifications';
  import { apiClient } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let loading = false;
  let result = '';

  // Form fields
  let commonName = '';
  let organization = '';
  let organizationalUnit = '';
  let city = '';
  let state = '';
  let country = '';
  let email = '';
  let validDays = 365;
  let keySize = 2048;

  async function generateCertificate() {
    if (!commonName) {
      notifications.error('Validation Error', 'Common Name is required');
      return;
    }

    loading = true;
    result = '';

    try {
      const response = await apiClient.post('/api/v1/openssl/certificates/generate', {
        commonName,
        organization,
        organizationalUnit,
        city,
        state,
        country,
        email,
        validDays,
        keySize
      });

      if (response.success && response.data) {
        result = JSON.stringify(response.data, null, 2);
        notifications.success('Success', 'Certificate generated successfully!');
      } else {
        notifications.error('Error', response.error || 'Failed to generate certificate');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function copyToClipboard() {
    navigator.clipboard.writeText(result);
    notifications.success('Copied', 'Certificate copied to clipboard');
  }

  function downloadCertificate() {
    const blob = new Blob([result], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `certificate-${Date.now()}.json`;
    a.click();
    URL.revokeObjectURL(url);
    notifications.success('Downloaded', 'Certificate downloaded');
  }

  function reset() {
    commonName = '';
    organization = '';
    organizationalUnit = '';
    city = '';
    state = '';
    country = '';
    email = '';
    validDays = 365;
    keySize = 2048;
    result = '';
  }
</script>

<svelte:head>
  <title>Generate Certificate - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-gray-900">Generate SSL Certificate</h1>
    <p class="mt-1 text-sm text-gray-500">
      Create a self-signed SSL/TLS certificate with custom parameters
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <!-- Form -->
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-semibold mb-6">Certificate Details</h2>

      <form on:submit|preventDefault={generateCertificate} class="space-y-4">
        <Input
          id="commonName"
          label="Common Name (CN) *"
          placeholder="example.com"
          bind:value={commonName}
          required
          disabled={loading}
        />

        <Input
          id="organization"
          label="Organization (O)"
          placeholder="My Company"
          bind:value={organization}
          disabled={loading}
        />

        <Input
          id="organizationalUnit"
          label="Organizational Unit (OU)"
          placeholder="IT Department"
          bind:value={organizationalUnit}
          disabled={loading}
        />

        <div class="grid grid-cols-2 gap-4">
          <Input
            id="city"
            label="City (L)"
            placeholder="San Francisco"
            bind:value={city}
            disabled={loading}
          />

          <Input
            id="state"
            label="State (ST)"
            placeholder="California"
            bind:value={state}
            disabled={loading}
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <Input
            id="country"
            label="Country (C)"
            placeholder="US"
            bind:value={country}
            maxlength="2"
            disabled={loading}
          />

          <Input
            id="email"
            label="Email"
            type="email"
            placeholder="admin@example.com"
            bind:value={email}
            disabled={loading}
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="validDays" class="block text-sm font-medium text-gray-700 mb-1">
              Valid Days
            </label>
            <input
              id="validDays"
              type="number"
              bind:value={validDays}
              min="1"
              max="3650"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          <div>
            <label for="keySize" class="block text-sm font-medium text-gray-700 mb-1">
              Key Size
            </label>
            <select
              id="keySize"
              bind:value={keySize}
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            >
              <option value={2048}>2048 bits</option>
              <option value={3072}>3072 bits</option>
              <option value={4096}>4096 bits</option>
            </select>
          </div>
        </div>

        <div class="flex space-x-3 pt-4">
          <Button type="submit" {loading} disabled={loading} fullWidth>
            {loading ? 'Generating...' : 'Generate Certificate'}
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
          <div class="flex space-x-2">
            <button
              on:click={copyToClipboard}
              class="text-sm text-blue-600 hover:text-blue-700"
              title="Copy to clipboard"
            >
              Copy
            </button>
            <button
              on:click={downloadCertificate}
              class="text-sm text-blue-600 hover:text-blue-700"
              title="Download"
            >
              Download
            </button>
          </div>
        {/if}
      </div>

      {#if result}
        <div class="bg-gray-50 rounded-md p-4 overflow-auto max-h-[600px]">
          <pre class="text-xs font-mono whitespace-pre-wrap">{result}</pre>
        </div>
      {:else}
        <div class="bg-gray-50 rounded-md p-8 text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
          </svg>
          <p class="mt-2 text-sm text-gray-500">
            Generated certificate will appear here
          </p>
        </div>
      {/if}
    </div>
  </div>
</div>