# 👥 User Domain

The **User** domain handles user management functionality for the Live Studio API.

## 📋 Overview

This module provides:
- User CRUD operations
- User listing with pagination
- Superadmin user management
- Multi-tenant user isolation

## 📁 Structure

```
user/
├── controller/
│   ├── user_controller.go       # Controller interface
│   └── user_controller_impl.go  # Controller implementation
├── entity/
│   └── user_entity.go           # User database model
├── params/
│   ├── request.go               # Request DTOs
│   └── response.go              # Response DTOs
├── repository/
│   ├── user_repository.go       # Repository interface
│   └── user_repository_impl.go  # Repository implementation
├── service/
│   ├── user_service.go          # Service interface
│   └── user_service_impl.go     # Service implementation
├── route.go                     # Route definitions
├── module.go                    # FX module
└── README.md                    # This file
```

## 🗃️ Entity

```go
type User struct {
    ID        *uuid.UUID     // Primary key (UUID v4)
    Name      string         // User name (max 100 chars)
    Email     string         // User email (max 100 chars)
    Password  string         // Hashed password
    RoleID    uint           // Foreign key to Role
    Role      Role           // Role relationship
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt // Soft delete
    TenantBase               // Multi-tenant support
}
```

## 🌐 API Endpoints

### Standard User Routes

Base path: `/api/users`

> **Middleware**: Tenant middleware required

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all users |
| `POST` | `/` | Create new user |
| `GET` | `/:id` | Get user by ID |
| `PUT` | `/:id` | Update user |
| `DELETE` | `/:id` | Delete user (soft delete) |

### Superadmin Routes

Base path: `/api/superadmin/user`

> **Requires**: `superadmin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all users (all tenants) |
| `POST` | `/` | Create new user |
| `GET` | `/:id` | Get user by ID |
| `PUT` | `/:id` | Update user |
| `DELETE` | `/:id` | Delete user |

## 📝 Request/Response Examples

### Create User

**Request:**
```json
POST /api/users
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securePassword123",
  "role_id": 2
}
```

**Response:**
```json
{
  "message": "User created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": {
      "id": 2,
      "name": "admin"
    },
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Users

**Request:**
```http
GET /api/users
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": {
        "id": 2,
        "name": "admin"
      }
    }
  ]
}
```

### Update User

**Request:**
```json
PUT /api/users/550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

## 🔗 Dependencies

- Role domain for user roles
- Password hashing utilities
- Tenant middleware for multi-tenancy
