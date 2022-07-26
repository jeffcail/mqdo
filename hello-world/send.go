package main

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare q queue")

	body := "hello world"
	err = ch.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish e message")

	log.Printf(" [x] Sent %s\n", body)
}
