package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9093"
	}

	topic := "orders"

	// START CONNECTION
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	defer conn.Close()
	if err != nil {
		log.Fatal("🔴 Failed to dial leader:", err)
	}

	// START READING
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	defer batch.Close()

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(b[:n], &payload); err != nil {
			log.Println("🔴 Invalid message:", err)
			continue
		}

		if _, ok := payload["action"]; !ok {
			log.Println("🔴 Invalid message:", err)
			continue
		}
		fmt.Println("🟢 Parsed message, action:", payload["action"])

		if payload["action"] == "healthcheck" {
			fmt.Printf("🟢 Healthcheck successful: %s\n", payload["body"])
		}
		if payload["action"] == "create_order" {
			// forward original message bytes as JSON POST to back service
			client := &http.Client{Timeout: 5 * time.Second}
			req, err := http.NewRequest("POST", "http://back:8080/orders", bytes.NewReader(b[:n]))
			if err != nil {
				log.Fatal("🔴 Failed to build request:", err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("🔴 Failed to send request:", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				fmt.Println("🟢 Successful request to server")
			} else {
				log.Printf("🔴 Server returned status: %s", resp.Status)
			}
		}

	}
}
