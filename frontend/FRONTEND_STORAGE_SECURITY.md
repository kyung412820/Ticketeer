# Frontend Storage & Security Strategy

## Goal
TICKETEER 프론트는 브라우저 저장소를 사용할 때
**무엇을 어디에 저장할지**, 그리고 **언제 정리할지**를 분명히 하는 것을 기준으로 설계했습니다.

프론트는 보안을 보장하는 계층이 아니기 때문에,
저장 전략의 목표는 **사용자 흐름을 안정적으로 유지하되, 저장 범위를 최소화하는 것**입니다.

---

## 1. Core Principle

프론트에서 다루는 값은 언제든 조작될 수 있다고 가정합니다.

즉,
- `client_id`
- `queue_token`
- `event_id`
- `seat_id`

같은 값은 프론트에서 보관하더라도 **신뢰하지 않습니다**.

실제 검증은 모두 백엔드가 수행합니다.

### Frontend responsibility
- 사용자 흐름 유지
- 중복 행동 방지
- 잘못된 진입 최소화
- 세션성 데이터 정리

### Backend responsibility
- queue token 검증
- seat / event 관계 검증
- hold / booking 정합성 검증
- 권한과 유효성 최종 판단

---

## 2. Why `client_id` Uses localStorage

`client_id`는 브라우저 단위 식별 보조값으로 사용합니다.

### 저장 위치
- `localStorage`

### 이유
- 같은 브라우저에서 대기열을 반복 진입할 때 동일 사용자를 어느 정도 구분하기 쉬움
- 새로고침 후에도 유지되어 사용자 흐름이 덜 끊김
- 로그인 기반 사용자 식별이 없는 MVP 구조에서 최소한의 브라우저 단위 식별 역할 가능

### 한계
- 사용자가 직접 삭제/변경 가능
- 다른 브라우저/기기와는 연결되지 않음
- 보안 토큰이 아니라 UX 보조용 식별값일 뿐임

즉 `client_id`는 **신뢰 가능한 인증 정보가 아니라, 흐름 유지를 위한 보조 식별값**입니다.

---

## 3. Why `queue_token` Uses sessionStorage

`queue_token`은 예매 흐름에만 필요한 세션성 데이터입니다.

### 저장 위치
- `sessionStorage`

### 이유
- queue token은 영구 저장할 필요가 없음
- 브라우저 탭/세션 범위 안에서만 유지되는 것이 더 적절함
- 예매 흐름 종료 후 자연스럽게 정리하기 쉬움
- `localStorage`보다 범위를 좁게 가져갈 수 있음

### 장점
- 새 탭/새 세션에 자동으로 퍼지지 않음
- 흐름 종료 후 제거 기준을 적용하기 쉬움

즉 `queue_token`은 **짧은 흐름에 필요한 세션성 값**이라서 `sessionStorage`가 더 적합합니다.

---

## 4. Why Query Params Are Not the Primary Source

현재 구현은 query param을 보조적으로 사용할 수 있지만,
기본 방향은 **storage 중심**입니다.

### 이유
- query string은 URL에 그대로 드러남
- 복사 / 공유 / 히스토리 저장 시 노출 가능성이 있음
- 세션성 토큰을 URL에 과하게 의존하는 것은 바람직하지 않음

그래서 현재 전략은 다음과 같습니다.

- query param이 있으면 초기 진입 보조값으로 활용 가능
- 실제 지속 상태는 `sessionStorage` 기준으로 관리
- 가능하면 토큰은 세션성 저장소 기준으로 흘려보내는 방향 유지

---

## 5. Token Cleanup Rules

세션성 데이터는 **언제 지울지**가 중요합니다.

현재 정리 기준:

### queue token 제거 시점
- queue 상태가 `EXPIRED`가 되었을 때
- queue status 조회 에러가 발생했을 때
- 예매가 성공적으로 완료되었을 때

### 이유
- 이미 쓸모없는 토큰을 오래 남기지 않기 위해
- 잘못된 재진입과 stale state를 줄이기 위해
- 다음 예매 흐름에서 이전 세션 값이 섞이지 않게 하기 위해

---

## 6. Security Viewpoint

프론트 저장 전략은 “안전하게 숨긴다”가 아니라
**민감도를 낮게 유지하고, 신뢰를 백엔드에 두는 것**이 핵심입니다.

### What we intentionally avoid
- queue token을 장기 저장하지 않음
- localStorage를 모든 토큰 저장소로 사용하지 않음
- 프론트 값만 믿고 예매를 확정하지 않음

### What we rely on instead
- 백엔드 queue token 검증
- event / seat / client 관계 검증
- hold / booking 단계의 서버 측 유효성 검증

---

## 7. Summary

현재 저장 전략은 다음과 같습니다.

- `client_id` → `localStorage`
  - 브라우저 단위 식별 보조값
- `queue_token` → `sessionStorage`
  - 예매 흐름용 세션성 값
- query param
  - 보조 진입 수단
- 실제 보안과 정합성 검증
  - 모두 백엔드 책임

이 전략은
**프론트에서 저장 범위를 최소화하면서도 사용자 흐름을 유지하기 위한 선택**입니다.
