# Frontend Retry / Timeout Patch

## Added
- `FRONTEND_RETRY_TIMEOUT_POLICY.md`
- `api-client.ts` timeout / limited retry support

## Why
프론트도 네트워크 실패 대응 전략이 필요합니다.
하지만 티켓팅 도메인에서는 모든 요청을 자동 재시도하면 안 됩니다.

## Portfolio Point
- 조회성 요청과 mutation 요청의 retry 정책을 분리
- 예매 확정은 자동 재시도보다 UX 신뢰와 중복 행동 방지를 우선
