# Frontend Storage / Security Strategy Patch

## Added
- `FRONTEND_STORAGE_SECURITY.md`

## Why
프론트에서 값을 저장했다는 사실보다,
- 왜 localStorage인지
- 왜 sessionStorage인지
- 왜 query param을 기본 수단으로 삼지 않는지
- 언제 토큰을 정리하는지
를 설명할 수 있어야 설계력이 보입니다.

## Portfolio Point
- 프론트는 값 자체를 신뢰하지 않고 UX 흐름만 유지
- 실제 보안/정합성 검증은 백엔드가 담당
- 저장 전략도 세션 범위와 민감도를 기준으로 나눠 설계
