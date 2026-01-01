# Authentication Flow Visualization

This document visualizes the current authentication architecture and flows in the `devnorth-back` application.

## Architecture Component Diagram

The application follows a clean architecture with dependency injection.

```mermaid
graph TD
    Client[Client] --> Router[HTTP Router (Chi)]
    Router --> AuthHandler[Auth Handler]
    
    subgraph "Domain Layer"
        UserUseCase[User UseCase]
        UserRepoInterface[User Repository Interface]
    end
    
    subgraph "Infrastructure Layer"
        UserRepoImpl[User Repository Implementation]
        HashService[Bcrypt Hasher]
        TokenService[JWT Generator]
        DB[(PostgreSQL)]
    end

    AuthHandler --> UserUseCase
    UserUseCase --> UserRepoInterface
    UserUseCase --> HashService
    UserUseCase --> TokenService
    
    UserRepoInterface -.-> UserRepoImpl
    UserRepoImpl --> DB
```

## Registration Flow (`POST /api/v1/auth/register`)

```mermaid
sequenceDiagram
    participant Client
    participant Handler as AuthHandler
    participant UseCase as UserUseCase
    participant Repo as UserRepository
    participant Hasher as PasswordHasher
    participant Token as TokenGenerator
    participant DB as Database

    Client->>Handler: POST /register {email, password}
    Handler->>Handler: Validate Request Body
    
    Handler->>UseCase: Register(email, password)
    
    UseCase->>UseCase: Validate Email & Password format
    
    UseCase->>Repo: GetByEmail(email)
    Repo->>DB: Select user query
    DB-->>Repo: Result
    Repo-->>UseCase: User or nil
    
    alt Email already exists
        UseCase-->>Handler: Error (EmailExists)
        Handler-->>Client: 409 Conflict
    end
    
    UseCase->>Hasher: Hash(password)
    Hasher-->>UseCase: hashedPassword
    
    UseCase->>Repo: Create(email, hashedPassword)
    Repo->>DB: Insert user query
    DB-->>Repo: Created User
    Repo-->>UseCase: User model
    UseCase-->>Handler: Created User

    note over Handler: Auto-login after registration
    Handler->>UseCase: Login(email, password)
    
    UseCase->>Repo: GetByEmail(email)
    Repo-->>UseCase: User
    UseCase->>Token: Generate(user)
    Token-->>UseCase: JWT Token
    
    UseCase-->>Handler: Token
    Handler-->>Client: 201 Created {user, token}
```

## Login Flow (`POST /api/v1/auth/login`)

```mermaid
sequenceDiagram
    participant Client
    participant Handler as AuthHandler
    participant UseCase as UserUseCase
    participant Repo as UserRepository
    participant Hasher as PasswordHasher
    participant Token as TokenGenerator
    participant DB as Database

    Client->>Handler: POST /login {email, password}
    Handler->>Handler: Validate Request Body
    
    Handler->>UseCase: Login(email, password)
    
    UseCase->>Repo: GetByEmail(email)
    Repo->>DB: Select user query
    DB-->>Repo: User or nil
    
    alt User not found
        UseCase-->>Handler: Error (InvalidCredentials)
        Handler-->>Client: 401 Unauthorized
    end
    
    UseCase->>Hasher: Compare(hashedPassword, inputPassword)
    
    alt Password mismatch
        Hasher-->>UseCase: Error
        UseCase-->>Handler: Error (InvalidCredentials)
        Handler-->>Client: 401 Unauthorized
    end
    
    UseCase->>Token: Generate(user)
    Token-->>UseCase: JWT Token
    
    UseCase-->>Handler: Token + User
    Handler-->>Client: 200 OK {user, token}
```

## Current State Notes

- **Token Generation**: JWT tokens are generated upon successful Login and Registration.
- **Password Security**: Passwords are hashed using Bcrypt (cost 10) before storage.
- **Middleware**: Currently, there is no HTTP middleware implemented to verify the JWT token on protected routes.
