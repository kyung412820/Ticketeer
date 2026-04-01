# Frontend Accessibility Notes

## Goal
TICKETEER 프론트는 기본적인 시각적 UI를 넘어서,
상태 변화와 상호작용이 보조기술에도 전달되도록 최소 접근성 보강을 적용했습니다.

## What was improved

### 1. Error message delivery
- 에러 배너에 `role="alert"` 적용
- `aria-live="assertive"`로 즉시 전달되도록 구성

### 2. Loading / status updates
- 로딩 문구와 상태 전환 영역에 `aria-live="polite"` 적용
- queue 상태 변화가 스크린리더에도 전달될 수 있도록 보강

### 3. Seat selection buttons
- 각 좌석 버튼에 `aria-label` 추가
- 좌석 번호 / 구역 / 가격 / 상태를 읽을 수 있게 구성
- `aria-pressed`로 현재 선택 상태 전달
- disabled 상태에 `aria-disabled` 반영

### 4. Action buttons
- 좌석 홀드 / 예매 확정 버튼에 `aria-disabled` 반영

## Why this matters
티켓팅 UI는 상태 변화가 많고,
실패 / 대기 / 선택 / 완료 같은 상태가 자주 바뀝니다.
이런 정보가 시각적으로만 보이면 접근성이 떨어지므로,
상태 변화가 보조기술에도 전달되도록 처리했습니다.

## Future work
- 키보드 포커스 이동 개선
- 예매 완료 후 포커스 이동 전략
- 더 정교한 landmark / heading 구조 개선
