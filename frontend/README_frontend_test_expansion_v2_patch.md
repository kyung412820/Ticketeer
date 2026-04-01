# Frontend Test Expansion V2 Patch

## Added / Changed
- stable test ids for error/loading state
- stricter and more reliable Playwright selectors
- new interaction state E2E test

## Why
테스트도 유지보수성을 고려해야 합니다.
Next 내부 announcer 같은 요소와 충돌하지 않도록
selector 전략을 보강했습니다.

## Portfolio Point
- 프론트 테스트도 "동작만 함"이 아니라
  안정적인 selector 전략까지 고려
