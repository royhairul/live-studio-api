# 🔐 Auth Domain

The **Auth** domain handles all authentication and authorization functionality for the Live Studio API.

## 📋 Overview

This module provides:
- User login/logout
- User registration
- Password reset flow (forgot password, OTP verification, reset)
- Current user information retrieval

## 📁 Structure

```
auth/
├── controller/
│   ├── auth_controller.go       # Controller interface
│   └── auth_controller_impl.go  # Controller implementation
├── entity/
│   └── auth_entity.go           # Auth-related entities
├── params/
│   ├── auth_request.go          # Request DTOs
│   └── auth_response.go         # Response DTOs
├── repository/
│   ├── auth_repository.go       # Repository interface
│   └── auth_repository_impl.go  # Repository implementation
├── service/
│   ├── auth_service.go          # Service interface
│   └── auth_service_impl.go     # Service implementation
├── route.go                     # Route definitions
├── module.go                    # FX module
└── README.md                    # This file
```

## 🌐 API Endpoints

Base path: `/api/auth`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/login` | Authenticate user and get JWT token | ❌ |
| `POST` | `/register` | Register a new user | ❌ |
| `POST` | `/forgot-password` | Request password reset OTP | ❌ |
| `POST` | `/verify-otp` | Verify OTP code | ❌ |
| `POST` | `/reset-password` | Reset password with verified OTP | ❌ |
| `GET` | `/me` | Get current authenticated user info | ✅ (superadmin, admin, host) |

## 🔒 Authorization

- Most endpoints are public (no authentication required)
- `/me` endpoint requires authentication and one of these roles: `superadmin`, `admin`, `host`

## 📝 Request/Response Examples

### Login

**Request:**
```json
POST /api/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "user@example.com",
      "role": "admin"
    }
  }
}
```

### Register

**Request:**
```json
POST /api/auth/register
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123"
}
```

### Get Current User

**Request:**
```http
GET /api/auth/me
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Success",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "user@example.com",
    "role": {
      "id": 1,
      "name": "admin"
    }
  }
}
```

## 🔗 Dependencies

- JWT package for token generation (`github.com/golang-jwt/jwt/v5`)
- User repository for user data
- Gomail for SMTP email delivery (`gopkg.in/gomail.v2`)
- Password hashing utilities (`golang.org/x/crypto/bcrypt`)

---

## 📧 Email/SMTP Implementation

The Auth domain uses **SMTP** to send password reset OTP emails via the `gomail.v2` library.

### Environment Variables Required

```env
EMAIL_FROM=your-email@gmail.com
EMAIL_PASSWORD=xxxx xxxx xxxx xxxx    # Gmail App Password
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
```

### How Password Reset Works

```
┌─────────────────────────────────────────────────────────────────┐
│                    Password Reset Flow                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. POST /api/auth/forgot-password                              │
│     └── User submits email                                      │
│         └── Server generates OTP                                │
│             └── OTP sent via SMTP email                         │
│                                                                 │
│  2. POST /api/auth/verify-otp                                   │
│     └── User submits email + OTP                                │
│         └── Server validates OTP                                │
│             └── Returns success if valid                        │
│                                                                 │
│  3. POST /api/auth/reset-password                               │
│     └── User submits email + new password                       │
│         └── Server updates password                             │
│             └── Returns success                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Email Template (HTML)

The OTP email is sent with the following HTML format:

```html
Hello <b>{user_name}</b>,
<br>
Your OTP is: <b>{otp_code}</b>
```

### Gmail SMTP Setup

1. Enable **2-Factor Authentication** on your Google account
2. Generate an **App Password**:
   - Go to [Google Account Security](https://myaccount.google.com/security)
   - Navigate to "2-Step Verification" → "App passwords"
   - Select "Mail" and generate password
3. Use the 16-character password in `EMAIL_PASSWORD`

### Supported SMTP Providers

| Provider | Host | Port |
|----------|------|------|
| Gmail | `smtp.gmail.com` | `587` |
| Outlook | `smtp-mail.outlook.com` | `587` |
| Yahoo | `smtp.mail.yahoo.com` | `587` |
| SendGrid | `smtp.sendgrid.net` | `587` |

---

## 📌 Notes

- OTP codes are stored in the database with expiration
- Previous OTPs are deleted when a new one is generated
- Passwords are hashed using bcrypt before storage
- JWT tokens contain user ID, name, role, and tenant information
