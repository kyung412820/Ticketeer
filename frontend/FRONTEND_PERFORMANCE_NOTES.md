# Frontend Performance Optimization Notes

## Goal
좌석 수가 적을 때는 체감이 적지만,
티켓팅 UI는 좌석 수가 많아질수록 렌더링 비용과 불필요한 re-render를 고민해야 합니다.

이번 단계에서는 지금 구조에서 효과가 큰 최적화를 먼저 적용했습니다.

---

## What changed

### 1. Seat grid memoization
- `SeatSectionGrid`를 `memo` 적용
- section 단위 컴포넌트 분리
- seat card 단위까지 `memo` 적용
- section grouping은 `useMemo`로 계산

### 2. Stable callbacks
- `handleSelectSeat`
- `handleHold`
- `handleBooking`
- `refetchSeats`

를 `useCallback`으로 고정해
하위 컴포넌트에 불필요한 함수 reference 변경이 덜 전파되게 했습니다.

### 3. Derived state memoization
- `selectionDisabled`
- 좌석 결과 요약 값
- hold 만료 표시

같은 파생값은 `useMemo`로 관리했습니다.

### 4. Summary component memoization
- `SeatSelectionSummary`도 `memo` 적용

---

## Why these first

지금 단계에서 가장 현실적인 최적화는
"구조를 크게 바꾸지 않으면서 render churn을 줄이는 것"입니다.

특히 좌석 UI는
- 좌석 수가 늘어날 수 있고
- 선택 상태가 자주 바뀌며
- hold / booking 이후 refetch가 발생합니다.

그래서 section/card 단위 memoization이 효과가 큽니다.

---

## What we did NOT add yet

### Virtualization
아직 좌석 수가 아주 많지 않은 MVP 기준이라 바로 도입하지 않았습니다.

도입 후보 상황:
- 500석 이상
- section당 카드 수가 매우 많을 때
- 스크롤 성능 저하가 체감될 때

### Heavy client cache tuning
React Query를 도입했지만,
지금은 기본적인 staleTime / invalidate 정도만 적용했습니다.
대규모 화면이 되면 query cache 정책을 더 세분화할 수 있습니다.

---

## Summary
이번 단계는
"프론트도 성능 이슈를 본다"를 말로만 하지 않고,

- memoization
- stable callbacks
- derived state optimization

을 통해 실제 렌더링 비용을 줄이는 방향으로 구조를 보강한 단계입니다.
