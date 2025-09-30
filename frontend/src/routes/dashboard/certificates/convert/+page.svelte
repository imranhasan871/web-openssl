<script lang="ts">
  import { apiClient } from '$lib/api/client';
  import { notifications } from '$lib/stores/notifications';
  import Button from '$lib/components/ui/Button.svelte';

  let formData = {
    inputFormat: 'PEM',
    outputFormat: 'DER',
    certificateInput: '',
    password: ''
  };

  let loading = false;
  let result = '';
  let resultBinary = false;

  async function convertCertificate() {
    if (!formData.certificateInput.trim()) {
      notifications.error('Validation Error', 'Please enter a certificate');
      return;
    }

    loading = true;
    try {
      const response = await apiClient.post('/api/v1/openssl/certificates/convert', {
        certificate: formData.certificateInput,
        inputFormat: formData.inputFormat,
        outputFormat: formData.outputFormat,
        password: formData.password || undefined
      });

      if (response.success && response.data) {
        result = response.data.certificate || '';
        resultBinary = formData.outputFormat === 'DER' || formData.outputFormat === 'PKCS12';
        notifications.success('Success', 'Certificate converted successfully');
      } else {
        notifications.error('Error', response.error || 'Failed to convert certificate');
      }
    } catch (error) {
      notifications.error('Error', 'An unexpected error occurred');
    } finally {
      loading = false;
    }
  }

  function copyToClipboard() {
    if (!resultBinary) {
      navigator.clipboard.writeText(result);
      notifications.success('Copied', 'Converted certificate copied to clipboard');
    }
  }

  function downloadResult() {
    const extension = formData.outputFormat.toLowerCase();
    const filename = `certificate.${extension === 'pkcs12' ? 'p12' : extension}`;

    if (resultBinary) {
      // For binary formats, the result should be base64 encoded
      const binaryData = atob(result);
      const bytes = new Uint8Array(binaryData.length);
      for (let i = 0; i < binaryData.length; i++) {
        bytes[i] = binaryData.charCodeAt(i);
      }
      const blob = new Blob([bytes], { type: 'application/octet-stream' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } else {
      const blob = new Blob([result], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    }
    notifications.success('Downloaded', 'Certificate downloaded');
  }

  function clearForm() {
    formData.certificateInput = '';
    formData.password = '';
    result = '';
  }

  $: needsPassword = formData.outputFormat === 'PKCS12';
</script>

<svelte:head>
  <title>Convert Certificate - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <Button href="/dashboard/certificates" variant="outline" size="sm">
      ← Back to Certificates
    </Button>
    <h1 class="text-2xl font-bold text-gray-900 mt-4">Convert Certificate Format</h1>
    <p class="mt-1 text-sm text-gray-500">
      Convert certificates between PEM, DER, and PKCS12 formats
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">Conversion Settings</h2>

      <form on:submit|preventDefault={convertCertificate} class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="inputFormat" class="block text-sm font-medium text-gray-700">Input Format</label>
            <select
              id="inputFormat"
              bind:value={formData.inputFormat}
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            >
              <option value="PEM">PEM</option>
              <option value="DER">DER</option>
              <option value="PKCS12">PKCS12</option>
            </select>
          </div>

          <div>
            <label for="outputFormat" class="block text-sm font-medium text-gray-700">Output Format</label>
            <select
              id="outputFormat"
              bind:value={formData.outputFormat}
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            >
              <option value="PEM">PEM</option>
              <option value="DER">DER</option>
              <option value="PKCS12">PKCS12</option>
            </select>
          </div>
        </div>

        <div>
          <label for="certificate" class="block text-sm font-medium text-gray-700 mb-2">
            Certificate
          </label>
          <textarea
            id="certificate"
            bind:value={formData.certificateInput}
            rows="12"
            placeholder={formData.inputFormat === 'PEM'
              ? '-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----'
              : 'Paste your certificate content here'}
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm font-mono"
            required
          ></textarea>
          <p class="mt-1 text-xs text-gray-500">
            {#if formData.inputFormat === 'PEM'}
              PEM format (Base64 encoded with BEGIN/END headers)
            {:else if formData.inputFormat === 'DER'}
              DER format (Binary, base64 encoded for input)
            {:else}
              PKCS12 format (Base64 encoded .p12/.pfx file)
            {/if}
          </p>
        </div>

        {#if needsPassword}
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
            <input
              type="password"
              id="password"
              bind:value={formData.password}
              placeholder="Enter password for PKCS12"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
            <p class="mt-1 text-xs text-gray-500">
              Required for PKCS12 output format
            </p>
          </div>
        {/if}

        <div class="flex gap-2">
          <Button type="submit" disabled={loading} class="flex-1">
            {loading ? 'Converting...' : 'Convert Certificate'}
          </Button>
          <Button type="button" variant="outline" on:click={clearForm}>
            Clear
          </Button>
        </div>
      </form>

      <div class="mt-6 bg-blue-50 border border-blue-200 rounded-md p-4">
        <h3 class="text-sm font-medium text-blue-900 mb-2">Format Information</h3>
        <dl class="text-sm text-blue-700 space-y-2">
          <div>
            <dt class="font-medium">PEM:</dt>
            <dd class="ml-4">Base64 encoded, human-readable, most common</dd>
          </div>
          <div>
            <dt class="font-medium">DER:</dt>
            <dd class="ml-4">Binary format, compact, used in Java</dd>
          </div>
          <div>
            <dt class="font-medium">PKCS12:</dt>
            <dd class="ml-4">Binary archive containing certificate and private key</dd>
          </div>
        </dl>
      </div>
    </div>

    {#if result}
      <div class="bg-white shadow rounded-lg p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Converted Certificate</h2>

        <div class="mb-4">
          <div class="flex justify-between items-center mb-2">
            <label class="block text-sm font-medium text-gray-700">
              {formData.outputFormat} Format
            </label>
            <div class="flex gap-2">
              {#if !resultBinary}
                <button
                  on:click={copyToClipboard}
                  class="text-sm text-blue-600 hover:text-blue-800"
                >
                  Copy
                </button>
              {/if}
              <button
                on:click={downloadResult}
                class="text-sm text-blue-600 hover:text-blue-800"
              >
                Download
              </button>
            </div>
          </div>

          {#if resultBinary}
            <div class="bg-gray-50 border border-gray-300 rounded-md p-4">
              <div class="flex items-center justify-center py-8">
                <div class="text-center">
                  <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
                  </svg>
                  <p class="mt-2 text-sm text-gray-600">Binary file format</p>
                  <p class="text-xs text-gray-500">Click Download to save the file</p>
                </div>
              </div>
            </div>
          {:else}
            <textarea
              readonly
              value={result}
              rows="15"
              class="block w-full rounded-md border-gray-300 shadow-sm font-mono text-xs bg-gray-50"
            ></textarea>
          {/if}
        </div>

        <div class="bg-green-50 border border-green-200 rounded-md p-4">
          <div class="flex">
            <svg class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-green-900">Conversion Complete</h3>
              <p class="text-sm text-green-700 mt-1">
                Your certificate has been converted to {formData.outputFormat} format.
              </p>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>