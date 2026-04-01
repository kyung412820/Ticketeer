# TICKETEER k6 Scripts

## Prerequisites
- Backend server running at `http://localhost:8080`
- PostgreSQL and Redis running
- At least one OPEN event in DB

## Install k6
See: https://k6.io/docs/get-started/installation/

## Run examples

### Queue enter load test
```bash
k6 run scripts/queue_enter.js
```

### Seat hold load test
Before running this, update `EVENT_ID`, `SEAT_ID`, and optionally `QUEUE_TOKEN`.
```bash
k6 run -e EVENT_ID=1 -e SEAT_ID=1 -e QUEUE_TOKEN=your_token_here scripts/seat_hold.js
```

### Booking load test
Before running this, make sure the seat is already held by the same client and token.
```bash
k6 run -e EVENT_ID=1 -e SEAT_ID=1 -e QUEUE_TOKEN=your_token_here scripts/booking.js
```

## Notes
- `seat_hold.js` and `booking.js` are most useful for targeted contention tests.
- `queue_enter.js` is useful for burst traffic simulation.
- For real measurement, keep backend logs and DB state clean before each run.
