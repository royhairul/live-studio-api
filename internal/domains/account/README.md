# 🔑 Account Domain

The **Account** domain manages platform accounts (e.g., Shopee accounts) for the Live Studio API.

## 📋 Overview

This module provides:
- Account CRUD operations
- Platform account management (Shopee, etc.)
- Session/cookie storage for API access
- Multi-tenant account isolation

## 📁 Structure

```
account/
├── controller/
│   ├── account_controller.go       # Controller interface
│   └── account_controller_impl.go  # Controller implementation
├── entity/
│   └── account_entity.go           # Account database model
├── params/
│   ├── request.go                  # Request DTOs
│   └── response.go                 # Response DTOs
├── repository/
│   ├── account_repository.go       # Repository interface
│   └── account_repository_impl.go  # Repository implementation
├── service/
│   ├── account_service.go          # Service interface
│   └── account_service_impl.go     # Service implementation
├── route.go                        # Route definitions
├── module.go                       # FX module
└── README.md                       # This file
```

## 🗃️ Entity

```go
type Account struct {
    gorm.Model               // ID, CreatedAt, UpdatedAt, DeletedAt
    Name      string         // Account display name (max 100 chars)
    UniqueID  string         // Platform unique ID (max 20 chars)
    Username  string         // Platform username (max 100 chars)
    Password  string         // Account password (optional)
    Email     string         // Account email (max 100 chars)
    Platform  string         // Platform name: "shopee", etc.
    Cookie    string         // Session cookie for API access
    Device    string         // Device information
    StudioID  uint16         // Foreign key to Studio
    Studio    Studio         // Studio relationship
    TenantBase               // Multi-tenant support
}
```

## 🌐 API Endpoints

Base path: `/api/account`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all accounts |
| `POST` | `/` | Create or update account |
| `GET` | `/:id` | Get account by ID |
| `PUT` | `/:id` | Update account |
| `DELETE` | `/:id` | Delete account (soft delete) |

## 📝 Request/Response Examples

### Create Account

**Request:**
```json
POST /api/account
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Shopee Store A",
  "unique_id": "12345678",
  "username": "shopee_user",
  "email": "store@example.com",
  "platform": "shopee",
  "cookie": "session_cookie_value",
  "studio_id": 1
}
```

**Response:**
```json
{
  "message": "Account created successfully",
  "data": {
    "id": 1,
    "name": "Shopee Store A",
    "unique_id": "12345678",
    "username": "shopee_user",
    "email": "store@example.com",
    "platform": "shopee",
    "studio": {
      "id": 1,
      "name": "Studio A"
    },
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Accounts

**Request:**
```http
GET /api/account
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Shopee Store A",
      "unique_id": "12345678",
      "username": "shopee_user",
      "platform": "shopee",
      "studio": {
        "id": 1,
        "name": "Studio A"
      }
    }
  ]
}
```

### Update Account

**Request:**
```json
PUT /api/account/1
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Shopee Store A - Updated",
  "cookie": "new_session_cookie"
}
```

## 🔗 Dependencies

- Studio domain for studio assignment
- Account Session domain for session management
- Transaction domain references accounts
- Tenant middleware for multi-tenancy

## 📌 Notes

- Accounts store platform credentials and session data
- Cookie field stores session tokens for API authentication
- Sensitive data (passwords, cookies) should be handled securely
- Account credentials are used to interact with platform APIs (Shopee)
