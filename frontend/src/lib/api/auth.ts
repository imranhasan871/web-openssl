import { apiClient, type ApiResponse } from './client';
import type { User } from '$lib/stores/auth';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface RefreshTokenResponse {
  token: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  password: string;
}

export interface VerifyEmailRequest {
  token: string;
}

export const authAPI = {
  // Login user
  login: async (credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
    return apiClient.post('/api/v1/auth/login', credentials);
  },

  // Register new user
  register: async (userData: RegisterRequest): Promise<ApiResponse<LoginResponse>> => {
    return apiClient.post('/api/v1/auth/register', userData);
  },

  // Refresh token
  refreshToken: async (): Promise<ApiResponse<RefreshTokenResponse>> => {
    return apiClient.post('/api/v1/auth/refresh');
  },

  // Forgot password
  forgotPassword: async (data: ForgotPasswordRequest): Promise<ApiResponse<void>> => {
    return apiClient.post('/api/v1/auth/forgot-password', data);
  },

  // Reset password
  resetPassword: async (data: ResetPasswordRequest): Promise<ApiResponse<void>> => {
    return apiClient.post('/api/v1/auth/reset-password', data);
  },

  // Verify email
  verifyEmail: async (data: VerifyEmailRequest): Promise<ApiResponse<void>> => {
    return apiClient.post('/api/v1/auth/verify-email', data);
  }
};