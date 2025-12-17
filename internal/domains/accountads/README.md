# 📺 Account Ads Domain

The **Account Ads** domain manages advertising accounts for the Live Studio API.

## 📋 Overview

This module provides:
- Advertising account CRUD operations
- Ad account management for platforms
- Multi-tenant ad account isolation

## 📁 Structure

```
accountads/
├── controller/
│   ├── accountads_controller.go       # Controller interface
│   └── accountads_controller_impl.go  # Controller implementation
├── entity/
│   └── accountads_entity.go           # Account Ads database model
├── params/
│   ├── request.go                     # Request DTOs
│   └── response.go                    # Response DTOs
├── repository/
│   ├── accountads_repository.go       # Repository interface
│   └── accountads_repository_impl.go  # Repository implementation
├── service/
│   ├── accountads_service.go          # Service interface
│   └── accountads_service_impl.go     # Service implementation
├── route.go                           # Route definitions
├── module.go                          # FX module
└── README.md                          # This file
```

## 🌐 API Endpoints

Base path: `/api/accountads`

> **Requires**: `superadmin` or `admin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all ad accounts |
| `POST` | `/` | Create new ad account |
| `GET` | `/:id` | Get ad account by ID |
| `PUT` | `/:id` | Update ad account |
| `DELETE` | `/:id` | Delete ad account (soft delete) |

## 📝 Request/Response Examples

### Create Ad Account

**Request:**
```json
POST /api/accountads
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "Shopee Ads Account",
  "account_id": 1,
  "platform": "shopee",
  "status": "active"
}
```

**Response:**
```json
{
  "message": "Ad account created successfully",
  "data": {
    "id": 1,
    "name": "Shopee Ads Account",
    "account_id": 1,
    "platform": "shopee",
    "status": "active",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Ad Accounts

**Request:**
```http
GET /api/accountads
Authorization: Bearer <jwt_token>
```

## 🔗 Dependencies

- Account domain for main account reference
- Tenant middleware for multi-tenancy

## 📌 Notes

- Ad accounts are linked to main platform accounts
- Used for managing advertising campaigns
- Tracks ad spend and performance
