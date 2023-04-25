package messaging

import (
	"encoding/json"
	"log"

	"github.com/S-A-RB05/TestManager/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			log.Printf("Received a message")
			// Unmarshal the byte slice into a struct
			var strat models.Strategy
			err := json.Unmarshal(d.Body, &strat)
			if err != nil {
				panic(err)
			}
			var bStrat models.StrategyRequest
			id, err := primitive.ObjectIDFromHex(strat.Id)
			if err != nil {
				panic(err)
			}
			log.Printf("iD: " + strat.Id)
			bStrat.Id = id
			bStrat.Name = strat.Name
			bStrat.Ex = strat.Ex
			bStrat.Created = strat.Created
			callback(bStrat)
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
