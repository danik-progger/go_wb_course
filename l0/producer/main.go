package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/segmentio/kafka-go"
)

func generateFakeOrder() OrderMessage {
	var items []Item
	trackNumber, err := gofakeit.Generate("{hex:10}")
	if err != nil {
		log.Fatal("🔴 Failed to generate track number")
	}
	// Generate between 1 and 4 items
	numItems := rand.Intn(4) + 1

	for range numItems {
		sale := gofakeit.Number(0, 70)
		price := gofakeit.Number(1000, 20000)
		totalPrice := price * (100 - sale) / 100
		item := Item{
			ChrtID:      gofakeit.Number(1000000, 9999999),
			TrackNumber: trackNumber,
			Price:       price,
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        sale,
			Size:        "0",
			TotalPrice:  totalPrice,
			NmID:        gofakeit.Number(1000000, 9999999),
			Brand:       gofakeit.Company(),
			Status:      202,
		}
		items = append(items, item)
	}

	delivery := Delivery{
		Name:    gofakeit.Name(),
		Phone:   gofakeit.Phone(),
		Zip:     gofakeit.Zip(),
		City:    gofakeit.City(),
		Address: gofakeit.Street(),
		Region:  gofakeit.State(),
		Email:   gofakeit.Email(),
	}

	var goodsTotal int
	for _, item := range items {
		goodsTotal += item.TotalPrice
	}
	deliveryCost := gofakeit.Number(500, 2000)

	payment := Payment{
		Transaction:  gofakeit.UUID(),
		RequestID:    "",
		Currency:     gofakeit.CurrencyShort(),
		Provider:     "wbpay",
		Amount:       goodsTotal + deliveryCost,
		PaymentDt:    time.Now().Unix(),
		Bank:         gofakeit.Company(),
		DeliveryCost: deliveryCost,
		GoodsTotal:   goodsTotal,
		CustomFee:    0,
	}

	order := Order{
		OrderUID:        gofakeit.UUID(),
		TrackNumber:     trackNumber,
		Entry:           "WBIL",
		Delivery:        delivery,
		Payment:         payment,
		Items:           items,
		Locale:          "en",
		CustomerID:      gofakeit.Username(),
		DeliveryService: "meest",
		Shardkey:        fmt.Sprintf("%d", gofakeit.Number(1, 10)),
		SmID:            gofakeit.Number(1, 100),
		DateCreated:     time.Now(),
		OofShard:        "1",
	}

	return OrderMessage{
		Action: "create_order",
		Body:   order,
	}
}

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9093"
	}
	topic := "orders"

	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	if err != nil {
		log.Fatal("🔴 Failed to dial leader:", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("🔴 Failed to close Kafka connection:", err)
		}
	}()
	fmt.Println("🟢 Kafka Connection opened")

	gofakeit.Seed(0)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("✨ Starting dynamic order generation...")
	for {
		select {
		case <-quit:
			log.Println("🟢 Shutting down producer...")
			return
		default:
			orderMsg := generateFakeOrder()

			msgBytes, err := json.Marshal(orderMsg)
			if err != nil {
				log.Printf("🔴 Failed to marshal JSON: %v", err)
				continue
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, err = conn.WriteMessages(kafka.Message{Value: msgBytes})
			if err != nil {
				log.Printf("🔴 Failed to write message: %v", err)
			} else {
				fmt.Printf("🟢 Sent order %s", orderMsg.Body.OrderUID)
			}

			// Wait for a random duration between 1 and 5 seconds
			time.Sleep(time.Duration(rand.Intn(4)+1) * time.Second)
		}
	}
}
