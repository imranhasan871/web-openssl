import { apiClient, type ApiResponse } from './client';

// Certificate Interfaces
export interface GenerateCertificateRequest {
  commonName: string;
  organization?: string;
  organizationalUnit?: string;
  country?: string;
  state?: string;
  locality?: string;
  emailAddress?: string;
  subjectAltNames?: string[];
  keySize?: number;
  validityDays?: number;
  keyUsage?: string[];
  extendedKeyUsage?: string[];
}

export interface GenerateKeyRequest {
  type: 'rsa' | 'ecdsa' | 'ed25519';
  size?: number;
  curve?: string;
  format?: 'pem' | 'der';
}

export interface GenerateCSRRequest {
  commonName: string;
  organization?: string;
  organizationalUnit?: string;
  country?: string;
  state?: string;
  locality?: string;
  emailAddress?: string;
  subjectAltNames?: string[];
  privateKey?: string;
  keySize?: number;
}

export interface ParseCertificateRequest {
  certificate: string;
  format?: 'pem' | 'der';
}

export interface VerifyCertificateRequest {
  certificate: string;
  caCertificate?: string;
  checkRevocation?: boolean;
}

export interface ConvertCertificateRequest {
  input: string;
  inputFormat: 'pem' | 'der' | 'p7b' | 'p12';
  outputFormat: 'pem' | 'der' | 'p7b' | 'p12';
  password?: string;
}

// Encryption Interfaces
export interface SymmetricEncryptRequest {
  data: string;
  algorithm: 'aes-256-cbc' | 'aes-256-gcm' | 'aes-192-cbc' | 'aes-128-cbc';
  key?: string;
  iv?: string;
}

export interface AsymmetricEncryptRequest {
  data: string;
  publicKey: string;
  algorithm?: 'rsa-oaep' | 'rsa-pkcs1';
}

export interface DecryptRequest {
  encryptedData: string;
  key: string;
  algorithm: string;
  iv?: string;
}

export interface GenerateHashRequest {
  data: string;
  algorithm: 'md5' | 'sha1' | 'sha256' | 'sha384' | 'sha512';
}

export interface GenerateHMACRequest {
  data: string;
  key: string;
  algorithm: 'sha256' | 'sha384' | 'sha512';
}

export interface VerifyHashRequest {
  data: string;
  hash: string;
  algorithm: string;
}

// SSL Analysis Interfaces
export interface TestSSLConnectionRequest {
  hostname: string;
  port?: number;
  protocol?: 'tls1.2' | 'tls1.3';
}

export interface AnalyzeSSLCertificateRequest {
  hostname: string;
  port?: number;
}

export const opensslAPI = {
  // Certificate Operations
  generateCertificate: async (data: GenerateCertificateRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/certificates/generate', data);
  },

  generateKey: async (data: GenerateKeyRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/keys/generate', data);
  },

  generateCSR: async (data: GenerateCSRRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/certificates/csr', data);
  },

  parseCertificate: async (data: ParseCertificateRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/certificates/parse', data);
  },

  verifyCertificate: async (data: VerifyCertificateRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/certificates/verify', data);
  },

  convertCertificate: async (data: ConvertCertificateRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/certificates/convert', data);
  },

  parseKey: async (data: { key: string; format?: string }): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/keys/parse', data);
  },

  convertKey: async (data: ConvertCertificateRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/keys/convert', data);
  },

  // Encryption Operations
  symmetricEncrypt: async (data: SymmetricEncryptRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/encrypt/symmetric', data);
  },

  asymmetricEncrypt: async (data: AsymmetricEncryptRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/encrypt/asymmetric', data);
  },

  decrypt: async (data: DecryptRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/encrypt/decrypt', data);
  },

  generateHash: async (data: GenerateHashRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/hash/generate', data);
  },

  generateHMAC: async (data: GenerateHMACRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/hash/hmac', data);
  },

  verifyHash: async (data: VerifyHashRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/hash/verify', data);
  },

  // SSL Analysis
  testSSLConnection: async (data: TestSSLConnectionRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/ssl/test-connection', data);
  },

  analyzeSSLCertificate: async (data: AnalyzeSSLCertificateRequest): Promise<ApiResponse> => {
    return apiClient.post('/api/v1/openssl/ssl/analyze-certificate', data);
  }
};