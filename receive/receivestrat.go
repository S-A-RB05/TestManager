package receive

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Strat() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"strat_queue", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare queue")

	err = ch.QueueBind(
		q.Name,           // queue name
		"mykey",          // routing key
		"strat_exchange", // exchange name
		false,
		nil,
	)
	failOnError(err, "Failed to bind to queue")

	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer name
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to consume message")

	for msg := range msgs {
		err = ioutil.WriteFile("myfile.txt", msg.Body, 0644)
		failOnError(err, "Failed to write file")

		log.Printf("File saved: %s", msg.Body)
	}

}
