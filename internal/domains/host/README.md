# 🎙️ Host Domain

The **Host** domain manages live streaming hosts (talents) for the Live Studio API.

## 📋 Overview

This module provides:
- Host CRUD operations
- Host-to-studio assignment
- Hosts grouped by studio view
- Multi-tenant host isolation

## 📁 Structure

```
host/
├── controller/
│   ├── host_controller.go       # Controller interface
│   └── host_controller_impl.go  # Controller implementation
├── entity/
│   └── host_entity.go           # Host database model
├── params/
│   ├── request.go               # Request DTOs
│   └── response.go              # Response DTOs
├── repository/
│   ├── host_repository.go       # Repository interface
│   └── host_repository_impl.go  # Repository implementation
├── service/
│   ├── host_service.go          # Service interface
│   └── host_service_impl.go     # Service implementation
├── route.go                     # Route definitions
├── module.go                    # FX module
└── README.md                    # This file
```

## 🗃️ Entity

```go
type Host struct {
    ID        *uuid.UUID     // Primary key (UUID v4)
    Name      string         // Host name (max 100 chars)
    Phone     string         // Phone number (max 15 chars)
    UserID    uint           // Optional linked user account
    StudioID  uint           // Foreign key to Studio
    Studio    Studio         // Studio relationship
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt // Soft delete
    TenantBase               // Multi-tenant support
}
```

## 🌐 API Endpoints

Base path: `/api/host`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all hosts |
| `POST` | `/` | Create new host |
| `GET` | `/group-by-studio` | Get all hosts grouped by studio |
| `GET` | `/:id` | Get host by ID |
| `PUT` | `/:id` | Update host |
| `DELETE` | `/:id` | Delete host (soft delete) |

## 📝 Request/Response Examples

### Create Host

**Request:**
```json
POST /api/host
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Jane Live",
  "phone": "081234567890",
  "studio_id": 1
}
```

**Response:**
```json
{
  "message": "Host created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Jane Live",
    "phone": "081234567890",
    "studio": {
      "id": 1,
      "name": "Studio A"
    },
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get Hosts Grouped by Studio

**Request:**
```http
GET /api/host/group-by-studio
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "studio_id": 1,
      "studio_name": "Studio A",
      "hosts": [
        {
          "id": "uuid-1",
          "name": "Jane Live",
          "phone": "081234567890"
        },
        {
          "id": "uuid-2",
          "name": "John Host",
          "phone": "081234567891"
        }
      ]
    },
    {
      "studio_id": 2,
      "studio_name": "Studio B",
      "hosts": [...]
    }
  ]
}
```

### Get All Hosts

**Request:**
```http
GET /api/host
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Jane Live",
      "phone": "081234567890",
      "studio": {
        "id": 1,
        "name": "Studio A"
      }
    }
  ]
}
```

## 🔗 Dependencies

- Studio domain for studio assignment
- User domain for optional user linking
- Tenant middleware for multi-tenancy
