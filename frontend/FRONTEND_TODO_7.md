# Frontend TODO 7 - 컴포넌트 구조 정리 / 최종 UI polish

## What changed
- 공통 UI 컴포넌트 분리
  - `ErrorBanner`
  - `LoadingState`
  - `BookingResultCard`
  - `SeatSelectionSummary`
- booking 페이지 책임 분리
- 예매 완료 상태와 일반 에러 상태를 분리해서 표시
- 최종 UI 문구와 영역 구성을 정리

## Why
기능이 늘어날수록 한 페이지에 모든 UI를 몰아넣으면
읽기 어렵고 수정 비용이 커집니다.
핵심 영역을 컴포넌트로 분리해 유지보수성과 가독성을 높였습니다.

## Next
- docker-compose
- Kubernetes
- GitHub Actions
