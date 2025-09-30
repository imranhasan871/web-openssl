/**
 * UserService Implementation
 * Single Responsibility: Handle user profile operations
 */

import type { IUserService, UpdateProfileRequest, ApiResponse } from './interfaces';
import type { User } from '$lib/stores/auth';
import type { IHttpClient } from '$lib/repositories/http-client';
import { httpClient } from '$lib/repositories/http-client';

export class UserService implements IUserService {
  private httpClient: IHttpClient;

  constructor(client: IHttpClient = httpClient) {
    this.httpClient = client;
  }

  async getProfile(): Promise<ApiResponse<User>> {
    return this.httpClient.get<User>('/api/v1/users/me');
  }

  async updateProfile(data: UpdateProfileRequest): Promise<ApiResponse<User>> {
    return this.httpClient.put<User>('/api/v1/users/me', data);
  }

  async deleteAccount(): Promise<ApiResponse<void>> {
    return this.httpClient.delete<void>('/api/v1/users/me');
  }

  async generateAPIKey(): Promise<ApiResponse<{ apiKey: string }>> {
    return this.httpClient.post<{ apiKey: string }>('/api/v1/users/api-key');
  }
}

export const userService = new UserService();