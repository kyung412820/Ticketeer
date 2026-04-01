# Frontend React Query Adoption

## What changed
이번 단계에서는 문서화만 하던 상태에서 한 단계 더 나아가,
React Query를 실제로 도입했습니다.

적용 대상:
- queue status polling
- seats fetch / refetch

## Why these first
티켓팅 프론트에서 가장 서버 상태 성격이 강한 값은 다음 두 가지입니다.

- queue status
- seats list

둘 다
- refetch가 자주 필요하고
- stale state를 줄여야 하며
- 상태 변화가 UI에 직접 연결됩니다.

## What improved
- queue polling을 React Query 기준으로 관리
- READY / EXPIRED 시 refetch interval 자동 중단
- 좌석 목록 캐시 키 분리
- hold / booking 후 invalidateQueries로 재동기화
- 페이지 로컬 fetch 코드를 일부 줄이고 서버 상태 관리 책임을 React Query로 이동

## Why not everything
현재 단계에서는 공연 목록/상세까지 한 번에 React Query로 옮기지 않았습니다.

이유:
- `/events`, `/events/[id]`는 서버 중심 렌더링 성격이 강함
- 먼저 클라이언트 상호작용이 많은 queue / booking 영역부터 적용하는 것이 자연스러움

## Summary
이번 단계는
"React Query를 쓸 수 있다" 수준이 아니라,
**어디에 먼저 도입해야 가장 효과적인지 판단하고, 실제로 적용한 단계**입니다.
