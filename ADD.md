# Architecture Decision Documents (ADD)

This document tracks key architectural decisions made during the development of DevNorth Backend.

## Format

Each decision should include:
- **Date**: When the decision was made
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

### 5. Database Configuration with Connection Pooling and Timeouts
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
- Monitor connection pool metrics (active connections, wait duration, idle connections)
- Tune pool settings based on actual application load and database capacity
- Add connection retry logic with exponential backoff
- Implement query logging and slow query detection using QueryTimeout
- Add database health checks and connection validation
- Consider adding read replica configuration for read scaling
- Evaluate connection pool settings under load testing
- Add monitoring/alerting for connection pool exhaustion
- Document recommended settings for different deployment sizes

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
