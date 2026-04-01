# Frontend Test Expansion V2

## What changed
- `error-banner`, `loading-state`에 test id 추가
- Playwright locator를 더 안정적으로 수정
- invalid direct entry 테스트를 strict-mode friendly 하게 수정
- interaction state 검증 테스트 추가
  - 잘못된 진입 시 버튼 disabled 확인
  - queue 페이지에서 loading / error / queue content 중 하나가 보이는지 확인

## Why
프론트 테스트는 셀렉터가 애매하면 UI가 맞아도 쉽게 깨집니다.
그래서 접근성 role만 믿기보다,
테스트 목적에 맞는 안정적인 selector를 일부 명시적으로 추가했습니다.

## Summary
이번 단계는 단순 smoke 테스트 추가를 넘어서,
**실제 테스트가 덜 깨지고 의미 있는 상태를 검증하도록 보강한 단계**입니다.
