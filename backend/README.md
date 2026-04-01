# TICKETEER Backend with Booking API

## Run

```bash
cd backend
go mod tidy
go run ./cmd/api
```

## Seed

Apply `seed.sql` to your PostgreSQL database after the tables are auto-created.

## Endpoints

- `GET /api/health`
- `GET /api/events`
- `GET /api/events/:id`
- `GET /api/events/:id/seats`
- `POST /api/queue/enter`
- `GET /api/queue/status/:token`
- `POST /api/seats/:seatId/hold`
- `POST /api/bookings`
