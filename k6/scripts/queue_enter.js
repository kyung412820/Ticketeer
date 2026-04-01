import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const EVENT_ID = Number(__ENV.EVENT_ID || 1);

export const options = {
  scenarios: {
    burst_queue_enter: {
      executor: 'constant-vus',
      vus: 50,
      duration: '10s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<1000'],
  },
};

export default function () {
  const clientId = `k6-client-${__VU}-${__ITER}-${Date.now()}`;

  const payload = JSON.stringify({
    event_id: EVENT_ID,
    client_id: clientId,
  });

  const res = http.post(`${BASE_URL}/api/queue/enter`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'queue enter status is 200 or 400': (r) => r.status === 200 || r.status === 400,
  });

  sleep(0.2);
}
