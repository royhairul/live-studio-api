# 🛒 Order Domain

The **Order** domain manages order tracking for the Live Studio API.

## 📋 Overview

This module provides:
- Order CRUD operations
- Order tracking and management
- Integration with transactions

## 📁 Structure

```
order/
├── controller/
│   ├── order_controller.go       # Controller interface
│   └── order_controller_impl.go  # Controller implementation
├── entity/
│   └── order_entity.go           # Order database model
├── params/
│   ├── request.go                # Request DTOs
│   └── response.go               # Response DTOs
├── repository/
│   ├── order_repository.go       # Repository interface
│   └── order_repository_impl.go  # Repository implementation
├── service/
│   ├── order_service.go          # Service interface
│   └── order_service_impl.go     # Service implementation
├── route.go                      # Route definitions
├── module.go                     # FX module
└── README.md                     # This file
```

## 🌐 API Endpoints

Base path: `/api/order`

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all orders |
| `POST` | `/` | Create new order |
| `GET` | `/:id` | Get order by ID |
| `PUT` | `/:id` | Update order |
| `DELETE` | `/:id` | Delete order (soft delete) |

## 📝 Request/Response Examples

### Create Order

**Request:**
```json
POST /api/order
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "transaction_id": 1234567890,
  "product_id": 1,
  "quantity": 2,
  "price": 150000,
  "status": "pending"
}
```

**Response:**
```json
{
  "message": "Order created successfully",
  "data": {
    "id": 1,
    "transaction_id": 1234567890,
    "product": {
      "id": 1,
      "name": "Premium T-Shirt"
    },
    "quantity": 2,
    "price": 150000,
    "total": 300000,
    "status": "pending",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Orders

**Request:**
```http
GET /api/order
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "transaction_id": 1234567890,
      "product": {
        "id": 1,
        "name": "Premium T-Shirt"
      },
      "quantity": 2,
      "price": 150000,
      "total": 300000,
      "status": "completed"
    }
  ]
}
```

## 🔗 Dependencies

- Transaction domain (orders belong to transactions)
- Product domain for product reference

## 📌 Notes

- Orders are line items within transactions
- Order status tracks the fulfillment process
- Orders are linked to products and transactions
