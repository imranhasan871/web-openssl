<script lang="ts">
  import { goto } from '$app/navigation';
</script>

<svelte:head>
  <title>Certificates - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-gray-900">Certificate Operations</h1>
    <p class="mt-1 text-sm text-gray-500">
      Manage SSL/TLS certificates and certificate signing requests
    </p>
  </div>

  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <a
      href="/dashboard/certificates/generate"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-blue-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Generate Certificate</h3>
        <p class="mt-2 text-sm text-gray-500">Create a new self-signed SSL certificate</p>
      </div>
    </a>

    <a
      href="/dashboard/certificates/csr"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-green-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Generate CSR</h3>
        <p class="mt-2 text-sm text-gray-500">Create a certificate signing request</p>
      </div>
    </a>

    <a
      href="/dashboard/certificates/parse"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-purple-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 21h7a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v11m0 5l4.879-4.879m0 0a3 3 0 104.243-4.242 3 3 0 00-4.243 4.242z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Parse Certificate</h3>
        <p class="mt-2 text-sm text-gray-500">View certificate details and information</p>
      </div>
    </a>

    <a
      href="/dashboard/certificates/verify"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-yellow-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Verify Certificate</h3>
        <p class="mt-2 text-sm text-gray-500">Validate certificate authenticity</p>
      </div>
    </a>

    <a
      href="/dashboard/certificates/keys"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-red-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Generate Keys</h3>
        <p class="mt-2 text-sm text-gray-500">Create RSA, ECDSA, or ED25519 keys</p>
      </div>
    </a>

    <a
      href="/dashboard/certificates/convert"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-indigo-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Convert Format</h3>
        <p class="mt-2 text-sm text-gray-500">Convert between PEM, DER, PKCS12 formats</p>
      </div>
    </a>
  </div>
</div>