import { writable } from 'svelte/store';

export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message: string;
  duration?: number;
  dismissible?: boolean;
}

function createNotificationStore() {
  const { subscribe, update } = writable<Notification[]>([]);

  const store = {
    subscribe,

    // Add notification
    add: (notification: Omit<Notification, 'id'>) => {
      const id = crypto.randomUUID();
      const newNotification: Notification = {
        ...notification,
        id,
        duration: notification.duration ?? 5000,
        dismissible: notification.dismissible ?? true
      };

      update(notifications => [...notifications, newNotification]);

      // Auto-dismiss after duration
      const duration = newNotification.duration;
      if (duration && duration > 0) {
        setTimeout(() => {
          store.remove(id);
        }, duration);
      }

      return id;
    },

    // Remove notification
    remove: (id: string) => {
      update(notifications => notifications.filter(n => n.id !== id));
    },

    // Clear all notifications
    clear: () => {
      update(() => []);
    },

    // Convenience methods
    success: (title: string, message: string, duration?: number) => {
      return store.add({ type: 'success', title, message, duration });
    },

    error: (title: string, message: string, duration?: number) => {
      return store.add({ type: 'error', title, message, duration: duration ?? 8000 });
    },

    warning: (title: string, message: string, duration?: number) => {
      return store.add({ type: 'warning', title, message, duration });
    },

    info: (title: string, message: string, duration?: number) => {
      return store.add({ type: 'info', title, message, duration });
    }
  };

  return store;
}

export const notifications = createNotificationStore();