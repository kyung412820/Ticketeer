# Frontend Test Strategy Update

## Added / Changed
- Playwright selector stabilization
- stable test ids for error/loading state
- interaction state E2E test
- frontend test strategy doc

## Why
프론트도 단순 렌더링이 아니라,
실제 브라우저 기준으로 핵심 흐름이 검증된다는 점을 보여주기 위함입니다.

## README Paragraph
### 프론트 테스트 전략
Playwright 기반 smoke / E2E 테스트를 추가해 실제 브라우저 기준으로 핵심 흐름을 검증했습니다.

- 공연 목록 렌더링 및 상세 이동
- queue 화면 상태 렌더링
- booking 직접 진입 차단
- 잘못된 진입 시 버튼 disabled 상태 확인
- alert / live region / 버튼 노출 확인

또한 Next 내부 route announcer와의 selector 충돌을 피하기 위해
`error-banner`, `loading-state`에 `data-testid`를 추가해 테스트 안정성을 높였습니다.
