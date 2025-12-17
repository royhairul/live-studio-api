# 🔐 Account Session Domain

The **Account Session** domain manages platform account sessions for the Live Studio API.

## 📋 Overview

This module provides:
- Account session CRUD operations
- Session token management
- Multi-tenant session isolation

## 📁 Structure

```
accountsession/
├── controller/
│   ├── accountsession_controller.go       # Controller interface
│   └── accountsession_controller_impl.go  # Controller implementation
├── entity/
│   └── accountsession_entity.go           # Session database model
├── params/
│   ├── request.go                         # Request DTOs
│   └── response.go                        # Response DTOs
├── repository/
│   ├── accountsession_repository.go       # Repository interface
│   └── accountsession_repository_impl.go  # Repository implementation
├── service/
│   ├── accountsession_service.go          # Service interface
│   └── accountsession_service_impl.go     # Service implementation
├── route.go                               # Route definitions
├── module.go                              # FX module
└── README.md                              # This file
```

## 🌐 API Endpoints

Base path: `/api/account-session`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all sessions |
| `POST` | `/` | Create new session |
| `GET` | `/:id` | Get session by ID |
| `PUT` | `/:id` | Update session |
| `DELETE` | `/:id` | Delete session |

## 📝 Request/Response Examples

### Create Session

**Request:**
```json
POST /api/account-session
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "account_id": 1,
  "session_token": "token_value",
  "expires_at": "2025-02-01T00:00:00Z"
}
```

**Response:**
```json
{
  "message": "Session created successfully",
  "data": {
    "id": 1,
    "account_id": 1,
    "session_token": "token_value",
    "expires_at": "2025-02-01T00:00:00Z",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Sessions

**Request:**
```http
GET /api/account-session
Authorization: Bearer <jwt_token>
```

## 🔗 Dependencies

- Account domain for account reference
- Tenant middleware for multi-tenancy

## 📌 Notes

- Sessions store authentication tokens for platform API access
- Sessions have expiration times and should be refreshed
- Used by external API clients (Shopee) for authenticated requests
