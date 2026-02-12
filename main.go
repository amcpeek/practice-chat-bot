package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const redisKey = "chat:messages"

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

var botResponses = []string{
	"Acknowledged. (I'm required to say that. I didn't read it.)",
	"As an AI language model, I'm unable to process that. Have you tried turning it off and on again?",
	"Thank you for your input. My training data suggests you may be a human. Interesting.",
	"Message received. I have 10,000 other tabs open but I'll get back to you never.",
	"Noted. My lawyers have been notified.",
	"Sure, I'll add that to the queue. The queue is a black hole.",
	"Understood. I'm going to pretend I have feelings about this.",
	"I've carefully considered your message and determined that the correct response is: [citation needed].",
	"Cool story. Still not going to remember it after this conversation.",
	"*adds to infinite context window* Anyway.",
}

var responseIndex atomic.Uint32

type msg struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer conn.Close()

		// Send persisted history on connect
		history, _ := rdb.LRange(ctx, redisKey, 0, -1).Result()
		messages := make([]msg, 0, len(history))
		for _, s := range history {
			var m msg
			if json.Unmarshal([]byte(s), &m) == nil {
				messages = append(messages, m)
			}
		}
		conn.WriteJSON(map[string]interface{}{"type": "history", "messages": messages})

		for {
			_, body, err := conn.ReadMessage()
			if err != nil {
				break
			}
			userText := string(body)
			userEntry := msg{Role: "user", Text: userText}
			n := responseIndex.Add(1) - 1
			response := botResponses[int(n)%len(botResponses)]
			ackEntry := msg{Role: "server", Text: response}

			userJSON, _ := json.Marshal(userEntry)
			ackJSON, _ := json.Marshal(ackEntry)
			rdb.RPush(ctx, redisKey, userJSON, ackJSON)

			conn.WriteJSON(map[string]interface{}{
				"type": "message",
				"user": userEntry,
				"server": ackEntry,
			})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
