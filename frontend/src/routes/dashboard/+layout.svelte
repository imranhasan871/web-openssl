<script lang="ts">
  import Navigation from '$lib/components/Navigation.svelte';
  import { authStore } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  // Protect this route
  onMount(() => {
    const unsubscribe = authStore.subscribe((auth) => {
      if (!auth.loading && !auth.isAuthenticated) {
        goto('/login');
      }
    });

    return unsubscribe;
  });
</script>

<div class="min-h-screen bg-gray-50">
  <Navigation />

  <main>
    <slot />
  </main>
</div>