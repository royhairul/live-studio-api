# 📊 Dashboard Domain

The **Dashboard** domain provides aggregated dashboard data for the Live Studio API.

## 📋 Overview

This module provides:
- Aggregated business metrics
- Summary statistics
- Quick overview data for dashboards

## 📁 Structure

```
dashboard/
├── controller/
│   ├── dashboard_controller.go       # Controller interface
│   └── dashboard_controller_impl.go  # Controller implementation
├── params/
│   ├── request.go                    # Request DTOs
│   └── response.go                   # Response DTOs
├── service/
│   ├── dashboard_service.go          # Service interface
│   └── dashboard_service_impl.go     # Service implementation
├── route.go                          # Route definitions
├── module.go                         # FX module
└── README.md                         # This file
```

## 🌐 API Endpoints

Base path: `/api/dashboard`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get dashboard summary data |

## 📝 Request/Response Examples

### Get Dashboard Data

**Request:**
```http
GET /api/dashboard
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": {
    "summary": {
      "total_hosts": 25,
      "total_studios": 5,
      "total_accounts": 10,
      "active_lives_today": 3
    },
    "performance_today": {
      "total_sales": 15000000,
      "total_orders": 150,
      "total_views": 25000,
      "avg_conversion_rate": 3.5
    },
    "performance_this_month": {
      "total_sales": 350000000,
      "total_orders": 3500,
      "total_views": 750000,
      "target_achievement": 75
    },
    "top_hosts": [
      {
        "id": "uuid-1",
        "name": "Jane Live",
        "total_sales": 75000000,
        "total_orders": 450
      }
    ],
    "top_accounts": [
      {
        "id": 1,
        "name": "Shopee Store A",
        "total_sales": 150000000
      }
    ],
    "recent_activities": [
      {
        "type": "live_ended",
        "message": "Live session ended for Studio A",
        "timestamp": "2025-01-15T12:00:00Z"
      }
    ]
  }
}
```

## 🔗 Dependencies

- Host domain for host statistics
- Studio domain for studio statistics
- Account domain for account statistics
- Live domain for live session data
- Transaction domain for sales data
- Performance domain for metrics
- Tenant middleware for multi-tenancy

## 📌 Notes

- Dashboard data is aggregated from multiple domains
- Provides quick overview for admin dashboard
- Data can be filtered by date range
- Used for at-a-glance business intelligence
