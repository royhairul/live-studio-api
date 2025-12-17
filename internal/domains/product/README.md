# 📦 Product Domain

The **Product** domain manages product catalog for the Live Studio API.

## 📋 Overview

This module provides:
- Product CRUD operations
- Product catalog management
- Multi-tenant product isolation

## 📁 Structure

```
product/
├── controller/
│   ├── product_controller.go       # Controller interface
│   └── product_controller_impl.go  # Controller implementation
├── entity/
│   └── product_entity.go           # Product database model
├── params/
│   ├── request.go                  # Request DTOs
│   └── response.go                 # Response DTOs
├── repository/
│   ├── product_repository.go       # Repository interface
│   └── product_repository_impl.go  # Repository implementation
├── service/
│   ├── product_service.go          # Service interface
│   └── product_service_impl.go     # Service implementation
├── route.go                        # Route definitions
├── module.go                       # FX module
└── README.md                       # This file
```

## 🌐 API Endpoints

Base path: `/api/product`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all products |
| `POST` | `/` | Create new product |
| `GET` | `/:id` | Get product by ID |
| `PUT` | `/:id` | Update product |
| `DELETE` | `/:id` | Delete product (soft delete) |

## 📝 Request/Response Examples

### Create Product

**Request:**
```json
POST /api/product
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Premium T-Shirt",
  "sku": "TS-001",
  "price": 150000,
  "description": "High quality cotton t-shirt"
}
```

**Response:**
```json
{
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "name": "Premium T-Shirt",
    "sku": "TS-001",
    "price": 150000,
    "description": "High quality cotton t-shirt",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Products

**Request:**
```http
GET /api/product
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Premium T-Shirt",
      "sku": "TS-001",
      "price": 150000
    }
  ]
}
```

## 🔗 Dependencies

- Order domain references products
- Tenant middleware for multi-tenancy

## 📌 Notes

- Products can be linked to orders
- Products are soft-deleted for data integrity
- Used for tracking items sold during live sessions
