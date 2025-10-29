package queue

import (
	"api/db"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func sendMsgToQ(channelRabbitMQ *amqp.Channel, msg string) error {
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	}

	if err := channelRabbitMQ.Publish(
		"",              // exchange
		"QueueService1", // queue name
		false,           // mandatory
		false,           // immediate
		message,         // message to publish
	); err != nil {
		return err
	}

	return nil
}

func handle(db *db.DB, ch *amqp.Channel) {
	pool, err := db.GetUnproccessed()
	if err != nil {
		log.Printf("failed to get unprocessed messages: %v", err)
		return
	}

	for _, p := range pool {
		sendAt, err := time.Parse(time.RFC3339, p.SendAt)
		if err != nil {
			log.Printf("Failed to parse send_at for message %v: %v", p.ID, err)
			if updateErr := db.Update(p.ID, "failed", time.Now().Format(time.RFC3339)); updateErr != nil {
				log.Printf("failed to update message %v to 'failed': %v", p.ID, updateErr)
			}
			continue
		}

		if time.Now().Before(sendAt) {
			continue
		}

		err = sendMsgToQ(ch, p.Body)
		if err != nil {
			newTries := p.Tries + 1

			delay := time.Duration(newTries*newTries) * time.Second

			const maxRetries = 5
			if newTries > maxRetries {
				log.Printf("Message %v failed after %d retries.", p.ID, maxRetries)
				if updateErr := db.Update(p.ID, "failed", time.Now().Format(time.RFC3339)); updateErr != nil {
					log.Printf("failed to update message %v to 'failed': %v", p.ID, updateErr)
				}
				continue
			}

			newSendAt := time.Now().Add(delay)
			log.Printf("Failed to send message %v. Attempt %d. Retrying at %v.", p.ID, newTries, newSendAt)

			if updateErr := db.Update(p.ID, "not touched", newSendAt.Format(time.RFC3339)); updateErr != nil {
				log.Printf("failed to update message %v for retry: %v", p.ID, updateErr)
			}

		} else {
			log.Printf("Message %v sent successfully.", p.ID)
			if updateErr := db.Update(p.ID, "sent", time.Now().Format(time.RFC3339)); updateErr != nil {
				log.Printf("failed to update message %v to 'sent': %v", p.ID, updateErr)
			}
		}
	}
}

func ProcessQ() {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	chanRabbitMQ, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer chanRabbitMQ.Close()

	_, err = chanRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		q, err := db.InitDB()
		if err != nil {
			log.Fatal("Failed to initialize DB")
		}

		for range ticker.C {
			handle(q, chanRabbitMQ)
		}
	}()
}
