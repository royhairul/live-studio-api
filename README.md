# 🎬 Live Studio Management API

A comprehensive backend management system built in **Go (Golang)** designed to manage live streaming studio operations, specifically tailored for managing live streaming activities on e-commerce platforms like Shopee.

## 📋 Table of Contents

- [Overview](#-overview)
- [Tech Stack](#️-tech-stack)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
- [API Endpoints](#-api-endpoints)
- [Domain Modules](#-domain-modules)
- [Project Structure](#-project-structure)
- [Environment Variables](#-environment-variables)
- [License](#-license)

---

## 🎯 Overview

**Live Studio API** provides a complete solution for managing live streaming studio operations including:

- 👥 **User & Role Management** - Multi-tenant user authentication with role-based access control
- 🏢 **Studio Management** - Manage multiple physical streaming studios
- 🎙️ **Host Management** - Register and manage live streaming talents
- 📅 **Scheduling** - Schedule hosts for shifts and track attendance
- 📺 **Live Session Tracking** - Track live sessions with comprehensive metrics
- 💰 **Finance & Transactions** - Monitor orders, transactions, and commissions
- 📊 **Performance Analytics** - Track performance by host, account, and studio
- 🎯 **Target Management** - Set and monitor sales/performance targets

---

## ⚙️ Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.23+ |
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **Dependency Injection** | [Uber FX](https://github.com/uber-go/fx) |
| **Database** | PostgreSQL with [GORM](https://gorm.io/) ORM |
| **Authentication** | JWT (JSON Web Tokens) |
| **External Integration** | Shopee API (Creator, Seller, Affiliate) |
| **ID Generation** | Snowflake IDs + UUID |
| **API Documentation** | Swagger UI / OpenAPI |
| **Email Service** | SMTP via [Gomail](https://github.com/go-gomail/gomail) |

---

## 🏗️ Architecture

The project follows a **clean/domain-driven architecture** with modular organization using Uber FX for dependency injection.

```
live-studio-api/
├── cmd/                    # Application commands (API, migrate, seed)
│   ├── api/                # Main API server entry
│   ├── migrate/            # Database migration commands
│   └── seed/               # Database seeding commands
├── config/                 # Configuration management
├── database/               # Database initialization & migrations
│   ├── migrations.go       # Auto-migrations setup
│   └── seeders/            # Data seeders
├── docs/                   # OpenAPI/Swagger documentation
├── internal/               # Core application logic
│   ├── aggregator/         # Data aggregation services
│   ├── clients/            # External API clients (Shopee)
│   ├── domains/            # Domain-specific modules
│   ├── jobs/               # Background jobs/scheduled tasks
│   ├── middleware/         # HTTP middleware
│   ├── pkg/                # Shared utilities
│   └── server/             # Application bootstrap
├── routes/                 # API route definitions
├── scripts/                # Utility scripts
└── logs/                   # Application logs
```

### Domain Module Structure

Each domain follows a consistent structure:

```
domain/
├── controller/             # HTTP request handlers
│   ├── controller.go       # Controller interface
│   └── controller_impl.go  # Controller implementation
├── entity/                 # Database models (GORM)
│   └── entity.go
├── params/                 # Request/Response DTOs
│   ├── request.go
│   └── response.go
├── repository/             # Data access layer
│   ├── repository.go       # Repository interface
│   └── repository_impl.go  # Repository implementation
├── service/                # Business logic
│   ├── service.go          # Service interface
│   └── service_impl.go     # Service implementation
├── route.go                # Route registration
└── module.go               # FX module definition
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 14+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/royhairul/live-studio-api.git
   cd live-studio-api
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Configure environment variables**
   Edit `.env` file with your database credentials and other configurations.

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Create PostgreSQL database**
   ```sql
   CREATE DATABASE db_livestudio;
   ```

6. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

### API Documentation

Once the server is running, access the Swagger UI documentation at:
- **Swagger UI**: `http://localhost:8080/docs`
- **OpenAPI JSON**: `http://localhost:8080/docs/openapi.json`

---

## 📌 API Endpoints

All API endpoints are prefixed with `/api`.

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/healthcheck` | Check API health status |

### 🔐 Authentication (`/api/auth`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/api/auth/login` | User login | ❌ |
| `POST` | `/api/auth/register` | User registration | ❌ |
| `POST` | `/api/auth/forgot-password` | Request password reset | ❌ |
| `POST` | `/api/auth/verify-otp` | Verify OTP code | ❌ |
| `POST` | `/api/auth/reset-password` | Reset password | ❌ |
| `GET` | `/api/auth/me` | Get current user info | ✅ (superadmin, admin, host) |

### 👥 Users (`/api/users`)

> **Requires**: Tenant middleware

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/users` | Get all users | ✅ |
| `POST` | `/api/users` | Create new user | ✅ |
| `GET` | `/api/users/:id` | Get user by ID | ✅ |
| `PUT` | `/api/users/:id` | Update user | ✅ |
| `DELETE` | `/api/users/:id` | Delete user | ✅ |

#### Superadmin Routes (`/api/superadmin/user`)

> **Requires**: `superadmin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/superadmin/user` | Get all users |
| `POST` | `/api/superadmin/user` | Create new user |
| `GET` | `/api/superadmin/user/:id` | Get user by ID |
| `PUT` | `/api/superadmin/user/:id` | Update user |
| `DELETE` | `/api/superadmin/user/:id` | Delete user |

### 🏢 Studios (`/api/studio`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/studio` | Get all studios |
| `POST` | `/api/studio` | Create new studio |
| `GET` | `/api/studio/:id` | Get studio by ID |
| `PUT` | `/api/studio/:id` | Update studio |
| `DELETE` | `/api/studio/:id` | Delete studio |

### 🎙️ Hosts (`/api/host`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/host` | Get all hosts |
| `POST` | `/api/host` | Create new host |
| `GET` | `/api/host/group-by-studio` | Get hosts grouped by studio |
| `GET` | `/api/host/:id` | Get host by ID |
| `PUT` | `/api/host/:id` | Update host |
| `DELETE` | `/api/host/:id` | Delete host |

### ⏰ Shifts (`/api/shift`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/shift` | Get all shifts |
| `POST` | `/api/shift` | Create new shift |
| `GET` | `/api/shift/:id` | Get shift by ID |
| `PUT` | `/api/shift/:id` | Update shift |
| `DELETE` | `/api/shift/:id` | Delete shift |

### 📅 Schedules (`/api/schedule`)

> **Requires**: `superadmin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/schedule` | Get all schedules |
| `POST` | `/api/schedule` | Create new schedule |
| `GET` | `/api/schedule/scheduled` | Get schedule by shift and date |
| `GET` | `/api/schedule/:id` | Get schedule by ID |
| `PUT` | `/api/schedule/:id` | Update schedule |
| `DELETE` | `/api/schedule/:id` | Delete schedule |

### ✅ Attendance (`/api/attendance`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/attendance` | Get all attendance records |
| `GET` | `/api/attendance/unchecked-out` | Get unchecked-out records |
| `POST` | `/api/attendance/check-in` | Check in |
| `POST` | `/api/attendance/check-out` | Check out |

### 🔑 Accounts (`/api/account`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/account` | Get all accounts |
| `POST` | `/api/account` | Create or update account |
| `GET` | `/api/account/:id` | Get account by ID |
| `PUT` | `/api/account/:id` | Update account |
| `DELETE` | `/api/account/:id` | Delete account |

### 🔐 Account Sessions (`/api/account-session`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/account-session` | Get all sessions |
| `POST` | `/api/account-session` | Create new session |
| `GET` | `/api/account-session/:id` | Get session by ID |
| `PUT` | `/api/account-session/:id` | Update session |
| `DELETE` | `/api/account-session/:id` | Delete session |

### 📺 Account Ads (`/api/accountads`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/accountads` | Get all ad accounts |
| `POST` | `/api/accountads` | Create new ad account |
| `GET` | `/api/accountads/:id` | Get ad account by ID |
| `PUT` | `/api/accountads/:id` | Update ad account |
| `DELETE` | `/api/accountads/:id` | Delete ad account |

### 📺 Live Sessions (`/api/live`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/live/shopee` | Get Shopee live sessions | ✅ |
| `GET` | `/api/live/shopee/:id/:sessionId` | Get live session details | ✅ |

### 📊 Performance (`/api/performa`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/performa/host` | Get all hosts performance |
| `GET` | `/api/performa/host/:id` | Get host performance by ID |
| `GET` | `/api/performa/account` | Get all accounts performance |
| `GET` | `/api/performa/studio` | Get all studios performance |
| `GET` | `/api/performa/studio/:id` | Get studio performance by ID |

### 💰 Finance (`/api/finance`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/finance` | Get all financial records |

### 📦 Products (`/api/product`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/product` | Get all products |
| `POST` | `/api/product` | Create new product |
| `GET` | `/api/product/:id` | Get product by ID |
| `PUT` | `/api/product/:id` | Update product |
| `DELETE` | `/api/product/:id` | Delete product |

### 🛒 Orders (`/api/order`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/order` | Get all orders |
| `POST` | `/api/order` | Create new order |
| `GET` | `/api/order/:id` | Get order by ID |
| `PUT` | `/api/order/:id` | Update order |
| `DELETE` | `/api/order/:id` | Delete order |

### 💳 Transactions (`/api/transaction`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/transaction` | Get all transactions |
| `POST` | `/api/transaction` | Create new transaction |
| `GET` | `/api/transaction/grouped` | Get transactions grouped |
| `GET` | `/api/transaction/:id` | Get transaction by ID |
| `PUT` | `/api/transaction/:id` | Update transaction |
| `DELETE` | `/api/transaction/:id` | Delete transaction |

### 🎯 Targets (`/api/target`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/target` | Get all targets |
| `POST` | `/api/target` | Create new target |
| `GET` | `/api/target/:id` | Get target by ID |
| `PUT` | `/api/target/:id` | Update target |
| `DELETE` | `/api/target/:id` | Delete target |

### 📊 Dashboard (`/api/dashboard`)

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/dashboard` | Get dashboard data |

### 🛡️ Roles (`/api/role`)

> **Requires**: `superadmin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/role` | Get all roles |
| `POST` | `/api/role` | Create new role |
| `GET` | `/api/role/:id` | Get role by ID |
| `PUT` | `/api/role/:id` | Update role |
| `DELETE` | `/api/role/:id` | Delete role |

### 🔐 Permissions (`/api/permission`)

> **Requires**: `superadmin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/permission` | Get all permissions |
| `GET` | `/api/permission/grouped` | Get permissions grouped |
| `POST` | `/api/permission` | Create new permission |
| `GET` | `/api/permission/:id` | Get permission by ID |
| `PUT` | `/api/permission/:id` | Update permission |
| `DELETE` | `/api/permission/:id` | Delete permission |

---

## 📦 Domain Modules

| Module | Description | Documentation |
|--------|-------------|---------------|
| **auth** | User authentication & authorization | [View Details](./internal/domains/auth/README.md) |
| **user** | User management | [View Details](./internal/domains/user/README.md) |
| **host** | Live streaming hosts (talents) | [View Details](./internal/domains/host/README.md) |
| **studio** | Physical studios management | [View Details](./internal/domains/studio/README.md) |
| **shift** | Work shifts configuration | [View Details](./internal/domains/shift/README.md) |
| **schedule** | Host scheduling | [View Details](./internal/domains/schedule/README.md) |
| **attendance** | Host attendance tracking | [View Details](./internal/domains/attendance/README.md) |
| **account** | Platform accounts (Shopee) | [View Details](./internal/domains/account/README.md) |
| **accountsession** | Account session management | [View Details](./internal/domains/accountsession/README.md) |
| **accountads** | Advertising accounts | [View Details](./internal/domains/accountads/README.md) |
| **live** | Live streaming sessions | [View Details](./internal/domains/live/README.md) |
| **performa** | Performance analytics | [View Details](./internal/domains/performa/README.md) |
| **finance** | Financial management | [View Details](./internal/domains/finance/README.md) |
| **product** | Product catalog | [View Details](./internal/domains/product/README.md) |
| **order** | Order management | [View Details](./internal/domains/order/README.md) |
| **transaction** | Transaction tracking | [View Details](./internal/domains/transaction/README.md) |
| **target** | Sales/performance targets | [View Details](./internal/domains/target/README.md) |
| **dashboard** | Dashboard aggregation | [View Details](./internal/domains/dashboard/README.md) |
| **role** | Role-based access control | [View Details](./internal/domains/role/README.md) |
| **permission** | Permission management | [View Details](./internal/domains/permission/README.md) |

---

## 🔧 Internal Packages

| Package | Description | Documentation |
|---------|-------------|---------------|
| **aggregator** | Data aggregation services (performance metrics) | [View Details](./internal/aggregator/README.md) |
| **clients** | External API clients (Shopee integration) | [View Details](./internal/clients/README.md) |
| **pkg** | Shared utilities and packages | [View Details](./internal/pkg/README.md) |

---

## 🔧 Environment Variables

### Application Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `LiveStudio` |
| `APP_ENV` | Environment (development/production/staging/maintenance) | `development` |
| `APP_DEBUG` | Enable debug logs | `true` |
| `APP_MAINTENANCE` | Enable maintenance mode | `false` |
| `GO_MODULE` | Go module path | `github.com/royhairul/live-studio-api` |

### Server Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host | `localhost` |
| `SERVER_PORT` | Server port | `8080` |
| `API_URL` | Full API URL for OpenAPI docs (optional) | Auto-generated from host:port |
| `ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | `http://localhost:5173,...` |

### Database Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | - |
| `DB_NAME` | PostgreSQL database name | `db_livestudio` |
| `DB_PORT` | PostgreSQL port | `5432` |

### JWT Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | JWT signing secret (use a strong random string) | - |
| `JWT_EXPIRED_AT` | JWT expiration time (e.g., `24h`, `7d`) | `24h` |

### Superadmin Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SUPERADMIN_NAME` | Default superadmin display name | - |
| `SUPERADMIN_EMAIL` | Default superadmin email address | - |
| `SUPERADMIN_PASS` | Default superadmin password | - |

### 📧 Email/SMTP Configuration

The application uses **SMTP** for sending emails (password reset OTP). It uses the `gomail.v2` library.

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `EMAIL_FROM` | Sender email address | - | `noreply@example.com` |
| `EMAIL_PASSWORD` | SMTP password or App Password | - | `xxxx xxxx xxxx xxxx` |
| `EMAIL_HOST` | SMTP server host | `smtp.gmail.com` | `smtp.gmail.com` |
| `EMAIL_PORT` | SMTP server port | `587` | `587` (TLS) |

#### Gmail SMTP Setup

To use Gmail as your SMTP provider:

1. **Enable 2-Factor Authentication** on your Google account
2. **Generate an App Password**:
   - Go to [Google Account Security](https://myaccount.google.com/security)
   - Navigate to "2-Step Verification" → "App passwords"
   - Select "Mail" and "Other (Custom name)"
   - Copy the generated 16-character password
3. **Configure environment variables**:
   ```env
   EMAIL_FROM=your-email@gmail.com
   EMAIL_PASSWORD=xxxx xxxx xxxx xxxx
   EMAIL_HOST=smtp.gmail.com
   EMAIL_PORT=587
   ```

#### Other SMTP Providers

| Provider | Host | Port (TLS) | Port (SSL) |
|----------|------|------------|------------|
| Gmail | `smtp.gmail.com` | `587` | `465` |
| Outlook/Hotmail | `smtp-mail.outlook.com` | `587` | - |
| Yahoo | `smtp.mail.yahoo.com` | `587` | `465` |
| SendGrid | `smtp.sendgrid.net` | `587` | `465` |

> ⚠️ **Note**: For production, consider using a dedicated email service like SendGrid or AWS SES for better deliverability and monitoring.

---

### 🔒 API Documentation Security

API documentation access is controlled by `APP_ENV`:

| APP_ENV | Docs Access | Authentication |
|---------|-------------|----------------|
| `development` | ✅ Open | None |
| `staging` | ✅ Protected | Basic Auth (uses `SUPERADMIN_NAME` / `SUPERADMIN_PASS`) |
| `production` | ❌ Disabled | Returns 404 |
| `maintenance` | ❌ Disabled | Returns 404 |

#### Examples

**Development** (docs open):
```env
APP_ENV=development
```

**Staging** (docs with login):
```env
APP_ENV=staging
SUPERADMIN_NAME=your-name
SUPERADMIN_PASS=your-password
```
> Browser will prompt for username/password when accessing `/docs`

**Production** (docs hidden):
```env
APP_ENV=production
```

---

## 📜 License

MIT License © 2025

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📞 Support

For support, please open an issue in the GitHub repository.
