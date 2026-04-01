# Frontend Performance Patch

## Added / Changed
- `SeatSectionGrid` memoization
- section / seat card component memoization
- stable callbacks in booking page
- derived state memoization
- `SeatSelectionSummary` memoization
- `FRONTEND_PERFORMANCE_NOTES.md`

## Why
좌석 수가 적을 때는 티가 덜 나지만,
티켓팅 UI는 좌석이 많아질수록 render churn을 줄이는 구조가 중요합니다.

## Portfolio Point
- 프론트도 성능 최적화를 고려
- memoization과 stable callback으로 좌석 UI re-render 비용 감소
