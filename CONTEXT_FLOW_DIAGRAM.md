# Context Flow Diagram

## Visual Representation

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              cmd/api/main.go                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. ROOT CONTEXT CREATED                                                    │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ ctx := context.Background()  ❌ NEVER CANCELS              │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  2. PASSED TO APP INITIALIZATION                                            │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ application, err := app.NewApp(ctx, cfg)                    │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  3. DEFERRED CLEANUP (same context)                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ defer application.Close(ctx)  ❌ Uses non-cancellable ctx   │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  4. SIGNAL RECEIVED (SIGINT/SIGTERM)                                       │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ <- quit                                                       │         │
│     │ application.Close(ctx)  ❌ Still using background context   │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
└──────────────────────────────┼───────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           internal/app/app.go                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  5. APP INITIALIZATION                                                      │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ func NewApp(ctx context.Context, cfg *config.Config)       │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  6. DATABASE INITIALIZATION                                                │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ db, err := initDatabase(ctx, cfg.Database, logger)         │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  7. APP SHUTDOWN                                                            │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ func (a *App) Close(ctx context.Context)                   │         │
│     │   └─> a.Server.Shutdown(ctx)  ❌ Non-cancellable ctx      │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
└──────────────────────────────┼───────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    internal/delivery/http/server.go                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  8. SERVER SHUTDOWN                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ func (s *Server) Shutdown(ctx context.Context) error         │         │
│     │   shutdownCtx, cancel := context.WithTimeout(ctx, 30s)  ✓   │         │
│     │   s.httpServer.Shutdown(shutdownCtx)                         │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
└──────────────────────────────┼───────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HTTP REQUEST CONTEXT FLOW (Separate)                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  9. HTTP REQUEST ARRIVES                                                    │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ r.Context()  ✓ Request-scoped context (cancellable)         │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  10. HANDLER (auth_handler.go)                                              │
│      ┌─────────────────────────────────────────────────────────────┐         │
│      │ h.userUseCase.Register(r.Context(), email, password)  ✓    │         │
│      └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  11. USE CASE (user_usecase.go)                                             │
│      ┌─────────────────────────────────────────────────────────────┐         │
│      │ func (uc *userUseCase) Register(ctx context.Context, ...)   │         │
│      │   uc.userRepo.GetByEmail(ctx, email)  ✓                     │         │
│      │   uc.userRepo.Create(ctx, email, hashedPassword, role)  ✓   │         │
│      │   ctxWithTimeout, cancel := context.WithTimeout(ctx, 5s) ✓ │         │
│      │   uc.tokenGenerator.Generate(ctxWithTimeout, user)  ✓      │         │
│      └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  12. REPOSITORY (user_repository.go)                                       │
│      ┌─────────────────────────────────────────────────────────────┐         │
│      │ func (r *userRepository) Create(ctx context.Context, ...)   │         │
│      │   r.queries.CreateUser(ctx, params)  ✓                      │         │
│      └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
│                              ▼                                               │
│  13. SQLC QUERIES (db/sqlc/users.sql.go)                                    │
│      ┌─────────────────────────────────────────────────────────────┐         │
│      │ func (q *Queries) CreateUser(ctx context.Context, ...)      │         │
│      └─────────────────────────────────────────────────────────────┘         │
│                              │                                               │
└──────────────────────────────┼───────────────────────────────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   DATABASE (Postgres)│
                    └─────────────────────┘
```

## Context Flow Summary

### Initialization Flow (Root Context)

```
main.go (context.Background())
    ↓
app.NewApp(ctx)
    ↓
initDatabase(ctx)
    ↓
database.NewConnection(ctx)
```

### Shutdown Flow (Root Context)

```
main.go (signal received)
    ↓
app.Close(ctx)  [same background context]
    ↓
server.Shutdown(ctx)
    ↓
context.WithTimeout(ctx, 30s)  [creates new timeout context]
    ↓
httpServer.Shutdown(shutdownCtx)
```

### Request Flow (Request Context)

```
HTTP Request
    ↓
r.Context()  [request-scoped context]
    ↓
Handler
    ↓
UseCase
    ↓
Repository
    ↓
SQLC Queries
    ↓
Database
```

---

## IDENTIFIED ISSUES

### 🔴 CRITICAL ISSUE #1: Non-Cancellable Root Context

**Location:** [`cmd/api/main.go:23`](cmd/api/main.go:23)

**Problem:**

```go
ctx := context.Background()
```

**Why it's a problem:**

- `context.Background()` never cancels and has no deadline
- When the application receives a shutdown signal (SIGINT/SIGTERM), the same non-cancellable context is used
- Database operations during initialization cannot be cancelled if they hang
- No way to propagate cancellation signals during the application lifecycle

**Impact:**

- If database connection hangs during startup, the application cannot be gracefully terminated
- Shutdown operations may not respect timeout properly at the root level
- No mechanism to cancel long-running initialization tasks

---

### 🔴 CRITICAL ISSUE #2: No Context for Signal Handling

**Location:** [`cmd/api/main.go:52-58`](cmd/api/main.go:52-58)

**Problem:**

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<- quit

application.Logger.Info("Shutting down application...")
application.Close(ctx)  // Still using the same background context
```

**Why it's a problem:**

- When a shutdown signal is received, a new cancellable context should be created
- The background context is reused, which doesn't provide cancellation semantics
- No timeout at the application level for graceful shutdown

**Impact:**

- Cannot enforce a maximum shutdown time at the application level
- If shutdown hangs, the process may not terminate properly
- No way to cancel shutdown operations if they take too long

---

### 🟡 MEDIUM ISSUE #3: No Timeout for Initialization

**Location:** [`cmd/api/main.go:26`](cmd/api/main.go:26)

**Problem:**

```go
application, err := app.NewApp(ctx, cfg)
```

**Why it's a problem:**

- No timeout for application initialization
- If database connection or other dependencies hang, the application will wait indefinitely
- No way to fail fast if initialization takes too long

**Impact:**

- Application may hang indefinitely during startup
- No feedback to operators about initialization progress
- Difficult to diagnose startup issues

---

### 🟡 MEDIUM ISSUE #4: Context Not Stored in App Struct

**Location:** [`internal/app/app.go:16-23`](internal/app/app.go:16-23)

**Problem:**

```go
type App struct {
    Config      *config.Config
    Logger      *slog.Logger
    DB          *pgxpool.Pool
    UserRepo    domain.UserRepository
    UserUseCase domain.UserUseCase
    Server      *httpDelivery.Server
    // ❌ No context field
}
```

**Why it's a problem:**

- The context is not stored in the App struct
- Other methods that might need the context cannot access it
- Cannot use the context for application-level operations

**Impact:**

- Limited flexibility for context usage across the application
- Cannot implement context-aware operations in App methods
- Harder to add features that require application-level context

---

### 🟢 MINOR ISSUE #5: Double Close Call

**Location:** [`cmd/api/main.go:30`](cmd/api/main.go:30) and [`cmd/api/main.go:58`](cmd/api/main.go:58)

**Problem:**

```go
defer application.Close(ctx)  // Line 30
// ...
application.Close(ctx)  // Line 58 - called again on signal
```

**Why it's a problem:**

- `Close()` is called twice: once via defer and once explicitly on signal
- The defer will execute after the explicit call, potentially causing issues
- The `Close()` method should be idempotent, but this is not guaranteed

**Impact:**

- Potential for double-closing resources
- May cause panic or unexpected behavior if Close() is not idempotent
- Confusing code flow

---

### 🟢 MINOR ISSUE #6: Context Not Used in All Database Operations

**Location:** Various repository methods

**Problem:**
While the context is properly passed through the request flow, it's not clear if all database operations consistently use the context for cancellation.

**Why it's a problem:**

- Some database operations might not respect context cancellation
- Inconsistent context usage across the codebase

**Impact:**

- Some operations may not cancel properly during shutdown
- Potential for resource leaks

---

## RECOMMENDATIONS

### For Critical Issues:

1. **Create a cancellable context for the application lifecycle:**

   ```go
   ctx, cancel := context.WithCancel(context.Background())
   defer cancel()
   ```

2. **Create a timeout context on shutdown signal:**

   ```go
   <- quit
   shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer shutdownCancel()
   application.Close(shutdownCtx)
   ```

3. **Add timeout for initialization:**
   ```go
   initCtx, initCancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer initCancel()
   application, err := app.NewApp(initCtx, cfg)
   ```

### For Medium Issues:

4. **Store context in App struct** (if needed for application-level operations)

5. **Remove the defer Close()** and only call it on signal

### For Minor Issues:

6. **Ensure all database operations use context** consistently

7. **Make Close() idempotent** to handle multiple calls safely

---

## POSITIVE ASPECTS ✓

1. **Request context flow is correct:** HTTP handlers properly use `r.Context()` and pass it down through use cases and repositories
2. **Server shutdown creates timeout context:** The server's Shutdown method properly creates a timeout context
3. **Token generation uses timeout:** The use case creates a timeout context for token generation
4. **Context is passed through all layers:** The request context flows correctly from handler → use case → repository → database
5. **SQLC queries accept context:** All database queries properly accept and use context
