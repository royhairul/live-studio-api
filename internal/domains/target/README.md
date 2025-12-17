# 🎯 Target Domain

The **Target** domain manages sales and performance targets for the Live Studio API.

## 📋 Overview

This module provides:
- Target CRUD operations
- Sales target management
- Performance goal tracking
- Multi-tenant target isolation

## 📁 Structure

```
target/
├── controller/
│   ├── target_controller.go       # Controller interface
│   └── target_controller_impl.go  # Controller implementation
├── entity/
│   └── target_entity.go           # Target database model
├── params/
│   ├── request.go                 # Request DTOs
│   └── response.go                # Response DTOs
├── repository/
│   ├── target_repository.go       # Repository interface
│   └── target_repository_impl.go  # Repository implementation
├── service/
│   ├── target_service.go          # Service interface
│   └── target_service_impl.go     # Service implementation
├── route.go                       # Route definitions
├── module.go                      # FX module
└── README.md                      # This file
```

## 🌐 API Endpoints

Base path: `/api/target`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all targets |
| `POST` | `/` | Create new target |
| `GET` | `/:id` | Get target by ID |
| `PUT` | `/:id` | Update target |
| `DELETE` | `/:id` | Delete target (soft delete) |

## 📝 Request/Response Examples

### Create Target

**Request:**
```json
POST /api/target
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "January 2025 Sales Target",
  "type": "sales",
  "target_value": 100000000,
  "period_start": "2025-01-01",
  "period_end": "2025-01-31",
  "host_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**
```json
{
  "message": "Target created successfully",
  "data": {
    "id": 1,
    "name": "January 2025 Sales Target",
    "type": "sales",
    "target_value": 100000000,
    "current_value": 0,
    "progress_percentage": 0,
    "period_start": "2025-01-01",
    "period_end": "2025-01-31",
    "host": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Jane Live"
    },
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Targets

**Request:**
```http
GET /api/target
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "January 2025 Sales Target",
      "type": "sales",
      "target_value": 100000000,
      "current_value": 45000000,
      "progress_percentage": 45,
      "period_start": "2025-01-01",
      "period_end": "2025-01-31",
      "status": "in_progress"
    }
  ]
}
```

### Update Target

**Request:**
```json
PUT /api/target/1
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "target_value": 120000000
}
```

## 🔗 Dependencies

- Host domain for host-specific targets
- Studio domain for studio-specific targets
- Performance domain for progress tracking
- Tenant middleware for multi-tenancy

## 📌 Notes

- Targets can be set for hosts, studios, or accounts
- Progress is calculated based on actual performance data
- Target types: sales, orders, views, followers, etc.
- Used for goal setting and performance monitoring
