# TICKETEER

대기열 기반 티켓 예매 시스템입니다.  
동일 좌석에 대한 **중복 홀드 / 중복 예매를 방지**하고,  
실시간 상태 변화가 많은 예매 흐름에서 **신뢰 가능한 화면 정합성**을 유지하는 데 초점을 맞췄습니다.

---

## 프로젝트 소개

티켓팅 서비스는 단순한 CRUD를 넘어, 짧은 시간에 많은 사용자가 몰리고 동일 좌석에 대한 경쟁이 발생하는 환경을 다뤄야 합니다.

TICKETEER는 이러한 문제를 해결하기 위해 다음 두 가지를 핵심 과제로 설정했습니다.

- 백엔드: 동시 요청 환경에서도 좌석 데이터 정합성 보장
- 프론트엔드: 실시간 상태 변화 속에서도 거짓말하지 않는 UI 제공

---

## 핵심 기능

- 공연 목록 / 상세 조회
- 대기열 진입 및 상태 조회
- 좌석 조회 및 임시 홀드
- 예매 확정
- 중복 요청 방지
- 예매 흐름 검증 테스트

---

## 기술 스택

### Frontend
- Next.js
- TypeScript
- Tailwind CSS
- React Query
- Playwright

### Backend
- Go
- Gin
- GORM

### Database / Cache
- PostgreSQL
- Redis

---

## 아키텍처 핵심 설계

### 1. 좌석 홀드와 최종 예매 책임 분리
- **Redis**: `SETNX + TTL` 기반 임시 홀드 처리
- **PostgreSQL**: transaction + row-level lock 기반 최종 예매 확정

임시 점유와 최종 확정의 역할을 분리해, 빠른 응답성과 최종 정합성을 동시에 확보했습니다.

### 2. 동시성 제어
동일 좌석에 여러 요청이 동시에 들어와도 최종적으로 **1건만 예매 성공**하도록 설계했습니다.

- Redis로 중복 홀드 방지
- PostgreSQL `SELECT ... FOR UPDATE`로 좌석 row lock 획득
- transaction 내부에서 상태 재검증 후 booking 생성

### 3. 프론트엔드 화면 정합성
예매 도메인에서는 잘못된 optimistic update가 사용자 신뢰를 해칠 수 있기 때문에,  
액션 이후 서버 기준 재조회 방식을 우선했습니다.

- hold 성공/실패 후 좌석 재조회
- booking 성공/실패 후 좌석 재조회
- 중복 클릭 방지
- polling 상태에 따른 UI 제어

---

## 렌더링 및 상태 관리 전략

### 렌더링 전략
- `/events`, `/events/[id]`  
  조회 중심 페이지로 서버 렌더링 기반 구성
- `/queue/[eventId]`, `/booking/[eventId]`  
  실시간 상호작용 중심 페이지로 클라이언트 렌더링 기반 구성

### 상태 관리
- React Query를 활용해 queue status polling 관리
- seats list를 query로 관리하고 `invalidateQueries` 기반 재동기화
- 서버 상태와 UI 상태의 역할을 분리해 stale state를 줄였습니다

---

## 테스트

### Backend
- 정상 예매 흐름 검증
- 홀드 없이 예매 차단
- 잘못된 queue token 차단
- 동일 좌석 동시 hold 경쟁 제어
- 동일 좌석 동시 booking 경쟁 제어

### Frontend
- 공연 목록 렌더링
- 상세 페이지 이동
- queue 상태 확인
- booking 직접 진입 차단
- 버튼 disabled / alert / live region 검증

### Load Test
- k6 기반 부하 테스트 및 경쟁 시나리오 검증
- 동일 좌석 경쟁 시 1건만 성공하고 나머지는 정상 거절되는지 확인

---

## 트러블슈팅

### 1. Redis 홀드만으로는 최종 정합성 보장이 어려웠던 문제
초기에는 Redis 중심으로 좌석 점유를 처리했지만,  
최종 예매 확정까지 Redis만으로 처리하면 동일 좌석 경쟁 상황에서 정합성을 강하게 보장하기 어려웠습니다.

**해결**
- 홀드: Redis `SETNX + TTL`
- 확정: PostgreSQL transaction
- 경쟁 제어: `SELECT ... FOR UPDATE`
- DB 제약조건 추가

**결과**
- 동일 좌석 동시 hold 요청은 1건만 성공
- 동일 좌석 동시 booking 요청도 최종 1건만 성공

---

### 2. 프론트에서 stale state가 쉽게 발생하던 문제
좌석 상태가 빠르게 변하기 때문에, 화면 상태와 서버 상태가 쉽게 어긋났습니다.

**해결**
- optimistic update 최소화
- 액션 후 재조회
- React Query invalidate 적용

**결과**
- 좌석 상태를 서버 기준으로 재동기화하는 구조로 개선

---

### 3. 비동기 UI에서 중복 요청과 잘못된 진입이 발생하기 쉬웠던 문제
대기열 polling, 좌석 선택, 예매 확정 흐름에서 중복 클릭과 잘못된 접근을 제어할 필요가 있었습니다.

**해결**
- 요청 중 버튼 비활성화
- polling 종료 조건 명확화
- booking 직접 진입 차단
- query 기반 상태 관리

**결과**
- 사용자 신뢰를 해치지 않는 예매 흐름 구성

---

### 4. Playwright 테스트가 UI는 맞는데 불안정하게 깨지는 문제
초기에는 role selector가 내부 요소와 충돌해 테스트가 불안정했습니다.

**해결**
- `data-testid` 추가
- stable selector 전략 적용

**결과**
- alert, loading, queue content 검증이 더 안정적으로 동작

---

## 프로젝트에서 얻은 점

- Redis와 PostgreSQL의 역할을 분리해 설명 가능한 구조를 설계했습니다
- 동시성 환경에서 최종 1건만 성공하도록 검증 가능한 예매 흐름을 구현했습니다
- 프론트엔드에서는 stale state, 중복 행동, polling, 접근성을 함께 고려한 UI를 설계했습니다
- Go test, Playwright, k6를 통해 기능/브라우저 흐름/부하 상황을 각각 검증했습니다

---

## 향후 개선 방향

- docker-compose 기반 실행 환경 통합
- CI/CD 자동화
- 결제 흐름 연동
- 예매 완료 E2E 강화
- 좌석 수 증가 시 virtualization 검토
