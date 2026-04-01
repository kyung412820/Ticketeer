INSERT INTO events
(title, description, venue, event_at, booking_open_at, booking_close_at, status, created_at, updated_at)
VALUES
(
  'Coldplay Live in Seoul',
  'World tour in Seoul',
  'Olympic Stadium',
  '2026-07-01 19:00:00',
  '2026-06-01 20:00:00',
  '2026-06-30 23:59:59',
  'OPEN',
  NOW(),
  NOW()
),
(
  'Jazz Night',
  'Special live jazz performance',
  'Blue Square',
  '2026-08-10 18:00:00',
  '2026-07-15 20:00:00',
  '2026-08-09 23:59:59',
  'OPEN_PENDING',
  NOW(),
  NOW()
);

INSERT INTO seats
(event_id, seat_no, section, price, status, created_at, updated_at)
VALUES
(1, 'A-1', 'R', 150000, 'AVAILABLE', NOW(), NOW()),
(1, 'A-2', 'R', 150000, 'AVAILABLE', NOW(), NOW()),
(1, 'A-3', 'R', 150000, 'BOOKED', NOW(), NOW()),
(1, 'A-4', 'R', 150000, 'AVAILABLE', NOW(), NOW()),
(1, 'B-1', 'S', 120000, 'AVAILABLE', NOW(), NOW()),
(1, 'B-2', 'S', 120000, 'AVAILABLE', NOW(), NOW()),
(1, 'B-3', 'S', 120000, 'BOOKED', NOW(), NOW());
