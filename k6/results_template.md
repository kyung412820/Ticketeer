# k6 Result Template

## Test target
- API:
- Scenario:
- Date:

## Environment
- Backend version:
- Database:
- Cache:
- Local / Docker / Server:

## k6 options
- VUs:
- Duration / Iterations:
- Thresholds:

## Result summary
- Avg response time:
- p95 response time:
- Failure rate:
- Success rate:

## Observations
- Did duplicate booking occur?
- Did duplicate hold occur?
- Any unexpected 5xx errors?
- Any queue bottlenecks?

## Portfolio summary sentence
- Example:
  - `k6로 동일 좌석 동시 booking 요청을 검증한 결과, 20개 경쟁 요청 중 1건만 성공하고 나머지는 정상적으로 실패 처리되어 중복 예매가 발생하지 않았다.`
