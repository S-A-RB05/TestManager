package messaging

import (
	"encoding/json"
	"log"

	"github.com/S-A-RB05/TestManager/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Callback func(models.StrategyRequest)

func ConsumeMessage(queue string, callback Callback) {
	conn, err := amqp.Dial("amqps://tnhdeowx:tInXH7wKtKdyn-v97fZ_HGM5XmHsDTNl@rattlesnake.rmq.cloudamqp.com/tnhdeowx")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// Unmarshal the byte slice into a struct
			var strat models.Strategy
			err := json.Unmarshal(d.Body, &strat)
			if err != nil {
				panic(err)
			}

			var mStrat models.StrategyRequest

			mStrat.Name = strat.Name
			mStrat.Ex = strat.Ex
			mStrat.Created = strat.Created
			
			callback(mStrat)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
