package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/segmentio/kafka-go"
)

func sendMessageFromFile(conn *kafka.Conn, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = conn.WriteMessages(
		kafka.Message{Value: data},
	)

	return err
}

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9093"
	}
	topic := "orders"

	// START KAFKA CONNECTION
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	if err != nil {
		log.Fatal("🔴 Failed to dial leader:", err)
	}
	fmt.Println("🟢 Connection opened")

	// SEND HEALTHCHECK
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = sendMessageFromFile(conn, "healthcheck.json")
	if err != nil {
		log.Fatal("🔴 Failed to send healthcheck:", err)
	}
	fmt.Println("🟢 Healthcheck sent")

	// SEND ALL ORDERS from orders/ directory
	entries, err := os.ReadDir("orders")
	if err != nil {
		log.Fatal("🔴 Failed to read orders directory:", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filePath := filepath.Join("orders", e.Name())
		if err := sendMessageFromFile(conn, filePath); err != nil {
			log.Printf("🔴 Failed to send %s: %v", filePath, err)
			continue
		}
		fmt.Printf("🟢 Sent %s\n", filePath)
		time.Sleep(300 * time.Millisecond)
	}

	// CLOSE CONNECTION
	if err := conn.Close(); err != nil {
		log.Fatal("🔴 Failed to close writer:", err)
	}
	fmt.Println("🟢 Connection closed")
}
