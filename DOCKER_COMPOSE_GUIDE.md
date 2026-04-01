# Docker Compose Guide

## Included services
- postgres
- redis
- backend
- frontend

## Why this setup
- postgres / redis are infrastructure dependencies
- backend waits for postgres / redis healthchecks
- frontend talks to:
  - browser 기준: `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api`
  - 서버 렌더링 기준: `INTERNAL_API_BASE_URL=http://backend:8080/api`

즉 Next.js가 서버 안에서 fetch할 때는 `backend` 서비스 이름으로 붙고,
브라우저에서는 `localhost:8080`으로 붙도록 분리했습니다.

## Run
프로젝트 루트에서:

```bash
docker compose up --build
```

백그라운드:

```bash
docker compose up -d --build
```

## Access
- frontend: http://localhost:3000
- backend: http://localhost:8080
- postgres: localhost:5432
- redis: localhost:6379

## Notes
- `backend/seed.sql`은 postgres 최초 초기화 시 자동 실행됩니다.
- postgres volume이 이미 있으면 seed가 다시 실행되지 않습니다.
- 초기화부터 다시 하고 싶으면:

```bash
docker compose down -v
docker compose up --build
```
