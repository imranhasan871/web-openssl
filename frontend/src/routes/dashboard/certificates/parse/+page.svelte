<script lang="ts">
  import { apiClient } from '$lib/api/client';
  import { notifications } from '$lib/stores/notifications';
  import Button from '$lib/components/ui/Button.svelte';

  let certificateInput = '';
  let loading = false;
  let result: any = null;

  async function parseCertificate() {
    if (!certificateInput.trim()) {
      notifications.error('Validation Error', 'Please enter a certificate');
      return;
    }

    loading = true;
    try {
      const response = await apiClient.post('/api/v1/openssl/certificates/parse', {
        certificate: certificateInput
      });

      if (response.success && response.data) {
        result = response.data;
        notifications.success('Success', 'Certificate parsed successfully');
      } else {
        notifications.error('Error', response.error || 'Failed to parse certificate');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function loadSampleCertificate() {
    certificateInput = `-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKL0UG+mRKSzMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTkwNzE2MDQzMzM0WhcNMjkwNzEzMDQzMzM0WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAy4gPtMTKU7IrWwLwIGfmHnIcWOLVtR0LJQJ9yZwGQAAMEJC9uOXp5TaS
-----END CERTIFICATE-----`;
  }

  function clearForm() {
    certificateInput = '';
    result = null;
  }
</script>

<svelte:head>
  <title>Parse Certificate - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <Button href="/dashboard/certificates" variant="outline" size="sm">
      ← Back to Certificates
    </Button>
    <h1 class="text-2xl font-bold text-gray-900 mt-4">Parse Certificate</h1>
    <p class="mt-1 text-sm text-gray-500">
      View detailed information from an SSL/TLS certificate
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">Certificate Input</h2>

      <form on:submit|preventDefault={parseCertificate} class="space-y-4">
        <div>
          <div class="flex justify-between items-center mb-2">
            <label for="certificate" class="block text-sm font-medium text-gray-700">
              Certificate (PEM format)
            </label>
            <button
              type="button"
              on:click={loadSampleCertificate}
              class="text-sm text-blue-600 hover:text-blue-800"
            >
              Load Sample
            </button>
          </div>
          <textarea
            id="certificate"
            bind:value={certificateInput}
            rows="15"
            placeholder="-----BEGIN CERTIFICATE-----&#10;...&#10;-----END CERTIFICATE-----"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm font-mono"
            required
          ></textarea>
          <p class="mt-1 text-xs text-gray-500">
            Paste your certificate in PEM format (including BEGIN/END lines)
          </p>
        </div>

        <div class="flex gap-2">
          <Button type="submit" disabled={loading} class="flex-1">
            {loading ? 'Parsing...' : 'Parse Certificate'}
          </Button>
          <Button type="button" variant="outline" on:click={clearForm}>
            Clear
          </Button>
        </div>
      </form>
    </div>

    {#if result}
      <div class="bg-white shadow rounded-lg p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Certificate Details</h2>

        <div class="space-y-4">
          <div>
            <h3 class="text-sm font-medium text-gray-700 mb-2">Subject</h3>
            <div class="bg-gray-50 rounded-md p-3 space-y-1">
              {#if result.subject}
                {#each Object.entries(result.subject) as [key, value]}
                  <div class="text-sm">
                    <span class="font-medium">{key}:</span>
                    <span class="text-gray-600 ml-2">{value}</span>
                  </div>
                {/each}
              {:else}
                <p class="text-sm text-gray-500">No subject information</p>
              {/if}
            </div>
          </div>

          <div>
            <h3 class="text-sm font-medium text-gray-700 mb-2">Issuer</h3>
            <div class="bg-gray-50 rounded-md p-3 space-y-1">
              {#if result.issuer}
                {#each Object.entries(result.issuer) as [key, value]}
                  <div class="text-sm">
                    <span class="font-medium">{key}:</span>
                    <span class="text-gray-600 ml-2">{value}</span>
                  </div>
                {/each}
              {:else}
                <p class="text-sm text-gray-500">No issuer information</p>
              {/if}
            </div>
          </div>

          <div>
            <h3 class="text-sm font-medium text-gray-700 mb-2">Validity</h3>
            <div class="bg-gray-50 rounded-md p-3 space-y-1">
              <div class="text-sm">
                <span class="font-medium">Not Before:</span>
                <span class="text-gray-600 ml-2">{result.notBefore || 'N/A'}</span>
              </div>
              <div class="text-sm">
                <span class="font-medium">Not After:</span>
                <span class="text-gray-600 ml-2">{result.notAfter || 'N/A'}</span>
              </div>
            </div>
          </div>

          <div>
            <h3 class="text-sm font-medium text-gray-700 mb-2">Technical Details</h3>
            <div class="bg-gray-50 rounded-md p-3 space-y-1">
              <div class="text-sm">
                <span class="font-medium">Serial Number:</span>
                <span class="text-gray-600 ml-2 font-mono text-xs">{result.serialNumber || 'N/A'}</span>
              </div>
              <div class="text-sm">
                <span class="font-medium">Version:</span>
                <span class="text-gray-600 ml-2">{result.version || 'N/A'}</span>
              </div>
              <div class="text-sm">
                <span class="font-medium">Signature Algorithm:</span>
                <span class="text-gray-600 ml-2">{result.signatureAlgorithm || 'N/A'}</span>
              </div>
              <div class="text-sm">
                <span class="font-medium">Public Key Algorithm:</span>
                <span class="text-gray-600 ml-2">{result.publicKeyAlgorithm || 'N/A'}</span>
              </div>
            </div>
          </div>

          {#if result.extensions}
            <div>
              <h3 class="text-sm font-medium text-gray-700 mb-2">Extensions</h3>
              <div class="bg-gray-50 rounded-md p-3 space-y-1">
                {#each Object.entries(result.extensions) as [key, value]}
                  <div class="text-sm">
                    <span class="font-medium">{key}:</span>
                    <span class="text-gray-600 ml-2">{value}</span>
                  </div>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>