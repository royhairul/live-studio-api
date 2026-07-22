# 🔗 Host Account Domain

The **Host Account** domain records which Shopee accounts belong to which host, and
for how long.

## 📋 Overview

This is the bridge the performa endpoints use to answer *whose numbers are these?*

Rows in `lives` belong to an **account**, not to a host. Something has to connect the
two. That used to be attendance — `attendance → account_sessions` — but attendance
records **who was present**, which is a different question from **who owns an
account's results**. A host who forgets to check in still owns their accounts'
sessions, and under the old scheme those sessions silently vanished from the report.

This module answers the ownership question directly, and does it *historically*.

## 📁 Structure

```
hostaccount/
├── controller/
│   ├── hostaccount_controller.go       # Controller interface
│   └── hostaccount_controller_impl.go  # Controller implementation
├── entity/
│   └── hostaccount_entity.go           # Assignment database model
├── params/
│   ├── hostaccount_filter.go           # Query filter
│   ├── hostaccount_request.go          # Request DTOs
│   └── hostaccount_response.go         # Response DTOs + HostAccountSpan
├── repository/
│   ├── hostaccount_repository.go       # Repository interface
│   └── hostaccount_repository_impl.go  # Repository implementation
├── service/
│   ├── hostaccount_service.go          # Service interface
│   └── hostaccount_service_impl.go     # Service implementation
├── route.go                            # Route definitions
├── module.go                           # FX module
└── README.md                           # This file
```

## 🗃️ Entity

```go
type HostAccount struct {
    gorm.Model

    HostID    uuid.UUID  // FK to Host
    AccountID uint       // FK to Account

    ValidFrom time.Time  // inclusive
    ValidTo   *time.Time // exclusive; nil = still running

    Note string
    tenantdb.TenantBase
}
```

Studio is **not** stored here — it resolves through the account, the same way the Live
and Transaction domains do it.

### Why versioned instead of a column on `accounts`

A plain `accounts.host_id` would be simpler, but it only knows the *present*. Move an
account to a new host and every past report silently rewrites itself, because the old
numbers would follow the new owner.

`ValidFrom`/`ValidTo` keep each period intact, so a report over July still resolves to
whoever held the account in July. That is also why **ending an assignment and deleting
one are different operations** — see the API section.

### One account, one host at a time

The schema permits overlapping rows; the service rejects them. A live session cannot
belong to two hosts at once, so the ambiguity is blocked at write time rather than
guessed at read time:

```
400 {"error":"account already assigned",
     "details":"account is held by host Budi for an overlapping period"}
```

The same host holding an account across adjacent spans is fine — only a *different*
host overlapping is refused.

## 🌐 API Endpoints

Base path: `/api/host-account` · Auth: Bearer (`superadmin`, `admin`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/host-account` | List assignments; filter by `host`, `account`, `studio`, `active=true` |
| `POST` | `/host-account` | Assign an account to a host |
| `GET` | `/host-account/:id` | One assignment |
| `PUT` | `/host-account/:id` | Patch an assignment — typically to end it |
| `DELETE` | `/host-account/:id` | Remove the assignment entirely |

### Ending ≠ deleting

|  | Ending (`PUT valid_to`) | Deleting (`DELETE`) |
|--|--|--|
| Row | kept | removed |
| Sessions before the end date | still attributed to the host | no longer attributed |
| Past reports | unchanged | **change retroactively** |

Ending is nearly always what you want. Deleting is for correcting an assignment that
should never have existed.

### Dates accept two formats

`valid_from` and `valid_to` take either `YYYY-MM-DD` or a full RFC3339 timestamp.

The timestamp form is not decoration. To end an assignment *now*, a plain date will
not do: today's date parses to **midnight**, which can fall before the assignment
started, and the range check rejects it:

```
400 {"error":"invalid date range","details":"valid_to must be after valid_from"}
```

So the UI sends `new Date().toISOString()` when the operator presses "Akhiri".

### Two defaults that exist to prevent surprises

- **`valid_from` defaults to midnight today**, not `time.Now()`. Assigning an account
  at 3pm still counts that morning's sessions, which is what assigning it in the UI is
  taken to mean.
- **Same-day handovers start where the last one ended.** Because of the default above,
  assigning an account that another host released at 10am would otherwise overlap
  midnight–10am and be refused. When `valid_from` is omitted, the service starts the
  new assignment at the moment of release instead.

## 📝 Request/Response Examples

### Assign an account

**Request:**
```http
POST /api/host-account
Authorization: Bearer <jwt_token>

{ "host_id": "b15d206d-8ebf-4cae-b216-7f9e9b6fdc02", "account_id": 6 }
```

**Response:** `201`
```json
{
  "message": "created host account successfully",
  "data": {
    "id": 1,
    "host_id": "b15d206d-8ebf-4cae-b216-7f9e9b6fdc02",
    "host_name": "Ayu",
    "account_id": 6,
    "account_name": "Toko Meycan",
    "studio_id": 11,
    "studio_name": "Studio A",
    "valid_from": "2026-07-22T00:00:00+07:00",
    "valid_to": null,
    "is_active": true,
    "note": ""
  }
}
```

### End an assignment

```http
PUT /api/host-account/1
Authorization: Bearer <jwt_token>

{ "valid_to": "2026-07-22T15:04:05Z" }
```

Sending `{"valid_to": ""}` clears it again, reopening the assignment. Omitting the
field leaves it untouched.

### List what a host currently holds

```bash
curl "/api/host-account?host=b15d206d-8ebf-4cae-b216-7f9e9b6fdc02&active=true" \
  -H "Authorization: Bearer $TOKEN"
```

## 🧮 How reporting consumes this

`PerformaAggregator` never queries this table directly. It calls
`AccountIDsForHost(ctx, hostID, start, end)`, which returns one
`params.HostAccountSpan` per assignment overlapping the report window, **clamped** to
that window:

```go
type HostAccountSpan struct {
    AccountID uint
    From      time.Time
    To        time.Time
}
```

Clamping is what stops an assignment that began before the range — or ended inside it —
from dragging in sessions it never covered. The aggregator then keeps a session only if
its `start_time` falls inside one of that account's spans, deduplicating by
`session_id` so overlapping spans cannot double-count.

A host with no assignment reports **zeros, not an error**. There is deliberately no
fallback to attendance: a silent fallback would reintroduce the dependency this domain
exists to remove, and would make the origin of a number impossible to reason about.

## 🔗 Dependencies

- Host domain — the assignee
- Account domain — the assigned account, and the studio behind it
- Consumed by `internal/aggregator/performa`

## 📌 Notes

- Assignments are tenant-scoped like every other domain.
- `valid_to` is **exclusive**; `valid_from` is inclusive.
- Deleting a host or account cascades to its assignments (`OnDelete:CASCADE`).
- Filling this table is a prerequisite for host performa figures — until an account is
  assigned, `/performa/host/{id}` legitimately reports zeros.
