# 📊 Performa (Performance) Domain

The **Performa** domain provides performance analytics for the Live Studio API.

## 📋 Overview

This module provides:
- Host performance metrics
- Account performance metrics
- Studio performance metrics
- Aggregated analytics data

## 📁 Structure

```
performa/
├── controller/
│   ├── performa_controller.go       # Controller interface
│   └── performa_controller_impl.go  # Controller implementation
├── entity/
│   └── performa_entity.go           # Performance database model
├── params/
│   ├── request.go                   # Request DTOs
│   └── response.go                  # Response DTOs
├── repository/
│   ├── performa_repository.go       # Repository interface
│   └── performa_repository_impl.go  # Repository implementation
├── service/
│   ├── performa_service.go          # Service interface
│   └── performa_service_impl.go     # Service implementation
├── route.go                         # Route definitions
├── module.go                        # FX module
└── README.md                        # This file
```

## 🌐 API Endpoints

Base path: `/api/performa`

> **Requires**: `superadmin` or `admin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/host` | Get all hosts performance |
| `GET` | `/host/:id` | Get specific host performance |
| `GET` | `/account` | Get all accounts performance |
| `GET` | `/studio` | Get all studios performance |
| `GET` | `/studio/:id` | Get specific studio performance |

## 📝 Request/Response Examples

### Get All Hosts Performance

**Request:**
```http
GET /api/performa/host
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "host_id": "uuid-1",
      "host_name": "Jane Live",
      "total_lives": 25,
      "total_duration": 3000,
      "total_views": 125000,
      "avg_peak_viewers": 650,
      "total_sales": 75000000,
      "total_orders": 450,
      "conversion_rate": 3.6
    }
  ]
}
```

### Get Host Performance by ID

**Request:**
```http
GET /api/performa/host/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": {
    "host_id": "550e8400-e29b-41d4-a716-446655440000",
    "host_name": "Jane Live",
    "period": {
      "start": "2025-01-01",
      "end": "2025-01-31"
    },
    "metrics": {
      "total_lives": 25,
      "total_duration_minutes": 3000,
      "total_views": 125000,
      "avg_peak_viewers": 650,
      "total_comments": 8000,
      "total_likes": 45000,
      "followers_growth": 1200,
      "total_placed_orders": 500,
      "total_placed_sales": 85000000,
      "total_confirmed_orders": 450,
      "total_confirmed_sales": 75000000,
      "conversion_rate": 3.6
    }
  }
}
```

### Get All Accounts Performance

**Request:**
```http
GET /api/performa/account
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "account_id": 1,
      "account_name": "Shopee Store A",
      "total_lives": 50,
      "total_sales": 150000000,
      "total_orders": 900,
      "total_commission": 7500000
    }
  ]
}
```

### Get All Studios Performance

**Request:**
```http
GET /api/performa/studio
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
      "host_count": 5,
      "total_lives": 100,
      "total_sales": 300000000,
      "total_orders": 1800,
      "avg_conversion_rate": 3.2
    }
  ]
}
```

### Get Studio Performance by ID

**Request:**
```http
GET /api/performa/studio/1
Authorization: Bearer <jwt_token>
```

## 🔗 Dependencies

- Host domain for host data
- Account domain for account data
- Studio domain for studio data
- Live domain for live session metrics
- Transaction domain for sales data

## 📌 Notes

- Performance data is aggregated from live sessions and transactions
- Metrics can be filtered by date range (query parameters)
- Conversion rate = (Orders / Views) * 100
- Used for business analytics and reporting
