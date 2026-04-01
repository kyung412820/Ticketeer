# Frontend TODO 2 - 에러 메시지 개선

## What changed
- 백엔드 에러 코드별 사용자 메시지 매핑 추가
- 네트워크 실패 시 기술 문구 대신 안내 메시지 출력
- queue / booking 페이지에서 공통 에러 메시지 사용
- hold / booking 실패 후 좌석 목록 재조회 유지

## Why
티켓팅 UI에서는 `Failed to fetch` 같은 기술 문구보다,
사용자 행동 기준으로 이해 가능한 메시지가 중요합니다.

## Next
- polling cleanup 추가 보강
- 좌석 상태 재동기화 세분화
- 상태 저장 전략 정리
