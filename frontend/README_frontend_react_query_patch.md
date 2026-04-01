# Frontend React Query Patch

## Added
- `@tanstack/react-query`
- QueryProvider
- queue status query hook
- seats query hook
- React Query adoption note

## Why
문서화만이 아니라 실제로
서버 상태 관리 책임 일부를 React Query로 이동했습니다.

## Portfolio Point
- queue status polling을 React Query로 관리
- seats refetch를 invalidateQueries 기반으로 정리
- SSR/CSR 전략과 맞물리게 도입 범위를 선택
