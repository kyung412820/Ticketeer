# Frontend TODO 4 - 상태 저장 전략 정리

## What changed
- client_id는 localStorage에서 생성/재사용
- queue_token은 eventId 기준으로 sessionStorage에 저장
- queue 페이지 재방문 시 기존 queue_token 재사용
- queue EXPIRED / 조회 에러 / 예매 완료 시 queue_token 정리
- booking 페이지는 query param이 없어도 storage 기반으로 진입 가능

## Why
티켓팅 흐름에서는 queueToken, clientId 같은 상태를
어디에 저장하고 언제 정리할지 기준이 분명해야 재진입/새로고침 시 덜 꼬입니다.

## Next
- 예매 완료 UX 개선
- 좌석 section/grouping 개선
- 컴포넌트 구조 정리
