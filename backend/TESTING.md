# Backend API Tests

## Run
```bash
cd backend
go test ./tests -v
```

## Covered
- happy path booking flow
- booking not open
- seat already booked
- booking without hold
- invalid queue token
- seat already held
- hold missing on booking
- concurrent booking on same seat: exactly one success
- concurrent hold on same seat: exactly one success
