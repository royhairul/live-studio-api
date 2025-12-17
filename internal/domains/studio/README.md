# 🏢 Studio Domain

The **Studio** domain manages physical streaming studio locations for the Live Studio API.

## 📋 Overview

This module provides:
- Studio CRUD operations
- Studio listing with pagination
- Multi-tenant studio isolation

## 📁 Structure

```
studio/
├── controller/
│   ├── studio_controller.go       # Controller interface
│   └── studio_controller_impl.go  # Controller implementation
├── entity/
│   └── studio_entity.go           # Studio database model
├── params/
│   ├── request.go                 # Request DTOs
│   └── response.go                # Response DTOs
├── repository/
│   ├── studio_repository.go       # Repository interface
│   └── studio_repository_impl.go  # Repository implementation
├── service/
│   ├── studio_service.go          # Service interface
│   └── studio_service_impl.go     # Service implementation
├── route.go                       # Route definitions
├── module.go                      # FX module
└── README.md                      # This file
```

## 🗃️ Entity

```go
type Studio struct {
    gorm.Model               // ID, CreatedAt, UpdatedAt, DeletedAt
    Name       string        // Studio name (max 100 chars)
    Address    string        // Studio address (max 255 chars)
    TenantBase               // Multi-tenant support
}
```

## 🌐 API Endpoints

Base path: `/api/studio`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all studios |
| `POST` | `/` | Create new studio |
| `GET` | `/:id` | Get studio by ID |
| `PUT` | `/:id` | Update studio |
| `DELETE` | `/:id` | Delete studio (soft delete) |

## 📝 Request/Response Examples

### Create Studio

**Request:**
```json
POST /api/studio
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Studio A",
  "address": "123 Main Street, Jakarta"
}
```

**Response:**
```json
{
  "message": "Studio created successfully",
  "data": {
    "id": 1,
    "name": "Studio A",
    "address": "123 Main Street, Jakarta",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Studios

**Request:**
```http
GET /api/studio
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Studio A",
      "address": "123 Main Street, Jakarta"
    },
    {
      "id": 2,
      "name": "Studio B",
      "address": "456 Second Avenue, Bandung"
    }
  ]
}
```

### Update Studio

**Request:**
```json
PUT /api/studio/1
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Studio A - Updated",
  "address": "123 Main Street, Jakarta (Updated)"
}
```

### Delete Studio

**Request:**
```http
DELETE /api/studio/1
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Studio deleted successfully"
}
```

## 🔗 Dependencies

- Tenant middleware for multi-tenancy
- Referenced by: Host, Account, Live, Performance

## 📌 Notes

- Studios are soft-deleted (can be recovered)
- Deleting a studio does not automatically delete associated hosts or accounts
- Studios are used as a grouping mechanism for hosts, accounts, and performance metrics
