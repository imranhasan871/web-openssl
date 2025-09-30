import type { Handle } from '@sveltejs/kit';
import { config } from '$lib/config';

export const handle: Handle = async ({ event, resolve }) => {
  // Proxy API requests to backend
  if (event.url.pathname.startsWith('/api/')) {
    const backendUrl = `${config.API_URL}${event.url.pathname}${event.url.search}`;

    try {
      const response = await fetch(backendUrl, {
        method: event.request.method,
        headers: event.request.headers,
        body: event.request.method !== 'GET' ? await event.request.arrayBuffer() : undefined,
      });

      return new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
      });
    } catch (error) {
      return new Response(JSON.stringify({ error: 'Backend unavailable' }), {
        status: 503,
        headers: { 'Content-Type': 'application/json' }
      });
    }
  }

  const response = await resolve(event);
  return response;
};