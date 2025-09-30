/**
 * AuthService Implementation
 * Single Responsibility: Handle authentication business logic
 * Open/Closed: Can extend without modifying
 */

import type { IAuthService, LoginRequest, RegisterRequest, AuthResponse, ApiResponse } from './interfaces';
import type { IHttpClient } from '$lib/repositories/http-client';
import { httpClient } from '$lib/repositories/http-client';

export class AuthService implements IAuthService {
  private httpClient: IHttpClient;

  // Dependency Injection via constructor
  constructor(client: IHttpClient = httpClient) {
    this.httpClient = client;
  }

  async login(credentials: LoginRequest): Promise<ApiResponse<AuthResponse>> {
    return this.httpClient.post<AuthResponse>('/api/v1/auth/login', credentials);
  }

  async register(userData: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
    return this.httpClient.post<AuthResponse>('/api/v1/auth/register', userData);
  }

  async logout(): Promise<void> {
    // Could call backend logout endpoint if needed
    // For now, just clear local state (handled by authStore)
  }

  async refreshToken(): Promise<ApiResponse<{ accessToken: string }>> {
    return this.httpClient.post<{ accessToken: string }>('/api/v1/auth/refresh');
  }

  async forgotPassword(email: string): Promise<ApiResponse<void>> {
    return this.httpClient.post<void>('/api/v1/auth/forgot-password', { email });
  }

  async resetPassword(token: string, password: string): Promise<ApiResponse<void>> {
    return this.httpClient.post<void>('/api/v1/auth/reset-password', { token, password });
  }
}

// Export singleton instance
// Can be replaced with mock for testing
export const authService = new AuthService();