/**
 * HTTP Client Repository
 * Following Repository Pattern
 * Single Responsibility: Handle HTTP communication
 */

import { config } from '$lib/config';
import { authStore } from '$lib/stores/auth';
import { get } from 'svelte/store';
import type { ApiResponse } from '$lib/services/interfaces';

/**
 * IHttpClient interface for HTTP operations
 * Allows easy mocking and testing
 */
export interface IHttpClient {
  get<T>(endpoint: string): Promise<ApiResponse<T>>;
  post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>>;
  put<T>(endpoint: string, data?: any): Promise<ApiResponse<T>>;
  delete<T>(endpoint: string): Promise<ApiResponse<T>>;
  uploadFile<T>(endpoint: string, file: File, additionalData?: Record<string, string>): Promise<ApiResponse<T>>;
}

/**
 * HttpClient implementation
 * Dependency Inversion: Depends on interface, not concrete implementation
 */
export class HttpClient implements IHttpClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = endpoint.startsWith('http') ? endpoint : `${this.baseUrl}${endpoint}`;

    // Get auth token from store
    const auth = get(authStore);

    const headers = new Headers(options.headers);
    headers.set('Content-Type', 'application/json');

    if (auth.token) {
      headers.set('Authorization', `Bearer ${auth.token}`);
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      const data = await response.json();

      if (!response.ok) {
        // Handle auth errors
        if (response.status === 401) {
          authStore.logout();
        }

        return {
          success: false,
          error: data.error || `HTTP ${response.status}`,
          data: data
        };
      }

      return {
        success: true,
        data: data
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Network error'
      };
    }
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }

  async uploadFile<T>(endpoint: string, file: File, additionalData?: Record<string, string>): Promise<ApiResponse<T>> {
    const url = endpoint.startsWith('http') ? endpoint : `${this.baseUrl}${endpoint}`;
    const auth = get(authStore);

    const formData = new FormData();
    formData.append('file', file);

    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value);
      });
    }

    const headers = new Headers();
    if (auth.token) {
      headers.set('Authorization', `Bearer ${auth.token}`);
    }

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers,
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        if (response.status === 401) {
          authStore.logout();
        }

        return {
          success: false,
          error: data.error || `HTTP ${response.status}`,
          data: data
        };
      }

      return {
        success: true,
        data: data
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Network error'
      };
    }
  }
}

// Singleton instance
// Can be replaced with mock for testing
export const httpClient: IHttpClient = new HttpClient(config.API_URL);