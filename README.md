<h1 align="center">
  <a href="https://gofiber.io">
    <picture>
      <source height="125" media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo-dark.svg">
      <img height="125" alt="Fiber" src="https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo.svg">
    </picture>
  </a>
  <br>
  <a href="https://pkg.go.dev/github.com/gofiber/fiber/v3#pkg-overview">
    <img src="https://img.shields.io/badge/%F0%9F%93%9A%20godoc-pkg-00ACD7.svg?color=00ACD7&style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/gofiber/fiber/v3">
    <img src="https://img.shields.io/badge/%F0%9F%93%9D%20goreport-A%2B-75C46B?style=flat-square">
  </a>
  <a href="https://codecov.io/gh/gofiber/fiber" >
   <img alt="Codecov" src="https://img.shields.io/codecov/c/github/gofiber/fiber?token=3Cr92CwaPQ&style=flat-square&logo=codecov&label=codecov">
 </a>
  <a href="https://github.com/gofiber/fiber/actions?query=workflow%3ATest">
    <img src="https://img.shields.io/github/actions/workflow/status/gofiber/fiber/test.yml?branch=main&label=%F0%9F%A7%AA%20tests&style=flat-square&color=75C46B">
  </a>
    <a href="https://docs.gofiber.io">
    <img src="https://img.shields.io/badge/%F0%9F%92%A1%20fiber-docs-00ACD7.svg?style=flat-square">
  </a>
  <a href="https://gofiber.io/discord">
    <img src="https://img.shields.io/discord/704680098577514527?style=flat-square&label=%F0%9F%92%AC%20discord&color=00ACD7">
  </a>
</h1>
<p align="center">
  <em><b>Fiber</b> is an <a href="https://github.com/expressjs/express">Express</a> inspired <b>web framework</b> built on top of <a href="https://github.com/valyala/fasthttp">Fasthttp</a>, the <b>fastest</b> HTTP engine for <a href="https://go.dev/doc/">Go</a>. Designed to <b>ease</b> things up for <b>fast</b> development with <a href="https://docs.gofiber.io/#zero-allocation"><b>zero memory allocation</b></a> and <b>performance</b> in mind.</em>
</p>>

# Fiber REST API Boilerplate

## Project Description
This repository provides a production-ready backend boilerplate for building RESTful APIs and microservices in Go. It combines Fiber for fast HTTP handling, Uber Fx for dependency injection, and GORM for database access to deliver a modular architecture that scales from small services to larger SaaS backends.

**Problems it solves**
- Eliminates repetitive setup for configuration, dependency wiring, and routing.
- Provides a consistent module pattern for features such as auth and user management.
- Bakes in infrastructure wiring (database, Redis, email, token handling) so teams can focus on business logic.

**Target use cases**
- Public APIs and internal microservices
- SaaS backends with authentication/authorization needs
- Rapid prototyping for Go/Fiber services

## Key Features
- **Modular architecture with Fx** for clear dependency boundaries and lifecycle management.
- **Fiber v2** for high-performance HTTP routing and middleware.
- **GORM + PostgreSQL** with migration and seeding entry points.
- **Redis cache + throttling** utilities via Fiber storage adapters.
- **JWT-based auth flows** and role-based access control scaffolding.
- **Centralized configuration** with Viper and environment validation.
- **Structured logging** using Logrus.
- **Security & scalability considerations**
  - JWT access/refresh secrets and token verification services.
  - Middleware layering for access control and request metadata.
  - Stateless services with Redis-backed cache and blacklist checks.

## Project Structure
```
.
├── cmd
│   ├── api              # API entrypoint
│   ├── migrate          # DB migration entrypoint
│   └── seed             # DB seeding entrypoint
├── config               # Environment configuration and module wiring
├── internal             # Application modules and infrastructure
│   ├── auth             # Auth domain (routes, handlers, services)
│   ├── user             # User domain (entities, repository, services)
│   ├── http             # HTTP routing contracts, middleware, error handling
│   ├── infra            # Infrastructure services (db, redis, email, token, logger)
│   ├── identity         # Current user and claims helpers
│   ├── enums            # Typed enums for roles, tokens, access levels
│   └── apperr           # Centralized error types
├── pkg                  # Shared helpers (env, pagination, http responses, crypto)
├── .env.example         # Sample environment variables
├── Dockerfile           # Container build
├── docker-compose.yml   # Local dev dependencies
├── Makefile             # Local dev commands
└── go.mod               # Go module definition
```

### Folder responsibilities
- **cmd/**: Entry points that bootstrap the application (API server, migrations, seeding). Each command sets up configuration and executes its workflow.
- **config/**: Environment loading and validation. Defines the `Env` struct and Fx module setup for config binding.
- **internal/**: Core application logic. Feature modules live here and are wired together using Fx.
  - **auth/** and **user/**: Example domain modules demonstrating handlers, services, DTOs, repositories, and routes.
  - **http/**: Router contracts, middleware (JWT, role checks, request metadata), and error handling.
  - **infra/**: Infrastructure services (database, Redis cache, token service, validators, logging).
- **pkg/**: Reusable utility packages shared across modules (pagination, password hashing, HTTP response helpers, slice utilities).

## Getting Started
### Prerequisites
- Go 1.24+
- PostgreSQL (or a compatible database)
- Redis (optional but required for cache/blacklist features)
- Make (optional but convenient)

### Environment configuration
1. Copy the example environment file:
   ```bash
   cp .env.example .env.local
   ```
2. Fill in database, Redis, JWT, and SMTP credentials as needed.
3. The configuration loader reads `.env.local` when `APP_ENV=development`, otherwise it falls back to `.env`.

### Run locally
```bash
make run
```
Or, without Make:
```bash
export $(cat .env.local | xargs)
go run ./cmd/api
```

### Migrations and seeding
```bash
make migrate
make seed
```

## Best Practices
- **Follow the module pattern**: Each feature module should expose its own `Module` that provides handlers, services, repositories, and routes.
- **Register new modules in `internal/feature_modules.go`** to keep the Fx wiring centralized.
- **Keep DTOs and entities separate**: Use DTOs for request/response shapes and entities for persistence models.
- **Lean middleware**: Put cross-cutting concerns in `internal/http/middleware` and keep handlers thin.
- **Prefer shared helpers in `pkg/`**: Pagination, hashing, and response formatting should be reused rather than duplicated.
- **Environment-driven configuration**: Never hardcode secrets; always use the `Env` struct and `.env` files.

## Extending the Boilerplate
1. Create a new module folder under `internal/` (e.g., `internal/order`).
2. Add domain layers (entity, repository, service, handler, route).
3. Implement a `Module` that provides dependencies and registers routes.
4. Add the module to `internal/feature_modules.go`.

---
This boilerplate is intentionally opinionated to encourage clean separation of concerns and scalable growth for production services.
