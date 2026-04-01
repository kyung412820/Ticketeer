import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const EVENT_ID = Number(__ENV.EVENT_ID || 1);
const SEAT_ID = Number(__ENV.SEAT_ID || 1);
const PROVIDED_QUEUE_TOKEN = __ENV.QUEUE_TOKEN || '';

export const options = {
  scenarios: {
    same_seat_booking_contention: {
      executor: 'shared-iterations',
      vus: 20,
      iterations: 20,
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<1500'],
  },
};

function enterQueue(clientId) {
  const payload = JSON.stringify({
    event_id: EVENT_ID,
    client_id: clientId,
  });

  const res = http.post(`${BASE_URL}/api/queue/enter`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  if (res.status !== 200) {
    return '';
  }

  try {
    return res.json('data.queue_token');
  } catch (_) {
    return '';
  }
}

function holdSeat(clientId, queueToken) {
  const payload = JSON.stringify({
    event_id: EVENT_ID,
    client_id: clientId,
    queue_token: queueToken,
  });

  return http.post(`${BASE_URL}/api/seats/${SEAT_ID}/hold`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });
}

export default function () {
  const clientId = `booking-client-${__VU}-${__ITER}-${Date.now()}`;
  const queueToken = PROVIDED_QUEUE_TOKEN || enterQueue(clientId);

  if (!queueToken) {
    sleep(0.1);
    return;
  }

  if (!PROVIDED_QUEUE_TOKEN) {
    const holdRes = holdSeat(clientId, queueToken);
    if (holdRes.status !== 200) {
      sleep(0.1);
      return;
    }
  }

  const payload = JSON.stringify({
    event_id: EVENT_ID,
    seat_id: SEAT_ID,
    client_id: clientId,
    queue_token: queueToken,
  });

  const res = http.post(`${BASE_URL}/api/bookings`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'booking status is 200 or 400/404': (r) => [200, 400, 404].includes(r.status),
  });

  sleep(0.1);
}
