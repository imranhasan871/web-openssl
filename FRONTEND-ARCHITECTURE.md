# Frontend Architecture - Clean Architecture with SOLID Principles

## Overview

The frontend follows **Clean Architecture** principles adapted for SvelteKit, with emphasis on **SOLID principles**, **separation of concerns**, and **testability**.

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                  Presentation Layer                      │
│  (Svelte Components - UI)                               │
│  - Routes (+page.svelte)                                │
│  - Components                                            │
│  - Layouts                                               │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                  Composables Layer                       │
│  (Reusable Logic - Hooks)                               │
│  - useAuth, useUser, useOperations                      │
│  - Business logic orchestration                          │
│  - State management                                      │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                   Service Layer                          │
│  (Business Logic)                                        │
│  - AuthService, UserService, OperationService           │
│  - Domain logic                                          │
│  - Validation                                            │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                  Repository Layer                        │
│  (Data Access)                                          │
│  - HttpClient (API communication)                       │
│  - LocalStorage (Cache)                                 │
│  - IndexedDB (Offline storage)                          │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│               Infrastructure Layer                       │
│  (External Services)                                    │
│  - Backend API                                          │
│  - Browser APIs                                         │
└─────────────────────────────────────────────────────────┘
```

## SOLID Principles in Frontend

### 1. Single Responsibility Principle (SRP)

Each module has ONE reason to change:

**❌ Bad: Component doing everything**
```svelte
<script>
  // Component handles: UI, API calls, validation, state management
  async function login() {
    if (!email || !password) return; // validation
    const response = await fetch('/api/login', { // API call
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
    const data = await response.json();
    localStorage.setItem('token', data.token); // state management
  }
</script>
```

**✅ Good: Separated concerns**
```svelte
<script>
  import { useAuth } from '$lib/composables';

  const { login, loading } = useAuth();

  async function handleLogin() {
    await login({ email, password });
  }
</script>

<!-- Component only handles UI -->
<button on:click={handleLogin} disabled={$loading}>
  Sign In
</button>
```

### 2. Open/Closed Principle (OCP)

Open for extension, closed for modification:

```typescript
// Interface - Open for extension
interface IAuthService {
  login(credentials: LoginRequest): Promise<ApiResponse>;
}

// Implementation - Can extend without modifying
class AuthService implements IAuthService {
  async login(credentials: LoginRequest) {
    return this.httpClient.post('/auth/login', credentials);
  }
}

// Extension - Add logging without modifying AuthService
class LoggedAuthService implements IAuthService {
  constructor(private authService: IAuthService) {}

  async login(credentials: LoginRequest) {
    console.log('Login attempt:', credentials.email);
    const result = await this.authService.login(credentials);
    console.log('Login result:', result.success);
    return result;
  }
}
```

### 3. Liskov Substitution Principle (LSP)

Implementations can be swapped:

```typescript
// Can use real HTTP client
const service = new AuthService(new HttpClient(API_URL));

// Can use mock for testing
const service = new AuthService(new MockHttpClient());

// Can use cached version
const service = new AuthService(new CachedHttpClient());
```

### 4. Interface Segregation Principle (ISP)

Small, focused interfaces:

```typescript
// ❌ Bad: Fat interface
interface IApiClient {
  get(url: string): Promise<any>;
  post(url: string, data: any): Promise<any>;
  uploadFile(file: File): Promise<any>;
  downloadFile(url: string): Promise<Blob>;
  streamData(url: string): AsyncIterator<any>;
}

// ✅ Good: Segregated interfaces
interface IHttpClient {
  get<T>(endpoint: string): Promise<ApiResponse<T>>;
  post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>>;
}

interface IFileClient {
  uploadFile<T>(endpoint: string, file: File): Promise<ApiResponse<T>>;
  downloadFile(url: string): Promise<Blob>;
}
```

### 5. Dependency Inversion Principle (DIP)

Depend on abstractions:

```typescript
// ❌ Bad: Depends on concrete implementation
class AuthService {
  private httpClient = new HttpClient(); // Hard dependency
}

// ✅ Good: Depends on abstraction
class AuthService implements IAuthService {
  constructor(private httpClient: IHttpClient) {} // Injected dependency
}
```

## Project Structure

```
frontend/src/
├── lib/
│   ├── api/                    # Legacy (to be removed)
│   ├── components/             # UI Components
│   │   ├── ui/                # Reusable UI components
│   │   │   ├── Button.svelte
│   │   │   ├── Input.svelte
│   │   │   └── Modal.svelte
│   │   └── Navigation.svelte  # Feature components
│   ├── composables/           # ✨ Reusable Logic (NEW)
│   │   ├── useAuth.ts        # Authentication logic
│   │   ├── useUser.ts        # User management logic
│   │   ├── useOperations.ts  # Operations logic
│   │   └── index.ts
│   ├── services/              # ✨ Business Logic (NEW)
│   │   ├── interfaces.ts     # Service contracts
│   │   ├── auth.service.ts   # Auth business logic
│   │   ├── user.service.ts   # User business logic
│   │   ├── operation.service.ts
│   │   └── index.ts
│   ├── repositories/          # ✨ Data Access (NEW)
│   │   └── http-client.ts    # HTTP abstraction
│   ├── stores/                # Global State
│   │   ├── auth.ts           # Auth state
│   │   └── notifications.ts  # Notification state
│   ├── types/                 # TypeScript types
│   └── config.ts              # Configuration
└── routes/                    # Pages
    ├── +layout.svelte
    ├── +page.svelte          # Home page
    ├── login/
    ├── register/
    └── dashboard/
        ├── +layout.svelte
        ├── +page.svelte
        ├── certificates/
        ├── encryption/
        └── operations/
```

## Design Patterns

### 1. Composable Pattern (Hooks)

Reusable logic extraction:

```typescript
// composables/useAuth.ts
export function useAuth() {
  const { login, register, logout } = authService;

  async function handleLogin(credentials: LoginRequest) {
    const result = await login(credentials);
    if (result.success) {
      authStore.login(result.data.user, result.data.accessToken);
      notifications.success('Logged in!');
      goto('/dashboard');
    }
    return result;
  }

  return {
    login: handleLogin,
    register,
    logout,
    store: authStore
  };
}
```

**Usage in component:**
```svelte
<script>
  import { useAuth } from '$lib/composables';

  const { login } = useAuth();

  async function handleSubmit() {
    await login({ email, password });
  }
</script>
```

### 2. Repository Pattern

Abstract data access:

```typescript
interface IHttpClient {
  get<T>(endpoint: string): Promise<ApiResponse<T>>;
  post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>>;
}

class HttpClient implements IHttpClient {
  // Implementation details
}
```

**Benefits:**
- Easy to mock for testing
- Can swap implementations (REST → GraphQL)
- Centralized error handling
- Consistent API interface

### 3. Service Layer Pattern

Business logic encapsulation:

```typescript
class AuthService implements IAuthService {
  constructor(private httpClient: IHttpClient) {}

  async login(credentials: LoginRequest): Promise<ApiResponse<AuthResponse>> {
    // Business logic
    // Validation
    // Error handling
    return this.httpClient.post('/auth/login', credentials);
  }
}
```

### 4. Store Pattern (State Management)

Centralized state:

```typescript
// stores/auth.ts
export const authStore = writable<AuthState>({
  user: null,
  token: null,
  isAuthenticated: false,
  loading: true
});
```

## Component Composition

### Container/Presentational Pattern

**Container (Smart Component):**
```svelte
<!-- routes/dashboard/+page.svelte -->
<script>
  import { useOperations } from '$lib/composables';
  import OperationsList from '$lib/components/OperationsList.svelte';

  const { operations, loading, loadOperations } = useOperations();

  onMount(() => loadOperations());
</script>

<OperationsList {operations} {loading} />
```

**Presentational (Dumb Component):**
```svelte
<!-- components/OperationsList.svelte -->
<script>
  export let operations: any[];
  export let loading: boolean;
</script>

{#if loading}
  <div>Loading...</div>
{:else}
  {#each operations as operation}
    <OperationItem {operation} />
  {/each}
{/if}
```

## Testing Strategy

### Unit Tests (Services)

```typescript
describe('AuthService', () => {
  let authService: AuthService;
  let mockHttpClient: IHttpClient;

  beforeEach(() => {
    mockHttpClient = {
      post: vi.fn().mockResolvedValue({
        success: true,
        data: { user: mockUser, accessToken: 'token' }
      })
    };
    authService = new AuthService(mockHttpClient);
  });

  it('should login successfully', async () => {
    const result = await authService.login({ email: 'test@example.com', password: 'pass' });

    expect(result.success).toBe(true);
    expect(mockHttpClient.post).toHaveBeenCalledWith('/api/v1/auth/login', {
      email: 'test@example.com',
      password: 'pass'
    });
  });
});
```

### Integration Tests (Composables)

```typescript
describe('useAuth', () => {
  it('should handle login flow', async () => {
    const { login } = useAuth();

    const result = await login({ email: 'test@example.com', password: 'pass' });

    expect(result.success).toBe(true);
    // Verify store updated
    // Verify notification shown
    // Verify navigation occurred
  });
});
```

### E2E Tests (Playwright)

```typescript
test('user can login', async ({ page }) => {
  await page.goto('/login');
  await page.fill('[name=email]', 'test@example.com');
  await page.fill('[name=password]', 'password');
  await page.click('button[type=submit]');

  await expect(page).toHaveURL('/dashboard');
  await expect(page.locator('text=Welcome')).toBeVisible();
});
```

## Best Practices

### 1. Component Structure

```svelte
<script lang="ts">
  // 1. Imports
  import { onMount } from 'svelte';
  import { useAuth } from '$lib/composables';

  // 2. Props
  export let title: string;

  // 3. Composables/Hooks
  const { login } = useAuth();

  // 4. Local state
  let email = '';
  let password = '';

  // 5. Derived state
  $: isValid = email && password;

  // 6. Functions
  async function handleSubmit() {
    await login({ email, password });
  }

  // 7. Lifecycle
  onMount(() => {
    // Initialization
  });
</script>

<!-- 8. Template -->
<div>
  <h1>{title}</h1>
  <form on:submit|preventDefault={handleSubmit}>
    <!-- ... -->
  </form>
</div>

<!-- 9. Styles (scoped) -->
<style>
  div {
    /* Component styles */
  }
</style>
```

### 2. Error Handling

```typescript
// Centralized in services
async function login(credentials: LoginRequest) {
  try {
    const response = await httpClient.post('/auth/login', credentials);
    if (response.success) {
      return response;
    }
    throw new Error(response.error || 'Login failed');
  } catch (error) {
    // Log to error tracking service
    console.error('Login error:', error);
    throw error;
  }
}
```

### 3. Type Safety

```typescript
// Use interfaces for contracts
interface LoginRequest {
  email: string;
  password: string;
}

interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

// Use generic types
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}
```

### 4. Reactive State

```svelte
<script>
  import { authStore } from '$lib/stores/auth';

  // Subscribe to store reactively
  $: user = $authStore.user;
  $: isAuthenticated = $authStore.isAuthenticated;
</script>

{#if isAuthenticated}
  <p>Welcome, {user.firstName}!</p>
{/if}
```

## Migration Guide

### From Old to New Architecture

**Before (Old):**
```svelte
<script>
  import { apiClient } from '$lib/api/client';

  async function login() {
    const response = await apiClient.post('/auth/login', { email, password });
    if (response.success) {
      authStore.login(response.data.user, response.data.accessToken);
    }
  }
</script>
```

**After (New):**
```svelte
<script>
  import { useAuth } from '$lib/composables';

  const { login } = useAuth();

  async function handleLogin() {
    await login({ email, password });
  }
</script>
```

## Advantages

### 1. Maintainability
- Clear separation of concerns
- Easy to locate code
- Predictable structure

### 2. Testability
- Services can be unit tested
- Components can be tested in isolation
- Easy to mock dependencies

### 3. Reusability
- Composables shared across components
- Services reusable across features
- UI components composable

### 4. Type Safety
- TypeScript interfaces
- Compile-time error checking
- Better IDE support

### 5. Scalability
- Easy to add new features
- Can split into packages
- Team collaboration friendly

## Conclusion

This architecture provides:
- ✅ Clean separation of concerns
- ✅ SOLID principles applied
- ✅ High testability
- ✅ Easy maintenance
- ✅ Type safety
- ✅ Reusable components
- ✅ Scalable structure