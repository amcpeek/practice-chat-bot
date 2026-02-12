package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const redisKey = "chat:messages"

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

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
			ackEntry := msg{Role: "server", Text: "Acknowledged"}

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
