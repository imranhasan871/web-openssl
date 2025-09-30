<script lang="ts">
  import { notifications, type Notification } from '$lib/stores/notifications';
  import { fly } from 'svelte/transition';

  function getNotificationIcon(type: Notification['type']) {
    switch (type) {
      case 'success':
        return {
          path: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
          color: 'text-green-400'
        };
      case 'error':
        return {
          path: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
          color: 'text-red-400'
        };
      case 'warning':
        return {
          path: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z',
          color: 'text-yellow-400'
        };
      case 'info':
        return {
          path: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
          color: 'text-blue-400'
        };
    }
  }

  function getBgColor(type: Notification['type']) {
    switch (type) {
      case 'success':
        return 'bg-green-50 border-green-200';
      case 'error':
        return 'bg-red-50 border-red-200';
      case 'warning':
        return 'bg-yellow-50 border-yellow-200';
      case 'info':
        return 'bg-blue-50 border-blue-200';
    }
  }
</script>

<div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm w-full">
  {#each $notifications as notification (notification.id)}
    <div
      transition:fly={{ x: 400, duration: 300 }}
      class="rounded-md border p-4 shadow-lg {getBgColor(notification.type)}"
    >
      <div class="flex">
        <div class="flex-shrink-0">
          <svg
            class="h-5 w-5 {getNotificationIcon(notification.type).color}"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d={getNotificationIcon(notification.type).path}
            />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-sm font-medium text-gray-800">
            {notification.title}
          </h3>
          <div class="mt-1 text-sm text-gray-700">
            {notification.message}
          </div>
        </div>
        {#if notification.dismissible}
          <div class="ml-auto pl-3">
            <button
              type="button"
              class="inline-flex rounded-md text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
              on:click={() => notifications.remove(notification.id)}
            >
              <span class="sr-only">Dismiss</span>
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/each}
</div>