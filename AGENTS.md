# CHATBOT – Agent context

## Stack
- **Backend:** Go 1.22, `net/http`, `gorilla/websocket`, `redis/go-redis/v9`
- **Frontend:** Single static HTML/JS in `static/index.html`
- **Infra:** Docker + docker-compose (app + Redis)

## Layout
- `main.go` – HTTP server, WebSocket handler, Redis read/write; serves `/` (HTML) and `/ws` (WebSocket)
- `static/index.html` – UI and WebSocket client
- `Dockerfile` – multi-stage Go build, serve on 8080
- `docker-compose.yml` – `app` (this service) + `redis`

## Run
- **Docker:** `docker compose up --build` then open http://localhost:8080
- **Local (need Redis):** `REDIS_ADDR=localhost:6379 go run .` then open http://localhost:8080

## Conventions
- Prefer standard library; add deps only when needed.
- Keep handlers in `main.go` unless the file grows large; then consider `internal/` packages.
- Redis key for message list: `chat:messages` (list of JSON `{"role","text"}`).

## Decisions
- See `DECISIONS.md` for design choices and interview notes.
