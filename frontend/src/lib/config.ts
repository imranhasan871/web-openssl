import { dev } from '$app/environment';

// Configuration
export const config = {
  API_URL: dev ? 'http://localhost:8080' : 'http://backend:8080',
  APP_NAME: 'OpenSSL UI',
  APP_VERSION: '1.0.0'
};