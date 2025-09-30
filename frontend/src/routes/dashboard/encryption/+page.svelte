<script lang="ts">
</script>

<svelte:head>
  <title>Encryption - OpenSSL UI</title>
</svelte:head>

<div class="px-4 sm:px-6 lg:px-8 py-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-gray-900">Encryption Operations</h1>
    <p class="mt-1 text-sm text-gray-500">
      Encrypt, decrypt, and hash data using OpenSSL
    </p>
  </div>

  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <a
      href="/dashboard/encryption/symmetric"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-green-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-green-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Symmetric Encryption</h3>
        <p class="mt-2 text-sm text-gray-500">AES, DES, and other symmetric algorithms</p>
      </div>
    </a>

    <a
      href="/dashboard/encryption/asymmetric"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-blue-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Asymmetric Encryption</h3>
        <p class="mt-2 text-sm text-gray-500">RSA public/private key encryption</p>
      </div>
    </a>

    <a
      href="/dashboard/encryption/hash"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-yellow-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-yellow-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2h3a1 1 0 011 1v1a1 1 0 01-1 1h-1v9a2 2 0 01-2 2H7a2 2 0 01-2-2V7H4a1 1 0 01-1-1V5a1 1 0 011-1h3z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Hash Generator</h3>
        <p class="mt-2 text-sm text-gray-500">MD5, SHA-1, SHA-256, SHA-512, and HMAC</p>
      </div>
    </a>

    <a
      href="/dashboard/encryption/decrypt"
      class="relative group bg-white p-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-purple-500 rounded-lg shadow hover:shadow-md transition-shadow"
    >
      <div>
        <span class="rounded-lg inline-flex p-3 bg-purple-500 text-white">
          <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"/>
          </svg>
        </span>
      </div>
      <div class="mt-4">
        <h3 class="text-lg font-medium text-gray-900">Decrypt Data</h3>
        <p class="mt-2 text-sm text-gray-500">Decrypt previously encrypted data</p>
      </div>
    </a>
  </div>
</div>