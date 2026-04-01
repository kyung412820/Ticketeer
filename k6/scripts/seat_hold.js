import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const EVENT_ID = Number(__ENV.EVENT_ID || 1);
const SEAT_ID = Number(__ENV.SEAT_ID || 1);
const PROVIDED_QUEUE_TOKEN = __ENV.QUEUE_TOKEN || '';

export const options = {
  scenarios: {
    same_seat_hold_contention: {
      executor: 'constant-vus',
      vus: 20,
      duration: '5s',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<1000'],
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

export default function () {
  const clientId = `hold-client-${__VU}-${__ITER}-${Date.now()}`;
  const queueToken = PROVIDED_QUEUE_TOKEN || enterQueue(clientId);

  if (!queueToken) {
    sleep(0.2);
    return;
  }

  const payload = JSON.stringify({
    event_id: EVENT_ID,
    client_id: clientId,
    queue_token: queueToken,
  });

  const res = http.post(`${BASE_URL}/api/seats/${SEAT_ID}/hold`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'hold status is 200 or 400/404': (r) => [200, 400, 404].includes(r.status),
  });

  sleep(0.2);
}
