/**
 * Service Layer Interfaces
 * Following Interface Segregation Principle (ISP)
 * Each service has a specific responsibility
 */

import type { User } from '$lib/stores/auth';

// ============= DTOs (Data Transfer Objects) =============

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface UpdateProfileRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
}

export interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// ============= Service Interfaces =============

/**
 * IAuthService - Authentication operations
 * Single Responsibility: Handle user authentication
 */
export interface IAuthService {
  login(credentials: LoginRequest): Promise<ApiResponse<AuthResponse>>;
  register(userData: RegisterRequest): Promise<ApiResponse<AuthResponse>>;
  logout(): Promise<void>;
  refreshToken(): Promise<ApiResponse<{ accessToken: string }>>;
  forgotPassword(email: string): Promise<ApiResponse<void>>;
  resetPassword(token: string, password: string): Promise<ApiResponse<void>>;
}

/**
 * IUserService - User profile operations
 * Single Responsibility: Handle user profile management
 */
export interface IUserService {
  getProfile(): Promise<ApiResponse<User>>;
  updateProfile(data: UpdateProfileRequest): Promise<ApiResponse<User>>;
  deleteAccount(): Promise<ApiResponse<void>>;
  generateAPIKey(): Promise<ApiResponse<{ apiKey: string }>>;
}

/**
 * IOperationService - Operation tracking
 * Single Responsibility: Handle operation history
 */
export interface IOperationService {
  getOperations(limit?: number): Promise<ApiResponse<any[]>>;
  getStats(): Promise<ApiResponse<any>>;
  deleteOperation(id: number): Promise<ApiResponse<void>>;
}

/**
 * ICertificateService - Certificate operations
 * Single Responsibility: Handle certificate generation/management
 */
export interface ICertificateService {
  generate(params: any): Promise<ApiResponse<any>>;
  generateCSR(params: any): Promise<ApiResponse<any>>;
  parse(certPEM: string): Promise<ApiResponse<any>>;
  verify(certPEM: string, caPEM: string): Promise<ApiResponse<any>>;
  convert(input: string, fromFormat: string, toFormat: string): Promise<ApiResponse<any>>;
}

/**
 * IEncryptionService - Encryption operations
 * Single Responsibility: Handle encryption/decryption
 */
export interface IEncryptionService {
  encryptSymmetric(data: string, password: string, algorithm: string): Promise<ApiResponse<string>>;
  encryptAsymmetric(data: string, publicKey: string): Promise<ApiResponse<string>>;
  decrypt(encryptedData: string, key: string): Promise<ApiResponse<string>>;
  generateHash(data: string, algorithm: string): Promise<ApiResponse<string>>;
  generateHMAC(data: string, key: string, algorithm: string): Promise<ApiResponse<string>>;
}