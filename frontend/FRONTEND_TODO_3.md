# Frontend TODO 3 - polling / 상태 재동기화 보강

## What changed
- queue polling이 READY / EXPIRED / 에러 시 즉시 중단되도록 보강
- booking 페이지에서 queueToken / clientId 없는 잘못된 접근 사전 차단
- 좌석 재조회 시 최신 요청만 반영하도록 request sequence 적용
- 좌석 재조회 후 선택 좌석 상태를 최신 서버 상태로 동기화
- 예매 버튼은 선택 좌석이 AVAILABLE일 때만 활성화

## Why
티켓팅 UI에서는 화면에 보이는 상태와 실제 서버 상태가 쉽게 어긋날 수 있습니다.
따라서 polling 종료 조건, 최신 응답만 반영, 좌석 상태 재동기화가 중요합니다.

## Next
- 상태 저장 전략 정리
- 예매 완료 UX 개선
- 좌석 section/grouping 개선
