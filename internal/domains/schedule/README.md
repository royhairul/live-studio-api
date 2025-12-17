# 📅 Schedule Domain

The **Schedule** domain handles host scheduling functionality for the Live Studio API.

## 📋 Overview

This module provides:
- Schedule CRUD operations
- Assign hosts to specific shifts and dates
- Query schedules by shift and date
- Multi-tenant schedule isolation

## 📁 Structure

```
schedule/
├── controller/
│   ├── schedule_controller.go       # Controller interface
│   └── schedule_controller_impl.go  # Controller implementation
├── entity/
│   └── schedule_entity.go           # Schedule database model
├── params/
│   ├── request.go                   # Request DTOs
│   └── response.go                  # Response DTOs
├── repository/
│   ├── schedule_repository.go       # Repository interface
│   └── schedule_repository_impl.go  # Repository implementation
├── service/
│   ├── schedule_service.go          # Service interface
│   └── schedule_service_impl.go     # Service implementation
├── route.go                         # Route definitions
├── module.go                        # FX module
└── README.md                        # This file
```

## 🌐 API Endpoints

Base path: `/api/schedule`

> **Requires**: `superadmin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all schedules |
| `POST` | `/` | Create new schedule |
| `GET` | `/scheduled` | Get schedules by shift and date |
| `GET` | `/:id` | Get schedule by ID |
| `PUT` | `/:id` | Update schedule |
| `DELETE` | `/:id` | Delete schedule (soft delete) |

## 📝 Request/Response Examples

### Create Schedule

**Request:**
```json
POST /api/schedule
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "host_id": "550e8400-e29b-41d4-a716-446655440000",
  "shift_id": 1,
  "date": "2025-01-15"
}
```

**Response:**
```json
{
  "message": "Schedule created successfully",
  "data": {
    "id": 1,
    "host": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Jane Live"
    },
    "shift": {
      "id": 1,
      "name": "Morning Shift",
      "start_time": "08:00",
      "end_time": "12:00"
    },
    "date": "2025-01-15",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get Schedules by Shift and Date

**Request:**
```http
GET /api/schedule/scheduled?shift_id=1&date=2025-01-15
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "host": {
        "id": "uuid-1",
        "name": "Jane Live"
      },
      "shift": {
        "id": 1,
        "name": "Morning Shift"
      },
      "date": "2025-01-15"
    }
  ]
}
```

### Get All Schedules

**Request:**
```http
GET /api/schedule
Authorization: Bearer <jwt_token>
```

## 🔗 Dependencies

- Host domain for host assignment
- Shift domain for shift assignment
- Tenant middleware for multi-tenancy

## 📌 Notes

- Schedules link hosts to specific shifts on specific dates
- Used in conjunction with attendance tracking
- Only superadmin can manage schedules
