package messaging

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ProduceMessage(message string, queue string) {
	err := godotenv.Load()
	FailOnError(err, "Error loading .env file")

	conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		queue, // queue name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := []byte(message)

	// Publish the message to the queue
	err = ch.PublishWithContext(ctx,
		"",     // exchange name
		q.Name, // routing key
		false,  // mandatory flag
		false,  // immediate flag
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
