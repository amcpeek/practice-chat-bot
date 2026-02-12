// One-off WebSocket test. Run: go run ./cmd/wscheck
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer c.Close()

	// Read first message (history or first response)
	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Fatalf("read: %v", err)
	}
	var out map[string]interface{}
	json.Unmarshal(msg, &out)
	fmt.Printf("Received: %s\n", string(msg))

	// Send a message
	if err := c.WriteMessage(websocket.TextMessage, []byte("test from script")); err != nil {
		log.Fatalf("write: %v", err)
	}

	_, msg, err = c.ReadMessage()
	if err != nil {
		log.Fatalf("read2: %v", err)
	}
	fmt.Printf("After send: %s\n", string(msg))
}
