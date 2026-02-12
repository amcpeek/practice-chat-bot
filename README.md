# Chat Bot (Challenge)

Minimal fullstack chatbot: send a message in the browser, server responds "Acknowledged". Messages persist in Redis and are shown in order when you reopen the browser.

## Requirements met
- User sends text in a web browser
- Server receives and responds with "Acknowledged"
- After closing the browser, reopening shows messages in send/receive order
- Uses WebSockets, Redis, and Docker

## Run with Docker
```bash
docker compose up --build
```
Open http://localhost:8080. Send messages; close the tab and reopen to see history.

To stop: `docker compose down`.

## Run locally (Redis required)
```bash
# Terminal 1: Redis
docker run -p 6379:6379 redis:7-alpine

# Terminal 2: App
REDIS_ADDR=localhost:6379 go run .
```
Open http://localhost:8080.

## Tech
- **Go** – HTTP + WebSocket server, Redis client
- **Redis** – Single list for message history
- **Static HTML/JS** – WebSocket client, no build step

## Optional: CLI WebSocket test
With the app running (Docker or local), from the project root:
```bash
go run ./cmd/wscheck
```
Connects to `/ws`, sends one message, and prints the history and response.

---

Design and implementation decisions are in [DECISIONS.md](DECISIONS.md).
