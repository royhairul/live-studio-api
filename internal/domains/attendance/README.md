# ✅ Attendance Domain

The **Attendance** domain manages host attendance tracking for the Live Studio API.

## 📋 Overview

This module provides:
- Check-in and check-out functionality
- Attendance records listing
- Unchecked-out records tracking
- Multi-tenant attendance isolation

## 📁 Structure

```
attendance/
├── controller/
│   ├── attendance_controller.go       # Controller interface
│   └── attendance_controller_impl.go  # Controller implementation
├── entity/
│   └── attendance_entity.go           # Attendance database model
├── params/
│   ├── request.go                     # Request DTOs
│   └── response.go                    # Response DTOs
├── repository/
│   ├── attendance_repository.go       # Repository interface
│   └── attendance_repository_impl.go  # Repository implementation
├── service/
│   ├── attendance_service.go          # Service interface
│   └── attendance_service_impl.go     # Service implementation
├── route.go                           # Route definitions
├── module.go                          # FX module
└── README.md                          # This file
```

## 🌐 API Endpoints

Base path: `/api/attendance`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all attendance records |
| `GET` | `/unchecked-out` | Get hosts who haven't checked out |
| `POST` | `/check-in` | Record host check-in |
| `POST` | `/check-out` | Record host check-out |

## 📝 Request/Response Examples

### Check In

**Request:**
```json
POST /api/attendance/check-in
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "host_id": "550e8400-e29b-41d4-a716-446655440000",
  "schedule_id": 1
}
```

**Response:**
```json
{
  "message": "Check-in successful",
  "data": {
    "id": 1,
    "host_id": "550e8400-e29b-41d4-a716-446655440000",
    "host_name": "Jane Live",
    "check_in_time": "2025-01-15T08:00:00Z",
    "check_out_time": null,
    "status": "checked_in"
  }
}
```

### Check Out

**Request:**
```json
POST /api/attendance/check-out
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "attendance_id": 1
}
```

**Response:**
```json
{
  "message": "Check-out successful",
  "data": {
    "id": 1,
    "host_id": "550e8400-e29b-41d4-a716-446655440000",
    "host_name": "Jane Live",
    "check_in_time": "2025-01-15T08:00:00Z",
    "check_out_time": "2025-01-15T12:00:00Z",
    "duration_minutes": 240,
    "status": "completed"
  }
}
```

### Get Unchecked Out

**Request:**
```http
GET /api/attendance/unchecked-out
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 2,
      "host": {
        "id": "uuid-2",
        "name": "John Host"
      },
      "check_in_time": "2025-01-15T14:00:00Z",
      "check_out_time": null
    }
  ]
}
```

### Get All Attendance Records

**Request:**
```http
GET /api/attendance
Authorization: Bearer <jwt_token>
```

## 🔗 Dependencies

- Host domain for host reference
- Schedule domain for schedule reference
- Tenant middleware for multi-tenancy

## 📌 Notes

- Attendance is linked to schedules
- Tracks actual working hours vs scheduled hours
- `unchecked-out` endpoint helps identify hosts who forgot to check out
