# Frontend Accessibility Deeper Patch

## Added / Changed
- focus management on selected seat
- focus management on booking completion
- landmark / section labeling improvements
- current state exposure with `aria-current`
- deeper accessibility note

## Why
접근성은 aria 속성 몇 개 추가로 끝나지 않습니다.
상태 전환이 많은 티켓팅 UI에서는
포커스 흐름과 구조 정보도 같이 고민해야 합니다.

## Portfolio Point
- 상태 변화가 많은 UI에서 focus management까지 고려
- landmark 구조와 live region을 함께 설계
