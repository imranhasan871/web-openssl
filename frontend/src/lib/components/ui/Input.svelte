<script lang="ts">
  export let type = 'text';
  export let value = '';
  export let placeholder = '';
  export let disabled = false;
  export let required = false;
  export let label = '';
  export let error = '';
  export let hint = '';
  export let id = '';

  let inputElement: HTMLInputElement;

  export function focus() {
    inputElement?.focus();
  }
</script>

<div class="space-y-1">
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700">
      {label}
      {#if required}
        <span class="text-red-500">*</span>
      {/if}
    </label>
  {/if}

  <input
    bind:this={inputElement}
    bind:value
    {type}
    {placeholder}
    {disabled}
    {required}
    {id}
    class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed
      {error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}"
    on:input
    on:change
    on:focus
    on:blur
    on:keydown
  />

  {#if error}
    <p class="text-sm text-red-600" role="alert">
      {error}
    </p>
  {:else if hint}
    <p class="text-sm text-gray-500">
      {hint}
    </p>
  {/if}
</div>