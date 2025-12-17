# 📦 Pkg (Shared Packages)

The **Pkg** directory contains shared utilities and packages used across the application.

## 📋 Overview

This directory provides:
- Common utilities and helpers
- Shared types and interfaces
- Cross-cutting concerns (validation, error handling, etc.)
- Infrastructure components

## 📁 Structure

```
pkg/
├── constants/          # Application constants
├── errorhandler/       # Error handling utilities
├── httpclient/         # HTTP client abstraction
├── paramhandler/       # Request parameter handling
├── response/           # API response utilities
├── snowflakeid/        # Snowflake ID generator
├── tenantdb/           # Multi-tenant database utilities
├── timehandler/        # Time/date utilities
├── utils/              # General utilities
└── validators/         # Input validation
```

---

## 🔧 Package Details

### 📊 constants

Application-wide constants and configuration values.

---

### ❌ errorhandler

Error handling and formatting utilities.

**Key Functions:**

```go
// Format validation errors into a user-friendly map
func FormatValidationError(err validator.ValidationErrors) map[string]string

// Custom error messages for validation tags
func msgForTag(tag string, param string) string
```

**Validation Messages:**

| Tag | Message |
|-----|---------|
| `required` | "field is required" |
| `email` | "invalid email format" |
| `min` | "value is too short (minimum is X characters)" |
| `max` | "value is too long (maximum is X characters)" |
| `isFile` | "field must be a file" |
| `image` | "field must be an image file" |
| `fileSize` | "file size is too large (maximum is X MB)" |

---

### 🌐 httpclient

HTTP client abstraction for making external API requests.

**Interface:**

```go
type Client interface {
    // Create a new HTTP request
    NewRequest(
        method, rawURL string, 
        query map[string]string, 
        body io.Reader, 
        headers map[string]string,
    ) (*http.Request, error)
    
    // Execute request and parse response
    DoRequest(req *http.Request, out interface{}) error
}
```

**Usage:**
- Used by Shopee client for external API calls
- Abstracts HTTP logic for testability
- Handles JSON serialization/deserialization

---

### 📝 paramhandler

Request parameter handling and parsing utilities.

---

### 📤 response

Standardized API response formatting.

**Base Response Structure:**

```go
type BaseResponse struct {
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

// Create a new response
func NewBaseResponse(message string, data any) *BaseResponse
```

**Example Response:**

```json
{
  "message": "Success",
  "data": {
    "id": 1,
    "name": "Example"
  }
}
```

**Features:**
- Handles nil slices (returns empty array instead of null)
- Omits data field when nil
- Consistent response format across all endpoints

---

### ❄️ snowflakeid

Distributed unique ID generation using Twitter's Snowflake algorithm.

**Initialization:**

```go
var Node *snowflake.Node

func InitSnowflake() (*snowflake.Node, error) {
    Node, err = snowflake.NewNode(1)  // Node ID = 1
    return Node, err
}
```

**Usage:**

```go
// Generate a new Snowflake ID
id := snowflakeid.Node.Generate().Int64()
```

**Benefits:**
- Time-ordered IDs
- Unique across distributed systems
- 64-bit integers (database-friendly)
- Used for Transaction IDs

---

### 🏢 tenantdb

Multi-tenant database utilities using GORM hooks.

**TenantBase Struct:**

```go
type TenantBase struct {
    TenantID string `json:"tenant_id" gorm:"index;nullable"`
}
```

**Automatic Hooks:**

1. **BeforeCreate** - Auto-populates `TenantID` from context:
   ```go
   func (t *TenantBase) BeforeCreate(tx *gorm.DB) error {
       if t.TenantID == "" {
           if ctxTenant, ok := tx.Statement.Context.Value("tenant_id").(string); ok {
               t.TenantID = ctxTenant
           }
       }
       return nil
   }
   ```

2. **BeforeFind** - Auto-filters queries by tenant:
   ```go
   func (t *TenantBase) BeforeFind(tx *gorm.DB) error {
       if ctxTenant, ok := tx.Statement.Context.Value("tenant_id").(string); ok {
           tx.Statement.AddClause(clause.Where{
               Exprs: []clause.Expression{
                   clause.Eq{Column: "tenant_id", Value: ctxTenant},
               },
           })
       }
       return nil
   }
   ```

**How It Works:**
1. Embed `TenantBase` in any entity that needs tenant isolation
2. Middleware sets `tenant_id` in request context
3. GORM hooks automatically filter/populate tenant data

---

### ⏰ timehandler

Time and date handling utilities.

**Files:**
- Time parsing and formatting
- Date range calculations
- Timezone handling

---

### 🛠️ utils

General utility functions.

**JWT Token Utilities:**

```go
// JWT Payload structure
type JWTPayloadDTO struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Role uint   `json:"role"`
    TenantBase
    jwt.RegisteredClaims
}

// Generate JWT token for a user
func GenerateTokenJWT(user *entity.User) (string, error)

// Verify and parse JWT token
func VerifyTokenJWT(tokenStr string, secret string) (*JWTPayloadDTO, error)

// Sign a token with secret key
func SignToken(payload *JWTPayloadDTO, secretKey string) (string, error)
```

**Other Utilities:**
- `display_name.go` - Display name formatting
- `is_today.go` - Date comparison helpers
- `parse_time.go` - Time parsing utilities
- `ptr_time.go` - Pointer helpers for time
- `random_duration.go` - Random duration generation

---

### ✅ validators

Input validation using go-playground/validator.

**Initialization:**

```go
func InitValidator() *validator.Validate
```

**Features:**
- Custom validation rules
- Integration with Gin binding
- Error message formatting

---

## 🔗 Package Dependencies

```
pkg/
├── errorhandler  ← validator/v10
├── httpclient    ← net/http
├── response      ← reflect
├── snowflakeid   ← bwmarrin/snowflake
├── tenantdb      ← gorm.io/gorm
├── utils         ← golang-jwt/jwt/v5, tenantdb
└── validators    ← go-playground/validator/v10
```

---

## 📌 Usage Guidelines

1. **Import Path**: `github.com/royhairul/live-studio-api/internal/pkg/<package>`
2. **Dependency Injection**: Most packages are initialized via Uber FX
3. **No Circular Dependencies**: Pkg packages should not import domain packages
4. **Testing**: Each package should be independently testable
