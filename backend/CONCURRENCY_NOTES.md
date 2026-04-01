# Concurrency hardening patch

## What changed
- `bookings.seat_id` unique constraint added
- `seats(event_id, seat_no)` composite unique constraint added
- booking confirmation now uses:
  - PostgreSQL transaction
  - row-level lock (`FOR UPDATE`) on seat row
  - seat status update + booking insert + queue state update in one transaction
- Redis hold key is deleted only after transaction commit succeeds

## Why
This makes the final booking step much closer to a real ticketing system:
- fast temporary hold in Redis
- final consistency guaranteed in PostgreSQL
- duplicate booking prevention at both application and DB levels

## Next
- add concurrent booking tests with goroutines
- add k6 load scripts
