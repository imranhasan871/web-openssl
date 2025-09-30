<script lang="ts">
  import { onMount } from 'svelte';
  import { config } from '$lib/config';
  import { goto } from '$app/navigation';
  import { authStore } from '$lib/stores/auth';

  let apiStatus = 'Checking...';
  let isAuthenticated = false;

  onMount(async () => {
    // Check auth status
    authStore.subscribe((auth) => {
      isAuthenticated = auth.isAuthenticated;
    });

    // Check backend health
    try {
      const response = await fetch(`${config.API_URL}/health`);
      if (response.ok) {
        apiStatus = 'Connected';
      } else {
        apiStatus = 'Error';
      }
    } catch (error) {
      apiStatus = 'Offline';
    }
  });

  const features = [
    {
      title: 'Certificate Generation',
      description: 'Generate self-signed SSL/TLS certificates, private keys, and certificate signing requests with customizable parameters.',
      icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z',
      gradient: 'from-blue-500 to-blue-600'
    },
    {
      title: 'Certificate Management',
      description: 'Parse, verify, and analyze existing certificates. View detailed information, expiration dates, and validation status.',
      icon: 'M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4',
      gradient: 'from-green-500 to-green-600'
    },
    {
      title: 'Encryption & Hashing',
      description: 'Symmetric (AES, DES) and asymmetric (RSA) encryption, digital signatures, and hash generation (SHA-256, MD5, HMAC).',
      icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z',
      gradient: 'from-purple-500 to-purple-600'
    },
    {
      title: 'SSL/TLS Testing',
      description: 'Test SSL connections, analyze certificate chains, verify SSL/TLS configurations, and troubleshoot security issues.',
      icon: 'M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9',
      gradient: 'from-yellow-500 to-orange-600'
    },
    {
      title: 'Format Conversion',
      description: 'Convert between different certificate and key formats: PEM, DER, PKCS12, JKS, and more with ease.',
      icon: 'M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4',
      gradient: 'from-red-500 to-pink-600'
    },
    {
      title: 'Secure & Private',
      description: 'User authentication, role-based access control, operation history, and enterprise-grade security features.',
      icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
      gradient: 'from-indigo-500 to-indigo-600'
    }
  ];
</script>

<svelte:head>
  <title>OpenSSL UI - Professional Certificate Management Platform</title>
  <meta name="description" content="Modern web interface for OpenSSL operations. Generate, manage, and analyze SSL/TLS certificates with professional-grade tools." />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50">
  <!-- Navigation -->
  <nav class="border-b border-gray-200 bg-white/80 backdrop-blur-sm sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
            </svg>
          </div>
          <span class="text-xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">OpenSSL UI</span>
        </div>
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2 text-sm">
            <div class="w-2 h-2 rounded-full {apiStatus === 'Connected' ? 'bg-green-500' : 'bg-red-500'} animate-pulse"></div>
            <span class="text-gray-600">{apiStatus}</span>
          </div>
          {#if isAuthenticated}
            <a href="/dashboard" class="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-blue-600 to-indigo-600 rounded-lg hover:from-blue-700 hover:to-indigo-700 transition-all">
              Dashboard
            </a>
          {:else}
            <a href="/login" class="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">
              Sign In
            </a>
            <a href="/register" class="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-blue-600 to-indigo-600 rounded-lg hover:from-blue-700 hover:to-indigo-700 transition-all">
              Get Started
            </a>
          {/if}
        </div>
      </div>
    </div>
  </nav>

  <!-- Hero Section -->
  <section class="relative overflow-hidden">
    <div class="absolute inset-0 bg-gradient-to-br from-blue-600/5 via-purple-600/5 to-indigo-600/5"></div>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 relative">
      <div class="text-center max-w-4xl mx-auto">
        <div class="inline-flex items-center space-x-2 bg-blue-100 text-blue-700 px-4 py-2 rounded-full text-sm font-medium mb-8">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
          </svg>
          <span>Professional SSL/TLS Certificate Management</span>
        </div>

        <h1 class="text-5xl md:text-6xl lg:text-7xl font-bold text-gray-900 mb-6 leading-tight">
          Simplify Your
          <span class="bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent"> Certificate Management</span>
        </h1>

        <p class="text-xl md:text-2xl text-gray-600 mb-12 leading-relaxed">
          Modern web interface for OpenSSL operations. Generate, manage, and analyze SSL/TLS certificates with enterprise-grade tools and intuitive workflows.
        </p>

        <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
          <a href="/register" class="w-full sm:w-auto px-8 py-4 text-lg font-medium text-white bg-gradient-to-r from-blue-600 to-indigo-600 rounded-xl hover:from-blue-700 hover:to-indigo-700 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all">
            Get Started Free
          </a>
          <a href="/login" class="w-full sm:w-auto px-8 py-4 text-lg font-medium text-gray-700 bg-white border-2 border-gray-300 rounded-xl hover:border-gray-400 hover:bg-gray-50 shadow-md hover:shadow-lg transition-all">
            Sign In
          </a>
        </div>

        <p class="mt-6 text-sm text-gray-500">
          ✨ No credit card required • ⚡ Start in seconds • 🔒 Secure & private
        </p>
      </div>
    </div>
  </section>

  <!-- Features Grid -->
  <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
    <div class="text-center mb-16">
      <h2 class="text-4xl font-bold text-gray-900 mb-4">Everything You Need</h2>
      <p class="text-xl text-gray-600">Powerful tools for modern certificate management</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      {#each features as feature}
        <div class="group relative bg-white rounded-2xl p-8 shadow-md hover:shadow-2xl transition-all duration-300 border border-gray-100 hover:border-transparent hover:-translate-y-1">
          <div class="absolute inset-0 bg-gradient-to-br {feature.gradient} opacity-0 group-hover:opacity-5 rounded-2xl transition-opacity"></div>

          <div class="relative">
            <div class="w-14 h-14 bg-gradient-to-br {feature.gradient} rounded-xl flex items-center justify-center mb-6 shadow-lg">
              <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={feature.icon}/>
              </svg>
            </div>

            <h3 class="text-xl font-bold text-gray-900 mb-3">{feature.title}</h3>
            <p class="text-gray-600 leading-relaxed">{feature.description}</p>
          </div>
        </div>
      {/each}
    </div>
  </section>

  <!-- CTA Section -->
  <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
    <div class="relative bg-gradient-to-br from-blue-600 to-indigo-600 rounded-3xl shadow-2xl overflow-hidden">
      <div class="absolute inset-0 bg-grid-white/10"></div>
      <div class="relative px-8 py-16 sm:px-16 text-center">
        <h2 class="text-4xl md:text-5xl font-bold text-white mb-6">
          Ready to Get Started?
        </h2>
        <p class="text-xl text-blue-100 mb-10 max-w-2xl mx-auto">
          Join thousands of developers and security professionals managing their SSL certificates with confidence.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <a href="/register" class="px-8 py-4 text-lg font-medium text-blue-600 bg-white rounded-xl hover:bg-gray-50 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all">
            Create Free Account
          </a>
          <a href="/login" class="px-8 py-4 text-lg font-medium text-white border-2 border-white rounded-xl hover:bg-white/10 transition-all">
            Sign In
          </a>
        </div>
      </div>
    </div>
  </section>

  <!-- Footer -->
  <footer class="border-t border-gray-200 bg-white/50 backdrop-blur-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center space-x-3 mb-4 md:mb-0">
          <div class="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
            </svg>
          </div>
          <span class="text-lg font-bold text-gray-900">OpenSSL UI</span>
        </div>
        <div class="text-sm text-gray-600">
          © 2025 OpenSSL UI. Professional certificate management platform.
        </div>
      </div>
    </div>
  </footer>
</div>