# DevNorth Backend - AI Assistant Guide

## Project Context

This is the backend for **DevNorth**, a proof-of-concept project using **Clean Architecture** with Go, PostgreSQL, SQLC, and go-migrate.

**Key Principle**: Balance speed (POC mindset) with code quality (SOLID principles). When in doubt about trade-offs, ask the user.

## Technology Stack

- **Language**: Go 1.25.2
- **Architecture**: Clean Architecture
- **Database**: PostgreSQL
- **Query Builder**: SQLC (type-safe SQL → Go code)
- **Migrations**: go-migrate
- **Config**: Environment variables (.env)

## Project Structure

```
devnorth-back/
├── cmd/api/                    # Application entry point (main.go)
├── internal/                   # Private application code
│   ├── app/                    # App initialization (wires DB, logger, config)
│   ├── database/               # Database connection pool setup
│   ├── domain/                 # Entities, interfaces & business logic (no external dependencies)
│   ├── repository/             # Repository implementations (wraps SQLC, converts models)
│   ├── usecase/                # Business logic (orchestrates domain + repo)
│   └── delivery/http/          # HTTP handlers, routes, middleware
├── db/
│   ├── migrations/             # SQL migration files (.sql)
│   ├── queries/                # SQL query files (.sql for SQLC)
│   └── sqlc/                   # SQLC-generated code (models, queries)
├── config/                     # Configuration structs & loaders
├── pkg/                        # Public reusable packages (if needed)
├── .env.example                # Environment variables template
├── Makefile                    # Common development commands
└── sqlc.yaml                   # SQLC configuration
```

**Dependency Flow**: `Delivery → UseCase → Repository → Domain ← (interface implemented by) Repository`
- Inner layers (domain) never import outer layers
- Domain defines repository interfaces; Repository layer implements them
- Repository converts SQLC models to domain models

## Development Workflow

### Adding a New Feature

Follow this process for any non-trivial feature:

1. **Plan & Get Approval**
   - Outline the approach (which layers, files, trade-offs)
   - Ask clarifying questions if requirements are unclear
   - Get user approval before implementing

2. **Database Changes** (if needed)
   ```bash
   make migrate-create name=add_users_table
   # Edit db/migrations/000001_add_users_table.up.sql
   # Edit db/migrations/000001_add_users_table.down.sql
   make migrate-up
   ```

3. **Add Queries** (if needed)
   ```bash
   # Create db/queries/users.sql with SQL queries
   make sqlc  # Generates Go code in db/sqlc
   ```

4. **Implement in Layers** (bottom-up)
   - **Domain**: Define entities (`internal/domain/user.go`) and repository interface (`internal/domain/user_repository.go`)
   - **Repository**: Implement domain interface, wrap SQLC, convert models (`internal/repository/user_repository.go`)
   - **UseCase**: Implement business logic using domain interfaces (`internal/usecase/user_usecase.go`)
   - **Delivery**: Add HTTP handlers (`internal/delivery/http/user_handler.go`)

5. **Test & Verify**
   ```bash
   make test
   make run
   ```

### Common Commands

See `Makefile` for all commands:
- `make run` - Run the application
- `make build` - Build binary to `bin/api`
- `make test` - Run tests
- `make sqlc` - Generate Go code from SQL queries
- `make migrate-up` - Apply migrations
- `make migrate-down` - Rollback last migration
- `make migrate-create name=<name>` - Create new migration

## Working Principles

### 1. Plan Before Acting

Always follow: **Plan → Seek Approval → Implement Incrementally → Check-in Between Steps**

Never make large changes without user visibility and approval.

### 2. Ask Questions First

If requirements are unclear or multiple approaches exist, **ask before implementing**. Clarify:
- Which approach to take (best practice vs. POC speed)
- Expected behavior and edge cases
- Trade-offs and priorities

### 3. Provide Honest Feedback

Speak up when you identify:
- Flaws or potential issues
- Better alternatives
- Trade-offs the user should consider
- Technical debt being created

Be direct and objective - honest feedback improves outcomes.

### 4. POC Mindset with Code Quality

**This is a POC** - prioritize speed and simplicity:
- Focus on core functionality, not edge cases
- Avoid over-engineering (no premature abstractions)
- Perfect is the enemy of done

**But maintain code quality**:
- Follow SOLID principles
- Write clean, readable code
- Use meaningful names
- Keep it simple (simple ≠ sloppy)

**Balance**: When facing trade-offs, present options and let the user decide.

## Go Code Conventions

- **Naming**: Use Go conventions (`UserService`, not `user_service`)
- **Errors**: Return errors, don't panic (except in main setup)
- **Interfaces**: Keep small and focused (ISP - Interface Segregation)
- **Context**: Pass `context.Context` as first parameter for I/O operations
- **Nil checks**: Check errors immediately after function calls
- **Exports**: Only export what's needed (use internal/ for private code)

## SQLC Guidelines

1. **Write SQL in `db/queries/`**: Create `.sql` files with queries
2. **Use SQLC annotations**:
   ```sql
   -- name: GetUser :one
   SELECT * FROM users WHERE id = $1;

   -- name: ListUsers :many
   SELECT * FROM users ORDER BY created_at DESC;

   -- name: CreateUser :one
   INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;
   ```
3. **Generate code**: Run `make sqlc` to generate type-safe Go code in `db/sqlc/`
4. **Wrap in repository layer**: Create repository implementations in `internal/repository/` that:
   - Implement domain repository interfaces
   - Wrap SQLC's generated `Queries` type
   - Convert between SQLC models and domain models

## Repository Pattern

**Architecture**: Domain defines interfaces, Repository implements them

1. **Define repository interface in domain layer**:
   ```go
   // internal/domain/user_repository.go
   type UserRepository interface {
       Create(ctx context.Context, email, hashedPassword string, role UserRole) (*User, error)
       GetByEmail(ctx context.Context, email string) (*User, error)
   }
   ```

2. **Implement in repository layer**:
   ```go
   // internal/repository/user_repository.go
   type userRepository struct {
       queries *sqlc.Queries
   }

   func NewUserRepository(db *sql.DB) domain.UserRepository {
       return &userRepository{queries: sqlc.New(db)}
   }

   func (r *userRepository) Create(...) (*domain.User, error) {
       sqlcUser, err := r.queries.CreateUser(ctx, params)
       return toDomainUser(sqlcUser), err // Convert SQLC → Domain
   }
   ```

3. **Use in upper layers**:
   ```go
   // UseCase or handler
   userRepo := repository.NewUserRepository(app.DB)
   user, err := userRepo.GetByEmail(ctx, "user@example.com")
   ```

**Key Points**:
- Domain models (`domain.User`) are separate from SQLC models (`sqlc.User`)
- Repository converts between them using helper functions (e.g., `toDomainUser()`)
- Domain layer has zero database dependencies
- Easy to mock for testing (just implement the interface)

## Migration Guidelines

1. **Create migration**: `make migrate-create name=descriptive_name`
2. **Edit both files**:
   - `000X_name.up.sql` - Apply changes
   - `000X_name.down.sql` - Revert changes (always make migrations reversible)
3. **Apply**: `make migrate-up`
4. **Rollback if needed**: `make migrate-down`

**Important**: Migrations are sequential and irreversible in production. Test thoroughly.

## Files to Modify vs. Not Modify

**Modify freely**:
- Application code in `cmd/`, `internal/`
- SQL queries in `db/queries/`
- Migration files in `db/migrations/`
- Config files (`.env.example`, `Makefile`)

**Don't modify**:
- SQLC-generated code in `db/sqlc/` (regenerate with `make sqlc` instead)
- `.gitignore`, `go.mod`, `go.sum` (unless explicitly requested)

**Ask before modifying**:
- `sqlc.yaml` (SQLC configuration)
- Project structure changes

## Testing Philosophy

- **POC**: Ask user about test coverage expectations
- **Focus**: Test business logic (use cases), not trivial getters/setters
- **Integration tests**: Consider for critical database operations
- **When in doubt**: Ask "Should I add tests for this?"
