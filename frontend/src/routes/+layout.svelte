<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import Notifications from '$lib/components/ui/Notifications.svelte';

  onMount(() => {
    // Initialize auth store
    authStore.init();
  });
</script>

<main>
  <slot />
</main>

<!-- Global notifications -->
<Notifications />