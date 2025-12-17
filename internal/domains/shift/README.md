# ⏰ Shift Domain

The **Shift** domain manages work shift configurations for the Live Studio API.

## 📋 Overview

This module provides:
- Shift CRUD operations
- Work schedule time slot definitions
- Multi-tenant shift isolation

## 📁 Structure

```
shift/
├── controller/
│   ├── shift_controller.go       # Controller interface
│   └── shift_controller_impl.go  # Controller implementation
├── entity/
│   └── shift_entity.go           # Shift database model
├── params/
│   ├── request.go                # Request DTOs
│   └── response.go               # Response DTOs
├── repository/
│   ├── shift_repository.go       # Repository interface
│   └── shift_repository_impl.go  # Repository implementation
├── service/
│   ├── shift_service.go          # Service interface
│   └── shift_service_impl.go     # Service implementation
├── route.go                      # Route definitions
├── module.go                     # FX module
└── README.md                     # This file
```

## 🌐 API Endpoints

Base path: `/api/shift`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all shifts |
| `POST` | `/` | Create new shift |
| `GET` | `/:id` | Get shift by ID |
| `PUT` | `/:id` | Update shift |
| `DELETE` | `/:id` | Delete shift (soft delete) |

## 📝 Request/Response Examples

### Create Shift

**Request:**
```json
POST /api/shift
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Morning Shift",
  "start_time": "08:00",
  "end_time": "12:00"
}
```

**Response:**
```json
{
  "message": "Shift created successfully",
  "data": {
    "id": 1,
    "name": "Morning Shift",
    "start_time": "08:00",
    "end_time": "12:00",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Shifts

**Request:**
```http
GET /api/shift
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Morning Shift",
      "start_time": "08:00",
      "end_time": "12:00"
    },
    {
      "id": 2,
      "name": "Afternoon Shift",
      "start_time": "12:00",
      "end_time": "18:00"
    },
    {
      "id": 3,
      "name": "Evening Shift",
      "start_time": "18:00",
      "end_time": "22:00"
    }
  ]
}
```

## 🔗 Dependencies

- Referenced by: Schedule domain
- Tenant middleware for multi-tenancy

## 📌 Notes

- Shifts define time slots that can be assigned to hosts via schedules
- Shifts are used in conjunction with schedules to track host working hours
