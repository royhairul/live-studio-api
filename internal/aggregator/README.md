# 📊 Aggregator

The **Aggregator** package contains services that aggregate data from multiple domains to provide consolidated metrics and analytics.

## 📋 Overview

Aggregators are responsible for:
- Combining data from multiple sources
- Calculating complex metrics
- Providing summary statistics
- Cross-domain data processing

## 📁 Structure

```
aggregator/
└── performa/
    ├── module.go                    # FX module definition
    ├── performa_aggregator.go       # Interface definition
    └── performa_aggregator_impl.go  # Implementation
```

---

## 📈 Performa Aggregator

The **Performa Aggregator** calculates performance metrics across hosts, accounts, and studios.

### Interface

```go
type PerformaAggregator interface {
    // Calculate overall performance for all studios
    Calculate(ctx context.Context, startDate, endDate *time.Time) (
        []PerformaStudioDetailItemResponse, 
        TotalPerformaAccount, 
        error,
    )
    
    // Calculate performance for all hosts
    CalculateByHosts(ctx context.Context, startDate, endDate *time.Time) (
        []*PerformaHostSummaryResponse, 
        error,
    )
    
    // Calculate performance for a specific host
    CalculateByHost(ctx context.Context, hostID string, startDate, endDate *time.Time) (
        PerformaHostDetailResponse, 
        error,
    )
    
    // Calculate performance for a specific studio
    CalculateByStudio(ctx context.Context, studioID string, startDate, endDate *time.Time) (
        []PerformaStudioDetailItemResponse, 
        TotalPerformaAccount, 
        error,
    )
}
```

### Data Types

#### TotalPerformaHost

| Field | Type | Description |
|-------|------|-------------|
| `Duration` | `int64` | Total live duration in minutes |
| `GMVSales` | `int64` | Gross Merchandise Value (Sales) |
| `GMVPaid` | `int64` | Confirmed/Paid GMV |
| `AvgSales` | `int64` | Average sales per session |
| `AvgPaid` | `int64` | Average paid amount per session |

#### TotalPerformaAccount

| Field | Type | Description |
|-------|------|-------------|
| `GMV` | `int64` | Total Gross Merchandise Value |
| `Ads` | `int64` | Total advertising spend |
| `CommissionTotal` | `int64` | Total commission earned |
| `CommissionPaid` | `int64` | Commission already paid |
| `CommissionPending` | `int64` | Commission pending payment |
| `Income` | `int64` | Net income (GMV - Ads) |

### Usage

The Performa Aggregator is used by:
- `/api/performa/host` - Host performance endpoints
- `/api/performa/studio` - Studio performance endpoints
- `/api/performa/account` - Account performance endpoints
- `/api/dashboard` - Dashboard summary data

---

## 🔗 Dependencies

- Transaction domain for sales data
- Account domain for account information
- Host domain for host data
- Studio domain for studio data

## 📌 Notes

- Aggregators are injected via Uber FX dependency injection
- All calculations support date range filtering
- Results are cached where appropriate for performance
