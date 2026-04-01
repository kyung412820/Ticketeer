# Frontend Rendering Strategy

## Goal
TICKETEER 프론트는 페이지 성격에 따라 **조회 중심 화면**과 **실시간 상호작용 화면**을 분리해서 렌더링 전략을 잡았습니다.

---

## Page Strategy

### `/events`
**전략:** 서버 중심 렌더링

**이유**
- 공연 목록은 초기 진입 시 바로 보여주는 것이 중요함
- SEO와 첫 화면 응답이 상대적으로 중요함
- 사용자별 실시간 상호작용보다 조회 성격이 강함

**현재 구현**
- 서버 fetch 기반
- `force-dynamic`으로 최신 데이터 반영

---

### `/events/[id]`
**전략:** 서버 중심 렌더링

**이유**
- 공연 상세 역시 조회 중심 화면
- 예매 오픈 시간, 공연 정보 같은 기본 정보는 서버 기준으로 보여주는 것이 자연스러움
- 사용자가 이 페이지에서 바로 대기열 입장 여부를 판단함

**현재 구현**
- 서버 fetch 기반
- 존재하지 않는 공연은 `notFound()` 처리

---

### `/queue/[eventId]`
**전략:** 클라이언트 중심 렌더링

**이유**
- queue token 발급, polling, 상태 변화가 모두 브라우저 상호작용 중심
- `WAITING -> READY -> EXPIRED` 상태 전환을 주기적으로 갱신해야 함
- local/session storage와 연결된 사용자 흐름 데이터 사용

**현재 구현**
- client component
- polling 기반 상태 갱신
- queue token / client id 활용

---

### `/booking/[eventId]`
**전략:** 클라이언트 중심 렌더링

**이유**
- 좌석 선택, 홀드, 예매 확정은 모두 사용자의 즉시 행동과 연결됨
- 요청 중 버튼 상태, 에러 메시지, 재동기화가 중요함
- queue token, client id, 선택 좌석 상태를 브라우저에서 다룸

**현재 구현**
- client component
- 좌석 재조회
- 요청 중 상태 제어
- 결과 UI 반영

---

## Summary

### 서버 중심 페이지
- `/events`
- `/events/[id]`

### 클라이언트 중심 페이지
- `/queue/[eventId]`
- `/booking/[eventId]`

이렇게 분리한 이유는,
**조회 중심 화면은 서버 기준으로 안정적으로 보여주고, 실시간 상호작용 화면은 클라이언트 상태 중심으로 제어하기 위해서**입니다.
