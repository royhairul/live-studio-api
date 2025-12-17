# 💰 Finance Domain

The **Finance** domain manages financial data and reporting for the Live Studio API.

## 📋 Overview

This module provides:
- Financial records listing
- Revenue and commission tracking
- Multi-tenant financial data

## 📁 Structure

```
finance/
├── controller/
│   ├── finance_controller.go       # Controller interface
│   └── finance_controller_impl.go  # Controller implementation
├── entity/
│   └── finance_entity.go           # Finance database model
├── params/
│   ├── request.go                  # Request DTOs
│   └── response.go                 # Response DTOs
├── repository/
│   ├── finance_repository.go       # Repository interface
│   └── finance_repository_impl.go  # Repository implementation
├── service/
│   ├── finance_service.go          # Service interface
│   └── finance_service_impl.go     # Service implementation
├── route.go                        # Route definitions
├── module.go                       # FX module
└── README.md                       # This file
```

## 🌐 API Endpoints

Base path: `/api/finance`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all financial records |

## 📝 Request/Response Examples

### Get All Financial Records

**Request:**
```http
GET /api/finance
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "period": "2025-01",
      "account": {
        "id": 1,
        "name": "Shopee Store A"
      },
      "total_sales": 150000000,
      "total_commission": 7500000,
      "total_orders": 900,
      "settlement_status": "completed"
    }
  ]
}
```

## 🔗 Dependencies

- Account domain for account reference
- Transaction domain for sales data
- Tenant middleware for multi-tenancy

## 📌 Notes

- Financial data is aggregated from transactions
- Used for accounting and settlement purposes
- Tracks commission earnings from sales
