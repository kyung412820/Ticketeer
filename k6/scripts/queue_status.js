import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const EVENT_ID = Number(__ENV.EVENT_ID || 1);

export const options = {
  scenarios: {
    polling_status: {
      executor: 'constant-vus',
      vus: 30,
      duration: '10s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<800'],
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
  const clientId = `status-client-${__VU}-${Date.now()}`;
  const queueToken = enterQueue(clientId);

  if (!queueToken) {
    sleep(0.2);
    return;
  }

  for (let i = 0; i < 5; i++) {
    const res = http.get(`${BASE_URL}/api/queue/status/${queueToken}`);
    check(res, {
      'queue status is 200': (r) => r.status === 200,
    });
    sleep(0.5);
  }
}
