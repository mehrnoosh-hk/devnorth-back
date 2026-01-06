# Architecture Decision Documents (ADD)

This document tracks key architectural decisions made during the development of DevNorth Backend.

## Format

Each decision should include:
- **Date**: When the decision was made
- **Status**: [Proposed|Accepted|Deprecated|Superseded]
- **Context**: What is the issue we're trying to solve?
- **Decision**: What did we decide to do?
- **Consequences**: What are the positive and negative outcomes?
- **POC → Production Steps**: What needs to be done before production

---

## Decisions

### 1. Clean Architecture with Go

**Date**: 2025-12-27

**Context**: Need a scalable, maintainable architecture for the backend while keeping POC speed.

**Decision**: Implement Clean Architecture with clear layer separation (Domain → Repository → UseCase → Delivery).

**Consequences**:
- **Positive**: Clear separation of concerns, testable code, easier to maintain
- **Negative**: Slightly more boilerplate than flat structure
- **Trade-off**: Balanced by POC mindset - avoid over-engineering
**POC → Production Steps**:
- Add comprehensive test coverage (unit + integration tests)
- Implement proper error handling and logging
- Add monitoring and observability
- Review and refactor any shortcuts taken during POC

---
### 2. SQLC for Database Queries

**Date**: 2025-12-27

**Context**: Need type-safe database access without heavy ORM overhead.

**Decision**: Use SQLC to generate type-safe Go code from SQL queries.

**Consequences**:
- **Positive**: Type safety, write plain SQL, no runtime overhead
- **Negative**: Requires code generation step, learning curve
- **Trade-off**: Better than ORM for this project's needs

**POC → Production Steps**:
- Add query performance monitoring
- Implement connection pooling optimization
- Add prepared statement caching
- Review all queries for N+1 issues and optimization opportunities

---
### 3. Simple Go-Idiomatic Dependency Injection

**Date**: 2025-12-27
**Status**: Accepted

**Context**: For POC, there is no need for a complex DI framework. The builder pattern and external DI tools would create unnecessary complexity for a proof-of-concept.

**Decision**: Use simple construction functions to wire dependencies manually. Only define interfaces where abstraction is needed (e.g., Repository layer for database operations). Dependencies are passed explicitly through function parameters.

**Consequences**:
- **Positive**: Minimal boilerplate, explicit dependency flow, no external dependencies, easier to understand, faster development
- **Negative**: Manual wiring of dependencies, less flexibility for swapping implementations, could become unwieldy as the app grows
- **Trade-off**: For POC, simplicity and development speed outweigh the flexibility of a full DI framework. Go's explicit style makes manual wiring acceptable for small-to-medium projects.

**POC → Production Steps**:
- Evaluate if Google Wire or similar DI tool is needed as dependency graph grows
- Evaluate Builder pattern as an option for DI
- Implement interfaces for all major dependencies (UseCase, external services), not just Repository
- Add factory functions for creating test doubles and mocks
- Consider lightweight DI container if manual wiring becomes difficult to maintain
- Document dependency graph and initialization order

---

### 4. User Authentication Schema with CITEXT and Role-Based Access
**Date**: 2025-12-28
**Status**: Accepted

**Context**: Need a user authentication system for the POC. Users need email/password authentication with role-based access control (USER vs ADMIN). Email lookup must be case-insensitive to prevent duplicate accounts like "user@email.com" and "User@Email.com".

**Decision**: Implement users table with:
- PostgreSQL CITEXT extension for case-insensitive email storage
- ENUM type for user roles (USER, ADMIN) with USER as default
- Standard fields: id (SERIAL), email (CITEXT UNIQUE), hashed_password (TEXT), role, created_at, updated_at
- Two SQLC queries for POC: CreateUser and GetUserByEmail

**Consequences**:
- **Positive**: Case-insensitive email prevents duplicate accounts, ENUM enforces valid roles, CITEXT is cleaner than LOWER() functions everywhere, SERIAL id is simple for POC
- **Negative**: CITEXT requires PostgreSQL extension (no MySQL/SQLite portability), SERIAL id may not be ideal for distributed systems, storing role in users table couples auth with authorization
- **Trade-off**: CITEXT is PostgreSQL-specific but provides better developer experience. SERIAL is sufficient for POC (can migrate to UUID later). Simple role field is adequate for two-role system.

**POC → Production Steps**:
- Evaluate switching from SERIAL to UUID for distributed scalability
- Add email validation and sanitization at application layer
- Implement password complexity requirements and hashing verification
- Add user status field (active, suspended, deleted) for account management
- Consider separate roles/permissions table if RBAC becomes more complex
- Add indexes for common query patterns (created_at for pagination, role for filtering)
- Implement rate limiting for GetUserByEmail to prevent enumeration attacks
- Add audit logging for sensitive operations (password changes, role updates)

---

### 5. Implement minimum required config
**Date**: 2025-12-28
**Status**: Accepted

**Context**: Need Configuration for database, server and logger setup

**Decision**: The config should address important db, server, logger and env setup, No hardcoded config is acceptable.

**Consequences**:
- **Positive**: Comprehensive config for db and server
- **Negative**: More time consuming than hardcoded default config
- **Trade-off**: Hardcoding the config would require more clean up and refactor later, also for testing the performance for POC, we should test it with different config.

**POC → Production Steps**:
- Add Validation for config, e.g. AUTH_SECRET should never be empty
- Add config for logger, e.g. log level, log format
- Add config for health check, e.g. health check interval
- More defensive Load function, which fails in case of invalid config or not provided config

### 6. Database Configuration with Connection Pooling and Timeouts
**Date**: 2025-12-28
**Status**: Accepted

**Context**: Need production-ready database configuration for PostgreSQL while maintaining POC simplicity. Database connections are expensive resources that need proper management to prevent connection exhaustion, slow queries, and resource leaks.

**Decision**: Implement comprehensive database configuration with:
- **Connection Pool Settings**: MaxOpenConns (25), MaxIdleConns (5), ConnMaxLifetime (5 min), ConnMaxIdleTime (5 min)
- **Timeout Settings**: ConnectTimeout (10s), QueryTimeout (30s)
- **Migration Settings**: MigrationsPath (db/migrations), MigrationsTable (schema_migrations)
- Separate config structs for Server, Database, and App environment
- Environment variable-based configuration with sensible defaults

**Consequences**:
- **Positive**: Production-ready connection pooling prevents resource exhaustion, timeout settings prevent hanging connections and runaway queries, migration config enables programmatic migration control, sensible defaults make local development easy
- **Negative**: Slightly more configuration complexity than minimal setup, requires understanding of connection pool tuning, incorrect pool settings could impact performance
- **Trade-off**: The added complexity is justified by production readiness. Default values are conservative and safe for POC. Can be tuned later based on actual load patterns.

**POC → Production Steps**:
- Add configuration options for:
  - LogLevel - Database query logging level (e.g., none, error, warn, info)
  - AutoMigrate - Whether to run migrations on startup (risky for production)
  - RequireSSL - Force SSL/TLS connections (production requirement)
  - CACertPath - Path to CA certificate for SSL verification
  - HealthCheckPeriod - Make health check interval configurable (currently hard-coded to 1 minute)
- Monitor connection pool metrics (active connections, wait duration, idle connections)
- Tune pool settings based on actual application load and database capacity
- Add connection retry logic with exponential backoff for initial connection attempts
- Implement query logging and slow query detection using QueryTimeout
- Add database health checks and connection validation
- Consider adding read replica configuration for read scaling
- Evaluate connection pool settings under load testing
- Add monitoring/alerting for connection pool exhaustion
- Document recommended settings for different deployment sizes

---

### 7. Repository Pattern with Domain Models Separate from SQLC Models
**Date**: 2025-12-29
**Status**: Accepted

**Context**: SQLC generates database models in `db/sqlc/` package, but Clean Architecture requires domain models to be independent of infrastructure. Need a way to use SQLC without coupling domain layer to database implementation details.

**Decision**: Implement repository pattern with model conversion:
- **Domain Layer** (`internal/domain/`): Define pure business entities (User) and repository interfaces (UserRepository)
- **Repository Layer** (`internal/repository/`): Implement domain interfaces by wrapping SQLC-generated code
- **Model Conversion**: Repository layer converts between `sqlc.User` ↔ `domain.User`
- **Database Package** (`internal/database/`): Contains connection pool initialization, separate from repositories
- **UserRepository Interface**: Defines `Create()` and `GetByEmail()` methods in domain layer

**Consequences**:
- **Positive**: Domain layer remains database-agnostic and testable, can swap SQLC for other implementations without changing domain, follows Dependency Inversion Principle, clear separation between business logic and data access
- **Negative**: Requires model conversion boilerplate (toDomainUser function), slight duplication between domain.User and sqlc.User structs, extra layer of indirection
- **Trade-off**: The architectural purity and testability justify the minimal conversion overhead. For POC with simple models, conversion is trivial (6 field mappings).

**POC → Production Steps**:
- Add repository methods as needed (Update, Delete, List, etc.)
- Implement pagination for List() operations to avoid loading all records in memory
- Implement transaction support using SQLC's WithTx() method
- Separate User model to User and UserForAuth. User should not pass hashed password around
- Add repository-level caching for frequently accessed data
- Create repository interfaces for other domain entities (projects, tasks, etc.)
- Add integration tests for repositories against real database
- Consider using generics to reduce conversion boilerplate if models grow complex
- Add repository method documentation with examples
- Implement soft delete pattern if needed for audit trail

---

### 8. pgxpool for Database Connection Pooling
**Date**: 2025-12-29
**Status**: Accepted

**Context**: Need to choose between `database/sql` with `lib/pq` driver versus pgx native library for PostgreSQL connections. The application is committed to PostgreSQL and needs production-grade performance and connection pooling.

**Decision**: Use `pgxpool` (pgx/v5) directly instead of `database/sql` with `lib/pq`:
- **Database Layer**: Use `pgxpool.Pool` for connection pooling
- **SQLC Configuration**: Set `sql_package: "pgx/v5"` to generate pgx-compatible code
- **Type Handling**: Use pgx native types (e.g., `pgtype.Timestamp`) with conversion to domain types
- **Connection Pool**: Configure via `pgxpool.Config` with MaxConns, MinConns, health checks
- **App Structure**: Store `*pgxpool.Pool` in App struct instead of `*sql.DB`

**Consequences**:
- **Positive**: 2-3x faster performance than lib/pq, actively maintained (lib/pq is maintenance mode only), better PostgreSQL-specific feature support (LISTEN/NOTIFY, better JSON handling, COPY protocol), superior connection pool management optimized for PostgreSQL, built-in prepared statement caching, native batch operations support, automatic health checks
- **Negative**: PostgreSQL-specific (no database portability), slightly different API from database/sql (but SQLC abstracts this), pgx types (pgtype.Timestamp) require conversion to standard Go types (time.Time)
- **Trade-off**: Loss of theoretical database portability is acceptable since we're committed to PostgreSQL. Performance gains and PostgreSQL-native features justify the choice for production readiness.

**POC → Production Steps**:
- Add connection pool metrics monitoring and expose via observability endpoints
  - Use pgxpool.Pool.Stat() to track: AcquireCount, AcquiredConns, IdleConns, MaxConns
  - Export metrics to Prometheus/monitoring system
  - Set up alerts for pool exhaustion
- Implement AfterConnect hooks for session parameter configuration
  - Set timezone, statement_timeout, application_name
  - Configure search_path for multi-tenant applications
  - Enable query logging for debugging
- Add batch insert operations for bulk data using pgx.Batch
- Configure statement cache size based on query patterns
- Implement LISTEN/NOTIFY for real-time features if needed
- Add pgx query logger for slow query detection
- Configure SSL/TLS certificate verification for production
- Test connection pool behavior under load (exhaustion, recovery)
- Add context-based query timeouts for long-running queries
- Consider using pgx.CopyFrom for bulk imports

---

### 9. JWT Authentication with HS256 for User Login
**Date**: 2025-12-29
**Status**: Accepted

**Context**: After implementing user registration with password hashing, need authentication mechanism for users to log in and access protected resources. JWT (JSON Web Tokens) is industry standard for stateless authentication in RESTful APIs. Need to choose between symmetric (HS256) vs asymmetric (RS256) signing, and determine appropriate token expiration for POC testing.

**Decision**: Implement JWT-based authentication with the following:
- **Library**: `github.com/golang-jwt/jwt/v5` (most popular, well-maintained Go JWT library)
- **Signing Algorithm**: HS256 (HMAC with SHA-256) for POC
- **Token Expiration**: 15 minutes (short enough to test expiration during manual testing)
- **Token Claims**: Include user ID, email, role, plus standard JWT claims (exp, iat, nbf)
- **Security Pattern**: Return generic "invalid credentials" error for both wrong email and wrong password (prevents email enumeration)
- **Architecture**: `TokenGenerator` interface in domain layer, JWT implementation in security layer
- **Use Case Flow**: Login validates credentials → generates JWT → returns token + user (without password)

**Consequences**:
- **Positive**: Industry-standard authentication, stateless (no session storage needed), works well with REST APIs, easy to implement and test, 15-minute expiration allows actual testing of token refresh flow, generic error messages prevent user enumeration attacks, JWT can be validated by any service with the secret key, includes user claims for authorization decisions
- **Negative**: HS256 requires shared secret across all services (single point of compromise), tokens cannot be revoked before expiration without additional infrastructure (blacklist/whitelist), 15-minute expiration may be too short for production UX (requires frequent re-authentication), token size larger than opaque tokens (contains claims), vulnerable to XSS if stored in localStorage (client-side concern)
- **Trade-off**: HS256 is simpler for POC with single backend service. Short expiration is acceptable for POC testing and actually beneficial for validating refresh token logic. The security vs usability balance is appropriate for proof-of-concept stage.

**POC → Production Steps**:
- Evaluate switching to RS256 (asymmetric) for distributed systems and microservices
  - RS256 allows token verification without sharing private key
  - Public key can be distributed to all services
  - More secure for multi-service architectures
- Implement Auth-Related config: Auth related config like cost should be set through config
- Implement refresh token mechanism
  - Long-lived refresh tokens (7-30 days)
  - Store refresh tokens securely (database with expiration)
  - Endpoint to exchange refresh token for new access token
  - Rotate refresh tokens on use (one-time use pattern)
- Increase access token expiration to 1-4 hours for better UX
- Add token revocation/blacklist system
  - Redis-based blacklist for logged-out tokens
  - Check blacklist on protected route middleware
  - Cleanup expired entries automatically
- Implement rate limiting for login endpoint
  - Prevent brute force attacks (5 attempts per 15 minutes per IP)
  - Add exponential backoff after failed attempts
  - Consider CAPTCHA after N failed attempts
- Add account lockout mechanism
  - Lock account after N failed login attempts (e.g., 10 attempts)
  - Require email verification or admin unlock
  - Log all failed attempts for security monitoring
- Store JWT secret in secure secret management system (HashiCorp Vault, AWS Secrets Manager)
  - Never commit secrets to version control
  - Rotate secrets periodically
  - Use different secrets per environment
- Add audit logging for authentication events
  - Log all login attempts (success and failure)
  - Track token generation and validation
  - Monitor for suspicious patterns (credential stuffing, account takeover)
- Consider implementing Multi-Factor Authentication (MFA)
  - TOTP (Time-based One-Time Password) support
  - SMS/Email verification codes
  - Backup codes for account recovery
- Add session management and device tracking
  - Track active sessions per user
  - Allow users to view and revoke sessions
  - Detect suspicious login locations/devices
- Implement token fingerprinting to prevent token theft
  - Bind tokens to specific user agent/IP (with caution for mobile users)
  - Detect token reuse from different locations
- Add CSRF protection for cookie-based token storage
- Consider implementing OAuth 2.0/OIDC for third-party integrations

---

### 10. Request Timeout Middleware with http.TimeoutHandler
**Date**: 2025-12-31 (Updated: 2026-01-06)
**Status**: Accepted

**Context**: Database operations could hang indefinitely without timeouts. The context flows from handlers → usecases → repository → database, but no timeout was being set anywhere in the chain, risking indefinite blocking on slow queries or network issues. Initial implementation using only `context.WithTimeout` was incomplete - it created a timeout context but didn't handle the timeout event (no response to client, race conditions, handler continues running).

**Decision**: Implement handler-level timeout middleware using `http.TimeoutHandler`:
- **Middleware Pattern**: Use Go's built-in `http.TimeoutHandler` which properly handles timeout scenarios
- **Configuration**: `SERVER_HANDLER_TIMEOUT=10` (seconds) in environment config
- **Placement**: Applied first in middleware chain before logger and CORS
- **Timeout Behavior**: When timeout occurs, sends `503 Service Unavailable` response, cancels handler context (stops context-aware operations like DB queries), prevents writes to response writer after timeout
- **Context Flow**: Timeout context automatically propagates to all downstream operations (usecases, repositories, database queries)

**Consequences**:
- **Positive**: All requests get automatic timeout protection with proper error responses, context-based cancellation works with pgx natively, prevents race conditions and incomplete responses, configurable per environment, prevents resource exhaustion from hanging connections, battle-tested standard library implementation handles edge cases
- **Negative**: Single timeout applies to all routes (simple and complex operations share same limit), requires choosing appropriate timeout value that works for slowest legitimate operation, 503 status code is generic (doesn't distinguish between different timeout causes)
- **Trade-off**: Handler-level timeout with `http.TimeoutHandler` is simpler than custom async timeout handling and sufficient for POC. Single global timeout is acceptable for POC since all current endpoints have similar performance profiles. Can add route-specific timeouts later if needed.

**POC → Production Steps**:
- **Implement route-specific timeouts using chi router groups** (Recommended approach):
  - Remove global timeout middleware from router initialization
  - Group routes by performance characteristics (fast: 5s, standard: 15s, long-running: 60s, no-timeout: streaming)
  - Apply different timeout middleware to each group using `r.Group()` with `r.Use(middleware.Timeout(duration))`
  - Example: health checks 5s, auth endpoints 15s, report generation 60s, webhooks/SSE no timeout
- **Alternative approaches for fine-grained control**:
  - Use `r.With(middleware.Timeout(duration))` for per-route inline timeouts
  - Move timeout configuration to environment variables or config file (e.g., `TIMEOUT_HEALTH=5s`, `TIMEOUT_STANDARD=15s`)
  - Implement handler-level timeouts using `context.WithTimeout` for operations needing stricter limits than route timeout
- **Enhance timeout error responses**:
  - Replace generic "Request timeout" message with structured JSON error responses
  - Include request ID and support contact information for debugging
  - Log timeout events with request details for monitoring
- **Add timeout monitoring and observability**:
  - Track timeout occurrences per endpoint using metrics
  - Alert on operations approaching timeout limits (e.g., >80% of timeout duration)
  - Log slow requests that consume significant portion of timeout budget
  - Monitor P95/P99 response times to tune timeout values
- **Context-aware timeout handling**:
  - Ensure all I/O operations (DB, HTTP client, external APIs) respect context cancellation
  - Add explicit context checks in long-running loops or computations
  - Implement graceful degradation for partial results when timeout occurs
- **Consider advanced patterns**:
  - Add circuit breaker pattern for downstream services
  - Implement request hedging (send duplicate requests after timeout threshold)
  - Add per-user timeout policies (premium users get longer timeouts)

---

### 11. In-Memory Sliding Window Rate Limiter
**Date**: 2026-01-03
**Status**: Accepted

**Context**: Need to protect API endpoints from abuse, brute force attacks, and excessive traffic. Without rate limiting, attackers can attempt unlimited login attempts, enumerate user emails, or overwhelm the server with requests. Rate limiting is especially critical for authentication endpoints (/login, /register) and public APIs.

**Decision**: Implement custom in-memory sliding window rate limiter in `pkg/ratelimit`:
- **Algorithm**: Sliding window with timestamp-based request tracking per IP and path
- **Data Structure**: `map[IP]map[Path][]time.Time` stores request history
- **Concurrency**: Protected by `sync.Mutex` for thread-safety
- **Cleanup Strategy**: Lazy filtering - remove expired requests when checking limits (no background goroutines)
- **Configuration**: Hardcoded rate limits in middleware initialization for POC
- **Limits Structure**: Per-path configuration with `MaxRequests` and `Window` duration
- **Architecture**: Generic limiter in `pkg/ratelimit`, HTTP middleware wrapper in `internal/delivery/http/middleware`

**Consequences**:
- **Positive**: Simple implementation with no external dependencies, accurate sliding window (not fixed window), thread-safe and suitable for concurrent requests, lazy cleanup prevents goroutine bloat, reusable package design (pkg/ vs internal/), low latency (in-memory lookups)
- **Negative**: Rate limit state lost on server restart, not suitable for distributed/multi-instance deployments, hardcoded limits require code changes, unbounded memory growth if many unique IPs (no IP cleanup implemented), no persistent storage or cross-instance coordination
- **Trade-off**: In-memory approach is acceptable for POC with single server instance. Hardcoded limits allow faster development and testing. Can migrate to Redis-backed distributed limiter for production.

**POC → Production Steps**:
- **Configuration Management**:
  - Move rate limits from hardcoded values to environment variables (e.g., `RATE_LIMIT_AUTH_LOGIN=10/5m`)
  - Or use YAML/JSON config file for complex multi-endpoint configuration
  - Support runtime configuration updates without restart
- **Distributed Rate Limiting**:
  - Replace in-memory store with Redis for multi-instance deployments
  - Use Redis sorted sets or specialized libraries (e.g., go-redis rate limiter)
  - Implement sliding window using Redis ZREMRANGEBYSCORE + ZADD
- **Account Lockout for Failed Login Attempts**:
  - Track failed login attempts per user (not just per IP)
  - Lock account after N failed attempts (e.g., 3-5 attempts)
  - Implement unlock mechanisms: time-based expiration, email verification, or admin unlock
  - Store lockout state in database with expiration timestamp
  - Differentiate between rate limiting (IP-based, temporary) and account lockout (user-based, persistent)
- **Memory Management**:
  - Implement IP cleanup to remove entries with no recent requests (e.g., older than 24 hours)
  - Add background goroutine for periodic cleanup or implement LRU eviction
  - Set maximum number of tracked IPs to prevent memory exhaustion
- **Observability and Monitoring**:
  - Add metrics for rate limit hits (rejected requests) per endpoint
  - Log rate limit violations with IP, path, and timestamp
  - Monitor memory usage of limiter store
  - Alert on unusual rate limit patterns (potential attacks)
- **Enhanced Rate Limit Strategies**:
  - Implement token bucket or leaky bucket for bursty traffic
  - Add user-based rate limiting (authenticated users get higher limits)
  - Support whitelist/blacklist for trusted IPs or known attackers
  - Implement adaptive rate limiting based on server load
- **HTTP Response Headers**:
  - Add standard rate limit headers (X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset)
  - Return 429 Too Many Requests with Retry-After header
  - Provide clear error messages to clients
- **Testing**:
  - Add unit tests for concurrent access and edge cases
  - Load test rate limiter under high concurrent request volume
  - Verify no memory leaks with sustained traffic
- **Security Enhancements**:
  - Use X-Forwarded-For or X-Real-IP for IP detection behind proxies/load balancers
  - Validate and sanitize IP addresses to prevent spoofing
  - Consider fingerprinting (IP + User-Agent) for more accurate tracking
  - Implement gradual backoff (exponential delay) for repeated violations

---

## Template for New Decisions

```markdown
### N. [Decision Title]

**Date**: YYYY-MM-DD
**Status**: [Proposed|Accepted|Deprecated|Superseded]

**Context**: [What is driving this decision?]

**Decision**: [What are we doing?]

**Consequences**:
- **Positive**: [Good outcomes]
- **Negative**: [Costs/drawbacks]
- **Trade-off**: [Why is this acceptable?]

**POC → Production Steps**:
- [Action item 1: What needs to be hardened/improved]
- [Action item 2: What technical debt needs addressing]
- [Action item 3: What additional features/safeguards are needed]
- [Action item N: Security, performance, or reliability considerations]
```
