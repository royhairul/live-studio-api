# 📺 Live Domain

The **Live** domain manages live streaming sessions for the Live Studio API.

## 📋 Overview

This module provides:
- Live session data from Shopee
- Live session details with metrics
- Integration with Shopee Live API

## 📁 Structure

```
live/
├── controller/
│   ├── live_controller.go       # Controller interface
│   └── live_controller_impl.go  # Controller implementation
├── entity/
│   └── live_entity.go           # Live session database model
├── params/
│   ├── request.go               # Request DTOs
│   └── response.go              # Response DTOs
├── repository/
│   ├── live_repository.go       # Repository interface
│   └── live_repository_impl.go  # Repository implementation
├── service/
│   ├── live_service.go          # Service interface
│   └── live_service_impl.go     # Service implementation
├── route.go                     # Route definitions
├── module.go                    # FX module
└── README.md                    # This file
```

## 🗃️ Entity

Persisted by the sync endpoints. Columns mirror exactly what Shopee's
`liveList/v2` returns — there is **no end time**, only `StartTime` + `Duration`.

```go
type Live struct {
    ID        int64      // Primary key (snowflake)
    SessionID int64      // Shopee live session ID (unique — the upsert key)
    Title     string
    CoverImage string
    Status    int        // 2 = ended
    StartTime *time.Time // converted from Shopee epoch milliseconds
    Duration  int64      // milliseconds

    // Audience
    Views, Viewers, PeakViews, AvgViewsDuration int
    Comments, Likes, FollowersGrowth            int
    EngagedUV, AvgEngagedCCU, ThirtyMinsCount   int

    // Commerce
    Atc, ProductClicks                 int
    ConversionRate                     float64
    PlacedOrders, PlacedItemSold       int
    PlacedSales                        float64
    ConfirmedOrders, ConfirmedItemSold int
    ConfirmedSales                     float64
    PaidOrders                         int
    PaidSales                          float64

    AccountID uint    // FK to Account; studio is reached via Account.Studio
    tenantdb.TenantBase
}
```

**Two behaviours worth knowing:**
- Shopee sends JSON `null` for `views`, `peakViews`, `likes`, `followersGrowth`,
  `productClicks`, `conversionRate`, `paidOrders` and `paidSales`. These persist as
  `0` — a `0` means *not reported*, not necessarily *zero*.
- Studio is **not** denormalised; it resolves through the account. Moving an account
  to a different studio therefore also moves its historical sessions. This matches
  how the Transaction domain already behaves.

## 🌐 API Endpoints

Base path: `/api/live`

The domain splits cleanly in two: **preview** is live and comes from Shopee over a
WebSocket; **history** is past and comes from the database.

| Method | Endpoint | Source | Auth |
|--------|----------|--------|------|
| `GET` | `/preview` | **WebSocket.** Shopee — sessions on air now, pushed every 2–10s | `?token=` |
| `GET` | `/preview/detail/:id/:sessionId` | **WebSocket.** Shopee — one session's dashboard, every 3–8s | `?token=` |
| `GET` | `/history` | **Database.** Stored sessions, all accounts | Bearer |
| `GET` | `/history/:id` | **Database.** Stored sessions, one account | Bearer |
| `POST` | `/history/sync` | Shopee → database, every account | Bearer |
| `POST` | `/history/:id/sync` | Shopee → database, one account | Bearer |

> ⚠️ The two `/preview` routes are **WebSocket** endpoints, not plain REST. They
> authenticate via a `?token=` query parameter because browsers cannot set headers on a
> WebSocket handshake, and the server never sends pings — a client that stays silent for
> 60s is disconnected, so send a keepalive message every ~30s.

### Preview shows only sessions that are still on air

Shopee's `realtime/sessionList` also returns sessions that have already finished, and its
`status` field reads `2` for both ongoing and ended ones — so it cannot be used to tell them
apart. What does distinguish them is that **`duration` advances while a stream runs**, so
`startTime + duration` tracks the present until the stream ends and then freezes.
`utils.IsLive` uses that, with a 5-minute tolerance for polling lag.

For this reason Shopee's `duration` must **never** be overwritten with
`now - startTime` — doing so makes every finished session look like it is still on air.

In the response, `total` counts everything Shopee returned and `relive` counts what is
actually live:

```jsonc
{"id": "8", "name": "nalanavfsh", "total": 1, "relive": 0, "reportLive": []}
// Shopee returned 1 session; it had already ended, so preview reports nothing live.
```

### History comes from the database

`GET /history` reads the `lives` table and **never contacts Shopee**, so it keeps working
when an account's cookie expires. Filter with `account`, `studio`, `startDate`/`endDate`
(which must be sent together), and paginate with `page`/`pageSize`.

```bash
curl "/api/live/history?studio=11&startDate=2025-08-01&endDate=2025-08-31" -H "Authorization: Bearer $TOKEN"
```

### Syncing Shopee → database

The sync endpoints are **idempotent**: sessions are matched on `session_id` and upserted,
so re-running refreshes figures that settle after a stream ends (notably `confirmedSales`)
instead of creating duplicates.

```bash
curl -X POST /api/live/history/6/sync -H "Authorization: Bearer $TOKEN"
# -> {"account_id":"6","account_name":"Toko Meycan","fetched":3,"created":3,"updated":0}
# re-run:
# -> {"account_id":"6","account_name":"Toko Meycan","fetched":3,"created":0,"updated":3}
```

Only sessions inside the `timeDim` window ending at `endDate` (default: today) are synced,
so **backfilling older history means calling sync repeatedly with earlier `endDate` values**:

```bash
curl -X POST "/api/live/history/6/sync?timeDim=1m&endDate=2025-08-31" -H "Authorization: Bearer $TOKEN"
```

There is deliberately **no scheduler** — trigger it from a UI button or manually.

## 📝 Request/Response Examples

### Preview (WebSocket)

```
ws://localhost:8080/api/live/preview?token=<jwt>
```

```json
[
  { "id": "6", "name": "Toko Meycan", "total": 0, "relive": 0, "reportLive": [] },
  { "id": "8", "name": "nalanavfsh", "total": 1, "relive": 0, "reportLive": [] }
]
```

`total` is what Shopee returned, `relive` is what is genuinely on air. A live session
appears in `reportLive` carrying `isLive: true` and a derived `omsetPerHour`.

### Stored history

**Request:**
```http
GET /api/live/history?pageSize=2&studio=11
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "retrieved live history successfully",
  "data": {
    "page": 1,
    "pageSize": 2,
    "total": 10,
    "totalPage": 5,
    "history": [
      {
        "session_id": 216553201,
        "title": "AMBYAR KASUR BUSA",
        "status": 2,
        "start_time": "2026-07-03T08:58:39.634+07:00",
        "duration": 17602747,
        "viewers": 120,
        "comments": 59,
        "atc": 9,
        "placed_orders": 1,
        "placed_sales": 83000,
        "confirmed_orders": 1,
        "confirmed_sales": 83000,
        "omset_per_hour": 16974.62,
        "account_id": 6,
        "account_name": "Toko Meycan",
        "studio_id": 11,
        "studio_name": "Studio A",
        "synced_at": "2026-07-22T07:09:12.506+07:00"
      }
    ]
  }
}
```

`synced_at` is when the row was last refreshed by a sync — useful for spotting stale data,
since nothing here auto-updates.

## 🔗 Dependencies

- Account domain — supplies the Shopee cookie and, via `Account.Studio`, the studio
- Shopee client (`internal/clients/shopee`) — `realtime/sessionList` for preview,
  `liveList/v2` for history

## 📌 Notes

- Preview data is never persisted; only the sync endpoints write to `lives`.
- History is only as fresh as the last sync — there is no scheduler by design.
- Shopee's field names differ between the two endpoints for the same concept
  (`peakViewers` in realtime vs `peakViews` in history). Both are mapped as Shopee sends
  them; do not assume they match.
