# 🔐 Permission Domain

The **Permission** domain manages granular permissions for the Live Studio API.

## 📋 Overview

This module provides:
- Permission CRUD operations
- Grouped permission views
- Fine-grained access control

## 📁 Structure

```
permission/
├── controller/
│   ├── permission_controller.go       # Controller interface
│   └── permission_controller_impl.go  # Controller implementation
├── entity/
│   └── permission_entity.go           # Permission database model
├── params/
│   ├── request.go                     # Request DTOs
│   └── response.go                    # Response DTOs
├── repository/
│   ├── permission_repository.go       # Repository interface
│   └── permission_repository_impl.go  # Repository implementation
├── service/
│   ├── permission_service.go          # Service interface
│   └── permission_service_impl.go     # Service implementation
├── route.go                           # Route definitions
├── module.go                          # FX module
└── README.md                          # This file
```

## 🌐 API Endpoints

Base path: `/api/permission`

> **Requires**: `superadmin` role

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Get all permissions |
| `GET` | `/grouped` | Get permissions grouped by module |
| `POST` | `/` | Create new permission |
| `GET` | `/:id` | Get permission by ID |
| `PUT` | `/:id` | Update permission |
| `DELETE` | `/:id` | Delete permission |

## 📝 Request/Response Examples

### Create Permission

**Request:**
```json
POST /api/permission
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "name": "hosts.delete",
  "display_name": "Delete Host",
  "description": "Can delete hosts",
  "module": "hosts"
}
```

**Response:**
```json
{
  "message": "Permission created successfully",
  "data": {
    "id": 10,
    "name": "hosts.delete",
    "display_name": "Delete Host",
    "description": "Can delete hosts",
    "module": "hosts",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### Get All Permissions

**Request:**
```http
GET /api/permission
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "users.read",
      "display_name": "View Users",
      "module": "users"
    },
    {
      "id": 2,
      "name": "users.create",
      "display_name": "Create User",
      "module": "users"
    }
  ]
}
```

### Get Permissions Grouped

**Request:**
```http
GET /api/permission/grouped
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": {
    "users": [
      {"id": 1, "name": "users.read", "display_name": "View Users"},
      {"id": 2, "name": "users.create", "display_name": "Create User"},
      {"id": 3, "name": "users.update", "display_name": "Update User"},
      {"id": 4, "name": "users.delete", "display_name": "Delete User"}
    ],
    "hosts": [
      {"id": 5, "name": "hosts.read", "display_name": "View Hosts"},
      {"id": 6, "name": "hosts.create", "display_name": "Create Host"},
      {"id": 7, "name": "hosts.update", "display_name": "Update Host"},
      {"id": 8, "name": "hosts.delete", "display_name": "Delete Host"}
    ],
    "studios": [
      {"id": 9, "name": "studios.read", "display_name": "View Studios"},
      {"id": 10, "name": "studios.create", "display_name": "Create Studio"}
    ]
  }
}
```

## 🔗 Dependencies

- Role domain (roles have permissions)

## 📌 Permission Naming Convention

Permissions follow the format: `{module}.{action}`

| Action | Description |
|--------|-------------|
| `read` | View/list resources |
| `create` | Create new resources |
| `update` | Modify existing resources |
| `delete` | Delete resources |

## 📌 Notes

- Only superadmin can manage permissions
- Permissions are assigned to roles, not directly to users
- Grouped endpoint organizes permissions by module for UI rendering
