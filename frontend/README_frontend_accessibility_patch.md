# Frontend Accessibility Enhancements

## Why
접근성은 UI에서 중요한 요소입니다. 상태 변화가 있을 때, 특히 좌석 선택이나 예매 완료 상태에서 사용자 경험을 고려해 접근성을 보강하였습니다.

## What
### 1. 포커스 이동
좌석 선택 후 해당 좌석으로 포커스를 이동시켜 사용자에게 명확한 시각적 피드백을 제공합니다.

### 2. 예매 완료 후 결과 카드
예매 완료 후 결과 카드로 포커스를 이동하여 스크린리더 사용자가 예약 정보를 쉽게 확인할 수 있도록 했습니다.

### 3. Landmark 구조 및 aria-labelledby
각각의 주요 섹션에 대해 `section`, `aside`, `aria-labelledby`를 추가하여 스크린리더 사용자가 페이지 구조를 쉽게 이해할 수 있도록 했습니다.

### 4. 좌석 선택 상태 강조
좌석 선택 상태에 `aria-pressed`와 `aria-label`을 활용하여 선택된 좌석과 선택되지 않은 좌석을 명확하게 구분합니다.

## Testing
- Playwright를 이용한 접근성 테스트는 `aria-live`, `aria-labelledby`, `aria-pressed` 속성의 활용을 확인하며 진행되었습니다.
- `focus` 이동과 `screen reader`로 확인된 상태 변화는 스크린리더 사용자가 쉽게 인식할 수 있도록 구현되었습니다.
