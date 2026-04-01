# Frontend Deep TODO 3 - 접근성 보강

## What changed
- 에러 배너에 `role="alert"` / `aria-live="assertive"` 적용
- 로딩 / 상태 변화 문구에 `aria-live="polite"` 적용
- 좌석 버튼에 `aria-label`, `aria-pressed`, `aria-disabled` 적용
- 액션 버튼에 `aria-disabled` 반영
- 접근성 정리 문서 추가

## Why
티켓팅 UI는 상태 변화가 많기 때문에,
시각적인 표시만으로는 부족합니다.
에러 / 로딩 / 좌석 선택 상태가 보조기술에도 전달되어야 합니다.

## Next
- 저장 전략 보안 관점 정리
