# Design & Implementation Decisions

Notes for interview discussion. Built to satisfy: user sends message in browser → server responds "Acknowledged" → messages persist across browser close/reopen, in order. Stack: WebSockets, Redis, Docker.

---

## 1. Language & Stack

- **Go** – Single binary, no framework needed for this scope. Standard library `net/http` + two deps: `gorilla/websocket`, `redis/go-redis`. Keeps the solution small and easy to run in Docker.
- **WebSockets** – Requirement; also a good fit for real-time send/receive without polling.
- **Redis** – Requirement; used as a single list to append messages and replay history. No DB schema or migrations.
- **Docker** – Requirement; `docker-compose` runs the app and Redis so reviewers can run it with one command.

---

## 2. Persistence Model

- **Single Redis list** (`chat:messages`). Each element is a JSON object: `{"role":"user"|"server","text":"..."}`.
- **RPUSH** on every user message (we push both the user message and the "Acknowledged" server message so order is preserved).
- **LRANGE 0 -1** on WebSocket connect to get full history and send it to the client.
- **Why a list?** Order is guaranteed (FIFO). No need for timestamps or sorted sets; list order is send/receive order. Minimal commands and minimal code.

---

## 3. Protocol (WebSocket)

- **On connect:** Server sends one message: `{"type":"history","messages":[...]}` with all stored messages. Client renders them in order.
- **On user send:** Client sends plain text. Server stores user + server messages in Redis, then sends `{"type":"message","user":{...},"server":{...}}`. Client appends that pair to the UI.
- **Why JSON?** Simple to parse in browser and in Go; easy to extend with `type` later. No extra schema or ID field needed for this scope.

---

## 4. Frontend

- **Single static HTML file** (`static/index.html`). No build step, no React/Vue. Vanilla JS: open WebSocket, on `history` render all messages, on `message` append one user + one server line. Form submit sends input value and clears the field.
- **WebSocket URL:** `ws://` or `wss://` + `location.host` + `/ws` so it works locally and behind the same host in Docker.

---

## 5. Go Structure

- **One `main.go`** – No extra packages. Handlers: `/` serves the HTML file, `/ws` upgrades to WebSocket. In the WebSocket loop: send history once, then read messages, persist to Redis, write response. Keeps everything in one place for “as little code as possible.”
- **Redis client** – Created at startup, `Ping` to fail fast if Redis is down. No connection pooling config; default is fine for this size.
- **CORS / origin** – `CheckOrigin` returns true so the same HTML can be opened from file or from the server; for a small challenge this avoids origin issues. In production you’d restrict origin.

---

## 6. Docker

- **Dockerfile:** Multi-stage. Stage 1: `golang:1.22-alpine`, copy `go.mod`, download deps, copy source, build static binary. Stage 2: `alpine`, copy binary + `static/`, run on 8080. Small image, no Go toolchain in final image.
- **docker-compose:** Two services. `redis`: official Redis Alpine image, port 6379. `app`: build from Dockerfile, `REDIS_ADDR=redis:6379`, `depends_on: redis`. One `docker compose up` gives a working app and persistence.

---

## 7. What I’d Do Differently at Scale

- **Auth / multi-user:** Add some notion of user/session; store messages per user or per room in Redis (e.g. keys like `chat:user:<id>:messages`).
- **Origin:** Restrict `Upgrader.CheckOrigin` to the real frontend origin.
- **Errors:** Return structured WebSocket close/error frames and surface them in the UI.
- **Redis:** Consider TTL or max length on the list if history can grow unbounded; or move old messages to cold storage.
- **Health:** Add a `/health` or `/ready` that checks Redis and use it in orchestration (e.g. Kubernetes readiness).

---

## 8. How Cursor Helped (Learning Go)

- **Go idioms:** Asked for “minimal Go” and “standard library first”; Cursor suggested `net/http`, `encoding/json`, and small, focused handlers.
- **Libraries:** Picked `gorilla/websocket` and `redis/go-redis` from Cursor’s suggestions; confirmed usage (Upgrader, ReadMessage/WriteJSON, Redis list commands) in-context.
- **Structure:** Kept one file and one HTML to satisfy “as little code as possible” while still making the flow clear for an interview walkthrough.
