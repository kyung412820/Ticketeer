# Frontend State Management Strategy

## Goal
TICKETEER 프론트는 상태를 **UI 상태**와 **서버 상태**로 나눠서 다루는 것을 기준으로 설계했습니다.

---

## 1. Current Strategy

현재는 별도의 전역 상태 라이브러리 없이 다음 조합을 사용합니다.

- **local state (`useState`)**
- **공통 fetch wrapper (`lib/api-client.ts`)**
- **storage (`localStorage`, `sessionStorage`)**
- **필요 시 재조회(refetch)**

이 구조를 선택한 이유는,
현재 프로젝트 범위에서는 상태 복잡도가 아직 크지 않고,
핵심은 “복잡한 상태 도구 도입”보다 **예매 흐름의 정합성과 UX 안정성 확보**였기 때문입니다.

---

## 2. State Separation

### UI State
브라우저 안에서만 필요한 상태입니다.

예:
- 현재 선택한 좌석
- 버튼 로딩 상태 (`isHolding`, `isBooking`)
- 에러 메시지
- 예매 완료 카드 표시 여부

이 상태는 컴포넌트 로컬 state로 관리합니다.

---

### Server State
서버 기준으로 달라질 수 있는 상태입니다.

예:
- 공연 목록
- 공연 상세
- 대기열 상태
- 좌석 목록
- 예매 결과

이 상태는 API 호출 결과로 가져오고,
중요한 액션 후에는 **서버 기준으로 다시 조회**해서 맞춥니다.

---

## 3. Why We Did Not Use Optimistic Update

티켓팅 도메인에서는 낙관적 업데이트를 조심해야 합니다.

예를 들어:
- 화면에서는 좌석을 AVAILABLE로 보고 있음
- 사용자가 홀드/예매 버튼을 누름
- 실제 서버에서는 이미 다른 사용자가 먼저 처리했을 수 있음

그래서 이 프로젝트에서는
**클릭 즉시 상태를 확정해서 보여주는 대신,
서버 응답 이후에 상태를 반영하는 방향**을 선택했습니다.

즉,
- 홀드 성공 후 좌석 재조회
- 홀드 실패 후 좌석 재조회
- 예매 성공 후 좌석 재조회
- 예매 실패 후 좌석 재조회

흐름으로 구성했습니다.

---

## 4. Why We Did Not Introduce React Query Yet

React Query는 서버 상태 관리에 매우 유용하지만,
이번 프로젝트에서는 우선순위를 다음처럼 두었습니다.

### 먼저 해결한 것
- 중복 요청 방지
- polling cleanup
- 좌석 상태 재동기화
- 사용자 친화적 에러 메시지
- queue token / client id 저장 전략
- 예매 완료 UX

### 아직 React Query를 넣지 않은 이유
- 프로젝트 초기 범위에서는 페이지 수와 상태 종류가 제한적이었음
- 직접 구현한 fetch wrapper로도 현재 범위는 충분히 제어 가능했음
- 먼저 티켓팅 핵심 문제인 정합성과 경쟁 상황 대응을 우선 해결하고자 했음

---

## 5. When React Query Would Be Worth Adding

프로젝트가 더 커지면 React Query 도입 가치가 커집니다.

### 도입 후보 1: 공연 목록 / 상세
- 캐시 관리
- 중복 fetch 감소
- 로딩 / 에러 패턴 통일

### 도입 후보 2: 좌석 목록
- refetch 제어
- stale time 설정
- mutate 후 invalidation

### 도입 후보 3: queue status polling
- polling 관리 표준화
- query key 기준 상태 관리
- 상태 전환 후 refetch 중단 제어

---

## 6. Current Trade-off

현재 구조의 장점:
- 단순하고 읽기 쉬움
- 흐름 추적이 쉬움
- 예매 핵심 로직과 UX 제어에 집중 가능

현재 구조의 한계:
- 서버 상태 관리가 페이지 단위로 흩어져 있음
- 캐시 전략이 약함
- 페이지 수가 많아지면 중복 코드가 증가할 수 있음

즉,
현재는 **작은 범위에서 명확함을 우선한 구조**이고,
규모가 커지면 React Query 같은 도구를 도입하는 것이 자연스럽습니다.

---

## 7. Summary

현재 TICKETEER 프론트의 상태관리 전략은 다음과 같습니다.

- UI 상태는 local state로 관리
- 서버 상태는 fetch + 재조회 중심으로 관리
- 티켓팅 도메인 특성상 optimistic update는 최소화
- 현재 범위에서는 단순한 구조를 유지
- 추후 확장 시 React Query 도입 여지가 큼

이 전략은
**티켓팅 시스템에서 중요한 “화면 정합성”과 “잘못된 사용자 행동 방지”를 먼저 해결하기 위한 선택**입니다.
