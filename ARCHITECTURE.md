# OpenSSL UI - Software Architecture

## Overview

This project follows **Clean Architecture** principles with a strong emphasis on **SOLID principles** for maintainability, testability, and changeability.

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│  (HTTP Handlers - Gin Controllers)                      │
│  - Handles HTTP requests/responses                       │
│  - Input validation                                      │
│  - Authentication/Authorization                          │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                   Application Layer                      │
│  (Services - Business Logic)                            │
│  - Core business rules                                   │
│  - Use case implementations                              │
│  - Domain logic orchestration                            │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                   Data Access Layer                      │
│  (Repositories - Database Operations)                   │
│  - CRUD operations                                       │
│  - Query implementations                                 │
│  - Database abstractions                                 │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                  Infrastructure Layer                    │
│  (Database, External Services)                          │
│  - PostgreSQL (GORM)                                    │
│  - Redis                                                │
│  - Stripe API                                           │
└─────────────────────────────────────────────────────────┘
```

## SOLID Principles Applied

### 1. Single Responsibility Principle (SRP)

Each class/struct has ONE reason to change:

- **Handlers**: Only handle HTTP concerns (request/response)
- **Services**: Only contain business logic
- **Repositories**: Only handle data persistence
- **Models**: Only represent data structures

**Example:**
```go
// ❌ BAD: Handler doing too much
func (h *Handler) CreateUser(c *gin.Context) {
    // Parsing request
    // Validating business rules
    // Hashing password
    // Database operations
    // Sending email
}

// ✅ GOOD: Separated concerns
// Handler - HTTP concerns
func (h *UserHandler) Register(c *gin.Context) {
    var req services.RegisterRequest
    c.ShouldBindJSON(&req)
    user, err := h.userService.Register(c.Request.Context(), req)
    c.JSON(http.StatusCreated, user)
}

// Service - Business logic
func (s *userService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
    // Business rules
    // Validation
    // Orchestration
}

// Repository - Data access
func (r *userRepository) Create(ctx context.Context, user *User) error {
    // Database operations only
}
```

### 2. Open/Closed Principle (OCP)

Open for extension, closed for modification. Use interfaces and composition.

```go
// Define interface
type UserService interface {
    Register(ctx context.Context, req RegisterRequest) (*User, error)
    Authenticate(ctx context.Context, email, password string) (*User, error)
}

// Implement interface
type userService struct {
    userRepo repository.UserRepository
}

// Easy to extend with decorators/wrappers
type cachedUserService struct {
    UserService
    cache Cache
}
```

### 3. Liskov Substitution Principle (LSP)

Implementations can be substituted without breaking functionality:

```go
// Any implementation of UserRepository can be used
var repo repository.UserRepository

// Can use real database
repo = repository.NewUserRepository(db)

// Can use mock for testing
repo = &MockUserRepository{}

// Can use cached version
repo = &CachedUserRepository{repo, cache}
```

### 4. Interface Segregation Principle (ISP)

Clients shouldn't depend on interfaces they don't use:

```go
// ❌ BAD: Fat interface
type Repository interface {
    CreateUser(user *User) error
    GetUser(id uint) (*User, error)
    CreateOperation(op *Operation) error
    GetOperation(id uint) (*Operation, error)
    // ... many more methods
}

// ✅ GOOD: Segregated interfaces
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uint) (*User, error)
}

type OperationRepository interface {
    Create(ctx context.Context, operation *Operation) error
    GetByID(ctx context.Context, id uint) (*Operation, error)
}
```

### 5. Dependency Inversion Principle (DIP)

Depend on abstractions, not concretions:

```go
// ❌ BAD: Direct dependency on concrete implementation
type UserService struct {
    db *gorm.DB  // Depends on concrete DB
}

// ✅ GOOD: Depend on interface
type UserService struct {
    userRepo UserRepository  // Depends on abstraction
}
```

## Design Patterns Used

### 1. Repository Pattern

Abstracts data access logic:

```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    // ...
}
```

**Benefits:**
- Centralized data access logic
- Easy to mock for testing
- Database-agnostic
- Can swap implementations (SQL → NoSQL)

### 2. Service Layer Pattern

Encapsulates business logic:

```go
type UserService interface {
    Register(ctx context.Context, req RegisterRequest) (*User, error)
    Authenticate(ctx context.Context, email, password string) (*User, error)
}
```

**Benefits:**
- Clear separation of concerns
- Business logic reusable
- Easy to test
- Transaction management

### 3. Dependency Injection (DI)

Constructor injection for loose coupling:

```go
type Container struct {
    UserRepo      UserRepository
    UserService   UserService
    OperationService OperationService
}

func NewContainer(cfg *config.Config, db *gorm.DB) *Container {
    container := &Container{}

    // Wire dependencies
    container.UserRepo = repository.NewUserRepository(db)
    container.UserService = services.NewUserService(
        container.UserRepo,
        authService,
    )

    return container
}
```

**Benefits:**
- Easy to test (inject mocks)
- Loose coupling
- Single point of configuration
- Clear dependency graph

### 4. Factory Pattern

For creating complex objects:

```go
func NewUserService(
    userRepo repository.UserRepository,
    authService *auth.Service,
) UserService {
    return &userService{
        userRepo:    userRepo,
        authService: authService,
    }
}
```

## Project Structure

```
backend/
├── cmd/
│   ├── server/         # Application entry point
│   └── seed/           # Database seeding
├── internal/
│   ├── config/         # Configuration management
│   ├── container/      # Dependency injection container
│   ├── handlers/       # HTTP handlers (controllers)
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Domain models
│   ├── repository/     # Data access layer
│   │   ├── interfaces.go
│   │   ├── user_repository.go
│   │   └── operation_repository.go
│   └── services/       # Business logic layer
│       ├── interfaces.go
│       ├── user_service.go
│       └── operation_service.go
├── pkg/                # Shared packages
│   ├── auth/           # Authentication utilities
│   └── billing/        # Billing integration
└── tests/              # Test files
```

## Testing Strategy

### Unit Tests

Test individual components in isolation:

```go
func TestUserService_Register(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := services.NewUserService(mockRepo, mockAuthService)

    // Act
    user, err := service.Register(ctx, req)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### Integration Tests

Test component interactions:

```go
func TestUserAPI_Integration(t *testing.T) {
    // Use real database (test DB)
    container := setupTestContainer()
    defer container.Close()

    // Test full flow
    resp := makeRequest("/api/v1/auth/register", registerData)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## Advantages of This Architecture

### 1. Maintainability
- Clear separation of concerns
- Each layer has specific responsibilities
- Easy to locate and fix bugs

### 2. Testability
- Dependencies injected via interfaces
- Easy to mock
- Test in isolation
- Fast unit tests

### 3. Changeability
- Can swap implementations (e.g., PostgreSQL → MongoDB)
- Change UI framework without touching business logic
- Add new features without modifying existing code

### 4. Scalability
- Stateless services
- Horizontal scaling friendly
- Can split into microservices later

### 5. Team Collaboration
- Clear boundaries between layers
- Multiple developers can work simultaneously
- Less merge conflicts

## Best Practices

### 1. Always Use Context
```go
func (s *userService) Register(ctx context.Context, req RegisterRequest) (*User, error)
```
- Enables request tracing
- Timeout/cancellation support
- Passes request-scoped data

### 2. Error Handling
```go
// Return business errors, not implementation details
if err != nil {
    return nil, errors.New("failed to register user")
}
```

### 3. DTOs (Data Transfer Objects)
```go
// Use DTOs for API boundaries
type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
}
```

### 4. Immutability
```go
// Return copies, not references to internal state
func (s *userService) GetProfile(ctx context.Context, userID uint) (*User, error) {
    user, _ := s.userRepo.GetByID(ctx, userID)
    userCopy := *user
    userCopy.Password = "" // Remove sensitive data
    return &userCopy, nil
}
```

## Migration Path

To migrate existing code to this architecture:

1. **Create interfaces first** (repository, service)
2. **Implement repositories** (extract DB code from handlers)
3. **Implement services** (extract business logic from handlers)
4. **Refactor handlers** (thin layer, delegate to services)
5. **Setup DI container** (wire everything together)
6. **Add tests** (now easy with mocks)

## Conclusion

This architecture provides:
- ✅ Clean separation of concerns
- ✅ High testability
- ✅ Easy to maintain and extend
- ✅ Follows industry best practices
- ✅ Ready for scaling