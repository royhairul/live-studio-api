# 💳 Transaction Domain

The **Transaction** domain manages transaction tracking for the Live Studio API.

## 📋 Overview

This module provides:
- Transaction CRUD operations
- Transaction grouping and aggregation
- Commission tracking
- Multi-tenant transaction isolation

## 📁 Structure

```
transaction/
├── controller/
│   ├── transaction_controller.go       # Controller interface
│   └── transaction_controller_impl.go  # Controller implementation
├── entity/
│   └── transaction_entity.go           # Transaction database model
├── params/
│   ├── request.go                      # Request DTOs
│   └── response.go                     # Response DTOs
├── repository/
│   ├── transaction_repository.go       # Repository interface
│   └── transaction_repository_impl.go  # Repository implementation
├── service/
│   ├── transaction_service.go          # Service interface
│   └── transaction_service_impl.go     # Service implementation
├── route.go                            # Route definitions
├── module.go                           # FX module
└── README.md                           # This file
```

## 🗃️ Entity

```go
type Transaction struct {
    ID                      int64      // Snowflake ID (primary key)
    UniqueID                string     // Platform unique transaction ID
    Status                  string     // Transaction status
    TotalPurchase           int64      // Total purchase amount
    TotalCommission         int64      // Commission earned
    TotalCommissionWithMCN  int64      // Commission with MCN fee
    PurchaseTime            *time.Time // When purchase was made
    CompleteTime            *time.Time // When transaction completed
    Orders                  []Order    // Associated orders
    AccountID               uint       // Foreign key to Account
    Account                 Account    // Account relationship
    CreatedAt               time.Time
    UpdatedAt               time.Time
    DeletedAt               gorm.DeletedAt
    TenantBase                         // Multi-tenant support
}
```

## 🌐 API Endpoints

Base path: `/api/transaction`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all transactions |
| `POST` | `/` | Create new transaction |
| `GET` | `/grouped` | Get transactions grouped (by account, period, etc.) |
| `GET` | `/:id` | Get transaction by ID |
| `PUT` | `/:id` | Update transaction |
| `DELETE` | `/:id` | Delete transaction (soft delete) |

## 📝 Request/Response Examples

### Create Transaction

**Request:**
```json
POST /api/transaction
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "unique_id": "TRX123456789",
  "status": "completed",
  "total_purchase": 500000,
  "total_commission": 25000,
  "total_commission_with_mcn": 20000,
  "purchase_time": "2025-01-15T10:00:00Z",
  "complete_time": "2025-01-15T12:00:00Z",
  "account_id": 1
}
```

**Response:**
```json
{
  "message": "Transaction created successfully",
  "data": {
    "id": 1234567890123456,
    "unique_id": "TRX123456789",
    "status": "completed",
    "total_purchase": 500000,
    "total_commission": 25000,
    "total_commission_with_mcn": 20000,
    "account": {
      "id": 1,
      "name": "Shopee Store A"
    },
    "created_at": "2025-01-15T12:00:00Z"
  }
}
```

### Get All Transactions

**Request:**
```http
GET /api/transaction
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1234567890123456,
      "unique_id": "TRX123456789",
      "status": "completed",
      "total_purchase": 500000,
      "total_commission": 25000,
      "account": {
        "id": 1,
        "name": "Shopee Store A"
      },
      "orders": [
        {
          "id": 1,
          "product_name": "Premium T-Shirt",
          "quantity": 2,
          "price": 250000
        }
      ]
    }
  ]
}
```

### Get Transactions Grouped

**Request:**
```http
GET /api/transaction/grouped?group_by=account&period=2025-01
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
      "period": "2025-01",
      "transaction_count": 150,
      "total_purchase": 75000000,
      "total_commission": 3750000
    }
  ]
}
```

## 🔗 Dependencies

- Account domain for account reference
- Order domain for order items
- Tenant middleware for multi-tenancy
- Snowflake ID generator for unique IDs

## 📌 Notes

- Transactions use Snowflake IDs for unique identification
- Commission calculation includes MCN (Multi-Channel Network) fees
- Transactions can have multiple orders
- Used for tracking affiliate/creator earnings
