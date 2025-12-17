# 🛡️ Role Domain

The **Role** domain manages user roles for the Live Studio API.

## 📋 Overview

This module provides:
- Role CRUD operations
- Role-based access control foundation
- Permission assignment to roles
- Multi-tenant role isolation

## 📁 Structure

```
role/
├── controller/
│   ├── role_controller.go       # Controller interface
│   └── role_controller_impl.go  # Controller implementation
├── entity/
│   └── role_entity.go           # Role database model
├── params/
│   ├── request.go               # Request DTOs
│   └── response.go              # Response DTOs
├── repository/
│   ├── role_repository.go       # Repository interface
│   └── role_repository_impl.go  # Repository implementation
├── service/
│   ├── role_service.go          # Service interface
│   └── role_service_impl.go     # Service implementation
├── route.go                     # Route definitions
├── module.go                    # FX module
└── README.md                    # This file
```

## 🌐 API Endpoints

Base path: `/api/role`

> **Requires**: `superadmin` role + Tenant middleware

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all roles |
| `POST` | `/` | Create new role |
| `GET` | `/:id` | Get role by ID |
| `PUT` | `/:id` | Update role |
| `DELETE` | `/:id` | Delete role (soft delete) |

## 📝 Request/Response Examples

### Create Role

**Request:**
```json
POST /api/role
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "manager",
  "display_name": "Manager",
  "description": "Studio manager with limited admin access",
  "permissions": [1, 2, 3, 5, 8]
}
```

**Response:**
```json
{
  "message": "Role created successfully",
  "data": {
    "id": 3,
    "name": "manager",
    "display_name": "Manager",
    "description": "Studio manager with limited admin access",
    "permissions": [
      {"id": 1, "name": "users.read"},
      {"id": 2, "name": "users.create"},
      {"id": 3, "name": "users.update"},
      {"id": 5, "name": "hosts.read"},
      {"id": 8, "name": "schedules.read"}
    ],
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Roles

**Request:**
```http
GET /api/role
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "superadmin",
      "display_name": "Super Administrator",
      "description": "Full system access"
    },
    {
      "id": 2,
      "name": "admin",
      "display_name": "Administrator",
      "description": "Tenant administrator"
    },
    {
      "id": 3,
      "name": "host",
      "display_name": "Host",
      "description": "Live streaming host"
    }
  ]
}
```

## 🔗 Dependencies

- Permission domain for role-permission assignment
- User domain (users have roles)
- Tenant middleware for multi-tenancy

## 📌 Default Roles

| Role | Description |
|------|-------------|
| `superadmin` | Full system access across all tenants |
| `admin` | Tenant-level administrator |
| `host` | Live streaming host with limited access |

## 📌 Notes

- Only superadmin can manage roles
- Roles define what actions users can perform
- Roles are linked to permissions for fine-grained access control
