package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

var validate = validator.New()

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9093"
	}

	topic := "orders"
	groupID := "order-consumers"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 1e6,  // 1MB
	})
	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	fmt.Println("🟢 Kafka reader started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-quit
		log.Println("🟢 Shutting down consumer...")
		cancel()
	}()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				break
			}
			log.Println("🔴 Error reading message:", err)
			break
		}

		var msg OrderMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Println("🔴 Invalid message:", err)
			continue
		}

		if msg.Action == "create_order" {
			var orderBody Order
			bodyBytes, err := json.Marshal(msg.Body)
			if err != nil {
				log.Println("🔴 Failed to marshal message body:", err)
				continue
			}
			if err := json.Unmarshal(bodyBytes, &orderBody); err != nil {
				log.Println("🔴 Invalid order message body:", err)
				continue
			}
			msg.Body = orderBody

			if err := validate.Struct(msg); err != nil {
				log.Println("🔴 Invalid message content:", err)
				continue
			}
		}

		fmt.Printf("🟢 Parsed message, action: %s, offset: %d\n", msg.Action, m.Offset)

		if msg.Action == "healthcheck" {
			fmt.Printf("🟢 Healthcheck successful: %#v\n", msg.Body)
		}
		if msg.Action == "create_order" {
			client := &http.Client{Timeout: 5 * time.Second}
			req, err := http.NewRequest("POST", "http://back:8080/orders", bytes.NewReader(m.Value))
			if err != nil {
				log.Println("🔴 Failed to build request:", err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Println("🔴 Failed to send request:", err)
				continue
			}

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				fmt.Println("🟢 Successful request to server")
				if err := r.CommitMessages(context.Background(), m); err != nil {
					log.Println("🔴 Failed to commit messages:", err)
				}
			} else {
				log.Printf("🔴 Server returned status: %s", resp.Status)
			}

			if err := resp.Body.Close(); err != nil {
				log.Println("🔴 Failed to close response body:", err)
			}
		}
	}
}
