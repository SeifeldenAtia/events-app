# Events Ticket Booking Service

A REST API (Go + Gin + SQLite) that lets users book tickets for events and lets owners manage events.

Author: Seifelden Atia

## How to run

```bash
go run main.go
```

Go will fetch the dependencies listed in `go.mod`/`go.sum` automatically on first run. The server starts on `http://localhost:8080`, and an `api.db` SQLite file is created automatically on first run (using the pure-Go [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) driver, so no C compiler/CGO is required).

## Example requests

Ready-to-use requests are in [api-test/](api-test/) (REST Client `.http` files). A few examples:

**Create an event** (owner)

```http
POST http://localhost:8080/api/v1/events
content-type: application/json
X-Owner-Id: 1

{
  "name": "Test event",
  "description": "A test event",
  "location": "A test location",
  "dateTime": "2025-01-01T15:30:00.000Z",
  "availableTickets": 100
}
```

**See all available events**

```http
GET http://localhost:8080/api/v1/events
```

**Book a ticket**

```http
POST http://localhost:8080/api/v1/events/1/bookings
content-type: application/json
X-User-Id: 42

{
  "quantity": 2
}
```

**See my booked tickets**

```http
GET http://localhost:8080/api/v1/bookings
X-User-Id: 42
```

**See tickets booked for an event** (owner)

```http
GET http://localhost:8080/api/v1/events/1/bookings/count
X-Owner-Id: 1
```

## Additional features

- `PUT /events/:id` to update an event.
- Bookings take a `quantity`; overselling is prevented atomically.
- Owner actions check `X-Owner-Id` against the event's real owner (`403` otherwise).
- Validation errors include field-level detail. Ready-to-run requests in [api-test/](api-test/).

## Assumptions

- No user management/auth: a user is identified via the `X-User-Id` header, and an event owner via `X-Owner-Id`. Both are just integers, the owner check above only guards against using the _wrong_ id, not against someone guessing a valid one.
- Deleting an event with existing bookings is blocked (`409`).

Future improvements are marked with `TODO` comments in the code (e.g. JWT auth, pagination, API versioning, Rate Limiting).
