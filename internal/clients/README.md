# 🔌 Clients

The **Clients** package contains external API client implementations for third-party services.

## 📋 Overview

Clients are responsible for:
- Communicating with external APIs
- Handling authentication and sessions
- Parsing API responses
- Error handling for external services

## 📁 Structure

```
clients/
└── shopee/
    ├── client.go               # Base Shopee client
    ├── provider.go             # Client factory/providers
    ├── params/                 # Request/Response DTOs
    │   └── *.go               # Various parameter structs
    └── service/                # Service implementations
        ├── sp_account_service.go
        ├── sp_account_service_impl.go
        ├── sp_checkout_service.go
        ├── sp_checkout_service_impl.go
        ├── sp_finance_service.go
        ├── sp_finance_service_impl.go
        ├── sp_live_service.go
        └── sp_live_service_impl.go
```

---

## 🛒 Shopee Client

The Shopee client provides integration with Shopee's internal APIs for live streaming, affiliate, and seller data.

### Base Client

```go
type ShopeeClient struct {
    BaseURL string
    Client  httpclient.Client
}

// Create a new request with Shopee-specific headers
func (c *ShopeeClient) NewShopeeRequest(
    method, endpoint string, 
    query map[string]string, 
    body io.Reader, 
    cookie string,
) (*http.Request, error)

// Execute request and parse response
func (c *ShopeeClient) DoShopeeRequest(req *http.Request, out interface{}) error
```

### Client Providers

The package provides factory functions for different Shopee platforms:

| Provider | Base URL | Purpose |
|----------|----------|---------|
| `NewShopeeCreatorClient` | `https://creator.shopee.co.id` | Creator/Streamer dashboard |
| `NewShopeeSellerClient` | `https://seller.shopee.co.id` | Seller center |
| `NewShopeeAffiliateClient` | `https://affiliate.shopee.co.id` | Affiliate program |
| `NewShopeeDefaultClient` | `https://shopee.co.id` | Main Shopee site |

### Default Headers

All Shopee requests include:

```go
headers := map[string]string{
    "User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36...",
    "Content-Type": "application/json",
    "Cookie":       "<session_cookie>",  // From account
}
```

---

## 🎬 Live Service

Service for fetching live streaming data.

### Interface

```go
type ShopeeLiveService interface {
    // Get real-time live sessions
    GetLiveSessionRT(cookie string) ([]ShopeeLiveReportItemRT, error)
    
    // Get dashboard overview for a session
    GetDashboardOverviewRT(cookie, sessionID string) (ShopeeLiveOverviewResponse, error)
    
    // Get buyer analytics
    GetDashboardBuyerRT(cookie, sessionID string) ([]ShopeeLiveAudienceAnalyticsResponse, error)
    
    // Get viewer analytics
    GetDashboardViewerRT(cookie, sessionID string) ([]ShopeeLiveAudienceAnalyticsResponse, error)
    
    // Get viewer source analytics
    GetDashboardViewerSourceRT(cookie, sessionID string) (ShopeeLiveAudienceAnalyticsResponse, error)
    
    // Get product list for a live session
    GetDashboardProductListRT(cookie, sessionID string, page, pageSize int) (
        ShopeeApiPaginationResult[ShopeeLiveProductResponse], 
        error,
    )
}
```

---

## 💰 Finance Service

Service for fetching financial/commission data.

```go
type ShopeeFinanceService interface {
    GetCommissionData(cookie string) (CommissionResponse, error)
}
```

---

## 🛍️ Checkout Service

Service for fetching order/checkout data.

```go
type ShopeeCheckoutService interface {
    GetCheckoutData(cookie string) (CheckoutResponse, error)
}
```

---

## 👤 Account Service

Service for account-related operations.

```go
type ShopeeAccountService interface {
    GetAccountInfo(cookie string) (AccountResponse, error)
}
```

---

## 🔧 Dependency Injection

Clients are registered with Uber FX using named annotations:

```go
fx.Provide(
    fx.Annotate(shopee.NewShopeeCreatorClient, fx.ResultTags(`name:"creatorShopeeClient"`)),
    fx.Annotate(shopee.NewShopeeSellerClient, fx.ResultTags(`name:"sellerShopeeClient"`)),
    fx.Annotate(shopee.NewShopeeAffiliateClient, fx.ResultTags(`name:"affiliateShopeeClient"`)),
    fx.Annotate(shopee.NewShopeeDefaultClient, fx.ResultTags(`name:"defaultShopeeClient"`)),
)
```

---

## 📌 Notes

- Shopee APIs require valid session cookies from logged-in accounts
- Cookies are stored in the Account entity and passed to client methods
- API responses are parsed into strongly-typed structs
- Error handling includes network and API-level errors

## ⚠️ Important

- Shopee's internal APIs may change without notice
- Cookie sessions expire and need to be refreshed
- Rate limiting may apply to API requests
