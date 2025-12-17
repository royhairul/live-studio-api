# 📺 Live Domain

The **Live** domain manages live streaming sessions for the Live Studio API.

## 📋 Overview

This module provides:
- Live session data from Shopee
- Live session details with metrics
- Integration with Shopee Live API

## 📁 Structure

```
live/
├── controller/
│   ├── live_controller.go       # Controller interface
│   └── live_controller_impl.go  # Controller implementation
├── entity/
│   └── live_entity.go           # Live session database model
├── params/
│   ├── request.go               # Request DTOs
│   └── response.go              # Response DTOs
├── repository/
│   ├── live_repository.go       # Repository interface
│   └── live_repository_impl.go  # Repository implementation
├── service/
│   ├── live_service.go          # Service interface
│   └── live_service_impl.go     # Service implementation
├── route.go                     # Route definitions
├── module.go                    # FX module
└── README.md                    # This file
```

## 🗃️ Entity

```go
type Live struct {
    ID               uuid.UUID  // Primary key (UUID v4)
    SessionID        int64      // Shopee live session ID (unique)
    Title            string     // Live session title
    StartTime        time.Time  // When the live started
    Duration         int        // Duration in minutes
    Views            int        // Total views
    PeakViewers      int        // Maximum concurrent viewers
    AvgViewDuration  float64    // Average view duration in seconds
    Comments         int        // Total comments
    Likes            int        // Total likes
    FollowersGrowth  int        // New followers gained
    PlacedOrders     int        // Orders placed during live
    PlacedSales      float64    // Sales amount from placed orders
    ConfirmedOrders  int        // Confirmed orders
    ConfirmedSales   float64    // Sales amount from confirmed orders
    HostID           uuid.UUID  // Foreign key to Host
    StudioID         uuid.UUID  // Foreign key to Studio
    AccountID        uuid.UUID  // Foreign key to Account
}
```

## 🌐 API Endpoints

Base path: `/api/live`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/shopee` | Get Shopee live sessions list | ✅ |
| `GET` | `/shopee/:id/:sessionId` | Get specific live session details | ✅ |

## 📝 Request/Response Examples

### Get Shopee Live Sessions

**Request:**
```http
GET /api/live/shopee
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "session_id": 1234567890,
      "title": "Flash Sale Friday!",
      "start_time": "2025-01-15T19:00:00Z",
      "duration": 120,
      "views": 5000,
      "peak_viewers": 850,
      "host": {
        "id": "uuid",
        "name": "Jane Live"
      },
      "account": {
        "id": 1,
        "name": "Shopee Store A"
      }
    }
  ]
}
```

### Get Live Session Details

**Request:**
```http
GET /api/live/shopee/1/1234567890
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "session_id": 1234567890,
    "title": "Flash Sale Friday!",
    "start_time": "2025-01-15T19:00:00Z",
    "duration": 120,
    "views": 5000,
    "peak_viewers": 850,
    "avg_view_duration": 45.5,
    "comments": 320,
    "likes": 1500,
    "followers_growth": 75,
    "placed_orders": 45,
    "placed_sales": 15000000,
    "confirmed_orders": 38,
    "confirmed_sales": 12500000,
    "host": {
      "id": "uuid",
      "name": "Jane Live"
    },
    "studio": {
      "id": "uuid",
      "name": "Studio A"
    },
    "account": {
      "id": 1,
      "name": "Shopee Store A"
    }
  }
}
```

## 🔗 Dependencies

- Host domain for host reference
- Studio domain for studio reference
- Account domain for account reference
- Shopee client for API integration

## 📌 Notes

- Live sessions are fetched from Shopee API
- Session metrics are updated periodically
- Data includes comprehensive performance metrics
- Used by Performance domain for analytics
