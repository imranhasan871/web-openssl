<script lang="ts">
  import { apiClient } from '$lib/api/client';
  import { notifications } from '$lib/stores/notifications';
  import Button from '$lib/components/ui/Button.svelte';

  let certificateInput = '';
  let caChainInput = '';
  let loading = false;
  let result: any = null;

  async function verifyCertificate() {
    if (!certificateInput.trim()) {
      notifications.error('Validation Error', 'Please enter a certificate');
      return;
    }

    loading = true;
    try {
      const response = await apiClient.post('/api/v1/openssl/certificates/verify', {
        certificate: certificateInput,
        caChain: caChainInput || undefined
      });

      if (response.success && response.data) {
        result = response.data;
        if (result.valid) {
          notifications.success('Valid', 'Certificate is valid');
        } else {
          notifications.error('Invalid', result.error || 'Certificate verification failed');
        }
      } else {
        notifications.error('Error', response.error || 'Failed to verify certificate');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function clearForm() {
    certificateInput = '';
    caChainInput = '';
    result = null;
  }
</script>

<svelte:head>
  <title>Verify Certificate - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <Button href="/dashboard/certificates" variant="outline" size="sm">
      ← Back to Certificates
    </Button>
    <h1 class="text-2xl font-bold text-gray-900 mt-4">Verify Certificate</h1>
    <p class="mt-1 text-sm text-gray-500">
      Validate the authenticity and integrity of an SSL/TLS certificate
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">Certificate Verification</h2>

      <form on:submit|preventDefault={verifyCertificate} class="space-y-4">
        <div>
          <label for="certificate" class="block text-sm font-medium text-gray-700 mb-2">
            Certificate (PEM format) *
          </label>
          <textarea
            id="certificate"
            bind:value={certificateInput}
            rows="10"
            placeholder="-----BEGIN CERTIFICATE-----&#10;...&#10;-----END CERTIFICATE-----"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm font-mono"
            required
          ></textarea>
          <p class="mt-1 text-xs text-gray-500">
            The certificate to verify
          </p>
        </div>

        <div>
          <label for="caChain" class="block text-sm font-medium text-gray-700 mb-2">
            CA Chain (optional)
          </label>
          <textarea
            id="caChain"
            bind:value={caChainInput}
            rows="8"
            placeholder="-----BEGIN CERTIFICATE-----&#10;...&#10;-----END CERTIFICATE-----&#10;-----BEGIN CERTIFICATE-----&#10;...&#10;-----END CERTIFICATE-----"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm font-mono"
          ></textarea>
          <p class="mt-1 text-xs text-gray-500">
            Intermediate and root CA certificates (optional)
          </p>
        </div>

        <div class="flex gap-2">
          <Button type="submit" disabled={loading} class="flex-1">
            {loading ? 'Verifying...' : 'Verify Certificate'}
          </Button>
          <Button type="button" variant="outline" on:click={clearForm}>
            Clear
          </Button>
        </div>
      </form>

      <div class="mt-6 bg-blue-50 border border-blue-200 rounded-md p-4">
        <h3 class="text-sm font-medium text-blue-900 mb-2">Verification Checks</h3>
        <ul class="list-disc list-inside text-sm text-blue-700 space-y-1">
          <li>Certificate signature validation</li>
          <li>Certificate chain verification</li>
          <li>Validity period check</li>
          <li>Issuer validation</li>
          <li>Certificate format validation</li>
        </ul>
      </div>
    </div>

    {#if result}
      <div class="bg-white shadow rounded-lg p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Verification Result</h2>

        <div class="space-y-4">
          <div class={`rounded-md p-4 ${result.valid ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'}`}>
            <div class="flex items-center">
              <svg
                class={`h-6 w-6 ${result.valid ? 'text-green-600' : 'text-red-600'}`}
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                {#if result.valid}
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                {:else}
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                {/if}
              </svg>
              <h3 class={`ml-3 text-lg font-medium ${result.valid ? 'text-green-900' : 'text-red-900'}`}>
                {result.valid ? 'Certificate is Valid' : 'Certificate is Invalid'}
              </h3>
            </div>
          </div>

          {#if result.message}
            <div class="bg-gray-50 rounded-md p-4">
              <h4 class="text-sm font-medium text-gray-700 mb-2">Details</h4>
              <p class="text-sm text-gray-600">{result.message}</p>
            </div>
          {/if}

          {#if result.error}
            <div class="bg-red-50 rounded-md p-4">
              <h4 class="text-sm font-medium text-red-700 mb-2">Error</h4>
              <p class="text-sm text-red-600">{result.error}</p>
            </div>
          {/if}

          {#if result.checks}
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-3">Verification Checks</h4>
              <div class="space-y-2">
                {#each Object.entries(result.checks) as [check, passed]}
                  <div class="flex items-center">
                    <svg
                      class={`h-5 w-5 ${passed ? 'text-green-500' : 'text-red-500'}`}
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      {#if passed}
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                      {:else}
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                      {/if}
                    </svg>
                    <span class="ml-2 text-sm text-gray-700">{check}</span>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          {#if result.subject}
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Certificate Subject</h4>
              <div class="bg-gray-50 rounded-md p-3">
                <p class="text-sm text-gray-600 font-mono">{result.subject}</p>
              </div>
            </div>
          {/if}

          {#if result.issuer}
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Certificate Issuer</h4>
              <div class="bg-gray-50 rounded-md p-3">
                <p class="text-sm text-gray-600 font-mono">{result.issuer}</p>
              </div>
            </div>
          {/if}

          {#if result.validFrom || result.validTo}
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Validity Period</h4>
              <div class="bg-gray-50 rounded-md p-3 space-y-1">
                {#if result.validFrom}
                  <div class="text-sm">
                    <span class="font-medium">From:</span>
                    <span class="text-gray-600 ml-2">{result.validFrom}</span>
                  </div>
                {/if}
                {#if result.validTo}
                  <div class="text-sm">
                    <span class="font-medium">To:</span>
                    <span class="text-gray-600 ml-2">{result.validTo}</span>
                  </div>
                {/if}
              </div>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>