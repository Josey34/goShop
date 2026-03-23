# GoShop — Learn Go + AWS with Clean Architecture

> **Build a serverless order processing system in Go — 100% local, 100% free.**
>
> This is a **learning guide**, not a code dump. Follow each phase, use Claude Code to implement, and build real understanding of Go and Clean Architecture.

---

## What You'll Build

An e-commerce order processing system called **GoShop**:

1. REST API receives orders from customers
2. Orders are queued via SQS for async processing
3. Lambda functions handle individual processing steps
4. Step Functions orchestrate the full workflow: validate → calculate → charge → fulfill → notify

```
Client → REST API → SQS Queue → Worker → Step Functions → Lambda chain → Done
              ↕                                                    ↕
           SQLite                                              Results
```

Everything runs locally using LocalStack, SAM CLI, and Step Functions Local. Zero AWS spending.

---

## Architecture Overview

This project uses **Clean Architecture** (Uncle Bob) adapted for Go. The core rule: **dependencies always point inward** — outer layers depend on inner layers, never the reverse.

```
┌─────────────────────────────────────────────────────────┐
│  DELIVERY        handler, middleware, dto, mapper,       │
│                  router, worker, lambda handlers         │
├─────────────────────────────────────────────────────────┤
│  SERVICE         orchestrates usecases, transactions,    │
│                  event publishing, cross-cutting logic    │
├─────────────────────────────────────────────────────────┤
│  USECASE         single-responsibility business actions  │
│                  (CreateOrder, Login, ValidateOrder)      │
├─────────────────────────────────────────────────────────┤
│  REPOSITORY      interfaces (ports) + implementations    │
│                  (adapters): sqlite, sqs, stepfunctions  │
├─────────────────────────────────────────────────────────┤
│  DOMAIN          entity, valueobject, factory, event,    │
│  (innermost)     errors — ZERO external dependencies     │
├─────────────────────────────────────────────────────────┤
│  INFRASTRUCTURE  config, database, pkg (jwt, hasher,     │
│                  logger, idgen, validator, response)      │
└─────────────────────────────────────────────────────────┘
```

### Why These Layers?

| Layer | Why It Exists | Go Principle |
|-------|---------------|--------------|
| **Domain** | Business rules that never change regardless of framework, DB, or delivery method | Zero imports — pure Go only |
| **Repository** | Interfaces decouple usecases from infrastructure. Swap SQLite → Postgres without touching business logic | Interfaces defined by consumer, not producer |
| **Usecase** | Each file = one action. Easy to test, easy to find, easy to reuse from different delivery mechanisms | Single Responsibility |
| **Service** | When one action needs to orchestrate multiple usecases or manage transactions, the service layer handles it | Composition over inheritance |
| **Delivery** | HTTP, Lambda, SQS worker — all just different entry points calling the same services | Thin adapters |
| **DTO + Mapper** | API contract evolves independently from domain. Request validation stays in delivery layer | Separation of concerns |
| **Factory** | Complex entity construction (ID generation, password hashing, event attachment) lives here, not in handlers | Encapsulated creation |
| **Config** | Single source of truth for environment variables, loaded once at startup | 12-factor app |
| **pkg/** | Reusable utilities that don't belong to any layer — JWT, hashing, logging, ID generation | Shared infrastructure |

---

## How Layers Connect (Request Flow)

```
HTTP Request
  → router              registers route + middleware chain
    → middleware         auth (JWT), logging, CORS, request ID, error mapping
      → handler          parses request body → DTO
        → mapper         converts DTO → usecase input
          → service      orchestrates, manages transactions, publishes events
            → usecase    executes single business action
              → repository interface (PORT)
                → repository adapter (sqlite/sqs/etc.)
                  → domain entity (returned)
          → mapper       converts entity → response DTO
        → response       serialized JSON
  ← HTTP Response
```

---

## Complete Project Structure

```
goshop/
├── go.mod
├── go.sum
├── .env / .env.example
├── Makefile
├── docker-compose.yml
├── template.yaml                          # SAM template (Phase 3)
│
├── domain/                                # ── INNERMOST: Zero external deps ──
│   ├── entity/
│   │   ├── product.go                     # Product + methods (ReduceStock, UpdatePrice)
│   │   ├── order.go                       # Order + methods (CalculateTotal, TransitionTo)
│   │   ├── order_item.go                  # OrderItem + LineTotal()
│   │   ├── customer.go                    # Customer entity
│   │   └── payment.go                     # Payment entity
│   ├── valueobject/
│   │   ├── money.go                       # int64 cents, Add/Multiply/Percentage
│   │   ├── email.go                       # Self-validating, immutable
│   │   ├── phone.go                       # Indonesian format (08xx / +628xx)
│   │   ├── address.go                     # Street, City, Province, PostalCode
│   │   ├── order_status.go                # Enum + state machine transitions
│   │   ├── payment_status.go              # Enum
│   │   └── pagination.go                  # Page, Limit, Offset
│   ├── factory/
│   │   ├── order_factory.go               # Builds Order from items, calculates total, attaches events
│   │   └── customer_factory.go            # Builds Customer with hashed password
│   ├── event/
│   │   ├── event.go                       # DomainEvent interface
│   │   ├── order_created.go
│   │   ├── order_status_changed.go
│   │   ├── payment_processed.go
│   │   └── stock_reduced.go
│   └── errors/
│       ├── errors.go                      # Base DomainError + type checkers
│       ├── not_found.go
│       ├── validation.go                  # Supports field-level errors
│       ├── conflict.go
│       ├── unauthorized.go
│       ├── insufficient_stock.go
│       └── invalid_transition.go
│
├── repository/                            # ── PORTS + ADAPTERS ──
│   ├── product_repository.go              # Interface (PORT)
│   ├── order_repository.go                # Interface (PORT)
│   ├── customer_repository.go             # Interface (PORT)
│   ├── order_queue_repository.go          # Interface (PORT)
│   ├── workflow_repository.go             # Interface (PORT)
│   ├── unit_of_work.go                    # Transaction interface (PORT)
│   ├── sqlite/                            # SQLite adapter
│   │   ├── product_repo.go
│   │   ├── order_repo.go
│   │   ├── customer_repo.go
│   │   └── unit_of_work.go
│   ├── memory/                            # In-memory adapter (testing + Phase 2)
│   │   ├── product_repo.go
│   │   ├── order_repo.go
│   │   ├── customer_repo.go
│   │   └── order_queue.go                 # Logs to stdout, swap to SQS in Phase 4
│   ├── sqs/                               # SQS adapter (Phase 4)
│   │   └── order_queue.go
│   └── stepfunctions/                     # Step Functions adapter (Phase 5)
│       └── workflow_orchestrator.go
│
├── usecase/                               # ── SINGLE-RESPONSIBILITY ACTIONS ──
│   ├── product/
│   │   ├── create_product.go
│   │   ├── get_product.go
│   │   ├── list_products.go
│   │   ├── update_product.go
│   │   └── delete_product.go
│   ├── order/
│   │   ├── create_order.go
│   │   ├── get_order.go
│   │   ├── list_orders.go
│   │   ├── update_order_status.go
│   │   ├── validate_order.go              # Used by Lambda
│   │   └── calculate_total.go             # Used by Lambda
│   ├── auth/
│   │   ├── register.go
│   │   └── login.go
│   ├── payment/
│   │   └── process_payment.go             # Used by Lambda
│   └── workflow/
│       ├── start_order_workflow.go
│       └── get_workflow_status.go
│
├── service/                               # ── ORCHESTRATION ──
│   ├── product_service.go                 # Wraps product usecases
│   ├── order_service.go                   # Wraps order usecases + transaction + events
│   ├── auth_service.go                    # Wraps auth usecases + token generation
│   ├── payment_service.go
│   └── workflow_service.go
│
├── delivery/                              # ── OUTERMOST ──
│   ├── http/
│   │   ├── server.go                      # Gin server + graceful shutdown
│   │   ├── router/
│   │   │   ├── router.go                  # Mounts all route groups + global middleware
│   │   │   ├── product_routes.go
│   │   │   ├── order_routes.go
│   │   │   └── auth_routes.go
│   │   ├── handler/
│   │   │   ├── product_handler.go
│   │   │   ├── order_handler.go
│   │   │   ├── auth_handler.go
│   │   │   └── health_handler.go
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go          # JWT validation → sets customer_id in context
│   │   │   ├── error_middleware.go         # DomainError → HTTP status code mapping
│   │   │   ├── logger_middleware.go        # Request logging with slog
│   │   │   ├── cors_middleware.go
│   │   │   ├── recovery_middleware.go      # Panic recovery
│   │   │   └── request_id.go              # X-Request-ID injection
│   │   └── dto/
│   │       ├── request/
│   │       │   ├── create_product_request.go
│   │       │   ├── update_product_request.go
│   │       │   ├── create_order_request.go
│   │       │   ├── register_request.go
│   │       │   ├── login_request.go
│   │       │   └── pagination_request.go
│   │       ├── response/
│   │       │   ├── product_response.go
│   │       │   ├── order_response.go
│   │       │   ├── auth_response.go
│   │       │   ├── error_response.go
│   │       │   ├── paginated_response.go
│   │       │   └── health_response.go
│   │       └── mapper/
│   │           ├── product_mapper.go       # Entity ↔ DTO
│   │           ├── order_mapper.go
│   │           └── auth_mapper.go
│   ├── lambda/                             # Phase 3
│   │   ├── validate_order/
│   │   │   ├── main.go
│   │   │   ├── handler.go                 # Event → usecase input → response
│   │   │   └── event.go                   # Lambda event/response types
│   │   ├── calculate_total/
│   │   ├── process_payment/
│   │   ├── fulfill_order/
│   │   └── send_notification/
│   └── worker/                             # Phase 4
│       ├── consumer.go                     # Generic SQS polling loop
│       └── order_worker.go                 # Order-specific message handler
│
├── config/
│   └── config.go                           # Load from .env, structs for App/DB/JWT/SQS/SFN
│
├── database/
│   ├── sqlite.go                           # NewSQLiteDB() with WAL + foreign keys
│   ├── migrator.go                         # Reads and runs .sql files in order
│   └── migrations/
│       ├── 001_create_products.sql
│       ├── 002_create_customers.sql
│       ├── 003_create_orders.sql
│       └── 004_create_order_items.sql
│
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                          # Generate + Validate tokens
│   ├── hasher/
│   │   └── bcrypt.go                       # PasswordHasher interface + bcrypt impl
│   ├── logger/
│   │   └── logger.go                       # slog wrapper (JSON in prod, text in dev)
│   ├── idgen/
│   │   └── uuid.go                         # IDGenerator interface + UUID impl
│   ├── response/
│   │   └── json.go                         # Success(), Error(), Paginated() helpers
│   └── validator/
│       └── validator.go                    # Custom validation helpers
│
├── cmd/
│   ├── api/
│   │   └── main.go                         # Wire everything → start HTTP server
│   ├── worker/
│   │   └── main.go                         # Wire everything → start SQS consumer
│   ├── migrate/
│   │   └── main.go                         # Run migrations standalone
│   └── seed/
│       └── main.go                         # Seed test data
│
├── stepfunctions/
│   └── order-workflow.asl.json             # State machine (Phase 5)
│
├── events/                                 # Lambda test events
│   ├── validate-order.json
│   ├── validate-order-invalid.json
│   ├── calculate-total.json
│   └── process-payment.json
│
└── scripts/
    ├── setup-localstack.sh
    └── seed-data.sh
```

---

## Where Does Code Belong? (Decision Guide)

Use this table when you're building and unsure where something goes:

| If the code... | Put it in | Example |
|----------------|-----------|---------|
| Is a business rule true regardless of DB/framework | `domain/entity/` | `Order.CalculateTotal()`, `Product.ReduceStock()` |
| Is an immutable type that validates itself on creation | `domain/valueobject/` | `Money`, `Email`, `OrderStatus` |
| Constructs an entity with injected dependencies (ID gen, hashing) | `domain/factory/` | `OrderFactory.Create()` |
| Is something that happened in the domain worth reacting to | `domain/event/` | `OrderCreated`, `StockReduced` |
| Is a domain-specific error condition | `domain/errors/` | `ErrInsufficientStock`, `ErrInvalidTransition` |
| Defines what data access looks like (not how) | `repository/*.go` | `ProductRepository` interface |
| Actually talks to a database/queue/external service | `repository/sqlite/`, `sqs/`, etc. | `sqlite.ProductRepo.Create()` |
| Is a single application action (one verb) | `usecase/` | `CreateProduct`, `Login`, `ValidateOrder` |
| Orchestrates multiple usecases or adds transaction/logging | `service/` | `OrderService.PlaceOrder()` |
| Shapes what the HTTP request/response looks like | `delivery/http/dto/` | `CreateProductRequest`, `OrderResponse` |
| Converts between DTO and domain entity | `delivery/http/dto/mapper/` | `ToProductResponse()`, `ToCreateOrderInput()` |
| Handles an HTTP endpoint (parse → call service → respond) | `delivery/http/handler/` | `ProductHandler.Create()` |
| Is a cross-cutting HTTP concern | `delivery/http/middleware/` | JWT auth, CORS, error mapping, logging |
| Registers routes and applies middleware groups | `delivery/http/router/` | `SetupProductRoutes()` |
| Handles a Lambda event | `delivery/lambda/` | `validate_order/handler.go` |
| Polls SQS and processes messages | `delivery/worker/` | `consumer.go`, `order_worker.go` |
| Is app configuration from env vars | `config/` | `Config.Load()` |
| Is DB connection + migration | `database/` | `NewSQLiteDB()`, `RunMigrations()` |
| Is a reusable utility not tied to any layer | `pkg/` | JWT service, bcrypt hasher, UUID generator |
| Wires all layers together (dependency injection) | `cmd/` | `main.go` — the ONLY file that imports everything |

---

## Dependency Rules

These are **strict** — verify them before every commit:

1. **`domain/`** imports NOTHING from other project packages. Only Go stdlib.
2. **`usecase/`** imports `domain/` and `repository/` (interfaces only).
3. **`service/`** imports `usecase/` and `domain/`.
4. **`repository/sqlite/`** etc. imports `domain/` and `repository/` (to implement interfaces).
5. **`delivery/`** imports `service/`, `domain/errors/`, `dto/`, `pkg/`. Never imports `repository/` directly.
6. **`pkg/`** imports nothing from the project (only Go stdlib + external libs).
7. **`cmd/main.go`** is the ONLY file that imports ALL layers to wire them together.

---

## Prerequisites

| Tool | Purpose | Install |
|------|---------|---------|
| **Go 1.22+** | Language | [go.dev/dl](https://go.dev/dl/) |
| **Docker** | LocalStack + Step Functions | [docker.com](https://www.docker.com/get-started/) |
| **AWS SAM CLI** | Local Lambda | [SAM install guide](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) |
| **AWS CLI** | Interact with LocalStack | [AWS CLI install](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) |

```bash
# Verify
go version        # 1.22+
docker --version  # 24+
sam --version     # 1.x
aws --version     # 2.x

# Dummy AWS credentials for LocalStack
aws configure
# Access Key: test | Secret: test | Region: ap-southeast-1 | Output: json
```

---

## Phase 1: Domain Layer (Week 1)

> **Goal:** Build the entire domain — entity, value object, factory, event, errors. Zero external deps. 100% unit testable without mocks.

### What You'll Learn
- Go basics: structs, methods, interfaces, error handling, packages
- Value objects (immutable, self-validating types)
- Entity design (constructors that guarantee validity)
- Factory pattern (complex construction with injected dependencies)
- Domain events (things that happened worth reacting to)
- Custom error types
- Table-driven tests

### What to Build (in order)

**1. Domain Errors** (`domain/errors/`)
- Base `DomainError` struct with Code, Message, optional Fields (for field-level validation)
- Separate files for each error type: NotFound, Validation (with field errors), Conflict, Unauthorized, InsufficientStock, InvalidTransition
- Type-checking helpers: `IsNotFound(err)`, `IsValidation(err)`, etc.

**2. Value Objects** (`domain/valueobject/`)
- `Money` — uses int64 (not float64!) to avoid precision bugs. IDR has no subunit, so 1 = Rp 1. Methods: `Add`, `Multiply`, `Percentage`, `IsZero`. Constructor validates non-negative.
- `Email` — lowercased, trimmed, regex validated. Immutable.
- `Phone` — Indonesian format: 08xx or +628xx, 10-15 digits.
- `Address` — Street, City, Province, PostalCode.
- `OrderStatus` — string enum with a **state machine**: `validTransitions` map defines which transitions are legal. `CanTransitionTo()` and `TransitionTo()` enforce the machine. States: pending → validated → calculated → paid → fulfilled → completed. Any state can go to failed. Pending can go to cancelled.
- `PaymentStatus` — simple enum: pending, success, failed, refunded.
- `Pagination` — Page, Limit, Offset. Constructor validates limit 1-100.

**3. Entities** (`domain/entity/`)
- `Product` — constructor validates name, price > 0, stock >= 0. Methods: `ReduceStock(qty)` (returns InsufficientStock error), `RestoreStock(qty)`, `UpdatePrice()`, `UpdateName()`. All methods update `UpdatedAt`.
- `OrderItem` — ProductID, Quantity, Price (snapshot at order time). Method: `LineTotal()` returns Price × Quantity.
- `Order` — ID, CustomerID, Items, Subtotal/Tax/Total (Money), Status (OrderStatus). Constructor validates non-empty items. Methods: `CalculateTotal()` (subtotal + PPN 11% tax), `TransitionTo(status)` (uses OrderStatus state machine), `AddEvent()` / `PullEvents()` for domain events.
- `Customer` — ID, Name, Email (value object), Phone (value object), PasswordHash. Constructor validates required fields.
- `Payment` — ID, OrderID, Amount, Status, Method, PaidAt.

**4. Factories** (`domain/factory/`)
- `OrderFactory` — depends on `IDGenerator` interface (injected). Method: `Create(customerID, items)` → generates IDs, builds OrderItems, creates Order, calculates total, attaches `OrderCreated` event.
- `CustomerFactory` — depends on `IDGenerator` + `PasswordHasher` interfaces. Method: `Create(name, email, phone, password)` → validates, hashes password, builds Customer.

**5. Domain Events** (`domain/event/`)
- `DomainEvent` interface: `EventName() string`, `OccurredAt() time.Time`
- Concrete events: `OrderCreated`, `OrderStatusChanged`, `PaymentProcessed`, `StockReduced`

### Key Design Decisions
- **Money uses int64, not float64** — Rp 50,000 is stored as `50000`. This avoids floating-point precision bugs that would silently corrupt totals.
- **OrderStatus has a state machine** — instead of allowing any status change, the domain enforces valid transitions. A paid order can't go back to pending.
- **Factories take interfaces, not implementations** — `OrderFactory` takes `IDGenerator`, not `uuid.Generator`. This means the domain stays pure.
- **Entities collect events, services publish them** — `order.AddEvent(OrderCreated{...})` stores the event. The service layer calls `order.PullEvents()` after persistence and publishes them.

### Testing Strategy
- Domain tests need **zero mocks** — the domain has no external dependencies
- Use table-driven tests for value objects (many input/output combinations)
- Test entity state transitions exhaustively
- Test factory with a simple stub IDGenerator that returns predictable IDs

### Verify
```bash
go test ./domain/... -v -count=1
```

---

## Phase 2: Everything Else for HTTP API (Week 2–3)

> **Goal:** Build all remaining layers — config, database, pkg, repository, usecase, service, dto, mapper, handler, middleware, router, server, main.go.

### What You'll Learn
- Config from env vars (.env file)
- SQLite with `database/sql` + migrations
- Repository pattern (interface + implementation)
- Usecase pattern (single-responsibility)
- Service pattern (orchestration)
- DTO + Mapper pattern (decouple API shape from domain)
- Gin framework (handlers, middleware, router)
- JWT authentication
- Dependency injection in Go (no framework — just constructors in main.go)

### Dependencies to Install
```bash
go get github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
```

### What to Build (in order)

**1. Config** (`config/`)
- Struct with nested configs: App (env, port), DB (path), JWT (secret, expiry), SQS (endpoint, queue URL), SFN (endpoint, ARN), AWS (region)
- Load from env vars with sensible defaults
- `.env.example` with all vars documented

**2. Database** (`database/`)
- `NewSQLiteDB(path)` — opens connection with WAL mode + foreign keys enabled
- `RunMigrations(db, dir)` — reads `.sql` files from a directory, sorts by name, executes in order
- Migration files: 001_create_products, 002_create_customers, 003_create_orders, 004_create_order_items
- **Important**: store Money as `INTEGER` (price_amount) + `TEXT` (price_currency), not as REAL

**3. Shared Packages** (`pkg/`)
- `pkg/idgen/` — `IDGenerator` interface + UUID implementation. This implements the same interface the factory expects.
- `pkg/hasher/` — `PasswordHasher` interface + bcrypt impl with `Hash()` and `Compare()`. Also implements the factory's interface.
- `pkg/jwt/` — JWT service with `Generate(customerID, email)` and `Validate(token) → Claims`
- `pkg/logger/` — thin slog wrapper. JSON handler in production, text handler in development.
- `pkg/response/` — Gin JSON helpers: `Success()`, `Error()`, `ErrorWithFields()`, `Paginated()`
- `pkg/validator/` — custom validation helpers if needed

**4. Repository Interfaces** (`repository/*.go`)
- `ProductRepository` — Create, FindByID, FindAll(pagination), Update, Delete
- `OrderRepository` — Create, FindByID, FindByCustomerID(pagination), UpdateStatus
- `CustomerRepository` — Create, FindByID, FindByEmail, ExistsByEmail
- `OrderQueue` — Enqueue(order)
- `WorkflowOrchestrator` — StartOrderWorkflow, GetWorkflowStatus
- `UnitOfWork` — Begin, Commit, Rollback + access to repos within a transaction

**5. Repository Implementations**
- `repository/sqlite/` — implements all repo interfaces using `database/sql`. Map between domain entities and SQL columns. Return domain errors (e.g., `sql.ErrNoRows` → `domainErr.NewNotFound()`).
- `repository/memory/` — in-memory implementations for testing. `InMemoryOrderQueue` just logs — swap to SQS in Phase 4.

**6. Usecases** (`usecase/`)
- Each usecase is a struct with a constructor (`New*UseCase(deps...)`) and an `Execute(ctx, input) (output, error)` method.
- `product/`: CreateProduct, GetProduct, ListProducts, UpdateProduct, DeleteProduct
- `order/`: CreateOrder (validates stock → factory creates → persist → enqueue), GetOrder, ListOrders, UpdateOrderStatus, ValidateOrder, CalculateTotal
- `auth/`: Register (check duplicate → factory creates → persist), Login (find by email → compare password → generate token)
- **Usecase imports**: only `domain/` and `repository/` interfaces

**7. Services** (`service/`)
- `ProductService` — wraps product usecases, adds logging
- `OrderService` — wraps order usecases, adds logging, publishes domain events after persistence
- `AuthService` — wraps auth usecases, adds logging
- Services are what handlers call. They're the "public API" of the backend.

**8. DTOs** (`delivery/http/dto/`)
- `request/` — one struct per endpoint: CreateProductRequest, UpdateProductRequest, CreateOrderRequest, RegisterRequest, LoginRequest, PaginationRequest. Use Gin's `binding:"required"` tags.
- `response/` — one struct per response shape: ProductResponse, OrderResponse, AuthResponse, ErrorResponse, PaginatedResponse, HealthResponse

**9. Mappers** (`delivery/http/dto/mapper/`)
- `product_mapper.go` — `ToCreateProductInput(req) → usecase input`, `ToProductResponse(entity) → response`, `ToProductListResponse([]entity) → []response`
- `order_mapper.go` — same pattern
- `auth_mapper.go` — same pattern
- Mappers keep conversion logic out of handlers

**10. Middleware** (`delivery/http/middleware/`)
- `auth_middleware.go` — extract Bearer token, validate with JWT service, set `customer_id` in Gin context
- `error_middleware.go` — map `DomainError` code to HTTP status (NOT_FOUND→404, VALIDATION→400, CONFLICT→409, UNAUTHORIZED→401, INSUFFICIENT_STOCK→422)
- `logger_middleware.go` — log method, path, status, duration using slog
- `cors_middleware.go` — allow all origins for local dev
- `recovery_middleware.go` — catch panics, return 500
- `request_id.go` — generate/forward X-Request-ID header

**11. Handlers** (`delivery/http/handler/`)
- Each handler takes a service in its constructor
- Pattern: parse request → DTO → mapper → service call → mapper → response
- `product_handler.go` — Create, List, GetByID, Update, Delete
- `order_handler.go` — Create, List, GetByID
- `auth_handler.go` — Register, Login
- `health_handler.go` — simple status check

**12. Router** (`delivery/http/router/`)
- `router.go` — main setup: applies global middleware, mounts route groups
- `product_routes.go` — `/api/v1/products` group
- `order_routes.go` — `/api/v1/orders` group
- `auth_routes.go` — `/auth` group (public, no JWT middleware)

**13. Server** (`delivery/http/server.go`)
- Creates Gin engine, calls router setup, starts HTTP server with graceful shutdown (listens for SIGINT/SIGTERM)

**14. Main** (`cmd/api/main.go`)
- This is where **dependency injection** happens — the ONLY file that imports all layers
- Order: config → logger → database → migrations → pkg utilities → factories → repos → usecases → services → handlers → router → server

### Key Pattern: Swapping Implementations

The whole point of clean architecture shows up here. In Phase 2, your `main.go` wires:
```
orderQueue := memory.NewInMemoryOrderQueue()  // logs to stdout
```

In Phase 4, you change ONE line:
```
orderQueue := sqs.NewOrderQueue(sqsClient, queueURL)  // sends to real SQS
```

Zero changes to usecases, services, or handlers. That's the payoff.

### Testing Strategy
- **Repository tests**: test SQLite repos with a real in-memory SQLite DB (`":memory:"`)
- **Usecase tests**: inject mock/in-memory repositories
- **Handler tests**: use `httptest.NewRecorder()` + a real Gin context
- **Integration test**: start the server, make real HTTP calls

### API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/auth/register` | No | Register customer |
| POST | `/auth/login` | No | Login → JWT token |
| GET | `/api/v1/products` | Yes | List products (paginated) |
| GET | `/api/v1/products/:id` | Yes | Get product by ID |
| POST | `/api/v1/products` | Yes | Create product |
| PUT | `/api/v1/products/:id` | Yes | Update product |
| DELETE | `/api/v1/products/:id` | Yes | Delete product |
| POST | `/api/v1/orders` | Yes | Create order |
| GET | `/api/v1/orders` | Yes | List customer orders |
| GET | `/api/v1/orders/:id` | Yes | Get order by ID |

### Verify
```bash
go run cmd/api/main.go

# Test
curl http://localhost:8080/health
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Josey","email":"josey@test.com","phone":"08123456789","password":"secret123"}'
```

---

## Phase 3: AWS Lambda — Local with SAM (Week 3–4)

> **Goal:** Each Lambda is a **delivery adapter** — a thin layer that maps an event to a usecase call.

### What You'll Learn
- Lambda handler pattern in Go
- SAM template (template.yaml)
- `sam local invoke` and `sam local start-lambda`
- How Lambdas reuse domain logic (same entities, same usecases)

### What to Build

**1. Lambda Event/Response Types** (`delivery/lambda/*/event.go`)
- Each Lambda has its own input/output types matching the Step Functions workflow

**2. Lambda Handlers** (`delivery/lambda/*/handler.go`)
- Map incoming event → usecase input → call usecase → return response
- `validate_order` — uses `entity.NewOrder()` for validation, returns is_valid + errors
- `calculate_total` — uses `Order.CalculateTotal()`, returns total with PPN 11%
- `process_payment` — simulates payment (random success/fail), returns payment_status
- `fulfill_order` — transitions order to fulfilled
- `send_notification` — logs notification (simulated)

**3. Lambda Entry Points** (`delivery/lambda/*/main.go`)
- Each has a `main()` that calls `lambda.Start(handler)`

**4. SAM Template** (`template.yaml`)
- Defines all 5 Lambda functions with `provided.al2023` runtime
- Each uses a Makefile build method for Go cross-compilation

**5. Test Events** (`events/`)
- JSON files for each Lambda: valid input, invalid input, edge cases

### Running
```bash
sam build
sam local invoke ValidateOrderFunction --event events/validate-order.json
sam local start-lambda --port 3001   # for Step Functions to call
```

---

## Phase 4: SQS with LocalStack (Week 4–5)

> **Goal:** Implement `OrderQueue` with real SQS. Zero changes to usecases — just swap the adapter.

### What You'll Learn
- SQS concepts: queues, messages, visibility timeout, dead letter queues
- `aws-sdk-go-v2`
- LocalStack (free, runs in Docker)
- The clean architecture payoff: swap in-memory → SQS with one line change

### What to Build

**1. Docker Compose** (`docker-compose.yml`)
- LocalStack service (SQS, S3, DynamoDB) on port 4566
- Step Functions Local on port 8083

**2. Setup Script** (`scripts/setup-localstack.sh`)
- Create DLQ: `goshop-orders-dlq`
- Create main queue: `goshop-orders` with DLQ redrive policy (max 3 retries)
- Create notification queue: `goshop-notifications`

**3. SQS Adapter** (`repository/sqs/order_queue.go`)
- Implements `repository.OrderQueue` interface
- Uses `aws-sdk-go-v2/service/sqs` with custom endpoint for LocalStack
- Serializes order to JSON, sends as SQS message

**4. Worker** (`delivery/worker/`)
- `consumer.go` — generic SQS polling loop: receive → handle → delete. Long polling (20s). Graceful shutdown.
- `order_worker.go` — deserializes message → calls `StartOrderWorkflow` usecase (Phase 5) or logs for now

**5. Worker Entry Point** (`cmd/worker/main.go`)
- Wire config → SQS client → consumer → start polling

**6. Update API main.go**
- Replace `memory.NewInMemoryOrderQueue()` with `sqs.NewOrderQueue(endpoint, queueURL)`

### Running
```bash
docker compose up -d localstack
./scripts/setup-localstack.sh
go run cmd/api/main.go     # terminal 1
go run cmd/worker/main.go  # terminal 2
# Create an order → watch worker pick it up
```

---

## Phase 5: Step Functions — Local Orchestration (Week 5–6)

> **Goal:** Orchestrate the full order workflow with a state machine.

### What You'll Learn
- Amazon States Language (ASL)
- State types: Task, Choice, Parallel, Wait, Fail, Succeed
- Error handling + retry in workflows
- Step Functions Local (Docker)

### What to Build

**1. State Machine** (`stepfunctions/order-workflow.asl.json`)
```
ValidateOrder → IsValid? → CalculateTotal → ProcessPayment → PaymentOK?
  → FulfillAndNotify (Parallel: FulfillOrder + SendNotification) → OrderCompleted
  → OrderFailed (on any error)
```
- Retry on TaskFailed (2 attempts, backoff)
- Catch all errors → OrderFailed

**2. Step Functions Adapter** (`repository/stepfunctions/workflow_orchestrator.go`)
- Implements `repository.WorkflowOrchestrator`
- Uses `aws-sdk-go-v2/service/sfn` with custom endpoint

**3. Update Worker**
- When worker receives an order message, call `StartOrderWorkflow` usecase
- Which calls the `WorkflowOrchestrator` adapter

### Running
```bash
docker compose up -d
sam build && sam local start-lambda --port 3001

aws stepfunctions create-state-machine \
  --endpoint-url http://localhost:8083 \
  --name GoShopOrderWorkflow \
  --definition file://stepfunctions/order-workflow.asl.json \
  --role-arn "arn:aws:iam::000000000000:role/DummyRole"

go run cmd/worker/main.go  # starts workflow when order arrives
```

---

## Phase 6: Capstone — Full Integration (Week 6–7)

> **Goal:** Everything works end-to-end with one `make run`.

### Checklist
- [ ] API accepts orders → persists to SQLite + sends to SQS
- [ ] Worker polls SQS → starts Step Function execution
- [ ] Step Function calls: Validate → Calculate → Pay → (Fulfill ∥ Notify)
- [ ] Each Lambda step runs via SAM Local
- [ ] Failed orders go to DLQ
- [ ] All domain logic tested without infrastructure
- [ ] Architecture rules verified (no wrong-direction imports)

### Makefile
```makefile
build:        sam build
infra-up:     docker compose up -d && sleep 5 && ./scripts/setup-localstack.sh
api:          go run cmd/api/main.go
worker:       go run cmd/worker/main.go
lambdas:      sam local start-lambda --port 3001
test:         go test ./... -v -count=1
test-domain:  go test ./domain/... -v
stop:         docker compose down
```

### End-to-End Test
```bash
make infra-up && make build
# Terminal 1: make lambdas
# Terminal 2: make api
# Terminal 3: make worker

# Register + Login + Create Order → watch the entire pipeline flow
```

---

## Architecture Verification Checklist

Run before every commit:

- [ ] `domain/` has ZERO imports from other project packages
- [ ] `usecase/` only imports `domain/` and `repository/` (interfaces)
- [ ] `service/` only imports `usecase/` and `domain/`
- [ ] `delivery/` never imports `repository/` directly
- [ ] `pkg/` imports nothing from the project
- [ ] `cmd/main.go` is the only file that imports all layers
- [ ] All entity construction uses constructors or factories
- [ ] All tests pass: `go test ./... -count=1`

---

## Bonus Challenges

1. **S3 (local)** — `repository.FileStorage` interface, implement with LocalStack S3 for product images
2. **DynamoDB (local)** — `repository.AuditLogger` interface for order processing logs
3. **gRPC delivery** — `delivery/grpc/` calling the same services
4. **Outbox pattern** — ensure order persistence + SQS message are atomic
5. **CQRS** — separate read and write repository interfaces
6. **OpenTelemetry** — distributed tracing across HTTP → SQS → Lambda → Step Functions

---

## Helpful Links

| Resource | URL |
|----------|-----|
| Go by Example | [gobyexample.com](https://gobyexample.com/) |
| Effective Go | [go.dev/doc/effective_go](https://go.dev/doc/effective_go) |
| Go Clean Arch Example | [github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch) |
| Gin Framework | [gin-gonic.com](https://gin-gonic.com/) |
| AWS SDK Go v2 | [aws.github.io/aws-sdk-go-v2](https://aws.github.io/aws-sdk-go-v2/) |
| SAM CLI | [docs.aws.amazon.com/sam](https://docs.aws.amazon.com/serverless-application-model/) |
| LocalStack | [docs.localstack.cloud](https://docs.localstack.cloud/) |
| Step Functions Local | [docs.aws.amazon.com/step-functions](https://docs.aws.amazon.com/step-functions/latest/dg/sfn-local.html) |