# Frontend Accessibility - Deeper Improvements

## Goal
기본적인 aria 속성 추가를 넘어서,
실제 사용자 흐름에서 더 자연스러운 접근성을 제공하는 방향으로 보강했습니다.

---

## What changed

### 1. Focus management
- 좌석 선택 시 해당 좌석 버튼으로 focus 이동
- 예매 완료 시 결과 카드 heading으로 focus 이동

이유:
- 스크린리더 / 키보드 사용자에게 현재 상태 전환이 더 명확히 전달됨
- "선택됨" / "완료됨" 시점의 컨텍스트를 잃지 않게 함

### 2. Landmark / section structure
- booking 페이지에 명확한 `section`, `aside`, `aria-labelledby`, `aria-label` 부여
- 좌석 목록과 예매 요약 영역의 의미를 더 분명하게 구성

이유:
- 화면 리더 사용자가 페이지 구조를 더 쉽게 파악할 수 있음

### 3. Current state exposure
- 선택된 좌석에 `aria-current`
- 좌석 목록, 범례, 구역 목록에 의미 있는 레이블 추가

이유:
- 단순히 pressed 여부뿐 아니라 "현재 선택 중인 항목"을 더 명확히 전달하기 위함

---

## What remains for future work
- 좌석 선택용 roving tabindex 패턴
- 방향키 기반 좌석 탐색
- skip link
- 더 정교한 포커스 복귀 전략

---

## Summary
이번 단계는 "aria 속성을 조금 붙였다" 수준이 아니라,
**상태 전환 시 포커스와 구조까지 고려한 접근성 보강** 단계입니다.
