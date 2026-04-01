# TICKETEER 프로젝트 총정리 문서
부제: 자소서 · 포트폴리오 · 면접 답변용 상세 기술 정리

---

## 1. 문서 목적

이 문서는 TICKETEER 프로젝트를 단순 구현 결과가 아니라,  
**문제 정의 → 기술 선택 → 설계 의사결정 → 검증 → 결과 → 회고**까지 일관되게 설명하기 위한 정리 문서다.

활용 목적은 다음과 같다.

- 자소서에서 프로젝트 설명 근거로 활용
- 포트폴리오 발표 시 설계 의사결정 설명 자료로 활용
- 면접에서 기술 선택 이유와 트러블슈팅을 구조적으로 답변하기 위한 기준 문서로 활용

즉, 이 문서는 “무엇을 만들었는가”보다  
**“왜 그렇게 설계했고, 어떤 문제를 어떻게 해결했는가”**에 초점을 둔다.

---

## 2. 프로젝트 한 줄 정의

TICKETEER는 **대기열 기반 티켓 예매 시스템**으로,  
동일 좌석에 대한 **중복 홀드 / 중복 예매 방지**와  
실시간 상태 변화가 많은 예매 흐름에서의 **화면 정합성 유지**를 핵심 문제로 두고 구현한 프로젝트다.

---

## 3. 왜 이 프로젝트를 만들었는가

티켓팅 서비스는 겉으로 보면 단순한 예매 서비스처럼 보이지만, 실제로는 일반 CRUD보다 훨씬 까다로운 문제를 가진다.

대표적으로:

- 많은 사용자가 짧은 시간 안에 동시에 진입한다.
- 여러 사용자가 같은 좌석을 동시에 선택한다.
- 임시 점유와 최종 확정의 성격이 다르다.
- 프론트 화면이 서버 상태보다 늦게 반영되면 잘못된 안내를 줄 수 있다.
- “정말 예매가 완료되었는가”, “중복 요청은 없었는가”, “화면이 거짓말하지 않는가”가 서비스 신뢰와 직결된다.

이 프로젝트는 단순히 예매 기능을 만드는 것이 아니라,
다음 두 문제를 해결하는 것을 목표로 했다.

1. **백엔드 관점**  
   동시 요청 환경에서도 데이터 정합성을 보장할 수 있는가

2. **프론트엔드 관점**  
   실시간 상태 변화 속에서도 사용자가 신뢰할 수 있는 UI를 제공할 수 있는가

---

## 4. 핵심 사용자 흐름

핵심 사용자 흐름은 다음과 같다.

1. 공연 목록 조회
2. 공연 상세 조회
3. 대기열 진입
4. 대기열 상태 조회
5. 좌석 목록 조회
6. 좌석 홀드
7. 예매 확정

이 흐름에서 가장 중요했던 점은  
단순히 단계가 이어지는 것이 아니라,  
**경쟁 상황에서도 시스템이 무너지지 않아야 한다**는 점이었다.

---

## 5. 전체 아키텍처 개요

### 5-1. 시스템 구성
- **Frontend**: Next.js 기반 웹 애플리케이션
- **Backend**: Go + Gin 기반 REST API 서버
- **PostgreSQL**: 최종 예매 정합성 보장
- **Redis**: 임시 좌석 홀드와 대기열 상태 관리
- **Playwright / Go test / k6**: 기능, 브라우저 흐름, 부하/경쟁 상황 검증
- **Docker Compose / Kubernetes**: 실행 환경 표준화 및 배포 단위 분리

### 5-2. 역할 분리 원칙
핵심 설계 원칙은 **역할 분리**였다.

- Redis: 빠르고 짧은 임시 상태 처리
- PostgreSQL: 최종 확정과 정합성 보장
- Frontend: 사용자 행동 제어와 화면 정합성 유지
- Backend: 동시성 환경에서의 비즈니스 규칙과 최종 상태 제어

---

## 6. 기술 스택과 선택 이유

## 6-1. Frontend

### Next.js
**선택 이유**
- 페이지별로 SSR/CSR 전략을 나눠 설명하기 좋다.
- 라우팅 구조가 명확하다.
- 단순 SPA보다 렌더링 전략에 대한 설계 설명이 가능하다.

**왜 적합했는가**
이 프로젝트는 페이지마다 성격이 다르다.

- `/events`, `/events/[id]`: 조회 중심
- `/queue/[eventId]`, `/booking/[eventId]`: 실시간 상태 변화 / 사용자 상호작용 중심

즉, 단순 React SPA보다 **페이지 성격에 맞춘 렌더링 전략**을 설명하기 좋은 구조가 필요했다.

---

### TypeScript
**선택 이유**
- 예매 흐름은 API 응답 구조가 많고 상태가 자주 바뀐다.
- 타입이 없으면 응답 구조 오해, nullable 처리 누락, 리팩터링 시 안정성 저하가 크다.

**적용 효과**
- API 응답 타입 일관화
- 상태 전환 로직 안정화
- 컴포넌트 props 구조 명확화

---

### Tailwind CSS
**선택 이유**
- 이 프로젝트의 목적은 화려한 디자인보다 상태 변화가 명확한 UI였다.
- disabled / selected / held / booked 같은 상태 차이를 빠르게 표현하기에 적합하다.

**적용 효과**
- 빠른 개발 속도
- 상태별 시각 차이 표현 용이
- 일관된 UI 유지

---

### React Query
**선택 이유**
초기에는 로컬 state + fetch wrapper로도 충분했지만,
대기열 상태와 좌석 목록은 명확한 **서버 상태(server state)**였다.

특히:
- queue status는 polling 필요
- seats list는 hold / booking 후 재동기화 필요
- stale state 관리가 중요

**적용 이유**
- polling을 컴포넌트 로직이 아니라 query 레벨에서 관리
- invalidate 기반 서버 상태 재동기화
- loading / error / stale 상태 관리 단순화

**적용 범위**
- queue status polling
- seats list fetch / invalidate

**왜 전면 도입이 아니라 부분 도입했는가**
신입 포트폴리오에서 중요한 것은 라이브러리 남용이 아니라
**서버 상태와 로컬 상태를 구분해 적절한 영역에 도입했는가**다.
가장 서버 상태 성격이 강한 `queue`, `seats`부터 도입하는 것이 가장 자연스러웠다.

---

### Playwright
**선택 이유**
- 실제 브라우저 기준으로 사용자 흐름을 검증할 수 있다.
- 접근성 role과 live region도 함께 검증 가능하다.
- 티켓팅처럼 “흐름이 중요한 UI”에 적합하다.

**적용 목적**
- 공연 목록 렌더링 및 상세 이동 검증
- 대기열 화면 상태 검증
- booking 직접 진입 차단 검증
- disabled 상태, alert/live region 검증

---

## 6-2. Backend

### Go
**선택 이유**
- 명확한 문법과 구조
- API 서버 구현에 적합
- goroutine 기반 동시성 테스트 작성이 용이
- 성능과 단순성이 모두 필요한 상황에 적합

**왜 적합했는가**
이 프로젝트는 복잡한 프레임워크보다
**도메인 문제와 동시성 제어에 집중**하는 것이 중요했다.

---

### Gin
**선택 이유**
- 라우팅이 단순하다.
- REST API 구현이 빠르다.
- 불필요한 추상화 없이 빠르게 구조를 잡을 수 있다.

**적용 효과**
- 예매 흐름 API를 빠르게 구성
- 핸들러 / 서비스 / 레포지토리 계층 분리
- 핵심 문제인 좌석 홀드 / 예매 확정 / 대기열 상태에 집중 가능

---

### GORM
**선택 이유**
- 모델 기반 개발 속도 향상
- CRUD 및 관계 구조를 빠르게 정리 가능
- 반복적인 쿼리 작업을 줄일 수 있음

**주의한 점**
정합성이 중요한 지점에서는 ORM에만 기대지 않았다.

- transaction
- row-level lock
- unique constraint

처럼 핵심 정합성 문제는 **DB 중심 사고**로 설계했다.

---

## 6-3. Database / Cache

### PostgreSQL
**선택 이유**
- transaction 지원
- row-level lock 지원
- unique constraint 활용 가능
- 관계형 데이터 구조에 적합

**프로젝트 내 역할**
PostgreSQL은 단순 저장소가 아니라,
**최종 예매 확정과 데이터 정합성을 책임지는 시스템**으로 사용했다.

---

### Redis
**선택 이유**
- `SETNX`를 통한 중복 홀드 방지
- TTL 기반 임시 홀드 만료 처리
- 빠른 응답

**프로젝트 내 역할**
Redis는 **임시 좌석 홀드와 대기열 상태 관리**를 담당했다.
최종 예매 확정은 Redis가 아니라 PostgreSQL에서 처리했다.

---

## 7. 백엔드 설계 상세

## 7-1. 핵심 문제: 동일 좌석 중복 예매 방지

티켓팅 시스템에서 가장 중요한 문제는
여러 사용자가 같은 좌석을 동시에 요청할 때
**최종 예매가 반드시 1건만 성공해야 한다**는 점이다.

이를 위해 다음 구조를 적용했다.

### 좌석 홀드 단계
- Redis `SETNX + TTL`
- 동일 좌석에 대해 먼저 온 사용자만 홀드 가능
- TTL 만료 시 자동 해제

### 예매 확정 단계
- queue token 검증
- Redis hold 상태 확인
- PostgreSQL transaction 시작
- `SELECT ... FOR UPDATE`로 seat row lock 획득
- seat 상태 재검증
- booking row 생성
- queue 상태 갱신
- commit 성공 이후 Redis hold key 제거

즉, 예매 확정 전체를 하나의 transaction 단위로 묶어
중간 실패 시 rollback 가능하도록 설계했다.

---

## 7-2. 왜 Redis만으로 끝내지 않았는가

홀드 단계는 Redis로 빠르게 처리할 수 있지만,
최종 예매 확정은 Redis만으로 충분히 신뢰하기 어렵다.

왜냐하면:
- 네트워크 지연과 복수 요청 경쟁 상황
- 프로세스 실패
- 상태 꼬임
- 최종 데이터 일관성 검증 부족

때문이다.

그래서 역할을 분리했다.

- Redis: **빠른 임시 선점**
- PostgreSQL: **최종 확정과 정합성**

이렇게 설계한 이유는
티켓팅의 핵심 문제를 “빠른 응답”과 “최종 일관성”으로 분리해 해결하기 위해서였다.

---

## 7-3. DB 차원의 방어선

애플리케이션 로직 외에도 DB 차원의 방어선을 추가했다.

- `seats(event_id, seat_no)` unique
- `bookings.seat_id` unique

이유:
애플리케이션 로직은 버그가 있을 수 있지만,
DB 제약조건은 **최종 방어선** 역할을 한다.
즉, 코드 실수나 극단적인 경쟁 상황에서도 DB가 마지막으로 정합성을 막아주도록 설계했다.

---

## 8. 프론트엔드 설계 상세

## 8-1. 핵심 문제: 화면 정합성과 사용자 흐름 보호

프론트엔드에서 가장 중요했던 질문은 다음과 같았다.

> 서버 상태가 빠르게 바뀌는 티켓팅 화면에서, 사용자가 믿을 수 있는 UI를 어떻게 유지할 것인가?

이 문제를 해결하기 위해 다음을 집중적으로 설계했다.

---

## 8-2. 중복 행동 방지

예매 UI에서는 버튼 연타만으로도 큰 혼란이 생긴다.

적용한 것:
- 대기열 진입 중복 호출 방지
- 좌석 홀드 버튼 연타 방지
- 예매 확정 버튼 연타 방지
- 요청 중 버튼 비활성화
- 요청 중 상태 문구 표시

**왜 중요했는가**
사용자 입장에서는 “버튼을 여러 번 눌렀는데 요청이 몇 번 갔는지”를 모른다.
티켓팅처럼 민감한 도메인에서는 이런 불확실성이 곧 UX 신뢰 하락으로 이어진다.

---

## 8-3. stale state 방지

티켓팅에서는 화면과 서버 상태가 쉽게 어긋난다.

예:
- 화면에는 `AVAILABLE`
- 실제 서버는 이미 `HELD` 또는 `BOOKED`

이를 막기 위해:
- 홀드 성공 후 좌석 재조회
- 홀드 실패 후 좌석 재조회
- 예매 성공 후 좌석 재조회
- 예매 실패 후 좌석 재조회

구조를 적용했다.

또한 티켓팅 도메인 특성상
**optimistic update를 최소화**하고,
서버 응답 이후 상태를 반영하는 쪽을 택했다.

**왜 이렇게 했는가**
소셜 피드나 좋아요 기능에서는 optimistic update가 UX에 도움이 되지만,
예매에서는 잘못된 optimistic update가 곧 **거짓 UI**가 된다.
따라서 이 도메인에서는 “빠른 착시”보다 “정확한 반영”이 더 중요하다고 판단했다.

---

## 8-4. polling과 비동기 상태 관리

대기열 화면은 polling 기반으로 동작하기 때문에 다음을 고려했다.

- `READY`, `EXPIRED` 시 polling 중단
- 페이지 이탈 시 interval cleanup
- 중복 polling 방지
- 잘못된 접근 차단 (`queueToken`, `clientId` 없음)

**왜 중요했는가**
대기열 화면은 단순 조회가 아니라
시간에 따라 상태가 변하고, 사용자가 기다리는 동안 UI가 계속 갱신된다.
이때 polling이 과하게 중첩되면:
- 불필요한 요청 증가
- 상태 꼬임
- 메모리 누수
- 사용자 혼란

으로 이어질 수 있기 때문이다.

---

## 8-5. 상태 저장 전략

저장 전략도 역할에 따라 분리했다.

### `client_id` → `localStorage`
- 브라우저 단위 식별 보조값
- 새로고침 후에도 흐름 유지 가능
- 로그인 없는 MVP 구조에서 최소한의 식별 역할

### `queue_token` → `sessionStorage`
- 예매 흐름에만 필요한 세션성 값
- 장기 저장 불필요
- 예매 완료 / 만료 시 제거 쉬움
- localStorage보다 노출 범위를 줄일 수 있음

추가 원칙:
- queue 만료 시 토큰 제거
- 예매 완료 시 토큰 제거
- query param은 보조 수단으로만 사용
- 프론트 저장값을 절대 신뢰하지 않음

**왜 이렇게 나눴는가**
모든 값을 localStorage에 몰아넣는 방식은 단순하지만,
세션성 데이터와 장기 식별 보조값의 성격이 다르다.
이 차이를 구분하는 것이 더 현실적인 설계라고 판단했다.

---

## 8-6. 렌더링 전략

페이지 성격에 따라 렌더링 전략을 나눴다.

### `/events`, `/events/[id]`
- 조회 중심 화면
- 서버 중심 렌더링

**선택 이유**
- 상대적으로 정적인 데이터
- 초기 로딩과 문서성 페이지에 적합
- SEO와 첫 화면 노출에 유리

### `/queue/[eventId]`, `/booking/[eventId]`
- 실시간 상태 변화
- 사용자 상호작용 중심
- 클라이언트 중심 렌더링

**선택 이유**
- polling 필요
- 브라우저 상태 (`queueToken`, `clientId`) 활용
- 액션 후 재조회와 즉시 UI 업데이트 필요

**왜 중요한가**
이 구분은 “Next.js를 써봤다”가 아니라,
**어떤 페이지를 왜 서버 중심 / 클라이언트 중심으로 나눴는지 설명할 수 있게 해준다.**

---

## 8-7. 네트워크 실패 대응

요청 성격에 따라 timeout / retry 정책을 다르게 적용했다.

### retry 적용
- `fetchSeats`
- `getQueueStatus`

**이유**
조회성 요청은 한두 번 재시도하는 것이 UX에 도움이 된다.

### retry 미적용
- `enterQueue`
- `holdSeat`
- `createBooking`

**이유**
상태 변경 요청, 특히 booking은 자동 재시도 시
중복 요청과 사용자 불신을 초래할 수 있다.
예매 확정에서 중요한 것은 자동 retry가 아니라,
- 중복 행동 방지
- 명확한 처리 상태 표시
- 실패 시 재조회와 올바른 안내

였다.

---

## 8-8. 접근성

실시간 상태 변화가 많은 UI인 만큼 기본 접근성과 한 단계 더 깊은 접근성까지 반영했다.

### 기본 적용
- 에러 배너 `role="alert"`
- 상태 문구 `aria-live`
- 좌석 버튼 `aria-label`
- `aria-pressed`
- `aria-disabled`

### 추가 적용
- 선택된 좌석으로 focus 이동
- 예매 완료 후 결과 카드 heading으로 focus 이동
- `section`, `aside`, `aria-labelledby` 구조 정리

**왜 중요했는가**
티켓팅 UI는
- 선택
- 실패
- 대기
- 완료

처럼 상태 변화가 많다.
이때 시각적 변화만으로는 충분하지 않다.
스크린리더와 키보드 사용자도 현재 상태와 컨텍스트를 잃지 않도록 구성하는 것이 중요했다.

---

## 9. React Query 도입 배경과 결과

초기에는 로컬 state + fetch wrapper로도 충분해 보였지만,
프로젝트를 프론트 메인 포트폴리오로 발전시키는 과정에서 질문이 생겼다.

> 서버 상태를 정말 서버 상태답게 다루고 있는가?

그래서 React Query를 실도입했다.

### 왜 queue와 seats부터 적용했는가
이 둘이 가장 서버 상태 성격이 강했기 때문이다.

- `queue status`: polling 필요
- `seats list`: hold / booking 후 invalidate 필요

### 도입 후 변화
- queue polling을 query 기준으로 관리
- `READY`, `EXPIRED` 시 polling 자동 중단
- seats list를 query로 관리
- hold / booking 후 `invalidateQueries`로 재동기화
- 직접 fetch + setState 반복 감소

### 결과
React Query는 단순히 라이브러리 추가가 아니라,
**서버 상태와 UI 상태를 분리해 설계한 근거**가 되었다.

---

## 10. 성능 최적화

현재 좌석 수가 아주 많지는 않지만,
티켓팅 UI는 좌석 수가 늘어나면 프론트 렌더링 비용이 커질 수 있다.

그래서 지금 단계에서 의미 있는 최적화를 먼저 적용했다.

### 적용한 최적화
- `SeatSectionGrid` memoization
- section / seat card 단위 memoization
- section grouping `useMemo`
- `SeatSelectionSummary` memoization
- stable callback 유지 (`useCallback`)
- 파생값 `useMemo` 정리

### 왜 이 최적화를 먼저 했는가
아직 좌석 수가 수백~수천 석 수준은 아니므로,
지금 당장 virtualization까지 넣는 것은 과할 수 있다.

그래서 먼저:
- 불필요한 re-render 감소
- render churn 완화
- refetch 이후 하위 카드 흔들림 감소

처럼 **지금 구조에서 효과가 큰 최적화**를 적용했다.

---

## 11. 테스트 전략

## 11-1. 백엔드 테스트

Go 테스트로 다음을 검증했다.

- 정상 예매 흐름
- 예매 오픈 전 진입 차단
- 이미 예매된 좌석 차단
- 홀드 없이 예매 차단
- 잘못된 queue token 차단
- 이미 홀드된 좌석 차단
- 동시 hold 경쟁 제어
- 동시 booking 경쟁 제어

**의미**
단순 기능 동작이 아니라,
동시 요청 경쟁 상황에서도 비즈니스 규칙이 유지되는지를 검증했다.

---

## 11-2. 프론트 테스트

Playwright 기반 smoke / E2E 테스트를 적용했다.

검증 항목:
- 공연 목록 렌더링
- 상세 이동
- 대기열 화면 상태 확인
- booking 직접 진입 차단
- alert / live region / 버튼 노출 확인
- 버튼 disabled 상태 확인

### 테스트 안정화 과정
초기에는 `getByRole("alert")`가 Next 내부 route announcer와 충돌했다.
그래서:
- `error-banner`
- `loading-state`

에 `data-testid`를 추가하고,  
selector를 안정화했다.

**의미**
프론트 테스트도 단순히 “있다”가 아니라,
**덜 깨지는 테스트 구조를 설계했다**는 점이 중요했다.

---

## 11-3. 부하 테스트

k6로 주요 API의 부하와 경쟁 상황을 검증했다.

### queue enter
- burst traffic 시나리오
- 실패율 0%
- p95 응답시간 안정적으로 확인

### seat hold / booking 경쟁
이 테스트의 목적은 “전부 성공”이 아니다.
동일 좌석을 여러 요청이 동시에 노릴 때:
- 1건만 성공하고
- 나머지는 정상 거절되는지

를 확인하는 것이다.

### 측정 결과
- `queue_enter`: 실패율 **0%**
- `seat_hold`: p95 **16.22ms**
- `booking`: p95 **32.97ms**

**해석**
경쟁 시나리오에서 일부 요청만 성공하고 나머지가 거절된 것은
오류가 아니라 **중복 점유 / 중복 예매 방지 로직이 의도대로 동작한 결과**다.

---

## 12. Docker Compose 적용

초기에는 backend, frontend, postgres, redis를 각각 수동 실행했다.
이후 `docker-compose.yml`을 구성해 전체 서비스를 한 번에 띄울 수 있도록 정리했다.

### 왜 필요했는가
- 로컬 실행 절차 단순화
- 환경 재현성 향상
- 타인이 프로젝트를 쉽게 실행 가능
- 포트폴리오 전달력 향상

### 구성 요소
- postgres
- redis
- backend
- frontend

### 설계 포인트
- backend는 postgres / redis 준비 이후 기동
- frontend는
  - 브라우저 기준 `NEXT_PUBLIC_API_BASE_URL`
  - 서버 렌더링 기준 `INTERNAL_API_BASE_URL`
  를 분리

### 트러블슈팅
초기에는 backend가 DB에 연결하지 못했는데,
원인은 애플리케이션이 실제로 읽는 환경변수명이 `POSTGRES_HOST` 계열인데
compose에서는 `DB_HOST` 계열로 주입하고 있었기 때문이다.

이를 수정해:
- `POSTGRES_HOST=postgres`
- `REDIS_HOST=redis`

형태로 맞춘 뒤 전체 서비스가 정상 기동되었다.

---

## 13. Kubernetes 적용

프로젝트를 로컬 실행과 Docker Compose에만 머무르지 않고,
Kubernetes 배포 단위까지 확장했다.

### 적용 리소스
- Namespace
- ConfigMap
- Secret
- Deployment
- Service
- PVC
- Ingress

### 왜 적용했는가
- 환경변수와 비밀값 분리
- Service DNS 기반 통신
- 상태 저장 리소스 분리
- readiness / liveness probe 적용
- 배포 단위 관점 학습과 검증

### 실제 검증 내용
- backend pod 환경변수 확인
- PostgreSQL Service / Endpoint 확인
- backend pod 내부 DNS 해석 확인
- debug pod에서 `pg_isready` 확인
- postgres DB 존재 확인
- backend pod `Running 1/1` 상태 확인

### 적용 중 발생한 문제
배포 초기에는 backend가 `CrashLoopBackOff`에 빠졌다.

초기 원인:
- PostgreSQL이 준비되기 전에 backend가 먼저 기동
- connection refused 발생

이후:
- ConfigMap 환경변수명 정리
- Service 기반 이름 통신 확인
- DNS / Endpoint / DB 연결 상태 직접 검증

을 통해 현재는 정상 상태로 수렴했다.

### 현재 상태 요약
- `postgres`: Running
- `redis`: Running
- `frontend`: Running
- `backend`: Running
- backend `/api/health` 정상 응답
- Postgres Service 및 Endpoint 정상
- DB 접속 및 DNS 해석 정상

---

## 14. 주요 트러블슈팅

## 14-1. Redis 홀드만으로는 최종 예매 정합성을 보장하기 어려웠던 문제
**문제**  
좌석 홀드는 Redis로 빠르게 처리할 수 있었지만,
최종 예매 확정까지 Redis 중심으로 처리하면 동일 좌석 경쟁 상황에서 강한 정합성을 보장하기 어려웠다.

**왜 어려웠는가**  
티켓팅에서는 임시 선점과 최종 확정의 성격이 다르다.
홀드는 빠른 응답이 중요하지만, 최종 확정은 반드시 데이터 일관성이 더 중요하다.

**해결**
- 좌석 홀드: Redis `SETNX + TTL`
- 예매 확정: PostgreSQL transaction
- 경쟁 제어: `SELECT ... FOR UPDATE`
- DB unique 제약조건 추가

**결과**
- 동일 좌석 동시 hold 요청: 1건만 성공
- 동일 좌석 동시 booking 요청: 1건만 성공
- booking row 최종 1건만 생성

---

## 14-2. 정상 흐름 구현만으로는 티켓팅 핵심 문제를 설명하기 어려웠던 문제
**문제**  
정상 흐름이 동작한다고 해서 티켓팅 시스템의 핵심 문제를 해결했다고 보긴 어려웠다.

**해결**
- goroutine 기반 동시 hold 테스트
- goroutine 기반 동시 booking 테스트

**결과**
- 동일 좌석 경쟁 시 최종 성공 1건만 남는 것을 검증
- 단순 CRUD가 아니라 동시성 제어까지 설명 가능한 프로젝트가 됨

---

## 14-3. 프론트에서 stale state가 쉽게 생기는 문제
**문제**  
좌석 상태가 빠르게 바뀌는 도메인이라
화면과 서버 상태가 쉽게 어긋났다.

**해결**
- optimistic update 최소화
- 액션 후 재조회
- React Query invalidate 적용

**결과**
- 좌석 상태를 서버 기준으로 다시 맞추는 구조 확보
- 사용자가 오래된 상태를 믿게 되는 문제 완화

---

## 14-4. 비동기 UI가 쉽게 불안정해지는 문제
**문제**  
대기열은 polling, 좌석 선택은 재조회가 많아
중복 요청, stale response, 잘못된 직접 진입 문제가 생기기 쉬웠다.

**해결**
- 중복 요청 방지
- polling 종료 조건 명확화
- cleanup 처리
- 잘못된 진입 차단
- query 기반 상태 관리

**결과**
- queue status는 상태에 따라 polling 자동 중단
- booking 직접 진입 차단
- 액션 이후 상태 재동기화 구조 확보

---

## 14-5. k6 경쟁 테스트에서 실패율이 높아 보이는 문제
**문제**  
동일 좌석 경쟁 시나리오에서 4xx 비율이 높아 실패처럼 보였다.

**왜 어려웠는가**  
일반 부하 테스트 기준으로는 실패율이 낮아야 좋지만,
이 테스트는 일부러 동일 좌석에 요청을 충돌시키는 시나리오였다.

**해결**
- `queue_enter`는 일반 부하 테스트로 해석
- `seat_hold`, `booking`은 경쟁 제어 테스트로 해석
- 응답시간과 정상 거절 여부 중심으로 판단 기준 정리

**결과**
- `queue_enter`: 실패율 0%
- `seat_hold`: p95 16.22ms
- `booking`: p95 32.97ms

---

## 14-6. 프론트 테스트에서 UI는 맞는데 테스트가 깨지던 문제
**문제**  
Playwright에서 role selector가 Next 내부 요소와 충돌했다.

**해결**
- `data-testid` 추가
- stable selector 전략 사용

**결과**
- 직접 진입 차단
- 버튼 disabled
- loading / error / queue content 검증

을 더 안정적으로 테스트할 수 있게 되었다.

---

## 14-7. Docker Compose에서 backend가 DB에 연결하지 못하던 문제
**문제**  
backend가 PostgreSQL에 연결하지 못하고 `localhost:5432`를 바라보는 문제가 있었다.

**원인**
- 애플리케이션이 읽는 환경변수명과 compose에서 주입하는 이름이 불일치
- 결과적으로 기본값 `localhost` 사용

**해결**
- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- `REDIS_HOST`, `REDIS_PORT`

형태로 환경변수명 정리

**결과**
- docker-compose 환경에서 전체 서비스 정상 기동
- 실행 환경 재현성 확보

---

## 14-8. Kubernetes에서 backend가 CrashLoopBackOff에 빠지던 문제
**문제**  
배포 직후 backend pod만 재시작되며 CrashLoopBackOff 상태가 발생했다.

**원인**
- PostgreSQL 준비 타이밍 이전 연결 시도
- 초기 connection refused 발생

**해결**
- ConfigMap 환경변수 정리
- Service / Endpoint / DNS 직접 검증
- readiness / liveness probe 적용
- pod 내부 env / service / DB 상태 점검

**결과**
- backend pod `1/1 Running`
- health check 정상 통과
- 서비스 간 연결 정상화

---

## 15. 성과 정리

### 백엔드 관점 성과
- Redis와 PostgreSQL의 역할을 분리해 설명 가능한 구조를 만들었다.
- transaction과 row-level lock을 예매 흐름에 실제 적용했다.
- 동시성 테스트로 “최종 1건 성공”을 검증했다.

### 프론트 관점 성과
- 단순 UI 구현을 넘어 stale state, 중복 행동, polling, retry 정책, 접근성을 설계했다.
- SSR/CSR 전략을 페이지 성격에 맞게 분리했다.
- React Query를 실도입해 서버 상태를 관리했다.
- 테스트와 접근성까지 반영한 프론트 설계 근거를 만들었다.

### 운영 / 배포 관점 성과
- Docker Compose로 전체 실행 환경을 표준화했다.
- Kubernetes에서 Deployment / Service / ConfigMap / Secret / PVC / Ingress를 적용했다.
- 로컬 → 컨테이너 → 오케스트레이션 단계까지 확장 가능한 프로젝트가 되었다.

---

## 16. 핵심 도메인 모델

### Event
공연 단위 엔티티다.  
예매 가능 시간, 공연 정보, 좌석 집합의 상위 개념이다.

### Seat
공연에 속한 좌석 엔티티다.  
`seat_no`, `section`, `price`, `status`를 가지며, 티켓팅의 핵심 경쟁 자원이다.

### QueueEntry
대기열 상태 엔티티다.  
사용자의 대기열 진입, READY 상태 전환, 만료 여부를 표현한다.

### Booking
최종 예매 확정 엔티티다.  
좌석에 대한 최종 소유 상태를 나타내며, 동일 좌석에 대해 1건만 생성되도록 제약을 둔다.

### 모델링 원칙
- Seat는 경쟁 자원이므로 최종 확정 시 DB 락의 대상이 된다.
- QueueEntry는 예매 가능 상태를 통제하는 흐름 엔티티다.
- Booking은 결과 엔티티이며, 최종 정합성의 기준점이다.

---

## 17. API 설계 원칙

이 프로젝트에서는 조회 API와 상태 변경 API를 명확히 분리했다.

### 조회 API
- 공연 목록 조회
- 공연 상세 조회
- 좌석 조회
- 대기열 상태 조회

조회 API는 상대적으로 재시도에 안전하며, 프론트에서 제한적 retry를 적용했다.

### 상태 변경 API
- 대기열 진입
- 좌석 홀드
- 예매 확정

상태 변경 API는 중복 호출 자체가 문제를 만들 수 있으므로, 프론트에서는 자동 retry를 적용하지 않았다.

### hold와 booking을 분리한 이유
티켓팅에서는 “임시 선점”과 “최종 확정”의 성격이 다르다.
따라서 빠른 응답이 중요한 hold와, 강한 정합성이 필요한 booking을 하나의 단계로 합치지 않고 분리했다.

### 에러 응답 설계 방향
에러는 단순 500 처리보다, 프론트가 사용자 친화 메시지로 변환할 수 있도록 도메인별 코드를 유지하는 방향으로 설계했다.
이 방식은 프론트에서 에러 메시지를 일관성 있게 매핑하는 데 도움이 되었다.

---

## 18. 한계와 트레이드오프

이 프로젝트는 신입 포트폴리오 범위 안에서 핵심 문제를 선명하게 다루는 데 집중했다.

### 일부러 하지 않은 것
- 결제 연동
- Kafka / 이벤트 드리븐 아키텍처
- MSA 분리
- 좌석 수천 개 기준 virtualization
- 운영용 모니터링/알람 시스템

### 그렇게 판단한 이유
핵심 문제는 “예매 정합성”과 “화면 정합성”이었기 때문이다.
범위를 넓히기보다, 가장 중요한 문제를 끝까지 검증 가능한 수준으로 만드는 것을 우선했다.

### 트레이드오프
- Redis만으로 빠르게 끝내는 대신 PostgreSQL transaction을 추가해 최종 정합성을 강화했다.
- optimistic update로 빠른 UX를 만들기보다, 서버 재조회 중심으로 신뢰 가능한 UI를 택했다.
- 프론트 상태관리도 전역 상태 남용보다 서버 상태와 로컬 상태를 구분하는 쪽을 택했다.

---

## 19. 결과 요약

### 기능 / 정합성
- 동일 좌석 동시 hold 경쟁: 1건만 성공
- 동일 좌석 동시 booking 경쟁: 1건만 성공
- booking row 최종 1건만 생성

### 성능 / 부하
- `queue_enter`: 실패율 0%
- `seat_hold`: p95 16.22ms
- `booking`: p95 32.97ms

### 프론트 품질
- React Query 기반 queue / seats 서버 상태 관리
- Playwright E2E / smoke 테스트 적용
- 접근성 속성 및 포커스 흐름 반영

### 배포 / 운영
- Docker Compose로 전체 실행 환경 표준화
- Kubernetes에서 backend / frontend / postgres / redis 분리 배포
- ConfigMap / Secret / Service / PVC / Ingress 적용

---

## 20. 자소서 / 면접에서 강조할 핵심 문장

### 프로젝트 소개 한 줄
대기열 기반 티켓 예매 시스템을 구현하면서, 백엔드에서는 동일 좌석 동시 요청에 대한 정합성을, 프론트에서는 실시간 상태 변화 속 화면 정합성과 사용자 신뢰를 중심으로 설계했습니다.

### 백엔드 강조 문장
좌석 홀드는 Redis `SETNX + TTL`로 처리하고, 최종 예매는 PostgreSQL transaction과 row-level lock으로 처리해 동일 좌석에 대한 중복 예매를 방지했습니다.

### 프론트 강조 문장
프론트는 단순 화면 구현보다 stale state, 중복 행동, polling, retry 정책, 접근성까지 포함해 실시간 예매 흐름에서 신뢰 가능한 UI를 만드는 데 초점을 맞췄습니다.

### 검증 강조 문장
기능 구현만으로 끝내지 않고, Go 테스트로 동시성 시나리오를, Playwright로 실제 브라우저 흐름을, k6로 부하와 경쟁 상황을 검증했습니다.

### 배포 강조 문장
Docker Compose로 전체 실행 환경을 표준화했고, Kubernetes에서는 ConfigMap/Secret, Service DNS, readiness/liveness probe, PVC를 적용해 애플리케이션을 배포 단위로 분리해봤습니다.

---

## 21. 이 프로젝트가 신입 포트폴리오에서 가지는 의미

이 프로젝트는 단순히 기술 스택을 나열한 포트폴리오가 아니라,
**문제를 정의하고, 기술을 목적에 맞게 선택하고, 검증까지 연결한 프로젝트**라는 점에서 의미가 있다.

특히 다음 질문에 답할 수 있게 해준다.

- 왜 Redis와 PostgreSQL을 같이 썼는가
- 왜 lock과 transaction이 필요한가
- 왜 optimistic update를 최소화했는가
- 왜 queue는 CSR이고 events는 서버 중심 렌더링인가
- 왜 booking은 retry하지 않았는가
- 왜 React Query를 queue/seats부터 도입했는가
- 왜 접근성과 Playwright를 넣었는가
- 왜 Docker Compose와 Kubernetes까지 확장했는가

즉 이 프로젝트는
“기능을 구현해봤다” 수준이 아니라,
**설계 이유와 검증 과정을 설명할 수 있는 포트폴리오**다.

---

## 22. 향후 확장 아이디어

### 백엔드
- 결제 흐름 추가
- queue 정책 고도화
- CI/CD 구축
- 운영 로그 / 모니터링 도입

### 프론트
- 정상 queue → booking 완료까지의 더 정교한 E2E
- 디자인 시스템화
- 좌석 수 증가 시 virtualization 검토
- 더 정교한 접근성 패턴 보강

### 인프라
- 이미지 레지스트리 연동
- Kubernetes 환경에서 Secret 관리 고도화
- Helm 또는 Kustomize 적용 검토

---

## 23. 최종 요약

TICKETEER는 티켓팅 서비스의 핵심 문제를
다음 두 축으로 나눠 해결한 프로젝트다.

### 백엔드
- Redis 기반 임시 홀드
- PostgreSQL transaction
- row-level lock
- unique constraint
- 동시성 검증

### 프론트
- stale state 방지
- 중복 행동 방지
- SSR/CSR 전략 분리
- React Query 도입
- retry/timeout 정책 설계
- 접근성 강화
- 브라우저 E2E 검증

### 운영
- Docker Compose 실행 환경 표준화
- Kubernetes 배포 단위 분리
- ConfigMap / Secret / Service / PVC / Ingress 적용

즉, 이 프로젝트의 핵심은 기능 구현 자체보다  
**왜 그렇게 설계했고, 어떤 문제를 해결했고, 그것을 어떻게 검증했는가**에 있다.
