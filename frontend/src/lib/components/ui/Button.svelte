<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' | 'outline' = 'primary';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let disabled = false;
  export let loading = false;
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let fullWidth = false;

  $: classes = [
    'inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed',

    // Size variants
    size === 'sm' ? 'px-3 py-1.5 text-sm' : '',
    size === 'md' ? 'px-4 py-2 text-sm' : '',
    size === 'lg' ? 'px-6 py-3 text-base' : '',

    // Color variants
    variant === 'primary' ? 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500' : '',
    variant === 'secondary' ? 'bg-gray-100 text-gray-900 hover:bg-gray-200 focus:ring-gray-500' : '',
    variant === 'danger' ? 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500' : '',
    variant === 'outline' ? 'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-blue-500' : '',

    // Full width
    fullWidth ? 'w-full' : '',
  ].filter(Boolean).join(' ');
</script>

<button
  class={classes}
  {type}
  {disabled}
  on:click
  on:mouseenter
  on:mouseleave
  on:focus
  on:blur
>
  {#if loading}
    <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  {/if}
  <slot />
</button>