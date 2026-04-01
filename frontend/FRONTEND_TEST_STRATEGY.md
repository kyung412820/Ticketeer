# Frontend Test Strategy

## Goal
프론트도 "화면이 보인다" 수준에서 끝내지 않고,
핵심 사용자 흐름이 실제 브라우저 기준으로 깨지지 않는지 자동화 테스트로 검증합니다.

## Why Playwright
- 실제 브라우저 기준 검증 가능
- 라우팅 / 렌더링 / 버튼 상태 / 접근성 역할을 함께 확인 가능
- 티켓팅처럼 사용자 흐름이 중요한 UI에 적합

## Current Coverage
- 공연 목록 진입
- 상세 이동
- queue 화면 상태 확인
- booking 직접 진입 차단
- 잘못된 직접 진입 시 버튼 disabled 확인
- 접근성용 alert / live region / 버튼 노출 확인

## Selector Strategy
- Next 내부 announcer와 충돌하지 않도록
  - `data-testid="error-banner"`
  - `data-testid="loading-state"`
  를 추가해 안정적인 locator를 사용합니다.

## Next Candidates
- 예매 완료 카드 렌더링 검증
- 정상 queue -> booking 흐름 E2E
- API mock 기반 deterministic 테스트
