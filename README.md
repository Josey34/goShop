# GoShop — Learn Go + AWS with Clean Architecture

> **Build a serverless order processing system in Go — 100% local, 100% free.**
>
> This is a **learning guide**, not a code dump. Follow each phase, use Claude Code to implement, and build real understanding of Go and Clean Architecture.

---

## What You'll Build

An e-commerce order processing system called **GoShop**:

1. REST API receives orders from customers
2. Product images are stored in S3 (local via LocalStack)
3. Orders are queued via SQS for async processing
4. Lambda functions handle individual processing steps
5. Step Functions orchestrate the full workflow: validate → calculate → charge → fulfill → notify

```
Client → REST API → SQS Queue → Worker → Step Functions → Lambda chain → Done
              ↕          ↕
           SQLite       S3
          (data)     (images)
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
│                  (adapters): sqlite, sqs, s3, stepfnc    │
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
│   │   ├── product.go                     # Product + methods (ReduceStock, UpdatePrice, SetImage)
│   │   ├── product_test.go
│   │   ├── order.go                       # Order + methods (CalculateTotal, TransitionTo)
│   │   ├── order_test.go
│   │   ├── order_item.go                  # OrderItem + LineTotal()
│   │   ├── order_item_test.go
│   │   ├── customer.go                    # Customer entity
│   │   ├── customer_test.go
│   │   └── payment.go                     # Payment entity
│   ├── valueobject/
│   │   ├── money.go                       # int64 cents, Add/Multiply/Percentage
│   │   ├── money_test.go
│   │   ├── email.go                       # Self-validating, immutable
│   │   ├── email_test.go
│   │   ├── phone.go                       # Indonesian format (08xx / +628xx)
│   │   ├── phone_test.go
│   │   ├── address.go                     # Street, City, Province, PostalCode
│   │   ├── order_status.go                # Enum + state machine transitions
│   │   ├── order_status_test.go
│   │   ├── image_url.go                   # Validated URL for product images
│   │   ├── payment_status.go              # Enum
│   │   └── pagination.go                  # Page, Limit, Offset
│   ├── factory/
│   │   ├── order_factory.go               # Builds Order from items, calculates total, attaches events
│   │   ├── order_factory_test.go
│   │   ├── customer_factory.go            # Builds Customer with hashed password
│   │   └── customer_factory_test.go
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
│   ├── file_storage_repository.go         # Interface (PORT) — Upload, GetURL, Delete
│   ├── order_queue_repository.go          # Interface (PORT)
│   ├── workflow_repository.go             # Interface (PORT)
│   ├── unit_of_work.go                    # Transaction interface (PORT)
│   ├── sqlite/                            # SQLite adapter
│   │   ├── product_repo.go
│   │   ├── product_repo_test.go
│   │   ├── order_repo.go
│   │   ├── order_repo_test.go
│   │   ├── customer_repo.go
│   │   ├── customer_repo_test.go
│   │   └── unit_of_work.go
│   ├── s3/                                # S3 adapter (LocalStack)
│   │   ├── file_storage.go               # implements FileStorage — upload, presigned URL, delete
│   │   └── file_storage_test.go
│   ├── memory/                            # In-memory adapter (testing + Phase 2)
│   │   ├── product_repo.go
│   │   ├── order_repo.go
│   │   ├── customer_repo.go
│   │   ├── order_queue.go                 # Logs to stdout, swap to SQS in Phase 4
│   │   └── file_storage.go               # Stores files in map, swap to S3 in Phase 4
│   ├── sqs/                               # SQS adapter (Phase 4)
│   │   └── order_queue.go
│   └── stepfunctions/                     # Step Functions adapter (Phase 5)
│       └── workflow_orchestrator.go
│
├── usecase/                               # ── SINGLE-RESPONSIBILITY ACTIONS ──
│   ├── product/
│   │   ├── create_product.go
│   │   ├── create_product_test.go
│   │   ├── get_product.go
│   │   ├── list_products.go
│   │   ├── list_products_test.go
│   │   ├── update_product.go
│   │   ├── delete_product.go
│   │   ├── upload_product_image.go        # Upload image → S3, save URL to product
│   │   └── upload_product_image_test.go
│   ├── order/
│   │   ├── create_order.go
│   │   ├── create_order_test.go
│   │   ├── get_order.go
│   │   ├── list_orders.go
│   │   ├── update_order_status.go
│   │   ├── validate_order.go              # Used by Lambda
│   │   ├── validate_order_test.go
│   │   └── calculate_total.go             # Used by Lambda
│   ├── auth/
│   │   ├── register.go
│   │   ├── register_test.go
│   │   ├── login.go
│   │   └── login_test.go
│   ├── payment/
│   │   └── process_payment.go             # Used by Lambda
│   └── workflow/
│       ├── start_order_workflow.go
│       └── get_workflow_status.go
│
├── service/                               # ── ORCHESTRATION ──
│   ├── product_service.go                 # Wraps product usecases
│   ├── product_service_test.go
│   ├── order_service.go                   # Wraps order usecases + transaction + events
│   ├── order_service_test.go
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
│   │   │   ├── product_handler_test.go
│   │   │   ├── order_handler.go
│   │   │   ├── order_handler_test.go
│   │   │   ├── auth_handler.go
│   │   │   ├── auth_handler_test.go
│   │   │   └── health_handler.go
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go          # JWT validation → sets customer_id in context
│   │   │   ├── auth_middleware_test.go
│   │   │   ├── error_middleware.go         # DomainError → HTTP status code mapping
│   │   │   ├── logger_middleware.go        # Request logging with slog
│   │   │   ├── cors_middleware.go
│   │   │   ├── recovery_middleware.go      # Panic recovery
│   │   │   └── request_id.go              # X-Request-ID injection
│   │   └── dto/
│   │       ├── request/
│   │       │   ├── create_product_request.go
│   │       │   ├── update_product_request.go
│   │       │   ├── upload_image_request.go
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
│   │   ├── jwt.go                          # Generate + Validate tokens
│   │   └── jwt_test.go
│   ├── hasher/
│   │   ├── bcrypt.go                       # PasswordHasher interface + bcrypt impl
│   │   └── bcrypt_test.go
│   ├── logger/
│   │   └── logger.go                       # slog wrapper (JSON in prod, text in dev)
│   ├── idgen/
│   │   └── uuid.go                         # IDGenerator interface + UUID impl
│   ├── response/
│   │   └── json.go                         # Success(), Error(), Paginated() helpers
│   └── validator/
│       └── validator.go                    # Custom validation helpers
│
├── test/                                    # ── INTEGRATION TESTS + HELPERS ──
│   ├── integration_test.go                  # Full API flow tests (register→login→order)
│   ├── testutil/
│   │   ├── setup.go                        # Wires full app with in-memory SQLite for testing
│   │   ├── fixtures.go                     # Reusable test data builders
│   │   └── assertions.go                   # Custom test assertion helpers
│   └── mock/
│       ├── product_repo_mock.go            # Hand-written mock for ProductRepository
│       ├── order_repo_mock.go
│       ├── customer_repo_mock.go
│       ├── order_queue_mock.go
│       └── file_storage_mock.go            # Mock for FileStorage
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
| Actually talks to a database/queue/external service | `repository/sqlite/`, `sqs/`, `s3/`, etc. | `sqlite.ProductRepo.Create()`, `s3.FileStorage.Upload()` |
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
- `ImageURL` — validated URL string for product images. Optional (can be empty). Constructor validates URL format when non-empty.
- `OrderStatus` — string enum with a **state machine**: `validTransitions` map defines which transitions are legal. `CanTransitionTo()` and `TransitionTo()` enforce the machine. States: pending → validated → calculated → paid → fulfilled → completed. Any state can go to failed. Pending can go to cancelled.
- `PaymentStatus` — simple enum: pending, success, failed, refunded.
- `Pagination` — Page, Limit, Offset. Constructor validates limit 1-100.

**3. Entities** (`domain/entity/`)
- `Product` — constructor validates name, price > 0, stock >= 0. Fields include `ImageURL` (optional, value object). Methods: `ReduceStock(qty)` (returns InsufficientStock error), `RestoreStock(qty)`, `UpdatePrice()`, `UpdateName()`, `SetImageURL(url)`. All methods update `UpdatedAt`.
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
- Struct with nested configs: App (env, port), DB (path), JWT (secret, expiry), S3 (endpoint, bucket name, region), SQS (endpoint, queue URL), SFN (endpoint, ARN), AWS (region)
- Load from env vars with sensible defaults
- `.env.example` with all vars documented

**2. Database** (`database/`)
- `NewSQLiteDB(path)` — opens connection with WAL mode + foreign keys enabled
- `RunMigrations(db, dir)` — reads `.sql` files from a directory, sorts by name, executes in order
- Migration files: 001_create_products (include `image_url TEXT DEFAULT ''`), 002_create_customers, 003_create_orders, 004_create_order_items
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
- `FileStorage` — Upload(ctx, key, data, contentType) → URL, GetPresignedURL(ctx, key) → URL, Delete(ctx, key). This is the interface for S3. Key = file path like `products/{id}/image.jpg`.
- `OrderQueue` — Enqueue(order)
- `WorkflowOrchestrator` — StartOrderWorkflow, GetWorkflowStatus
- `UnitOfWork` — Begin, Commit, Rollback + access to repos within a transaction

**5. Repository Implementations**
- `repository/sqlite/` — implements all repo interfaces using `database/sql`. Map between domain entities and SQL columns. Return domain errors (e.g., `sql.ErrNoRows` → `domainErr.NewNotFound()`).
- `repository/s3/` — implements `FileStorage` using `aws-sdk-go-v2/service/s3`. Uploads files to a LocalStack S3 bucket, generates presigned URLs for reading. Uses custom endpoint for LocalStack (`http://localhost:4566`).
- `repository/memory/` — in-memory implementations for testing. `InMemoryOrderQueue` just logs — swap to SQS in Phase 4. `InMemoryFileStorage` stores bytes in a map — swap to S3 in Phase 4.

**6. Usecases** (`usecase/`)
- Each usecase is a struct with a constructor (`New*UseCase(deps...)`) and an `Execute(ctx, input) (output, error)` method.
- `product/`: CreateProduct, GetProduct, ListProducts, UpdateProduct, DeleteProduct, UploadProductImage (receives file bytes → uploads via FileStorage → saves URL on product entity → persists)
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
- `response/` — one struct per response shape: ProductResponse (includes `image_url`), OrderResponse, AuthResponse, ErrorResponse, PaginatedResponse, HealthResponse

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
- `product_handler.go` — Create, List, GetByID, Update, Delete, UploadImage (multipart form file upload)
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
orderQueue := memory.NewInMemoryOrderQueue()    // logs to stdout
fileStorage := memory.NewInMemoryFileStorage()  // stores in map
```

In Phase 4, you change TWO lines:
```
orderQueue := sqs.NewOrderQueue(sqsClient, queueURL)         // sends to real SQS
fileStorage := s3.NewFileStorage(s3Client, bucketName)        // uploads to real S3
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
| POST | `/api/v1/products/:id/image` | Yes | Upload product image (multipart) |
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

## Phase 4: SQS + S3 with LocalStack (Week 4–5)

> **Goal:** Implement `OrderQueue` with real SQS and `FileStorage` with real S3. Zero changes to usecases — just swap the adapters.

### What You'll Learn
- SQS concepts: queues, messages, visibility timeout, dead letter queues
- S3 concepts: buckets, objects, presigned URLs, content types
- `aws-sdk-go-v2` — works for both SQS and S3
- LocalStack (free, runs in Docker)
- The clean architecture payoff: swap in-memory → real AWS services with adapter changes only

### What to Build

**1. Docker Compose** (`docker-compose.yml`)
- LocalStack service (SQS, S3, DynamoDB) on port 4566
- Step Functions Local on port 8083

**2. Setup Script** (`scripts/setup-localstack.sh`)
- Create S3 bucket: `goshop-images`
- Create DLQ: `goshop-orders-dlq`
- Create main queue: `goshop-orders` with DLQ redrive policy (max 3 retries)
- Create notification queue: `goshop-notifications`

**3. S3 Adapter** (`repository/s3/file_storage.go`)
- Implements `repository.FileStorage` interface
- Uses `aws-sdk-go-v2/service/s3` with custom endpoint for LocalStack
- `Upload(ctx, key, data, contentType)` — puts object in bucket, returns the object URL
- `GetPresignedURL(ctx, key)` — generates a time-limited presigned URL for reading (e.g., 15 min expiry)
- `Delete(ctx, key)` — removes object from bucket
- Key format: `products/{product_id}/{filename}` — organized by entity
- Validate content type: only allow `image/jpeg`, `image/png`, `image/webp`
- Validate file size: max 5MB

**4. SQS Adapter** (`repository/sqs/order_queue.go`)
- Implements `repository.OrderQueue` interface
- Uses `aws-sdk-go-v2/service/sqs` with custom endpoint for LocalStack
- Serializes order to JSON, sends as SQS message

**5. Worker** (`delivery/worker/`)
- `consumer.go` — generic SQS polling loop: receive → handle → delete. Long polling (20s). Graceful shutdown.
- `order_worker.go` — deserializes message → calls `StartOrderWorkflow` usecase (Phase 5) or logs for now

**6. Worker Entry Point** (`cmd/worker/main.go`)
- Wire config → SQS client → consumer → start polling

**7. Update API main.go**
- Replace `memory.NewInMemoryOrderQueue()` with `sqs.NewOrderQueue(endpoint, queueURL)`
- Replace `memory.NewInMemoryFileStorage()` with `s3.NewFileStorage(s3Client, bucketName)`

### How Product Image Upload Works (End-to-End)

```
1. Client POSTs multipart file to /api/v1/products/:id/image
2. Handler extracts file from request, validates size + content type
3. Handler calls ProductService.UploadImage()
4. Service calls UploadProductImage usecase
5. Usecase calls FileStorage.Upload(ctx, "products/{id}/image.jpg", fileBytes, "image/jpeg")
6. S3 adapter uploads to LocalStack S3 bucket
7. Usecase saves returned URL on product entity (product.ImageURL)
8. Usecase calls ProductRepository.Update() to persist the URL
9. Response returns product with image_url field
10. Client can fetch the image via the presigned URL or direct LocalStack URL
```

### Running
```bash
docker compose up -d localstack
./scripts/setup-localstack.sh
go run cmd/api/main.go     # terminal 1
go run cmd/worker/main.go  # terminal 2

# Upload a product image
curl -X POST http://localhost:8080/api/v1/products/{id}/image \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@product-photo.jpg"

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
- [ ] Product images upload to S3 via LocalStack
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

## Testing Strategy (All Levels)

Testing follows the same layered architecture. Each layer has a different testing approach, different dependencies, and different goals.

### File Naming & Location Convention

- Test files live **next to the code they test** — not in a separate `tests/` folder
- Naming: `*_test.go` (Go convention)
- Package: use `_test` suffix for black-box testing (e.g., `package entity_test` tests `package entity` from the outside)

```
domain/entity/order.go           → domain/entity/order_test.go
domain/valueobject/money.go      → domain/valueobject/money_test.go
domain/factory/order_factory.go  → domain/factory/order_factory_test.go
usecase/product/create_product.go → usecase/product/create_product_test.go
repository/sqlite/product_repo.go → repository/sqlite/product_repo_test.go
delivery/http/handler/product_handler.go → delivery/http/handler/product_handler_test.go
```

### Test Commands

```bash
go test ./...                    # run everything
go test ./domain/... -v          # domain only (verbose)
go test ./usecase/... -v         # usecases only
go test ./delivery/... -v        # handlers only
go test ./repository/sqlite/...  # repo only
go test -run TestOrder -v        # specific test name pattern
go test -count=1 ./...           # disable test caching
go test -cover ./...             # show coverage %
go test -coverprofile=cover.out ./... && go tool cover -html=cover.out  # coverage report
```

---

### Level 1: Domain Tests (Zero mocks, zero dependencies)

**Goal:** Verify all business rules — validation, calculations, state transitions, factory construction.

**Why zero mocks?** The domain has no external dependencies. Entities, value objects, and factories only depend on each other and Go stdlib. This is the biggest benefit of clean architecture — your most important logic is the easiest to test.

**What to test:**

**Value Objects:**
- `Money` — creation with valid/invalid amounts, Add with same/different currencies, Multiply with positive/negative/zero qty, Percentage edge cases (0%, 100%), IsZero, String formatting
- `Email` — valid emails, invalid formats, empty string, trimming, lowercasing
- `Phone` — valid Indonesian formats (08xx, +628xx, 628xx), invalid formats, stripping spaces/dashes
- `OrderStatus` — every valid transition, every invalid transition, terminal state checks, CanTransitionTo vs TransitionTo
- `Pagination` — valid page/limit, defaults for invalid values, limit bounds (1-100)

**Entities:**
- `Product` — constructor rejects empty name / zero price / negative stock. ReduceStock with valid/insufficient/zero/negative qty. RestoreStock. UpdatePrice with zero.
- `Order` — constructor rejects empty customer ID / empty items. CalculateTotal with single item, multiple items, verifies subtotal + PPN 11% tax + total. TransitionTo valid/invalid paths.
- `OrderItem` — constructor rejects empty product ID / zero qty / zero price. LineTotal calculation.
- `Customer` — constructor rejects empty name / invalid email.

**Factories:**
- Create a stub `IDGenerator` that returns predictable IDs ("id-1", "id-2", etc.) and a stub `PasswordHasher` that returns `"hashed_" + password`
- `OrderFactory.Create` — verify it assigns IDs, calculates totals, attaches OrderCreated event
- `CustomerFactory.Create` — verify it validates, hashes password, assigns ID

**Testing pattern: table-driven tests**
```
Use Go's table-driven test pattern for all value object and entity tests:
  tests := []struct{ name string; input; wantOutput; wantErr bool }
  for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { ... }) }
```

This gives you exhaustive coverage with minimal code repetition. Aim for 5-10 test cases per function covering: happy path, boundary values, invalid inputs, edge cases.

---

### Level 2: Usecase Tests (Mock repositories)

**Goal:** Verify application logic — the usecase calls the right repos in the right order with the right data.

**How to mock:** Create mock implementations of repository interfaces. Two approaches:

**Option A: Hand-written mocks** (simpler, recommended for learning)
- Create `usecase/product/mock_test.go` with a `MockProductRepository` struct
- Each method stores what it received and returns pre-configured results
- Example: `Create()` stores the product it received; `FindByID()` returns a product from a map

**Option B: Use `repository/memory/` adapters as test doubles**
- The in-memory implementations you built for local dev also work as test doubles
- Advantage: they actually behave like a real repo (store/retrieve data)

**What to test per usecase:**

**CreateProduct:**
- Happy path: returns created product with generated ID
- Validation error: propagates domain error when invalid input
- Repo error: propagates when repo.Create fails

**CreateOrder:**
- Happy path: resolves product prices, reduces stock, creates order via factory, persists, enqueues
- Product not found: returns NotFound error
- Insufficient stock: returns InsufficientStock error
- Queue failure: order still created (queue failure is non-fatal)

**Register:**
- Happy path: creates customer, hashes password, persists
- Duplicate email: returns Conflict error

**Login:**
- Happy path: returns token
- Wrong email: returns Unauthorized
- Wrong password: returns Unauthorized

**Verify these behaviors:**
- Correct repository methods called with correct arguments
- Domain validation errors propagated (not swallowed)
- Side effects happened (stock reduced, event attached, etc.)

---

### Level 3: Repository Tests (Real SQLite, in-memory DB)

**Goal:** Verify SQL queries work correctly — inserts, selects, joins, edge cases.

**How:** Use a real SQLite database with `":memory:"` path. This gives you a fresh DB per test with zero cleanup needed.

**Setup pattern:**
```
Each test (or TestMain):
  1. Open SQLite with ":memory:"
  2. Run migrations
  3. Seed test data if needed
  4. Run test
  5. DB disappears automatically (in-memory)
```

**What to test:**

**ProductRepo:**
- Create + FindByID roundtrip — verify all fields survive persistence (especially Money → price_amount/price_currency → Money)
- FindAll with pagination — verify limit, offset, total count, ordering
- Update — verify fields changed, updated_at changed
- Delete — verify gone; delete non-existent returns NotFound
- FindByID not found — returns domain NotFound error (not sql.ErrNoRows)

**OrderRepo:**
- Create order with items — verify order + order_items both persisted
- FindByID — verify items loaded with correct prices
- FindByCustomerID — verify only that customer's orders returned
- UpdateStatus — verify status changed

**CustomerRepo:**
- Create + FindByEmail roundtrip
- ExistsByEmail — true when exists, false when not
- Duplicate email — returns error (UNIQUE constraint)

**Key thing to verify:** The repo correctly maps between domain entities (Money value object) and SQL columns (price_amount INTEGER + price_currency TEXT). This mapping is where bugs hide.

---

### Level 4: Handler Tests (httptest + Gin)

**Goal:** Verify HTTP contract — request parsing, response shape, status codes, error mapping.

**How:** Use Go's `net/http/httptest` package with a real Gin engine. Mock the service layer.

**Setup pattern:**
```
Each test:
  1. Create mock service (or use a real service with mock repos)
  2. Create handler with mock service
  3. Create Gin engine + register route
  4. Create httptest.NewRecorder + http.NewRequest
  5. Serve the request
  6. Assert: status code, response body (JSON), headers
```

**What to test:**

**Product Handler:**
- POST /products with valid body → 201 + product JSON
- POST /products with missing name → 400 + error message
- POST /products with negative price → 400 + error message
- GET /products → 200 + paginated list
- GET /products/:id with valid ID → 200 + product
- GET /products/:id with unknown ID → 404 + not found error
- DELETE /products/:id → 200 or 204

**Order Handler:**
- POST /orders with valid items → 201 + order JSON with calculated total
- POST /orders with empty items → 400
- POST /orders with insufficient stock → 422

**Auth Handler:**
- POST /register with valid body → 201
- POST /register with duplicate email → 409
- POST /register with invalid email → 400
- POST /login with valid credentials → 200 + token
- POST /login with wrong password → 401

**Middleware Tests:**
- Auth middleware: valid token → sets customer_id in context + continues
- Auth middleware: missing header → 401 + aborts
- Auth middleware: invalid token → 401 + aborts
- Error middleware: DomainError(NOT_FOUND) → 404
- Error middleware: DomainError(VALIDATION) → 400
- Error middleware: unknown error → 500

**Key things to verify:**
- Response JSON structure matches your DTO response types exactly
- Error responses include the `error` field with a useful message
- Paginated responses include `meta` with page, limit, total, total_pages
- Auth-protected routes reject requests without a token

---

### Level 5: Integration Tests (Full API flow)

**Goal:** Verify the entire stack works together — HTTP → handler → service → usecase → repo → DB and back.

**How:** Start a real server (or use Gin's test mode), use a real in-memory SQLite DB, and make actual HTTP requests.

**Setup pattern:**
```
TestMain or each test:
  1. Load config (test overrides: in-memory DB, test JWT secret)
  2. Wire everything (same as cmd/api/main.go but with test config)
  3. Create httptest.Server with the Gin engine
  4. Run tests against the server URL
  5. Server shuts down after test
```

**What to test (end-to-end flows):**

**Flow 1: Auth → Product CRUD**
1. Register a customer → 201
2. Login → 200 + token
3. Create product with token → 201
4. List products → 200 + contains created product
5. Get product by ID → 200
6. Update product → 200
7. Delete product → 200 or 204
8. Get deleted product → 404

**Flow 2: Auth → Order Flow**
1. Register + Login
2. Create 2 products (with stock)
3. Create order with both products → 201 + verify total includes PPN 11%
4. Get order → verify status is "pending", items match, total correct
5. Create order that exceeds stock → 422 InsufficientStock
6. Verify product stock was reduced after successful order

**Flow 3: Error Handling**
1. Access protected endpoint without token → 401
2. Access with expired/invalid token → 401
3. Create product with missing fields → 400 with field errors
4. Register with duplicate email → 409

**Key things to verify:**
- Data actually persists across requests (create → get returns same data)
- Stock is correctly reduced after order creation
- PPN 11% tax calculation is correct end-to-end
- Auth flow works: register → login → use token → access protected resources
- Error codes match what the API contract promises

---

### Testing Pyramid Summary

```
        /  Integration  \        ← Few tests, slow, high confidence
       /    Handlers      \      ← HTTP contract, status codes, JSON shape
      /     Usecases        \    ← Business flow, mock repos
     /    Repositories        \  ← SQL correctness, in-memory SQLite
    /       Domain              \ ← Most tests, fastest, zero deps
```

| Level | # of Tests | Speed | Mocks Needed | What Breaks If It Fails |
|-------|-----------|-------|--------------|------------------------|
| Domain | Many (50+) | <1s | None | Business rules wrong |
| Repository | Medium (20+) | <2s | None (real in-memory DB) | Data persistence broken |
| Usecase | Medium (20+) | <1s | Mock repos | App logic flow wrong |
| Handler | Medium (15+) | <2s | Mock services | API contract broken |
| Integration | Few (5-10) | <5s | None (real stack) | Stack doesn't work together |

### Minimum Viable Test Coverage

Before moving to the next phase, ensure:

**Phase 1 (Domain):** All entity constructors, all value object creation/methods, all factory methods, all status transitions — aim for 90%+ coverage on `domain/`

**Phase 2 (API):** At least happy path + one error case per usecase. At least happy path + one error case per handler endpoint. One full integration flow (register → login → create product → create order).

**Phase 3+ (AWS):** Lambda handlers tested with sample event JSON. Worker message processing tested with mock queue.

---

## Bonus Challenges

1. **DynamoDB (local)** — `repository.AuditLogger` interface for order processing logs
2. **gRPC delivery** — `delivery/grpc/` calling the same services
3. **Outbox pattern** — ensure order persistence + SQS message are atomic
4. **CQRS** — separate read and write repository interfaces
5. **OpenTelemetry** — distributed tracing across HTTP → SQS → Lambda → Step Functions
6. **Image resizing Lambda** — trigger a Lambda on S3 upload that creates thumbnail versions

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