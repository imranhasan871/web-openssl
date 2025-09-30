/**
 * OperationService Implementation
 * Single Responsibility: Handle operation tracking
 */

import type { IOperationService, ApiResponse } from './interfaces';
import type { IHttpClient } from '$lib/repositories/http-client';
import { httpClient } from '$lib/repositories/http-client';

export class OperationService implements IOperationService {
  private httpClient: IHttpClient;

  constructor(client: IHttpClient = httpClient) {
    this.httpClient = client;
  }

  async getOperations(limit: number = 50): Promise<ApiResponse<any[]>> {
    const response = await this.httpClient.get<{ operations: any[] }>(`/api/v1/operations/?limit=${limit}`);

    if (response.success && response.data) {
      return {
        success: true,
        data: response.data.operations || []
      };
    }

    return {
      success: false,
      error: response.error,
      data: []
    };
  }

  async getStats(): Promise<ApiResponse<any>> {
    return this.httpClient.get<any>('/api/v1/operations/stats');
  }

  async deleteOperation(id: number): Promise<ApiResponse<void>> {
    return this.httpClient.delete<void>(`/api/v1/operations/${id}`);
  }
}

export const operationService = new OperationService();